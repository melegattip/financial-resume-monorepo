import React, { useEffect, useRef, useCallback } from 'react';

/**
 * Componente para manejo avanzado de foco y navegación por teclado
 */
export const FocusManager = ({ children, trapFocus = false, autoFocus = false, restoreFocus = false }) => {
  const containerRef = useRef(null);
  const previousActiveElement = useRef(null);

  // Elementos que pueden recibir foco
  const focusableSelector = [
    'button:not([disabled])',
    'input:not([disabled])',
    'select:not([disabled])',
    'textarea:not([disabled])',
    'a[href]',
    '[tabindex]:not([tabindex="-1"])',
    '[contenteditable]'
  ].join(', ');

  const getFocusableElements = useCallback(() => {
    if (!containerRef.current) return [];
    return Array.from(containerRef.current.querySelectorAll(focusableSelector));
  }, [focusableSelector]);

  const handleKeyDown = useCallback((event) => {
    if (!trapFocus) return;

    const focusableElements = getFocusableElements();
    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    if (event.key === 'Tab') {
      if (event.shiftKey) {
        // Shift + Tab - navegar hacia atrás
        if (document.activeElement === firstElement) {
          event.preventDefault();
          lastElement?.focus();
        }
      } else {
        // Tab - navegar hacia adelante
        if (document.activeElement === lastElement) {
          event.preventDefault();
          firstElement?.focus();
        }
      }
    }

    if (event.key === 'Escape') {
      // Escapar del contenedor con Escape
      previousActiveElement.current?.focus();
    }
  }, [trapFocus, getFocusableElements]);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    // Guardar elemento activo anterior
    if (restoreFocus) {
      previousActiveElement.current = document.activeElement;
    }

    // Auto-focus al primer elemento
    if (autoFocus) {
      const focusableElements = getFocusableElements();
      focusableElements[0]?.focus();
    }

    // Agregar event listener
    container.addEventListener('keydown', handleKeyDown);

    return () => {
      container.removeEventListener('keydown', handleKeyDown);
      
      // Restaurar foco anterior
      if (restoreFocus && previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, [autoFocus, restoreFocus, handleKeyDown, getFocusableElements]);

  return (
    <div ref={containerRef} className="focus-manager">
      {children}
    </div>
  );
};

/**
 * Hook para manejo de navegación por flechas en listas
 */
export const useArrowNavigation = (itemsCount, initialIndex = 0) => {
  const [activeIndex, setActiveIndex] = React.useState(initialIndex);

  const handleKeyDown = useCallback((event) => {
    switch (event.key) {
      case 'ArrowDown':
        event.preventDefault();
        setActiveIndex((prev) => (prev + 1) % itemsCount);
        break;
      case 'ArrowUp':
        event.preventDefault(); 
        setActiveIndex((prev) => (prev - 1 + itemsCount) % itemsCount);
        break;
      case 'Home':
        event.preventDefault();
        setActiveIndex(0);
        break;
      case 'End':
        event.preventDefault();
        setActiveIndex(itemsCount - 1);
        break;
    }
  }, [itemsCount]);

  return { activeIndex, setActiveIndex, handleKeyDown };
};

/**
 * Componente para anuncios de screen reader
 */
export const ScreenReaderAnnouncement = ({ message, priority = 'polite', children }) => {
  const [announcement, setAnnouncement] = React.useState('');

  React.useEffect(() => {
    if (message) {
      setAnnouncement(message);
      // Limpiar después de un breve delay para que el screen reader lo lea
      const timer = setTimeout(() => setAnnouncement(''), 1000);
      return () => clearTimeout(timer);
    }
  }, [message]);

  return (
    <>
      <div
        role="status"
        aria-live={priority}
        aria-atomic="true"
        className="sr-only"
      >
        {announcement}
      </div>
      {children}
    </>
  );
};

export default FocusManager; 