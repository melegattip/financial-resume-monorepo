import React, { useState, useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import { FaUser, FaBell, FaLock, FaDownload, FaTrash, FaShieldAlt, FaCamera, FaTimes, FaUserCog, FaHistory, FaCog, FaExclamationTriangle, FaChevronDown } from 'react-icons/fa';
import toast from 'react-hot-toast';
import TwoFASetup from '../components/TwoFASetup';
import ChangePasswordModal from '../components/ChangePasswordModal';
import { useAuth } from '../contexts/AuthContext';
import { getAvatarUrl } from '../utils/avatarUtils';
import TenantSettings from './TenantSettings';
import AuditLogs from './AuditLogs';
import { useTenant } from '../contexts/TenantContext';
import authService from '../services/authService';
import tenantService from '../services/tenantService';

// ─── Notifications Tab ────────────────────────────────────────────────────────

function Toggle({ on, disabled, onChange, label }) {
  return (
    <button
      onClick={onChange}
      disabled={disabled}
      aria-label={label}
      className={`relative flex-shrink-0 w-11 h-6 rounded-full transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 ${
        on ? 'bg-blue-600' : 'bg-gray-200 dark:bg-gray-700'
      } ${disabled ? 'opacity-40 cursor-not-allowed' : ''}`}
    >
      <span
        className={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow transition-transform duration-200 ${
          on ? 'translate-x-5' : 'translate-x-0'
        }`}
      />
    </button>
  );
}

function NotificationsTab() {
  const [notifSettings, setNotifSettings] = useState(null);
  const [saving, setSaving] = useState(null);

  useEffect(() => {
    authService.getNotifications()
      .then(data => setNotifSettings(data))
      .catch(() => toast.error('Error cargando preferencias de notificaciones'));
  }, []);

  const handleToggle = async (key) => {
    if (!notifSettings || saving) return;
    const newValue = !notifSettings[key];
    setNotifSettings(prev => ({ ...prev, [key]: newValue }));
    setSaving(key);
    try {
      await authService.updateNotifications({ ...notifSettings, [key]: newValue });
    } catch {
      setNotifSettings(prev => ({ ...prev, [key]: !newValue }));
      toast.error('Error actualizando preferencia');
    } finally {
      setSaving(null);
    }
  };

  if (!notifSettings) {
    return <div className="py-8 text-center text-gray-500 dark:text-gray-400 text-sm">Cargando…</div>;
  }

  const emailOn = notifSettings.email_notifications;

  return (
    <div className="space-y-1">
      <h3 className="text-lg font-semibold text-fr-gray-900 dark:text-gray-100 mb-4">Notificaciones</h3>

      {/* Parent row */}
      <div className="flex items-center justify-between gap-4 py-4 border-b border-gray-100 dark:border-gray-700">
        <div className="min-w-0">
          <p className="text-sm font-medium text-gray-900 dark:text-gray-100">Notificaciones por email</p>
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Recibí alertas importantes en tu casilla de correo.</p>
        </div>
        <Toggle on={emailOn} disabled={saving === 'email_notifications'} onChange={() => handleToggle('email_notifications')} label="Notificaciones por email" />
      </div>

      {/* Sub-rows — indented, disabled when parent is off */}
      <div className="pl-5 border-l-2 border-gray-100 dark:border-gray-700 space-y-0">
        <div className={`flex items-center justify-between gap-4 py-3.5 border-b border-gray-100 dark:border-gray-700 transition-opacity ${!emailOn ? 'opacity-40' : ''}`}>
          <div className="min-w-0">
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">Alertas de presupuesto</p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Avisá cuando estés cerca o superes el límite de un presupuesto.</p>
          </div>
          <Toggle on={notifSettings.budget_alerts} disabled={!emailOn || saving === 'budget_alerts'} onChange={() => handleToggle('budget_alerts')} label="Alertas de presupuesto" />
        </div>

        <div className={`flex items-center justify-between gap-4 py-3.5 transition-opacity ${!emailOn ? 'opacity-40' : ''}`}>
          <div className="min-w-0">
            <p className="text-sm font-medium text-gray-900 dark:text-gray-100">Nuevos dispositivos conectados</p>
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Notificá cada vez que tu cuenta inicia sesión desde un dispositivo.</p>
          </div>
          <Toggle on={notifSettings.login_notifications} disabled={!emailOn || saving === 'login_notifications'} onChange={() => handleToggle('login_notifications')} label="Nuevos dispositivos conectados" />
        </div>
      </div>
    </div>
  );
}

// ─── Deletion Section ─────────────────────────────────────────────────────────

function DeletionSection({ user, myRole, currentTenant }) {
  const [expanded, setExpanded] = useState(false);
  const [confirmText, setConfirmText] = useState('');
  const [password, setPassword] = useState('');
  const [deleting, setDeleting] = useState(false);
  const [memberCount, setMemberCount] = useState(null); // null = loading
  const [alsoDeleteSpace, setAlsoDeleteSpace] = useState(false);

  const isOwner = myRole === 'owner' && !!currentTenant;
  const confirmTarget = user?.email || '';
  const canDelete = confirmText === confirmTarget && password.length > 0;

  useEffect(() => {
    if (!isOwner) { setMemberCount(0); return; }
    tenantService.listMembers()
      .then(members => setMemberCount((members || []).length))
      .catch(() => setMemberCount(1));
  }, [isOwner]);

  const willTransferOwnership = isOwner && memberCount > 1;
  const willDeleteTenant      = isOwner && memberCount === 1;

  const getTitle = () => {
    if (willDeleteTenant || (willTransferOwnership && alsoDeleteSpace)) return 'Eliminar cuenta y espacio';
    return 'Eliminar tu cuenta';
  };

  const getDescription = () => {
    if (willTransferOwnership) {
      return (
        <>
          Sos el propietario de <span className="font-semibold">"{currentTenant.name}"</span>.
          {' '}Al eliminar tu cuenta, el miembro con mayor rango o antigüedad pasará a ser el nuevo propietario.
        </>
      );
    }
    if (willDeleteTenant) {
      return (
        <>
          Sos el único miembro de <span className="font-semibold">"{currentTenant.name}"</span>.
          {' '}Al eliminar tu cuenta también se eliminará el espacio y todos sus datos.
        </>
      );
    }
    return 'Se eliminarán permanentemente tu usuario y todos tus datos personales.';
  };

  const handleDelete = async () => {
    if (!canDelete || !password || deleting) return;
    if (!window.confirm('Esta acción es irreversible. ¿Confirmar eliminación?')) return;
    setDeleting(true);
    try {
      if (willTransferOwnership && alsoDeleteSpace) {
        await tenantService.deleteMyTenant();
      }
      await authService.deleteAccount(password);
      toast.success('Cuenta eliminada. Hasta luego.');
      await authService.logout();
      window.location.href = '/login';
    } catch (err) {
      const msg = err?.response?.data?.error || err?.message || 'Error al eliminar la cuenta';
      if (msg.includes('invalid') || msg.includes('password') || msg.includes('contraseña')) {
        toast.error('Contraseña incorrecta');
      } else {
        toast.error(msg);
      }
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="border border-red-200 dark:border-red-800 rounded-fr overflow-hidden">
      {/* Collapsible header */}
      <button
        onClick={() => setExpanded(v => !v)}
        className="w-full flex items-center justify-between px-4 py-3 bg-red-50 dark:bg-red-900/20 text-left hover:bg-red-100 dark:hover:bg-red-900/30 transition-colors"
      >
        <span className="flex items-center gap-2 text-sm font-medium text-red-700 dark:text-red-400">
          <FaExclamationTriangle className="w-4 h-4 flex-shrink-0" />
          Eliminar cuenta
        </span>
        <FaChevronDown className={`w-3.5 h-3.5 text-red-500 transition-transform duration-200 ${expanded ? 'rotate-180' : ''}`} />
      </button>

      {/* Expandable body */}
      {expanded && (
        <div className="px-4 pb-4 pt-3 bg-red-50 dark:bg-red-900/20 border-t border-red-200 dark:border-red-800 space-y-4">
          {/* Description */}
          <div className="flex items-start gap-2.5">
            <FaExclamationTriangle className="w-4 h-4 text-red-500 mt-0.5 flex-shrink-0" />
            <div className="space-y-1">
              <p className="text-sm font-semibold text-red-900 dark:text-red-100">{getTitle()}</p>
              <p className="text-sm text-red-700 dark:text-red-300">{getDescription()}</p>
              <p className="text-xs text-red-600 dark:text-red-400 opacity-80">Esta acción es irreversible.</p>
            </div>
          </div>

          {/* Option to also delete space (only for transfer case) */}
          {willTransferOwnership && isOwner && memberCount !== null && (
            <label className="flex items-start gap-2.5 cursor-pointer group">
              <input
                type="checkbox"
                checked={alsoDeleteSpace}
                onChange={e => setAlsoDeleteSpace(e.target.checked)}
                className="mt-0.5 accent-red-600 w-4 h-4 flex-shrink-0"
              />
              <span className="text-sm text-red-700 dark:text-red-300">
                También eliminar el espacio <span className="font-semibold">"{currentTenant.name}"</span> y todos sus datos
              </span>
            </label>
          )}

          {/* Confirm inputs */}
          <div className="space-y-3">
            <div>
              <label className="block text-xs text-red-700 dark:text-red-400 mb-1">
                Escribí <span className="font-mono font-semibold">{confirmTarget}</span> para confirmar:
              </label>
              <input
                type="text"
                value={confirmText}
                onChange={(e) => setConfirmText(e.target.value)}
                placeholder={confirmTarget}
                className="w-full px-3 py-2 border border-red-300 dark:border-red-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-red-500 text-sm"
              />
            </div>
            <div>
              <label className="block text-xs text-red-700 dark:text-red-400 mb-1">
                Contraseña actual:
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full px-3 py-2 border border-red-300 dark:border-red-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-red-500 text-sm"
              />
            </div>
            <button
              onClick={handleDelete}
              disabled={!canDelete || deleting}
              className="flex items-center gap-2 px-4 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded-lg disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <FaTrash className="w-3.5 h-3.5" />
              {deleting ? 'Eliminando…' : getTitle()}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

const Settings = () => {
  const location = useLocation();
  const { hasPermission, currentTenant, myRole } = useTenant();
  const [activeTab, setActiveTab] = useState('ajustes');
  const [show2FAModal, setShow2FAModal] = useState(false);
  const [showChangePasswordModal, setShowChangePasswordModal] = useState(false);
  const { user, updateProfile, uploadAvatar } = useAuth();
  const [loading, setLoading] = useState(false);
  const fileInputRef = useRef(null);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [showEmailWarning, setShowEmailWarning] = useState(false);


  
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
    
    if (tabParam && ['ajustes', 'security', 'espacio', 'actividad'].includes(tabParam)) {
      setActiveTab(tabParam);
    }
  }, [location.search]);

  const tabs = [
    { id: 'ajustes', label: 'Ajustes', icon: FaCog },
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

  const emailChanged = settings.profile.email.trim() !== (user?.email || '');

  const handleSave = async () => {
    if (!settings.profile.name.trim()) {
      toast.error('El nombre es obligatorio');
      return;
    }
    // If email changed, show warning first
    if (emailChanged) {
      setShowEmailWarning(true);
      return;
    }
    await doSave();
  };

  const doSave = async () => {
    setShowEmailWarning(false);
    setLoading(true);
    try {
      const [first_name, ...rest] = settings.profile.name.trim().split(' ');
      const last_name = rest.join(' ');

      if (settings.profile.avatar instanceof File) {
        try {
          const result = await uploadAvatar(settings.profile.avatar);
          if (!result.success) throw new Error(result.error || 'Error subiendo avatar');
        } catch (error) {
          toast.error(`Error subiendo avatar: ${error.message}`);
        }
      }

      const profileData = {
        first_name,
        last_name,
        phone: settings.profile.phone,
        email: settings.profile.email.trim(),
      };

      const result = await updateProfile(profileData);
      if (result.success) {
        if (result.emailChanged) {
          toast.success('Email actualizado. Revisá tu casilla para verificarlo. Tu sesión se cerrará.');
          setTimeout(async () => {
            await authService.logout();
            window.location.href = '/login';
          }, 3000);
        } else {
          toast.success('Perfil actualizado');
          setAvatarPreview(null);
        }
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
      <div className="flex overflow-x-auto scrollbar-hide border-b border-gray-200 dark:border-gray-700">
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
        {activeTab === 'ajustes' && (
          <div className="space-y-8">
            {/* ── Perfil ── */}
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
                  {emailChanged && (
                    <span className="ml-2 text-xs font-normal text-amber-600 dark:text-amber-400">
                      · cambiar email cerrará tu sesión
                    </span>
                  )}
                </label>
                <input
                  type="email"
                  value={settings.profile.email}
                  onChange={(e) => updateSetting('profile', 'email', e.target.value)}
                  className={`input ${emailChanged ? 'border-amber-400 dark:border-amber-500 focus:ring-amber-400' : ''}`}
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

            {/* ── Divider ── */}
            <hr className="border-gray-200 dark:border-gray-700" />

            {/* ── Notificaciones ── */}
            <NotificationsTab />
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

              {/* ── Zona de peligro ── */}
              <DeletionSection user={user} myRole={myRole} currentTenant={currentTenant} />
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

      {/* Email change warning modal */}
      {showEmailWarning && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4">
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-sm w-full p-6 space-y-4">
            <div className="flex items-start gap-3">
              <FaBell className="w-5 h-5 text-amber-500 flex-shrink-0 mt-0.5" />
              <div>
                <h3 className="text-base font-semibold text-gray-900 dark:text-gray-100">Cambiar email</h3>
                <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                  Tu sesión se cerrará y se enviará un correo de verificación a{' '}
                  <span className="font-medium text-gray-900 dark:text-gray-100">{settings.profile.email}</span>.
                  Necesitarás verificarlo para volver a iniciar sesión.
                </p>
              </div>
            </div>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setShowEmailWarning(false)}
                className="px-4 py-2 text-sm text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
              >
                Cancelar
              </button>
              <button
                onClick={doSave}
                className="px-4 py-2 text-sm bg-amber-500 hover:bg-amber-600 text-white rounded-lg transition-colors"
              >
                Confirmar y cerrar sesión
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Settings; 