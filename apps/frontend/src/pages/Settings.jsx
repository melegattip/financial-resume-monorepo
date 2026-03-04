import React, { useState, useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import { FaUser, FaBell, FaLock, FaDownload, FaTrash, FaShieldAlt, FaCamera, FaTimes, FaUserCog, FaHistory } from 'react-icons/fa';
import toast from 'react-hot-toast';
import TwoFASetup from '../components/TwoFASetup';
import ChangePasswordModal from '../components/ChangePasswordModal';
import { useAuth } from '../contexts/AuthContext';
import { getAvatarUrl } from '../utils/avatarUtils';
import TenantSettings from './TenantSettings';
import AuditLogs from './AuditLogs';
import { useTenant } from '../contexts/TenantContext';

const Settings = () => {
  const location = useLocation();
  const { hasPermission } = useTenant();
  const [activeTab, setActiveTab] = useState('profile');
  const [show2FAModal, setShow2FAModal] = useState(false);
  const [showChangePasswordModal, setShowChangePasswordModal] = useState(false);
  const { user, updateProfile, uploadAvatar } = useAuth();
  const [loading, setLoading] = useState(false);
  const fileInputRef = useRef(null);
  const [avatarPreview, setAvatarPreview] = useState(null);


  
  // Usar datos reales del usuario autenticado
  const [settings, setSettings] = useState({
    profile: {
      name: user ? `${user.first_name || ''} ${user.last_name || ''}`.trim() : '',
      email: user?.email || '',
      phone: user?.phone || '',
      avatar: user?.avatar || null,
    },
    notifications: {
      emailNotifications: true,
      pushNotifications: false,
      weeklyReports: true,
      expenseAlerts: true,
    },
  });

  // Actualizar settings cuando cambie el usuario
  useEffect(() => {
    if (user) {
      console.log('🔧 [Settings] Usuario actualizado:', user);
      console.log('🔧 [Settings] Avatar del usuario:', user.avatar);
      setSettings(prev => ({
        ...prev,
        profile: {
          name: `${user.first_name || ''} ${user.last_name || ''}`.trim(),
          email: user.email || '',
          phone: user.phone || '',
          avatar: user.avatar || null,
        }
      }));
    }
  }, [user]);

  // Manejar parámetros de URL para abrir pestaña específica
  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const tabParam = searchParams.get('tab');
    
    if (tabParam && ['profile', 'notifications', 'security', 'espacio', 'actividad'].includes(tabParam)) {
      setActiveTab(tabParam);
    }
  }, [location.search]);

  const tabs = [
    { id: 'profile', label: 'Perfil', icon: FaUser },
    { id: 'notifications', label: 'Notificaciones', icon: FaBell, disabled: true },
    { id: 'security', label: 'Seguridad', icon: FaLock },
    { id: 'espacio', label: 'Espacio', icon: FaUserCog },
    ...(hasPermission('view_audit_logs') ? [{ id: 'actividad', label: 'Actividad', icon: FaHistory }] : []),
  ];

  const handleAvatarChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      // Validar tipo de archivo
      if (!file.type.startsWith('image/')) {
        toast.error('Por favor selecciona una imagen válida');
        return;
      }
      
      // Validar tamaño (máximo 5MB)
      if (file.size > 5 * 1024 * 1024) {
        toast.error('La imagen debe ser menor a 5MB');
        return;
      }
      
      // Crear preview
      const reader = new FileReader();
      reader.onload = (e) => {
        setAvatarPreview(e.target.result);
        setSettings(prev => ({
          ...prev,
          profile: {
            ...prev.profile,
            avatar: file
          }
        }));
      };
      reader.readAsDataURL(file);
    }
  };

  const removeAvatar = () => {
    setAvatarPreview(null);
    setSettings(prev => ({
      ...prev,
      profile: {
        ...prev.profile,
        avatar: null
      }
    }));
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSave = async () => {
    // Validación mínima
    if (!settings.profile.name.trim()) {
      toast.error('El nombre es obligatorio');
      return;
    }
    setLoading(true);
    try {
      // Separar nombre completo en first_name y last_name
      const [first_name, ...rest] = settings.profile.name.trim().split(' ');
      const last_name = rest.join(' ');
      
      // Subir avatar si hay un archivo nuevo
      if (settings.profile.avatar instanceof File) {
        try {
          console.log('🔧 [Settings] Subiendo avatar usando authService...');
          console.log('🔧 [Settings] Archivo a subir:', {
            name: settings.profile.avatar.name,
            size: settings.profile.avatar.size,
            type: settings.profile.avatar.type
          });
          
          const result = await uploadAvatar(settings.profile.avatar);
          console.log('🔧 [Settings] Resultado completo del uploadAvatar:', result);
          
          if (result.success) {
            console.log('✅ [Settings] Avatar subido exitosamente:', result);
          } else {
            console.error('❌ [Settings] Upload falló:', result.error);
            throw new Error(result.error || 'Error subiendo avatar');
          }
        } catch (error) {
          console.error('❌ [Settings] Error subiendo avatar:', error);
          toast.error(`Error subiendo avatar: ${error.message}`);
        }
      }
      
      // Actualizar datos del perfil (sin avatar, ya que se maneja por separado)
      const profileData = {
        id: user?.id,
        email: user?.email,
        first_name,
        last_name,
        phone: settings.profile.phone,
      };
      
      const result = await updateProfile(profileData);
      if (result.success) {
        toast.success('Perfil actualizado');
        // Limpiar preview después de guardar
        setAvatarPreview(null);
      } else {
        toast.error(result.error || 'Error actualizando perfil');
      }
    } catch (e) {
      toast.error('Error actualizando perfil');
    } finally {
      setLoading(false);
    }
  };

  const handleExportData = () => {
    toast.success('Exportación iniciada. Recibirás un email con tus datos.');
  };

  const handleDeleteAccount = () => {
    if (window.confirm('¿Estás seguro de que quieres eliminar tu cuenta? Esta acción no se puede deshacer.')) {
      toast.error('Funcionalidad no implementada en demo');
    }
  };

  const updateSetting = (section, key, value) => {
    setSettings(prev => ({
      ...prev,
      [section]: {
        ...prev[section],
        [key]: value,
      },
    }));
  };

  return (
    <div className="space-y-4">
      {/* Tabs — mismo estilo que el dashboard */}
      <div className="flex border-b border-gray-200 dark:border-gray-700">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              onClick={() => !tab.disabled && setActiveTab(tab.id)}
              disabled={tab.disabled}
              className={`flex items-center gap-1.5 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors -mb-px whitespace-nowrap ${
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
              } ${tab.disabled ? 'opacity-40 cursor-not-allowed' : ''}`}
            >
              <Icon className="w-4 h-4" />
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* Contenido de tabs */}
      {activeTab === 'espacio' && <TenantSettings widgetMode />}
      {activeTab === 'actividad' && <AuditLogs />}
      <div className={`card${activeTab === 'espacio' || activeTab === 'actividad' ? ' hidden' : ''}`}>
        {activeTab === 'profile' && (
          <div className="space-y-6">
            <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100">Información del Perfil</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                  Nombre completo
                </label>
                <input
                  type="text"
                  value={settings.profile.name}
                  onChange={(e) => updateSetting('profile', 'name', e.target.value)}
                  className="input"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                  Email
                </label>
                <input
                  type="email"
                  value={settings.profile.email}
                  onChange={(e) => updateSetting('profile', 'email', e.target.value)}
                  className="input"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                  Teléfono
                </label>
                <input
                  type="tel"
                  value={settings.profile.phone}
                  onChange={(e) => updateSetting('profile', 'phone', e.target.value)}
                  className="input"
                />
              </div>

              <div className="col-span-full md:col-span-1">
                <label className="block text-sm font-medium text-fr-gray-700 dark:text-gray-300 mb-2">
                  Avatar
                </label>
                <div className="flex items-center space-x-4">
                  <div className="relative w-16 h-16 rounded-full overflow-hidden bg-fr-gray-200 dark:bg-fr-gray-700">
                    {avatarPreview ? (
                      <img src={avatarPreview} alt="Avatar Preview" className="w-full h-full object-cover" />
                    ) : user?.avatar ? (
                      <img 
                        src={getAvatarUrl(user.avatar)} 
                        alt="Avatar" 
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          console.error('❌ Error cargando avatar:', e.target.src);
                          console.error('❌ Avatar path original:', user.avatar);
                          e.target.style.display = 'none';
                        }}
                        onLoad={() => {
                          console.log('✅ Avatar cargado exitosamente');
                        }}
                      />
                    ) : (
                      <div className="flex items-center justify-center h-full text-fr-gray-500 dark:text-gray-400">
                        <FaUser className="w-8 h-8" />
                      </div>
                    )}
                    
                    {/* Overlay solo cuando NO hay avatar */}
                    {!avatarPreview && !user?.avatar && (
                      <label htmlFor="avatar-upload" className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 rounded-full cursor-pointer">
                        <FaCamera className="w-6 h-6 text-white" />
                        <input
                          type="file"
                          id="avatar-upload"
                          accept="image/*"
                          onChange={handleAvatarChange}
                          className="hidden"
                          ref={fileInputRef}
                        />
                      </label>
                    )}
                    
                    {/* Input oculto para cuando HAY avatar */}
                    {(avatarPreview || user?.avatar) && (
                      <input
                        type="file"
                        id="avatar-upload-hidden"
                        accept="image/*"
                        onChange={handleAvatarChange}
                        className="hidden"
                        ref={fileInputRef}
                      />
                    )}
                  </div>
                  
                  {/* Botón para cambiar avatar cuando HAY avatar */}
                  {(avatarPreview || user?.avatar) && (
                    <button
                      onClick={() => fileInputRef.current?.click()}
                      className="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                    >
                      Cambiar avatar
                    </button>
                  )}
                  {avatarPreview && (
                    <button
                      onClick={removeAvatar}
                      className="flex items-center space-x-2 text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-500"
                    >
                      <FaTimes className="w-4 h-4" />
                      <span>Eliminar Avatar</span>
                    </button>
                  )}
                </div>
              </div>
            </div>

            <div className="flex justify-end">
              <button onClick={handleSave} className="btn-primary" disabled={loading}>
                {loading ? 'Guardando...' : 'Guardar Cambios'}
              </button>
            </div>
          </div>
        )}

        {activeTab === 'notifications' && (
          <div className="space-y-6">
            <div className="flex items-center space-x-3 mb-6">
              <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100">Configuración de Notificaciones</h3>
              <span className="px-2 py-1 text-xs bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200 rounded-full">
                Próximamente
              </span>
            </div>
            
            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-6 text-center">
              <FaBell className="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" />
              <h4 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                Configuración de Notificaciones
              </h4>
              <p className="text-gray-600 dark:text-gray-400 mb-4">
                Esta funcionalidad estará disponible próximamente. Podrás configurar notificaciones por email, push y alertas personalizadas.
              </p>
              <div className="space-y-2 text-sm text-gray-500 dark:text-gray-400">
                <p>• Notificaciones por email</p>
                <p>• Notificaciones push en tiempo real</p>
                <p>• Reportes semanales automáticos</p>
                <p>• Alertas de gastos y presupuestos</p>
              </div>
            </div>
          </div>
        )}



        {activeTab === 'security' && (
          <div className="space-y-6">
            <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100">Seguridad y Privacidad</h3>
            
            <div className="space-y-6">
              <div className="border border-fr-gray-200 dark:border-gray-600 rounded-fr p-4">
                <h4 className="font-medium text-fr-gray-900 dark:text-gray-100 mb-2">Cambiar contraseña</h4>
                <p className="text-sm text-fr-gray-500 dark:text-gray-300 mb-4">
                  Actualiza tu contraseña regularmente para mantener tu cuenta segura
                </p>
                <button 
                  onClick={() => setShowChangePasswordModal(true)}
                  className="btn-outline"
                >
                  Cambiar Contraseña
                </button>
              </div>

              <div className="border border-fr-gray-200 dark:border-gray-600 rounded-fr p-4">
                <h4 className="font-medium text-fr-gray-900 dark:text-gray-100 mb-2">Autenticación de dos factores</h4>
                <p className="text-sm text-fr-gray-500 dark:text-gray-300 mb-4">
                  Agrega una capa extra de seguridad a tu cuenta
                </p>
                <button 
                  onClick={() => setShow2FAModal(true)}
                  className="btn-outline flex items-center space-x-2"
                >
                  <FaShieldAlt className="w-4 h-4" />
                  <span>Configurar 2FA</span>
                </button>
              </div>

              <div className="border border-fr-gray-200 dark:border-gray-600 rounded-fr p-4">
                <h4 className="font-medium text-fr-gray-900 dark:text-gray-100 mb-2">Exportar datos</h4>
                <p className="text-sm text-fr-gray-500 dark:text-gray-300 mb-4">
                  Descarga una copia de todos tus datos financieros
                </p>
                <button onClick={handleExportData} className="btn-outline flex items-center space-x-2">
                  <FaDownload className="w-4 h-4" />
                  <span>Exportar Datos</span>
                </button>
              </div>

              <div className="border border-red-200 dark:border-red-800 rounded-fr p-4 bg-red-50 dark:bg-red-900/20">
                <h4 className="font-medium text-red-900 dark:text-red-100 mb-2">Zona de peligro</h4>
                <p className="text-sm text-red-700 dark:text-red-200 mb-4">
                  Una vez que elimines tu cuenta, no hay vuelta atrás. Por favor, ten cuidado.
                </p>
                <button 
                  onClick={handleDeleteAccount}
                  className="bg-red-600 hover:bg-red-700 dark:bg-red-700 dark:hover:bg-red-800 text-white font-medium py-2 px-4 rounded-fr transition-colors flex items-center space-x-2"
                >
                  <FaTrash className="w-4 h-4" />
                  <span>Eliminar Cuenta</span>
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Modal de configuración 2FA */}
      {show2FAModal && (
        <TwoFASetup
          isOpen={show2FAModal}
          onClose={() => setShow2FAModal(false)}
        />
      )}

      {/* Modal de cambio de contraseña */}
      <ChangePasswordModal
        isOpen={showChangePasswordModal}
        onClose={() => setShowChangePasswordModal(false)}
      />
    </div>
  );
};

export default Settings; 