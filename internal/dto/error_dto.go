package dto

// ErrorResponse represents a simple error message response.
// The actual message varies by status (e.g. 400: "invalid user ID format", 404: "user not found").
type ErrorResponse struct {
	Error string `json:"error" example:"user not found"`
}

// MessageResponse represents a simple message response (e.g. delete success).
type MessageResponse struct {
	Message string `json:"message" example:"User deleted successfully"`
}

// ErrorResponseValidation represents error response for validation (e.g. generate AI 400).
type ErrorResponseValidation struct {
	Error string `json:"error" example:"Validation error"`
}

// ErrorResponseServer represents error response for server errors (e.g. generate AI 500).
type ErrorResponseServer struct {
	Error string `json:"error" example:"Internal server error"`
}
