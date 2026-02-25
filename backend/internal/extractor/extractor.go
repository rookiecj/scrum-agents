package extractor

import (
	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// Extractor defines the interface for content extractors.
type Extractor interface {
	Extract(url string) (*model.ExtractedContent, error)
}
