package image

import (
	"api.us4ever/internal/config"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/ent"
	"api.us4ever/internal/ent/image"
)

const (
	taskLimit = 1 // Process one image per task run
)

type OCRResponse struct {
	Result struct {
		Errcode     int    `json:"errcode"`
		Height      int    `json:"height"`
		Width       int    `json:"width"`
		Imgpath     string `json:"imgpath"`
		OCRResponse []struct {
			Text   string  `json:"text"`
			Left   float64 `json:"left"`
			Top    float64 `json:"top"`
			Right  float64 `json:"right"`
			Bottom float64 `json:"bottom"`
			Rate   float64 `json:"rate"`
		} `json:"ocr_response"`
	} `json:"result"`
}

type OCRRequest struct {
	Image string `json:"image"` // Base64 encoded image data
}

type ExtraData struct {
	OCRResponse []struct {
		Text   string  `json:"text"`
		Left   float64 `json:"left"`
		Top    float64 `json:"top"`
		Right  float64 `json:"right"`
		Bottom float64 `json:"bottom"`
		Rate   float64 `json:"rate"`
	} `json:"ocr_response,omitempty"`
	// Include other potential fields from ExtraData if known
}

// ProcessImageOCR finds images needing OCR, processes them, and updates the database.
func ProcessImageOCR(db database.Service) {
	log.Println("Starting ProcessImageOCR task...")
	ctx := context.Background()

	// Find images that have an original file ID.
	// We will filter based on ExtraData content after fetching.
	imagesToCheck, err := db.Client().Image.Query().
		Where(
			image.OriginalIDNEQ(""),
			image.DescriptionEQ(""),
			image.SizeLTE(1_000_000),
		).
		// 统一放进一个 Selector，用 OR 连接
		Where(func(s *sql.Selector) {
			s.Where(
				sqljson.ValueEQ(
					image.FieldExtraData,
					json.RawMessage(`{}`),
				),
			)
		}).
		WithOriginal(func(q *ent.FileQuery) { // Eager load the original file
			q.WithBucket() // Eager load the bucket associated with the file
		}).
		Order(ent.Desc(image.FieldUpdatedAt)). // Process recently updated first, or choose another order
		Limit(taskLimit).                      // Fetch a slightly larger batch to filter
		All(ctx)

	if err != nil {
		log.Printf("Error querying images for OCR check: %v", err)
		return
	}

	if len(imagesToCheck) == 0 {
		log.Println("No images with original files found.")
		return
	}

	// 遍历检查的图片
	for _, img := range imagesToCheck {
		if !needsOCRProcessing(img) {
			log.Println("Image ID:", img.ID, "does not need OCR processing.")
			continue
		}

		log.Printf("Processing OCR for image ID: %s", img.ID)
		err := ProcessSingleImageOCR(ctx, db, img.ID)
		if err != nil {
			log.Printf("Failed to process image %s: %v", img.ID, err)
			continue
		}
	}

	log.Println("ProcessImageOCR task finished.")
}

// needsOCRProcessing checks if the image's ExtraData indicates OCR is needed.
func needsOCRProcessing(img *ent.Image) bool {
	var currentExtraData map[string]interface{}
	if err := json.Unmarshal(img.ExtraData, &currentExtraData); err != nil {
		log.Printf("Image ID %s: Failed to unmarshal existing ExtraData for checking: %v. Assuming needs processing.", img.ID, err)
		return true // Error unmarshalling, assume it needs processing to be safe
	}

	_, hasKey := currentExtraData["ocr_response"]
	return !hasKey // Needs processing if the key doesn't exist
}

func downloadImage(url string) ([]byte, error) {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	// 读取响应体
	return io.ReadAll(resp.Body)
}

func callOCRAPI(base64Image string) (*OCRResponse, error) {
	requestPayload := OCRRequest{Image: base64Image}
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OCR request: %w", err)
	}

	// 从配置中获取 endpoint 和 apiKey
	appConfig := config.GetAppConfig()
	endPoint := appConfig.OCR.Endpoint
	if endPoint == "" {
		return nil, fmt.Errorf("failed to get OCR endpoint from config")
	}

	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create OCR API request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second} // Add a timeout
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to OCR API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close OCR API response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OCR API returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var ocrResponse OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResponse); err != nil {
		return nil, fmt.Errorf("failed to decode OCR API response: %w", err)
	}

	return &ocrResponse, nil
}

func cleanDescription(description string) string {
	reNewline := regexp.MustCompile(`[\r\n]+`)
	description = reNewline.ReplaceAllString(description, "")

	reSpaces := regexp.MustCompile(`\s+`)
	description = reSpaces.ReplaceAllString(description, " ")

	// 3. 去除首尾空格
	return strings.TrimSpace(description)
}

func updateImageExtraData(ctx context.Context, db database.Service, img *ent.Image, ocrData []struct {
	Text   string  `json:"text"`
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
	Rate   float64 `json:"rate"`
}) error {
	var currentExtraData map[string]interface{}

	// Check if ExtraData is nil or empty before trying to unmarshal
	if len(img.ExtraData) > 0 {
		if err := json.Unmarshal(img.ExtraData, &currentExtraData); err != nil {
			// If unmarshalling fails, log it but proceed with a new map
			log.Printf("Image ID %s: Failed to unmarshal existing ExtraData '%s': %v. Overwriting with new data.", img.ID, string(img.ExtraData), err)
			currentExtraData = make(map[string]interface{}) // Initialize if unmarshal fails
		}
	} else {
		currentExtraData = make(map[string]interface{}) // Initialize if nil or empty
	}

	// Add or update the ocr_response field
	currentExtraData["ocr_response"] = ocrData

	// 合并所有OCR文本到description
	var textParts []string
	for _, item := range ocrData {
		if item.Text != "" && item.Rate >= 0.7 { // 只使用置信度大于70%的文本
			textParts = append(textParts, item.Text)
		}
	}
	description := strings.Join(textParts, " ")
	// 去除\r\n和多余空格
	description = cleanDescription(description)
	log.Println("Combined OCR text for image ID:", img.ID, ":", description)

	// Marshal the updated map back to json.RawMessage
	updatedExtraDataBytes, err := json.Marshal(currentExtraData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated ExtraData: %w", err)
	}

	// Save the updated image record with both ExtraData and Description
	_, err = db.Client().Image.UpdateOne(img).
		SetExtraData(updatedExtraDataBytes).
		SetDescription(description).
		SetUpdatedAt(time.Now()). // Update the timestamp
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to save updated image data: %w", err)
	}

	return nil
}

// ProcessSingleImageOCR processes OCR for a specific image ID
func ProcessSingleImageOCR(ctx context.Context, db database.Service, imageID string) error {
	img, err := db.Client().Image.Query().
		Where(image.ID(imageID)).
		WithOriginal(func(q *ent.FileQuery) {
			q.WithBucket()
		}).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("image not found: %s", imageID)
		}
		return fmt.Errorf("failed to query image: %v", err)
	}

	if img.Edges.Original == nil || img.Edges.Original.Edges.Bucket == nil {
		return fmt.Errorf("image %s: missing original file or bucket information", imageID)
	}

	originalFile := img.Edges.Original
	bucketInfo := originalFile.Edges.Bucket

	if bucketInfo.PublicUrl == "" || originalFile.Path == "" {
		return fmt.Errorf("image %s: missing PublicUrl or Path for original file %s", imageID, originalFile.ID)
	}

	imageURL := fmt.Sprintf("%s/%s", bucketInfo.PublicUrl, originalFile.Path)

	// Download image
	imageData, err := downloadImage(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}

	// Detect content type and create base64 image with proper MIME type
	contentType := http.DetectContentType(imageData)
	base64Image := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(imageData)

	// Call OCR API
	ocrResult, err := callOCRAPI(base64Image)
	if err != nil {
		return fmt.Errorf("failed to call OCR API: %v", err)
	}

	if ocrResult.Result.Errcode != 0 {
		log.Println("OCR API returned error code:", ocrResult.Result.Errcode)
	}

	// Update Image ExtraData
	err = updateImageExtraData(ctx, db, img, ocrResult.Result.OCRResponse)
	if err != nil {
		return fmt.Errorf("failed to update ExtraData: %v", err)
	}

	return nil
}
