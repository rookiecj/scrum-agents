import { Component, type ErrorInfo, type ReactNode } from 'react'
import { logger } from '../utils/logger'

interface ErrorBoundaryProps {
  children: ReactNode
  /** Optional fallback UI to render when an error is caught. */
  fallback?: ReactNode
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

/**
 * React ErrorBoundary that catches unhandled errors in the component tree
 * and logs them at Error level with the component stack trace.
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    logger.error('Unhandled React error caught by ErrorBoundary', {
      errorMessage: error.message,
      errorName: error.name,
      componentStack: errorInfo.componentStack ?? undefined,
    })
  }

  render(): ReactNode {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback
      }

      return (
        <div
          role="alert"
          style={{
            padding: '1.5rem',
            margin: '1rem 0',
            backgroundColor: '#fed7d7',
            border: '1px solid #fc8181',
            borderRadius: '8px',
            color: '#c53030',
          }}
        >
          <h2 style={{ margin: '0 0 0.5rem', fontSize: '1.1rem' }}>Something went wrong</h2>
          <p style={{ margin: 0, fontSize: '0.9rem' }}>
            {this.state.error?.message ?? 'An unexpected error occurred.'}
          </p>
        </div>
      )
    }

    return this.props.children
  }
}
