// Utilidades SOLO para formateo de presentación (UI)
// NO contiene lógica de negocio

export const formatCurrency = (amount) => {
  const numericAmount = Number(amount);
  if (isNaN(numericAmount) || amount === null || amount === undefined) {
    return '$0,00';
  }
  
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
  }).format(numericAmount);
};

export const formatDate = (date) => {
  if (!date || date === null || date === undefined) {
    return 'Fecha no disponible';
  }
  
  const dateObj = new Date(date);
  if (isNaN(dateObj.getTime())) {
    return 'Fecha inválida';
  }
  
  return new Intl.DateTimeFormat('es-AR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(dateObj);
};

export const formatPercentage = (percentage) => {
  const numericPercentage = Number(percentage);
  if (isNaN(numericPercentage) || percentage === null || percentage === undefined) {
    return '0.0%';
  }
  
  return `${numericPercentage.toFixed(1)}%`;
};

export const formatShortDate = (date) => {
  if (!date) return 'N/A';
  
  const dateObj = new Date(date);
  if (isNaN(dateObj.getTime())) return 'Fecha inválida';
  
  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(dateObj);
};

export const formatAmount = (amount, hideBalance = false) => {
  if (hideBalance) return '••••••';
  return formatCurrency(amount);
};

export const truncateText = (text, maxLength = 50) => {
  if (!text || text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
}; 