import environments from '../config/environments';

/**
 * Construye la URL completa para un avatar del usuario
 * @param {string} avatarPath - Path del avatar que puede venir con o sin /uploads/
 * @returns {string|null} - URL completa del avatar o null si no hay path
 */
export const getAvatarUrl = (avatarPath) => {
  if (!avatarPath) {
    console.log('🔧 [avatarUtils] Avatar path es null/undefined');
    return null;
  }

  const baseUrl = environments.USERS_API_URL;
  let cleanPath = avatarPath;

  // El avatarPath del backend ya viene con /uploads/, solo asegurar que empiece con /
  if (!cleanPath.startsWith('/')) {
    cleanPath = `/${cleanPath}`;
  }

  // Agregar cache busting para evitar problemas de cache del navegador
  const cacheBuster = Date.now();
  const fullUrl = `${baseUrl}${cleanPath}?v=${cacheBuster}`;
  
  console.log('🔧 [avatarUtils] Avatar URL construida:', { 
    avatarPath, 
    cleanPath, 
    baseUrl, 
    fullUrl 
  });
  
  return fullUrl;
};

/**
 * Verifica si un path de avatar es válido
 * @param {string} avatarPath - Path del avatar
 * @returns {boolean} - true si es válido
 */
export const isValidAvatarPath = (avatarPath) => {
  if (!avatarPath || typeof avatarPath !== 'string') {
    return false;
  }
  
  // Verificar que termine en una extensión de imagen válida
  const validExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp'];
  const lowerPath = avatarPath.toLowerCase();
  
  return validExtensions.some(ext => lowerPath.endsWith(ext));
};