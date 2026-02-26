import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { LoginForm } from './LoginForm'

const mockFetch = vi.fn()

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch)
})

afterEach(() => {
  vi.restoreAllMocks()
})

describe('LoginForm', () => {
  it('renders login form by default', () => {
    render(<LoginForm onLogin={vi.fn()} />)

    expect(screen.getByRole('heading', { name: 'Log In' })).toBeInTheDocument()
    expect(screen.getByLabelText('Email')).toBeInTheDocument()
    expect(screen.getByLabelText(/Password/)).toBeInTheDocument()
  })

  it('switches to signup mode when Sign Up is clicked', () => {
    render(<LoginForm onLogin={vi.fn()} />)

    fireEvent.click(screen.getByText('Sign Up'))

    expect(screen.getByRole('heading', { name: 'Create Account' })).toBeInTheDocument()
  })

  it('calls onLogin with token on successful login', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ token: 'test-jwt-token' }),
    })

    const onLogin = vi.fn()
    render(<LoginForm onLogin={onLogin} />)

    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'alice@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/Password/), {
      target: { value: 'password123' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Log In' }))

    await waitFor(() => {
      expect(onLogin).toHaveBeenCalledWith('test-jwt-token')
    })
  })

  it('shows error on invalid credentials', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: 'invalid email or password' }),
    })

    render(<LoginForm onLogin={vi.fn()} />)

    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'alice@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/Password/), {
      target: { value: 'wrongpassword' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Log In' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('invalid email or password')
    })
  })

  it('performs signup then auto-login', async () => {
    // First call: signup success
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ id: 1, email: 'alice@example.com' }),
    })
    // Second call: login success
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ token: 'new-user-token' }),
    })

    const onLogin = vi.fn()
    render(<LoginForm onLogin={onLogin} />)

    // Switch to signup
    fireEvent.click(screen.getByText('Sign Up'))

    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'alice@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/Password/), {
      target: { value: 'password123' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))

    await waitFor(() => {
      expect(onLogin).toHaveBeenCalledWith('new-user-token')
    })

    // Verify both signup and login were called
    expect(mockFetch).toHaveBeenCalledTimes(2)
  })

  it('shows error on signup with duplicate email', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: 'email already registered' }),
    })

    render(<LoginForm onLogin={vi.fn()} />)

    fireEvent.click(screen.getByText('Sign Up'))
    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'existing@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/Password/), {
      target: { value: 'password123' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Sign Up' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('email already registered')
    })
  })

  it('shows network error on fetch failure', async () => {
    mockFetch.mockRejectedValueOnce(new Error('network error'))

    render(<LoginForm onLogin={vi.fn()} />)

    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'alice@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/Password/), {
      target: { value: 'password123' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Log In' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('Network error')
    })
  })

  it('clears error when switching between login and signup', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: 'invalid email or password' }),
    })

    render(<LoginForm onLogin={vi.fn()} />)

    fireEvent.change(screen.getByLabelText('Email'), {
      target: { value: 'alice@example.com' },
    })
    fireEvent.change(screen.getByLabelText(/Password/), {
      target: { value: 'wrongpassword' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Log In' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument()
    })

    // Switch to signup - error should clear
    fireEvent.click(screen.getByText('Sign Up'))

    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })
})
