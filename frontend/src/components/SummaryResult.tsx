import type { SummarizeResponse } from '../types/api'

interface SummaryResultProps {
  result: SummarizeResponse
}

const categoryColors: Record<string, string> = {
  '원리소개': '#805ad5',
  '사용기': '#dd6b20',
  '생각정리': '#38a169',
  '기술소개': '#3182ce',
  '튜토리얼': '#d53f8c',
  '뉴스/분석': '#e53e3e',
}

export function SummaryResult({ result }: SummaryResultProps) {
  if (result.error) {
    return (
      <div
        role="alert"
        style={{
          padding: '1rem',
          backgroundColor: '#fed7d7',
          border: '1px solid #fc8181',
          borderRadius: '8px',
          color: '#c53030',
        }}
      >
        <strong>Error:</strong> {result.error}
      </div>
    )
  }

  const badgeColor = categoryColors[result.classification?.primary] || '#718096'

  return (
    <div
      style={{
        border: '1px solid #e2e8f0',
        borderRadius: '12px',
        padding: '1.5rem',
        backgroundColor: '#fff',
      }}
    >
      <div style={{ marginBottom: '1rem' }}>
        <h2 style={{ margin: '0 0 0.5rem', fontSize: '1.25rem' }}>
          {result.link_info.title || result.link_info.url}
        </h2>
        <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
          <span
            style={{
              display: 'inline-block',
              padding: '0.25rem 0.75rem',
              borderRadius: '9999px',
              fontSize: '0.75rem',
              fontWeight: 'bold',
              color: 'white',
              backgroundColor: badgeColor,
            }}
          >
            {result.classification?.primary}
          </span>
          <span
            style={{
              display: 'inline-block',
              padding: '0.25rem 0.75rem',
              borderRadius: '9999px',
              fontSize: '0.75rem',
              color: '#4a5568',
              backgroundColor: '#edf2f7',
            }}
          >
            {result.link_info.link_type}
          </span>
          {result.classification?.confidence && (
            <span
              style={{
                display: 'inline-block',
                padding: '0.25rem 0.75rem',
                borderRadius: '9999px',
                fontSize: '0.75rem',
                color: '#718096',
                backgroundColor: '#f7fafc',
              }}
            >
              Confidence: {Math.round(result.classification.confidence * 100)}%
            </span>
          )}
        </div>
      </div>

      <div
        style={{
          whiteSpace: 'pre-wrap',
          lineHeight: 1.7,
          color: '#2d3748',
          fontSize: '0.95rem',
        }}
      >
        {result.summary}
      </div>

      {result.link_info.author && (
        <div style={{ marginTop: '1rem', fontSize: '0.85rem', color: '#718096' }}>
          Author: {result.link_info.author}
        </div>
      )}
    </div>
  )
}
