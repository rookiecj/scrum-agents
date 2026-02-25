import { useState } from 'react'
import { UrlInput } from './components/UrlInput'
import { ProgressIndicator } from './components/ProgressIndicator'
import { SummaryResult } from './components/SummaryResult'
import type { SummarizeResponse, SummarizeStep } from './types/api'

export function App() {
  const [step, setStep] = useState<SummarizeStep>('done')
  const [result, setResult] = useState<SummarizeResponse | null>(null)

  const handleSubmit = async (url: string, provider: 'claude' | 'openai') => {
    setResult(null)

    try {
      // Step 1: Detect
      setStep('detecting')
      const detectRes = await fetch('/api/detect', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })
      const detectData = await detectRes.json()
      if (detectData.error) {
        setResult({ link_info: { url, link_type: 'unknown' }, classification: { primary: '기술소개', confidence: 0 }, summary: '', error: detectData.error })
        setStep('error')
        return
      }

      // Step 2: Extract
      setStep('extracting')
      const extractRes = await fetch('/api/extract', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      })
      const extractData = await extractRes.json()
      if (extractData.error && !extractData.content) {
        setResult({ link_info: extractData.link_info, classification: { primary: '기술소개', confidence: 0 }, summary: '', error: extractData.error })
        setStep('error')
        return
      }

      // Step 3: Classify
      setStep('classifying')
      const classifyRes = await fetch('/api/classify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: extractData.content, provider }),
      })
      const classifyData = await classifyRes.json()

      // Step 4: Summarize
      setStep('summarizing')
      const summarizeRes = await fetch('/api/summarize', {
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
        summary: summarizeData.summary || summarizeData.error || 'No summary generated',
      })
      setStep('done')
    } catch (err) {
      setResult({
        link_info: { url, link_type: 'unknown' },
        classification: { primary: '기술소개', confidence: 0 },
        summary: '',
        error: err instanceof Error ? err.message : 'An unexpected error occurred',
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
      <h1 style={{ fontSize: '1.75rem', marginBottom: '0.5rem' }}>Link Summarizer</h1>
      <p style={{ color: '#718096', marginBottom: '2rem' }}>
        Paste a link and get an optimized summary based on content type.
      </p>

      <UrlInput onSubmit={handleSubmit} isLoading={isLoading} />

      {isLoading && <ProgressIndicator currentStep={step} />}

      {result && <SummaryResult result={result} />}
    </div>
  )
}
