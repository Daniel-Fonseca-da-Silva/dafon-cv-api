package dto

// GenerateTranslationAIRequest represents the request structure for AI translation (runtime: flexible map).
type GenerateTranslationAIRequest struct {
	CurriculumData map[string]interface{} `json:"curriculum_data" binding:"required"`
	TargetLanguage string                 `json:"target_language" binding:"required,oneof=pt en es"`
}

// GenerateTranslationAIRequestDoc is the documented request for Swagger (same shape as curriculum + target_language).
type GenerateTranslationAIRequestDoc struct {
	CurriculumData CurriculumResponse `json:"curriculum_data" binding:"required"`
	TargetLanguage string             `json:"target_language" binding:"required,oneof=pt en es"`
}

// GenerateTranslationAIResponse represents the response structure for AI translation (runtime: flexible map).
type GenerateTranslationAIResponse struct {
	TranslatedCurriculum map[string]interface{} `json:"translated_curriculum"`
}
