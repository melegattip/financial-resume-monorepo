/**
 * Servicio para cargar configuración dinámica desde el backend
 */
class ConfigService {
  constructor() {
    this.config = null;
    this.loading = false;
    this.error = null;
    this.loadPromise = null;  // Para evitar cargas múltiples simultáneas
  }

  /**
   * Carga la configuración desde el backend
   * @param {string} fallbackUrl - URL de fallback si no se puede cargar la configuración
   * @returns {Promise<Object>} - Configuración cargada
   */
  async loadConfig(fallbackUrl = null) {
    // Si ya hay una promesa de carga en curso, reutilizarla
    if (this.loadPromise) {
      return this.loadPromise;
    }

    if (this.config) {
      return this.config;
    }

    // Crear y almacenar la promesa de carga
    this.loadPromise = this._performLoad(fallbackUrl);
    
    try {
      const result = await this.loadPromise;
      this.loadPromise = null;  // Limpiar la promesa
      return result;
    } catch (error) {
      this.loadPromise = null;  // Limpiar la promesa en caso de error
      throw error;
    }
  }

  async _performLoad(fallbackUrl = null) {
    this.loading = true;
    this.error = null;

    try {
      // Importar configuración dinámica
      const envConfig = (await import('../config/environments')).default;
      
      // Determinar ambiente actual
      const currentEnv = envConfig.ENVIRONMENT;
      
      console.log(`🔍 [configService] Ambiente detectado: ${currentEnv}`);
      console.log(`🔍 [configService] Variables de entorno:`, {
        REACT_APP_ENVIRONMENT: process.env.REACT_APP_ENVIRONMENT,
        REACT_APP_API_URL: process.env.REACT_APP_API_URL,
        NODE_ENV: process.env.NODE_ENV
      });
      
      // En desarrollo o Docker, usar directamente las variables de entorno sin intentar cargar desde API
      if (currentEnv === 'development' || currentEnv === 'docker') {
        console.log(`🔧 [configService] Ambiente ${currentEnv} detectado, usando configuración local`);
        
        const config = {
          api_base_url: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
          gamification_url: process.env.REACT_APP_GAMIFICATION_URL || 'http://localhost:8080/api/v1',
          ai_service_url: process.env.REACT_APP_AI_SERVICE_URL || 'http://localhost:8080/api/v1',
          users_service_url: process.env.REACT_APP_USERS_SERVICE_URL || 'http://localhost:8080/api/v1',
          environment: currentEnv,
          version: '1.0.0'
        };
        
        console.log(`✅ Configuración de ${currentEnv} cargada:`, config);
        this.config = config;
        return config;
      }
      
      // En producción (render o production), usar el monolito
      if (currentEnv === 'render' || currentEnv === 'production') {
        const monolithUrl = 'https://financial-resume-monorepo.onrender.com/api/v1';
        const config = {
          api_base_url: monolithUrl,
          gamification_url: monolithUrl,
          ai_service_url: monolithUrl,
          users_service_url: monolithUrl,
          environment: currentEnv,
          version: '1.0.0'
        };

        console.log(`✅ Configuración de ${currentEnv} (monolito):`, config);
        this.config = config;
        return config;
      }

      // Fallback general
      console.warn('⚠️ Ambiente desconocido, usando configuración de fallback');
      const config = {
        api_base_url: fallbackUrl || 'http://localhost:8080/api/v1',
        gamification_url: fallbackUrl || 'http://localhost:8080/api/v1',
        ai_service_url: fallbackUrl || 'http://localhost:8080/api/v1',
        users_service_url: fallbackUrl || 'http://localhost:8080/api/v1',
        environment: 'development',
        version: '1.0.0'
      };
      
      this.config = config;
      return config;

    } catch (error) {
      console.error('❌ Error cargando configuración:', error);
      this.error = error;
      
      // Usar configuración de fallback en caso de error
      const fallbackConfig = {
        api_base_url: fallbackUrl,
        environment: 'development',
        version: '1.0.0'
      };
      
      this.config = fallbackConfig;
      return fallbackConfig;
    } finally {
      this.loading = false;
    }
  }

  /**
   * Obtiene la URL base del API
   * @returns {string} - URL base del API
   */
  getApiBaseUrl() {
    // Si ya hay configuración cargada, usarla
    if (this.config?.api_base_url) {
      return this.config.api_base_url;
    }
    
    // Si hay variable de entorno, usarla
    if (process.env.REACT_APP_API_URL) {
      return process.env.REACT_APP_API_URL;
    }
    
    // Fallback basado en detección de ambiente
    const hostname = window.location.hostname;
    if (hostname.includes('onrender.com') || hostname === 'financial.niloft.com') {
      return 'https://financial-resume-monorepo.onrender.com/api/v1';  // Render monolith
    } else {
      return 'http://localhost:8080/api/v1';  // Development
    }
  }

  /**
   * Obtiene el entorno actual
   * @returns {string} - Entorno (development, production, etc.)
   */
  getEnvironment() {
    return this.config?.environment || 'development';
  }

  /**
   * Obtiene la versión de la aplicación
   * @returns {string} - Versión
   */
  getVersion() {
    return this.config?.version || '1.0.0';
  }

  /**
   * Limpia la configuración cargada (útil para testing)
   */
  clearConfig() {
    this.config = null;
    this.loading = false;
    this.error = null;
  }
}

// Crear instancia singleton
const configService = new ConfigService();

export default configService; 