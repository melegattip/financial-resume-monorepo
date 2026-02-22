import React, { useState, useEffect, useCallback } from 'react';
import { createPortal } from 'react-dom';
import { useLocation } from 'react-router-dom';
import { FaPlus, FaSearch, FaArrowDown, FaCalendar, FaEdit, FaTrash, FaCheckCircle, FaTimesCircle, FaDollarSign } from 'react-icons/fa';
import { formatCurrency, formatPercentage, expensesAPI as expensesAPIraw } from '../services/api';
import { usePeriod } from '../contexts/PeriodContext';
import { useGamification } from '../contexts/GamificationContext';
import { useAuth } from '../contexts/AuthContext';
import { useOptimizedAPI } from '../hooks/useOptimizedAPI';
import useDataRefresh from '../hooks/useDataRefresh';
import toast from 'react-hot-toast';
import ConfirmationModal from '../components/ConfirmationModal';
import ValidatedInput from '../components/ValidatedInput';
import { validateAmount, validateDescription, VALIDATION_RULES } from '../utils/validation';

const Expenses = () => {
  const location = useLocation();
  const [expenses, setExpenses] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [showOverpaymentModal, setShowOverpaymentModal] = useState(false);
  const [editingExpense, setEditingExpense] = useState(null);
  const [payingExpense, setPayingExpense] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterPaid, setFilterPaid] = useState('all');
  const [paymentAmount, setPaymentAmount] = useState('');
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingExpense, setDeletingExpense] = useState(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  // Estados para edición inline tipo Excel
  const [editingCell, setEditingCell] = useState(null); // { expenseId, field }
  const [editValues, setEditValues] = useState({});
  const [savingCell, setSavingCell] = useState(null);

  // Estados para nuevos filtros de ordenamiento
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState('desc');
  const [totalIncome, setTotalIncome] = useState(0);
  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    category_id: '',
    due_date: '',
    paid: false,
  });

  // Estados para validación del formulario
  const [formErrors, setFormErrors] = useState({});
  const [isFormValid, setIsFormValid] = useState(false);

  // Usar el contexto global de período
  const {
    selectedYear,
    selectedMonth,
    balancesHidden,
    updateAvailableData,
  } = usePeriod();

  // Usar el hook optimizado para operaciones API
  const {
    expenses: expensesAPI,
    categories: categoriesAPI,
    incomes: incomesAPI,
    dataService
  } = useOptimizedAPI();

  // Hook de gamificación para registrar acciones
  const { recordCreateExpense, recordUpdateExpense, recordDeleteExpense } = useGamification();

  // Hook de autenticación
  const { user } = useAuth();


  // Leer parámetros de URL y aplicar filtros automáticamente
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const statusParam = searchParams.get('status');

    if (statusParam) {
      // Mapear parámetros de URL a valores del filtro
      const filterMapping = {
        'pending': 'unpaid',
        'paid': 'paid',
        'all': 'all'
      };

      const newFilter = filterMapping[statusParam] || 'all';
      console.log(`🔍 [Expenses] Aplicando filtro desde URL: ${statusParam} → ${newFilter}`);
      setFilterPaid(newFilter);
    }
  }, [location.search]);

  const formatAmount = (amount) => {
    if (balancesHidden) return '••••••';
    return formatCurrency(amount);
  };

  // Función para obtener colores por categoría (consistente con Dashboard)
  const getCategoryColor = (categoryId) => {
    const colors = [
      { bg: 'bg-blue-100 dark:bg-blue-900/30', border: 'border-blue-400', text: 'text-blue-700 dark:text-blue-300' },
      { bg: 'bg-green-100 dark:bg-green-900/30', border: 'border-green-400', text: 'text-green-700 dark:text-green-300' },
      { bg: 'bg-yellow-100 dark:bg-yellow-900/30', border: 'border-yellow-400', text: 'text-yellow-700 dark:text-yellow-300' },
      { bg: 'bg-purple-100 dark:bg-purple-900/30', border: 'border-purple-400', text: 'text-purple-700 dark:text-purple-300' },
      { bg: 'bg-pink-100 dark:bg-pink-900/30', border: 'border-pink-400', text: 'text-pink-700 dark:text-pink-300' },
      { bg: 'bg-indigo-100 dark:bg-indigo-900/30', border: 'border-indigo-400', text: 'text-indigo-700 dark:text-indigo-300' },
      { bg: 'bg-cyan-100 dark:bg-cyan-900/30', border: 'border-cyan-400', text: 'text-cyan-700 dark:text-cyan-300' },
      { bg: 'bg-orange-100 dark:bg-orange-900/30', border: 'border-orange-400', text: 'text-orange-700 dark:text-orange-300' },
    ];

    if (!categoryId) {
      return { bg: 'bg-gray-100 dark:bg-gray-700', border: 'border-gray-400 dark:border-gray-500', text: 'text-gray-700 dark:text-gray-300' };
    }

    // Usar el hash del categoryId para asignar colores consistentes
    let hash = 0;
    for (let i = 0; i < categoryId.length; i++) {
      hash = categoryId.charCodeAt(i) + ((hash << 5) - hash);
    }
    const colorIndex = Math.abs(hash) % colors.length;
    return colors[colorIndex];
  };

  const formatPercentageAmount = (percentage) => {
    if (balancesHidden) return '••••';
    return formatPercentage(percentage);
  };

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      console.log('🔄 Cargando datos de gastos con API optimizada...');

      const [expensesResponse, categoriesResponse, incomesResponse] = await Promise.all([
        expensesAPI.list(),
        categoriesAPI.list(),
        incomesAPI.list(),
      ]);

      // Normalizar datos de respuesta
      const expensesData = expensesResponse.data?.expenses || expensesResponse.expenses || expensesResponse || [];
      const categoriesData = categoriesResponse.data?.data || categoriesResponse.data || categoriesResponse || [];
      const incomesData = incomesResponse.data?.incomes || incomesResponse.incomes || incomesResponse || [];

      setExpenses(Array.isArray(expensesData) ? expensesData : []);
      setCategories(Array.isArray(categoriesData) ? categoriesData : []);

      // Calcular total de ingresos para porcentajes
      const totalIncomeAmount = Array.isArray(incomesData) ?
        incomesData.reduce((sum, income) => sum + (income.amount || 0), 0) : 0;
      setTotalIncome(totalIncomeAmount);

      // Actualizar datos disponibles en el contexto de períodos
      updateAvailableData(expensesData, incomesData);


    } catch (error) {
      console.warn('⚠️ Error al cargar gastos:', error.message);

      // Establecer datos vacíos
      setExpenses([]);
      setCategories([]);
      setTotalIncome(0);

      // No mostrar toast aquí porque useOptimizedAPI ya lo maneja
    } finally {
      setLoading(false);
    }
  }, [expensesAPI, categoriesAPI, incomesAPI, updateAvailableData]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Hook para refrescar automáticamente cuando cambian los datos
  useDataRefresh(loadData, ['expense', 'recurring_transaction']);

  // Validar formulario completo
  const validateForm = useCallback(() => {
    const errors = {};
    let valid = true;

    // Validar descripción
    const descriptionValidation = validateDescription(formData.description);
    if (!descriptionValidation.isValid) {
      errors.description = descriptionValidation.error;
      valid = false;
    }

    // Validar monto
    const amountValidation = validateAmount(formData.amount);
    if (!amountValidation.isValid) {
      errors.amount = amountValidation.error;
      valid = false;
    }

    setFormErrors(errors);
    setIsFormValid(valid);
    return valid;
  }, [formData]);

  // Validar formulario cuando cambien los datos
  useEffect(() => {
    validateForm();
  }, [validateForm]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Validar antes de enviar
    if (!validateForm()) {
      toast.error('Por favor corrige los errores en el formulario');
      return;
    }

    try {
      // Convertir amount a número y fecha a RFC3339 antes de enviar
      const dataToSend = {
        ...formData,
        amount: parseFloat(formData.amount),
        due_date: formData.due_date ? `${formData.due_date}T00:00:00Z` : '',
      };

      // Lanzar la operación sin bloquear la UI
      const operationPromise = editingExpense
        ? expensesAPI.update(editingExpense.id, dataToSend)
        : expensesAPI.create(dataToSend);

      // Cerrar modal y limpiar estado inmediatamente para que la UI avance
      setShowModal(false);
      setEditingExpense(null);
      setFormData({
        description: '',
        amount: '',
        category_id: '',
        due_date: '',
        paid: false,
      });
      setFormErrors({});

      // Forzar limpiar cache de lista y recargar para reflejar cambios lo antes posible
      try {
        dataService?.clearCache?.('expenses_list');
      } catch { }
      await loadData();

      // Esperar resultado para registrar gamificación y hacer una recarga final
      const result = await operationPromise;

      if (editingExpense) {
        console.log(`🎯 [Expenses] Registrando actualización de expense: ${editingExpense.id}`);
        recordUpdateExpense(editingExpense.id, `Gasto actualizado: ${dataToSend.description}`);
      } else {
        const expenseId = result?.data?.id || `expense-${Date.now()}`;
        console.log(`🎯 [Expenses] Registrando creación de expense: ${expenseId}`);
        recordCreateExpense(expenseId, `Nuevo gasto: ${dataToSend.description}`);
      }

      // Recarga final para asegurar consistencia luego de invalidación
      await loadData();
    } catch (error) {
      // useOptimizedAPI ya maneja el error y muestra el toast
      console.error('Error en handleSubmit:', error);
      // Si la request tardó pero el backend creó el recurso, evitamos dejar el modal abierto
      if (error?.code === 'ECONNABORTED') {
        setShowModal(false);
        setEditingExpense(null);
        setFormData({
          description: '',
          amount: '',
          category_id: '',
          due_date: '',
          paid: false,
        });
        setFormErrors({});
        // Forzar recarga de datos para reflejar el gasto creado
        try {
          dataService?.clearCache?.('expenses_list');
        } catch { }
        await loadData();
      }
    }
  };

  const handleEdit = (expense) => {
    setEditingExpense(expense);
    setFormData({
      description: expense.description,
      amount: expense.amount.toString(),
      category_id: expense.category_id || '',
      due_date: expense.due_date || '',
      paid: expense.paid,
    });
    setShowModal(true);
  };

  const handleDelete = (expense) => {
    setDeletingExpense(expense);
    setShowDeleteModal(true);
  };

  const confirmDelete = async () => {
    if (!deletingExpense) return;

    try {
      setDeleteLoading(true);
      await expensesAPI.delete(deletingExpense.id);
      // useOptimizedAPI ya muestra el toast de éxito

      // 🎮 Registrar acción de gamificación
      console.log(`🎯 [Expenses] Registrando eliminación de expense: ${deletingExpense.id}`);
      recordDeleteExpense(deletingExpense.id, `Gasto eliminado: ${deletingExpense.description}`);

      // Recargar datos para mostrar cambios
      await loadData();
    } catch (error) {
      // useOptimizedAPI ya maneja el error
      console.error('Error en confirmDelete:', error);
    } finally {
      // ✅ Siempre cerrar modal y limpiar estado, sin importar si hay errores
      setDeleteLoading(false);
      setShowDeleteModal(false);
      setDeletingExpense(null);
    }
  };

  const cancelDelete = () => {
    setShowDeleteModal(false);
    setDeletingExpense(null);
  };

  const togglePaid = async (expense) => {
    if (expense.paid) {
      // Si ya está pagado, permitir marcarlo como no pagado
      try {
        const updateData = { paid: false };

        // ✨ Actualización optimista
        setExpenses(prevExpenses =>
          prevExpenses.map(exp =>
            exp.id === expense.id
              ? { ...exp, paid: false }
              : exp
          )
        );

        // Llamar API en background sin invalidar caché
        await expensesAPIraw.update(user.id, expense.id, updateData);

        toast.success('Pago anulado', { duration: 1000 });

      } catch (error) {
        console.error('Error en togglePaid:', error);
        toast.error('Error al anular pago');

        // Revertir en caso de error
        await loadData();
      }
    } else {
      // Si no está pagado, abrir modal de pago
      setPayingExpense(expense);
      const pendingAmount = expense.pending_amount || (expense.amount - (expense.amount_paid || 0));
      setPaymentAmount(pendingAmount.toString());
      setShowPaymentModal(true);
    }
  };

  const handlePayment = async (paymentType) => {
    try {
      if (paymentType === 'total') {
        // Pago total - marcar como pagado (prioridad absoluta)
        // Resetea cualquier pago parcial previo y marca como 100% pagado
        const updateData = {
          paid: true,
          amount_paid: payingExpense.amount,
          pending_amount: 0
        };

        // ✨ Actualización optimista
        setExpenses(prevExpenses =>
          prevExpenses.map(exp =>
            exp.id === payingExpense.id
              ? { ...exp, ...updateData }
              : exp
          )
        );

        // Cerrar modal antes de API call
        setShowPaymentModal(false);
        const expenseId = payingExpense.id;
        setPayingExpense(null);
        setPaymentAmount('');

        // Llamar API en background
        await expensesAPIraw.update(user.id, expenseId, updateData);

        toast.success('Gasto marcado como pagado completamente', { duration: 1000 });
      } else if (paymentType === 'partial') {
        // Pago parcial - validar monto
        const paymentAmt = parseFloat(paymentAmount);
        const pendingAmount = payingExpense.pending_amount || (payingExpense.amount - (payingExpense.amount_paid || 0));

        if (paymentAmt <= 0) {
          toast.error('El monto debe ser mayor a 0');
          return;
        }

        // Verificar si intenta pagar más del monto pendiente
        if (paymentAmt > pendingAmount) {
          // Mostrar modal de sobrepago
          setShowOverpaymentModal(true);
          return;
        }

        const updateData = { payment_amount: paymentAmt };

        // Calcular nuevo estado optimista
        const newAmountPaid = (payingExpense.amount_paid || 0) + paymentAmt;
        const newPendingAmount = payingExpense.amount - newAmountPaid;
        const isPaidNow = newPendingAmount <= 0;

        // ✨ Actualización optimista
        setExpenses(prevExpenses =>
          prevExpenses.map(exp =>
            exp.id === payingExpense.id
              ? {
                ...exp,
                amount_paid: newAmountPaid,
                pending_amount: isPaidNow ? 0 : newPendingAmount,
                paid: isPaidNow
              }
              : exp
          )
        );

        // Cerrar modal antes de API call
        setShowPaymentModal(false);
        const expenseId = payingExpense.id;
        setPayingExpense(null);
        setPaymentAmount('');

        // Llamar API en background
        await expensesAPIraw.update(user.id, expenseId, updateData);

        // Feedback apropiado
        if (isPaidNow) {
          toast.success('Gasto pagado completamente', { duration: 1000 });
        } else {
          toast.success(`Pago parcial registrado. Quedan ${formatCurrency(newPendingAmount)} pendientes`, { duration: 2000 });
        }
      }


    } catch (error) {
      console.error('Error en handlePayment:', error);
      toast.error('Error al registrar el pago');

      // Revertir en caso de error
      await loadData();
    }
  };

  const handleOverpaymentChoice = async (choice) => {
    try {
      const paymentAmt = parseFloat(paymentAmount);

      if (choice === 'increase_expense') {
        // Opción 1: Aumentar el gasto total al monto del pago y aplicar pago total
        const updateData = {
          amount: paymentAmt,  // Aumentar el monto total del gasto
          paid: true,          // Marcar como pagado
          amount_paid: paymentAmt,
          pending_amount: 0
        };
        await expensesAPI.update(payingExpense.id, updateData);
        toast.success(`Gasto actualizado a ${formatCurrency(paymentAmt)} y marcado como pagado completamente`);
      } else if (choice === 'apply_total_payment') {
        // Opción 2: Aplicar pago total con el monto original
        const updateData = {
          paid: true,
          amount_paid: payingExpense.amount,
          pending_amount: 0
        };
        await expensesAPI.update(payingExpense.id, updateData);
        toast.success('Gasto marcado como pagado completamente con el monto original');
      }

      // Cerrar modales y limpiar
      setShowOverpaymentModal(false);
      setShowPaymentModal(false);
      setPayingExpense(null);
      setPaymentAmount('');

      // Recargar datos para mostrar cambios
      await loadData();
    } catch (error) {
      console.error('Error en handleOverpaymentChoice:', error);
      toast.error('Error al procesar el pago');
    }
  };

  // ===== Funciones para Edición Inline Tipo Excel =====

  const startEditing = (expenseId, field, currentValue) => {
    setEditingCell({ expenseId, field });
    setEditValues({
      ...editValues,
      [`${expenseId}-${field}`]: currentValue
    });
  };

  const cancelInlineEdit = () => {
    setEditingCell(null);
    setEditValues({});
    setSavingCell(null);
  };

  const saveInlineEdit = async (expenseId, field) => {
    console.log('🚀 [saveInlineEdit] INICIANDO', { expenseId, field });

    const key = `${expenseId}-${field}`;
    const newValue = editValues[key];
    const expense = expenses.find(e => e.id === expenseId);

    console.log('🔍 [saveInlineEdit] Valores:', {
      key,
      newValue,
      hasExpense: !!expense,
      currentValue: expense?.[field]
    });

    if (!expense) return;

    // Validar que el valor haya cambiado (comparación apropiada según el tipo)
    let valueChanged = false;
    if (field === 'amount') {
      const currentAmount = parseFloat(expense.amount);
      const newAmount = parseFloat(newValue);
      valueChanged = newAmount !== currentAmount;
    } else {
      valueChanged = newValue !== expense[field];
    }

    if (!valueChanged) {
      cancelInlineEdit();
      return;
    }

    // Validaciones por campo
    if (field === 'description' && (!newValue || newValue.trim() === '')) {
      toast.error('La descripción no puede estar vacía');
      return;
    }

    if (field === 'amount') {
      const amountNum = parseFloat(newValue);
      if (isNaN(amountNum) || amountNum <= 0) {
        toast.error('El monto debe ser un número positivo');
        return;
      }
    }

    try {
      setSavingCell({ expenseId, field });

      const updateData = {};

      if (field === 'amount') {
        const newAmount = parseFloat(newValue);
        updateData.amount = newAmount;

        // Si hay pago parcial, recalcular pending_amount
        if (expense.amount_paid > 0 && !expense.paid) {
          const newPendingAmount = newAmount - expense.amount_paid;
          updateData.pending_amount = newPendingAmount;

          // Si el nuevo monto hace que pending sea 0 o negativo, marcar como pagado
          if (newPendingAmount <= 0) {
            updateData.paid = true;
            updateData.pending_amount = 0;
          }
        }
      } else if (field === 'category_id') {
        updateData.category_id = newValue;
      } else if (field === 'due_date') {
        // Convert YYYY-MM-DD from date input to RFC3339 for the backend
        updateData.due_date = newValue ? `${newValue}T00:00:00Z` : '';
      } else {
        updateData[field] = newValue;
      }

      // ✨ ACTUALIZACIÓN OPTIMISTA: actualizar UI inmediatamente
      setExpenses(prevExpenses =>
        prevExpenses.map(exp =>
          exp.id === expenseId
            ? { ...exp, ...updateData }
            : exp
        )
      );

      // Limpiar estado de edición antes de llamar API
      cancelInlineEdit();

      // ⚡ LLAMAR API DIRECTO SIN INVALIDACIÓN DE CACHÉ
      // Esto evita el refresh y mantiene la actualización optimista pura
      console.log('🔧 [saveInlineEdit] Llamando API con:', {
        userId: user.id,
        expenseId,
        updateData,
        field
      });

      const response = await expensesAPIraw.update(user.id, expenseId, updateData);

      console.log('✅ [saveInlineEdit] API response:', response.data);

      // 🎮 Registrar edición
      recordUpdateExpense(expenseId, `Campo ${field} editado inline`);

      // ✅ NO recargar datos, la UI ya está actualizada

      // Mostrar feedback breve
      toast.success('Guardado', { duration: 1000 });

    } catch (error) {
      console.error('❌ [saveInlineEdit] Error completo:', {
        error,
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
        field,
        expenseId
      });
      toast.error('Error al guardar el cambio');

      // En caso de error, revertir cambios recargando datos
      await loadData();

      setSavingCell(null);
    }
  };

  const handleInlineKeyDown = async (e, expenseId, field) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      await saveInlineEdit(expenseId, field);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      cancelInlineEdit();
    }
  };

  // ========== Fin funciones edición inline ==========

  const filteredExpenses = Array.isArray(expenses)
    ? expenses.filter(expense => {
      const matchesSearch = expense.description.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesFilter = filterPaid === 'all' ||
        (filterPaid === 'paid' && expense.paid) ||
        (filterPaid === 'unpaid' && !expense.paid);

      // Filtros de fecha — usar transaction_date (fecha funcional del gasto)
      const txDate = expense.due_date || expense.transaction_date || expense.created_at;
      const expenseDate = new Date(txDate);
      const matchesYear = !selectedYear || expenseDate.getFullYear().toString() === selectedYear;
      const matchesMonth = !selectedMonth || txDate.slice(0, 7) === selectedMonth;

      return matchesSearch && matchesFilter && matchesYear && matchesMonth;
    })
      .sort((a, b) => {
        let aValue, bValue;

        switch (sortBy) {
          case 'amount':
            aValue = a.amount;
            bValue = b.amount;
            break;
          case 'category':
            const aCat = categories.find(c => c.id === a.category_id);
            const bCat = categories.find(c => c.id === b.category_id);
            aValue = (aCat?.name || 'Sin categoría').toLowerCase();
            bValue = (bCat?.name || 'Sin categoría').toLowerCase();
            break;
          case 'due_date':
            aValue = a.due_date ? new Date(a.due_date).getTime() : 0;
            bValue = b.due_date ? new Date(b.due_date).getTime() : 0;
            break;
          case 'created_at':
          default:
            aValue = new Date(a.created_at).getTime();
            bValue = new Date(b.created_at).getTime();
            break;
        }

        if (sortOrder === 'asc') {
          return aValue > bValue ? 1 : aValue < bValue ? -1 : 0;
        } else {
          return aValue < bValue ? 1 : aValue > bValue ? -1 : 0;
        }
      })
    : [];

  const totalExpenses = filteredExpenses.reduce((sum, expense) => sum + expense.amount, 0);
  const unpaidExpenses = filteredExpenses.filter(e => !e.paid);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-fr-gray-600 dark:text-gray-400">Cargando gastos...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Title */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">Gastos</h1>
      </div>

      {/* Header con métricas */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-6">
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Total Gastos</p>
                  <p className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">{formatAmount(totalExpenses)}</p>
                </div>
                <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Cantidad</p>
                  <p className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">{filteredExpenses.length}</p>
                </div>
              </div>
            </div>
            <div className="flex-shrink-0 p-3 rounded-fr bg-gray-100 dark:bg-gray-700 ml-4">
              <FaArrowDown className="w-6 h-6 text-fr-gray-900 dark:text-gray-300" />
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <div className="flex items-center space-x-6">
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Pendientes</p>
                  <p className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">{unpaidExpenses.length}</p>
                </div>
                <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                <div>
                  <p className="text-sm font-medium text-fr-gray-600 dark:text-gray-400">Monto</p>
                  <p className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">
                    {formatAmount(unpaidExpenses.reduce((sum, e) => sum + e.amount, 0))}
                  </p>
                </div>
              </div>
            </div>
            <div className="flex-shrink-0 p-3 rounded-fr bg-red-100 dark:bg-red-900/30 ml-4">
              <FaTimesCircle className="w-6 h-6 text-fr-error dark:text-red-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Controles */}
      <div className="card">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="flex flex-col lg:flex-row space-y-4 lg:space-y-0 lg:space-x-4">
            {/* Primera fila: Búsqueda y Estado */}
            <div className="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
              {/* Búsqueda */}
              <div className="relative">
                <FaSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-fr-gray-400 dark:text-gray-500" />
                <input
                  type="text"
                  placeholder="Buscar gastos..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10 input w-full sm:w-64"
                />
              </div>

              {/* Filtro de estado */}
              <select
                value={filterPaid}
                onChange={(e) => setFilterPaid(e.target.value)}
                className="input w-full sm:w-auto"
              >
                <option value="all">Todos los gastos</option>
                <option value="paid">Pagados</option>
                <option value="unpaid">Pendientes</option>
              </select>
            </div>

            {/* Segunda fila: Ordenamiento */}
            <div className="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
              {/* Ordenar por */}
              <div className="flex items-center space-x-2">
                <label className="text-sm font-medium text-fr-gray-700 dark:text-gray-300 whitespace-nowrap">
                  Ordenar por:
                </label>
                <select
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value)}
                  className="input w-full sm:w-auto"
                >
                  <option value="created_at">Fecha de creación</option>
                  <option value="due_date">Fecha de vencimiento</option>
                  <option value="amount">Monto</option>
                  <option value="category">Categoría</option>
                </select>
              </div>

              {/* Orden */}
              <div className="flex items-center space-x-2">
                <label className="text-sm font-medium text-fr-gray-700 dark:text-gray-300 whitespace-nowrap">
                  Orden:
                </label>
                <select
                  value={sortOrder}
                  onChange={(e) => setSortOrder(e.target.value)}
                  className="input w-full sm:w-auto"
                >
                  <option value="desc">Descendente</option>
                  <option value="asc">Ascendente</option>
                </select>
              </div>
            </div>
          </div>

          <button
            onClick={() => setShowModal(true)}
            className="btn-primary flex items-center space-x-2"
          >
            <FaPlus className="w-4 h-4" />
            <span>Nuevo Gasto</span>
          </button>
        </div>
      </div>

      {/* Lista de gastos */}
      <div className="card">
        <div className="space-y-4">
          {filteredExpenses.length === 0 ? (
            <div className="text-center py-12">
              <FaArrowDown className="w-12 h-12 text-fr-gray-400 dark:text-gray-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-fr-gray-900 dark:text-gray-100 mb-2">No hay gastos</h3>
              <p className="text-fr-gray-500 dark:text-gray-400">Comienza agregando tu primer gasto</p>
            </div>
          ) : (
            filteredExpenses.map((expense) => {
              const category = categories.find(c => c.id === expense.category_id);
              const color = getCategoryColor(expense.category_id);
              const incomePercentage = totalIncome > 0 ? (expense.amount / totalIncome) * 100 : 0;

              return (
                <div key={expense.id} className="flex items-center gap-2 py-1.5 px-3 rounded-lg bg-fr-gray-50 dark:bg-gray-700 hover:bg-fr-gray-100 dark:hover:bg-gray-600 transition-colors">
                  {/* Estado de pago compacto */}
                  <div className="flex-shrink-0 w-6 h-6">
                    <button
                      onClick={() => togglePaid(expense)}
                      className={`w-full h-full rounded-md transition-colors flex items-center justify-center ${expense.paid
                        ? 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 hover:bg-green-200 dark:hover:bg-green-900/50'
                        : 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 hover:bg-red-200 dark:hover:bg-red-900/50'
                        }`}
                    >
                      {expense.paid ? (
                        <FaCheckCircle className="w-3 h-3" />
                      ) : (
                        <FaTimesCircle className="w-3 h-3" />
                      )}
                    </button>
                  </div>

                  {/* Descripción - Editable inline */}
                  <div className="flex-1 min-w-0">
                    {editingCell?.expenseId === expense.id && editingCell?.field === 'description' ? (
                      <input
                        type="text"
                        value={editValues[`${expense.id}-description`] ?? expense.description}
                        onChange={(e) => setEditValues({ ...editValues, [`${expense.id}-description`]: e.target.value })}
                        onBlur={() => saveInlineEdit(expense.id, 'description')}
                        onKeyDown={(e) => handleInlineKeyDown(e, expense.id, 'description')}
                        className="w-full px-2 py-1 text-sm font-medium bg-white dark:bg-gray-800 border-2 border-blue-400 dark:border-blue-500 rounded focus:outline-none"
                        autoFocus
                        disabled={savingCell?.expenseId === expense.id && savingCell?.field === 'description'}
                      />
                    ) : (
                      <h3
                        className="font-medium text-fr-gray-900 dark:text-gray-100 text-sm truncate cursor-pointer hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                        onClick={() => startEditing(expense.id, 'description', expense.description)}
                        title="Click para editar"
                      >
                        {expense.description}
                      </h3>
                    )}
                  </div>

                  {/* Categoría - Editable inline */}
                  <div className="flex-shrink-0 hidden sm:flex items-center text-left">
                    {editingCell?.expenseId === expense.id && editingCell?.field === 'category_id' ? (
                      <select
                        value={editValues[`${expense.id}-category_id`] || expense.category_id}
                        onChange={(e) => setEditValues({ ...editValues, [`${expense.id}-category_id`]: e.target.value })}
                        onBlur={() => saveInlineEdit(expense.id, 'category_id')}
                        onKeyDown={(e) => handleInlineKeyDown(e, expense.id, 'category_id')}
                        className="px-2 py-1 text-xs font-medium bg-white dark:bg-gray-800 border-2 border-blue-400 dark:border-blue-500 rounded focus:outline-none"
                        autoFocus
                        disabled={savingCell?.expenseId === expense.id && savingCell?.field === 'category_id'}
                      >
                        {categories.map(cat => (
                          <option key={cat.id} value={cat.id}>{cat.name}</option>
                        ))}
                      </select>
                    ) : (
                      category && (
                        <span
                          className={`px-1.5 py-0.5 rounded-full text-xs font-medium ${color.bg} ${color.text} border ${color.border} whitespace-nowrap cursor-pointer hover:opacity-80 transition-opacity`}
                          onClick={() => startEditing(expense.id, 'category_id', expense.category_id)}
                          title="Click para editar"
                        >
                          {category.name}
                        </span>
                      )
                    )}
                  </div>

                  {/* Fecha de vencimiento - Editable inline */}
                  <div className="flex-shrink-0 hidden md:block text-xs text-gray-500 dark:text-gray-400 text-center whitespace-nowrap">
                    {editingCell?.expenseId === expense.id && editingCell?.field === 'due_date' ? (
                      <input
                        type="date"
                        value={editValues[`${expense.id}-due_date`] || expense.due_date || ''}
                        onChange={(e) => setEditValues({ ...editValues, [`${expense.id}-due_date`]: e.target.value })}
                        onBlur={() => saveInlineEdit(expense.id, 'due_date')}
                        onKeyDown={(e) => handleInlineKeyDown(e, expense.id, 'due_date')}
                        className="px-2 py-1 text-xs bg-white dark:bg-gray-800 border-2 border-blue-400 dark:border-blue-500 rounded focus:outline-none"
                        autoFocus
                        disabled={savingCell?.expenseId === expense.id && savingCell?.field === 'due_date'}
                      />
                    ) : (
                      expense.due_date && (
                        <span
                          className="cursor-pointer hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                          onClick={() => startEditing(expense.id, 'due_date', expense.due_date)}
                          title="Click para editar"
                        >
                          Vto: {new Date(expense.due_date).toLocaleDateString('es-AR', {
                            day: 'numeric',
                            month: 'numeric'
                          })}
                        </span>
                      )
                    )}
                  </div>

                  {/* % de Ingreso */}
                  <div className="flex-shrink-0 hidden lg:block text-xs text-gray-500 dark:text-gray-400 text-center whitespace-nowrap">
                    {incomePercentage.toFixed(1)}%
                  </div>

                  {/* Monto - Editable inline (Total tachado + Saldo pendiente) */}
                  <div className="flex-shrink-0 text-right min-w-[100px]">
                    {editingCell?.expenseId === expense.id && editingCell?.field === 'amount' ? (
                      <input
                        type="number"
                        step="0.01"
                        value={editValues[`${expense.id}-amount`] || expense.amount}
                        onChange={(e) => setEditValues({ ...editValues, [`${expense.id}-amount`]: e.target.value })}
                        onBlur={() => saveInlineEdit(expense.id, 'amount')}
                        onKeyDown={(e) => handleInlineKeyDown(e, expense.id, 'amount')}
                        className="w-full px-2 py-1 text-sm font-semibold bg-white dark:bg-gray-800 border-2 border-blue-400 dark:border-blue-500 rounded focus:outline-none text-right"
                        autoFocus
                        disabled={savingCell?.expenseId === expense.id && savingCell?.field === 'amount'}
                      />
                    ) : (
                      <div
                        className="flex items-center justify-end gap-2 whitespace-nowrap cursor-pointer hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                        onClick={() => startEditing(expense.id, 'amount', expense.amount)}
                        title="Click para editar monto"
                      >
                        {expense.amount_paid > 0 && !expense.paid && (
                          <span className="text-xs text-gray-400 dark:text-gray-500 line-through">
                            -{formatAmount(expense.amount)}
                          </span>
                        )}
                        <span className="font-semibold text-fr-gray-900 dark:text-gray-100 text-sm">
                          -{formatAmount(expense.amount_paid > 0 && !expense.paid ? expense.amount - expense.amount_paid : expense.amount)}
                        </span>
                      </div>
                    )}
                  </div>

                  {/* Botones de acción compactos */}
                  <div className="flex space-x-0.5 flex-shrink-0">
                    {!expense.paid && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setPayingExpense(expense);
                          setShowPaymentModal(true);
                        }}
                        className="w-6 h-6 bg-green-100 dark:bg-green-900/30 rounded-md flex items-center justify-center hover:bg-green-200 dark:hover:bg-green-900/50 transition-colors"
                        title="Pagar"
                      >
                        <FaCheckCircle className="w-3 h-3 text-green-600 dark:text-green-400" />
                      </button>
                    )}

                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleEdit(expense);
                      }}
                      className="w-6 h-6 bg-gray-100 dark:bg-gray-700 rounded-md flex items-center justify-center hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                      title="Editar"
                    >
                      <FaEdit className="w-2.5 h-2.5 text-gray-600 dark:text-gray-400" />
                    </button>

                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDelete(expense);
                      }}
                      className="w-6 h-6 bg-red-100 dark:bg-red-900/30 rounded-md flex items-center justify-center hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors"
                      title="Eliminar"
                    >
                      <FaTrash className="w-2.5 h-2.5 text-red-600 dark:text-red-400" />
                    </button>
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>

      {/* Modal */}
      {showModal && createPortal(
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[9999] p-4">
          <div className="bg-white dark:bg-gray-800 rounded-fr-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-fr-gray-900 dark:text-gray-100 mb-6">
              {editingExpense ? 'Editar Gasto' : 'Nuevo Gasto'}
            </h2>

            <form onSubmit={handleSubmit} className="space-y-4">
              <ValidatedInput
                type="text"
                name="description"
                label="Descripción"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                validator={validateDescription}
                validateOnChange={true}
                required={true}
                placeholder="Ej: Compras del supermercado"
                helpText="Describe brevemente el gasto"
                maxLength={255}
              />

              <ValidatedInput
                type="currency"
                name="amount"
                label="Monto"
                value={formData.amount}
                onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                validator={(value) => validateAmount(value, { fieldName: 'monto' })}
                validateOnChange={true}
                required={true}
                placeholder="0.00"
                helpText="Ingresa el monto del gasto"
                icon={<FaDollarSign />}
                iconPosition="left"
                allowNegative={false}
                maxDecimals={2}
              />

              <div>
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                  Categoría
                </label>
                <select
                  value={formData.category_id}
                  onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                  className="input"
                >
                  <option value="">Seleccionar categoría</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>
                      {category.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                  Fecha del gasto
                </label>
                <input
                  type="date"
                  value={formData.due_date}
                  onChange={(e) => setFormData({ ...formData, due_date: e.target.value })}
                  className="input"
                />
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="paid"
                  checked={formData.paid}
                  onChange={(e) => setFormData({ ...formData, paid: e.target.checked })}
                  className="mr-2"
                />
                <label htmlFor="paid" className="text-sm font-medium text-fr-gray-700 dark:text-gray-300">
                  Marcar como pagado
                </label>
              </div>

              <div className="flex space-x-4 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false);
                    setEditingExpense(null);
                    setFormData({
                      description: '',
                      amount: '',
                      category_id: '',
                      due_date: '',
                      paid: false,
                    });
                    setFormErrors({});
                  }}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className={`btn-primary flex-1 ${!isFormValid ? 'opacity-50 cursor-not-allowed' : ''}`}
                  disabled={!isFormValid}
                >
                  {editingExpense ? 'Actualizar' : 'Crear'}
                </button>
              </div>
            </form>
          </div>
        </div>,
        document.body
      )}

      {/* Modal de Pago */}
      {showPaymentModal && payingExpense && createPortal(
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[9999] p-4">
          <div className="bg-white dark:bg-gray-800 rounded-fr-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-fr-gray-900 dark:text-gray-100 mb-6">
              Registrar Pago
            </h2>

            {/* Información del gasto */}
            <div className="bg-fr-gray-50 dark:bg-gray-700 rounded-fr p-4 mb-6">
              <h3 className="font-medium text-fr-gray-900 dark:text-gray-100 mb-2">{payingExpense.description}</h3>
              <div className="space-y-1">
                <p className="text-lg font-bold text-fr-gray-900 dark:text-gray-100">
                  Monto total: {formatCurrency(payingExpense.amount)}
                </p>
                {payingExpense.amount_paid > 0 && (
                  <>
                    <p className="text-sm text-fr-secondary dark:text-green-400">
                      Ya pagado: {formatCurrency(payingExpense.amount_paid)}
                    </p>
                    <p className="text-lg font-bold text-fr-accent dark:text-yellow-400">
                      Pendiente: {formatCurrency(payingExpense.pending_amount || (payingExpense.amount - payingExpense.amount_paid))}
                    </p>
                  </>
                )}
              </div>
              {payingExpense.due_date && (
                <p className="text-sm text-fr-gray-600 dark:text-gray-400 mt-1">
                  Vence: {new Date(payingExpense.due_date).toLocaleDateString('es-AR')}
                </p>
              )}
            </div>

            {/* Opciones de pago */}
            <div className="space-y-4">
              {/* Pago Total */}
              <button
                onClick={() => handlePayment('total')}
                className="w-full p-4 border-2 border-fr-secondary dark:border-green-600 rounded-fr hover:bg-green-50 dark:hover:bg-green-900/20 transition-colors text-left"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium text-fr-gray-900 dark:text-gray-100">💰 Pago Total</h4>
                    <p className="text-sm text-fr-gray-600 dark:text-gray-400">Marcar como completamente pagado</p>
                  </div>
                  <p className="font-bold text-fr-secondary dark:text-green-400">
                    {formatCurrency(payingExpense.pending_amount || (payingExpense.amount - (payingExpense.amount_paid || 0)))}
                  </p>
                </div>
              </button>

              {/* Pago Parcial */}
              <div className="border-2 border-fr-accent dark:border-yellow-600 rounded-fr p-4">
                <h4 className="font-medium text-fr-gray-900 dark:text-gray-100 mb-3">💸 Pago Parcial</h4>
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                      Monto a pagar
                    </label>
                    <input
                      type="number"
                      step="0.01"
                      max={payingExpense.pending_amount || (payingExpense.amount - (payingExpense.amount_paid || 0))}
                      value={paymentAmount}
                      onChange={(e) => setPaymentAmount(e.target.value)}
                      className="input"
                      placeholder="0.00"
                    />
                  </div>
                  <div className="text-sm text-fr-gray-600 dark:text-gray-400">
                    <p>Quedarían pendientes: <span className="font-medium">
                      {formatCurrency(Math.max(0, (payingExpense.pending_amount || (payingExpense.amount - (payingExpense.amount_paid || 0))) - (parseFloat(paymentAmount) || 0)))}
                    </span></p>
                  </div>
                  <button
                    onClick={() => handlePayment('partial')}
                    disabled={!paymentAmount || parseFloat(paymentAmount) <= 0}
                    className="w-full btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Registrar Pago Parcial
                  </button>
                </div>
              </div>
            </div>

            {/* Botones de acción */}
            <div className="flex space-x-4 pt-6">
              <button
                type="button"
                onClick={() => {
                  setShowPaymentModal(false);
                  setPayingExpense(null);
                  setPaymentAmount('');
                }}
                className="btn-outline flex-1"
              >
                Cancelar
              </button>
            </div>
          </div>
        </div>,
        document.body
      )}

      {/* Modal de Sobrepago */}
      {showOverpaymentModal && payingExpense && createPortal(
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[9999] p-4">
          <div className="bg-white dark:bg-gray-800 rounded-fr-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-fr-gray-900 dark:text-gray-100 mb-4">
              ⚠️ Monto Mayor al Pendiente
            </h2>

            <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-700 rounded-fr p-4 mb-6">
              <p className="text-sm text-yellow-800 dark:text-yellow-300 mb-3">
                Estás intentando pagar <strong>{formatCurrency(parseFloat(paymentAmount))}</strong> pero solo hay <strong>{formatCurrency(payingExpense.pending_amount || (payingExpense.amount - (payingExpense.amount_paid || 0)))}</strong> pendientes.
              </p>
              <p className="text-sm text-yellow-800 dark:text-yellow-300">
                ¿Qué quieres hacer?
              </p>
            </div>

            <div className="space-y-3">
              {/* Opción 1: Aumentar el gasto */}
              <button
                onClick={() => handleOverpaymentChoice('increase_expense')}
                className="w-full p-4 text-left border-2 border-blue-200 dark:border-blue-700 rounded-fr hover:bg-blue-50 dark:hover:bg-blue-900/20 transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium text-fr-gray-900 dark:text-gray-100">
                      📈 Aumentar el gasto total
                    </h4>
                    <p className="text-sm text-fr-gray-600 dark:text-gray-400 mt-1">
                      Cambiar el gasto de {formatCurrency(payingExpense.amount)} a {formatCurrency(parseFloat(paymentAmount))} y marcarlo como pagado
                    </p>
                  </div>
                </div>
              </button>

              {/* Opción 2: Pago total con monto original */}
              <button
                onClick={() => handleOverpaymentChoice('apply_total_payment')}
                className="w-full p-4 text-left border-2 border-green-200 dark:border-green-700 rounded-fr hover:bg-green-50 dark:hover:bg-green-900/20 transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium text-fr-gray-900 dark:text-gray-100">
                      💰 Aplicar pago total
                    </h4>
                    <p className="text-sm text-fr-gray-600 dark:text-gray-400 mt-1">
                      Marcar como pagado con el monto original de {formatCurrency(payingExpense.amount)}
                    </p>
                  </div>
                </div>
              </button>
            </div>

            {/* Botón cancelar */}
            <div className="flex space-x-4 pt-6">
              <button
                type="button"
                onClick={() => {
                  setShowOverpaymentModal(false);
                }}
                className="btn-outline flex-1"
              >
                Cancelar
              </button>
            </div>
          </div>
        </div>,
        document.body
      )}

      {/* Modal de Confirmación de Eliminación */}
      <ConfirmationModal
        isOpen={showDeleteModal}
        onClose={cancelDelete}
        onConfirm={confirmDelete}
        title="Eliminar Gasto"
        message={`¿Estás seguro de que quieres eliminar el gasto "${deletingExpense?.description}"? Esta acción no se puede deshacer.`}
        confirmText="Eliminar"
        cancelText="Cancelar"
        type="danger"
        loading={deleteLoading}
      />
    </div>
  );
};

export default Expenses; 