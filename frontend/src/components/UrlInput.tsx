import { useState, FormEvent } from 'react'

interface UrlInputProps {
  onSubmit: (url: string, provider: 'claude' | 'openai') => void
  isLoading: boolean
}

export function UrlInput({ onSubmit, isLoading }: UrlInputProps) {
  const [url, setUrl] = useState('')
  const [provider, setProvider] = useState<'claude' | 'openai'>('claude')

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (url.trim()) {
      onSubmit(url.trim(), provider)
    }
  }

  return (
    <form onSubmit={handleSubmit} style={{ marginBottom: '2rem' }}>
      <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.75rem' }}>
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com/article"
          required
          disabled={isLoading}
          aria-label="URL"
          style={{
            flex: 1,
            padding: '0.75rem 1rem',
            fontSize: '1rem',
            border: '2px solid #e2e8f0',
            borderRadius: '8px',
            outline: 'none',
          }}
        />
        <button
          type="submit"
          disabled={isLoading || !url.trim()}
          style={{
            padding: '0.75rem 1.5rem',
            fontSize: '1rem',
            backgroundColor: isLoading ? '#a0aec0' : '#3182ce',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            cursor: isLoading ? 'not-allowed' : 'pointer',
          }}
        >
          {isLoading ? 'Processing...' : 'Summarize'}
        </button>
      </div>
      <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
        <span style={{ fontSize: '0.875rem', color: '#718096' }}>LLM Provider:</span>
        <label style={{ fontSize: '0.875rem', cursor: 'pointer' }}>
          <input
            type="radio"
            name="provider"
            value="claude"
            checked={provider === 'claude'}
            onChange={() => setProvider('claude')}
            disabled={isLoading}
          />{' '}
          Claude
        </label>
        <label style={{ fontSize: '0.875rem', cursor: 'pointer' }}>
          <input
            type="radio"
            name="provider"
            value="openai"
            checked={provider === 'openai'}
            onChange={() => setProvider('openai')}
            disabled={isLoading}
          />{' '}
          OpenAI
        </label>
      </div>
    </form>
  )
}
