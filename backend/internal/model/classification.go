package model

// ContentCategory represents the classification of content.
type ContentCategory string

const (
	CategoryPrinciple  ContentCategory = "원리소개"   // Principle/concept explanation
	CategoryReview     ContentCategory = "사용기"     // Usage review/experience
	CategoryOpinion    ContentCategory = "생각정리"   // Opinion/essay
	CategoryTechIntro  ContentCategory = "기술소개"   // Technology introduction
	CategoryTutorial   ContentCategory = "튜토리얼"   // Step-by-step guide
	CategoryNews       ContentCategory = "뉴스/분석"  // News/analysis
)

// AllCategories returns all valid content categories.
func AllCategories() []ContentCategory {
	return []ContentCategory{
		CategoryPrinciple,
		CategoryReview,
		CategoryOpinion,
		CategoryTechIntro,
		CategoryTutorial,
		CategoryNews,
	}
}

// ClassificationResult holds the result of content classification.
type ClassificationResult struct {
	Primary    ContentCategory `json:"primary"`
	Confidence float64         `json:"confidence"`
	Secondary  ContentCategory `json:"secondary,omitempty"`
	SecondConf float64         `json:"secondary_confidence,omitempty"`
}
