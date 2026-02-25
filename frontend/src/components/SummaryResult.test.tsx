import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { SummaryResult } from './SummaryResult'
import type { SummarizeResponse } from '../types/api'

describe('SummaryResult', () => {
  const mockResult: SummarizeResponse = {
    link_info: {
      url: 'https://example.com/article',
      link_type: 'article',
      title: 'Test Article',
      author: 'Test Author',
    },
    classification: {
      primary: '기술소개',
      confidence: 0.92,
    },
    summary: 'This is a test summary of the article.',
  }

  it('renders title', () => {
    render(<SummaryResult result={mockResult} />)
    expect(screen.getByText('Test Article')).toBeInTheDocument()
  })

  it('renders category badge', () => {
    render(<SummaryResult result={mockResult} />)
    expect(screen.getByText('기술소개')).toBeInTheDocument()
  })

  it('renders link type badge', () => {
    render(<SummaryResult result={mockResult} />)
    expect(screen.getByText('article')).toBeInTheDocument()
  })

  it('renders confidence percentage', () => {
    render(<SummaryResult result={mockResult} />)
    expect(screen.getByText('Confidence: 92%')).toBeInTheDocument()
  })

  it('renders summary text', () => {
    render(<SummaryResult result={mockResult} />)
    expect(screen.getByText('This is a test summary of the article.')).toBeInTheDocument()
  })

  it('renders author', () => {
    render(<SummaryResult result={mockResult} />)
    expect(screen.getByText('Author: Test Author')).toBeInTheDocument()
  })

  it('renders error state', () => {
    const errorResult: SummarizeResponse = {
      ...mockResult,
      error: 'Something went wrong',
    }
    render(<SummaryResult result={errorResult} />)
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/Something went wrong/)).toBeInTheDocument()
  })

  it('uses URL when title is missing', () => {
    const noTitleResult: SummarizeResponse = {
      ...mockResult,
      link_info: { ...mockResult.link_info, title: undefined },
    }
    render(<SummaryResult result={noTitleResult} />)
    expect(screen.getByText('https://example.com/article')).toBeInTheDocument()
  })
})
