package workerpool

// Custom errors for PDF worker pool operations
var (
	ErrQueueFull         = &PDFError{message: "queue is full"}
	ErrPoolStopped       = &PDFError{message: "worker pool is stopped"}
	ErrWorkerUnavailable = &PDFError{message: "no worker available"}
)

// PDFError represents an error that can occur during PDF worker pool operations
type PDFError struct {
	message string
}

// Error returns the error message
func (e *PDFError) Error() string {
	return e.message
}
