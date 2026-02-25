/**
 * Structured logging utility for the frontend application.
 *
 * Provides debug/info/warn/error level methods that output structured data
 * (timestamp, level, message, context) to the browser console.
 *
 * Log level filtering is controlled by NODE_ENV:
 *   - development: all levels (debug and above)
 *   - production:  warn and above (debug and info are suppressed)
 */

export type LogLevel = 'debug' | 'info' | 'warn' | 'error'

export interface LogEntry {
  timestamp: string
  level: LogLevel
  message: string
  context?: Record<string, unknown>
}

/** Numeric priority for each log level (lower = more verbose). */
const LOG_LEVEL_PRIORITY: Record<LogLevel, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
}

/**
 * Return the minimum log level based on the current environment.
 * In production mode debug and info logs are suppressed.
 */
function getMinLevel(): LogLevel {
  if (typeof import.meta !== 'undefined' && import.meta.env?.MODE === 'production') {
    return 'warn'
  }
  // Default to 'debug' in development / test
  return 'debug'
}

export interface Logger {
  debug(message: string, context?: Record<string, unknown>): void
  info(message: string, context?: Record<string, unknown>): void
  warn(message: string, context?: Record<string, unknown>): void
  error(message: string, context?: Record<string, unknown>): void
}

/**
 * Create a structured log entry object.
 */
function buildEntry(level: LogLevel, message: string, context?: Record<string, unknown>): LogEntry {
  const entry: LogEntry = {
    timestamp: new Date().toISOString(),
    level,
    message,
  }
  if (context !== undefined) {
    entry.context = context
  }
  return entry
}

/**
 * Determine whether the given level should be emitted based on the minimum level.
 */
function shouldLog(level: LogLevel, minLevel: LogLevel): boolean {
  return LOG_LEVEL_PRIORITY[level] >= LOG_LEVEL_PRIORITY[minLevel]
}

/**
 * Write a log entry to the appropriate console method.
 */
function emit(entry: LogEntry): void {
  const consoleMethods: Record<LogLevel, (...args: unknown[]) => void> = {
    debug: console.debug,
    info: console.info,
    warn: console.warn,
    error: console.error,
  }
  const fn = consoleMethods[entry.level]
  fn(`[${entry.timestamp}] [${entry.level.toUpperCase()}] ${entry.message}`, entry.context ?? '')
}

/**
 * Create a Logger instance.
 *
 * @param overrideMinLevel - Override the environment-derived minimum level (useful for testing).
 */
export function createLogger(overrideMinLevel?: LogLevel): Logger {
  const minLevel = overrideMinLevel ?? getMinLevel()

  function log(level: LogLevel, message: string, context?: Record<string, unknown>): void {
    if (!shouldLog(level, minLevel)) {
      return
    }
    const entry = buildEntry(level, message, context)
    emit(entry)
  }

  return {
    debug: (message, context?) => log('debug', message, context),
    info: (message, context?) => log('info', message, context),
    warn: (message, context?) => log('warn', message, context),
    error: (message, context?) => log('error', message, context),
  }
}

/**
 * Default logger instance used throughout the application.
 *
 * The minimum level is derived from the environment automatically.
 */
export const logger: Logger = createLogger()
