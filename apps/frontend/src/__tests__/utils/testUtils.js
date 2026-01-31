import React from 'react';
import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import toast, { Toaster } from 'react-hot-toast';
import { AuthProvider } from '../../contexts/AuthContext';
import { ThemeProvider } from '../../contexts/ThemeContext';
import { PeriodProvider } from '../../contexts/PeriodContext';

// Custom render que incluye providers necesarios
export const renderWithRouter = (ui, options = {}) => {
  const Wrapper = ({ children }) => (
    <MemoryRouter initialEntries={['/']}>
      <ThemeProvider>
        <AuthProvider>
          <PeriodProvider>
            {children}
            <Toaster />
          </PeriodProvider>
        </AuthProvider>
      </ThemeProvider>
    </MemoryRouter>
  );

  return render(ui, { wrapper: Wrapper, ...options });
};

// Custom render sin router para cuando el componente ya tiene uno
export const renderWithoutRouter = (ui, options = {}) => {
  const Wrapper = ({ children }) => (
    <ThemeProvider>
      <AuthProvider>
        <PeriodProvider>
          {children}
          <Toaster />
        </PeriodProvider>
      </AuthProvider>
    </ThemeProvider>
  );

  return render(ui, { wrapper: Wrapper, ...options });
};

// Mock data para tests
export const mockExpense = {
  id: 1,
  user_id: 1,
  description: 'Test Expense',
  amount: 100.50,
  category_id: 1,
  paid: false,
  due_date: '2024-01-15',
  created_at: '2024-01-10T10:00:00Z',
  updated_at: '2024-01-10T10:00:00Z',
};

export const mockIncome = {
  id: 1,
  user_id: 1,
  description: 'Test Income',
  amount: 500.00,
  category_id: 2,
  created_at: '2024-01-10T10:00:00Z',
  updated_at: '2024-01-10T10:00:00Z',
};

export const mockCategory = {
  id: 1,
  name: 'Test Category',
  description: 'Category for testing',
  created_at: '2024-01-01T10:00:00Z',
  updated_at: '2024-01-01T10:00:00Z',
};

export const mockDashboardData = {
  totalIncome: 1000,
  totalExpenses: 600,
  balance: 400,
  expenses: [mockExpense],
  incomes: [mockIncome],
  categories: [mockCategory],
};

// Mock APIs
export const mockAPI = {
  expenses: {
    list: jest.fn(() => Promise.resolve({ data: [mockExpense] })),
    create: jest.fn(() => Promise.resolve({ data: mockExpense })),
    update: jest.fn(() => Promise.resolve({ data: mockExpense })),
    delete: jest.fn(() => Promise.resolve()),
  },
  incomes: {
    list: jest.fn(() => Promise.resolve({ data: [mockIncome] })),
    create: jest.fn(() => Promise.resolve({ data: mockIncome })),
    update: jest.fn(() => Promise.resolve({ data: mockIncome })),
    delete: jest.fn(() => Promise.resolve()),
  },
  categories: {
    list: jest.fn(() => Promise.resolve({ data: { data: [mockCategory] } })),
  },
  dashboard: {
    overview: jest.fn(() => Promise.resolve({ data: mockDashboardData })),
  },
};

// Utility para esperar por carga asíncrona
export const waitForLoading = async () => {
  await new Promise(resolve => setTimeout(resolve, 0));
};

// Mock para toast notifications
export const mockToast = {
  success: jest.fn(),
  error: jest.fn(),
  loading: jest.fn(),
};

// Setup para cada test
export const setupTest = () => {
  // Reset all mocks
  jest.clearAllMocks();
  
  // Reset toast
  toast.success = mockToast.success;
  toast.error = mockToast.error;
  toast.loading = mockToast.loading;
};

// Mock para Recharts components
export const mockRecharts = {
  ResponsiveContainer: ({ children }) => <div data-testid="responsive-container">{children}</div>,
  AreaChart: ({ children }) => <div data-testid="area-chart">{children}</div>,
  Area: () => <div data-testid="area" />,
  XAxis: () => <div data-testid="x-axis" />,
  YAxis: () => <div data-testid="y-axis" />,
  CartesianGrid: () => <div data-testid="cartesian-grid" />,
  Tooltip: () => <div data-testid="tooltip" />,
  PieChart: ({ children }) => <div data-testid="pie-chart">{children}</div>,
  Pie: () => <div data-testid="pie" />,
  Cell: () => <div data-testid="cell" />,
};

// Test dummy para evitar error de suite vacía
describe('Test Utils', () => {
  test('should export utilities correctly', () => {
    expect(renderWithRouter).toBeDefined();
    expect(mockAPI).toBeDefined();
    expect(setupTest).toBeDefined();
  });
}); 