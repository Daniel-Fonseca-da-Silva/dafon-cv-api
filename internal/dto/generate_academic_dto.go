package dto

// GenerateAcademicAIRequest represents the request structure for AI filtering
type GenerateAcademicAIRequest struct {
	Content string `json:"content" binding:"required,min=10,max=20000"`
}

// GenerateAcademicAIResponse represents the response structure for AI filtering
type GenerateAcademicAIResponse struct {
	FilteredContent string `json:"filtered_content"`
}
