import React, { useState, useEffect, useRef } from 'react';
import { FaCalculator } from 'react-icons/fa';

const buttons = [
  ['C', '±', '%', '÷'],
  ['7', '8', '9', '×'],
  ['4', '5', '6', '−'],
  ['1', '2', '3', '+'],
  ['0', '.', '='],
];

const SidebarCalculator = ({ collapsed }) => {
  const [open, setOpen] = useState(false);
  const [display, setDisplay] = useState('0');
  const [prev, setPrev] = useState(null);
  const [op, setOp] = useState(null);
  const [waitingForOperand, setWaitingForOperand] = useState(false);
  const panelRef = useRef(null);

  const handleNumber = (num) => {
    if (waitingForOperand) {
      setDisplay(num);
      setWaitingForOperand(false);
    } else {
      setDisplay(display === '0' ? num : display + num);
    }
  };

  const handleDot = () => {
    if (waitingForOperand) {
      setDisplay('0.');
      setWaitingForOperand(false);
      return;
    }
    if (!display.includes('.')) setDisplay(display + '.');
  };

  const handleOperator = (operator) => {
    const current = parseFloat(display);
    if (prev !== null && !waitingForOperand) {
      const result = calculate(prev, current, op);
      setDisplay(String(result));
      setPrev(result);
    } else {
      setPrev(current);
    }
    setOp(operator);
    setWaitingForOperand(true);
  };

  const calculate = (a, b, operator) => {
    switch (operator) {
      case '+': return parseFloat((a + b).toPrecision(10));
      case '−': return parseFloat((a - b).toPrecision(10));
      case '×': return parseFloat((a * b).toPrecision(10));
      case '÷': return b !== 0 ? parseFloat((a / b).toPrecision(10)) : 0;
      default: return b;
    }
  };

  const handleEquals = () => {
    if (prev === null || op === null) return;
    const current = parseFloat(display);
    const result = calculate(prev, current, op);
    setDisplay(String(result));
    setPrev(null);
    setOp(null);
    setWaitingForOperand(true);
  };

  const handleClear = () => {
    setDisplay('0');
    setPrev(null);
    setOp(null);
    setWaitingForOperand(false);
  };

  const handleToggleSign = () => {
    setDisplay(String(parseFloat(display) * -1));
  };

  const handlePercent = () => {
    setDisplay(String(parseFloat(display) / 100));
  };

  const handleButton = React.useCallback((btn) => {
    if (btn === 'C') return handleClear();
    if (btn === '±') return handleToggleSign();
    if (btn === '%') return handlePercent();
    if (btn === '.') return handleDot();
    if (btn === '=') return handleEquals();
    if (['+', '−', '×', '÷'].includes(btn)) return handleOperator(btn);
    handleNumber(btn);
  }, [display, prev, op, waitingForOperand]); // eslint-disable-line react-hooks/exhaustive-deps

  const isOperator = (btn) => ['+', '−', '×', '÷'].includes(btn);
  const isAction = (btn) => ['C', '±', '%'].includes(btn);

  useEffect(() => {
    if (!open) return;
    const keyMap = {
      '0': '0', '1': '1', '2': '2', '3': '3', '4': '4',
      '5': '5', '6': '6', '7': '7', '8': '8', '9': '9',
      '.': '.', ',': '.',
      '+': '+', '-': '−', '*': '×', '/': '÷',
      'Enter': '=', '=': '=',
      'Backspace': 'backspace',
      'Escape': 'C', 'Delete': 'C',
      '%': '%',
    };
    const onKeyDown = (e) => {
      const mapped = keyMap[e.key];
      if (!mapped) return;
      e.preventDefault();
      if (mapped === 'backspace') {
        setDisplay(prev => prev.length > 1 ? prev.slice(0, -1) : '0');
      } else {
        handleButton(mapped);
      }
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [open, handleButton]);

  // Collapsed: only icon toggle
  if (collapsed) {
    return (
      <button
        onClick={() => setOpen(!open)}
        title="Calculadora"
        className={`flex items-center justify-center p-3 rounded-xl transition-all duration-200 hover:bg-gray-50 dark:hover:bg-gray-700 ${
          open ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/30' : 'text-gray-400 dark:text-gray-500'
        }`}
      >
        <FaCalculator className="w-5 h-5" />
      </button>
    );
  }

  return (
    <div className="border-t border-gray-200 dark:border-gray-700">
      {/* Toggle button */}
      <button
        onClick={() => setOpen(!open)}
        className="w-full flex items-center space-x-3 px-4 py-3 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all duration-200"
      >
        <FaCalculator className="w-5 h-5 text-gray-400 dark:text-gray-500" />
        <span className="font-medium text-sm">Calculadora</span>
      </button>

      {/* Calculator panel */}
      {open && (
        <div className="px-3 pb-3">
          {/* Display */}
          <div className="bg-gray-900 dark:bg-gray-950 rounded-lg px-3 py-2 mb-2 text-right">
            <div className="text-xs text-gray-500 h-4 flex justify-between items-center">
              <span className="text-gray-600 text-xs">⌨</span>
              <span>{prev !== null ? `${prev} ${op}` : ''}</span>
            </div>
            <div className="text-white text-xl font-light truncate">
              {display.length > 10 ? parseFloat(display).toExponential(3) : display}
            </div>
          </div>

          {/* Buttons */}
          <div className="grid grid-cols-4 gap-1">
            {buttons.flat().map((btn, i) => (
              <button
                key={i}
                onClick={() => handleButton(btn)}
                className={`
                  ${btn === '0' ? 'col-span-2' : ''}
                  py-2 rounded-lg text-sm font-medium transition-all duration-100 active:scale-95
                  ${btn === '=' ? 'bg-blue-500 hover:bg-blue-600 text-white' :
                    isOperator(btn) ? 'bg-amber-400 hover:bg-amber-500 dark:bg-amber-500 dark:hover:bg-amber-600 text-white' :
                    isAction(btn) ? 'bg-gray-200 hover:bg-gray-300 dark:bg-gray-600 dark:hover:bg-gray-500 text-gray-800 dark:text-gray-100' :
                    'bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-100'
                  }
                `}
              >
                {btn}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default SidebarCalculator;
