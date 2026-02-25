import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import { UrlInput } from './UrlInput'

describe('UrlInput', () => {
  it('renders URL input and submit button', () => {
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    expect(screen.getByLabelText('URL')).toBeInTheDocument()
    expect(screen.getByText('Summarize')).toBeInTheDocument()
  })

  it('renders provider radio buttons', () => {
    render(<UrlInput onSubmit={vi.fn()} isLoading={false} />)

    expect(screen.getByLabelText('Claude')).toBeInTheDocument()
    expect(screen.getByLabelText('OpenAI')).toBeInTheDocument()
    expect(screen.getByLabelText('Gemini')).toBeInTheDocument()
  })

  it('calls onSubmit with URL and provider on form submit', () => {
    const onSubmit = vi.fn()
    render(<UrlInput onSubmit={onSubmit} isLoading={false} />)

    fireEvent.change(screen.getByLabelText('URL'), {
      target: { value: 'https://example.com/article' },
    })
    fireEvent.click(screen.getByText('Summarize'))

    expect(onSubmit).toHaveBeenCalledWith('https://example.com/article', 'claude')
  })

  it('shows Processing text when loading', () => {
    render(<UrlInput onSubmit={vi.fn()} isLoading={true} />)

    expect(screen.getByText('Processing...')).toBeInTheDocument()
  })

  it('disables input and button when loading', () => {
    render(<UrlInput onSubmit={vi.fn()} isLoading={true} />)

    expect(screen.getByLabelText('URL')).toBeDisabled()
    expect(screen.getByText('Processing...')).toBeDisabled()
  })

  it('switches provider to OpenAI', () => {
    const onSubmit = vi.fn()
    render(<UrlInput onSubmit={onSubmit} isLoading={false} />)

    fireEvent.click(screen.getByLabelText('OpenAI'))
    fireEvent.change(screen.getByLabelText('URL'), {
      target: { value: 'https://example.com' },
    })
    fireEvent.click(screen.getByText('Summarize'))

    expect(onSubmit).toHaveBeenCalledWith('https://example.com', 'openai')
  })

  it('switches provider to Gemini', () => {
    const onSubmit = vi.fn()
    render(<UrlInput onSubmit={onSubmit} isLoading={false} />)

    fireEvent.click(screen.getByLabelText('Gemini'))
    fireEvent.change(screen.getByLabelText('URL'), {
      target: { value: 'https://example.com' },
    })
    fireEvent.click(screen.getByText('Summarize'))

    expect(onSubmit).toHaveBeenCalledWith('https://example.com', 'gemini')
  })
})
