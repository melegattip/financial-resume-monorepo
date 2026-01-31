import React from 'react';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import Dashboard from '../../pages/Dashboard';
import { renderWithRouter, mockAPI, mockDashboardData, setupTest } from '../utils/testUtils';
import * as api from '../../services/api';

// Mock del módulo de API
jest.mock('../../services/api');

// Mock de Recharts
jest.mock('recharts', () => ({
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
}));

describe('Dashboard Component', () => {
  beforeEach(() => {
    setupTest();
    
    // Setup API mocks
    api.expensesAPI = mockAPI.expenses;
    api.incomesAPI = mockAPI.incomes;
    api.categoriesAPI = mockAPI.categories;
    api.dashboardAPI = mockAPI.dashboard;
    api.formatCurrency = jest.fn((amount) => `$${amount?.toFixed(2) || '0.00'}`);
    api.formatDate = jest.fn((date) => new Date(date).toLocaleDateString());
    api.formatPercentage = jest.fn((percentage) => `${percentage}%`);
  });

  test('muestra estado de loading inicial', () => {
    renderWithRouter(<Dashboard />);
    
    expect(screen.getByText('Cargando dashboard...')).toBeInTheDocument();
  });

  test('renderiza métricas financieras correctamente', async () => {
    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: mockDashboardData }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('Balance Total')).toBeInTheDocument();
      expect(screen.getByText('Total Ingresos')).toBeInTheDocument();
      expect(screen.getByText('Total Gastos')).toBeInTheDocument();
      expect(screen.getByText('Gastos Pendientes')).toBeInTheDocument();
    });

    // Verificar que los montos se muestran
    expect(screen.getByText('$400.00')).toBeInTheDocument(); // Balance
    expect(screen.getByText('$1000.00')).toBeInTheDocument(); // Ingresos
    expect(screen.getByText('$600.00')).toBeInTheDocument(); // Gastos
  });

  test('maneja errores de API correctamente', async () => {
    mockAPI.dashboard.overview.mockRejectedValueOnce(new Error('API Error'));

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('$0.00')).toBeInTheDocument();
    });
  });

  test('renderiza gráficos cuando hay datos', async () => {
    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: mockDashboardData }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByTestId('responsive-container')).toBeInTheDocument();
      expect(screen.getByTestId('area-chart')).toBeInTheDocument();
    });
  });

  test('muestra mensaje cuando no hay datos para gráficos', async () => {
    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: { ...mockDashboardData, expenses: [], incomes: [] } }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('No hay gastos por categorías')).toBeInTheDocument();
    });
  });

  test('calcula balance correctamente', async () => {
    const testData = {
      totalIncome: 1500,
      totalExpenses: 900,
      balance: 600,
      expenses: [],
      incomes: [],
      categories: [],
    };

    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: testData }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('$600.00')).toBeInTheDocument();
    });
  });

  test('muestra indicador de balance positivo', async () => {
    const positiveBalance = {
      ...mockDashboardData,
      balance: 500,
    };

    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: positiveBalance }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('Positivo')).toBeInTheDocument();
    });
  });

  test('muestra indicador de balance negativo', async () => {
    const negativeBalance = {
      ...mockDashboardData,
      balance: -200,
    };

    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: negativeBalance }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('Negativo')).toBeInTheDocument();
    });
  });

  test('fallback a endpoints legacy cuando nuevos fallan', async () => {
    // Primer llamado falla (nuevos endpoints)
    mockAPI.dashboard.overview.mockRejectedValueOnce(new Error('Not found'));
    
    // Endpoints legacy responden
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [mockDashboardData.expenses[0]] });
    mockAPI.incomes.list.mockResolvedValueOnce({ data: [mockDashboardData.incomes[0]] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: mockDashboardData.categories });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(mockAPI.expenses.list).toHaveBeenCalled();
      expect(mockAPI.incomes.list).toHaveBeenCalled();
      expect(mockAPI.categories.list).toHaveBeenCalled();
    });
  });

  test('muestra transacciones recientes cuando las hay', async () => {
    const dataWithTransactions = {
      ...mockDashboardData,
      expenses: [
        { ...mockDashboardData.expenses[0], description: 'Gasto reciente' }
      ],
      incomes: [
        { ...mockDashboardData.incomes[0], description: 'Ingreso reciente' }
      ]
    };

    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: dataWithTransactions }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      expect(screen.getByText('Transacciones Recientes')).toBeInTheDocument();
      expect(screen.getByText('Gasto reciente')).toBeInTheDocument();
      expect(screen.getByText('Ingreso reciente')).toBeInTheDocument();
    });
  });

  test('ordena transacciones por fecha correctamente', async () => {
    const dataWithMultipleTransactions = {
      ...mockDashboardData,
      expenses: [
        { ...mockDashboardData.expenses[0], created_at: '2024-01-15T10:00:00Z', description: 'Gasto nuevo' },
        { ...mockDashboardData.expenses[0], id: 2, created_at: '2024-01-10T10:00:00Z', description: 'Gasto viejo' }
      ]
    };

    mockAPI.dashboard.overview.mockResolvedValueOnce({
      data: { Metrics: dataWithMultipleTransactions }
    });

    renderWithRouter(<Dashboard />);

    await waitFor(() => {
      const transactions = screen.getAllByText(/Gasto/);
      expect(transactions[0]).toHaveTextContent('Gasto nuevo');
    });
  });
}); 