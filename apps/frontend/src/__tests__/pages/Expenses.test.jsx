import React from 'react';
import { screen, waitFor, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Expenses from '../../pages/Expenses';
import { renderWithRouter, mockAPI, mockExpense, mockCategory, setupTest } from '../utils/testUtils';
import * as api from '../../services/api';

// Mock del módulo de API
jest.mock('../../services/api');

describe('Expenses Component', () => {
  const user = userEvent.setup();

  beforeEach(() => {
    setupTest();
    
    // Setup API mocks
    api.expensesAPI = mockAPI.expenses;
    api.categoriesAPI = mockAPI.categories;
    api.formatCurrency = jest.fn((amount) => `$${amount?.toFixed(2) || '0.00'}`);
    api.formatPercentage = jest.fn((percentage) => `${percentage}%`);
  });

  test('muestra estado de loading inicial', () => {
    renderWithRouter(<Expenses />);
    
    expect(screen.getByText('Cargando gastos...')).toBeInTheDocument();
  });

  test('renderiza lista de gastos correctamente', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [mockExpense] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Test Expense')).toBeInTheDocument();
      expect(screen.getByText('$100.50')).toBeInTheDocument();
    });
  });

  test('muestra métricas de gastos correctamente', async () => {
    const multipleExpenses = [
      mockExpense,
      { ...mockExpense, id: 2, amount: 200, paid: true }
    ];

    mockAPI.expenses.list.mockResolvedValueOnce({ data: multipleExpenses });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Total Gastos')).toBeInTheDocument();
      expect(screen.getByText('Gastos Pendientes')).toBeInTheDocument();
      expect(screen.getByText('Monto Pendiente')).toBeInTheDocument();
    });
  });

  test('filtra gastos por búsqueda', async () => {
    const multipleExpenses = [
      { ...mockExpense, description: 'Comida' },
      { ...mockExpense, id: 2, description: 'Transporte' }
    ];

    mockAPI.expenses.list.mockResolvedValueOnce({ data: multipleExpenses });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Comida')).toBeInTheDocument();
      expect(screen.getByText('Transporte')).toBeInTheDocument();
    });

    const searchInput = screen.getByPlaceholderText('Buscar gastos...');
    await user.type(searchInput, 'Comida');

    expect(screen.getByText('Comida')).toBeInTheDocument();
    expect(screen.queryByText('Transporte')).not.toBeInTheDocument();
  });

  test('filtra gastos por estado de pago', async () => {
    const multipleExpenses = [
      { ...mockExpense, paid: false },
      { ...mockExpense, id: 2, paid: true }
    ];

    mockAPI.expenses.list.mockResolvedValueOnce({ data: multipleExpenses });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getAllByText('Test Expense')).toHaveLength(2);
    });

    const filterSelect = screen.getByDisplayValue('Todos los gastos');
    await user.selectOptions(filterSelect, 'paid');

    await waitFor(() => {
      expect(screen.getAllByText('Test Expense')).toHaveLength(1);
    });
  });

  test('abre modal de nuevo gasto al hacer clic en botón', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Nuevo Gasto')).toBeInTheDocument();
    });

    const newButton = screen.getByText('Nuevo Gasto');
    await user.click(newButton);

    expect(screen.getByRole('dialog', { hidden: true })).toBeInTheDocument();
    expect(screen.getByText('Descripción')).toBeInTheDocument();
    expect(screen.getByText('Monto')).toBeInTheDocument();
  });

  test('crea nuevo gasto correctamente', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });
    mockAPI.expenses.create.mockResolvedValueOnce({ data: mockExpense });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      const newButton = screen.getByText('Nuevo Gasto');
      user.click(newButton);
    });

    await waitFor(() => {
      expect(screen.getByRole('dialog', { hidden: true })).toBeInTheDocument();
    });

    // Llenar formulario
    const descriptionInput = screen.getByLabelText(/descripción/i);
    const amountInput = screen.getByLabelText(/monto/i);
    
    await user.type(descriptionInput, 'Test Expense');
    await user.type(amountInput, '100.50');

    const createButton = screen.getByText('Crear');
    await user.click(createButton);

    await waitFor(() => {
      expect(mockAPI.expenses.create).toHaveBeenCalledWith({
        description: 'Test Expense',
        amount: 100.50,
        category_id: '',
        due_date: '',
        paid: false,
      });
    });
  });

  test('edita gasto existente', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [mockExpense] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });
    mockAPI.expenses.update.mockResolvedValueOnce({ data: mockExpense });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Test Expense')).toBeInTheDocument();
    });

    const editButton = screen.getByLabelText(/editar/i);
    await user.click(editButton);

    await waitFor(() => {
      expect(screen.getByDisplayValue('Test Expense')).toBeInTheDocument();
    });

    const descriptionInput = screen.getByDisplayValue('Test Expense');
    await user.clear(descriptionInput);
    await user.type(descriptionInput, 'Updated Expense');

    const updateButton = screen.getByText('Actualizar');
    await user.click(updateButton);

    await waitFor(() => {
      expect(mockAPI.expenses.update).toHaveBeenCalled();
    });
  });

  test('elimina gasto con confirmación', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [mockExpense] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });
    mockAPI.expenses.delete.mockResolvedValueOnce({});

    // Mock window.confirm
    const confirmSpy = jest.spyOn(window, 'confirm').mockReturnValue(true);

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Test Expense')).toBeInTheDocument();
    });

    const deleteButton = screen.getByLabelText(/eliminar/i);
    await user.click(deleteButton);

    await waitFor(() => {
      expect(confirmSpy).toHaveBeenCalledWith('¿Estás seguro de que quieres eliminar este gasto?');
      expect(mockAPI.expenses.delete).toHaveBeenCalled();
    });

    confirmSpy.mockRestore();
  });

  test('cambia estado de pago de gasto', async () => {
    const unpaidExpense = { ...mockExpense, paid: false };
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [unpaidExpense] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });
    mockAPI.expenses.update.mockResolvedValueOnce({ data: { ...unpaidExpense, paid: true } });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('Test Expense')).toBeInTheDocument();
    });

    const payButton = screen.getByLabelText(/marcar como pagado/i);
    await user.click(payButton);

    await waitFor(() => {
      expect(screen.getByText('Pago Total')).toBeInTheDocument();
    });

    const payTotalButton = screen.getByText('Pago Total');
    await user.click(payTotalButton);

    await waitFor(() => {
      expect(mockAPI.expenses.update).toHaveBeenCalledWith(
        unpaidExpense.user_id,
        unpaidExpense.id,
        { paid: true }
      );
    });
  });

  test('maneja errores de carga de datos', async () => {
    mockAPI.expenses.list.mockRejectedValueOnce(new Error('API Error'));
    mockAPI.categories.list.mockRejectedValueOnce(new Error('API Error'));

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      expect(screen.getByText('No hay gastos')).toBeInTheDocument();
    });
  });

  test('valida campos requeridos en formulario', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      const newButton = screen.getByText('Nuevo Gasto');
      user.click(newButton);
    });

    await waitFor(() => {
      const createButton = screen.getByText('Crear');
      user.click(createButton);
    });

    // Los campos requeridos deben tener el atributo required
    expect(screen.getByLabelText(/descripción/i)).toBeRequired();
    expect(screen.getByLabelText(/monto/i)).toBeRequired();
  });

  test('cancela operación de modal', async () => {
    mockAPI.expenses.list.mockResolvedValueOnce({ data: [] });
    mockAPI.categories.list.mockResolvedValueOnce({ data: [mockCategory] });

    renderWithRouter(<Expenses />);

    await waitFor(() => {
      const newButton = screen.getByText('Nuevo Gasto');
      user.click(newButton);
    });

    await waitFor(() => {
      expect(screen.getByRole('dialog', { hidden: true })).toBeInTheDocument();
    });

    const cancelButton = screen.getByText('Cancelar');
    await user.click(cancelButton);

    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
  });
}); 