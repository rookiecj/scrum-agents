import { test, expect, Page } from '@playwright/test'

// Mock API response helpers
const mockDetectResponse = {
  link_info: {
    url: 'https://example.com/article',
    link_type: 'article',
    title: '',
  },
}

const mockExtractResponse = {
  link_info: {
    url: 'https://example.com/article',
    link_type: 'article',
    title: 'Test Article Title',
    author: 'John Doe',
  },
  content: 'This is the extracted article content for testing.',
}

const mockClassifyResponse = {
  classification: {
    primary: '기술소개',
    confidence: 0.92,
  },
}

const mockSummarizeResponse = {
  summary: 'This is a summarized version of the test article content.',
}

// Set up API mocks that respond after a short delay to allow progress steps to be visible
async function setupAPIMocks(page: Page, overrides?: {
  detect?: object | null
  extract?: object | null
  classify?: object | null
  summarize?: object | null
}) {
  await page.route('**/api/detect', async (route) => {
    if (overrides?.detect === null) {
      await route.abort()
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(overrides?.detect ?? mockDetectResponse),
    })
  })

  await page.route('**/api/extract', async (route) => {
    if (overrides?.extract === null) {
      await route.abort()
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(overrides?.extract ?? mockExtractResponse),
    })
  })

  await page.route('**/api/classify', async (route) => {
    if (overrides?.classify === null) {
      await route.abort()
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(overrides?.classify ?? mockClassifyResponse),
    })
  })

  await page.route('**/api/summarize', async (route) => {
    if (overrides?.summarize === null) {
      await route.abort()
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(overrides?.summarize ?? mockSummarizeResponse),
    })
  })
}

test.describe('Link Summarizer App', () => {

  test('shows the page title and URL input form', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole('heading', { name: 'Link Summarizer' })).toBeVisible()
    await expect(page.getByLabel('URL')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Summarize' })).toBeVisible()
    await expect(page.getByLabel('Claude')).toBeChecked()
  })

  test('submit button is disabled when URL is empty', async ({ page }) => {
    await page.goto('/')

    const button = page.getByRole('button', { name: 'Summarize' })
    await expect(button).toBeDisabled()
  })

  test('full happy path: submit URL and see results', async ({ page }) => {
    await setupAPIMocks(page)
    await page.goto('/')

    // Enter URL
    await page.getByLabel('URL').fill('https://example.com/article')

    // Submit
    await page.getByRole('button', { name: 'Summarize' }).click()

    // Wait for result to appear
    await expect(page.getByText('Test Article Title')).toBeVisible({ timeout: 10000 })

    // Verify classification badge
    await expect(page.getByText('기술소개')).toBeVisible()

    // Verify link type badge
    await expect(page.getByText('article', { exact: true })).toBeVisible()

    // Verify confidence
    await expect(page.getByText('Confidence: 92%')).toBeVisible()

    // Verify summary text
    await expect(page.getByText('This is a summarized version')).toBeVisible()

    // Verify author
    await expect(page.getByText('Author: John Doe')).toBeVisible()
  })

  test('provider selection: OpenAI sends correct provider', async ({ page }) => {
    let capturedClassifyBody: string = ''
    let capturedSummarizeBody: string = ''

    await page.route('**/api/detect', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockDetectResponse),
      })
    })

    await page.route('**/api/extract', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockExtractResponse),
      })
    })

    await page.route('**/api/classify', async (route) => {
      capturedClassifyBody = route.request().postData() || ''
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockClassifyResponse),
      })
    })

    await page.route('**/api/summarize', async (route) => {
      capturedSummarizeBody = route.request().postData() || ''
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockSummarizeResponse),
      })
    })

    await page.goto('/')

    // Select OpenAI
    await page.getByLabel('OpenAI').click()

    // Enter URL and submit
    await page.getByLabel('URL').fill('https://example.com/article')
    await page.getByRole('button', { name: 'Summarize' }).click()

    // Wait for completion
    await expect(page.getByText('Test Article Title')).toBeVisible({ timeout: 10000 })

    // Verify provider was sent in classify request
    expect(capturedClassifyBody).toContain('"provider":"openai"')

    // Verify provider was sent in summarize request
    expect(capturedSummarizeBody).toContain('"provider":"openai"')
  })

  test('shows error when detect API returns error', async ({ page }) => {
    await setupAPIMocks(page, {
      detect: { link_info: { url: '', link_type: 'unknown' }, error: 'invalid URL: parse error' },
    })

    await page.goto('/')
    await page.getByLabel('URL').fill('https://example.com/article')
    await page.getByRole('button', { name: 'Summarize' }).click()

    // Wait for error to appear
    await expect(page.getByRole('alert')).toBeVisible({ timeout: 10000 })
    await expect(page.getByText('invalid URL: parse error')).toBeVisible()
  })

  test('shows error when extract API fails', async ({ page }) => {
    await setupAPIMocks(page, {
      extract: { link_info: { url: 'https://example.com/article', link_type: 'article' }, error: 'extraction failed' },
    })

    await page.goto('/')
    await page.getByLabel('URL').fill('https://example.com/article')
    await page.getByRole('button', { name: 'Summarize' }).click()

    await expect(page.getByRole('alert')).toBeVisible({ timeout: 10000 })
    await expect(page.getByText('extraction failed')).toBeVisible()
  })

  test('shows error on network failure', async ({ page }) => {
    // Abort all API requests to simulate network failure
    await setupAPIMocks(page, {
      detect: null,
    })

    await page.goto('/')
    await page.getByLabel('URL').fill('https://example.com/article')
    await page.getByRole('button', { name: 'Summarize' }).click()

    await expect(page.getByRole('alert')).toBeVisible({ timeout: 10000 })
  })

  test('progress steps are shown during workflow', async ({ page }) => {
    // Use delayed responses so progress steps are visible
    await page.route('**/api/detect', async (route) => {
      await new Promise(r => setTimeout(r, 500))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockDetectResponse),
      })
    })

    await page.route('**/api/extract', async (route) => {
      await new Promise(r => setTimeout(r, 500))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockExtractResponse),
      })
    })

    await page.route('**/api/classify', async (route) => {
      await new Promise(r => setTimeout(r, 500))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockClassifyResponse),
      })
    })

    await page.route('**/api/summarize', async (route) => {
      await new Promise(r => setTimeout(r, 500))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockSummarizeResponse),
      })
    })

    await page.goto('/')
    await page.getByLabel('URL').fill('https://example.com/article')
    await page.getByRole('button', { name: 'Summarize' }).click()

    // Verify progress indicator appears with step labels
    await expect(page.getByText('Detecting link type')).toBeVisible({ timeout: 5000 })

    // Wait for final result
    await expect(page.getByText('Test Article Title')).toBeVisible({ timeout: 15000 })
  })

  test('button shows Processing state while loading', async ({ page }) => {
    await page.route('**/api/detect', async (route) => {
      await new Promise(r => setTimeout(r, 1000))
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(mockDetectResponse),
      })
    })
    // Set up other routes so they don't hang
    await page.route('**/api/extract', async (route) => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(mockExtractResponse) })
    })
    await page.route('**/api/classify', async (route) => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(mockClassifyResponse) })
    })
    await page.route('**/api/summarize', async (route) => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(mockSummarizeResponse) })
    })

    await page.goto('/')
    await page.getByLabel('URL').fill('https://example.com/article')
    await page.getByRole('button', { name: 'Summarize' }).click()

    // Button text changes during loading
    await expect(page.getByRole('button', { name: 'Processing...' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Processing...' })).toBeDisabled()
  })
})
