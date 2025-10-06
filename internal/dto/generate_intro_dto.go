package dto

// GenerateIntroAIRequest represents the request structure for AI filtering
type GenerateIntroAIRequest struct {
	Content string `json:"content" binding:"required,min=3,max=20000"`
}

// GenerateIntroAIResponse represents the response structure for AI filtering
type GenerateIntroAIResponse struct {
	FilteredContent string `json:"filtered_content"`
}
