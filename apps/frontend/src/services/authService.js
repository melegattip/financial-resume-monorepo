import axios from 'axios';
import toast from '../utils/notifications';
import configService from './configService';
import dataService from './dataService';

// Función para obtener la URL del users service dinámicamente
const getUsersServiceUrl = async () => {
  try {
    const config = await configService.loadConfig();
    console.log('🔧 [authService] getUsersServiceUrl - Config recibida:', config);
    console.log('🔧 [authService] getUsersServiceUrl - users_service_url:', config.users_service_url);
    return config.users_service_url;
  } catch (error) {
    console.error('Error getting users service URL:', error);
    
    // Fallback con detección de ambiente
    if (process.env.REACT_APP_USERS_SERVICE_URL) {
      return process.env.REACT_APP_USERS_SERVICE_URL;
    }
    
    // Usar el proxy del backend principal en lugar de conectar directamente
    const hostname = window.location.hostname;
    if (hostname.includes('onrender.com') || hostname === 'financial.niloft.com') {
      return 'https://financial-resume-monorepo.onrender.com/api/v1';  // Render monolith
    } else {
      return 'http://localhost:8080/api/v1';  // Development
    }
  }
};

// Función para determinar baseURL inicial del users service por ambiente
const getInitialAuthBaseURL = () => {
  if (process.env.REACT_APP_USERS_SERVICE_URL) {
    return process.env.REACT_APP_USERS_SERVICE_URL;
  }

  const hostname = window.location.hostname;
  if (hostname.includes('onrender.com') || hostname === 'financial.niloft.com') {
    return 'https://financial-resume-monorepo.onrender.com/api/v1';  // Render monolith
  } else {
    return 'http://localhost:8080/api/v1';  // Development
  }
};

// Crear instancia de axios para autenticación
const authAPI = axios.create({
  baseURL: getInitialAuthBaseURL(),
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Función para inicializar la configuración dinámica
let configInitialized = false;

const initializeConfig = async () => {
  if (configInitialized) return;
  
  try {
    // Initializing dynamic configuration
    const config = await configService.loadConfig();
    
    // Actualizar la baseURL de axios con la configuración del users service
    console.log('🔧 [authService] Configuración recibida del backend:', config);
    authAPI.defaults.baseURL = config.users_service_url;
    configInitialized = true;
    
    console.log('✅ [authService] Configuración dinámica inicializada:', {
      baseURL: authAPI.defaults.baseURL,
      environment: config.environment,
      version: config.version
    });
  } catch (error) {
    console.error('❌ [authService] Error inicializando configuración:', error);
    // Mantener la configuración por defecto
  }
};

// Inicializar configuración al cargar el módulo
initializeConfig();

const TOKEN_KEY = 'auth_token';
const REFRESH_TOKEN_KEY = 'auth_refresh_token';
const USER_KEY = 'auth_user';
const EXPIRES_AT_KEY = 'auth_expires_at';

/**
 * Servicio de autenticación
 */
class AuthService {
  constructor() {
    this.token = localStorage.getItem(TOKEN_KEY);
    this.refreshToken_value = localStorage.getItem(REFRESH_TOKEN_KEY);
    this.expiresAt = localStorage.getItem(EXPIRES_AT_KEY);
    
    // Cargar usuario si existe token válido
    this.loadUserFromStorage();
    
    // Configurar interceptor de axios para agregar token automáticamente
    this.setupAuthInterceptor();
    
    // Configurar renovación automática de tokens
    this.setupTokenRenewal();
  }

  /**
   * Configura la renovación automática de tokens
   */
  setupTokenRenewal() {
    // Verificar si el token va a expirar en los próximos 10 minutos
    setInterval(() => {
      if (this.isAuthenticated() && this.expiresAt) {
        const now = Math.floor(Date.now() / 1000);
        const expirationTime = parseInt(this.expiresAt);
        const timeUntilExpiry = expirationTime - now;
        
        // Si el token expira en los próximos 10 minutos, renovarlo automáticamente
        if (timeUntilExpiry > 0 && timeUntilExpiry < 600) { // 10 minutos
          console.log('🔄 [authService] Token expira pronto, renovando automáticamente...');
          this.refreshToken().catch(error => {
            console.error('❌ [authService] Error en renovación automática:', error);
          });
        }
      }
    }, 60000); // Verificar cada minuto
  }

  /**
   * Carga el usuario del localStorage si hay un token válido
   */
  loadUserFromStorage() {
    if (this.token && this.isTokenValid()) {
      const storedUser = localStorage.getItem(USER_KEY);
      if (storedUser) {
        try {
          this.user = JSON.parse(storedUser);
        } catch (error) {
          console.error('Error parseando usuario del localStorage:', error);
          this.clearAuthData();
        }
      }
    } else {
      // Si no hay token válido, limpiar datos
      if (this.token) {
        console.log('🧹 [authService] Token inválido, limpiando datos');
        this.clearAuthData();
      }
    }
  }

  /**
   * Configura el interceptor de Axios para agregar token automáticamente
   */
  setupAuthInterceptor() {
    authAPI.interceptors.request.use(
      async (config) => {
        // En desarrollo, evitar múltiples inicializaciones de configuración
        const isDevelopment = window.location.hostname === 'localhost';
        
        if (this.token && this.isTokenValid()) {
          config.headers.Authorization = `Bearer ${this.token}`;
          
          // Agregar X-Caller-ID si tenemos usuario
          if (this.user?.id) {
            config.headers['X-Caller-ID'] = this.user.id.toString();
          }
        }
        
        return config;
      },
      (error) => Promise.reject(error)
    );

    authAPI.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        // Solo intentar refresh si es un 401 y no es una petición de login/register/refresh
        const isAuthEndpoint = originalRequest.url?.includes('/auth/login') ||
                              originalRequest.url?.includes('/auth/register') ||
                              originalRequest.url?.includes('/auth/refresh');

        if (error.response?.status === 401 && !originalRequest._retry && !isAuthEndpoint) {
          originalRequest._retry = true;

          // Intentar refresh token si tenemos un refresh token válido
          if (this.refreshToken_value) {
            try {
              console.log('🔄 [authService] Intentando renovar token automáticamente...');
              await this.refreshToken();
              originalRequest.headers.Authorization = `Bearer ${this.token}`;
              return authAPI(originalRequest);
            } catch (refreshError) {
              console.error('❌ [authService] Error renovando token:', refreshError);
              // Solo hacer logout si realmente tenemos una sesión activa
              if (this.isAuthenticated()) {
                this.logout();
                window.location.href = '/login';
              }
              return Promise.reject(refreshError);
            }
          } else {
            // Si no tenemos refresh token, limpiar datos y rechazar
            this.clearAuthData();
            return Promise.reject(error);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  /**
   * Verifica si el token actual es válido (no expirado)
   */
  isTokenValid() {
    if (!this.token || !this.expiresAt) return false;
    
    const now = Math.floor(Date.now() / 1000); // Tiempo actual en segundos
    const expirationTime = parseInt(this.expiresAt);
    
    // Considerar token inválido si expira en los próximos 2 minutos
    return expirationTime > (now + 120);
  }

  /**
   * Registra un nuevo usuario
   * @param {Object} userData - Datos del usuario (email, password, firstName, lastName)
   */
  async register(userData) {
    try {
      // Transformar datos del frontend al formato que espera el backend
      const backendData = {
        email: userData.email,
        password: userData.password,
        first_name: userData.firstName,
        last_name: userData.lastName
      };
      
      console.log('🔧 Enviando datos de registro:', backendData);
      
      const response = await authAPI.post('/auth/register', backendData);
      const authData = response.data;
      
      // Verificar que la respuesta tenga la estructura esperada
      // El backend devuelve { user, tokens: { access_token, refresh_token, expires_at } }
      if ((!authData.access_token && !authData.tokens?.access_token) || !authData.user) {
        console.error('❌ [authService] Respuesta inválida del servidor:', authData);
        throw new Error('Respuesta inválida del servidor');
      }

      this.saveAuthData(authData);

      console.log('✅ [authService] Registro exitoso para usuario:', authData.user.first_name);
      toast.success('¡Registro exitoso! Bienvenido');
      return { success: true, data: authData };
    } catch (error) {
      console.error('❌ Error en registro:', error.response?.data || error.message);
      const message = error.response?.data?.error || error.response?.data?.message || 'Error en el registro';
      toast.error(message);
      throw new Error(message);
    }
  }

  /**
   * Inicia sesión con email y contraseña
   * @param {Object} credentials - Credenciales (email, password)
   */
  async login(credentials) {
    try {
      console.log('🔧 [authService] Intentando login con credenciales:', { email: credentials.email });
      const response = await authAPI.post('/auth/login', credentials);
      console.log('🔧 [authService] Respuesta del servidor:', response.data);
      
      const authData = response.data;
      
      // Verificar si el servidor requiere 2FA
      if (authData.error === '2FA code required' || authData.requires_2fa) {
        console.log('🔧 [authService] 2FA requerido para el usuario');
        throw new Error('2FA code required');
      }
      
      // Verificar que la respuesta tenga la estructura esperada
      // El backend devuelve { user, tokens: { access_token, refresh_token, expires_at } }
      if ((!authData.access_token && !authData.tokens?.access_token) || !authData.user) {
        console.error('❌ [authService] Respuesta inválida del servidor:', authData);
        throw new Error('Respuesta inválida del servidor');
      }

      this.saveAuthData(authData);

      console.log('✅ [authService] Login exitoso para usuario:', authData.user.first_name);
      toast.success(`¡Bienvenido de vuelta, ${authData.user.first_name}!`);
      return { success: true, data: authData };
    } catch (error) {
      console.error('❌ [authService] Error en login:', error);
      
      // Si es un error de 2FA, re-lanzar el error específico
      if (error.message === '2FA code required' || error.response?.data?.requires_2fa) {
        throw new Error('2FA code required');
      }
      
      const message = error.response?.data?.error || 'Error en el login';
      toast.error(message);
      throw new Error(message);
    }
  }

  /**
   * Cierra la sesión actual
   */
  async logout() {
    try {
      // Intentar notificar al servidor (opcional)
      if (this.isAuthenticated()) {
        await authAPI.post('/users/logout');
      }
    } catch (error) {
      // No importa si falla, igual limpiaremos el local storage
      console.warn('Error during logout:', error);
    } finally {
      this.clearAuthData();
      toast.success('Sesión cerrada correctamente');
    }
  }

  /**
   * Renueva el token JWT
   */
  async refreshToken() {
    try {
      console.log('🔄 [authService] Renovando token...');
      
      if (!this.refreshToken_value) {
        throw new Error('No hay refresh token disponible');
      }

      const response = await authAPI.post('/auth/refresh', {
        refresh_token: this.refreshToken_value
      });
      
      const authData = response.data;
      console.log('✅ [authService] Token renovado exitosamente');
      
      this.saveAuthData(authData);
      return authData;
    } catch (error) {
      console.error('❌ [authService] Error renovando token:', error);
      this.clearAuthData();
      throw new Error('Error renovando token');
    }
  }

  /**
   * Obtiene el perfil del usuario autenticado
   */
  async getProfile() {
    try {
      const response = await authAPI.get('/users/profile');
      console.log('🔧 [authService] Respuesta completa de /users/profile:', response.data);
      const userData = response.data.user || response.data; // El backend envía {user: {...}}
      console.log('🔧 [authService] Datos del usuario extraídos:', userData);
      console.log('🔧 [authService] Avatar en userData:', userData?.avatar);
      
      // Actualizar datos del usuario en memoria y storage
      this.user = userData;
      localStorage.setItem(USER_KEY, JSON.stringify(userData));
      
      return userData;
    } catch (error) {
      const message = error.response?.data?.error || 'Error obteniendo perfil';
      throw new Error(message);
    }
  }

  /**
   * Cambia la contraseña del usuario
   * @param {Object} passwordData - Datos de cambio de contraseña
   */
  async changePassword(passwordData) {
    try {
      console.log('🔧 [authService] Cambiando contraseña:', passwordData);
      
      // Usar el endpoint correcto del users-service
      const response = await authAPI.put('/users/change-password', passwordData);
      console.log('✅ [authService] Contraseña cambiada exitosamente');
      
      toast.success('Contraseña cambiada exitosamente');
      return { success: true };
    } catch (error) {
      console.error('❌ [authService] Error cambiando contraseña:', error);
      const message = error.response?.data?.error || 'Error cambiando contraseña';
      toast.error(message);
      throw new Error(message);
    }
  }

  /**
   * Actualiza el perfil del usuario autenticado
   * @param {Object} profileData - { first_name, last_name, phone }
   */
  async updateProfile(profileData) {
    try {
      console.log('🔧 [authService] Actualizando perfil con datos:', profileData);
      const response = await authAPI.put('/users/profile', profileData);
      const user = response.data.user;
      console.log('🔧 [authService] Respuesta del backend:', user);
      
      // Actualizar storage local con los datos reales del backend
      this.user = user;
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      console.log('✅ [authService] Usuario actualizado en localStorage:', user);
      
      toast.success('Perfil actualizado correctamente');
      return user;
    } catch (error) {
      console.error('❌ [authService] Error actualizando perfil:', error);
      const message = error.response?.data?.error || 'Error actualizando perfil';
      toast.error(message);
      throw new Error(message);
    }
  }

  /**
   * Sube un avatar para el usuario autenticado
   * @param {File} file - Archivo de imagen del avatar
   */
  async uploadAvatar(file) {
    try {
      console.log('🔧 [authService] Subiendo avatar...');
      console.log('🔧 [authService] Archivo:', { name: file.name, size: file.size, type: file.type });
      
      const formData = new FormData();
      formData.append('avatar', file);
      
      console.log('🔧 [authService] FormData creado, enviando petición...');
      
      const response = await authAPI.post('/users/profile/avatar', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      
      console.log('🔧 [authService] Response status:', response.status);
      console.log('🔧 [authService] Response headers:', response.headers);
      console.log('✅ [authService] Avatar subido exitosamente:', response.data);
      
      toast.success('Avatar actualizado correctamente');
      return response.data;
    } catch (error) {
      console.error('❌ [authService] Error subiendo avatar:', error);
      console.error('❌ [authService] Error response:', error.response?.data);
      console.error('❌ [authService] Error status:', error.response?.status);
      
      const message = error.response?.data?.error || 'Error subiendo avatar';
      toast.error(message);
      throw new Error(message);
    }
  }

  /**
   * Guarda los datos de autenticación en localStorage
   * @param {Object} authData - Datos de autenticación del servidor
   */
  saveAuthData(authData) {
    console.log('🔧 [authService] Guardando datos de autenticación:', authData);
    
    // Support both flat { access_token, user } and nested { tokens: { access_token }, user }
    const tokens = authData.tokens || authData;
    const access_token = tokens.access_token;
    const refresh_token = tokens.refresh_token;
    const expires_at = tokens.expires_at;
    const user = authData.user;

    if (!access_token) {
      console.error('❌ [authService] No se encontró access_token en la respuesta');
      throw new Error('Token de acceso no encontrado en la respuesta');
    }

    if (!user) {
      console.error('❌ [authService] No se encontró user en la respuesta');
      throw new Error('Datos de usuario no encontrados en la respuesta');
    }

    // Normalize expires_at to Unix timestamp (backend returns ISO8601 string from time.Time)
    let expiresTimestamp = expires_at;
    if (typeof expires_at === 'string' && isNaN(Number(expires_at))) {
      expiresTimestamp = Math.floor(new Date(expires_at).getTime() / 1000);
    }

    this.token = access_token;
    this.refreshToken_value = refresh_token;
    this.user = user;
    this.expiresAt = expiresTimestamp;

    localStorage.setItem(TOKEN_KEY, access_token);
    if (refresh_token) {
      localStorage.setItem(REFRESH_TOKEN_KEY, refresh_token);
    }
    localStorage.setItem(USER_KEY, JSON.stringify(user));
    localStorage.setItem(EXPIRES_AT_KEY, expiresTimestamp?.toString() ?? '');
    
    console.log('✅ [authService] Datos de autenticación guardados correctamente');
  }

  /**
   * Limpia todos los datos de autenticación y cache del usuario
   */
  clearAuthData() {
    this.token = null;
    this.refreshToken_value = null;
    this.user = null;
    this.expiresAt = null;
    
    // Limpiar datos de autenticación
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    localStorage.removeItem(EXPIRES_AT_KEY);
    
    // Limpiar cache de datos del usuario
    localStorage.removeItem('dataChanged');
    localStorage.removeItem('financial_gamification');
    localStorage.removeItem('gamification_analytics');
    localStorage.removeItem('ai_insights_cache');
    localStorage.removeItem('health_score_cache');
    
    // Limpiar cache del dataService
    try {
      dataService.clearCache();
    } catch (error) {
      console.warn('Error limpiando cache del dataService:', error);
    }
    
    console.log('🧹 Cache y datos del usuario limpiados completamente');
  }

  /**
   * Verifica si el usuario está autenticado
   */
  isAuthenticated() {
    return this.token && this.user && this.isTokenValid();
  }

  /**
   * Obtiene el usuario actual
   */
  getCurrentUser() {
    return this.user;
  }

  /**
   * Obtiene el token actual
   */
  getToken() {
    return this.token;
  }

  /**
   * Obtiene headers de autorización para requests manuales
   */
  getAuthHeaders() {
    if (this.token && this.isTokenValid()) {
      const headers = {
        Authorization: `Bearer ${this.token}`,
      };
      
      // Agregar X-Caller-ID si tenemos usuario
      // Intentar diferentes propiedades que podría tener el ID
      let userId = null;
      if (this.user) {
        userId = this.user.id || this.user.ID || this.user.user_id || this.user.userId;
      }
      
      if (userId) {
        headers['X-Caller-ID'] = userId.toString();
      }
      
      return headers;
    }
    return {};
  }

  /**
   * Verifica si el usuario tiene un rol específico (para futuras implementaciones)
   */
  hasRole(role) {
    return this.user?.roles?.includes(role) || false;
  }

  /**
   * Obtiene información de la sesión
   */
  getSessionInfo() {
    return {
      isAuthenticated: this.isAuthenticated(),
      user: this.user,
      expiresAt: this.expiresAt ? new Date(parseInt(this.expiresAt) * 1000) : null,
      timeUntilExpiry: this.expiresAt ? (parseInt(this.expiresAt) - Math.floor(Date.now() / 1000)) : 0,
    };
  }
}

// Crear instancia única
const authService = new AuthService();

// Funciones de utilidad para facilitar el uso
export const authAPI_instance = authAPI; // Para requests manuales
export const isAuthenticated = () => authService.isAuthenticated();
export const getCurrentUser = () => authService.getCurrentUser();
export const getAuthHeaders = () => authService.getAuthHeaders();


export default authService; 