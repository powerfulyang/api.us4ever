package task

import (
	"api.us4ever/internal/database"
	"api.us4ever/internal/task/image"
	"api.us4ever/internal/task/keep"
	"api.us4ever/internal/task/telegram"
)

// RegisterTasks 注册所有定时任务
func RegisterTasks(scheduler *Scheduler, getDB func() database.Service) error {
	// 每 1 分钟执行一次生成 title 和 summary 的任务
	err := scheduler.AddTaskWithDB("generate_keep_title_summary", "0 * * * * *", keep.GenerateTitleAndSummary, getDB)
	if err != nil {
		return err
	}

	// 每个整点执行一次 TriggerSyncTelegram
	err = scheduler.AddTask("trigger_sync_telegram", "0 0 * * * *", telegram.TriggerSyncTelegram)
	if err != nil {
		return err
	}

	//Add the image OCR task (runs every 5 seconds)
	err = scheduler.AddTaskWithDB("process_image_ocr", "*/5 * * * * *", image.ProcessImageOCR, getDB)
	if err != nil {
		return err
	}

	return nil
}
