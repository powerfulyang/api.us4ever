package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/ent"
	"api.us4ever/internal/ent/image"
)

const (
	ocrAPIURL = "http://localhost:5000/ocr"
	taskLimit = 1 // Process one image per task run
)

type OCRResponse struct {
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
}

type OCRRequest struct {
	Image string `json:"image"` // Base64 encoded image data
}

type ImageExtraData struct {
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
		Where(image.OriginalIDNEQ("")).       // Ensure there is an original file reference
		WithOriginal(func(q *ent.FileQuery) { // Eager load the original file
			q.WithBucket() // Eager load the bucket associated with the file
		}).
		Order(ent.Desc(image.FieldUpdatedAt)). // Process recently updated first, or choose another order
		Limit(taskLimit * 5).                  // Fetch a slightly larger batch to filter
		All(ctx)

	if err != nil {
		log.Printf("Error querying images for OCR check: %v", err)
		return
	}

	if len(imagesToCheck) == 0 {
		// log.Println("No images with original files found.")
		return
	}

	var imagesToProcess []*ent.Image
	for _, img := range imagesToCheck {
		if needsOCRProcessing(img) {
			imagesToProcess = append(imagesToProcess, img)
			if len(imagesToProcess) >= taskLimit {
				break // Limit the number of images processed per run
			}
		}
	}

	if len(imagesToProcess) == 0 {
		// log.Println("No images found requiring OCR processing in the checked batch.")
		return
	}

	log.Printf("Found %d image(s) to process for OCR.", len(imagesToProcess))

	for _, img := range imagesToProcess {
		if img.Edges.Original == nil || img.Edges.Original.Edges.Bucket == nil {
			log.Printf("Image ID %s: Missing original file or bucket information, skipping.", img.ID)
			continue
		}

		originalFile := img.Edges.Original
		bucketInfo := originalFile.Edges.Bucket

		if bucketInfo.PublicUrl == "" || originalFile.Path == "" {
			log.Printf("Image ID %s: Missing PublicUrl or Path for original file %s, skipping.", img.ID, originalFile.ID)
			continue
		}

		imageURL := fmt.Sprintf("%s/%s", bucketInfo.PublicUrl, originalFile.Path)
		log.Printf("Image ID %s: Processing image from URL: %s", img.ID, imageURL)

		// Download image
		imageData, err := downloadImage(imageURL)
		if err != nil {
			log.Printf("Image ID %s: Failed to download image: %v", img.ID, err)
			continue
		}

		// Base64 encode
		base64Image := base64.StdEncoding.EncodeToString(imageData)

		// Call OCR API
		ocrResult, err := callOCRAPI(base64Image)
		if err != nil {
			log.Printf("Image ID %s: Failed to call OCR API: %v", img.ID, err)
			continue
		}

		if ocrResult.Errcode != 0 {
			log.Printf("Image ID %s: OCR API returned error code %d", img.ID, ocrResult.Errcode)
			continue // Or handle specific error codes if needed
		}

		// Update Image ExtraData
		err = updateImageExtraData(ctx, db, img, ocrResult.OCRResponse)
		if err != nil {
			log.Printf("Image ID %s: Failed to update ExtraData: %v", img.ID, err)
			continue
		}

		log.Printf("Image ID %s: Successfully processed OCR and updated ExtraData.", img.ID)
	}

	log.Println("ProcessImageOCR task finished.")
}

// needsOCRProcessing checks if the image's ExtraData indicates OCR is needed.
func needsOCRProcessing(img *ent.Image) bool {
	if img.ExtraData == nil || len(img.ExtraData) == 0 || string(img.ExtraData) == "null" {
		return true // ExtraData is nil or empty, needs processing
	}

	var currentExtraData map[string]interface{}
	if err := json.Unmarshal(img.ExtraData, &currentExtraData); err != nil {
		log.Printf("Image ID %s: Failed to unmarshal existing ExtraData for checking: %v. Assuming needs processing.", img.ID, err)
		return true // Error unmarshalling, assume it needs processing to be safe
	}

	_, hasKey := currentExtraData["ocr_response"]
	return !hasKey // Needs processing if the key doesn't exist
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}
	return imageData, nil
}

func callOCRAPI(base64Image string) (*OCRResponse, error) {
	requestPayload := OCRRequest{Image: base64Image}
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OCR request: %w", err)
	}

	req, err := http.NewRequest("POST", ocrAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create OCR API request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second} // Add a timeout
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to OCR API: %w", err)
	}
	defer resp.Body.Close()

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
	if img.ExtraData != nil && len(img.ExtraData) > 0 && string(img.ExtraData) != "null" {
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

	// Marshal the updated map back to json.RawMessage
	updatedExtraDataBytes, err := json.Marshal(currentExtraData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated ExtraData: %w", err)
	}

	// Save the updated image record
	_, err = db.Client().Image.UpdateOne(img).
		SetExtraData(updatedExtraDataBytes).
		SetUpdatedAt(time.Now()). // Update the timestamp
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to save updated image ExtraData: %w", err)
	}

	return nil
}

// Helper predicate for checking if JSON field has a key (requires DB specific functions or raw SQL)
// ent doesn't directly support HasKey for JSONB out of the box for all drivers in a portable way.
// We use a placeholder here. The query above uses a simplified check.
// For PostgreSQL, you might use raw SQL like: `WHERE "extraData" ? 'ocr_response'`
// For simplicity in this example, we query for NULL or check after fetching.
// The current implementation fetches records where ExtraData IS NULL OR OriginalID is not empty,
// then checks the ExtraData content after fetching, which might be less efficient.
// A more robust solution might involve adding a dedicated 'ocr_status' field.
//
// The query `image.Not(image.ExtraDataHasKey("ocr_response"))` relies on ent generated code which might not exist
// or work as expected depending on ent version and DB driver features for JSON.
// Let's adjust the query logic slightly to be safer. We'll query for images
// where ExtraData is NULL *or* where ExtraData does *not* contain the key.
// Direct key checking in WHERE might need custom predicates or raw SQL depending on the DB.
// The provided `image.ExtraDataHasKey` predicate assumes ent generated it based on schema/db features.
// Let's assume it works for now.
