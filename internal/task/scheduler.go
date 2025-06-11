package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"api.us4ever/internal/logger"
	"api.us4ever/internal/server"
	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron/v3"
)

// Scheduler represents a task scheduler with cron jobs and worker pool
type Scheduler struct {
	cron    *cron.Cron
	pool    *ants.Pool
	tasks   map[string]cron.EntryID
	locks   map[string]*sync.Mutex
	logger  *logger.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	mu      sync.RWMutex
}

// NewScheduler creates a new task scheduler with improved error handling and logging
func NewScheduler() (*Scheduler, error) {
	const defaultPoolSize = 10

	pool, err := ants.NewPool(defaultPoolSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create worker pool: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	taskLogger, err := logger.New("scheduler")
	if err != nil {
		pool.Release()
		cancel()
		return nil, fmt.Errorf("failed to create task logger: %w", err)
	}

	scheduler := &Scheduler{
		cron:    cron.New(cron.WithSeconds()),
		pool:    pool,
		tasks:   make(map[string]cron.EntryID),
		locks:   make(map[string]*sync.Mutex),
		logger:  taskLogger,
		ctx:     ctx,
		cancel:  cancel,
		running: false,
	}

	taskLogger.Info("task scheduler created successfully", logger.Fields{
		"pool_size": defaultPoolSize,
	})

	return scheduler, nil
}

// Start starts the task scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		s.logger.Warn("scheduler is already running")
		return
	}

	s.cron.Start()
	s.running = true
	s.logger.Info("task scheduler started")
}

// Stop stops the task scheduler gracefully
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		s.logger.Warn("scheduler is not running")
		return
	}

	// Stop accepting new tasks
	s.cron.Stop()

	// Cancel context to signal running tasks to stop
	s.cancel()

	// Release the worker pool
	s.pool.Release()

	s.running = false
	s.logger.Info("task scheduler stopped")
}

// IsRunning returns whether the scheduler is currently running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

type FuncWithServer func(fiberServer *server.FiberServer) (int, error)
type FuncWithoutServer func() (int, error)

// AddTask adds a scheduled task without server dependency
func (s *Scheduler) AddTask(name, spec string, task FuncWithoutServer) error {
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}
	if spec == "" {
		return fmt.Errorf("task spec cannot be empty")
	}
	if task == nil {
		return fmt.Errorf("task function cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if task already exists
	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("task %s already exists", name)
	}

	// Create a lock for this task
	s.locks[name] = &sync.Mutex{}

	id, err := s.cron.AddFunc(spec, func() {
		// Submit task to worker pool
		submitErr := s.pool.Submit(func() {
			// Try to acquire lock, skip if task is already running
			if !s.locks[name].TryLock() {
				s.logger.Warn("task is already running, skipping execution", logger.Fields{
					"task": name,
				})
				return
			}
			defer s.locks[name].Unlock()

			startTime := time.Now()
			s.logger.Debug("starting task execution", logger.Fields{
				"task": name,
			})

			count, taskErr := task()
			duration := time.Since(startTime)

			if taskErr != nil {
				s.logger.Error("task execution failed", logger.Fields{
					"task":     name,
					"duration": duration,
					"error":    taskErr.Error(),
				})
			} else {
				s.logger.Info("task execution completed", logger.Fields{
					"task":     name,
					"duration": duration,
					"count":    count,
				})
			}
		})

		if submitErr != nil {
			s.logger.Error("failed to submit task to worker pool", logger.Fields{
				"task":  name,
				"error": submitErr.Error(),
			})
		}
	})

	if err != nil {
		delete(s.locks, name)
		return fmt.Errorf("failed to add task %s: %w", name, err)
	}

	s.tasks[name] = id
	s.logger.Info("task added successfully", logger.Fields{
		"task": name,
		"spec": spec,
	})

	return nil
}

// AddTaskWithServer adds a scheduled task that requires server dependency
func (s *Scheduler) AddTaskWithServer(name, spec string, task FuncWithServer, fiberServer *server.FiberServer) error {
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}
	if spec == "" {
		return fmt.Errorf("task spec cannot be empty")
	}
	if task == nil {
		return fmt.Errorf("task function cannot be nil")
	}
	if fiberServer == nil {
		return fmt.Errorf("fiber server cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if task already exists
	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("task %s already exists", name)
	}

	// Create a lock for this task
	s.locks[name] = &sync.Mutex{}

	id, err := s.cron.AddFunc(spec, func() {
		// Submit task to worker pool
		submitErr := s.pool.Submit(func() {
			// Try to acquire lock, skip if task is already running
			if !s.locks[name].TryLock() {
				s.logger.Warn("task is already running, skipping execution", logger.Fields{
					"task": name,
				})
				return
			}
			defer s.locks[name].Unlock()

			startTime := time.Now()
			s.logger.Debug("starting task execution with server", logger.Fields{
				"task": name,
			})

			count, taskErr := task(fiberServer)
			duration := time.Since(startTime)

			if taskErr != nil {
				s.logger.Error("task execution failed", logger.Fields{
					"task":     name,
					"duration": duration,
					"error":    taskErr.Error(),
				})
			} else {
				s.logger.Info("task execution completed", logger.Fields{
					"task":     name,
					"duration": duration,
					"count":    count,
				})
			}
		})

		if submitErr != nil {
			s.logger.Error("failed to submit task to worker pool", logger.Fields{
				"task":  name,
				"error": submitErr.Error(),
			})
		}
	})

	if err != nil {
		delete(s.locks, name)
		return fmt.Errorf("failed to add task %s: %w", name, err)
	}

	s.tasks[name] = id
	s.logger.Info("task with server added successfully", logger.Fields{
		"task": name,
		"spec": spec,
	})

	return nil
}

// RemoveTask removes a scheduled task
func (s *Scheduler) RemoveTask(name string) error {
	if name == "" {
		return fmt.Errorf("task name cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id, exists := s.tasks[name]
	if !exists {
		return fmt.Errorf("task %s does not exist", name)
	}

	s.cron.Remove(id)
	delete(s.tasks, name)
	delete(s.locks, name)

	s.logger.Info("task removed successfully", logger.Fields{
		"task": name,
	})

	return nil
}

// GetTaskCount returns the number of registered tasks
func (s *Scheduler) GetTaskCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks)
}

// ListTasks returns a list of all registered task names
func (s *Scheduler) ListTasks() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]string, 0, len(s.tasks))
	for name := range s.tasks {
		tasks = append(tasks, name)
	}
	return tasks
}
