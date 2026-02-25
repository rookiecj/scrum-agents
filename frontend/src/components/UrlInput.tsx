import { useState, useEffect, FormEvent } from 'react'
import type { ProviderInfo, ProviderName } from '../types/api'

interface UrlInputProps {
  onSubmit: (url: string, provider: ProviderName) => void
  isLoading: boolean
}

export function UrlInput({ onSubmit, isLoading }: UrlInputProps) {
  const [url, setUrl] = useState('')
  const [provider, setProvider] = useState<ProviderName>('claude')
  const [providers, setProviders] = useState<ProviderInfo[]>([])
  const [providersLoaded, setProvidersLoaded] = useState(false)

  useEffect(() => {
    fetch('/api/providers')
      .then((res) => res.json())
      .then((data: ProviderInfo[]) => {
        setProviders(data)
        setProvidersLoaded(true)

        // Auto-select the first available provider
        const firstAvailable = data.find((p) => p.available)
        if (firstAvailable) {
          setProvider(firstAvailable.name as ProviderName)
        }
      })
      .catch(() => {
        // If fetch fails, show all providers as available (fallback)
        setProviders([
          { name: 'claude', available: true, envVar: 'ANTHROPIC_API_KEY' },
          { name: 'openai', available: true, envVar: 'OPENAI_API_KEY' },
          { name: 'gemini', available: true, envVar: 'GOOGLE_API_KEY' },
        ])
        setProvidersLoaded(true)
      })
  }, [])

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    if (url.trim()) {
      onSubmit(url.trim(), provider)
    }
  }

  const isProviderDisabled = (name: string): boolean => {
    if (!providersLoaded) return false
    const info = providers.find((p) => p.name === name)
    return info ? !info.available : false
  }

  const getProviderTooltip = (name: string): string | undefined => {
    if (!providersLoaded) return undefined
    const info = providers.find((p) => p.name === name)
    if (info && !info.available) {
      return `API 키 설정이 필요합니다 (${info.envVar})`
    }
    return undefined
  }

  const hasAnyAvailable = providers.some((p) => p.available)

  const providerLabel = (name: string): string => {
    switch (name) {
      case 'claude':
        return 'Claude'
      case 'openai':
        return 'OpenAI'
      case 'gemini':
        return 'Gemini'
      default:
        return name
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
      {providersLoaded && !hasAnyAvailable && (
        <div
          role="alert"
          style={{
            padding: '0.75rem 1rem',
            marginBottom: '0.75rem',
            backgroundColor: '#fff5f5',
            color: '#c53030',
            border: '1px solid #feb2b2',
            borderRadius: '8px',
            fontSize: '0.875rem',
          }}
        >
          사용 가능한 LLM provider가 없습니다. .env 파일에 API 키를 설정해 주세요.
        </div>
      )}
      <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
        <span style={{ fontSize: '0.875rem', color: '#718096' }}>LLM Provider:</span>
        {(providersLoaded ? providers : [
          { name: 'claude', available: true, envVar: 'ANTHROPIC_API_KEY' },
          { name: 'openai', available: true, envVar: 'OPENAI_API_KEY' },
          { name: 'gemini', available: true, envVar: 'GOOGLE_API_KEY' },
        ]).map((p) => {
          const disabled = isLoading || isProviderDisabled(p.name)
          return (
            <label
              key={p.name}
              style={{
                fontSize: '0.875rem',
                cursor: disabled ? 'not-allowed' : 'pointer',
                opacity: isProviderDisabled(p.name) ? 0.5 : 1,
              }}
              title={getProviderTooltip(p.name)}
            >
              <input
                type="radio"
                name="provider"
                value={p.name}
                checked={provider === p.name}
                onChange={() => setProvider(p.name as ProviderName)}
                disabled={disabled}
              />{' '}
              {providerLabel(p.name)}
            </label>
          )
        })}
      </div>
    </form>
  )
}
