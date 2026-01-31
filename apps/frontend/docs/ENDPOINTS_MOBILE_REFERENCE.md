# 📱 Financial Resume Engine - Guía Completa para Desarrollo Mobile

## 🎯 Información General

**Financial Resume Engine** es un ecosistema financiero completo con arquitectura de microservicios, gamificación nativa e IA integrada. Este documento contiene toda la información necesaria para que una IA externa pueda desarrollar la versión mobile sin acceso a los repositorios.

### 🏗️ Arquitectura del Sistema

**Microservicios Desplegados:**
- **Financial Resume Engine** (Puerto 8080) - Servicio principal + API Gateway
- **Gamification Service** (Puerto 8081) - Microservicio independiente de gamificación
- **AI Service** (Puerto 8082) - Servicios especializados de IA

**URLs de Producción:**
- **Desarrollo**: `http://localhost:8080/api/v1`
- **Render**: `https://financial-resume-engine.onrender.com/api/v1`

---

## 🔐 Sistema de Autenticación

### Configuración JWT
```javascript
// Headers requeridos para todas las peticiones autenticadas
{
  "Authorization": "Bearer {jwt_token}",
  "Content-Type": "application/json",
  "X-Caller-ID": "{user_id}" // Opcional, para compatibilidad
}
```

### Estructura del Usuario
```typescript
interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  created_at: string;
  // Token se almacena en localStorage como 'auth_token'
  // Datos de usuario en localStorage como 'auth_user'
  // Expiración en localStorage como 'auth_expires_at'
}
```

### Manejo de Sesión
- **Token Storage**: `localStorage.getItem('auth_token')`
- **User Data**: `JSON.parse(localStorage.getItem('auth_user'))`
- **Auto-refresh**: Se maneja automáticamente en interceptores
- **Logout**: Limpiar localStorage y redirigir a login

---

## 📄 Mapeo de Endpoints por Página

### 1. **LOGIN/REGISTER** (`/login`, `/register`)

**Endpoints Consumidos:**
```javascript
// Autenticación
POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/logout
POST /api/v1/auth/refresh
```

**Request/Response:**
```typescript
// Login Request
{
  email: string;
  password: string;
}

// Login/Register Response
{
  success: boolean;
  data: {
    token: string;
    expires_at: number;
    user: User;
  }
}
```

### 2. **DASHBOARD** (`/dashboard`)

**Endpoints Consumidos:**
```javascript
// Dashboard principal
GET /api/v1/dashboard?year={year}&month={month}

// Analytics optimizados (endpoints paralelos)
GET /api/v1/analytics/expenses?year={year}&month={month}&sort=date&order=desc&limit=50
GET /api/v1/analytics/incomes?year={year}&month={month}&sort=date&order=desc&limit=50
GET /api/v1/analytics/categories?year={year}&month={month}

// Funcionalidades avanzadas
GET /api/v1/budgets/dashboard?year={year}&month={month}
GET /api/v1/savings-goals/dashboard
GET /api/v1/recurring-transactions/dashboard

// Categorías para dropdowns
GET /api/v1/categories
```

**Estructura de Response Dashboard:**
```typescript
interface DashboardData {
  totalIncome: number;
  totalExpenses: number;
  balance: number;
  expenses: Transaction[];
  incomes: Transaction[];
  categories: Category[];
  dashboardMetrics: {
    monthly_savings_rate: number;
    expense_growth_rate: number;
    top_category: string;
  };
  categoriesAnalytics: CategoryAnalytics[];
}
```

### 3. **EXPENSES** (`/expenses`)

**Endpoints Consumidos:**
```javascript
// CRUD Gastos
GET /api/v1/expenses
POST /api/v1/expenses
PATCH /api/v1/expenses/{id}
DELETE /api/v1/expenses/{id}

// Categorías para formularios
GET /api/v1/categories

// Gastos no pagados
GET /api/v1/expenses/unpaid
```

**Estructura Expense:**
```typescript
interface Expense {
  id: number;
  description: string;
  amount: number;
  category_id?: number;
  due_date?: string;
  paid: boolean;
  amount_paid?: number;
  pending_amount?: number;
  created_at: string;
  user_id: number;
}
```

### 4. **INCOMES** (`/incomes`)

**Endpoints Consumidos:**
```javascript
// CRUD Ingresos
GET /api/v1/incomes
POST /api/v1/incomes
PATCH /api/v1/incomes/{id}
DELETE /api/v1/incomes/{id}

// Categorías para formularios
GET /api/v1/categories
```

**Estructura Income:**
```typescript
interface Income {
  id: number;
  description: string;
  amount: number;
  category_id?: number;
  created_at: string;
  user_id: number;
}
```

### 5. **CATEGORIES** (`/categories`)

**Endpoints Consumidos:**
```javascript
// CRUD Categorías
GET /api/v1/categories
POST /api/v1/categories
PATCH /api/v1/categories/{id}
DELETE /api/v1/categories/{id}
```

**Estructura Category:**
```typescript
interface Category {
  id: number;
  name: string;
  created_at: string;
  user_id: number;
}
```

### 6. **BUDGETS** (`/budgets`)

**Endpoints Consumidos:**
```javascript
// CRUD Presupuestos
GET /api/v1/budgets?year={year}&month={month}&category_id={id}&status={status}
POST /api/v1/budgets
PUT /api/v1/budgets/{id}
DELETE /api/v1/budgets/{id}

// Dashboard y estado
GET /api/v1/budgets/dashboard?year={year}&month={month}
GET /api/v1/budgets/status?year={year}&month={month}

// Categorías
GET /api/v1/categories
```

**Estructura Budget:**
```typescript
interface Budget {
  id: number;
  name: string;
  amount: number;
  category_id?: number;
  period: 'weekly' | 'monthly' | 'yearly';
  alert_threshold: number;
  start_date?: string;
  end_date?: string;
  created_at: string;
}
```

### 7. **SAVINGS GOALS** (`/savings-goals`)

**Endpoints Consumidos:**
```javascript
// CRUD Metas de Ahorro
GET /api/v1/savings-goals?status={status}&category={category}&priority={priority}
POST /api/v1/savings-goals
PUT /api/v1/savings-goals/{id}
DELETE /api/v1/savings-goals/{id}

// Transacciones de ahorro
POST /api/v1/savings-goals/{id}/deposit
POST /api/v1/savings-goals/{id}/withdraw

// Dashboard
GET /api/v1/savings-goals/dashboard
```

**Estructura SavingsGoal:**
```typescript
interface SavingsGoal {
  id: number;
  name: string;
  description?: string;
  target_amount: number;
  current_amount: number;
  target_date?: string;
  category: 'emergency' | 'vacation' | 'investment' | 'purchase' | 'education';
  priority: 'low' | 'medium' | 'high';
  auto_save_amount?: number;
  auto_save_frequency?: 'weekly' | 'monthly';
  created_at: string;
}
```

### 8. **RECURRING TRANSACTIONS** (`/recurring-transactions`)

**Endpoints Consumidos:**
```javascript
// CRUD Transacciones Recurrentes
GET /api/v1/recurring-transactions?type={type}&frequency={frequency}&status={status}
POST /api/v1/recurring-transactions
PUT /api/v1/recurring-transactions/{id}
DELETE /api/v1/recurring-transactions/{id}

// Control de ejecución
POST /api/v1/recurring-transactions/{id}/pause
POST /api/v1/recurring-transactions/{id}/resume
POST /api/v1/recurring-transactions/{id}/execute

// Operaciones masivas
POST /api/v1/recurring-transactions/batch/process
POST /api/v1/recurring-transactions/batch/notify

// Análisis y proyección
GET /api/v1/recurring-transactions/dashboard
GET /api/v1/recurring-transactions/projection?months={months}

// Categorías
GET /api/v1/categories
```

**Estructura RecurringTransaction:**
```typescript
interface RecurringTransaction {
  id: number;
  description: string;
  amount: number;
  type: 'income' | 'expense';
  frequency: 'daily' | 'weekly' | 'monthly' | 'yearly';
  category_id?: number;
  next_date: string;
  end_date?: string;
  day_of_month?: number;
  day_of_week?: number;
  is_active: boolean;
  created_at: string;
}
```

### 9. **ACHIEVEMENTS** (`/achievements`)

**Endpoints Consumidos:**
```javascript
// Gamificación (via proxy a microservicio)
GET /api/v1/gamification/profile
GET /api/v1/gamification/achievements
GET /api/v1/gamification/stats
GET /api/v1/gamification/features

// Endpoints públicos
GET /api/v1/gamification/action-types
GET /api/v1/gamification/levels

// Challenges
GET /api/v1/gamification/challenges/daily
GET /api/v1/gamification/challenges/weekly
POST /api/v1/gamification/challenges/progress
```

**Estructuras Gamification:**
```typescript
interface UserProfile {
  user_id: string;
  total_xp: number;
  current_level: number;
  level_name: string;
  next_level_xp: number;
  current_level_xp: number;
}

interface Achievement {
  id: number;
  name: string;
  description: string;
  xp_reward: number;
  is_completed: boolean;
  completed_at?: string;
  progress?: number;
  target?: number;
}
```

### 10. **AI INSIGHTS** (`/ai-insights`)

**Endpoints Consumidos:**
```javascript
// IA Especializada (via proxy a microservicio)
GET /api/v1/ai/insights?year={year}&month={month}
POST /api/v1/ai/can-i-buy
POST /api/v1/ai/credit-improvement-plan?year={year}&month={month}
GET /api/v1/insights/financial-health

// Marcar insights como entendidos
POST /api/v1/insights/mark-understood
```

**Estructura AI Insights:**
```typescript
interface AIInsight {
  id: string;
  title: string;
  description: string;
  type: 'suggestion' | 'warning' | 'achievement' | 'tip';
  priority: 'low' | 'medium' | 'high';
  category: string;
  generated_at: string;
  is_understood: boolean;
}

interface HealthScore {
  score: number;
  level: 'poor' | 'fair' | 'good' | 'excellent';
  message: string;
  insights: AIInsight[];
}
```

### 11. **REPORTS** (`/reports`)

**Endpoints Consumidos:**
```javascript
// Generación de reportes
GET /api/v1/reports/generate?start_date={date}&end_date={date}
GET /api/v1/reports/download?start_date={date}&end_date={date}&format=pdf
```

**Estructura Report:**
```typescript
interface Report {
  total_income: number;
  total_expenses: number;
  transactions: Transaction[];
  category_summary: CategorySummary[];
  period: {
    start_date: string;
    end_date: string;
  };
}
```

---

## 🎮 Sistema de Gamificación

### Auto-triggers Implementados
La gamificación se activa automáticamente cuando el usuario:
- ✅ Crea/actualiza/elimina gastos e ingresos
- ✅ Crea/asigna categorías
- ✅ Ve el dashboard y analytics
- ✅ Usa funciones de IA
- ✅ Realiza login diario

### Registro Manual de Acciones
```javascript
// Endpoint para registrar acciones manuales
POST /api/v1/gamification/actions
{
  action_type: string; // 'view_dashboard', 'create_expense', etc.
  entity_type: string; // 'expense', 'income', 'category', etc.
  entity_id: string;   // ID de la entidad
  description: string; // Descripción opcional
}
```

### Feature Gates (Funcionalidades Desbloqueables)
```typescript
const FEATURE_GATES = {
  SAVINGS_GOALS: { requiredLevel: 3, xpThreshold: 200 },
  BUDGETS: { requiredLevel: 5, xpThreshold: 700 },
  AI_INSIGHTS: { requiredLevel: 7, xpThreshold: 1800 }
};
```

---

## 🛠️ Configuración Técnica

### Detección de Ambiente
```javascript
const getApiBaseUrl = () => {
  const hostname = window.location.hostname;
  
  if (hostname.includes('onrender.com')) {
    return 'https://financial-resume-engine.onrender.com/api/v1';
  } else {
    return 'http://localhost:8080/api/v1'; // Development
  }
};
```

### Interceptores de Axios
```javascript
// Request interceptor - Agregar token automáticamente
api.interceptors.request.use(async (config) => {
  const token = localStorage.getItem('auth_token');
  const user = JSON.parse(localStorage.getItem('auth_user') || '{}');
  
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  
  if (user?.id) {
    config.headers['X-Caller-ID'] = user.id.toString();
  }
  
  return config;
});

// Response interceptor - Manejo de errores globales
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Limpiar datos de autenticación y redirigir
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
      localStorage.removeItem('auth_expires_at');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### Cache y Optimización
```javascript
// Cache de 5 minutos para datos del dashboard
const CACHE_DURATION = 5 * 60 * 1000;

// Rate limiting en desarrollo (1 request/segundo por endpoint)
const REQUEST_THROTTLE = 1000;

// Invalidación de cache después de mutaciones
const invalidateCache = (dataType) => {
  // Limpiar cache relacionado después de create/update/delete
  localStorage.removeItem(`cache_${dataType}`);
};
```

---

## 📊 Estructuras de Datos Comunes

### Transaction Base
```typescript
interface Transaction {
  id: number;
  description: string;
  amount: number;
  category_id?: number;
  user_id: number;
  created_at: string;
  updated_at?: string;
}
```

### Filtros de Período
```typescript
interface PeriodFilter {
  year?: string;    // "2024"
  month?: string;   // "01", "02", etc.
}

// Se convierte a query params: ?year=2024&month=01
```

### Response Wrapper
```typescript
interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  error?: string;
}
```

### Paginación
```typescript
interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    current_page: number;
    total_pages: number;
    total_items: number;
    items_per_page: number;
  };
}
```

---

## 🔒 Validaciones y Reglas de Negocio

### Validaciones de Frontend
```typescript
// Montos
const validateAmount = (amount: string) => {
  const num = parseFloat(amount);
  return {
    isValid: !isNaN(num) && num > 0 && num <= 999999999,
    error: !isValid ? 'Monto debe ser un número positivo' : null
  };
};

// Descripciones
const validateDescription = (description: string) => {
  return {
    isValid: description.trim().length >= 3 && description.length <= 255,
    error: !isValid ? 'Descripción debe tener entre 3 y 255 caracteres' : null
  };
};

// Email
const validateEmail = (email: string) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return {
    isValid: emailRegex.test(email),
    error: !isValid ? 'Email no válido' : null
  };
};
```

### Estados de Carga
```typescript
interface LoadingState {
  loading: boolean;
  error: string | null;
  data: any | null;
}

// Pattern usado en todas las páginas
const [state, setState] = useState<LoadingState>({
  loading: true,
  error: null,
  data: null
});
```

---

## 🎨 Patrones de UI/UX

### Temas (Dark/Light Mode)
```typescript
// Configuración de tema almacenada en localStorage
const theme = localStorage.getItem('financial-resume-theme') || 'light';

// Classes de Tailwind para modo oscuro
const darkModeClasses = 'dark:bg-gray-800 dark:text-white';
```

### Formateo de Moneda
```javascript
const formatCurrency = (amount) => {
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    minimumFractionDigits: 2
  }).format(amount);
};

// Ocultación de saldos
const formatAmount = (amount, hidden = false) => {
  return hidden ? '••••••' : formatCurrency(amount);
};
```

### Notificaciones Toast
```javascript
import toast from 'react-hot-toast';

// Éxito
toast.success('Operación exitosa');

// Error
toast.error('Error en la operación');

// Información
toast.info('Información importante');
```

---

## 📱 Consideraciones Específicas para Mobile

### Responsive Design
- **Breakpoints**: `sm:640px`, `md:768px`, `lg:1024px`, `xl:1280px`
- **Mobile-first**: Diseño base para móvil, luego adaptación desktop
- **Touch targets**: Mínimo 44px para botones tocables

### Navegación Mobile
```typescript
// Stack principal de navegación
const MobileStack = {
  Dashboard: '/dashboard',
  Expenses: '/expenses', 
  Incomes: '/incomes',
  Categories: '/categories',
  Budgets: '/budgets',
  SavingsGoals: '/savings-goals',
  RecurringTransactions: '/recurring-transactions',
  Reports: '/reports',
  Achievements: '/achievements',
  AIInsights: '/ai-insights'
};

// Bottom Tab Navigation recomendada
const BottomTabs = [
  { name: 'Dashboard', icon: 'home', route: '/dashboard' },
  { name: 'Gastos', icon: 'minus-circle', route: '/expenses' },
  { name: 'Ingresos', icon: 'plus-circle', route: '/incomes' },
  { name: 'Análisis', icon: 'chart-bar', route: '/reports' },
  { name: 'Perfil', icon: 'user', route: '/achievements' }
];
```

### Gestión de Estado Mobile
```typescript
// Context providers necesarios
const AppProviders = ({ children }) => (
  <AuthProvider>
    <ThemeProvider>
      <PeriodProvider>
        <GamificationProvider>
          {children}
        </GamificationProvider>
      </PeriodProvider>
    </ThemeProvider>
  </AuthProvider>
);
```

### Offline Support (Recomendado)
```typescript
// Cache local para modo offline
interface CacheEntry {
  data: any;
  timestamp: number;
  expiry: number;
}

const setCache = (key: string, data: any, duration = 300000) => {
  const entry: CacheEntry = {
    data,
    timestamp: Date.now(),
    expiry: Date.now() + duration
  };
  localStorage.setItem(`mobile_cache_${key}`, JSON.stringify(entry));
};

const getCache = (key: string) => {
  const cached = localStorage.getItem(`mobile_cache_${key}`);
  if (!cached) return null;
  
  const entry: CacheEntry = JSON.parse(cached);
  if (Date.now() > entry.expiry) {
    localStorage.removeItem(`mobile_cache_${key}`);
    return null;
  }
  
  return entry.data;
};
```

---

## 🚀 Checklist de Implementación Mobile

### ✅ Funcionalidades Core (MVP)
- [ ] Sistema de autenticación JWT
- [ ] Dashboard con métricas principales
- [ ] CRUD de gastos e ingresos
- [ ] Gestión de categorías
- [ ] Filtros por período (año/mes)
- [ ] Modo oscuro/claro

### ✅ Funcionalidades Avanzadas
- [ ] Presupuestos con alertas
- [ ] Metas de ahorro con progreso
- [ ] Transacciones recurrentes
- [ ] Sistema de gamificación completo
- [ ] Insights de IA
- [ ] Reportes y analytics

### ✅ Características Mobile
- [ ] Navegación con tabs inferiores
- [ ] Gestos touch optimizados
- [ ] Cache offline básico
- [ ] Notificaciones push (opcional)
- [ ] Biometric authentication (opcional)

### ✅ Testing y Calidad
- [ ] Tests unitarios para servicios API
- [ ] Tests de integración con backend
- [ ] Validación de todos los formularios
- [ ] Manejo de errores de red
- [ ] Loading states en todas las operaciones

---

## 📋 Recursos Adicionales

### Variables de Entorno Necesarias
```bash
# Para desarrollo
REACT_APP_API_URL=http://localhost:8080/api/v1

# Para producción
REACT_APP_API_URL=https://financial-resume-engine.onrender.com/api/v1
```

### Dependencias Clave
```json
{
  "axios": "^1.x.x",
  "react-hot-toast": "^2.x.x",
  "@react-navigation/native": "^6.x.x", // Para React Native
  "@react-navigation/bottom-tabs": "^6.x.x"
}
```

### Constantes Importantes
```typescript
const API_TIMEOUT = 10000; // 10 segundos
const CACHE_DURATION = 300000; // 5 minutos
const JWT_REFRESH_THRESHOLD = 300; // 5 minutos antes de expirar
const MAX_RETRY_ATTEMPTS = 3;
const RETRY_DELAY = 1000; // 1 segundo
```

---

**Última actualización**: Enero 2025
**Versión del sistema**: 1.0.0
**Total de endpoints**: 78+ endpoints funcionales 