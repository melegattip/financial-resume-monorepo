import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { PeriodProvider } from './contexts/PeriodContext';
import { ThemeProvider } from './contexts/ThemeContext';
import { GamificationProvider } from './contexts/GamificationContext';
import ProtectedRoute, { PublicOnlyRoute } from './components/ProtectedRoute';
import FeatureGuard from './components/FeatureGuard';

// Páginas principales
import Resumen from './pages/Dashboard';
import FinancialInsights from './pages/FinancialInsights';
import Expenses from './pages/Expenses';
import Incomes from './pages/Incomes';
import Categories from './pages/Categories';
import Reports from './pages/Reports';
import Settings from './pages/Settings';

// Nuevas páginas de funcionalidades avanzadas
import Budgets from './pages/Budgets';
import SavingsGoals from './pages/SavingsGoals';
import RecurringTransactions from './pages/RecurringTransactions';
import Achievements from './pages/Achievements';

// Páginas de autenticación
import Login from './pages/Login';
import Register from './pages/Register';

// Layout components
import Layout from './components/Layout/Layout';


// Estilos
import './index.css';

// Componente de rutas que usa los contexts (exportado para tests)
export function AppContent() {
  return (
    <>
      <Routes>
        {/* Rutas públicas (solo para usuarios NO autenticados) */}
        <Route 
          path="/login" 
          element={
            <PublicOnlyRoute>
              <Login />
            </PublicOnlyRoute>
          } 
        />
        <Route 
          path="/register" 
          element={
            <PublicOnlyRoute>
              <Register />
            </PublicOnlyRoute>
          } 
        />

        {/* Rutas protegidas (requieren autenticación) */}
        <Route 
          path="/" 
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route path="dashboard" element={<Resumen />} />
          <Route 
            path="insights" 
            element={
              <FeatureGuard feature="AI_INSIGHTS">
                <FinancialInsights />
              </FeatureGuard>
            } 
          />
          <Route path="expenses" element={<Expenses />} />
          <Route path="incomes" element={<Incomes />} />
          <Route path="categories" element={<Categories />} />
          <Route path="reports" element={<Reports />} />
          <Route 
            path="budgets" 
            element={
              <FeatureGuard feature="BUDGETS">
                <Budgets />
              </FeatureGuard>
            } 
          />
          <Route 
            path="savings-goals" 
            element={
              <FeatureGuard feature="SAVINGS_GOALS">
                <SavingsGoals />
              </FeatureGuard>
            } 
          />
          <Route path="recurring-transactions" element={<RecurringTransactions />} />
          <Route path="achievements" element={<Achievements />} />
          <Route path="settings" element={<Settings />} />
          <Route index element={<Navigate to="/dashboard" replace />} />
        </Route>



        {/* Ruta 404 - página no encontrada */}
        <Route 
          path="*" 
          element={
            <div className="min-h-screen flex items-center justify-center bg-fr-gray-50 dark:bg-gray-900">
              <div className="text-center">
                <h1 className="text-4xl font-bold text-fr-gray-900 dark:text-gray-100 mb-4">404</h1>
                <p className="text-fr-gray-600 dark:text-gray-400 mb-6">Página no encontrada</p>
                <a 
                  href="/" 
                  className="btn-primary"
                >
                  Volver al inicio
                </a>
              </div>
            </div>
          } 
        />
      </Routes>
      
      
    </>
  );
}

// Componente principal de la aplicación
function App() {
  return (
    <Router>
      <ThemeProvider>
        <AuthProvider>
          <PeriodProvider>
            <GamificationProvider>
              <AppContent />
            </GamificationProvider>
          </PeriodProvider>
        </AuthProvider>
      </ThemeProvider>
    </Router>
  );
}

export default App; 