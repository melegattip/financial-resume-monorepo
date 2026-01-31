/**
 * Utilidades de Validación para Financial Resume Engine
 * Validaciones client-side para mejorar UX y reducir requests innecesarios
 */

// Regex patterns comunes
const PATTERNS = {
  // Permitir letras, números, espacios y caracteres especiales comunes en español
  DESCRIPTION: /^[a-zA-Z0-9\s\-.,áéíóúÁÉÍÓÚñÑüÜ()&@#$%/]+$/,
  // Email básico
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  // Solo números y punto decimal
  NUMERIC: /^[0-9]+(\.[0-9]{1,2})?$/,
  // Nombre de categoría (más restrictivo)
  CATEGORY_NAME: /^[a-zA-Z0-9\s\-áéíóúÁÉÍÓÚñÑüÜ]+$/,
};

/**
 * Valida montos monetarios
 * @param {string|number} value - Valor a validar
 * @param {Object} options - Opciones de validación
 * @returns {Object} { isValid: boolean, error: string|null }
 */
export const validateAmount = (value, options = {}) => {
  const {
    required = true,
    min = 0.01,
    max = 999999999,
    allowZero = false,
    fieldName = 'monto'
  } = options;

  // Verificar si es requerido
  if (required && (!value || value === '')) {
    return { isValid: false, error: `El ${fieldName} es requerido` };
  }

  // Si no es requerido y está vacío, es válido
  if (!required && (!value || value === '')) {
    return { isValid: true, error: null };
  }

  // Convertir a número
  const numericValue = typeof value === 'string' ? parseFloat(value.replace(/,/g, '')) : value;

  // Verificar si es un número válido
  if (isNaN(numericValue)) {
    return { isValid: false, error: `El ${fieldName} debe ser un número válido` };
  }

  // Verificar si es infinito
  if (!isFinite(numericValue)) {
    return { isValid: false, error: `El ${fieldName} no puede ser infinito` };
  }

  // Verificar valor mínimo
  const minValue = allowZero ? 0 : min;
  if (numericValue < minValue) {
    return { 
      isValid: false, 
      error: allowZero 
        ? `El ${fieldName} no puede ser negativo` 
        : `El ${fieldName} debe ser mayor a $${min}`
    };
  }

  // Verificar valor máximo
  if (numericValue > max) {
    return { 
      isValid: false, 
      error: `El ${fieldName} no puede exceder $${max.toLocaleString()}`
    };
  }

  // Verificar decimales (máximo 2)
  const decimalPlaces = (numericValue.toString().split('.')[1] || '').length;
  if (decimalPlaces > 2) {
    return { 
      isValid: false, 
      error: `El ${fieldName} no puede tener más de 2 decimales`
    };
  }

  return { isValid: true, error: null };
};

/**
 * Valida descripciones de transacciones
 * @param {string} value - Descripción a validar
 * @param {Object} options - Opciones de validación
 * @returns {Object} { isValid: boolean, error: string|null }
 */
export const validateDescription = (value, options = {}) => {
  const {
    required = true,
    minLength = 3,
    maxLength = 255,
    fieldName = 'descripción'
  } = options;

  // Verificar si value es undefined o null
  if (!value || typeof value !== 'string') {
    if (required) {
      return { isValid: false, error: `La ${fieldName} es requerida` };
    }
    return { isValid: true, error: null };
  }

  // Verificar si es requerido
  if (required && value.trim() === '') {
    return { isValid: false, error: `La ${fieldName} es requerida` };
  }

  // Si no es requerido y está vacío, es válido
  if (!required && value.trim() === '') {
    return { isValid: true, error: null };
  }

  const trimmedValue = value.trim();

  // Verificar longitud mínima
  if (trimmedValue.length < minLength) {
    return { 
      isValid: false, 
      error: `La ${fieldName} debe tener al menos ${minLength} caracteres`
    };
  }

  // Verificar longitud máxima
  if (trimmedValue.length > maxLength) {
    return { 
      isValid: false, 
      error: `La ${fieldName} no puede exceder ${maxLength} caracteres`
    };
  }

  // Verificar caracteres válidos
  if (!PATTERNS.DESCRIPTION.test(trimmedValue)) {
    return { 
      isValid: false, 
      error: `La ${fieldName} contiene caracteres no válidos`
    };
  }

  return { isValid: true, error: null };
};

/**
 * Valida nombres de categorías
 * @param {string} value - Nombre a validar
 * @param {Object} options - Opciones de validación
 * @returns {Object} { isValid: boolean, error: string|null }
 */
export const validateCategoryName = (value, options = {}) => {
  const {
    required = true,
    minLength = 2,
    maxLength = 50
  } = options;

  // Verificar si value es undefined o null
  if (!value || typeof value !== 'string') {
    if (required) {
      return { isValid: false, error: 'El nombre de la categoría es requerido' };
    }
    return { isValid: true, error: null };
  }

  // Verificar si es requerido
  if (required && value.trim() === '') {
    return { isValid: false, error: 'El nombre de la categoría es requerido' };
  }

  // Si no es requerido y está vacío, es válido
  if (!required && value.trim() === '') {
    return { isValid: true, error: null };
  }

  const trimmedValue = value.trim();

  // Verificar longitud mínima
  if (trimmedValue.length < minLength) {
    return { 
      isValid: false, 
      error: `El nombre debe tener al menos ${minLength} caracteres`
    };
  }

  // Verificar longitud máxima
  if (trimmedValue.length > maxLength) {
    return { 
      isValid: false, 
      error: `El nombre no puede exceder ${maxLength} caracteres`
    };
  }

  // Verificar caracteres válidos (más restrictivo para nombres)
  if (!PATTERNS.CATEGORY_NAME.test(trimmedValue)) {
    return { 
      isValid: false, 
      error: 'El nombre contiene caracteres no válidos'
    };
  }

  return { isValid: true, error: null };
};

/**
 * Valida fechas
 * @param {string} value - Fecha en formato YYYY-MM-DD
 * @param {Object} options - Opciones de validación
 * @returns {Object} { isValid: boolean, error: string|null }
 */
export const validateDate = (value, options = {}) => {
  const {
    required = true,
    allowPast = true,
    allowFuture = true,
    fieldName = 'fecha'
  } = options;

  // Verificar si value es undefined o null
  if (!value || typeof value !== 'string') {
    if (required) {
      return { isValid: false, error: `La ${fieldName} es requerida` };
    }
    return { isValid: true, error: null };
  }

  // Verificar si es requerido
  if (required && value.trim() === '') {
    return { isValid: false, error: `La ${fieldName} es requerida` };
  }

  // Si no es requerido y está vacío, es válido
  if (!required && value.trim() === '') {
    return { isValid: true, error: null };
  }

  // Verificar formato de fecha
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(value)) {
    return { 
      isValid: false, 
      error: `La ${fieldName} debe tener el formato YYYY-MM-DD`
    };
  }

  // Verificar si es una fecha válida
  const date = new Date(value);
  if (isNaN(date.getTime())) {
    return { isValid: false, error: `La ${fieldName} no es válida` };
  }

  // Verificar si la fecha coincide con el string (evita fechas como 2024-02-30)
  if (date.toISOString().split('T')[0] !== value) {
    return { isValid: false, error: `La ${fieldName} no es válida` };
  }

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const inputDate = new Date(value);
  inputDate.setHours(0, 0, 0, 0);

  // Verificar fechas pasadas
  if (!allowPast && inputDate < today) {
    return { 
      isValid: false, 
      error: `La ${fieldName} no puede ser anterior a hoy`
    };
  }

  // Verificar fechas futuras
  if (!allowFuture && inputDate > today) {
    return { 
      isValid: false, 
      error: `La ${fieldName} no puede ser posterior a hoy`
    };
  }

  return { isValid: true, error: null };
};

/**
 * Valida emails
 * @param {string} value - Email a validar
 * @param {Object} options - Opciones de validación
 * @returns {Object} { isValid: boolean, error: string|null }
 */
export const validateEmail = (value, options = {}) => {
  const { required = true } = options;

  // Verificar si value es undefined o null
  if (!value || typeof value !== 'string') {
    if (required) {
      return { isValid: false, error: 'El email es requerido' };
    }
    return { isValid: true, error: null };
  }

  // Verificar si es requerido
  if (required && value.trim() === '') {
    return { isValid: false, error: 'El email es requerido' };
  }

  // Si no es requerido y está vacío, es válido
  if (!required && value.trim() === '') {
    return { isValid: true, error: null };
  }

  const trimmedValue = value.trim();

  // Verificar formato de email
  if (!PATTERNS.EMAIL.test(trimmedValue)) {
    return { isValid: false, error: 'El email no tiene un formato válido' };
  }

  // Verificar longitud máxima
  if (trimmedValue.length > 254) {
    return { isValid: false, error: 'El email es demasiado largo' };
  }

  return { isValid: true, error: null };
};

/**
 * Valida porcentajes
 * @param {string|number} value - Porcentaje a validar
 * @param {Object} options - Opciones de validación
 * @returns {Object} { isValid: boolean, error: string|null }
 */
export const validatePercentage = (value, options = {}) => {
  const {
    required = true,
    min = 0,
    max = 100,
    fieldName = 'porcentaje'
  } = options;

  // Verificar si es requerido
  if (required && (!value || value === '')) {
    return { isValid: false, error: `El ${fieldName} es requerido` };
  }

  // Si no es requerido y está vacío, es válido
  if (!required && (!value || value === '')) {
    return { isValid: true, error: null };
  }

  // Convertir a número
  const numericValue = typeof value === 'string' ? parseFloat(value) : value;

  // Verificar si es un número válido
  if (isNaN(numericValue)) {
    return { isValid: false, error: `El ${fieldName} debe ser un número válido` };
  }

  // Verificar rango
  if (numericValue < min || numericValue > max) {
    return { 
      isValid: false, 
      error: `El ${fieldName} debe estar entre ${min}% y ${max}%`
    };
  }

  return { isValid: true, error: null };
};

/**
 * Valida formularios completos
 * @param {Object} formData - Datos del formulario
 * @param {Object} validationRules - Reglas de validación
 * @returns {Object} { isValid: boolean, errors: Object }
 */
export const validateForm = (formData, validationRules) => {
  const errors = {};
  let isValid = true;

  Object.keys(validationRules).forEach(fieldName => {
    const rule = validationRules[fieldName];
    const value = formData[fieldName];
    
    let validation;
    
    switch (rule.type) {
      case 'amount':
        validation = validateAmount(value, rule.options);
        break;
      case 'description':
        validation = validateDescription(value, rule.options);
        break;
      case 'categoryName':
        validation = validateCategoryName(value, rule.options);
        break;
      case 'date':
        validation = validateDate(value, rule.options);
        break;
      case 'email':
        validation = validateEmail(value, rule.options);
        break;
      case 'percentage':
        validation = validatePercentage(value, rule.options);
        break;
      default:
        validation = { isValid: true, error: null };
    }

    if (!validation.isValid) {
      errors[fieldName] = validation.error;
      isValid = false;
    }
  });

  return { isValid, errors };
};

/**
 * Formatea números para input
 * @param {string} value - Valor del input
 * @returns {string} Valor formateado
 */
export const formatNumericInput = (value) => {
  // Remover todo excepto números, punto y coma
  let formatted = value.replace(/[^0-9.,]/g, '');
  
  // Reemplazar comas por puntos para decimales
  formatted = formatted.replace(/,/g, '.');
  
  // Permitir solo un punto decimal
  const parts = formatted.split('.');
  if (parts.length > 2) {
    formatted = parts[0] + '.' + parts.slice(1).join('');
  }
  
  // Limitar decimales a 2 dígitos
  if (parts.length === 2 && parts[1].length > 2) {
    formatted = parts[0] + '.' + parts[1].substring(0, 2);
  }
  
  return formatted;
};

/**
 * Sanitiza texto para prevenir XSS básico
 * @param {string} value - Texto a sanitizar
 * @returns {string} Texto sanitizado
 */
export const sanitizeText = (value) => {
  if (typeof value !== 'string') return value;
  
  return value
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#x27;')
    .replace(/\//g, '&#x2F;');
};

// Exportar patterns para uso externo si es necesario
export { PATTERNS };

// Exportar funciones de validación específicas para casos comunes
export const VALIDATION_RULES = {
  // Reglas para gastos/ingresos
  TRANSACTION: {
    amount: {
      type: 'amount',
      options: { required: true, min: 0.01, max: 999999999, fieldName: 'monto' }
    },
    description: {
      type: 'description',
      options: { required: true, minLength: 3, maxLength: 255, fieldName: 'descripción' }
    }
  },
  
  // Reglas para categorías
  CATEGORY: {
    name: {
      type: 'categoryName',
      options: { required: true, minLength: 2, maxLength: 50 }
    },
    description: {
      type: 'description',
      options: { required: false, minLength: 0, maxLength: 500, fieldName: 'descripción' }
    }
  },
  
  // Reglas para presupuestos
  BUDGET: {
    amount: {
      type: 'amount',
      options: { required: true, min: 1, max: 999999999, fieldName: 'presupuesto' }
    },
    alertAt: {
      type: 'percentage',
      options: { required: false, min: 1, max: 100, fieldName: 'alerta' }
    }
  },
  
  // Reglas para metas de ahorro
  SAVINGS_GOAL: {
    name: {
      type: 'description',
      options: { required: true, minLength: 3, maxLength: 100, fieldName: 'nombre de la meta' }
    },
    targetAmount: {
      type: 'amount',
      options: { required: true, min: 1, max: 999999999, fieldName: 'monto objetivo' }
    },
    targetDate: {
      type: 'date',
      options: { required: true, allowPast: false, fieldName: 'fecha objetivo' }
    }
  }
}; 