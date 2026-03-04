import React, { useState, useRef, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useTenant } from '../../contexts/TenantContext';
import PeriodFilter from './PeriodFilter';
import ThemeToggle from '../ThemeToggle';
import GamificationWidget from '../GamificationWidget';
import { FaUser, FaSignOutAlt, FaHome, FaBrain, FaPlusCircle, FaMinusCircle, FaFolderOpen, FaFileAlt, FaCog, FaChartPie, FaBullseye, FaRedo, FaTrophy, FaBell, FaLock, FaChevronDown, FaUserCog, FaHistory, FaExchangeAlt, FaCheck } from 'react-icons/fa';
import { getAvatarUrl } from '../../utils/avatarUtils';
import toast from 'react-hot-toast';

const Header = () => {
  const { user, logout } = useAuth();
  const { currentTenant, myRole, hasPermission, availableTenants, switching, switchTenant } = useTenant();
  const location = useLocation();
  const navigate = useNavigate();
  const [showUserMenu, setShowUserMenu] = useState(false);
  const [showTenantMenu, setShowTenantMenu] = useState(false);
  const userMenuRef = useRef(null);



  // Cerrar menú al hacer click fuera
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target)) {
        setShowUserMenu(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Mapeo de rutas a información de página
  const getPageInfo = (pathname) => {
    const routes = {
      '/dashboard': {
        title: 'Cuentas',
        icon: FaHome
      },
      '/insights': { 
        title: 'Análisis Inteligente', 
        icon: FaBrain
      },
      '/expenses': { 
        title: 'Gastos', 
        icon: FaMinusCircle
      },
      '/incomes': { 
        title: 'Ingresos', 
        icon: FaPlusCircle
      },
      '/categories': { 
        title: 'Categorías', 
        icon: FaFolderOpen
      },
      '/reports': { 
        title: 'Reportes', 
        icon: FaFileAlt
      },
      '/budgets': { 
        title: 'Presupuestos', 
        icon: FaChartPie
      },
      '/savings-goals': { 
        title: 'Metas de Ahorro', 
        icon: FaBullseye
      },
      '/recurring-transactions': { 
        title: 'Transacciones Recurrentes', 
        icon: FaRedo
      },
      '/achievements': { 
        title: 'Logros y Progreso', 
        icon: FaTrophy
      },
      '/settings': {
        title: 'Configuración',
        icon: FaCog
      },
    };
    return routes[pathname] || { 
      title: 'Niloft', 
      subtitle: 'Gestión financiera inteligente',
      icon: FaHome
    };
  };

  const pageInfo = getPageInfo(location.pathname);
  const PageIcon = pageInfo.icon;

  const handleLogout = () => {
    logout();
  };

  const handleMenuClick = (action) => {
    setShowUserMenu(false);
    setShowTenantMenu(false);
    switch (action) {
      case 'profile':
        navigate('/settings?tab=profile');
        break;
      case 'security':
        navigate('/settings?tab=security');
        break;
      case 'logout':
        handleLogout();
        break;
      default:
        break;
    }
  };

  const handleSwitchTenant = async (tenantId) => {
    if (tenantId === currentTenant?.id) return;
    setShowUserMenu(false);
    setShowTenantMenu(false);
    try {
      await switchTenant(tenantId);
      toast.success('Espacio cambiado correctamente');
      window.location.href = '/dashboard';
    } catch (err) {
      toast.error(err?.response?.data?.error || 'Error al cambiar de espacio');
    }
  };

  return (
    <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-40 shadow-sm dark:shadow-gray-900/20 transition-colors duration-300">
      <div className="pl-14 pr-4 lg:px-6 xl:px-8">
        <div className="flex items-center justify-between h-12 sm:h-14">
          
          {/* Left Side - Page title and info */}
          <div className="flex items-center flex-1 min-w-0 mr-2 sm:mr-4">
            {/* Desktop/Tablet: Icono + título completo */}
            <div className="hidden sm:flex items-center space-x-3">
              <div className="p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
                <PageIcon className="w-5 h-5 text-blue-600 dark:text-blue-400" />
              </div>
              <div className="min-w-0">
                <h1 className="text-xl lg:text-2xl font-bold text-gray-900 dark:text-gray-100 truncate">
                  {pageInfo.title}
                </h1>
                <p className="text-sm text-gray-600 dark:text-gray-400 truncate hidden lg:block">
                  {pageInfo.subtitle}
                </p>
              </div>
            </div>
            
            {/* Mobile: Solo título compacto */}
            <div className="sm:hidden min-w-0 flex-1">
              <h1 className="text-sm font-bold text-gray-900 dark:text-gray-100 truncate">
                {pageInfo.title}
              </h1>
            </div>
          </div>

          {/* Right Side - Actions */}
          <div className="flex items-center space-x-1 sm:space-x-2 lg:space-x-3 flex-shrink-0">
            
            {/* Period Filter - Responsive sizing */}
            <div className="hidden sm:block">
              <PeriodFilter />
            </div>
            
            {/* Mobile Period Filter - Very compact */}
            <div className="sm:hidden">
              <div className="scale-75">
                <PeriodFilter />
              </div>
            </div>

            {/* Gamification Widget - Only desktop */}
            <div className="hidden lg:block">
              <GamificationWidget />
            </div>

            {/* Theme Toggle - Very compact on mobile */}
            <div className="scale-75 sm:scale-100">
              <ThemeToggle />
            </div>

            {/* User Menu - Responsive */}
            <div className="flex items-center space-x-0.5 sm:space-x-2" ref={userMenuRef}>
              
              {/* User info - Hidden on mobile */}
              <div className="hidden md:block text-right">
                <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate max-w-24 lg:max-w-32">
                  {user?.name || user?.email || 'Usuario'}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400 truncate max-w-24 lg:max-w-32">
                  {currentTenant ? `${currentTenant.name} · ${myRole || 'miembro'}` : (myRole || 'Miembro')}
                </p>
              </div>
              
              {/* Avatar with dropdown - Very compact on mobile */}
              <div className="relative">
                <button
                  onClick={() => setShowUserMenu(!showUserMenu)}
                  className="flex items-center space-x-1 p-1 sm:p-2 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
                  title="Menú de usuario"
                >
                  <div className="w-6 h-6 sm:w-8 sm:h-8 lg:w-10 lg:h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center">
                    {user?.avatar ? (
                      <img 
                        src={getAvatarUrl(user.avatar)} 
                        alt="Avatar" 
                        className="w-full h-full object-cover rounded-full"
                        onError={(e) => {
                          console.error('❌ [Header] Error cargando avatar:', e.target.src);
                          e.target.style.display = 'none';
                        }}
                      />
                    ) : (
                      <FaUser className="w-2.5 h-2.5 sm:w-4 sm:h-4 lg:w-5 lg:h-5 text-white" />
                    )}
                  </div>
                  <FaChevronDown className={`w-2.5 h-2.5 sm:w-3 sm:h-3 lg:w-4 lg:h-4 transition-transform ${showUserMenu ? 'rotate-180' : ''}`} />
                </button>
                
                {/* Dropdown Menu */}
                {showUserMenu && (
                  <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1 z-50">
                    {/* Profile Section */}
                    <div className="px-4 py-2 border-b border-gray-200 dark:border-gray-700">
                      <p className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        {user?.name || user?.email || 'Usuario'}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {user?.email}
                      </p>
                    </div>
                    
                    {/* Menu Items */}
                    <div className="py-1">
                      <button
                        onClick={() => handleMenuClick('profile')}
                        className="flex items-center w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                      >
                        <FaUser className="w-4 h-4 mr-3" />
                        Perfil
                      </button>

                      <button
                        onClick={() => handleMenuClick('security')}
                        className="flex items-center w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                      >
                        <FaLock className="w-4 h-4 mr-3" />
                        Seguridad
                      </button>

                      {/* Tenant Switcher — only shown when user belongs to multiple tenants */}
                      {availableTenants.length > 1 && (
                        <div className="border-t border-gray-100 dark:border-gray-700 mt-1 pt-1">
                          <button
                            onClick={() => setShowTenantMenu(v => !v)}
                            className="flex items-center justify-between w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                            disabled={switching}
                          >
                            <span className="flex items-center">
                              <FaExchangeAlt className="w-4 h-4 mr-3 text-blue-500" />
                              Cambiar espacio
                            </span>
                            <FaChevronDown className={`w-3 h-3 transition-transform ${showTenantMenu ? 'rotate-180' : ''}`} />
                          </button>
                          {showTenantMenu && (
                            <div className="bg-gray-50 dark:bg-gray-750 border-t border-gray-100 dark:border-gray-700">
                              {availableTenants.map(t => (
                                <button
                                  key={t.id}
                                  onClick={() => handleSwitchTenant(t.id)}
                                  disabled={switching || t.id === currentTenant?.id}
                                  className="flex items-center justify-between w-full px-6 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors disabled:opacity-60"
                                >
                                  <span className="text-left min-w-0">
                                    <span className="block truncate text-gray-800 dark:text-gray-200 font-medium">{t.name}</span>
                                    <span className="block text-xs text-gray-500 dark:text-gray-400">{t.role}</span>
                                  </span>
                                  {t.id === currentTenant?.id && (
                                    <FaCheck className="w-3 h-3 text-green-500 flex-shrink-0 ml-2" />
                                  )}
                                </button>
                              ))}
                            </div>
                          )}
                        </div>
                      )}

                      <button
                        onClick={() => { setShowUserMenu(false); navigate('/settings?tab=espacio'); }}
                        className="flex items-center w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                      >
                        <FaUserCog className="w-4 h-4 mr-3" />
                        Espacio
                      </button>

                      {hasPermission('view_audit_logs') && (
                        <button
                          onClick={() => { setShowUserMenu(false); navigate('/settings?tab=actividad'); }}
                          className="flex items-center w-full px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                        >
                          <FaHistory className="w-4 h-4 mr-3" />
                          Actividad
                        </button>
                      )}

                      <button
                        disabled
                        className="flex items-center w-full px-4 py-2 text-sm text-gray-400 dark:text-gray-500 cursor-not-allowed opacity-50"
                        title="Próximamente"
                      >
                        <FaBell className="w-4 h-4 mr-3" />
                        Notificaciones
                      </button>
                    </div>
                    
                    {/* Logout Section */}
                    <div className="border-t border-gray-200 dark:border-gray-700 pt-1">
                      <button
                        onClick={() => handleMenuClick('logout')}
                        className="flex items-center w-full px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors"
                      >
                        <FaSignOutAlt className="w-4 h-4 mr-3" />
                        Cerrar sesión
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header; 