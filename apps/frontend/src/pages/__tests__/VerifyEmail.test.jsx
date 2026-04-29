/**
 * VerifyEmail page tests
 *
 * Covers the critical fix: verification must NOT fire automatically on page
 * load (which lets email security scanners consume the token before the user
 * clicks). The page must start in 'idle' state and only call the API when the
 * user explicitly clicks "Verificar mi cuenta".
 */

// --- Mocks (hoisted before imports) ---

jest.mock('../../config/environments', () => ({
  __esModule: true,
  default: { USERS_API_URL: 'http://localhost:8080/api/v1' },
}));

jest.mock('../../services/authService', () => ({
  __esModule: true,
  default: { resendVerification: jest.fn() },
}));

jest.mock('../../components/Logo', () => ({
  __esModule: true,
  default: () => <div data-testid="logo" />,
}));

// --- Imports ---

import React from 'react';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import VerifyEmail from '../VerifyEmail';
import authService from '../../services/authService';

// Helper: render with a given ?token= search param
const renderWithToken = (token) => {
  const search = token !== undefined ? `?token=${token}` : '';
  return render(
    <MemoryRouter initialEntries={[`/verify-email${search}`]}>
      <Routes>
        <Route path="/verify-email" element={<VerifyEmail />} />
        <Route path="/login" element={<div>Login page</div>} />
      </Routes>
    </MemoryRouter>
  );
};

// ---------------------------------------------------------------------------
// Security: no automatic API call on load
// ---------------------------------------------------------------------------

describe('VerifyEmail — no automatic verification on load', () => {
  beforeEach(() => {
    global.fetch = jest.fn();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('does NOT call fetch automatically when the page loads with a valid token', () => {
    renderWithToken('some-valid-token');
    expect(global.fetch).not.toHaveBeenCalled();
  });

  it('shows "idle" state with a verify button — not a spinner or error', () => {
    renderWithToken('some-valid-token');

    expect(screen.getByText(/verificar mi cuenta/i)).toBeInTheDocument();
    expect(screen.queryByText(/verificando/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/enlace inválido/i)).not.toBeInTheDocument();
  });
});

// ---------------------------------------------------------------------------
// Missing token
// ---------------------------------------------------------------------------

describe('VerifyEmail — missing token', () => {
  it('shows error state immediately when there is no token in the URL', () => {
    renderWithToken(undefined);

    expect(screen.getByText(/enlace inválido/i)).toBeInTheDocument();
    expect(screen.queryByText(/verificar mi cuenta/i)).not.toBeInTheDocument();
  });
});

// ---------------------------------------------------------------------------
// Successful verification flow
// ---------------------------------------------------------------------------

describe('VerifyEmail — successful verification', () => {
  beforeEach(() => {
    global.fetch = jest.fn().mockResolvedValue({ ok: true });
  });

  afterEach(() => jest.restoreAllMocks());

  it('calls fetch with the correct endpoint when user clicks the verify button', async () => {
    renderWithToken('my-token-abc');

    fireEvent.click(screen.getByText(/verificar mi cuenta/i));

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/auth/verify-email/my-token-abc'
      );
    });
  });

  it('shows success state after a successful API response', async () => {
    renderWithToken('my-token-abc');

    fireEvent.click(screen.getByText(/verificar mi cuenta/i));

    await waitFor(() => {
      expect(screen.getByText(/¡cuenta verificada!/i)).toBeInTheDocument();
    });

    expect(screen.getByRole('link', { name: /iniciar sesión/i })).toBeInTheDocument();
  });

  it('shows loading spinner while the request is in flight', async () => {
    // Never resolves — keeps the request pending
    global.fetch = jest.fn(() => new Promise(() => {}));

    renderWithToken('my-token-abc');
    fireEvent.click(screen.getByText(/verificar mi cuenta/i));

    expect(screen.getByText(/verificando tu cuenta/i)).toBeInTheDocument();
  });
});

// ---------------------------------------------------------------------------
// Error flow (expired / invalid token)
// ---------------------------------------------------------------------------

describe('VerifyEmail — error flow', () => {
  afterEach(() => jest.restoreAllMocks());

  it('shows error state when API returns non-ok for an expired token', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: 'Verification token has expired' }),
    });

    renderWithToken('expired-token');
    fireEvent.click(screen.getByText(/verificar mi cuenta/i));

    await waitFor(() => {
      expect(screen.getByText(/enlace inválido/i)).toBeInTheDocument();
      expect(screen.getByText(/verification token has expired/i)).toBeInTheDocument();
    });
  });

  it('shows the resend form in the error state', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: 'Invalid verification token' }),
    });

    renderWithToken('bad-token');
    fireEvent.click(screen.getByText(/verificar mi cuenta/i));

    await waitFor(() => {
      expect(screen.getByPlaceholderText(/tu@email\.com/i)).toBeInTheDocument();
      expect(screen.getByText(/reenviar email de verificación/i)).toBeInTheDocument();
    });
  });
});

// ---------------------------------------------------------------------------
// Resend flow
// ---------------------------------------------------------------------------

describe('VerifyEmail — resend verification', () => {
  afterEach(() => jest.restoreAllMocks());

  const showErrorState = async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      json: () => Promise.resolve({ error: 'Invalid verification token' }),
    });
    renderWithToken('bad-token');
    fireEvent.click(screen.getByText(/verificar mi cuenta/i));
    await waitFor(() => screen.getByPlaceholderText(/tu@email\.com/i));
  };

  it('calls authService.resendVerification with the entered email', async () => {
    authService.resendVerification.mockResolvedValue({ success: true });
    await showErrorState();

    fireEvent.change(screen.getByPlaceholderText(/tu@email\.com/i), {
      target: { value: 'user@example.com' },
    });
    fireEvent.click(screen.getByText(/reenviar email de verificación/i));

    await waitFor(() => {
      expect(authService.resendVerification).toHaveBeenCalledWith('user@example.com');
    });
  });

  it('resend button is disabled when email input is empty', async () => {
    await showErrorState();

    const resendBtn = screen.getByText(/reenviar email de verificación/i);
    expect(resendBtn).toBeDisabled();
  });
});
