package task

import (
	"api.us4ever/internal/database"
	"api.us4ever/internal/task/keep"
)

// RegisterTasks 注册所有定时任务
func RegisterTasks(scheduler *Scheduler, getDB func() database.Service) error {
	// 每 1 分钟执行一次生成 title 和 summary 的任务
	err := scheduler.AddTaskWithDB("generate_keep_title_summary", "0 * * * * *", keep.GenerateTitleAndSummary, getDB)
	if err != nil {
		return err
	}

	return nil
}
