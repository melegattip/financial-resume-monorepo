import React, { createContext, useContext, useEffect, useState } from 'react';

const ThemeContext = createContext();

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme debe ser usado dentro de ThemeProvider');
  }
  return context;
};

export const ThemeProvider = ({ children }) => {
  const [theme, setTheme] = useState(() => {
    // Verificar localStorage primero
    try {
      const savedTheme = localStorage?.getItem('financial-resume-theme');
      if (savedTheme) {
        return savedTheme;
      }
    } catch (error) {
      // localStorage no disponible (tests)
    }
    
    // Si no hay tema guardado, verificar preferencia del sistema
    try {
      if (window?.matchMedia && window.matchMedia('(prefers-color-scheme: dark)')?.matches) {
        return 'dark';
      }
    } catch (error) {
      // matchMedia no disponible (tests)
    }
    
    return 'light';
  });

  useEffect(() => {
    try {
      const root = window?.document?.documentElement;
      
      if (root) {
        if (theme === 'dark') {
          root.classList.add('dark');
        } else {
          root.classList.remove('dark');
        }
      }
      
      // Guardar en localStorage
      localStorage?.setItem('financial-resume-theme', theme);
    } catch (error) {
      // Ignorar errores en entorno de tests
    }
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prevTheme => prevTheme === 'light' ? 'dark' : 'light');
  };

  const setLightTheme = () => setTheme('light');
  const setDarkTheme = () => setTheme('dark');
  const setSystemTheme = () => {
    try {
      const systemTheme = window?.matchMedia && window.matchMedia('(prefers-color-scheme: dark)')?.matches ? 'dark' : 'light';
      setTheme(systemTheme);
    } catch (error) {
      // Fallback a light si no se puede detectar
      setTheme('light');
    }
  };

  const value = {
    theme,
    toggleTheme,
    setLightTheme,
    setDarkTheme,
    setSystemTheme,
    isDark: theme === 'dark'
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
}; 