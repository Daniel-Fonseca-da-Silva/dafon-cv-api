package dto

// GenerateSkillAIRequest represents the request structure for AI filtering
type GenerateSkillAIRequest struct {
	Content string `json:"content" binding:"required,min=10,max=20000"`
}

// GenerateSkillAIResponse represents the response structure for AI filtering
type GenerateSkillAIResponse struct {
	FilteredContent string `json:"filtered_content"`
}
