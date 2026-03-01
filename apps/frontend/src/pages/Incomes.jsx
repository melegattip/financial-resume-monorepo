import React, { useState, useEffect, useCallback } from 'react';
import { createPortal } from 'react-dom';
import { FaPlus, FaSearch, FaArrowUp, FaEdit, FaTrash, FaDollarSign } from 'react-icons/fa';
import { formatCurrency } from '../services/api';
import dataService from '../services/dataService';
import { usePeriod } from '../contexts/PeriodContext';
import { useGamification } from '../contexts/GamificationContext';
import { useOptimizedAPI } from '../hooks/useOptimizedAPI';
import useDataRefresh from '../hooks/useDataRefresh';
import toast from 'react-hot-toast';
import ConfirmationModal from '../components/ConfirmationModal';
import ValidatedInput from '../components/ValidatedInput';
import { validateAmount, validateDescription } from '../utils/validation';

const Incomes = () => {
  const [incomes, setIncomes] = useState([]);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingIncome, setEditingIncome] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletingIncome, setDeletingIncome] = useState(null);
  const [deleteLoading, setDeleteLoading] = useState(false);
  
  // Estados para nuevos filtros de ordenamiento
  const [sortBy, setSortBy] = useState('category_priority');
  const [sortOrder, setSortOrder] = useState('desc');
  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    category_id: '',
    received_date: '',
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
    incomes: incomesAPI, 
    categories: categoriesAPI
  } = useOptimizedAPI();

  // Hook de gamificación para registrar acciones
  const { recordCreateIncome, recordUpdateIncome, recordDeleteIncome } = useGamification();

  const formatAmount = (amount) => {
    if (balancesHidden) return '••••••';
    return formatCurrency(amount);
  };

  // Función para obtener colores por categoría (consistente con Dashboard y Gastos)
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

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      console.log('🔄 Cargando datos de ingresos con API optimizada...');
      
      const [incomesResponse, categoriesResponse] = await Promise.all([
        incomesAPI.list(),
        categoriesAPI.list(),
      ]);
      
      // Normalizar datos de respuesta
      const incomesData = incomesResponse.data?.incomes || incomesResponse.incomes || incomesResponse || [];
      const categoriesData = categoriesResponse.data?.data || categoriesResponse.data || categoriesResponse || [];
      
      setIncomes(Array.isArray(incomesData) ? incomesData : []);
      setCategories(Array.isArray(categoriesData) ? categoriesData : []);
      
      // Actualizar datos disponibles en el contexto de períodos
      updateAvailableData([], incomesData);
      
      console.log('✅ Datos de ingresos cargados exitosamente:', {
        incomes: incomesData.length,
        categories: categoriesData.length
      });
      
    } catch (error) {
      console.warn('⚠️ Error al cargar ingresos:', error.message);
      
      // Establecer datos vacíos
      setIncomes([]);
      setCategories([]);
      
      // No mostrar toast aquí porque useOptimizedAPI ya lo maneja
    } finally {
      setLoading(false);
    }
  }, [incomesAPI, categoriesAPI, updateAvailableData]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Hook para refrescar automáticamente cuando cambian los datos
  useDataRefresh(loadData, ['income', 'recurring_transaction']);

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
    const amountValidation = validateAmount(formData.amount, { fieldName: 'monto' });
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
        received_date: formData.received_date ? `${formData.received_date}T00:00:00Z` : '',
      };

      if (editingIncome) {
        await incomesAPI.update(editingIncome.id, dataToSend);
        recordUpdateIncome(editingIncome.id, `Ingreso actualizado: ${dataToSend.description}`);
      } else {
        const result = await incomesAPI.create(dataToSend);
        const incomeId = result?.data?.id || `income-${Date.now()}`;
        recordCreateIncome(incomeId, `Nuevo ingreso: ${dataToSend.description}`);
      }

      dataService.invalidateAfterMutation('income');
      await loadData();
    } catch (error) {
      console.error('Error en handleSubmit:', error);
    } finally {
      setShowModal(false);
      setEditingIncome(null);
      setFormData({ description: '', amount: '', category_id: '', received_date: '' });
      setFormErrors({});
    }
  };

  const handleEdit = (income) => {
    setEditingIncome(income);
    setFormData({
      description: income.description,
      amount: income.amount.toString(),
      category_id: income.category_id || '',
      received_date: income.received_date ? income.received_date.slice(0, 10) : '',
    });
    setShowModal(true);
  };

  const handleDelete = (income) => {
    setDeletingIncome(income);
    setShowDeleteModal(true);
  };

  const confirmDelete = async () => {
    if (!deletingIncome) return;
    
    try {
      setDeleteLoading(true);
      await incomesAPI.delete(deletingIncome.id);
      // useOptimizedAPI ya muestra el toast de éxito
      
      // 🎮 Registrar acción de gamificación
      console.log(`🎯 [Incomes] Registrando eliminación de income: ${deletingIncome.id}`);
      recordDeleteIncome(deletingIncome.id, `Ingreso eliminado: ${deletingIncome.description}`);
      dataService.invalidateAfterMutation('income');
      await loadData();
    } catch (error) {
      // useOptimizedAPI ya maneja el error
      console.error('Error en confirmDelete:', error);
    } finally {
      // ✅ Siempre cerrar modal y limpiar estado, sin importar si hay errores
      setDeleteLoading(false);
      setShowDeleteModal(false);
      setDeletingIncome(null);
    }
  };

  const cancelDelete = () => {
    setShowDeleteModal(false);
    setDeletingIncome(null);
  };

  const filteredIncomes = Array.isArray(incomes) 
    ? incomes.filter(income => {
        const matchesSearch = income.description.toLowerCase().includes(searchTerm.toLowerCase());
        
        // Filtros de fecha — usar received_date (fecha funcional del ingreso)
        const rcvDate = income.received_date || income.created_at;
        const incomeDate = new Date(rcvDate);
        const matchesYear = !selectedYear || incomeDate.getFullYear().toString() === selectedYear;
        const matchesMonth = !selectedMonth || rcvDate.slice(0, 7) === selectedMonth;
        
        return matchesSearch && matchesYear && matchesMonth;
      })
      .sort((a, b) => {
        let aValue, bValue;
        
        switch (sortBy) {
          case 'category_priority': {
            const aCatP = categories.find(c => c.id === a.category_id);
            const bCatP = categories.find(c => c.id === b.category_id);
            const aPrio = aCatP?.priority || 0;
            const bPrio = bCatP?.priority || 0;
            // 0 (sin prioridad) siempre va al final, sin importar el orden
            if (aPrio === 0 && bPrio === 0) return 0;
            if (aPrio === 0) return 1;
            if (bPrio === 0) return -1;
            aValue = aPrio;
            bValue = bPrio;
            break;
          }
          case 'amount':
            aValue = a.amount;
            bValue = b.amount;
            break;
          case 'category': {
            const aCat = categories.find(c => c.id === a.category_id);
            const bCat = categories.find(c => c.id === b.category_id);
            aValue = (aCat?.name || 'Sin categoría').toLowerCase();
            bValue = (bCat?.name || 'Sin categoría').toLowerCase();
            break;
          }
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

  const totalIncomes = filteredIncomes.reduce((sum, income) => sum + income.amount, 0);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-fr-gray-600 dark:text-gray-400">Cargando ingresos...</span>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {/* Header con métricas compacto */}
      <div className="card py-3 px-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-6">
            <div>
              <p className="text-xs font-medium text-fr-gray-500 dark:text-gray-400">Total Ingresos</p>
              <p className="text-xl font-bold text-fr-secondary dark:text-green-400">{formatAmount(totalIncomes)}</p>
            </div>
            <div className="h-8 w-px bg-fr-gray-200 dark:bg-gray-600"></div>
            <div>
              <p className="text-xs font-medium text-fr-gray-500 dark:text-gray-400">Cantidad</p>
              <p className="text-xl font-bold text-fr-secondary dark:text-green-400">{filteredIncomes.length}</p>
            </div>
          </div>
          <div className="flex-shrink-0 p-2 rounded-fr bg-green-100 dark:bg-green-900/30">
            <FaArrowUp className="w-4 h-4 text-fr-secondary dark:text-green-400" />
          </div>
        </div>
      </div>

      {/* Controles + Lista unificados */}
      <div className="card p-0 overflow-hidden">
        {/* Toolbar compacta */}
        <div className="flex flex-wrap items-center gap-2 px-3 py-2 border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <div className="relative flex-1 min-w-[160px]">
            <FaSearch className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3 h-3 text-gray-400 dark:text-gray-500" />
            <input
              type="text"
              placeholder="Buscar ingresos..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-8 pr-3 py-1.5 text-sm border border-gray-200 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-500 w-full"
            />
          </div>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="text-sm border border-gray-200 dark:border-gray-600 rounded-lg px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="category_priority">Prioridad</option>
            <option value="created_at">Fecha creación</option>
            <option value="amount">Monto</option>
            <option value="category">Categoría</option>
          </select>
          <select
            value={sortOrder}
            onChange={(e) => setSortOrder(e.target.value)}
            className="text-sm border border-gray-200 dark:border-gray-600 rounded-lg px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="desc">↓ Desc</option>
            <option value="asc">↑ Asc</option>
          </select>
          <button
            onClick={() => setShowModal(true)}
            className="ml-auto btn-secondary flex items-center gap-1.5 py-1.5 px-3 text-sm"
          >
            <FaPlus className="w-3 h-3" />
            <span>Nuevo</span>
          </button>
        </div>

        {/* Lista de ingresos - estilo Gmail */}
        <div className="divide-y divide-gray-100 dark:divide-gray-700/50">
          {filteredIncomes.length === 0 ? (
            <div className="text-center py-12">
              <FaArrowUp className="w-12 h-12 text-fr-gray-400 dark:text-gray-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-fr-gray-900 dark:text-gray-100 mb-2">No hay ingresos</h3>
              <p className="text-fr-gray-500 dark:text-gray-400">Comienza agregando tu primer ingreso</p>
            </div>
          ) : (
            filteredIncomes.map((income) => {
              const category = categories.find(c => c.id === income.category_id);
              const color = getCategoryColor(income.category_id);
              
              return (
                <div key={income.id} className="flex items-center gap-2 py-1.5 px-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                  {/* Icono de ingreso */}
                  <div className="flex-shrink-0 w-6 h-6">
                    <div className="w-full h-full rounded-md bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
                      <FaArrowUp className="w-3 h-3 text-green-600 dark:text-green-400" />
                    </div>
                  </div>

                  {/* Descripción */}
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-fr-gray-900 dark:text-gray-100 text-sm truncate">
                      {income.description}
                    </h3>
                  </div>

                  {/* Categoría */}
                  <div className="flex-shrink-0 hidden sm:block text-left min-w-[80px]">
                    {category && (
                      <span className={`px-1.5 py-0.5 rounded-full text-xs font-medium ${color.bg} ${color.text} border ${color.border}`}>
                        {category.name}
                      </span>
                    )}
                  </div>

                  {/* Fecha del ingreso */}
                  <div className="flex-shrink-0 hidden md:block min-w-[100px] text-right">
                    <span className="text-xs text-gray-500 dark:text-gray-400">
                      {income.received_date ? new Date(income.received_date).toLocaleDateString('es-AR', { day: '2-digit', month: 'short', year: 'numeric' }) : '—'}
                    </span>
                  </div>

                  {/* Monto */}
                  <div className="flex-shrink-0 text-right min-w-[90px]">
                    <div className="font-semibold text-green-600 dark:text-green-400 text-sm">
                      +{formatAmount(income.amount)}
                    </div>
                  </div>
                  
                  {/* Botones de acción compactos */}
                  <div className="flex space-x-0.5 flex-shrink-0">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleEdit(income);
                      }}
                      className="w-6 h-6 bg-gray-100 dark:bg-gray-700 rounded-md flex items-center justify-center hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                      title="Editar"
                    >
                      <FaEdit className="w-3 h-3 text-gray-600 dark:text-gray-400" />
                    </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDelete(income);
                      }}
                      className="w-6 h-6 bg-red-100 dark:bg-red-900/30 rounded-md flex items-center justify-center hover:bg-red-200 dark:hover:bg-red-900/50 transition-colors"
                      title="Eliminar"
                    >
                      <FaTrash className="w-3 h-3 text-red-600 dark:text-red-400" />
                    </button>
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>

      {/* Modal */}
      {showModal && (
        createPortal(
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-fr-lg max-w-md w-full p-6">
              <h2 className="text-xl font-bold text-fr-gray-900 dark:text-gray-100 mb-6">
                {editingIncome ? 'Editar Ingreso' : 'Nuevo Ingreso'}
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
                  placeholder="Ej: Salario mensual, Freelance, etc."
                  helpText="Describe brevemente el ingreso"
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
                  helpText="Ingresa el monto del ingreso"
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
                    Fecha del ingreso
                  </label>
                  <input
                    type="date"
                    value={formData.received_date}
                    onChange={(e) => setFormData({ ...formData, received_date: e.target.value })}
                    className="input"
                  />
                </div>

                <div className="flex space-x-4 pt-4">
                  <button
                    type="button"
                    onClick={() => {
                      setShowModal(false);
                      setEditingIncome(null);
                      setFormData({ description: '', amount: '', category_id: '', received_date: '' });
                      setFormErrors({});
                    }}
                    className="btn-outline flex-1"
                  >
                    Cancelar
                  </button>
                  <button 
                    type="submit" 
                    className={`btn-secondary flex-1 ${!isFormValid ? 'opacity-50 cursor-not-allowed' : ''}`}
                    disabled={!isFormValid}
                  >
                    {editingIncome ? 'Actualizar' : 'Crear'}
                  </button>
                </div>
              </form>
            </div>
          </div>,
          document.body
        )
      )}

      {/* Modal de Confirmación de Eliminación */}
      <ConfirmationModal
        isOpen={showDeleteModal}
        onClose={cancelDelete}
        onConfirm={confirmDelete}
        title="Eliminar Ingreso"
        message={`¿Estás seguro de que quieres eliminar el ingreso "${deletingIncome?.description}"? Esta acción no se puede deshacer.`}
        confirmText="Eliminar"
        cancelText="Cancelar"
        type="danger"
        loading={deleteLoading}
      />
    </div>
  );
};

export default Incomes; 