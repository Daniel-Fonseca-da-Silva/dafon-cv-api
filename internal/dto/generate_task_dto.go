package dto

// GenerateTaskAIRequest represents the request structure for AI filtering
type GenerateTaskAIRequest struct {
	Content string `json:"content" binding:"required,min=3,max=20000"`
}

// GenerateTaskAIResponse represents the response structure for AI filtering
type GenerateTaskAIResponse struct {
	FilteredContent string `json:"filtered_content"`
}
