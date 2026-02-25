export type LinkType = 'article' | 'youtube' | 'pdf' | 'twitter' | 'newsletter' | 'unknown'

export type ContentCategory =
  | '원리소개'
  | '사용기'
  | '생각정리'
  | '기술소개'
  | '튜토리얼'
  | '뉴스/분석'

export interface LinkInfo {
  url: string
  link_type: LinkType
  title?: string
  author?: string
  date?: string
}

export interface ClassificationResult {
  primary: ContentCategory
  confidence: number
  secondary?: ContentCategory
  secondary_confidence?: number
}

export interface SummarizeRequest {
  url: string
  provider?: 'claude' | 'openai' | 'gemini'
}

export interface SummarizeResponse {
  link_info: LinkInfo
  classification: ClassificationResult
  summary: string
  error?: string
}

export type SummarizeStep = 'detecting' | 'extracting' | 'classifying' | 'summarizing' | 'done' | 'error'

export type ProviderName = 'claude' | 'openai' | 'gemini'

export interface ProviderInfo {
  name: string
  available: boolean
  envVar: string
}
