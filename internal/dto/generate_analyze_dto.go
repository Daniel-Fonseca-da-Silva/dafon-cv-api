package dto

// GenerateAnalyzeAIRequest represents the request structure for AI filtering
type GenerateAnalyzeAIRequest struct {
	Content string `json:"content" binding:"required,min=500,max=20000" example:"string"`
}

// GenerateAnalyzeAIResponse represents the response structure for AI curriculum analysis
type GenerateAnalyzeAIResponse struct {
	Score                 float64           `json:"score,omitempty"`
	Description           string            `json:"description,omitempty"`
	ImprovementPoints     []string          `json:"improvement_points,omitempty"`
	BestPractices         []string          `json:"best_practices,omitempty"`
	ATSCompatibility      *ATSCompatibility `json:"ats_compatibility,omitempty"`
	ProfessionalAlignment []string          `json:"professional_alignment,omitempty"`
	Strengths             []string          `json:"strengths,omitempty"`
	Recommendations       []string          `json:"recommendations,omitempty"`
}

// ATSCompatibility represents ATS-related evaluation
type ATSCompatibility struct {
	Assessment      string   `json:"assessment,omitempty"`
	Chance          string   `json:"chance,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}
