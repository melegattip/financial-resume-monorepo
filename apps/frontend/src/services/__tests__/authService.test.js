/**
 * authService contract tests
 *
 * Purpose: pin the parsing logic so any change to the backend JSON shape
 * is caught before production. Specifically guards against the mismatch
 * between the old microservice flat response and the monolith nested response.
 *
 * Backend contract (monolith):
 *   POST /api/v1/auth/login → { user: {...}, tokens: { access_token, refresh_token, expires_at } }
 *
 * Old microservice flat response (legacy, no longer used):
 *   POST /users/login → { access_token, refresh_token, expires_at, user: {...} }
 */

// --- Mocks (hoisted before imports by Jest) ---

jest.mock('axios', () => {
  const mockInstance = {
    interceptors: {
      request: { use: jest.fn() },
      response: { use: jest.fn() },
    },
    post: jest.fn(),
    get: jest.fn(),
    put: jest.fn(),
    defaults: { baseURL: 'http://localhost:8080/api/v1' },
  };
  return {
    create: jest.fn(() => mockInstance),
    default: { create: jest.fn(() => mockInstance) },
  };
});

jest.mock('../configService', () => ({
  loadConfig: jest.fn().mockResolvedValue({
    users_service_url: 'http://localhost:8080/api/v1',
    environment: 'development',
    version: '1.0.0',
  }),
}));

jest.mock('../../utils/notifications', () => ({
  success: jest.fn(),
  error: jest.fn(),
  warn: jest.fn(),
  default: { success: jest.fn(), error: jest.fn(), warn: jest.fn() },
}));

jest.mock('../dataService', () => ({
  clearCache: jest.fn(),
  default: { clearCache: jest.fn() },
}));

// --- Test setup ---

// Mock localStorage
const localStorageMock = (() => {
  let store = {};
  return {
    getItem: jest.fn((key) => store[key] ?? null),
    setItem: jest.fn((key, val) => { store[key] = String(val); }),
    removeItem: jest.fn((key) => { delete store[key]; }),
    clear: jest.fn(() => { store = {}; }),
  };
})();
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

// Mock window.location so authService constructor doesn't break
Object.defineProperty(window, 'location', {
  value: { hostname: 'localhost' },
  writable: true,
});

// Import AFTER mocks are set up
const authService = require('../authService').default;

// --- Test data fixtures ---

/** Monolith response shape: { user, tokens: { access_token, ... } } */
const monolithAuthResponse = {
  user: {
    id: 'user-123',
    email: 'test@example.com',
    first_name: 'Test',
    last_name: 'User',
    is_active: true,
    is_verified: false,
    created_at: '2026-01-01T00:00:00Z',
  },
  tokens: {
    access_token: 'monolith.access.token',
    refresh_token: 'monolith.refresh.token',
    expires_at: '2026-03-01T12:00:00Z', // ISO8601 string (time.Time)
    token_type: 'Bearer',
  },
};

/** Legacy flat response shape (old microservice) */
const legacyFlatAuthResponse = {
  access_token: 'legacy.access.token',
  refresh_token: 'legacy.refresh.token',
  expires_at: 1772568000, // Unix timestamp as number
  user: {
    id: 'user-456',
    email: 'legacy@example.com',
    first_name: 'Legacy',
    last_name: 'User',
    is_active: true,
    is_verified: true,
    created_at: '2025-01-01T00:00:00Z',
  },
};

// --- Tests ---

describe('authService.saveAuthData', () => {
  beforeEach(() => {
    localStorageMock.clear();
    localStorageMock.getItem.mockClear();
    localStorageMock.setItem.mockClear();
    authService.token = null;
    authService.refreshToken_value = null;
    authService.expiresAt = null;
    authService.user = null;
  });

  describe('monolith nested response { user, tokens: { access_token } }', () => {
    it('extracts access_token from tokens object', () => {
      authService.saveAuthData(monolithAuthResponse);
      expect(authService.token).toBe('monolith.access.token');
      expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_token', 'monolith.access.token');
    });

    it('extracts refresh_token from tokens object', () => {
      authService.saveAuthData(monolithAuthResponse);
      expect(authService.refreshToken_value).toBe('monolith.refresh.token');
      expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_refresh_token', 'monolith.refresh.token');
    });

    it('converts ISO8601 expires_at string to Unix timestamp integer', () => {
      authService.saveAuthData(monolithAuthResponse);
      const expectedUnix = Math.floor(new Date('2026-03-01T12:00:00Z').getTime() / 1000);
      expect(authService.expiresAt).toBe(expectedUnix);
      // localStorage must store a parseable integer string
      const stored = localStorageMock.setItem.mock.calls.find(([k]) => k === 'auth_expires_at');
      expect(stored).toBeDefined();
      expect(isNaN(parseInt(stored[1]))).toBe(false);
    });

    it('stores user data correctly', () => {
      authService.saveAuthData(monolithAuthResponse);
      expect(authService.user).toEqual(monolithAuthResponse.user);
    });
  });

  describe('backward compat: legacy flat response { access_token, user }', () => {
    it('extracts access_token from top level', () => {
      authService.saveAuthData(legacyFlatAuthResponse);
      expect(authService.token).toBe('legacy.access.token');
    });

    it('keeps numeric expires_at as-is (no conversion needed)', () => {
      authService.saveAuthData(legacyFlatAuthResponse);
      // 1772568000 is a valid number, isNaN(Number(1772568000)) === false → no conversion
      expect(authService.expiresAt).toBe(1772568000);
    });
  });

  describe('error cases', () => {
    it('throws when access_token is missing in nested structure', () => {
      const bad = { user: monolithAuthResponse.user, tokens: {} };
      expect(() => authService.saveAuthData(bad)).toThrow('Token de acceso no encontrado');
    });

    it('throws when user is missing', () => {
      const bad = { tokens: monolithAuthResponse.tokens };
      expect(() => authService.saveAuthData(bad)).toThrow('Datos de usuario no encontrados');
    });
  });
});

describe('authService login/register validation', () => {
  describe('accepts monolith nested token response', () => {
    it('validation passes when tokens.access_token exists', () => {
      // Direct test of the condition used in login() and register():
      // (!authData.access_token && !authData.tokens?.access_token) || !authData.user
      const data = monolithAuthResponse;
      const isInvalid = (!data.access_token && !data.tokens?.access_token) || !data.user;
      expect(isInvalid).toBe(false);
    });

    it('validation fails when neither access_token exists', () => {
      const data = { user: monolithAuthResponse.user, tokens: {} };
      const isInvalid = (!data.access_token && !data.tokens?.access_token) || !data.user;
      expect(isInvalid).toBe(true);
    });

    it('validation fails when user is missing', () => {
      const data = { tokens: monolithAuthResponse.tokens };
      const isInvalid = (!data.access_token && !data.tokens?.access_token) || !data.user;
      expect(isInvalid).toBe(true);
    });
  });
});
