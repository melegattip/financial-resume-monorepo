import React, { useState, useEffect } from 'react';
import { FaCalendar, FaDollarSign } from 'react-icons/fa';
import { recurringTransactionsAPI, categoriesAPI, formatCurrency } from '../services/api';
import { usePeriod } from '../contexts/PeriodContext';
import { formatAmount } from '../utils/formatters';
import toast from '../utils/notifications';
import dataService from '../services/dataService';
import ResponsiveTable from '../components/ResponsiveTable';

const RecurringTransactions = () => {
  const { balancesHidden } = usePeriod();
  const [transactions, setTransactions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [dashboard, setDashboard] = useState(null);
  const [projection, setProjection] = useState(null);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showProjectionModal, setShowProjectionModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState(null);
  const [deletingTransaction, setDeletingTransaction] = useState(null);
  const [projectionMonths, setProjectionMonths] = useState(6);
  const [filters, setFilters] = useState({
    type: '',
    frequency: '',
    status: '',
    category_id: '',
    sort_by: 'next_execution_date',
    sort_order: 'asc'
  });

  const getDefaultDate = () => {
    const nextMonth = new Date();
    nextMonth.setMonth(nextMonth.getMonth() + 1);
    nextMonth.setDate(1);
    return nextMonth.toISOString().split('T')[0];
  };

  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    type: 'expense',
    frequency: 'monthly',
    category_id: '',
    next_date: getDefaultDate(),
    end_date: '',
    is_active: true
  });

  useEffect(() => {
    loadData();
  }, [filters]);

  const loadData = async () => {
    try {
      setLoading(true);
      const [transactionsRes, categoriesRes, dashboardRes] = await Promise.all([
        recurringTransactionsAPI.list(filters),
        categoriesAPI.list(),
        recurringTransactionsAPI.getDashboard()
      ]);
      
      console.log('📊 Dashboard response:', dashboardRes.data);
      console.log('📊 Dashboard summary:', dashboardRes.data.data?.summary);
      
      setTransactions(transactionsRes.data.data?.transactions || []);
      setCategories(categoriesRes.data.data || []);
      setDashboard(dashboardRes.data.data);
    } catch (error) {
      console.error('Error loading recurring transactions:', error);
      toast.error('Error cargando transacciones recurrentes');
    } finally {
      setLoading(false);
    }
  };

  const loadProjection = async (months = projectionMonths) => {
    try {
      const projectionRes = await recurringTransactionsAPI.getProjection(months);
      setProjection(projectionRes.data.data);
      setShowProjectionModal(true);
    } catch (error) {
      console.error('Error loading projection:', error);
      toast.error('Error cargando proyección');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validación del frontend
    if (!formData.description || formData.description.trim() === '') {
      toast.error('La descripción es requerida');
      return;
    }
    
    if (!formData.amount || parseFloat(formData.amount) <= 0) {
      toast.error('El monto debe ser mayor a 0');
      return;
    }
    
    if (!formData.next_date) {
      toast.error('La fecha de próxima ejecución es requerida');
      return;
    }
    
    // Validar que la fecha de próxima ejecución no sea en el pasado
    const nextDateString = formData.next_date;
    const todayString = new Date().toISOString().split('T')[0];
    
    if (nextDateString < todayString) {
      toast.error('La fecha de próxima ejecución no puede ser anterior a hoy');
      return;
    }
    
    // Validar fecha de fin si se proporciona
    if (formData.end_date) {
      if (formData.end_date <= nextDateString) {
        toast.error('La fecha de fin debe ser posterior a la fecha de próxima ejecución');
        return;
      }
    }
    
    try {
      const data = {
        description: formData.description.trim(),
        amount: parseFloat(formData.amount),
        type: formData.type,
        frequency: formData.frequency,
        category_id: formData.category_id || undefined,
        next_date: formData.next_date,
        auto_create: true, // Por defecto habilitado
        notify_before: 1, // Notificar 1 día antes por defecto
        end_date: formData.end_date || undefined,
        max_executions: undefined // No implementado en el frontend aún
      };

      if (editingTransaction) {
        await recurringTransactionsAPI.update(editingTransaction.id, data);
        toast.success('Transacción recurrente actualizada exitosamente');
      } else {
        await recurringTransactionsAPI.create(data);
        toast.success('Transacción recurrente creada exitosamente');
      }

      setShowModal(false);
      setEditingTransaction(null);
      resetForm();
      loadData();
    } catch (error) {
      console.error('Error saving recurring transaction:', error);
      const errorMessage = error.response?.data?.error || error.message || 'Error guardando transacción recurrente';
      toast.error(errorMessage);
    }
  };

  const handleEdit = (transaction) => {
    setEditingTransaction(transaction);
    setFormData({
      description: transaction.description || '',
      amount: transaction.amount.toString(),
      type: transaction.type,
      frequency: transaction.frequency,
      category_id: transaction.category_id || '',
      next_date: transaction.next_date,
      end_date: transaction.end_date || '',
      is_active: transaction.is_active
    });
    setShowModal(true);
  };

  const handleDelete = (transaction) => {
    setDeletingTransaction(transaction);
    setShowDeleteModal(true);
  };

  const confirmDelete = async () => {
    if (deletingTransaction) {
      try {
        await recurringTransactionsAPI.delete(deletingTransaction.id);
        toast.success('Transacción recurrente eliminada exitosamente');
        loadData();
      } catch (error) {
        console.error('Error deleting recurring transaction:', error);
        toast.error('Error eliminando transacción recurrente');
      } finally {
        setShowDeleteModal(false);
        setDeletingTransaction(null);
      }
    }
  };

  const cancelDelete = () => {
    setShowDeleteModal(false);
    setDeletingTransaction(null);
  };

  const handlePause = async (id) => {
    try {
      await recurringTransactionsAPI.pause(id);
      toast.success('Transacción pausada exitosamente');
      loadData();
    } catch (error) {
      console.error('Error pausing transaction:', error);
      toast.error('Error pausando transacción');
    }
  };

  const handleResume = async (id) => {
    try {
      await recurringTransactionsAPI.resume(id);
      toast.success('Transacción reanudada exitosamente');
      loadData();
    } catch (error) {
      console.error('Error resuming transaction:', error);
      toast.error('Error reanudando transacción');
    }
  };

  const handleProcessPending = async () => {
    try {
      setLoading(true);
      const response = await recurringTransactionsAPI.processPending();
      const result = response.data;
      
      if (result.success_count > 0) {
        toast.success(`✅ ${result.success_count} transacciones ejecutadas exitosamente`);
        
        // Recargar datos locales primero
        await loadData();
        
        // Invalidar cache inmediatamente
        dataService.invalidateAfterMutation('recurring_transaction');
        
        // Forzar actualización adicional con delay para asegurar sincronización
        setTimeout(() => {
          console.log('🔄 Forzando actualización adicional después de ejecutar transacciones pendientes');
          dataService.invalidateAfterMutation('recurring_transaction');
        }, 1500);
      }
      
      if (result.failure_count > 0) {
        toast.error(`❌ ${result.failure_count} transacciones fallaron`);
      }
      
      if (result.processed_count === 0) {
        toast.info('ℹ️ No hay transacciones pendientes por ejecutar');
      }
      
    } catch (error) {
      console.error('Error processing pending transactions:', error);
      toast.error('Error procesando transacciones pendientes');
    } finally {
      setLoading(false);
    }
  };

  const handleExecuteTransaction = async (id) => {
    try {
      const response = await recurringTransactionsAPI.execute(id);
      const result = response.data.data;
      
      if (result.success) {
        toast.success(`✅ Transacción ejecutada exitosamente`);
        if (result.next_execution_date) {
          toast.info(`📅 Próxima ejecución: ${formatDate(result.next_execution_date)}`);
        }
        
        // Recargar datos locales primero
        await loadData();
        
        // Invalidar cache inmediatamente
        dataService.invalidateAfterMutation('recurring_transaction');
        
        // Forzar actualización adicional con delay para asegurar sincronización
        setTimeout(() => {
          console.log('🔄 Forzando actualización adicional después de ejecutar transacción individual');
          dataService.invalidateAfterMutation('recurring_transaction');
        }, 1500);
      } else {
        toast.error(`❌ Error: ${result.message}`);
      }
      
    } catch (error) {
      console.error('Error executing transaction:', error);
      toast.error('Error ejecutando transacción');
    }
  };

  const resetForm = () => {
    setFormData({
      description: '',
      amount: '',
      type: 'expense',
      frequency: 'monthly',
      category_id: '',
      next_date: getDefaultDate(),
      end_date: '',
      is_active: true
    });
  };

  const getTypeColor = (type) => {
    return type === 'income' 
      ? 'text-green-700 bg-green-100 dark:text-green-300 dark:bg-green-900/30' 
      : 'text-red-700 bg-red-100 dark:text-red-300 dark:bg-red-900/30';
  };

  const getTypeText = (type) => {
    return type === 'income' ? 'Ingreso' : 'Gasto';
  };

  const getFrequencyText = (frequency) => {
    const frequencies = {
      daily: 'Diaria',
      weekly: 'Semanal',
      monthly: 'Mensual',
      yearly: 'Anual'
    };
    return frequencies[frequency] || frequency;
  };

  const getCategoryName = (categoryId) => {
    const category = categories.find(c => c.id === categoryId);
    return category ? category.name : 'Sin categoría';
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    
    // Manejar tanto fechas ISO como fechas simples
    let date;
    if (dateString.includes('T')) {
      // Formato ISO
      date = new Date(dateString);
    } else {
      // Formato YYYY-MM-DD
      date = new Date(dateString + 'T00:00:00');
    }
    
    return new Intl.DateTimeFormat('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date);
  };

  const formatDateTime = (dateString) => {
    if (!dateString) return '';
    
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('es-ES', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  };

  const getDaysUntilNext = (nextDate) => {
    const today = new Date();
    const next = new Date(nextDate);
    const diffTime = next - today;
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays < 0) return 'Vencida';
    if (diffDays === 0) return 'Hoy';
    if (diffDays === 1) return 'Mañana';
    return `${diffDays} días`;
  };

  const isTransactionOverdue = (nextDate) => {
    const today = new Date();
    const next = new Date(nextDate);
    return next <= today;
  };

  const calculateMonthlyEquivalent = (transaction) => {
    switch (transaction.frequency) {
      case 'daily':
        return transaction.amount * 30;
      case 'weekly':
        return transaction.amount * 4.33;
      case 'monthly':
        return transaction.amount;
      case 'yearly':
        return transaction.amount / 12;
      default:
        return transaction.amount;
    }
  };

  // Configuración de columnas para ResponsiveTable
  const tableColumns = [
    {
      header: 'Transacción',
      accessor: 'description',
      render: (value) => (
        <div className="font-medium text-gray-900 dark:text-gray-100">
          {value}
        </div>
      )
    },
    {
      header: 'Tipo',
      accessor: 'type',
      hideOnTablet: true,
      render: (value) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getTypeColor(value)}`}>
          {getTypeText(value)}
        </span>
      )
    },
    {
      header: 'Monto',
      accessor: 'amount',
      align: 'right',
      render: (value) => formatAmount(value, balancesHidden)
    },
    {
      header: 'Frecuencia',
      accessor: 'frequency',
      hideOnMobile: true,
      render: (value) => getFrequencyText(value)
    },
    {
      header: 'Categoría',
      accessor: 'category_id',
      hideOnTablet: true,
      render: (value) => getCategoryName(value)
    },
    // Mostrar fecha de creación solo cuando se ordena por ella
    ...(filters.sort_by === 'created_at' ? [{
      header: 'Fecha de Creación',
      accessor: 'created_at',
      hideOnMobile: true,
      render: (value) => (
        <div className="text-sm text-gray-600 dark:text-gray-400">
          {formatDateTime(value)}
        </div>
      )
    }] : []),
    {
      header: 'Próxima Ejecución',
      accessor: 'next_date',
      render: (value, item) => (
        <div>
          <div className="text-sm text-gray-900 dark:text-gray-100">
            {formatDate(value)}
          </div>
          <div className={`text-xs font-medium ${
            isTransactionOverdue(value) 
              ? 'text-red-600 dark:text-red-400' 
              : 'text-gray-500 dark:text-gray-400'
          }`}>
            {isTransactionOverdue(value) && '⚠️ '}
            {getDaysUntilNext(value)}
          </div>
        </div>
      )
    },
    {
      header: 'Estado',
      accessor: 'is_active',
      hideOnTablet: true,
      render: (value) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
          value 
            ? 'text-green-700 bg-green-100 dark:text-green-300 dark:bg-green-900/30' 
            : 'text-yellow-700 bg-yellow-100 dark:text-yellow-300 dark:bg-yellow-900/30'
        }`}>
          {value ? 'Activa' : 'Pausada'}
        </span>
      )
    },
    {
      header: 'Acciones',
      accessor: 'actions',
      align: 'right',
      render: (_, transaction) => (
        <div className="flex justify-end space-x-2">
          {/* Botón de ejecutar para transacciones vencidas y activas */}
          {transaction.is_active && isTransactionOverdue(transaction.next_date) && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                handleExecuteTransaction(transaction.id);
              }}
              className="text-purple-600 hover:text-purple-900 dark:text-purple-400 dark:hover:text-purple-300 p-1"
              title="Ejecutar ahora (vencida)"
            >
              ⚡
            </button>
          )}
          {transaction.is_active ? (
            <button
              onClick={(e) => {
                e.stopPropagation();
                handlePause(transaction.id);
              }}
              className="text-yellow-600 hover:text-yellow-900 dark:text-yellow-400 dark:hover:text-yellow-300 p-1"
              title="Pausar"
            >
              ⏸️
            </button>
          ) : (
            <button
              onClick={(e) => {
                e.stopPropagation();
                handleResume(transaction.id);
              }}
              className="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300 p-1"
              title="Reanudar"
            >
              ▶️
            </button>
          )}
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleEdit(transaction);
            }}
            className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 p-1"
            title="Editar"
          >
            ⚙️
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleDelete(transaction);
            }}
            className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 p-1"
            title="Eliminar"
          >
            🗑️
          </button>
        </div>
      )
    }
  ];

  // Renderizado personalizado para cards móviles
  const renderMobileCard = (transaction, index) => (
    <div 
      key={index}
      className="card border-l-4 border-fr-primary transition-all duration-200 hover:shadow-lg"
    >
      <div className="space-y-3">
        {/* Header de la card */}
        <div className="flex justify-between items-start">
          <div className="flex-1 min-w-0">
            <h3 className="font-medium text-gray-900 dark:text-gray-100 text-sm">
              {transaction.description}
            </h3>
            <div className="flex items-center space-x-2 mt-1">
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getTypeColor(transaction.type)}`}>
                {getTypeText(transaction.type)}
              </span>
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                transaction.is_active 
                  ? 'text-green-700 bg-green-100 dark:text-green-300 dark:bg-green-900/30' 
                  : 'text-yellow-700 bg-yellow-100 dark:text-yellow-300 dark:bg-yellow-900/30'
              }`}>
                {transaction.is_active ? 'Activa' : 'Pausada'}
              </span>
            </div>
          </div>
          <div className="text-right">
            <div className="text-lg font-bold text-gray-900 dark:text-gray-100">
              {formatAmount(transaction.amount, balancesHidden)}
            </div>
          </div>
        </div>

        {/* Detalles */}
        <div className="grid grid-cols-2 gap-3 text-sm">
          <div>
            <span className="text-gray-600 dark:text-gray-400">Frecuencia:</span>
            <div className="font-medium text-gray-900 dark:text-gray-100">
              {getFrequencyText(transaction.frequency)}
            </div>
          </div>
          <div>
            <span className="text-gray-600 dark:text-gray-400">Categoría:</span>
            <div className="font-medium text-gray-900 dark:text-gray-100">
              {getCategoryName(transaction.category_id)}
            </div>
          </div>
        </div>

        {/* Fecha de creación - solo mostrar cuando se ordena por ella */}
        {filters.sort_by === 'created_at' && (
          <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-3">
            <div className="flex items-center space-x-2">
              <span className="text-blue-600 dark:text-blue-400">📅</span>
              <div>
                <span className="text-xs text-gray-600 dark:text-gray-400">Fecha de creación:</span>
                <div className="font-medium text-gray-900 dark:text-gray-100">
                  {formatDateTime(transaction.created_at)}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Próxima ejecución */}
        <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
          <div className="flex justify-between items-center">
            <div>
              <span className="text-xs text-gray-600 dark:text-gray-400">Próxima ejecución:</span>
              <div className="font-medium text-gray-900 dark:text-gray-100">
                {formatDate(transaction.next_date)}
              </div>
            </div>
            <div className={`text-xs font-medium ${
              isTransactionOverdue(transaction.next_date) 
                ? 'text-red-600 dark:text-red-400' 
                : 'text-gray-500 dark:text-gray-400'
            }`}>
              {isTransactionOverdue(transaction.next_date) && '⚠️ '}
              {getDaysUntilNext(transaction.next_date)}
            </div>
          </div>
        </div>

        {/* Acciones */}
        <div className="flex justify-end space-x-3 pt-2 border-t border-gray-200 dark:border-gray-600">
          {transaction.is_active && isTransactionOverdue(transaction.next_date) && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                handleExecuteTransaction(transaction.id);
              }}
              className="flex items-center space-x-1 px-3 py-1 text-xs font-medium text-purple-600 hover:text-purple-900 dark:text-purple-400 dark:hover:text-purple-300 bg-purple-50 dark:bg-purple-900/20 rounded-full transition-colors"
            >
              <span>⚡</span>
              <span>Ejecutar</span>
            </button>
          )}
          {transaction.is_active ? (
            <button
              onClick={(e) => {
                e.stopPropagation();
                handlePause(transaction.id);
              }}
              className="flex items-center space-x-1 px-3 py-1 text-xs font-medium text-yellow-600 hover:text-yellow-900 dark:text-yellow-400 dark:hover:text-yellow-300 bg-yellow-50 dark:bg-yellow-900/20 rounded-full transition-colors"
            >
              <span>⏸️</span>
              <span>Pausar</span>
            </button>
          ) : (
            <button
              onClick={(e) => {
                e.stopPropagation();
                handleResume(transaction.id);
              }}
              className="flex items-center space-x-1 px-3 py-1 text-xs font-medium text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300 bg-green-50 dark:bg-green-900/20 rounded-full transition-colors"
            >
              <span>▶️</span>
              <span>Reanudar</span>
            </button>
          )}
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleEdit(transaction);
            }}
            className="flex items-center space-x-1 px-3 py-1 text-xs font-medium text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 bg-blue-50 dark:bg-blue-900/20 rounded-full transition-colors"
          >
            <span>⚙️</span>
            <span>Editar</span>
          </button>
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleDelete(transaction);
            }}
            className="flex items-center space-x-1 px-3 py-1 text-xs font-medium text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 bg-red-50 dark:bg-red-900/20 rounded-full transition-colors"
          >
            <span>🗑️</span>
            <span>Eliminar</span>
          </button>
        </div>
      </div>
    </div>
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-gray-600 dark:text-gray-400">Cargando transacciones recurrentes...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Page Title */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-fr-gray-900 dark:text-gray-100">Transacciones Recurrentes</h1>
      </div>

      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <p className="text-gray-600 dark:text-gray-400">Gestiona tus ingresos y gastos automáticos</p>
        </div>
        <div className="flex space-x-3">
          {/* Botón discreto para desarrollo - ejecutar pendientes */}
          <button
            onClick={handleProcessPending}
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 px-2 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-sm"
            title="Ejecutar transacciones pendientes (desarrollo)"
            disabled={loading}
          >
            {loading ? '⏳' : '🔄'}
          </button>
          <button
            onClick={loadProjection}
            className="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 transition-colors"
          >
            Ver Proyección
          </button>
          <button
            onClick={() => setShowModal(true)}
            className="btn-primary"
          >
            Nueva Transacción
          </button>
        </div>
      </div>

      {/* Dashboard Summary */}
      {dashboard && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="card">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-6">
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Total Activas</p>
                    <p className="text-2xl font-bold text-gray-900 dark:text-gray-100">{dashboard.summary?.total_active || 0}</p>
                  </div>
                  <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Inactivas</p>
                    <p className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{dashboard.summary?.total_inactive || 0}</p>
                  </div>
                </div>
              </div>
              <div className="flex-shrink-0 p-3 rounded-fr bg-blue-100 dark:bg-blue-900/30 ml-4">
                <FaCalendar className="w-6 h-6 text-blue-600 dark:text-blue-400" />
              </div>
            </div>
          </div>
          <div className="card">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <div className="flex items-center space-x-6">
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Ingresos</p>
                    <p className="text-2xl font-bold text-green-600 dark:text-green-400">
                      {formatAmount(dashboard.summary?.monthly_income_total || 0, balancesHidden)}
                    </p>
                  </div>
                  <div className="h-12 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
                  <div>
                    <p className="text-sm font-medium text-gray-500 dark:text-gray-400">Gastos</p>
                    <p className="text-2xl font-bold text-red-600 dark:text-red-400">
                      {formatAmount(dashboard.summary?.monthly_expense_total || 0, balancesHidden)}
                    </p>
                  </div>
                </div>
              </div>
              <div className="flex-shrink-0 p-3 rounded-fr bg-purple-100 dark:bg-purple-900/30 ml-4">
                <FaDollarSign className="w-6 h-6 text-purple-600 dark:text-purple-400" />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="card">
        <div className="grid grid-cols-1 md:grid-cols-6 gap-4">
          <select
            value={filters.type}
            onChange={(e) => setFilters({...filters, type: e.target.value})}
            className="input"
          >
            <option value="">Todos los tipos</option>
            <option value="income">Ingresos</option>
            <option value="expense">Gastos</option>
          </select>
          
          <select
            value={filters.frequency}
            onChange={(e) => setFilters({...filters, frequency: e.target.value})}
            className="input"
          >
            <option value="">Todas las frecuencias</option>
            <option value="daily">Diaria</option>
            <option value="weekly">Semanal</option>
            <option value="monthly">Mensual</option>
            <option value="yearly">Anual</option>
          </select>

          <select
            value={filters.status}
            onChange={(e) => setFilters({...filters, status: e.target.value})}
            className="input"
          >
            <option value="">Todos los estados</option>
            <option value="active">Activa</option>
            <option value="paused">Pausada</option>
          </select>

          <select
            value={filters.category_id}
            onChange={(e) => setFilters({...filters, category_id: e.target.value})}
            className="input"
          >
            <option value="">Todas las categorías</option>
            {categories.map(category => (
              <option key={category.id} value={category.id}>{category.name}</option>
            ))}
          </select>

          <select
            value={filters.sort_by}
            onChange={(e) => setFilters({...filters, sort_by: e.target.value})}
            className="input"
          >
            <option value="next_execution_date">Próxima ejecución</option>
            <option value="description">Descripción</option>
            <option value="amount">Monto</option>
            <option value="created_at">Fecha de creación</option>
          </select>

          <select
            value={filters.sort_order}
            onChange={(e) => setFilters({...filters, sort_order: e.target.value})}
            className="input"
          >
            <option value="asc">Ascendente</option>
            <option value="desc">Descendente</option>
          </select>
        </div>
      </div>

      {/* Transactions List */}
      <div className="space-y-4">
        <ResponsiveTable
          data={transactions}
          columns={tableColumns}
          renderMobileCard={renderMobileCard}
          loading={loading}
          emptyMessage="No hay transacciones recurrentes. Crea tu primera transacción recurrente para automatizar tus finanzas."
          className="card"
        />
        
        {/* Botón para crear transacción cuando no hay datos */}
        {transactions.length === 0 && !loading && (
          <div className="text-center">
            <button
              onClick={() => setShowModal(true)}
              className="btn-primary"
            >
              Crear Transacción
            </button>
          </div>
        )}
      </div>

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md max-h-screen overflow-y-auto">
            <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
              {editingTransaction ? 'Editar Transacción Recurrente' : 'Nueva Transacción Recurrente'}
            </h2>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Descripción
                </label>
                <input
                  type="text"
                  value={formData.description}
                  onChange={(e) => setFormData({...formData, description: e.target.value})}
                  className="input"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Monto
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.amount}
                  onChange={(e) => setFormData({...formData, amount: e.target.value})}
                  className="input"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Tipo
                </label>
                <select
                  value={formData.type}
                  onChange={(e) => setFormData({...formData, type: e.target.value})}
                  className="input"
                  required
                >
                  <option value="expense">Gasto</option>
                  <option value="income">Ingreso</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Frecuencia
                </label>
                <select
                  value={formData.frequency}
                  onChange={(e) => setFormData({...formData, frequency: e.target.value})}
                  className="input"
                  required
                >
                  <option value="daily">Diaria</option>
                  <option value="weekly">Semanal</option>
                  <option value="monthly">Mensual</option>
                  <option value="yearly">Anual</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Categoría
                </label>
                <select
                  value={formData.category_id}
                  onChange={(e) => setFormData({...formData, category_id: e.target.value})}
                  className="input"
                >
                  <option value="">Seleccionar categoría</option>
                  {categories.map(category => (
                    <option key={category.id} value={category.id}>{category.name}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Próxima ejecución
                </label>
                <input
                  type="date"
                  value={formData.next_date}
                  onChange={(e) => setFormData({...formData, next_date: e.target.value})}
                  className="input"
                  min={new Date().toISOString().split('T')[0]}
                  required
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Puede ser hoy o cualquier fecha futura
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Fecha de fin (opcional)
                </label>
                <input
                  type="date"
                  value={formData.end_date}
                  onChange={(e) => setFormData({...formData, end_date: e.target.value})}
                  className="input"
                  min={formData.next_date || new Date().toISOString().split('T')[0]}
                />
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Debe ser posterior a la fecha de próxima ejecución
                </p>
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="is_active"
                  checked={formData.is_active}
                  onChange={(e) => setFormData({...formData, is_active: e.target.checked})}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label htmlFor="is_active" className="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                  Activa
                </label>
              </div>

              <div className="flex justify-end space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false);
                    setEditingTransaction(null);
                    resetForm();
                  }}
                  className="btn-secondary"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="btn-primary"
                >
                  {editingTransaction ? 'Actualizar' : 'Crear'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Projection Modal */}
      {showProjectionModal && projection && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-6xl max-h-screen overflow-y-auto">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100">
                📊 Proyección de Flujo de Caja - {projectionMonths} Meses
              </h2>
              <button
                onClick={() => setShowProjectionModal(false)}
                className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
              >
                ✕
              </button>
            </div>
            
            <div className="mb-4 p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
              <p className="text-sm text-blue-800 dark:text-blue-300">
                💡 Esta proyección se basa en tus transacciones recurrentes activas y calcula el flujo de caja esperado.
                Los cálculos consideran las frecuencias reales de cada transacción (diaria, semanal, mensual, anual).
              </p>
            </div>

            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
              <div className="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
                <h3 className="text-sm font-medium text-green-800 dark:text-green-300 mb-1">
                  Ingresos Mensuales Promedio
                </h3>
                <p className="text-2xl font-bold text-green-600 dark:text-green-400">
                  {formatAmount(projection.summary?.average_monthly_income || 0, balancesHidden)}
                </p>
              </div>
              
              <div className="bg-red-50 dark:bg-red-900/20 rounded-lg p-4">
                <h3 className="text-sm font-medium text-red-800 dark:text-red-300 mb-1">
                  Gastos Mensuales Promedio
                </h3>
                <p className="text-2xl font-bold text-red-600 dark:text-red-400">
                  {formatAmount(projection.summary?.average_monthly_expenses || 0, balancesHidden)}
                </p>
              </div>
              
              <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
                <h3 className="text-sm font-medium text-blue-800 dark:text-blue-300 mb-1">
                  Balance Neto Mensual
                </h3>
                <p className={`text-2xl font-bold ${
                  (projection.summary?.net_projected_amount || 0) >= 0 
                    ? 'text-green-600 dark:text-green-400' 
                    : 'text-red-600 dark:text-red-400'
                }`}>
                  {formatAmount(projection.summary?.net_projected_amount || 0, balancesHidden)}
                </p>
              </div>
            </div>

            {/* Monthly Projections Table */}
            <div className="bg-white dark:bg-gray-700 rounded-lg overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-600">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
                  Proyección Mensual Detallada
                </h3>
              </div>
              
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-600">
                  <thead className="bg-gray-50 dark:bg-gray-600">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                        Mes
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                        Ingresos
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                        Gastos
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                        Balance Neto
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                        Acumulado
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                    {projection.monthly_projections?.map((month, index) => (
                      <tr key={index} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-gray-100">
                          {month.month_display}
                        </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600 dark:text-green-400">
                          {formatAmount(month.income, balancesHidden)}
                        </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600 dark:text-red-400">
                          {formatAmount(month.expenses, balancesHidden)}
                        </td>
                        <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${
                          month.net_amount >= 0 
                            ? 'text-green-600 dark:text-green-400' 
                            : 'text-red-600 dark:text-red-400'
                        }`}>
                          {formatAmount(month.net_amount, balancesHidden)}
                        </td>
                        <td className={`px-6 py-4 whitespace-nowrap text-sm font-bold ${
                          month.cumulative_net >= 0 
                            ? 'text-green-600 dark:text-green-400' 
                            : 'text-red-600 dark:text-red-400'
                        }`}>
                          {formatAmount(month.cumulative_net, balancesHidden)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            {/* Insights */}
            <div className="mt-6 p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
              <h4 className="text-sm font-medium text-yellow-800 dark:text-yellow-300 mb-2">
                💡 Insights
              </h4>
              <ul className="text-sm text-yellow-700 dark:text-yellow-400 space-y-1">
                <li>• Esta proyección considera solo transacciones recurrentes activas</li>
                <li>• Los gastos únicos no están incluidos en esta proyección</li>
                <li>• Revisa regularmente tus transacciones recurrentes para mantener la proyección actualizada</li>
                <li>• Considera ajustar frecuencias o montos según tus necesidades financieras</li>
              </ul>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {showDeleteModal && deletingTransaction && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center mb-4">
              <div className="flex-shrink-0 w-10 h-10 mx-auto bg-red-100 dark:bg-red-900/30 rounded-full flex items-center justify-center">
                <svg className="w-6 h-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 18.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
            </div>
            
            <div className="text-center">
              <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                Eliminar Transacción Recurrente
              </h3>
              <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
                ¿Estás seguro de que quieres eliminar la transacción "{deletingTransaction.description}"? 
                Esta acción no se puede deshacer.
              </p>
            </div>

            <div className="flex justify-center space-x-3">
              <button
                type="button"
                onClick={cancelDelete}
                className="btn-secondary"
              >
                Cancelar
              </button>
              <button
                type="button"
                onClick={confirmDelete}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 dark:bg-red-700 dark:hover:bg-red-800 transition-colors"
              >
                Eliminar
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default RecurringTransactions;