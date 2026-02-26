import { useState, useCallback } from 'react'
import { UrlInput } from './components/UrlInput'
import { ProgressIndicator } from './components/ProgressIndicator'
import { SummaryResult } from './components/SummaryResult'
import { LoginForm } from './components/LoginForm'
import { logger } from './utils/logger'
import type { SummarizeResponse, SummarizeStep, ProviderName } from './types/api'

/**
 * Helper that performs a fetch with auth token and logs failures.
 */
function createAuthFetch(authToken: string, onUnauthorized: () => void) {
  return async function fetchWithAuth(url: string, init: RequestInit): Promise<Response> {
    const headers = new Headers(init.headers)
    headers.set('Authorization', `Bearer ${authToken}`)
    const res = await fetch(url, { ...init, headers })
    if (!res.ok) {
      logger.error('API call failed', { url, status: res.status, statusText: res.statusText })
      if (res.status === 401) {
        onUnauthorized()
      }
    }
    return res
  }
}

const AUTH_TOKEN_KEY = 'auth_token'

export function App() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(AUTH_TOKEN_KEY))
  const [step, setStep] = useState<SummarizeStep>('done')
  const [result, setResult] = useState<SummarizeResponse | null>(null)

  const handleLogin = useCallback((newToken: string) => {
    localStorage.setItem(AUTH_TOKEN_KEY, newToken)
    setToken(newToken)
    logger.info('User logged in')
  }, [])

  const handleLogout = useCallback(() => {
    localStorage.removeItem(AUTH_TOKEN_KEY)
    setToken(null)
    setResult(null)
    setStep('done')
    logger.info('User logged out')
  }, [])

  if (!token) {
    return <LoginForm onLogin={handleLogin} />
  }

  const fetchWithAuth = createAuthFetch(token, handleLogout)

  const handleSubmit = async (url: string, provider: ProviderName) => {
    setResult(null)
    logger.info('Starting summarization', { url, provider })

    try {
      // Step 1: Detect
      setStep('detecting')
      logger.debug('Step: detecting link type', { url })
      const detectRes = await fetchWithAuth('/api/detect', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })
      const detectData = await detectRes.json()
      if (detectData.error) {
        logger.warn('Detect step returned an error', { url, error: detectData.error })
        setResult({ link_info: { url, link_type: 'unknown' }, classification: { primary: '기술소개', confidence: 0 }, summary: '', error: detectData.error })
        setStep('error')
        return
      }

      // Step 2: Extract
      setStep('extracting')
      logger.debug('Step: extracting content', { url })
      const extractRes = await fetchWithAuth('/api/extract', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })
      const extractData = await extractRes.json()
      if (extractData.error && !extractData.content) {
        logger.warn('Extract step returned an error', { url, error: extractData.error })
        setResult({ link_info: extractData.link_info, classification: { primary: '기술소개', confidence: 0 }, summary: '', error: extractData.error })
        setStep('error')
        return
      }

      // Step 3: Classify
      setStep('classifying')
      logger.debug('Step: classifying content', { url })
      const classifyRes = await fetchWithAuth('/api/classify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: extractData.content, provider }),
      })
      const classifyData = await classifyRes.json()

      // Step 4: Summarize
      setStep('summarizing')
      logger.debug('Step: generating summary', { url })
      const summarizeRes = await fetchWithAuth('/api/summarize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          content: extractData.content,
          category: classifyData.classification?.primary,
          provider,
        }),
      })
      const summarizeData = await summarizeRes.json()

      setResult({
        link_info: extractData.link_info,
        classification: classifyData.classification || { primary: '기술소개', confidence: 0 },
        summary: summarizeData.result?.summary || summarizeData.error || 'No summary generated',
      })
      setStep('done')
      logger.info('Summarization complete', { url })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'An unexpected error occurred'
      logger.error('Summarization failed with exception', { url, error: errorMessage })
      setResult({
        link_info: { url, link_type: 'unknown' },
        classification: { primary: '기술소개', confidence: 0 },
        summary: '',
        error: errorMessage,
      })
      setStep('error')
    }
  }

  const isLoading = step !== 'done' && step !== 'error'

  return (
    <div
      style={{
        maxWidth: '720px',
        margin: '0 auto',
        padding: '2rem 1rem',
        fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ fontSize: '1.75rem', marginBottom: '0.5rem' }}>Link Summarizer</h1>
        <button
          type="button"
          onClick={handleLogout}
          style={{
            background: 'none',
            border: '1px solid #e2e8f0',
            borderRadius: '6px',
            padding: '0.375rem 0.75rem',
            fontSize: '0.875rem',
            color: '#718096',
            cursor: 'pointer',
          }}
        >
          Log Out
        </button>
      </div>
      <p style={{ color: '#718096', marginBottom: '2rem' }}>
        Paste a link and get an optimized summary based on content type.
      </p>

      <UrlInput onSubmit={handleSubmit} isLoading={isLoading} />

      {isLoading && <ProgressIndicator currentStep={step} />}

      {result && <SummaryResult result={result} />}
    </div>
  )
}
