import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { UrlInput } from './UrlInput'

// Mock fetch globally for provider API calls
const mockFetch = vi.fn()

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch)
})

afterEach(() => {
  vi.restoreAllMocks()
})

function mockProvidersResponse(providers: Array<{ name: string; available: boolean; envVar: string }>) {
  mockFetch.mockResolvedValueOnce({
    ok: true,
    json: async () => providers,
  })
}

const allAvailable = [
  { name: 'claude', available: true, envVar: 'ANTHROPIC_API_KEY' },
  { name: 'openai', available: true, envVar: 'OPENAI_API_KEY' },
  { name: 'gemini', available: true, envVar: 'GOOGLE_API_KEY' },
]

const onlyClaudeAvailable = [
  { name: 'claude', available: true, envVar: 'ANTHROPIC_API_KEY' },
  { name: 'openai', available: false, envVar: 'OPENAI_API_KEY' },
  { name: 'gemini', available: false, envVar: 'GOOGLE_API_KEY' },
]

const noneAvailable = [
  { name: 'claude', available: false, envVar: 'ANTHROPIC_API_KEY' },
  { name: 'openai', available: false, envVar: 'OPENAI_API_KEY' },
  { name: 'gemini', available: false, envVar: 'GOOGLE_API_KEY' },
]

describe('UrlInput', () => {
  it('renders URL input and submit button', async () => {
    mockProvidersResponse(allAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('URL')).toBeInTheDocument()
      expect(screen.getByText('Summarize')).toBeInTheDocument()
    })
  })

  it('renders provider radio buttons', async () => {
    mockProvidersResponse(allAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('Claude')).toBeInTheDocument()
      expect(screen.getByLabelText('OpenAI')).toBeInTheDocument()
      expect(screen.getByLabelText('Gemini')).toBeInTheDocument()
    })
  })

  it('calls onSubmit with URL and provider on form submit', async () => {
    mockProvidersResponse(allAvailable)
    const onSubmit = vi.fn()
    render(<UrlInput onSubmit={onSubmit} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('Claude')).not.toBeDisabled()
    })

    fireEvent.change(screen.getByLabelText('URL'), {
      target: { value: 'https://example.com/article' },
    })
    fireEvent.click(screen.getByText('Summarize'))

    expect(onSubmit).toHaveBeenCalledWith('https://example.com/article', 'claude')
  })

  it('shows Processing text when loading', async () => {
    mockProvidersResponse(allAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={true} />)

    await waitFor(() => {
      expect(screen.getByText('Processing...')).toBeInTheDocument()
    })
  })

  it('disables input and button when loading', async () => {
    mockProvidersResponse(allAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={true} />)

    await waitFor(() => {
      expect(screen.getByLabelText('URL')).toBeDisabled()
      expect(screen.getByText('Processing...')).toBeDisabled()
    })
  })

  it('switches provider to OpenAI', async () => {
    mockProvidersResponse(allAvailable)
    const onSubmit = vi.fn()
    render(<UrlInput onSubmit={onSubmit} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('OpenAI')).not.toBeDisabled()
    })

    fireEvent.click(screen.getByLabelText('OpenAI'))
    fireEvent.change(screen.getByLabelText('URL'), {
      target: { value: 'https://example.com' },
    })
    fireEvent.click(screen.getByText('Summarize'))

    expect(onSubmit).toHaveBeenCalledWith('https://example.com', 'openai')
  })

  it('switches provider to Gemini', async () => {
    mockProvidersResponse(allAvailable)
    const onSubmit = vi.fn()
    render(<UrlInput onSubmit={onSubmit} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('Gemini')).not.toBeDisabled()
    })

    fireEvent.click(screen.getByLabelText('Gemini'))
    fireEvent.change(screen.getByLabelText('URL'), {
      target: { value: 'https://example.com' },
    })
    fireEvent.click(screen.getByText('Summarize'))

    expect(onSubmit).toHaveBeenCalledWith('https://example.com', 'gemini')
  })

  it('disables unavailable providers', async () => {
    mockProvidersResponse(onlyClaudeAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('Claude')).not.toBeDisabled()
      expect(screen.getByLabelText('OpenAI')).toBeDisabled()
      expect(screen.getByLabelText('Gemini')).toBeDisabled()
    })
  })

  it('shows tooltip on disabled providers', async () => {
    mockProvidersResponse(onlyClaudeAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    await waitFor(() => {
      const openaiLabel = screen.getByLabelText('OpenAI').closest('label')
      expect(openaiLabel).toHaveAttribute('title', 'API 키 설정이 필요합니다 (OPENAI_API_KEY)')

      const geminiLabel = screen.getByLabelText('Gemini').closest('label')
      expect(geminiLabel).toHaveAttribute('title', 'API 키 설정이 필요합니다 (GOOGLE_API_KEY)')
    })
  })

  it('auto-selects first available provider', async () => {
    const onlyGeminiAvailable = [
      { name: 'claude', available: false, envVar: 'ANTHROPIC_API_KEY' },
      { name: 'openai', available: false, envVar: 'OPENAI_API_KEY' },
      { name: 'gemini', available: true, envVar: 'GOOGLE_API_KEY' },
    ]
    mockProvidersResponse(onlyGeminiAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByLabelText('Gemini')).toBeChecked()
    })
  })

  it('shows warning when no providers available', async () => {
    mockProvidersResponse(noneAvailable)
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent(
        '사용 가능한 LLM provider가 없습니다'
      )
    })
  })

  it('falls back gracefully when fetch fails', async () => {
    mockFetch.mockRejectedValueOnce(new Error('network error'))
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    // Should still render all providers as enabled
    await waitFor(() => {
      expect(screen.getByLabelText('Claude')).not.toBeDisabled()
      expect(screen.getByLabelText('OpenAI')).not.toBeDisabled()
      expect(screen.getByLabelText('Gemini')).not.toBeDisabled()
    })
  })
})
