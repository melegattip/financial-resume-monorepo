/**
 * Sistema simple de notificaciones sin dependencias externas
 * Para uso temporal hasta que se pueda instalar react-hot-toast
 */

// Container para notificaciones
let notificationContainer = null;

// Inicializar el container de notificaciones
const initNotificationContainer = () => {
  if (notificationContainer) return notificationContainer;

  notificationContainer = document.createElement('div');
  notificationContainer.id = 'notification-container';
  notificationContainer.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 9999;
    pointer-events: none;
    display: flex;
    flex-direction: column;
    gap: 10px;
  `;
  
  document.body.appendChild(notificationContainer);
  return notificationContainer;
};

// Crear una notificación
const createNotification = (message, type = 'info', duration = 4000) => {
  const container = initNotificationContainer();
  
  const notification = document.createElement('div');
  notification.style.cssText = `
    padding: 12px 16px;
    border-radius: 8px;
    color: white;
    font-family: system-ui, -apple-system, sans-serif;
    font-size: 14px;
    font-weight: 500;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    pointer-events: auto;
    cursor: pointer;
    transform: translateX(100%);
    transition: transform 0.3s ease-in-out, opacity 0.3s ease-in-out;
    max-width: 350px;
    word-wrap: break-word;
  `;

  // Colores según el tipo
  const colors = {
    success: '#10b981', // green-500
    error: '#ef4444',   // red-500
    warning: '#f59e0b', // amber-500
    info: '#3b82f6'     // blue-500
  };

  notification.style.backgroundColor = colors[type] || colors.info;
  notification.textContent = message;

  // Agregar al container
  container.appendChild(notification);

  // Animación de entrada
  setTimeout(() => {
    notification.style.transform = 'translateX(0)';
  }, 10);

  // Auto-remover después del tiempo especificado
  const timeoutId = setTimeout(() => {
    removeNotification(notification);
  }, duration);

  // Remover al hacer click
  notification.addEventListener('click', () => {
    clearTimeout(timeoutId);
    removeNotification(notification);
  });

  return notification;
};

// Remover notificación con animación
const removeNotification = (notification) => {
  notification.style.transform = 'translateX(100%)';
  notification.style.opacity = '0';
  
  setTimeout(() => {
    if (notification.parentNode) {
      notification.parentNode.removeChild(notification);
    }
  }, 300);
};

// API pública similar a react-hot-toast
const toast = {
  success: (message, options = {}) => {
    const duration = options.duration || 4000;
    return createNotification(message, 'success', duration);
  },
  
  error: (message, options = {}) => {
    const duration = options.duration || 6000; // Errores duran más
    return createNotification(message, 'error', duration);
  },
  
  warning: (message, options = {}) => {
    const duration = options.duration || 5000;
    return createNotification(message, 'warning', duration);
  },
  
  info: (message, options = {}) => {
    const duration = options.duration || 4000;
    return createNotification(message, 'info', duration);
  },

  // Método genérico
  default: (message, options = {}) => {
    const type = options.type || 'info';
    const duration = options.duration || 4000;
    return createNotification(message, type, duration);
  }
};

export default toast; 