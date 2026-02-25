import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { ErrorBoundary } from './ErrorBoundary'

// Mock the logger module to verify logging calls
vi.mock('../utils/logger', () => ({
  logger: {
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  },
}))

// Import the mocked logger so we can inspect calls
import { logger } from '../utils/logger'

/** A component that always throws an error when rendered. */
function ThrowingComponent({ message }: { message: string }): never {
  throw new Error(message)
}

/** A normal component that renders without error. */
function GoodComponent() {
  return <div>All good</div>
}

describe('ErrorBoundary', () => {
  let consoleErrorSpy: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    // React logs caught errors to console.error; suppress during tests
    consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    vi.clearAllMocks()
  })

  afterEach(() => {
    consoleErrorSpy.mockRestore()
  })

  it('renders children when no error occurs', () => {
    render(
      <ErrorBoundary>
        <GoodComponent />
      </ErrorBoundary>,
    )
    expect(screen.getByText('All good')).toBeInTheDocument()
  })

  it('renders default fallback UI when an error is thrown', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent message="boom" />
      </ErrorBoundary>,
    )
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText('Something went wrong')).toBeInTheDocument()
    expect(screen.getByText('boom')).toBeInTheDocument()
  })

  it('renders custom fallback when provided', () => {
    render(
      <ErrorBoundary fallback={<div>Custom error page</div>}>
        <ThrowingComponent message="oops" />
      </ErrorBoundary>,
    )
    expect(screen.getByText('Custom error page')).toBeInTheDocument()
  })

  it('logs the error at Error level via the logger', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent message="test error" />
      </ErrorBoundary>,
    )

    expect(logger.error).toHaveBeenCalledTimes(1)
    expect(logger.error).toHaveBeenCalledWith(
      'Unhandled React error caught by ErrorBoundary',
      expect.objectContaining({
        errorMessage: 'test error',
        errorName: 'Error',
      }),
    )
  })

  it('includes componentStack in the logged context', () => {
    render(
      <ErrorBoundary>
        <ThrowingComponent message="stack test" />
      </ErrorBoundary>,
    )

    const loggedContext = (logger.error as ReturnType<typeof vi.fn>).mock.calls[0][1] as Record<
      string,
      unknown
    >
    // componentStack should be present (React provides it as a string)
    expect(loggedContext).toHaveProperty('componentStack')
  })
})
