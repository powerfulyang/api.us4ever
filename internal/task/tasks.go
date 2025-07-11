package task

import (
	"api.us4ever/internal/server"
	"api.us4ever/internal/task/image"
	"api.us4ever/internal/task/keep"
	"api.us4ever/internal/task/telegram"
	"api.us4ever/internal/task/vector"
)

// RegisterTasks 注册所有定时任务
func RegisterTasks(scheduler *Scheduler, fiberServer *server.FiberServer) error {

	// 每 1 分钟执行一次生成 title 和 summary 的任务
	err := scheduler.AddTaskWithServer("generate_keep_title_summary", "0 * * * * *", keep.GenerateTitleAndSummary, fiberServer)
	if err != nil {
		return err
	}

	// 每 60s 执行一次 TriggerSyncTelegram
	err = scheduler.AddTask("trigger_sync_telegram", "0 * * * * *", telegram.TriggerSyncTelegram)
	if err != nil {
		return err
	}

	// Add the image OCR task (runs every 5 seconds)
	err = scheduler.AddTaskWithServer("process_image_ocr", "*/5 * * * * *", image.ProcessImageOCR, fiberServer)
	if err != nil {
		return err
	}

	// embedding moments task (runs every 60 seconds)
	err = scheduler.AddTaskWithServer("embedding_moments", "0 * * * * *", vector.EmbeddingMoments, fiberServer)
	if err != nil {
		return err
	}

	// embedding keeps task (runs every 60 seconds)
	err = scheduler.AddTaskWithServer("embedding_keeps", "0 * * * * *", vector.EmbeddingKeeps, fiberServer)
	if err != nil {
		return err
	}

	return nil
}
