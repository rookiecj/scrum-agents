package classifier

import (
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// Classifier classifies content into a category.
type Classifier interface {
	Classify(content string) (*model.ClassificationResult, error)
}
