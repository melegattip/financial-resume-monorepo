import React, { useState, useEffect, useCallback } from 'react';
import { FaPlus, FaSearch, FaTag, FaEdit, FaTrash, FaChartBar } from 'react-icons/fa';
import { useOptimizedAPI } from '../hooks/useOptimizedAPI';
import { useGamification } from '../contexts/GamificationContext';
import { usePeriod } from '../contexts/PeriodContext';
import ValidatedInput from '../components/ValidatedInput';
import { validateCategoryName } from '../utils/validation';
import { analyticsAPI, formatCurrency } from '../services/api';
import dataService from '../services/dataService';
import toast from 'react-hot-toast';

const CHART_COLORS = ['#009ee3', '#00a650', '#ff6900', '#e53e3e', '#6b7280', '#8b5cf6', '#f59e0b', '#06b6d4', '#ec4899'];

const Categories = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    name: '',
  });

  // Estados para validación del formulario
  const [formErrors, setFormErrors] = useState({});
  const [isFormValid, setIsFormValid] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Spending analytics
  const [spendingByCategory, setSpendingByCategory] = useState([]);
  const [spendingLoading, setSpendingLoading] = useState(false);

  // Usar el hook optimizado para operaciones API
  const {
    categories: categoriesAPI
  } = useOptimizedAPI();

  // Global period filter
  const { getFilterParams, getPeriodTitle } = usePeriod();

  // Hook de gamificación
  const { recordAction } = useGamification();

  const loadCategories = useCallback(async () => {
    try {
      setLoading(true);
      console.log('🔄 Cargando categorías con API optimizada...');
      
      const response = await categoriesAPI.list();
      console.log('🔍 [Categories] Respuesta completa del backend:', response);
      console.log('🔍 [Categories] response.data:', response.data);
      console.log('🔍 [Categories] response.data?.data:', response.data?.data);
      
      const categoriesData = response.data?.data || response.data || response || [];
      console.log('🔍 [Categories] Datos procesados:', categoriesData);
      console.log('🔍 [Categories] ¿Es array?:', Array.isArray(categoriesData));
      
      setCategories(Array.isArray(categoriesData) ? categoriesData : []);
      
      console.log('✅ Categorías cargadas exitosamente:', categoriesData.length);
    } catch (error) {
      console.error('❌ [Categories] Error loading categories:', error);
      // No mostrar toast aquí porque useOptimizedAPI ya lo maneja
      setCategories([]);
    } finally {
      setLoading(false);
    }
  }, [categoriesAPI]);

  const loadSpendingAnalytics = useCallback(async () => {
    setSpendingLoading(true);
    try {
      const analyticsParams = dataService.toAnalyticsDateParams(getFilterParams());
      const response = await analyticsAPI.categories(analyticsParams);
      const items = response?.data?.data || response?.data?.Categories || [];
      setSpendingByCategory(Array.isArray(items) ? items : []);
    } catch (err) {
      console.error('Error loading category spending:', err);
      setSpendingByCategory([]);
    } finally {
      setSpendingLoading(false);
    }
  }, [getFilterParams]);

  useEffect(() => {
    loadCategories();
    loadSpendingAnalytics();
  }, [loadCategories, loadSpendingAnalytics]);

  // Validar formulario completo
  const validateForm = useCallback(() => {
    const errors = {};
    let valid = true;

    // Validar nombre
    const nameValidation = validateCategoryName(formData.name);
    if (!nameValidation.isValid) {
      errors.name = nameValidation.error;
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
    
    // Prevenir doble click/submit
    if (isSubmitting) {
      console.log('🚫 [Categories] Submit ya en progreso, ignorando...');
      return;
    }
    
    // Validar antes de enviar
    if (!validateForm()) {
      toast.error('Por favor corrige los errores en el formulario');
      return;
    }

    setIsSubmitting(true);
    try {
      if (editingCategory) {
        // Para actualización, el backend espera el campo "new_name"
        const updateData = {
          new_name: formData.name
        };
        await categoriesAPI.update(editingCategory.id, updateData);
        // useOptimizedAPI ya muestra el toast de éxito
      } else {
        // Para creación, enviar solo el nombre (sin descripción por ahora)
        const createData = {
          name: formData.name
        };
        console.log('🚀 [Categories] Enviando datos para crear:', createData);
        const result = await categoriesAPI.create(createData);
        console.log('🔍 [Categories] Resultado de creación:', result);
        console.log('🔍 [Categories] result.data:', result.data);
        console.log('🔍 [Categories] result.status:', result.status);
        
        // 🎮 REGISTRAR GAMIFICACIÓN - Obtener ID de la respuesta
        if (result && result.data) {
          const categoryId = result.data.id || result.data.category_id;
          console.log('🔍 [Categories] ID extraído para gamificación:', categoryId);
          if (categoryId) {
            console.log('🎯 Registrando creación de categoría para gamificación:', categoryId);
            await recordAction('create_category', 'category', categoryId, `Created category: ${formData.name}`);
          } else {
            console.warn('⚠️ [Categories] No se pudo extraer el ID de la categoría creada');
          }
        } else {
          console.warn('⚠️ [Categories] Resultado de creación no tiene estructura esperada');
        }
        
        // useOptimizedAPI ya muestra el toast de éxito
      }
      
      setShowModal(false);
      setEditingCategory(null);
      setFormData({ name: '' });
      console.log('🔄 [Categories] Recargando categorías después de operación...');
      await loadCategories();
    } catch (error) {
      // useOptimizedAPI ya maneja el error
      console.error('Error en handleSubmit:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEdit = (category) => {
    setEditingCategory(category);
    setFormData({
      name: category.name,
    });
    setShowModal(true);
    setIsSubmitting(false);
  };

  const handleDelete = async (category) => {
    if (window.confirm('¿Estás seguro de que quieres eliminar esta categoría?')) {
      try {
        await categoriesAPI.delete(category.id);
        // useOptimizedAPI ya muestra el toast de éxito
        await loadCategories();
      } catch (error) {
        // useOptimizedAPI ya maneja el error
        console.error('Error en handleDelete:', error);
      }
    }
  };

  const filteredCategories = Array.isArray(categories) 
    ? categories.filter(category =>
        category.name.toLowerCase().includes(searchTerm.toLowerCase())
      )
    : [];

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="spinner"></div>
        <span className="ml-2 text-fr-gray-600 dark:text-gray-400">Cargando categorías...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="card">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="relative">
            <FaSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-fr-gray-400" />
            <input
              type="text"
              placeholder="Buscar categorías..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 input w-full sm:w-64"
            />
          </div>

          <button
            onClick={() => {
              setShowModal(true);
              setIsSubmitting(false);
            }}
            className="btn-primary flex items-center space-x-2"
          >
            <FaPlus className="w-4 h-4" />
            <span>Nueva Categoría</span>
          </button>
        </div>
      </div>

      {/* Gastos por categoría — período seleccionado */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-base font-semibold text-fr-gray-900 dark:text-gray-100 flex items-center space-x-2">
            <FaChartBar className="w-4 h-4 text-fr-primary dark:text-blue-400" />
            <span>Gastos por categoría — {getPeriodTitle()}</span>
          </h3>
        </div>

        {spendingLoading ? (
          <div className="space-y-2">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-8 bg-gray-200 dark:bg-gray-700 rounded animate-pulse" />
            ))}
          </div>
        ) : spendingByCategory.length === 0 ? (
          <p className="text-sm text-fr-gray-500 dark:text-gray-400 py-4 text-center">
            Sin gastos registrados con categoría en este período.
          </p>
        ) : (
          <div className="space-y-3">
            {spendingByCategory.map((cat, index) => {
              const pct = cat.percentage || 0;
              const color = CHART_COLORS[index % CHART_COLORS.length];
              return (
                <div key={cat.category_id || index}>
                  <div className="flex items-center justify-between text-sm mb-1">
                    <span className="font-medium text-fr-gray-800 dark:text-gray-200 truncate max-w-xs">
                      {cat.category_name || 'Sin nombre'}
                    </span>
                    <div className="flex items-center space-x-3 flex-shrink-0 ml-4">
                      <span className="text-fr-gray-500 dark:text-gray-400 text-xs">{pct.toFixed(1)}%</span>
                      <span className="font-semibold text-fr-gray-900 dark:text-gray-100">{formatCurrency(cat.amount || 0)}</span>
                    </div>
                  </div>
                  <div className="w-full bg-gray-100 dark:bg-gray-700 rounded-full h-2 overflow-hidden">
                    <div
                      className="h-2 rounded-full transition-all duration-500"
                      style={{ width: `${Math.min(pct, 100)}%`, backgroundColor: color }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Lista de categorías */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredCategories.length === 0 ? (
          <div className="col-span-full text-center py-12">
            <FaTag className="w-12 h-12 text-fr-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-fr-gray-900 dark:text-gray-100 mb-2">No hay categorías</h3>
            <p className="text-fr-gray-500 dark:text-gray-400">Comienza creando tu primera categoría</p>
          </div>
        ) : (
          filteredCategories.map((category) => (
            <div key={category.id} className="card-hover">
              <div className="flex items-start justify-between">
                <div className="flex items-center space-x-3">
                  <div className="p-2 rounded-fr bg-blue-100 dark:bg-blue-900/30">
                    <FaTag className="w-5 h-5 text-fr-primary dark:text-blue-400" />
                  </div>
                  <div>
                    <h3 className="font-medium text-fr-gray-900 dark:text-gray-100">{category.name}</h3>
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => handleEdit(category)}
                    className="p-2 rounded-fr text-fr-gray-600 dark:text-gray-400 hover:bg-fr-gray-200 dark:hover:bg-gray-700 transition-colors"
                  >
                    <FaEdit className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDelete(category)}
                    className="p-2 rounded-fr text-fr-error dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors"
                  >
                    <FaTrash className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-fr-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-fr-gray-900 dark:text-gray-100 mb-6">
              {editingCategory ? 'Editar Categoría' : 'Nueva Categoría'}
            </h2>

            <form onSubmit={handleSubmit} className="space-y-4">
              <ValidatedInput
                type="text"
                name="name"
                label="Nombre"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                validator={validateCategoryName}
                validateOnChange={true}
                required={true}
                placeholder="Ej: Alimentación, Transporte, Entretenimiento"
                helpText="Nombre único para identificar la categoría"
                maxLength={50}
              />



              <div className="flex space-x-4 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowModal(false);
                    setEditingCategory(null);
                    setFormData({ name: '' });
                    setIsSubmitting(false);
                  }}
                  className="btn-outline flex-1"
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  className={`btn-primary flex-1 ${(!isFormValid || isSubmitting) ? 'opacity-50 cursor-not-allowed' : ''}`}
                  disabled={!isFormValid || isSubmitting}
                >
                  {isSubmitting ? 'Enviando...' : (editingCategory ? 'Actualizar' : 'Crear')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Categories; 