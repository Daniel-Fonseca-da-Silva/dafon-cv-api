package handlers

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/usecases"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/utils"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/workerpool"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ModernPDFHandler handles HTTP requests for modern PDF operations.
type ModernPDFHandler struct {
	workerPool *workerpool.PDFWorkerPool
	logger     *zap.Logger
	mu         sync.RWMutex // Protects the workerPool during shutdown

	// Centralized monitoring system
	monitorCtx    context.Context
	monitorCancel context.CancelFunc
	monitorWg     sync.WaitGroup
}

// NewModernPDFHandler creates a new instance of ModernPDFHandler.
func NewModernPDFHandler(simplePDFUseCase usecases.SimplePDFUseCase, logger *zap.Logger, numWorkers int, queueSize int) *ModernPDFHandler {
	// Create the worker pool with custom configuration.
	workerPool := workerpool.NewPDFWorkerPoolWithQueueSize(numWorkers, queueSize, simplePDFUseCase, logger)
	workerPool.Start()

	// Create context for centralized monitoring.
	monitorCtx, monitorCancel := context.WithCancel(context.Background())

	handler := &ModernPDFHandler{
		workerPool:    workerPool,
		logger:        logger,
		monitorCtx:    monitorCtx,
		monitorCancel: monitorCancel,
	}

	return handler
}

// CreateModernPDF handles POST /pdf-modern/:id request
func (h *ModernPDFHandler) CreateModernPDF(c *gin.Context) {
	// Get the curriculum ID from the URL
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.HandleValidationError(c, errors.New("invalid curriculum ID format"))
		return
	}

	// Verifica se o handler está disponível
	if !h.IsAvailable() {
		utils.HandleValidationError(c, errors.New("PDF service is not available"))
		return
	}

	// Get thread-safe reference to the worker pool.
	h.mu.RLock()
	workerPool := h.workerPool
	h.mu.RUnlock()

	// Response with status accepted and the curriculum ID
	c.JSON(http.StatusAccepted, gin.H{
		"message":        "PDF generation started",
		"curriculum_id":  id.String(),
		"status":         "processing",
		"queue_size":     workerPool.GetQueueSize(),
		"active_workers": workerPool.GetActiveWorkers(),
		"active_jobs":    workerPool.GetActiveJobs(),
	})

	// Submete a tarefa ao worker pool
	resultChan, err := workerPool.Submit(id, c.ClientIP())
	if err != nil {
		h.logger.Error("Failed to submit job to worker pool",
			zap.String("curriculum_id", id.String()),
			zap.String("user_ip", c.ClientIP()),
			zap.Error(err),
		)
		return
	}

	// Monitor the result in background with proper cancellation.
	h.monitorWg.Add(1)
	go func() {
		defer h.monitorWg.Done()

		// Create a context with timeout for monitoring.
		monitorCtx, cancel := context.WithTimeout(h.monitorCtx, 5*time.Minute)
		defer cancel() // Ensure the context is canceled.

		select {
		case err := <-resultChan:
			if err != nil {
				if err == context.Canceled {
					h.logger.Warn("PDF generation canceled (context canceled)",
						zap.String("curriculum_id", id.String()),
						zap.String("user_ip", c.ClientIP()),
					)
				} else {
					h.logger.Error("PDF generation failed",
						zap.String("curriculum_id", id.String()),
						zap.String("user_ip", c.ClientIP()),
						zap.Error(err),
					)
				}
			} else {
				h.logger.Info("PDF generation completed successfully",
					zap.String("curriculum_id", id.String()),
					zap.String("user_ip", c.ClientIP()),
				)
			}
		case <-monitorCtx.Done():
			h.logger.Warn("PDF generation monitoring timed out",
				zap.String("curriculum_id", id.String()),
				zap.String("user_ip", c.ClientIP()),
				zap.Error(monitorCtx.Err()),
			)
		}
	}()
}

// GetPoolStatus return the status of the worker pool (optional endpoint for monitoring).
func (h *ModernPDFHandler) GetPoolStatus(c *gin.Context) {
	h.mu.RLock()
	workerPool := h.workerPool
	h.mu.RUnlock()

	if workerPool == nil {
		utils.HandleValidationError(c, errors.New("worker pool not available"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"queue_size":     workerPool.GetQueueSize(),
		"active_workers": workerPool.GetActiveWorkers(),
		"total_workers":  workerPool.GetTotalWorkers(),
		"active_jobs":    workerPool.GetActiveJobs(),
		"completed_jobs": workerPool.GetCompletedJobs(),
		"failed_jobs":    workerPool.GetFailedJobs(),
		"is_running":     workerPool.IsRunning(),
		"status":         "running",
	})
}

// IsAvailable check if the handler is available to process requests.
func (h *ModernPDFHandler) IsAvailable() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.workerPool != nil && h.workerPool.IsRunning()
}

// StopWorkerPool stop the worker pool (for graceful shutdown).
// NOTE: This method should be called during application shutdown to ensure
// all workers are properly stopped and resources are cleaned up.
// Currently not automatically called - needs to be integrated with signal handling in main.go
func (h *ModernPDFHandler) StopWorkerPool() {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Cancel the monitoring system.
	if h.monitorCancel != nil {
		h.logger.Info("Canceling PDF monitoring system")
		h.monitorCancel()
	}

	// Wait for all monitoring goroutines to finish.
	h.logger.Info("Waiting for monitoring goroutines to finish")
	h.monitorWg.Wait()

	// Stop the worker pool.
	if h.workerPool != nil {
		h.workerPool.Stop()
		h.workerPool = nil
	}

	h.logger.Info("PDF handler shutdown completed")
}
