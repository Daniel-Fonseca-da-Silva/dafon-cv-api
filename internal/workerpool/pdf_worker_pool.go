package workerpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PDFJob represents a PDF generation task.
type PDFJob struct {
	CurriculumID uuid.UUID
	Result       chan error
	UserIP       string
}

// PDFWorkerPool manages workers for PDF generation.
type PDFWorkerPool struct {
	workers    int
	jobQueue   chan PDFJob
	workerPool chan chan PDFJob
	quit       chan bool
	useCase    usecases.SimplePDFUseCase
	logger     *zap.Logger
	wg         sync.WaitGroup

	// Thread-safe state management
	mu            sync.RWMutex
	isRunning     bool
	isStopping    bool
	activeJobs    int32
	completedJobs int32
	failedJobs    int32
	activeWorkers int32 // Atomic counter of active workers
}

// NewPDFWorkerPoolWithQueueSize creates a new worker pool with configurable queue size.
func NewPDFWorkerPoolWithQueueSize(numWorkers, queueSize int, useCase usecases.SimplePDFUseCase, logger *zap.Logger) *PDFWorkerPool {
	pool := &PDFWorkerPool{
		workers:       numWorkers,
		jobQueue:      make(chan PDFJob, queueSize),
		workerPool:    make(chan chan PDFJob, numWorkers),
		quit:          make(chan bool),
		useCase:       useCase,
		logger:        logger,
		isRunning:     false,
		isStopping:    false,
		activeWorkers: 0,
	}

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		pool.startWorker(i)
	}

	logger.Info("PDF Worker Pool initialized",
		zap.Int("num_workers", numWorkers),
		zap.Int("queue_buffer", queueSize),
	)

	return pool
}

// startWorker starts an individual worker.
func (wp *PDFWorkerPool) startWorker(id int) {
	wp.wg.Add(1)
	jobChannel := make(chan PDFJob)

	go func() {
		defer wp.wg.Done()
		defer close(jobChannel)

		// Increment the active workers counter
		atomic.AddInt32(&wp.activeWorkers, 1)
		defer atomic.AddInt32(&wp.activeWorkers, -1)

		wp.logger.Info("Worker started", zap.Int("worker_id", id))

		for {
			// Check if the pool is stopping
			wp.mu.RLock()
			if wp.isStopping {
				wp.mu.RUnlock()
				wp.logger.Info("Worker stopping (pool stopping)", zap.Int("worker_id", id))
				return
			}
			wp.mu.RUnlock()

			// Register with the available worker pool
			select {
			case wp.workerPool <- jobChannel:
				// Worker is available
			case <-wp.quit:
				wp.logger.Info("Worker stopping", zap.Int("worker_id", id))
				return
			}

			select {
			case job := <-jobChannel:
				// Increment the active jobs counter
				atomic.AddInt32(&wp.activeJobs, 1)

				wp.logger.Info("Worker processing job",
					zap.Int("worker_id", id),
					zap.String("curriculum_id", job.CurriculumID.String()),
					zap.String("user_ip", job.UserIP),
				)

				// Create a context with timeout for this operation
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

				// Process the PDF generation
				err := wp.useCase.GeneratePDF(ctx, job.CurriculumID)

				// Cancel the context after processing
				cancel()

				// Decrement the active jobs counter
				atomic.AddInt32(&wp.activeJobs, -1)

				// Update the completed/failed jobs counters
				if err != nil {
					atomic.AddInt32(&wp.failedJobs, 1)
				} else {
					atomic.AddInt32(&wp.completedJobs, 1)
				}

				// Send the result
				select {
				case job.Result <- err:
					// Result sent successfully
				case <-time.After(10 * time.Second):
					wp.logger.Warn("Timeout sending result",
						zap.Int("worker_id", id),
						zap.String("curriculum_id", job.CurriculumID.String()),
					)
				}

				if err != nil {
					if err == context.Canceled {
						wp.logger.Warn("PDF generation canceled (context canceled)",
							zap.Int("worker_id", id),
							zap.String("curriculum_id", job.CurriculumID.String()),
						)
					} else {
						wp.logger.Error("PDF generation failed",
							zap.Int("worker_id", id),
							zap.String("curriculum_id", job.CurriculumID.String()),
							zap.Error(err),
						)
					}
				} else {
					wp.logger.Info("PDF generation completed",
						zap.Int("worker_id", id),
						zap.String("curriculum_id", job.CurriculumID.String()),
					)
				}

			case <-wp.quit:
				wp.logger.Info("Worker stopping", zap.Int("worker_id", id))
				return
			}
		}
	}()
}

// Submit adds a new task to the pool
func (wp *PDFWorkerPool) Submit(curriculumID uuid.UUID, userIP string) (chan error, error) {
	// Check if the pool is running
	wp.mu.RLock()
	if !wp.isRunning || wp.isStopping {
		wp.mu.RUnlock()
		return nil, ErrPoolStopped
	}
	wp.mu.RUnlock()

	result := make(chan error, 1)

	job := PDFJob{
		CurriculumID: curriculumID,
		Result:       result,
		UserIP:       userIP,
	}

	// Try to send the task to the queue with timeout
	select {
	case wp.jobQueue <- job:
		wp.logger.Info("Job submitted to queue",
			zap.String("curriculum_id", curriculumID.String()),
			zap.String("user_ip", userIP),
		)
		return result, nil
	case <-time.After(30 * time.Second):
		close(result)
		return nil, ErrQueueFull
	case <-wp.quit:
		close(result)
		return nil, ErrPoolStopped
	}
}

// dispatcher distributes jobs to available workers
func (wp *PDFWorkerPool) dispatcher() {
	wp.wg.Add(1)
	defer wp.wg.Done()

	wp.logger.Info("Dispatcher started")

	for {
		// Check if the pool is stopping
		wp.mu.RLock()
		if wp.isStopping {
			wp.mu.RUnlock()
			wp.logger.Info("Dispatcher stopping (pool stopping)")
			return
		}
		wp.mu.RUnlock()

		select {
		case job := <-wp.jobQueue:
			// Get an available worker
			select {
			case worker := <-wp.workerPool:
				// Send the job to the worker
				select {
				case worker <- job:
					wp.logger.Debug("Job dispatched to worker",
						zap.String("curriculum_id", job.CurriculumID.String()),
					)
				case <-time.After(5 * time.Second):
					wp.logger.Error("Failed to dispatch job to worker",
						zap.String("curriculum_id", job.CurriculumID.String()),
					)
					job.Result <- ErrWorkerUnavailable
				}
			case <-time.After(10 * time.Second):
				wp.logger.Error("No worker available, job rejected",
					zap.String("curriculum_id", job.CurriculumID.String()),
				)
				job.Result <- ErrWorkerUnavailable
			}

		case <-wp.quit:
			wp.logger.Info("Dispatcher stopping")
			return
		}
	}
}

// Start starts the dispatcher
func (wp *PDFWorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.isRunning {
		wp.logger.Warn("Worker pool is already running")
		return
	}

	wp.isRunning = true
	wp.isStopping = false
	go wp.dispatcher()

	wp.logger.Info("Worker pool started")
}

// Stop stops the worker pool
func (wp *PDFWorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.isStopping || !wp.isRunning {
		wp.logger.Warn("Worker pool is already stopping or not running")
		return
	}

	wp.logger.Info("Stopping PDF Worker Pool")
	wp.isStopping = true
	wp.isRunning = false

	close(wp.quit)
	wp.wg.Wait()

	wp.logger.Info("PDF Worker Pool stopped")
}

// GetQueueSize return the current queue size
func (wp *PDFWorkerPool) GetQueueSize() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if wp.isStopping || !wp.isRunning {
		return 0
	}

	return len(wp.jobQueue)
}

// GetActiveWorkers return the number of active workers
func (wp *PDFWorkerPool) GetActiveWorkers() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if wp.isStopping || !wp.isRunning {
		return 0
	}

	return int(atomic.LoadInt32(&wp.activeWorkers))
}

// GetActiveJobs return the number of active jobs
func (wp *PDFWorkerPool) GetActiveJobs() int32 {
	return atomic.LoadInt32(&wp.activeJobs)
}

// GetCompletedJobs return the number of completed jobs
func (wp *PDFWorkerPool) GetCompletedJobs() int32 {
	return atomic.LoadInt32(&wp.completedJobs)
}

// GetFailedJobs return the number of failed jobs
func (wp *PDFWorkerPool) GetFailedJobs() int32 {
	return atomic.LoadInt32(&wp.failedJobs)
}

// IsRunning return if the pool is running
func (wp *PDFWorkerPool) IsRunning() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.isRunning && !wp.isStopping
}

// GetTotalWorkers return the total number of workers in the pool
func (wp *PDFWorkerPool) GetTotalWorkers() int {
	return wp.workers
}
