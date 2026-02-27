import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { FaPlusCircle, FaMinusCircle, FaFolderOpen, FaBrain, FaFileAlt, FaCog, FaBars, FaTimes, FaHome, FaStar, FaChartPie, FaBullseye, FaRedo, FaTrophy, FaChevronLeft, FaChevronRight } from 'react-icons/fa';
import Brand from '../Brand';
import FeatureProgressIndicator from '../FeatureProgressIndicator';

const Sidebar = ({ isDesktopCollapsed = false, onDesktopToggle }) => {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);

  // Grupo 1: Transacciones principales
  const mainMenuItems = [
    { path: '/dashboard', icon: FaHome, label: 'Resumen' },
    { path: '/recurring-transactions', icon: FaRedo, label: 'Recurrentes' },
    { path: '/incomes', icon: FaPlusCircle, label: 'Ingresos' },
    { path: '/expenses', icon: FaMinusCircle, label: 'Gastos' },
    { path: '/categories', icon: FaFolderOpen, label: 'Categorías' }
  ];

  // Grupo 2: Análisis y planificación
  const analysisMenuItems = [
    { path: '/insights', icon: FaBrain, label: 'IA Financiero', hasSparkles: true, feature: 'AI_INSIGHTS' },
    { path: '/budgets', icon: FaChartPie, label: 'Presupuestos', subtitle: 'Controla tus límites', feature: 'BUDGETS' },
    { path: '/savings-goals', icon: FaBullseye, label: 'Metas de Ahorro', subtitle: 'Objetivos financieros', feature: 'SAVINGS_GOALS' },
    { path: '/achievements', icon: FaTrophy, label: 'Logros', subtitle: 'Progreso y gamificación' },
    { path: '/reports', icon: FaFileAlt, label: 'Reportes' }
  ];

  // Grupo 3: Configuración
  const settingsMenuItems = [
    { path: '/settings', icon: FaCog, label: 'Configuración' },
  ];

  const toggleSidebar = () => {
    setIsOpen(!isOpen);
  };

  const isActive = (path) => location.pathname === path;

  const renderMenuItem = (item) => {
    const Icon = item.icon;
    const active = isActive(item.path);

    const menuItem = isDesktopCollapsed ? (
      // Collapsed: icon only, centered
      <Link
        key={item.path}
        to={item.path}
        title={item.label}
        onClick={() => setIsOpen(false)}
        className={`
          group flex items-center justify-center p-3 rounded-xl transition-all duration-200 hover:bg-gray-50 dark:hover:bg-gray-700
          ${active
            ? 'bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 text-blue-700 dark:text-blue-400'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
          }
        `}
      >
        <Icon className={`w-5 h-5 ${
          active ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300'
        }`} />
      </Link>
    ) : (
      // Expanded: full item with label
      <Link
        key={item.path}
        to={item.path}
        onClick={() => setIsOpen(false)}
        className={`
          group flex flex-col px-4 py-3 rounded-xl transition-all duration-200 hover:bg-gray-50 dark:hover:bg-gray-700
          ${active
            ? 'bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-700 text-blue-700 dark:text-blue-400'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
          }
        `}
      >
        <div className="flex items-center space-x-3">
          <Icon className={`w-5 h-5 ${
            active ? 'text-blue-600 dark:text-blue-400' : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300'
          }`} />
          <span className="font-medium">{item.label}</span>
          {item.hasSparkles && (
            <FaStar className={`w-3 h-3 ${
              active ? 'text-blue-500 dark:text-blue-400' : 'text-gray-300 dark:text-gray-600'
            }`} />
          )}
        </div>
        {item.subtitle && (
          <span className={`text-xs ml-8 mt-1 ${
            active ? 'text-blue-600 dark:text-blue-400' : 'text-gray-500 dark:text-gray-400'
          }`}>
            {item.subtitle}
          </span>
        )}
      </Link>
    );

    // Si la item tiene una feature, envolverla con FeatureProgressIndicator
    if (item.feature) {
      return (
        <FeatureProgressIndicator key={item.path} feature={item.feature}>
          {menuItem}
        </FeatureProgressIndicator>
      );
    }

    return menuItem;
  };

  const renderSeparator = () => (
    <div className="my-4">
      <div className="h-px bg-gray-200 dark:bg-gray-700"></div>
    </div>
  );

  return (
    <>
      {/* Mobile menu button */}
      <button
        onClick={toggleSidebar}
        className="lg:hidden fixed top-3 left-3 z-50 p-2 bg-white/90 backdrop-blur-sm dark:bg-gray-800/90 rounded-lg shadow-lg dark:shadow-gray-900/30 border border-gray-200 dark:border-gray-700 text-gray-900 dark:text-gray-100"
        aria-label={isOpen ? 'Cerrar menú' : 'Abrir menú'}
      >
        {isOpen ? <FaTimes className="w-5 h-5" /> : <FaBars className="w-5 h-5" />}
      </button>

      {/* Overlay for mobile */}
      {isOpen && (
        <div
          className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-30"
          onClick={toggleSidebar}
        />
      )}

      {/* Sidebar */}
      <div className={`
        fixed lg:fixed inset-y-0 left-0 z-40 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transform transition-all duration-300 ease-in-out
        ${isDesktopCollapsed ? 'lg:w-16' : 'lg:w-52'}
        w-52
        ${isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
      `}>
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className={`flex items-center border-b border-gray-200 dark:border-gray-700 transition-all duration-300 ${isDesktopCollapsed ? 'p-2 flex-col gap-2' : 'px-4 py-3 justify-between'}`}>
            {!isDesktopCollapsed && <Brand size="sm" showTagline={false} />}
            {isDesktopCollapsed && (
              <div className="w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-white text-xs font-bold">
                N
              </div>
            )}
            {/* Desktop collapse toggle - top of sidebar */}
            <button
              onClick={onDesktopToggle}
              title={isDesktopCollapsed ? 'Expandir menú' : 'Colapsar menú'}
              className="hidden lg:flex p-1.5 rounded-lg text-gray-400 dark:text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-600 dark:hover:text-gray-300 transition-colors flex-shrink-0"
            >
              {isDesktopCollapsed
                ? <FaChevronRight className="w-3.5 h-3.5" />
                : <FaChevronLeft className="w-3.5 h-3.5" />
              }
            </button>
          </div>

          {/* Navigation */}
          <nav className={`flex-1 overflow-y-auto transition-all duration-300 ${isDesktopCollapsed ? 'p-2' : 'p-4'}`}>
            {/* Grupo 1: Transacciones principales */}
            <div className="space-y-1">
              {mainMenuItems.map(renderMenuItem)}
            </div>

            {renderSeparator()}

            {/* Grupo 2: Análisis y planificación */}
            <div className="space-y-1">
              {analysisMenuItems.map(renderMenuItem)}
            </div>

            {renderSeparator()}

            {/* Grupo 3: Configuración */}
            <div className="space-y-1">
              {settingsMenuItems.map(renderMenuItem)}
            </div>
          </nav>

        </div>
      </div>
    </>
  );
};

export default Sidebar;
