import React, { useState, useEffect, useCallback } from 'react';
import { savingsGoalsAPI, formatCurrency, formatPercentage } from '../services/api';
import { usePeriod } from '../contexts/PeriodContext';
import { formatAmount } from '../utils/formatters';
import toast from '../utils/notifications';
import TrialBanner from '../components/TrialBanner';
import ConfirmationModal from '../components/ConfirmationModal';

const SavingsGoals = () => {
  const { balancesHidden } = usePeriod();
  const [goals, setGoals] = useState([]);
  const [dashboard, setDashboard] = useState(null);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showTransactionModal, setShowTransactionModal] = useState(false);
  const [editingGoal, setEditingGoal] = useState(null);
  const [transactionGoal, setTransactionGoal] = useState(null);
  const [transactionType, setTransactionType] = useState('deposit');
  const [selectedGoal, setSelectedGoal] = useState(null);
  const [filters] = useState({
    status: '',
    category: '',
    priority: '',
    sort_by: 'created_at',
    sort_order: 'desc'
  });

  // Estado para movimientos de una meta seleccionada
  const [transactions, setTransactions] = useState([]);
  const [txLoading, setTxLoading] = useState(false);

  // Cargar movimientos cuando se abre el detalle de una meta
  useEffect(() => {
    const fetchDetailAndTransactions = async () => {
      if (!selectedGoal?.id) {
        setTransactions([]);
        return;
      }
      try {
        setTxLoading(true);
        // Refrescar detalle para que el current_amount sea coherente con los movimientos
        try {
          const detail = await savingsGoalsAPI.get(selectedGoal.id);
          const fresh = detail?.data?.data || detail?.data;
          if (fresh) setSelectedGoal(fresh);
        } catch (_) {}

        const res = await savingsGoalsAPI.getTransactions(selectedGoal.id, { limit: 100, offset: 0 });
        const txs = res?.data?.data?.transactions || res?.data?.transactions || [];
        setTransactions(txs);
      } catch (error) {
        console.error('Error loading transactions:', error);
      } finally {
        setTxLoading(false);
      }
    };

    fetchDetailAndTransactions();
  }, [selectedGoal?.id]);

  const [formData, setFormData] = useState({
    name: '',
    description: '',
    target_amount: '',
    current_amount: '',
    target_date: '',
    category: 'emergency',
    priority: 'medium',
    auto_save_amount: '',
    auto_save_frequency: 'monthly',
    icon: ''
  });

  const [transactionData, setTransactionData] = useState({
    amount: '',
    description: ''
  });

  // Estados para modal de confirmación de eliminación
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingGoal, setDeletingGoal] = useState(null);
  const [deleteLoading, setDeleteLoading] = useState(false);

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      const [goalsRes, dashboardRes] = await Promise.all([
        savingsGoalsAPI.list(filters),
        savingsGoalsAPI.getDashboard()
      ]);
      
      setGoals(goalsRes.data.data?.goals || []);
      setDashboard(dashboardRes.data.data);
    } catch (error) {
      console.error('Error loading savings goals:', error);
      toast.error('Error cargando metas de ahorro');
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      // Normalizar payload para backend
      const allowedCategories = new Set(['vacation','emergency','house','car','education','retirement','investment','other']);
      const normalizedCategory = (formData.category || 'other').toLowerCase();
      const category = allowedCategories.has(normalizedCategory) ? normalizedCategory : 'other';

      // target_date del input date viene como YYYY-MM-DD → convertir a ISO RFC3339
      let targetDateISO = undefined;
      if (formData.target_date) {
        const d = new Date(`${formData.target_date}T00:00:00`);
        if (!isNaN(d.getTime())) targetDateISO = d.toISOString();
      }

      const targetAmount = parseFloat(formData.target_amount);
      const currentAmount = parseFloat(formData.current_amount || 0);
      const autoSaveAmount = parseFloat(formData.auto_save_amount || 0);
      const priority = (formData.priority || 'medium').toLowerCase();

      // Para creación usamos payload completo; para edición, solo campos modificados para evitar conflictos
      const fullPayload = {
        name: formData.name,
        description: formData.description || '',
        target_amount: targetAmount,
        category,
        priority,
        ...(targetDateISO ? { target_date: targetDateISO } : {}),
        ...(autoSaveAmount > 0 ? { is_auto_save: true, auto_save_amount: autoSaveAmount, auto_save_frequency: formData.auto_save_frequency || 'monthly' } : { is_auto_save: false }),
        ...(formData.icon ? { image_url: `data:text/plain;charset=utf-8,${encodeURIComponent(formData.icon)}` } : {})
      };

      if (editingGoal) {
        // Construir payload parcial solo con cambios
        const partial = {};
        if (editingGoal.name !== fullPayload.name) partial.name = fullPayload.name;
        if ((editingGoal.description || '') !== (fullPayload.description || '')) partial.description = fullPayload.description;
        if (Number(editingGoal.target_amount) !== Number(fullPayload.target_amount)) partial.target_amount = fullPayload.target_amount;
        if (editingGoal.category !== fullPayload.category) partial.category = fullPayload.category;
        if (editingGoal.priority !== fullPayload.priority) partial.priority = fullPayload.priority;
        if (targetDateISO) partial.target_date = targetDateISO;
        if (fullPayload.is_auto_save !== undefined) {
          // Solo enviar si cambió el flag o parámetros asociados
          if (Boolean(editingGoal.is_auto_save) !== Boolean(fullPayload.is_auto_save)) partial.is_auto_save = fullPayload.is_auto_save;
          if (fullPayload.is_auto_save) {
            if (Number(editingGoal.auto_save_amount || 0) !== Number(fullPayload.auto_save_amount || 0)) partial.auto_save_amount = fullPayload.auto_save_amount;
            if ((editingGoal.auto_save_frequency || '') !== (fullPayload.auto_save_frequency || '')) partial.auto_save_frequency = fullPayload.auto_save_frequency;
          }
        }
        if (fullPayload.image_url) partial.image_url = fullPayload.image_url;

        const res = await savingsGoalsAPI.update(editingGoal.id, Object.keys(partial).length ? partial : { name: fullPayload.name });
        const updated = res?.data?.data || null;
        // Refrescar la vista de detalle inmediatamente si estamos viendo esta meta
        if (updated && selectedGoal && selectedGoal.id === editingGoal.id) {
          setSelectedGoal(updated);
        }
        toast.success('Meta de ahorro actualizada exitosamente');
      } else {
        const res = await savingsGoalsAPI.create(fullPayload);
        const createdId = res?.data?.data?.id || res?.data?.id;
        // Si el usuario indicó un monto inicial, lo registramos como depósito
        if (createdId && currentAmount > 0) {
          try {
            await savingsGoalsAPI.deposit(createdId, { amount: currentAmount, description: 'Depósito inicial' });
          } catch (err) {
            console.error('Error realizando depósito inicial:', err);
            toast.error('La meta se creó, pero falló el depósito inicial');
          }
        }
        toast.success('Meta de ahorro creada exitosamente');
      }

      setShowModal(false);
      setEditingGoal(null);
      resetForm();
      loadData();
    } catch (error) {
      console.error('Error saving goal:', error);
      toast.error('Error guardando meta de ahorro');
    }
  };

  const handleTransaction = async (e) => {
    e.preventDefault();
    try {
      const data = {
        amount: parseFloat(transactionData.amount),
        description: transactionData.description
      };

      if (transactionType === 'deposit') {
        await savingsGoalsAPI.deposit(transactionGoal.id, data);
        toast.success('Depósito realizado exitosamente');
      } else {
        await savingsGoalsAPI.withdraw(transactionGoal.id, data);
        toast.success('Retiro realizado exitosamente');
      }

      setShowTransactionModal(false);
      setTransactionGoal(null);
      setTransactionData({ amount: '', description: '' });
      loadData();

      // Si estamos en la vista de detalle de ESTA meta, refrescamos detalle y movimientos
      const justOperatedGoalId = transactionGoal?.id;
      const isInDetailForSameGoal = !!selectedGoal?.id && selectedGoal.id === justOperatedGoalId;
      if (isInDetailForSameGoal && justOperatedGoalId) {
        try {
          const detail = await savingsGoalsAPI.get(justOperatedGoalId);
          const newGoal = detail?.data?.data || detail?.data;
          if (newGoal) setSelectedGoal(newGoal);
        } catch (_) {}
        try {
          const txRes = await savingsGoalsAPI.getTransactions(justOperatedGoalId, { limit: 100, offset: 0 });
          const txs = txRes?.data?.data?.transactions || txRes?.data?.transactions || [];
          setTransactions(txs);
        } catch (_) {}
      }
    } catch (error) {
      console.error('Error processing transaction:', error);
      toast.error('Error procesando transacción');
    }
  };

  const decodeIcon = (imageUrl) => {
    if (!imageUrl) return '';
    try {
      const prefix = 'data:text/plain;charset=utf-8,';
      if (String(imageUrl).startsWith(prefix)) {
        return decodeURIComponent(String(imageUrl).slice(prefix.length));
      }
      return '';
    } catch (_) {
      return '';
    }
  };

  const handleEdit = (goal) => {
    setEditingGoal(goal);
    setFormData({
      name: goal.name,
      description: goal.description || '',
      target_amount: goal.target_amount.toString(),
      current_amount: goal.current_amount.toString(),
      // Para input type=date usamos YYYY-MM-DD
      target_date: goal.target_date ? String(goal.target_date).substring(0, 10) : '',
      category: goal.category,
      priority: goal.priority,
      auto_save_amount: goal.auto_save_amount?.toString() || '',
      auto_save_frequency: goal.auto_save_frequency || 'monthly',
      icon: decodeIcon(goal.image_url)
    });
    setShowModal(true);
  };

  const handleDelete = (goal) => {
    setDeletingGoal(goal);
    setShowDeleteModal(true);
  };

  const confirmDelete = async () => {
    if (!deletingGoal) return;
    
    try {
      setDeleteLoading(true);
      await savingsGoalsAPI.delete(deletingGoal.id);
      toast.success('Meta de ahorro eliminada exitosamente');
      loadData();
    } catch (error) {
      console.error('Error deleting goal:', error);
      toast.error('Error eliminando meta de ahorro');
    } finally {
      // ✅ Siempre cerrar modal y limpiar estado, sin importar si hay errores
      setDeleteLoading(false);
      setShowDeleteModal(false);
      setDeletingGoal(null);
    }
  };

  const cancelDelete = () => {
    setShowDeleteModal(false);
    setDeletingGoal(null);
    setDeleteLoading(false);
  };

  // (Pausar/Reanudar) no utilizados en esta vista actualmente

  const openTransactionModal = (goal, type) => {
    setTransactionGoal(goal);
    setTransactionType(type);
    setShowTransactionModal(true);
  };

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      target_amount: '',
      current_amount: '',
      target_date: '',
      // Usar categoría válida por defecto
      category: 'emergency',
      priority: 'medium',
      auto_save_amount: '',
      auto_save_frequency: 'monthly',
      icon: ''
    });
  };

  // Helpers de estado/prioridad removidos por no usarse actualmente

  const getCategoryText = (category) => {
    const icon = iconOptions.find(option => option.value === category);
    return icon ? icon.label : category;
  };

  // Opciones de iconos disponibles (alineadas con categorías válidas del backend)
  const iconOptions = [
    { value: 'car', emoji: '🚗', label: 'Auto' },
    { value: 'house', emoji: '🏠', label: 'Casa' },
    { value: 'vacation', emoji: '✈️', label: 'Vacaciones' },
    { value: 'education', emoji: '📚', label: 'Educación' },
    { value: 'emergency', emoji: '🏥', label: 'Emergencia' },
    { value: 'investment', emoji: '📈', label: 'Inversión' },
    { value: 'retirement', emoji: '🧓', label: 'Retiro' },
    { value: 'other', emoji: '💰', label: 'Otro' }
  ];

  const getCategoryIcon = (category) => {
    const icon = iconOptions.find(option => option.value === category);
    return icon ? icon.emoji : '💰';
  };

  const getGoalDisplayIcon = (goal) => {
    const icon = decodeIcon(goal.image_url);
    return icon || getCategoryIcon(goal.category);
  };

  const getProgressPercent = (progress, current, target) => {
    let p = 0;
    if (typeof progress === 'number') {
      p = progress * 100;
    } else if (Number(target) > 0) {
      p = (Number(current) / Number(target)) * 100;
    }
    if (!isFinite(p) || isNaN(p)) p = 0;
    return Math.max(0, Math.min(100, Math.round(p)));
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-gray-600 dark:text-gray-400">Cargando metas de ahorro...</span>
      </div>
    );
  }

  // Vista detalle de una meta específica
  if (selectedGoal) {
    return (
      <>
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
          <div className="bg-white dark:bg-gray-800 rounded-2xl mx-4 mt-4 p-6 shadow-sm border dark:border-gray-700">
            {/* Header con botón de regreso */}
            <div className="flex items-center mb-6">
              <button 
                onClick={() => setSelectedGoal(null)}
                className="flex items-center text-blue-500 dark:text-blue-400 font-medium"
              >
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                Ir a mis metas
              </button>
            </div>

            {/* Título y imagen */}
            <div className="text-center mb-8">
              <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">{selectedGoal.name}</h1>
              <div className="text-6xl mb-6">{getGoalDisplayIcon(selectedGoal)}</div>
              
              <div className="mb-4">
                <div className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-2">
                  {formatAmount(selectedGoal.current_amount, balancesHidden)}
                </div>
              </div>

              {/* Barra de progreso grande */}
              <div className="max-w-xl mx-auto w-full">
                <div className="w-full h-4 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-gradient-to-r from-blue-500 to-green-500"
                    style={{ width: `${getProgressPercent(selectedGoal.progress, selectedGoal.current_amount, selectedGoal.target_amount)}%` }}
                  />
                </div>
                <div className="flex justify-between text-sm mt-2 text-gray-600 dark:text-gray-400">
                  <span>{formatPercentage(getProgressPercent(selectedGoal.progress, selectedGoal.current_amount, selectedGoal.target_amount))}</span>
                  <span>
                    {formatAmount(selectedGoal.current_amount, balancesHidden)} / {formatAmount(selectedGoal.target_amount, balancesHidden)}
                  </span>
                </div>
              </div>

              {/* Botones de acción */}
              <div className="flex justify-center space-x-8 mb-8">
                <button 
                  onClick={() => openTransactionModal(selectedGoal, 'deposit')}
                  className="flex flex-col items-center relative group"
                  title="Depositar dinero en esta meta"
                >
                  <div className="w-16 h-16 bg-green-500 dark:bg-green-600 rounded-full flex items-center justify-center mb-2 transition-transform group-hover:scale-105">
                    <span className="text-2xl">💰</span>
                  </div>
                  <span className="text-gray-700 dark:text-gray-300 font-medium">Ahorrar</span>
                  {/* Tooltip */}
                  <div className="absolute bottom-full mb-2 px-3 py-2 bg-gray-800 dark:bg-gray-700 text-white dark:text-gray-100 text-sm rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10">
                    Depositar dinero en esta meta
                    <div className="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-800 dark:border-t-gray-700"></div>
                  </div>
                </button>
                
                <button 
                  onClick={() => openTransactionModal(selectedGoal, 'withdraw')}
                  className="flex flex-col items-center relative group"
                  title="Retirar dinero de esta meta"
                >
                  <div className="w-16 h-16 bg-orange-500 dark:bg-orange-600 rounded-full flex items-center justify-center mb-2 transition-transform group-hover:scale-105">
                    <span className="text-2xl">💸</span>
                  </div>
                  <span className="text-gray-700 dark:text-gray-300 font-medium">Retirar</span>
                  {/* Tooltip */}
                  <div className="absolute bottom-full mb-2 px-3 py-2 bg-gray-800 dark:bg-gray-700 text-white dark:text-gray-100 text-sm rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10">
                    Retirar dinero de esta meta
                    <div className="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-800 dark:border-t-gray-700"></div>
                  </div>
                </button>
                
                <button 
                  onClick={() => handleEdit(selectedGoal)}
                  className="flex flex-col items-center relative group"
                  title="Editar configuración de la meta"
                >
                  <div className="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center mb-2 transition-transform group-hover:scale-105">
                    <svg className="w-8 h-8 text-blue-500 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                  </div>
                  <span className="text-gray-700 dark:text-gray-300 font-medium">Configurar</span>
                  {/* Tooltip */}
                  <div className="absolute bottom-full mb-2 px-3 py-2 bg-gray-800 dark:bg-gray-700 text-white dark:text-gray-100 text-sm rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10">
                    Editar configuración de la meta
                    <div className="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-gray-800 dark:border-t-gray-700"></div>
                  </div>
                </button>
              </div>
            </div>
            {/* Historial de movimientos */}
            <div className="mt-6">
              <div className="flex items-center justify-between mb-3">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Movimientos</h2>
                {txLoading && (
                  <span className="text-sm text-gray-500 dark:text-gray-400">Cargando...</span>
                )}
              </div>
              {transactions.length === 0 ? (
                <div className="text-gray-600 dark:text-gray-400 text-sm">No hay movimientos aún</div>
              ) : (
                <ul className="divide-y divide-gray-200 dark:divide-gray-700">
                  {transactions.map((tx) => (
                    <li key={tx.id} className="py-3 flex items-center justify-between">
                      <div className="flex items-center space-x-3">
                        <div className={`w-9 h-9 rounded-full flex items-center justify-center ${tx.type === 'deposit' ? 'bg-green-100 dark:bg-green-900/30' : 'bg-yellow-100 dark:bg-yellow-900/30'}`}>
                          <span className="text-lg">{tx.type === 'deposit' ? '💰' : '💸'}</span>
                        </div>
                        <div>
                          <div className="text-sm font-medium text-gray-900 dark:text-gray-100">
                            {tx.description || (tx.type === 'deposit' ? 'Depósito' : 'Retiro')}
                          </div>
                          <div className="text-xs text-gray-500 dark:text-gray-400">
                            {new Date(tx.created_at).toLocaleString('es-AR')}
                          </div>
                        </div>
                      </div>
                      <div className={`text-sm font-semibold ${tx.type === 'deposit' ? 'text-green-600 dark:text-green-400' : 'text-yellow-600 dark:text-yellow-400'}`}>
                        {tx.type === 'deposit' ? '+' : '-'}{formatAmount(tx.amount, balancesHidden)}
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </div>

        {/* Create/Edit Modal */}
        {showModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md max-h-screen overflow-y-auto mx-4 border dark:border-gray-700">
              <div className="flex items-center justify-between mb-6">
                <button 
                  onClick={() => {
                    setShowModal(false);
                    setEditingGoal(null);
                    resetForm();
                  }}
                  className="flex items-center text-blue-500 dark:text-blue-400 font-medium"
                >
                  <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                  </svg>
                  Volver a mis metas
                </button>
              </div>

              <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6 text-center">
                Configurá tu meta de ahorro
              </h2>

              <div className="text-center mb-6">
                <div className="text-6xl mb-4">{getCategoryIcon(formData.category)}</div>
              </div>
              
              <form onSubmit={handleSubmit} className="space-y-6">
                <div>
                  <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                    Icono
                  </label>
                  <div className="grid grid-cols-4 gap-3 mb-4">
                    {iconOptions.map((option) => (
                      <button
                        key={option.value}
                        type="button"
                        onClick={() => setFormData({...formData, category: option.value})}
                        className={`p-3 rounded-lg border-2 transition-all ${
                          formData.category === option.value
                            ? 'border-blue-500 bg-blue-50 dark:border-blue-400 dark:bg-blue-900/30'
                            : 'border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500'
                        }`}
                      >
                        <div className="text-2xl mb-1">{option.emoji}</div>
                        <div className="text-xs text-gray-600 dark:text-gray-400">{option.label}</div>
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                    Nombre
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-3 text-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                    placeholder="Ej: Viaje a Europa, Auto nuevo, Casa propia"
                    required
                  />
                </div>

                <div>
                  <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                    Meta
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={formData.target_amount}
                    onChange={(e) => setFormData({...formData, target_amount: e.target.value})}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-3 text-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                    placeholder="Ingresá el monto que deseás alcanzar"
                    required
                  />
                </div>

                <div>
                  <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                    Plazo
                  </label>
                  <input
                    type="date"
                    value={formData.target_date}
                    onChange={(e) => setFormData({...formData, target_date: e.target.value})}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-3 text-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                    placeholder="Elegí la fecha para lograr tu objetivo"
                    required
                  />
                </div>

                <button
                  type="submit"
                  className="w-full bg-blue-500 dark:bg-blue-600 text-white py-4 rounded-lg text-lg font-medium hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors"
                >
                  {editingGoal ? 'Actualizar Meta' : 'Crear Meta'}
                </button>
              </form>
            </div>
          </div>
        )}

        {/* Transaction Modal */}
        {showTransactionModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md mx-4 border dark:border-gray-700">
              <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
                {transactionType === 'deposit' ? 'Depositar' : 'Retirar'} - {transactionGoal?.name}
              </h2>
              
              <form onSubmit={handleTransaction} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Monto
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={transactionData.amount}
                    onChange={(e) => setTransactionData({...transactionData, amount: e.target.value})}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Descripción
                  </label>
                  <textarea
                    value={transactionData.description}
                    onChange={(e) => setTransactionData({...transactionData, description: e.target.value})}
                    className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                    rows="3"
                  />
                </div>

                <div className="flex space-x-3">
                  <button
                    type="button"
                    onClick={() => setShowTransactionModal(false)}
                    className="flex-1 bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-200 py-2 rounded-lg hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors"
                  >
                    Cancelar
                  </button>
                  <button
                    type="submit"
                    className={`flex-1 py-2 rounded-lg text-white transition-colors ${
                      transactionType === 'deposit' 
                        ? 'bg-green-500 hover:bg-green-600 dark:bg-green-600 dark:hover:bg-green-700' 
                        : 'bg-yellow-500 hover:bg-yellow-600 dark:bg-yellow-600 dark:hover:bg-yellow-700'
                    }`}
                  >
                    {transactionType === 'deposit' ? 'Depositar' : 'Retirar'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Modal de Confirmación de Eliminación */}
        <ConfirmationModal
          isOpen={showDeleteModal}
          onClose={cancelDelete}
          onConfirm={confirmDelete}
          title="Eliminar Meta de Ahorro"
          message={`¿Estás seguro de que quieres eliminar la meta "${deletingGoal?.name}"? Esta acción no se puede deshacer y se perderá todo el historial asociado.`}
          confirmText="Eliminar"
          cancelText="Cancelar"
          type="danger"
          loading={deleteLoading}
        />
      </>
    );
  }

  // Vista principal de ahorros
  return (
    <>
      <div className="min-h-[calc(100vh-4rem)] bg-gray-50 dark:bg-gray-900 px-4 sm:px-6 lg:px-8 pt-4">
        <TrialBanner featureKey="SAVINGS_GOALS" />
        
        {/* Contenedor principal usando todo el ancho */}
        <div className="w-full">
          {/* Título de la sección */}
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-gray-100">Metas de Ahorro</h1>
          </div>

          {/* Botón de crear meta redondo y centrado */}
          <div className="flex justify-center mb-6">
            <button 
              onClick={() => setShowModal(true)}
              className="w-14 h-14 sm:w-16 sm:h-16 bg-blue-500 dark:bg-blue-600 text-white rounded-full flex items-center justify-center hover:bg-blue-600 dark:hover:bg-blue-700 hover:scale-105 transition-all duration-200 shadow-lg"
            >
              <svg className="w-6 h-6 sm:w-7 sm:h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </button>
          </div>

          {/* Header con total ahorrado */}
          <div className="mb-8 text-center sm:text-left">
            <h2 className="text-lg font-medium text-gray-600 dark:text-gray-400 mb-3">Total ahorrado</h2>
            <div className="mb-8">
              <span className="text-4xl sm:text-5xl font-bold text-gray-900 dark:text-gray-100">
                {formatAmount(dashboard?.summary?.total_saved || 0, balancesHidden)}
              </span>
            </div>
          </div>

                  {/* Lista de metas de ahorro */}
        <div className="space-y-4">
          {goals.length === 0 ? (
            <div className="text-center py-12">
              <div className="text-gray-400 dark:text-gray-500 mb-4">
                <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">No hay metas de ahorro</h3>
              <p className="text-gray-600 dark:text-gray-400 mb-4">Crea tu primera meta para empezar a ahorrar</p>
              <button
                onClick={() => setShowModal(true)}
                className="bg-blue-500 dark:bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors"
              >
                Crear Meta
              </button>
            </div>
          ) : (
            goals.map((goal) => (
              <div 
                key={goal.id} 
                className="p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-750 transition-all duration-200 cursor-pointer shadow-sm hover:shadow-md"
                onClick={() => setSelectedGoal(goal)}
              >
                {/* Header con icono, nombre y botones de acción */}
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center flex-1 min-w-0">
                    <div className="text-2xl mr-3 flex-shrink-0">
                      {getGoalDisplayIcon(goal)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h4 className="font-semibold text-gray-900 dark:text-gray-100 text-base truncate">{goal.name}</h4>
                      <p className="text-xs text-gray-500 dark:text-gray-400">{getCategoryText(goal.category)}</p>
                    </div>
                  </div>
                  
                  {/* Botones de acción */}
                  <div className="flex space-x-1.5 flex-shrink-0">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        openTransactionModal(goal, 'deposit');
                      }}
                      className="w-8 h-8 bg-green-500 dark:bg-green-600 rounded-full flex items-center justify-center hover:scale-110 transition-transform shadow-sm"
                      title="Depositar"
                    >
                      <span className="text-xs">💰</span>
                    </button>
                    
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        openTransactionModal(goal, 'withdraw');
                      }}
                      className="w-8 h-8 bg-orange-500 dark:bg-orange-600 rounded-full flex items-center justify-center hover:scale-110 transition-transform shadow-sm"
                      title="Retirar"
                    >
                      <span className="text-xs">💸</span>
                    </button>
                    
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDelete(goal);
                      }}
                      className="w-8 h-8 bg-red-500 dark:bg-red-600 rounded-full flex items-center justify-center hover:scale-110 transition-transform shadow-sm"
                      title="Eliminar"
                    >
                      <svg className="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                {/* Montos y progreso */}
                <div className="space-y-3">
                  {/* Montos */}
                  <div className="flex items-baseline justify-between">
                    <div className="flex-1">
                      <div className="text-xl font-bold text-gray-900 dark:text-gray-100">
                        {formatAmount(goal.current_amount, balancesHidden)}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        de {formatAmount(goal.target_amount, balancesHidden)}
                      </div>
                    </div>
                    <div className="text-right flex-shrink-0">
                      <div className="text-base font-semibold text-blue-600 dark:text-blue-400">
                        {formatPercentage(getProgressPercent(goal.progress, goal.current_amount, goal.target_amount))}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        completado
                      </div>
                    </div>
                  </div>
                  
                  {/* Barra de progreso */}
                  <div className="w-full">
                    <div className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-gradient-to-r from-blue-500 to-blue-600 dark:from-blue-400 dark:to-blue-500 transition-all duration-500 ease-out rounded-full"
                        style={{ width: `${Math.min(getProgressPercent(goal.progress, goal.current_amount, goal.target_amount), 100)}%` }}
                      />
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
          </div>
        </div>
      </div>

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md max-h-screen overflow-y-auto mx-4 border dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <button 
                onClick={() => {
                  setShowModal(false);
                  setEditingGoal(null);
                  resetForm();
                }}
                className="flex items-center text-blue-500 dark:text-blue-400 font-medium"
              >
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                Volver a mis metas
              </button>
            </div>

            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6 text-center">
              Configurá tu meta de ahorro
            </h2>

            <div className="text-center mb-6">
              <div className="text-6xl mb-4">{formData.icon || getCategoryIcon(formData.category)}</div>
            </div>
            
            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                  <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">Icono</label>
                <div className="grid grid-cols-4 gap-3 mb-4">
                  {iconOptions.map((option) => (
                    <button
                      key={option.value}
                      type="button"
                        onClick={() => setFormData({...formData, category: option.value, icon: option.emoji})}
                      className={`p-3 rounded-lg border-2 transition-all ${
                        formData.category === option.value
                          ? 'border-blue-500 bg-blue-50 dark:border-blue-400 dark:bg-blue-900/30'
                          : 'border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500'
                      }`}
                    >
                      <div className="text-2xl mb-1">{option.emoji}</div>
                      <div className="text-xs text-gray-600 dark:text-gray-400">{option.label}</div>
                    </button>
                  ))}
                </div>
                  <div className="flex items-center space-x-3 mb-2">
                    <span className="text-sm text-gray-600 dark:text-gray-400">Vista previa:</span>
                    <span className="text-2xl">{formData.icon || getCategoryIcon(formData.category)}</span>
                  </div>
              </div>

              <div>
                <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                  Nombre
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-3 text-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  placeholder="Ej: Viaje a Europa, Auto nuevo, Casa propia"
                  required
                />
              </div>

              <div>
                <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                  Meta
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={formData.target_amount}
                  onChange={(e) => setFormData({...formData, target_amount: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-3 text-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  placeholder="Ingresá el monto que deseás alcanzar"
                  required
                />
              </div>

              <div>
                <label className="block text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                  Plazo
                </label>
                <input
                  type="date"
                  value={formData.target_date}
                  onChange={(e) => setFormData({...formData, target_date: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-3 text-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  placeholder="Elegí la fecha para lograr tu objetivo"
                  required
                />
              </div>

              <button
                type="submit"
                className="w-full bg-blue-500 dark:bg-blue-600 text-white py-4 rounded-lg text-lg font-medium hover:bg-blue-600 dark:hover:bg-blue-700 transition-colors"
              >
                {editingGoal ? 'Actualizar Meta' : 'Crear Meta'}
              </button>
            </form>
          </div>
        </div>
      )}

      {/* Transaction Modal */}
      {showTransactionModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md mx-4 border dark:border-gray-700">
            <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
              {transactionType === 'deposit' ? 'Depositar' : 'Retirar'} - {transactionGoal?.name}
            </h2>
            
            <form onSubmit={handleTransaction} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Monto
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={transactionData.amount}
                  onChange={(e) => setTransactionData({...transactionData, amount: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Descripción
                </label>
                <textarea
                  value={transactionData.description}
                  onChange={(e) => setTransactionData({...transactionData, description: e.target.value})}
                  className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                  rows="3"
                />
              </div>

              <div className="flex space-x-3">
                <button
                  type="button"
                  onClick={() => setShowTransactionModal(false)}
                  className="flex-1 bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-200 py-2 rounded-lg hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className={`flex-1 py-2 rounded-lg text-white transition-colors ${
                    transactionType === 'deposit' 
                      ? 'bg-green-500 hover:bg-green-600 dark:bg-green-600 dark:hover:bg-green-700' 
                      : 'bg-yellow-500 hover:bg-yellow-600 dark:bg-yellow-600 dark:hover:bg-yellow-700'
                  }`}
                >
                  {transactionType === 'deposit' ? 'Depositar' : 'Retirar'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Modal de Confirmación de Eliminación */}
      <ConfirmationModal
        isOpen={showDeleteModal}
        onClose={cancelDelete}
        onConfirm={confirmDelete}
        title="Eliminar Meta de Ahorro"
        message={`¿Estás seguro de que quieres eliminar la meta "${deletingGoal?.name}"? Esta acción no se puede deshacer y se perderá todo el historial asociado.`}
        confirmText="Eliminar"
        cancelText="Cancelar"
        type="danger"
        loading={deleteLoading}
      />
    </>
  );
};

export default SavingsGoals; 