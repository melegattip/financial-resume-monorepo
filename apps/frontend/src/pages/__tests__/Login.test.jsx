/**
 * Login page tests — EMAIL_NOT_VERIFIED flow
 *
 * Guards the resend-verification UX: when the backend returns EMAIL_NOT_VERIFIED
 * the login page must show a "Reenviar email" button, not just a static error.
 */

// --- Mocks ---

jest.mock('../../contexts/AuthContext', () => ({
  useAuth: jest.fn(),
}));

jest.mock('../../services/authService', () => ({
  __esModule: true,
  default: { resendVerification: jest.fn() },
}));

jest.mock('../../components/Logo', () => ({
  __esModule: true,
  default: () => <div data-testid="logo" />,
}));

jest.mock('../../components/TwoFAModal', () => ({
  __esModule: true,
  default: () => null,
}));

jest.mock('../../config/environments', () => ({
  __esModule: true,
  default: { USERS_API_URL: 'http://localhost:8080/api/v1' },
}));

jest.mock('../../utils/validation', () => ({
  validateEmail: () => ({ isValid: true }),
  sanitizeText: (v) => v,
}));

// --- Imports ---

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import Login from '../Login';
import authService from '../../services/authService';

// Helper
const renderLogin = () =>
  render(
    <MemoryRouter initialEntries={['/login']}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/dashboard" element={<div>Dashboard</div>} />
      </Routes>
    </MemoryRouter>
  );

const mockUseAuth = (overrides = {}) => {
  useAuth.mockReturnValue({
    login: jest.fn(),
    isAuthenticated: false,
    isLoading: false,
    isInitialized: true,
    ...overrides,
  });
};

// Simulate the check-2FA preflight returning false (no 2FA needed)
const mockNo2FA = () => {
  global.fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ requires_2fa: false }),
  });
};

// ---------------------------------------------------------------------------
// EMAIL_NOT_VERIFIED error
// ---------------------------------------------------------------------------

describe('Login — EMAIL_NOT_VERIFIED', () => {
  afterEach(() => jest.restoreAllMocks());

  it('shows error message when login throws EMAIL_NOT_VERIFIED', async () => {
    mockNo2FA();
    const loginFn = jest.fn().mockRejectedValue(new Error('EMAIL_NOT_VERIFIED'));
    mockUseAuth({ login: loginFn });

    renderLogin();

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(screen.getByText(/tu correo no ha sido verificado/i)).toBeInTheDocument();
    });
  });

  it('shows the resend button when login fails with EMAIL_NOT_VERIFIED', async () => {
    mockNo2FA();
    const loginFn = jest.fn().mockRejectedValue(new Error('EMAIL_NOT_VERIFIED'));
    mockUseAuth({ login: loginFn });

    renderLogin();

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(screen.getByText(/reenviar email de verificación/i)).toBeInTheDocument();
    });
  });

  it('does NOT show resend button for generic login errors', async () => {
    mockNo2FA();
    const loginFn = jest.fn().mockRejectedValue(new Error('invalid email or password'));
    mockUseAuth({ login: loginFn });

    renderLogin();

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'wrongpass' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(screen.queryByText(/reenviar email de verificación/i)).not.toBeInTheDocument();
    });
  });

  it('calls authService.resendVerification with the email from the form', async () => {
    mockNo2FA();
    authService.resendVerification.mockResolvedValue({ success: true });
    const loginFn = jest.fn().mockRejectedValue(new Error('EMAIL_NOT_VERIFIED'));
    mockUseAuth({ login: loginFn });

    renderLogin();

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => screen.getByText(/reenviar email de verificación/i));
    fireEvent.click(screen.getByText(/reenviar email de verificación/i));

    await waitFor(() => {
      expect(authService.resendVerification).toHaveBeenCalledWith('user@example.com');
    });
  });
});

// ---------------------------------------------------------------------------
// Normal login flows
// ---------------------------------------------------------------------------

describe('Login — normal flows', () => {
  afterEach(() => jest.restoreAllMocks());

  it('redirects to dashboard on successful login', async () => {
    mockNo2FA();
    const loginFn = jest.fn().mockResolvedValue({ success: true, data: { user: { first_name: 'Test' } } });
    mockUseAuth({ login: loginFn, isAuthenticated: false });

    renderLogin();

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'user@example.com' } });
    fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

    await waitFor(() => {
      expect(loginFn).toHaveBeenCalled();
    });
  });

  it('renders the login form initially', () => {
    mockUseAuth();
    renderLogin();

    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/contraseña/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /iniciar sesión/i })).toBeInTheDocument();
  });
});
