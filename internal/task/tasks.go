package task

import "log"

// ExampleTask 示例任务
func ExampleTask() {
	log.Println("执行示例任务...")
}

// RegisterTasks 注册所有定时任务
func RegisterTasks(scheduler *Scheduler) error {
	// 示例：每分钟执行一次的任务
	err := scheduler.AddTask("example_task", "0 * * * * *", ExampleTask)
	if err != nil {
		return err
	}
	return nil
}
