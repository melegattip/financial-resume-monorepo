import { useState, useMemo, useCallback } from 'react';
import { expensesAPI, incomesAPI, categoriesAPI } from '../services/api';
import dataService from '../services/dataService';
import toast from 'react-hot-toast';
import { useAuth } from '../contexts/AuthContext';

/**
 * Hook personalizado para operaciones CRUD optimizadas con invalidación de cache
 */
export const useOptimizedAPI = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { user } = useAuth();

  // Función genérica para manejar operaciones con invalidación de cache
  const executeWithCacheInvalidation = useCallback(async (operation, cacheType, successMessage) => {
    try {
      setLoading(true);
      setError(null);
      
      const result = await operation();
      console.log(`🔍 [useOptimizedAPI] Resultado de operación ${cacheType}:`, result);
      
      // Invalidar cache después de la operación exitosa
      console.log(`🗑️ [useOptimizedAPI] Invalidando caché para tipo: ${cacheType}`);
      dataService.invalidateAfterMutation(cacheType);
      
      if (successMessage) {
        toast.success(successMessage);
      }
      
      return result;
    } catch (err) {
      console.error(`❌ [useOptimizedAPI] Error en operación ${cacheType}:`, err);
      setError(err);
      toast.error(err.message || `Error en operación ${cacheType}`);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // === OPERACIONES DE GASTOS ===
  const expenses = useMemo(() => ({
    create: async (data) => {
      return executeWithCacheInvalidation(
        () => expensesAPI.create(data),
        'expense',
        'Gasto creado exitosamente'
      );
    },

    update: async (id, data) => {
      if (!user?.id) throw new Error('Usuario no autenticado');
      return executeWithCacheInvalidation(
        () => expensesAPI.update(user.id, id, data),
        'expense',
        'Gasto actualizado exitosamente'
      );
    },

    delete: async (id) => {
      if (!user?.id) throw new Error('Usuario no autenticado');
      return executeWithCacheInvalidation(
        () => expensesAPI.delete(user.id, id),
        'expense',
        'Gasto eliminado exitosamente'
      );
    },

    list: async () => {
      // Para operaciones de lectura, usar el cache del dataService
      return dataService.getCachedData(
        dataService.getCacheKey('expenses_list', {}),
        () => expensesAPI.list()
      );
    }
  }), [executeWithCacheInvalidation]);

  // === OPERACIONES DE INGRESOS ===
  const incomes = useMemo(() => ({
    create: async (data) => {
      return executeWithCacheInvalidation(
        () => incomesAPI.create(data),
        'income',
        'Ingreso creado exitosamente'
      );
    },

    update: async (id, data) => {
      if (!user?.id) throw new Error('Usuario no autenticado');
      return executeWithCacheInvalidation(
        () => incomesAPI.update(user.id, id, data),
        'income',
        'Ingreso actualizado exitosamente'
      );
    },

    delete: async (id) => {
      if (!user?.id) throw new Error('Usuario no autenticado');
      return executeWithCacheInvalidation(
        () => incomesAPI.delete(user.id, id),
        'income',
        'Ingreso eliminado exitosamente'
      );
    },

    list: async () => {
      return dataService.getCachedData(
        dataService.getCacheKey('incomes_list', {}),
        () => incomesAPI.list()
      );
    }
  }), [executeWithCacheInvalidation]);

  // === OPERACIONES DE CATEGORÍAS ===
  const categories = useMemo(() => ({
    create: async (data) => {
      return executeWithCacheInvalidation(
        () => categoriesAPI.create(data),
        'category',
        'Categoría creada exitosamente'
      );
    },

    update: async (id, data) => {
      return executeWithCacheInvalidation(
        () => categoriesAPI.update(id, data),
        'category',
        'Categoría actualizada exitosamente'
      );
    },

    delete: async (id) => {
      return executeWithCacheInvalidation(
        () => categoriesAPI.delete(id),
        'category',
        'Categoría eliminada exitosamente'
      );
    },

    list: async () => {
      return dataService.getCachedData(
        dataService.getCacheKey('categories_list', {}),
        () => categoriesAPI.list()
      );
    }
  }), [executeWithCacheInvalidation]);

  // === UTILIDADES ===
  const utils = useMemo(() => ({
    clearAllCache: () => {
      dataService.clearCache();
      toast.success('Cache limpiado exitosamente');
    },

    clearCachePattern: (pattern) => {
      dataService.clearCache(pattern);
      toast.success(`Cache limpiado para: ${pattern}`);
    },

    refreshDashboard: async (filterParams = {}) => {
      try {
        setLoading(true);
        // Limpiar cache del dashboard para forzar recarga
        dataService.clearCache('dashboard');
        dataService.clearCache('analytics');
        
        // Recargar datos
        const data = await dataService.loadDashboardData(filterParams, true);
        toast.success('Dashboard actualizado');
        return data;
      } catch (err) {
        console.error('Error refrescando dashboard:', err);
        toast.error('Error refrescando dashboard');
        throw err;
      } finally {
        setLoading(false);
      }
    }
  }), []);

  return useMemo(() => ({
    // Estados
    loading,
    error,
    
    // Operaciones CRUD
    expenses,
    incomes,
    categories,
    
    // Utilidades
    utils,
    
    // Acceso directo al dataService para casos especiales
    dataService
  }), [loading, error, expenses, incomes, categories, utils]);
};

export default useOptimizedAPI; 