package classifier

import "fmt"

// ClassificationPrompt returns the prompt used to classify content.
func ClassificationPrompt(content string) string {
	return fmt.Sprintf(`You are a content classifier. Classify the following content into exactly one of these categories:

1. 원리소개 - Explains a principle, concept, or how something works (e.g., "How TCP works", "양자컴퓨팅 원리")
2. 사용기 - Product/tool/service usage review or experience (e.g., "M4 MacBook Pro 한달 사용기", "Cursor IDE 리뷰")
3. 생각정리 - Opinion, essay, or philosophical reflection (e.g., "AI가 개발자를 대체할까", "스타트업 문화에 대한 단상")
4. 기술소개 - Introduction of a new technology/tool/framework (e.g., "Introducing Bun 1.0", "Go 1.22 새 기능")
5. 튜토리얼 - Step-by-step guide or how-to (e.g., "React에서 상태관리 구현하기", "Docker 입문")
6. 뉴스/분석 - Industry news and trend analysis (e.g., "2024 AI 트렌드 리포트", "OpenAI DevDay 정리")

Respond ONLY with a JSON object in this exact format:
{"primary": "<category>", "confidence": <0.0-1.0>, "secondary": "<category>", "secondary_confidence": <0.0-1.0>}

Content to classify:
---
%s
---`, truncate(content, 4000))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
