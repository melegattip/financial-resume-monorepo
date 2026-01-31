/**
 * Production-safe logger utility
 * Automatically disables console logs in production environment
 */

import config from '../config/environments';

// Define log levels
const LOG_LEVELS = {
  ERROR: 0,
  WARN: 1,
  INFO: 2,
  DEBUG: 3
};

class Logger {
  constructor() {
    this.isProduction = config.IS_PRODUCTION;
    this.environment = config.ENVIRONMENT;
    
    // In production, only allow ERROR and WARN levels
    this.maxLevel = this.isProduction ? LOG_LEVELS.WARN : LOG_LEVELS.DEBUG;
    
    // Create bound methods to preserve context
    this.log = this.log.bind(this);
    this.error = this.error.bind(this);
    this.warn = this.warn.bind(this);
    this.info = this.info.bind(this);
    this.debug = this.debug.bind(this);
  }

  /**
   * Generic log method that checks environment and level
   */
  _shouldLog(level) {
    return level <= this.maxLevel;
  }

  /**
   * Format log message with timestamp and environment info
   */
  _formatMessage(level, message, ...args) {
    const timestamp = new Date().toISOString();
    const levelName = Object.keys(LOG_LEVELS)[level];
    
    if (this.isProduction) {
      // In production, minimal formatting
      return [message, ...args];
    } else {
      // In development, include more context
      return [`[${timestamp}] [${this.environment.toUpperCase()}] [${levelName}]`, message, ...args];
    }
  }

  /**
   * Error logging - always shown
   */
  error(message, ...args) {
    if (this._shouldLog(LOG_LEVELS.ERROR)) {
      console.error(...this._formatMessage(LOG_LEVELS.ERROR, message, ...args));
    }
  }

  /**
   * Warning logging - shown in production and development
   */
  warn(message, ...args) {
    if (this._shouldLog(LOG_LEVELS.WARN)) {
      console.warn(...this._formatMessage(LOG_LEVELS.WARN, message, ...args));
    }
  }

  /**
   * Info logging - only in development
   */
  info(message, ...args) {
    if (this._shouldLog(LOG_LEVELS.INFO)) {
      console.info(...this._formatMessage(LOG_LEVELS.INFO, message, ...args));
    }
  }

  /**
   * Debug logging - only in development
   */
  debug(message, ...args) {
    if (this._shouldLog(LOG_LEVELS.DEBUG)) {
      console.log(...this._formatMessage(LOG_LEVELS.DEBUG, message, ...args));
    }
  }

  /**
   * General log method - maps to debug level
   */
  log(message, ...args) {
    this.debug(message, ...args);
  }

  /**
   * Group logging for better organization
   */
  group(label) {
    if (!this.isProduction && console.group) {
      console.group(label);
    }
  }

  groupEnd() {
    if (!this.isProduction && console.groupEnd) {
      console.groupEnd();
    }
  }

  /**
   * Table logging for development
   */
  table(data) {
    if (!this.isProduction && console.table) {
      console.table(data);
    }
  }

  /**
   * Time logging for performance monitoring
   */
  time(label) {
    if (!this.isProduction && console.time) {
      console.time(label);
    }
  }

  timeEnd(label) {
    if (!this.isProduction && console.timeEnd) {
      console.timeEnd(label);
    }
  }
}

// Create and export singleton instance
const logger = new Logger();

// Override global console in production
if (logger.isProduction) {
  // Preserve original console methods for internal use
  const originalConsole = { ...console };
  
  // Override console methods to use our logger
  window.console = {
    ...originalConsole,
    log: logger.debug,
    info: logger.info,
    warn: logger.warn,
    error: logger.error,
    debug: logger.debug,
    group: logger.group,
    groupEnd: logger.groupEnd,
    table: logger.table,
    time: logger.time,
    timeEnd: logger.timeEnd,
    // Keep these methods unchanged
    clear: originalConsole.clear,
    count: originalConsole.count,
    countReset: originalConsole.countReset,
    dir: originalConsole.dir,
    dirxml: originalConsole.dirxml,
    trace: originalConsole.trace
  };
}

export default logger;

// Named exports for convenience
export const log = logger.log;
export const error = logger.error;
export const warn = logger.warn;
export const info = logger.info;
export const debug = logger.debug;
