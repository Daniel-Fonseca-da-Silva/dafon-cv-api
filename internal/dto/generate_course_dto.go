package dto

// GenerateCoursesAIRequest represents the request structure for AI filtering
type GenerateCoursesAIRequest struct {
	Content string `json:"content" binding:"required,min=3,max=20000"`
}

// GenerateCoursesAIResponse represents the response structure for AI filtering
type GenerateCoursesAIResponse struct {
	FilteredContent string `json:"filtered_content"`
}
