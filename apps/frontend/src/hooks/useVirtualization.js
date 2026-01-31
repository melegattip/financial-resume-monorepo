import { useState, useEffect, useMemo } from 'react';

/**
 * Hook para virtualización de listas grandes
 * @param {Array} items - Array de items a virtualizar
 * @param {number} itemHeight - Altura de cada item
 * @param {number} containerHeight - Altura del contenedor
 * @param {number} overscan - Número de items extra a renderizar (buffer)
 * @returns {Object} Objeto con items visibles y props de scroll
 */
export function useVirtualization({
  items = [],
  itemHeight = 50,
  containerHeight = 400,
  overscan = 5
}) {
  const [scrollTop, setScrollTop] = useState(0);

  const visibleRange = useMemo(() => {
    const totalItems = items.length;
    const visibleCount = Math.ceil(containerHeight / itemHeight);
    
    const startIndex = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan);
    const endIndex = Math.min(
      totalItems - 1,
      startIndex + visibleCount + overscan * 2
    );

    return { startIndex, endIndex };
  }, [scrollTop, items.length, itemHeight, containerHeight, overscan]);

  const visibleItems = useMemo(() => {
    return items.slice(visibleRange.startIndex, visibleRange.endIndex + 1)
      .map((item, index) => ({
        item,
        index: visibleRange.startIndex + index,
        style: {
          position: 'absolute',
          top: (visibleRange.startIndex + index) * itemHeight,
          height: itemHeight,
          width: '100%'
        }
      }));
  }, [items, visibleRange, itemHeight]);

  const totalHeight = items.length * itemHeight;

  const handleScroll = (event) => {
    setScrollTop(event.target.scrollTop);
  };

  return {
    visibleItems,
    totalHeight,
    handleScroll,
    containerProps: {
      style: {
        height: containerHeight,
        overflow: 'auto',
        position: 'relative'
      },
      onScroll: handleScroll
    }
  };
}

export default useVirtualization; 