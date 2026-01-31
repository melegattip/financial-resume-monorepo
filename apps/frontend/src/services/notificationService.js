/**
 * Servicio de notificaciones push avanzado
 */
class NotificationService {
  constructor() {
    this.isSupported = 'Notification' in window && 'serviceWorker' in navigator;
    this.permission = Notification.permission;
    this.registration = null;
    this.subscriptions = new Map();
  }

  /**
   * Inicializar el servicio de notificaciones
   */
  async init() {
    if (!this.isSupported) {
      console.warn('Notificaciones push no soportadas en este navegador');
      return false;
    }

    try {
      // Registrar service worker
      this.registration = await navigator.serviceWorker.register('/sw.js');
      console.log('Service Worker registrado:', this.registration);

      // Verificar estado del service worker
      if (this.registration.installing) {
        console.log('Service Worker instalando...');
      } else if (this.registration.waiting) {
        console.log('Service Worker esperando...');
      } else if (this.registration.active) {
        console.log('Service Worker activo');
      }

      return true;
    } catch (error) {
      console.error('Error registrando Service Worker:', error);
      return false;
    }
  }

  /**
   * Solicitar permisos de notificación
   */
  async requestPermission() {
    if (!this.isSupported) return false;

    if (this.permission === 'granted') {
      return true;
    }

    if (this.permission === 'denied') {
      console.warn('Permisos de notificación denegados por el usuario');
      return false;
    }

    try {
      const permission = await Notification.requestPermission();
      this.permission = permission;
      return permission === 'granted';
    } catch (error) {
      console.error('Error solicitando permisos:', error);
      return false;
    }
  }

  /**
   * Mostrar notificación local
   */
  async showNotification(title, options = {}) {
    if (!await this.requestPermission()) return;

    const defaultOptions = {
      icon: '/icon-192x192.png',
      badge: '/badge-72x72.png',
      vibrate: [200, 100, 200],
      requireInteraction: false,
      actions: []
    };

    const notificationOptions = { ...defaultOptions, ...options };

    if (this.registration) {
      // Usar service worker para notificaciones persistentes
      return this.registration.showNotification(title, notificationOptions);
    } else {
      // Fallback a notificaciones básicas
      return new Notification(title, notificationOptions);
    }
  }

  /**
   * Notificaciones específicas para la app financiera
   */
  async notifyExpenseAdded(expense) {
    await this.showNotification('Nuevo Gasto Registrado', {
      body: `${expense.description}: $${expense.amount}`,
      icon: '/icons/expense.png',
      tag: 'expense-added',
      data: { type: 'expense', id: expense.id }
    });
  }

  async notifyIncomeAdded(income) {
    await this.showNotification('Nuevo Ingreso Registrado', {
      body: `${income.description}: $${income.amount}`,
      icon: '/icons/income.png',
      tag: 'income-added',
      data: { type: 'income', id: income.id }
    });
  }

  async notifyExpenseDue(expense) {
    await this.showNotification('Gasto Pendiente', {
      body: `${expense.description} vence el ${expense.due_date}`,
      icon: '/icons/warning.png',
      tag: 'expense-due',
      requireInteraction: true,
      actions: [
        {
          action: 'pay',
          title: 'Marcar como Pagado'
        },
        {
          action: 'remind',
          title: 'Recordar Después'
        }
      ],
      data: { type: 'expense-due', id: expense.id }
    });
  }

  async notifyBudgetExceeded(category, amount, limit) {
    await this.showNotification('Presupuesto Excedido', {
      body: `Has gastado $${amount} en ${category} (límite: $${limit})`,
      icon: '/icons/budget-alert.png',
      tag: 'budget-exceeded',
      requireInteraction: true,
      data: { type: 'budget-alert', category, amount, limit }
    });
  }

  /**
   * Programar notificaciones recurrentes
   */
  scheduleRecurringNotifications() {
    // Recordatorio semanal para revisar gastos
    this.scheduleNotification('weekly-review', {
      title: 'Revisión Semanal',
      body: '¿Ya revisaste tus gastos de esta semana?',
      triggerAt: this.getNextWeeklyReview()
    });

    // Recordatorio mensual para análisis financiero
    this.scheduleNotification('monthly-analysis', {
      title: 'Análisis Mensual',
      body: 'Es momento de revisar tu resumen financiero mensual',
      triggerAt: this.getNextMonthlyAnalysis()
    });
  }

  /**
   * Programar notificación específica
   */
  scheduleNotification(id, { title, body, triggerAt, ...options }) {
    const timeUntilTrigger = triggerAt.getTime() - Date.now();
    
    if (timeUntilTrigger <= 0) return;

    const timeoutId = setTimeout(() => {
      this.showNotification(title, { body, ...options });
    }, timeUntilTrigger);

    this.subscriptions.set(id, timeoutId);
  }

  /**
   * Cancelar notificación programada
   */
  cancelScheduledNotification(id) {
    const timeoutId = this.subscriptions.get(id);
    if (timeoutId) {
      clearTimeout(timeoutId);
      this.subscriptions.delete(id);
    }
  }

  /**
   * Utilidades para fechas
   */
  getNextWeeklyReview() {
    const now = new Date();
    const nextSunday = new Date(now);
    nextSunday.setDate(now.getDate() + (7 - now.getDay()));
    nextSunday.setHours(19, 0, 0, 0); // 7 PM los domingos
    return nextSunday;
  }

  getNextMonthlyAnalysis() {
    const now = new Date();
    const nextMonth = new Date(now.getFullYear(), now.getMonth() + 1, 1);
    nextMonth.setHours(10, 0, 0, 0); // 10 AM el primer día del mes
    return nextMonth;
  }

  /**
   * Limpiar recursos
   */
  cleanup() {
    this.subscriptions.forEach((timeoutId) => {
      clearTimeout(timeoutId);
    });
    this.subscriptions.clear();
  }
}

// Instancia singleton
const notificationService = new NotificationService();

export default notificationService; 