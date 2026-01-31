import React from 'react';
import { screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AppContent } from '../../App';
import { renderWithRouter, mockAPI, mockExpense, mockIncome, mockCategory, setupTest } from '../utils/testUtils';
import * as api from '../../services/api';

// Mock del módulo de API completo
jest.mock('../../services/api');

describe('App Integration Tests', () => {
  const user = userEvent.setup();

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

  test('navegación completa entre páginas', async () => {
    // Setup datos mock
    mockAPI.dashboard.overview.mockResolvedValue({
      data: { 
        Metrics: {
          totalIncome: 1000,
          totalExpenses: 600,
          balance: 400,
          expenses: [mockExpense],
          incomes: [mockIncome],
          categories: [mockCategory]
        }
      }
    });
    mockAPI.expenses.list.mockResolvedValue({ data: [mockExpense] });
    mockAPI.incomes.list.mockResolvedValue({ data: [mockIncome] });
    mockAPI.categories.list.mockResolvedValue({ data: [mockCategory] });

    renderWithRouter(<AppContent />);

    // Inicialmente estamos en Dashboard
    await waitFor(() => {
      expect(screen.getByText('Balance Total')).toBeInTheDocument();
    });

    // Navegar a Gastos
    const expensesLink = screen.getByText('Gastos');
    await user.click(expensesLink);

    await waitFor(() => {
      expect(screen.getByText('Total Gastos')).toBeInTheDocument();
    });

    // Navegar a Ingresos
    const incomesLink = screen.getByText('Ingresos');
    await user.click(incomesLink);

    await waitFor(() => {
      expect(screen.getByText('Total Ingresos')).toBeInTheDocument();
    });

    // Navegar a Reportes
    const reportsLink = screen.getByText('Reportes');
    await user.click(reportsLink);

    await waitFor(() => {
      expect(screen.getByText('Reportes Financieros')).toBeInTheDocument();
    });

    // Volver al Dashboard
    const dashboardLink = screen.getByText('Dashboard');
    await user.click(dashboardLink);

    await waitFor(() => {
      expect(screen.getByText('Balance Total')).toBeInTheDocument();
    });
  });

  test('flujo completo CRUD de gastos', async () => {
    // Setup inicial
    mockAPI.expenses.list.mockResolvedValue({ data: [] });
    mockAPI.categories.list.mockResolvedValue({ data: [mockCategory] });
    
    renderWithRouter(<AppContent />);

    // Navegar a gastos
    const expensesLink = screen.getByText('Gastos');
    await user.click(expensesLink);

    await waitFor(() => {
      expect(screen.getByText('No hay gastos')).toBeInTheDocument();
    });

    // CREAR nuevo gasto
    mockAPI.expenses.create.mockResolvedValueOnce({ data: mockExpense });
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [mockExpense] });

    const newButton = screen.getByText('Nuevo Gasto');
    await user.click(newButton);

    // Llenar formulario
    const descriptionInput = screen.getByLabelText(/descripción/i);
    const amountInput = screen.getByLabelText(/monto/i);
    
    await user.type(descriptionInput, 'Test Expense');
    await user.type(amountInput, '100.50');

    const createButton = screen.getByText('Crear');
    await user.click(createButton);

    // Verificar que se creó
    await waitFor(() => {
      expect(mockAPI.expenses.create).toHaveBeenCalled();
    });

    // EDITAR gasto
    mockAPI.expenses.update.mockResolvedValueOnce({ 
      data: { ...mockExpense, description: 'Updated Expense' } 
    });
    mockAPI.expenses.list.mockResolvedValueOnce({ 
      data: [{ ...mockExpense, description: 'Updated Expense' }] 
    });

    await waitFor(() => {
      expect(screen.getByText('Test Expense')).toBeInTheDocument();
    });

    const editButton = screen.getByLabelText(/editar/i);
    await user.click(editButton);

    const editDescriptionInput = screen.getByDisplayValue('Test Expense');
    await user.clear(editDescriptionInput);
    await user.type(editDescriptionInput, 'Updated Expense');

    const updateButton = screen.getByText('Actualizar');
    await user.click(updateButton);

    // Verificar actualización
    await waitFor(() => {
      expect(mockAPI.expenses.update).toHaveBeenCalled();
    });

    // ELIMINAR gasto
    mockAPI.expenses.delete.mockResolvedValueOnce({});
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [] });

    const confirmSpy = jest.spyOn(window, 'confirm').mockReturnValue(true);

    const deleteButton = screen.getByLabelText(/eliminar/i);
    await user.click(deleteButton);

    await waitFor(() => {
      expect(mockAPI.expenses.delete).toHaveBeenCalled();
    });

    confirmSpy.mockRestore();
  });

  test('sincronización de datos entre dashboard y páginas específicas', async () => {
    const testExpense = { ...mockExpense, amount: 250 };
    
    // Setup inicial - Dashboard muestra datos
    mockAPI.dashboard.overview.mockResolvedValue({
      data: { 
        Metrics: {
          totalIncome: 1000,
          totalExpenses: 250,
          balance: 750,
          expenses: [testExpense],
          incomes: [mockIncome],
          categories: [mockCategory]
        }
      }
    });

    renderWithRouter(<AppContent />);

    // Verificar datos en Dashboard
    await waitFor(() => {
      expect(screen.getByText('$750.00')).toBeInTheDocument(); // Balance
      expect(screen.getByText('$250.00')).toBeInTheDocument(); // Gastos totales
    });

    // Navegar a Gastos y verificar consistencia
    mockAPI.expenses.list.mockResolvedValue({ data: [testExpense] });
    mockAPI.categories.list.mockResolvedValue({ data: [mockCategory] });

    const expensesLink = screen.getByText('Gastos');
    await user.click(expensesLink);

    await waitFor(() => {
      expect(screen.getByText('$250.00')).toBeInTheDocument(); // Mismo monto
    });

    // Crear nuevo gasto y verificar que se actualiza
    const newExpense = { ...mockExpense, id: 2, amount: 150 };
    mockAPI.expenses.create.mockResolvedValueOnce({ data: newExpense });
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [testExpense, newExpense] });

    const newButton = screen.getByText('Nuevo Gasto');
    await user.click(newButton);

    const descriptionInput = screen.getByLabelText(/descripción/i);
    const amountInput = screen.getByLabelText(/monto/i);
    
    await user.type(descriptionInput, 'New Expense');
    await user.type(amountInput, '150');

    const createButton = screen.getByText('Crear');
    await user.click(createButton);

    // Verificar que el total se actualiza
    await waitFor(() => {
      expect(mockAPI.expenses.create).toHaveBeenCalled();
    });
  });

  test('manejo de errores de red y recuperación', async () => {
    // Simular error de red inicial
    mockAPI.dashboard.overview.mockRejectedValueOnce(new Error('Network Error'));
    
    renderWithRouter(<AppContent />);

    // Debe mostrar valores por defecto
    await waitFor(() => {
      expect(screen.getByText('$0.00')).toBeInTheDocument();
    });

    // Navegar a gastos con error
    mockAPI.expenses.list.mockRejectedValueOnce(new Error('Network Error'));
    mockAPI.categories.list.mockRejectedValueOnce(new Error('Network Error'));

    const expensesLink = screen.getByText('Gastos');
    await user.click(expensesLink);

    await waitFor(() => {
      expect(screen.getByText('No hay gastos')).toBeInTheDocument();
    });

    // Simular recuperación de red
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [mockExpense] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    // Refrescar página (simular retry)
    window.location.reload = jest.fn();
    
    // Los datos deberían cargar correctamente ahora
    await waitFor(() => {
      expect(screen.getByText('Total Gastos')).toBeInTheDocument();
    });
  });

  test('filtrado y búsqueda consistente entre páginas', async () => {
    const expenses = [
      { ...mockExpense, description: 'Comida', amount: 50 },
      { ...mockExpense, id: 2, description: 'Transporte', amount: 30 },
      { ...mockExpense, id: 3, description: 'Comida rápida', amount: 20 }
    ];

    mockAPI.expenses.list.mockResolvedValue({ data: expenses });
    mockAPI.categories.list.mockResolvedValue({ data: [mockCategory] });

    renderWithRouter(<AppContent />);

    // Navegar a gastos
    const expensesLink = screen.getByText('Gastos');
    await user.click(expensesLink);

    await waitFor(() => {
      expect(screen.getByText('Comida')).toBeInTheDocument();
      expect(screen.getByText('Transporte')).toBeInTheDocument();
      expect(screen.getByText('Comida rápida')).toBeInTheDocument();
    });

    // Filtrar por búsqueda
    const searchInput = screen.getByPlaceholderText('Buscar gastos...');
    await user.type(searchInput, 'Comida');

    // Solo deben aparecer los gastos de comida
    expect(screen.getByText('Comida')).toBeInTheDocument();
    expect(screen.getByText('Comida rápida')).toBeInTheDocument();
    expect(screen.queryByText('Transporte')).not.toBeInTheDocument();

    // Limpiar búsqueda
    await user.clear(searchInput);

    // Todos los gastos deben volver a aparecer
    await waitFor(() => {
      expect(screen.getByText('Comida')).toBeInTheDocument();
      expect(screen.getByText('Transporte')).toBeInTheDocument();
      expect(screen.getByText('Comida rápida')).toBeInTheDocument();
    });
  });

  test('responsive behavior en diferentes tamaños', async () => {
    // Simular viewport móvil
    global.innerWidth = 375;
    global.innerHeight = 667;
    global.dispatchEvent(new Event('resize'));

    mockAPI.dashboard.overview.mockResolvedValue({
      data: { 
        Metrics: {
          totalIncome: 1000,
          totalExpenses: 600,
          balance: 400,
          expenses: [mockExpense],
          incomes: [mockIncome],
          categories: [mockCategory]
        }
      }
    });

    renderWithRouter(<AppContent />);

    await waitFor(() => {
      expect(screen.getByText('Balance Total')).toBeInTheDocument();
    });

    // Verificar que los componentes se adaptan
    const container = screen.getByTestId('main-container') || document.querySelector('.container');
    expect(container).toBeInTheDocument();

    // Simular viewport desktop
    global.innerWidth = 1920;
    global.innerHeight = 1080;
    global.dispatchEvent(new Event('resize'));

    // Los componentes deben seguir funcionando
    await waitFor(() => {
      expect(screen.getByText('Balance Total')).toBeInTheDocument();
    });
  });
}); 