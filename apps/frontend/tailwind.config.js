/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  darkMode: 'class', // Habilitar dark mode con clase
  theme: {
    extend: {
      colors: {
        // Colores principales del sistema
        'fr-blue': {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        },
        'fr-primary': '#009ee3',
        'fr-secondary': '#00a650',
        'fr-accent': '#ff6900',
        'fr-gray': {
          50: '#f9fafb',
          100: '#f3f4f6',
          200: '#e5e7eb',
          300: '#d1d5db',
          400: '#9ca3af',
          500: '#6b7280',
          600: '#4b5563',
          700: '#374151',
          800: '#1f2937',
          900: '#111827',
        },
        'fr-success': '#00a650',
        'fr-warning': '#ff6900',
        'fr-error': '#e53e3e',
        
        // Colores base para componentes
        border: '#e5e7eb',
        input: '#f9fafb',
        ring: '#009ee3',
        background: '#ffffff',
        foreground: '#111827',
        primary: {
          DEFAULT: '#009ee3',
          foreground: '#ffffff',
        },
        secondary: {
          DEFAULT: '#f3f4f6',
          foreground: '#374151',
        },
        destructive: {
          DEFAULT: '#e53e3e',
          foreground: '#ffffff',
        },
        muted: {
          DEFAULT: '#f3f4f6',
          foreground: '#6b7280',
        },
        accent: {
          DEFAULT: '#f3f4f6',
          foreground: '#374151',
        },
        popover: {
          DEFAULT: '#ffffff',
          foreground: '#111827',
        },
        card: {
          DEFAULT: '#ffffff',
          foreground: '#111827',
        },
      },
      fontFamily: {
        'sans': ['Inter', 'Proxima Nova', 'system-ui', 'sans-serif'],
      },
      fontSize: {
        'xs': ['0.65rem', { lineHeight: '0.95rem' }],
        'sm': ['0.8rem', { lineHeight: '1.1rem' }],
        'base': ['0.9rem', { lineHeight: '1.35rem' }],
        'lg': ['1rem', { lineHeight: '1.5rem' }],
        'xl': ['1.125rem', { lineHeight: '1.65rem' }],
        '2xl': ['1.35rem', { lineHeight: '1.9rem' }],
        '3xl': ['1.7rem', { lineHeight: '2.25rem' }],
        '4xl': ['2rem', { lineHeight: '2.5rem' }],
        '5xl': ['2.5rem', { lineHeight: '3rem' }],
        '6xl': ['3rem', { lineHeight: '3.5rem' }],
      },
      borderRadius: {
        'fr': '8px',
        'fr-lg': '12px',
        lg: '8px',
        md: '6px',
        sm: '4px',
      },
      boxShadow: {
        'fr': '0 2px 4px 0 rgba(0, 0, 0, 0.1)',
        'fr-lg': '0 4px 12px 0 rgba(0, 0, 0, 0.15)',
        'fr-xl': '0 8px 24px 0 rgba(0, 0, 0, 0.15)',
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'pulse-soft': 'pulseSoft 2s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        pulseSoft: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.8' },
        },
      },
    },
  },
  plugins: [],
} 