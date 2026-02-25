package urldetect

import (
	"net/url"
	"strings"

	"github.com/rookiecj/scrum-agents/backend/internal/model"
)

// Detect analyzes a URL and returns its detected LinkType.
func Detect(rawURL string) (model.LinkType, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return model.LinkTypeUnknown, err
	}

	host := strings.ToLower(u.Hostname())
	path := strings.ToLower(u.Path)

	switch {
	case isYouTube(host, path):
		return model.LinkTypeYouTube, nil
	case isTwitter(host):
		return model.LinkTypeTwitter, nil
	case isPDF(path):
		return model.LinkTypePDF, nil
	case isNewsletter(host):
		return model.LinkTypeNewsletter, nil
	default:
		return model.LinkTypeArticle, nil
	}
}

func isYouTube(host, _ string) bool {
	return strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be")
}

func isTwitter(host string) bool {
	return strings.Contains(host, "twitter.com") || strings.Contains(host, "x.com")
}

func isPDF(path string) bool {
	return strings.HasSuffix(path, ".pdf")
}

func isNewsletter(host string) bool {
	newsletters := []string{"substack.com", "medium.com", "beehiiv.com", "buttondown.email"}
	for _, n := range newsletters {
		if strings.Contains(host, n) {
			return true
		}
	}
	return false
}
