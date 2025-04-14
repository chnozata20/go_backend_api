package worker

import (
	"context"
	"sync"
	"sync/atomic"
)

type TaskProcessor struct {
	tasks    chan Task
	workers  int
	wg       sync.WaitGroup
	stats    *TaskStats
}

type Task interface {
	Execute(ctx context.Context) error
}

type TaskStats struct {
	TotalProcessed uint64
	TotalErrors    uint64
}

func NewTaskProcessor(workers int) *TaskProcessor {
	return &TaskProcessor{
		tasks:   make(chan Task, 1000),
		workers: workers,
		stats:   &TaskStats{},
	}
}

func (p *TaskProcessor) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.processTasks(ctx)
	}
}

func (p *TaskProcessor) Stop() {
	close(p.tasks)
	p.wg.Wait()
}

func (p *TaskProcessor) Submit(task Task) {
	p.tasks <- task
}

func (p *TaskProcessor) SubmitBatch(tasks []Task) {
	for _, task := range tasks {
		p.Submit(task)
	}
}

func (p *TaskProcessor) processTasks(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			
			if err := task.Execute(ctx); err != nil {
				atomic.AddUint64(&p.stats.TotalErrors, 1)
				continue
			}

			atomic.AddUint64(&p.stats.TotalProcessed, 1)
		}
	}
}

func (p *TaskProcessor) GetStats() TaskStats {
	return *p.stats
} 