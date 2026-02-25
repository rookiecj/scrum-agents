import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { ProgressIndicator } from './ProgressIndicator'

describe('ProgressIndicator', () => {
  it('renders progress steps when detecting', () => {
    render(<ProgressIndicator currentStep="detecting" />)

    expect(screen.getByText('Detecting link type')).toBeInTheDocument()
    expect(screen.getByText('Extracting content')).toBeInTheDocument()
    expect(screen.getByText('Classifying content')).toBeInTheDocument()
    expect(screen.getByText('Generating summary')).toBeInTheDocument()
  })

  it('returns null when step is done', () => {
    const { container } = render(<ProgressIndicator currentStep="done" />)
    expect(container.firstChild).toBeNull()
  })

  it('returns null when step is error', () => {
    const { container } = render(<ProgressIndicator currentStep="error" />)
    expect(container.firstChild).toBeNull()
  })

  it('has progressbar role', () => {
    render(<ProgressIndicator currentStep="extracting" />)
    expect(screen.getByRole('progressbar')).toBeInTheDocument()
  })
})
