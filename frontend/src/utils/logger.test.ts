import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { createLogger, type LogLevel } from './logger'

describe('logger', () => {
  let debugSpy: ReturnType<typeof vi.spyOn>
  let infoSpy: ReturnType<typeof vi.spyOn>
  let warnSpy: ReturnType<typeof vi.spyOn>
  let errorSpy: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    debugSpy = vi.spyOn(console, 'debug').mockImplementation(() => {})
    infoSpy = vi.spyOn(console, 'info').mockImplementation(() => {})
    warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})
    errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('createLogger with debug level', () => {
    it('outputs debug messages', () => {
      const log = createLogger('debug')
      log.debug('test debug')
      expect(debugSpy).toHaveBeenCalledTimes(1)
    })

    it('outputs info messages', () => {
      const log = createLogger('debug')
      log.info('test info')
      expect(infoSpy).toHaveBeenCalledTimes(1)
    })

    it('outputs warn messages', () => {
      const log = createLogger('debug')
      log.warn('test warn')
      expect(warnSpy).toHaveBeenCalledTimes(1)
    })

    it('outputs error messages', () => {
      const log = createLogger('debug')
      log.error('test error')
      expect(errorSpy).toHaveBeenCalledTimes(1)
    })
  })

  describe('createLogger with warn level (production-like)', () => {
    it('suppresses debug messages', () => {
      const log = createLogger('warn')
      log.debug('should not appear')
      expect(debugSpy).not.toHaveBeenCalled()
    })

    it('suppresses info messages', () => {
      const log = createLogger('warn')
      log.info('should not appear')
      expect(infoSpy).not.toHaveBeenCalled()
    })

    it('outputs warn messages', () => {
      const log = createLogger('warn')
      log.warn('should appear')
      expect(warnSpy).toHaveBeenCalledTimes(1)
    })

    it('outputs error messages', () => {
      const log = createLogger('warn')
      log.error('should appear')
      expect(errorSpy).toHaveBeenCalledTimes(1)
    })
  })

  describe('createLogger with error level', () => {
    it('suppresses debug messages', () => {
      const log = createLogger('error')
      log.debug('no')
      expect(debugSpy).not.toHaveBeenCalled()
    })

    it('suppresses info messages', () => {
      const log = createLogger('error')
      log.info('no')
      expect(infoSpy).not.toHaveBeenCalled()
    })

    it('suppresses warn messages', () => {
      const log = createLogger('error')
      log.warn('no')
      expect(warnSpy).not.toHaveBeenCalled()
    })

    it('outputs error messages', () => {
      const log = createLogger('error')
      log.error('yes')
      expect(errorSpy).toHaveBeenCalledTimes(1)
    })
  })

  describe('structured output', () => {
    it('includes timestamp, level, and message in the log output', () => {
      const log = createLogger('debug')
      log.info('hello world')

      expect(infoSpy).toHaveBeenCalledTimes(1)
      const firstArg = infoSpy.mock.calls[0][0] as string
      // Format: [<ISO timestamp>] [INFO] hello world
      expect(firstArg).toMatch(/^\[.+\] \[INFO\] hello world$/)
    })

    it('includes context as the second argument', () => {
      const log = createLogger('debug')
      const ctx = { url: '/api/test', status: 500 }
      log.error('request failed', ctx)

      expect(errorSpy).toHaveBeenCalledTimes(1)
      const secondArg = errorSpy.mock.calls[0][1] as Record<string, unknown>
      expect(secondArg).toEqual(ctx)
    })

    it('passes empty string when no context is provided', () => {
      const log = createLogger('debug')
      log.info('no context')

      expect(infoSpy).toHaveBeenCalledTimes(1)
      expect(infoSpy.mock.calls[0][1]).toBe('')
    })

    it('includes a valid ISO timestamp', () => {
      const log = createLogger('debug')
      log.debug('timestamp check')

      const firstArg = debugSpy.mock.calls[0][0] as string
      // Extract the timestamp from [<timestamp>] [DEBUG] ...
      const match = firstArg.match(/^\[(.+?)\]/)
      expect(match).not.toBeNull()
      const ts = match![1]
      expect(new Date(ts).toISOString()).toBe(ts)
    })
  })

  describe('level labels in output', () => {
    const levels: LogLevel[] = ['debug', 'info', 'warn', 'error']
    const spyMap = () => ({
      debug: debugSpy,
      info: infoSpy,
      warn: warnSpy,
      error: errorSpy,
    })

    it.each(levels)('outputs uppercase level label for %s', (level) => {
      const log = createLogger('debug')
      log[level]('msg')
      const spy = spyMap()[level]
      const firstArg = spy.mock.calls[0][0] as string
      expect(firstArg).toContain(`[${level.toUpperCase()}]`)
    })
  })
})
