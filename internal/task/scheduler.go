package task

import (
	"api.us4ever/internal/server"
	"log"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron  *cron.Cron
	pool  *ants.Pool
	tasks map[string]cron.EntryID
	locks map[string]*sync.Mutex
}

// NewScheduler 创建新的调度器
func NewScheduler() (*Scheduler, error) {
	pool, err := ants.NewPool(10)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		cron:  cron.New(cron.WithSeconds()),
		pool:  pool,
		tasks: make(map[string]cron.EntryID),
		locks: make(map[string]*sync.Mutex),
	}, nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.pool.Release()
}

type FuncWithServer func(fiberServer *server.FiberServer) (int, error)
type FuncWithoutServer func() (int, error)

// AddTask 添加定时任务
func (s *Scheduler) AddTask(name, spec string, task FuncWithoutServer) error {
	// 为每个任务创建一个锁
	s.locks[name] = &sync.Mutex{}

	id, err := s.cron.AddFunc(spec, func() {

		// 提交任务到协程池
		err := s.pool.Submit(func() {
			// 尝试获取锁，如果任务正在执行则跳过本次执行
			if !s.locks[name].TryLock() {
				log.Printf("任务 %s 正在执行中，跳过本次执行", name)
				return
			}
			defer s.locks[name].Unlock()
			startTime := time.Now()
			count, err := task()
			if err != nil {
				log.Printf("任务 %s 执行出错: %v", name, err)
			} else if count > 0 {
				log.Printf("任务 %s 执行完成，耗时: %v，共处理: %v 条数据", name, time.Since(startTime), count)
			}
		})
		if err != nil {
			log.Printf("提交任务 %s 失败: %v", name, err)
		}
	})
	if err != nil {
		return err
	}
	s.tasks[name] = id
	return nil
}

// AddTaskWithServer 添加需要数据库连接的定时任务
func (s *Scheduler) AddTaskWithServer(name, spec string, task FuncWithServer, fiberServer *server.FiberServer) error {
	// 为每个任务创建一个锁
	s.locks[name] = &sync.Mutex{}

	id, err := s.cron.AddFunc(spec, func() {
		// 提交任务到协程池
		err := s.pool.Submit(func() {
			// 尝试获取锁，如果任务正在执行则跳过本次执行
			if !s.locks[name].TryLock() {
				log.Printf("任务 %s 正在执行中，跳过本次执行", name)
				return
			}
			defer s.locks[name].Unlock()
			startTime := time.Now()
			count, err := task(fiberServer)
			if err != nil {
				log.Printf("任务 %s 执行出错: %v", name, err)
			} else if count > 0 {
				log.Printf("任务 %s 执行完成，耗时: %v，共处理: %v 条数据", name, time.Since(startTime), count)
			}
		})
		if err != nil {
			log.Printf("提交任务 %s 失败: %v", name, err)
		}
	})
	if err != nil {
		return err
	}
	s.tasks[name] = id
	return nil
}

// RemoveTask 移除定时任务
func (s *Scheduler) RemoveTask(name string) {
	if id, exists := s.tasks[name]; exists {
		s.cron.Remove(id)
		delete(s.tasks, name)
		delete(s.locks, name)
	}
}
