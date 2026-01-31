const path = require('path');

module.exports = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.js'],
  
  // Módulos que se deben transformar
  transform: {
    '^.+\\.(js|jsx|ts|tsx)$': ['babel-jest', {
      presets: [
        ['@babel/preset-env', { targets: { node: 'current' } }],
        ['@babel/preset-react', { runtime: 'automatic' }]
      ]
    }]
  },

  // Archivos que no se deben transformar (CORREGIDO)
  transformIgnorePatterns: [
    'node_modules/(?!(axios|recharts|d3-|@react-hook/|react-hot-toast)/)'
  ],

  // Module mapping para assets estáticos
  moduleNameMapping: {
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$': '<rootDir>/src/__mocks__/fileMock.js'
  },

  // Path mapping según configuración del proyecto
  moduleDirectories: ['node_modules', '<rootDir>/src'],

  // Coverage configuration
  collectCoverage: true,
  collectCoverageFrom: [
    'src/**/*.{js,jsx}',
    '!src/index.js',
    '!src/reportWebVitals.js',
    '!src/setupTests.js',
    '!src/**/*.test.{js,jsx}',
    '!src/**/__tests__/**',
    '!src/**/__mocks__/**'
  ],
  
  coveragePathIgnorePatterns: [
    '/node_modules/',
    '/build/',
    '/public/'
  ],

  coverageReporters: [
    'text',
    'lcov',
    'html',
    'json-summary'
  ],

  coverageThreshold: {
    global: {
      branches: 50,
      functions: 50,
      lines: 50,
      statements: 50
    }
  },

  // Test patterns
  testMatch: [
    '<rootDir>/src/**/__tests__/**/*.{js,jsx}',
    '<rootDir>/src/**/*.{test,spec}.{js,jsx}'
  ],

  // Performance settings
  maxWorkers: '50%',
  testTimeout: 10000,

  // Verbose output para debugging
  verbose: false,

  // Configuraciones adicionales
  clearMocks: true,
  restoreMocks: true,
  resetMocks: true,

  // Utilidades globales para tests
  globals: {
    'process.env.NODE_ENV': 'test'
  }
}; 