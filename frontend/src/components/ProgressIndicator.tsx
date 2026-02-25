import type { SummarizeStep } from '../types/api'

interface ProgressIndicatorProps {
  currentStep: SummarizeStep
}

const steps: { key: SummarizeStep; label: string }[] = [
  { key: 'detecting', label: 'Detecting link type' },
  { key: 'extracting', label: 'Extracting content' },
  { key: 'classifying', label: 'Classifying content' },
  { key: 'summarizing', label: 'Generating summary' },
]

function getStepIndex(step: SummarizeStep): number {
  const idx = steps.findIndex((s) => s.key === step)
  return idx >= 0 ? idx : steps.length
}

export function ProgressIndicator({ currentStep }: ProgressIndicatorProps) {
  if (currentStep === 'done' || currentStep === 'error') return null

  const currentIndex = getStepIndex(currentStep)

  return (
    <div style={{ marginBottom: '1.5rem' }} role="progressbar" aria-label="Summarization progress">
      {steps.map((step, i) => {
        const isActive = i === currentIndex
        const isDone = i < currentIndex

        return (
          <div
            key={step.key}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '0.75rem',
              padding: '0.5rem 0',
              color: isActive ? '#3182ce' : isDone ? '#38a169' : '#a0aec0',
              fontWeight: isActive ? 'bold' : 'normal',
            }}
          >
            <span style={{ fontSize: '1.25rem' }}>
              {isDone ? '\u2713' : isActive ? '\u25CB' : '\u25CB'}
            </span>
            <span>{step.label}</span>
            {isActive && (
              <span style={{ fontSize: '0.75rem', color: '#718096' }}>...</span>
            )}
          </div>
        )
      })}
    </div>
  )
}
