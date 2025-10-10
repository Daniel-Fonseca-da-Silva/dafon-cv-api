package dto

// GenerateTranslationAIRequest represents the request structure for AI translation
type GenerateTranslationAIRequest struct {
	CurriculumData map[string]interface{} `json:"curriculum_data" binding:"required"`
	TargetLanguage string                 `json:"target_language" binding:"required,oneof=pt en es"`
}

// GenerateTranslationAIResponse represents the response structure for AI translation
type GenerateTranslationAIResponse struct {
	TranslatedCurriculum map[string]interface{} `json:"translated_curriculum"`
}
