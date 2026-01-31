-- =====================================================
-- FINANCIAL RESUME ENGINE - ESTRUCTURA DE BASE DE DATOS
-- Script que crea solo tablas, índices y relaciones
-- Los usuarios se gestionan en el microservicio users-service
-- =====================================================

-- Eliminar tablas existentes si existen (en orden correcto por dependencias)
DROP TABLE IF EXISTS user_actions CASCADE;
DROP TABLE IF EXISTS achievements CASCADE;
DROP TABLE IF EXISTS user_gamification CASCADE;
DROP TABLE IF EXISTS savings_transactions CASCADE;
DROP TABLE IF EXISTS savings_goals CASCADE;
DROP TABLE IF EXISTS recurring_transaction_executions CASCADE;
DROP TABLE IF EXISTS recurring_transaction_notifications CASCADE;
DROP TABLE IF EXISTS recurring_transactions CASCADE;
DROP TABLE IF EXISTS budget_notifications CASCADE;
DROP TABLE IF EXISTS budgets CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS incomes CASCADE;
DROP TABLE IF EXISTS expenses CASCADE;
-- ✅ USERS TABLE REMOVIDA - Ahora manejada por users-service independiente

-- =====================================================
-- ✅ USERS TABLE REMOVIDA - MANEJADA POR USERS-SERVICE
-- =====================================================
-- Los usuarios ahora se manejan en el microservicio independiente
-- Base de datos: users_db (puerto 5434)
-- API: users-service (puerto 8083)
-- Ver: users-service/scripts/user_init.sql para estructura completa

-- =====================================================
-- TABLAS DE GAMIFICACIÓN
-- =====================================================

-- Tabla principal de gamificación por usuario
CREATE TABLE IF NOT EXISTS user_gamification (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE, -- Referencia externa al microservicio de usuarios
    total_xp INTEGER NOT NULL DEFAULT 0,
    current_level INTEGER NOT NULL DEFAULT 0,
    insights_viewed INTEGER NOT NULL DEFAULT 0,
    actions_completed INTEGER NOT NULL DEFAULT 0,
    achievements_count INTEGER NOT NULL DEFAULT 0,
    current_streak INTEGER NOT NULL DEFAULT 0,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de achievements/logros
CREATE TABLE IF NOT EXISTS achievements (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL, -- Referencia externa al microservicio de usuarios
    type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    progress INTEGER NOT NULL DEFAULT 0,
    target INTEGER NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    unlocked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);

-- Tabla de acciones del usuario para tracking de XP
CREATE TABLE IF NOT EXISTS user_actions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL, -- Referencia externa al microservicio de usuarios
    action_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    xp_earned INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);

-- =====================================================
-- TABLAS PRINCIPALES DEL SISTEMA FINANCIERO
-- =====================================================

-- Crear la tabla de categorías
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    icon VARCHAR(10) DEFAULT '📂',
    color VARCHAR(7) DEFAULT '#3B82F6',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Crear la tabla de ingresos
CREATE TABLE IF NOT EXISTS incomes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL, -- Referencia externa al microservicio de usuarios
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    category_id VARCHAR(50),
    source VARCHAR(50),
    percentage DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Crear la tabla de gastos
CREATE TABLE IF NOT EXISTS expenses (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    amount_paid DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    description TEXT,
    category_id VARCHAR(50),
    due_date DATE,
    paid BOOLEAN DEFAULT FALSE,
    percentage DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- =====================================================
-- SISTEMA DE PRESUPUESTOS
-- =====================================================

-- Tabla principal de presupuestos
CREATE TABLE IF NOT EXISTS budgets (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL, -- Referencia externa al microservicio de usuarios
    category_id VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    spent_amount DECIMAL(15,2) DEFAULT 0.00 CHECK (spent_amount >= 0),
    period VARCHAR(20) NOT NULL CHECK (period IN ('monthly', 'weekly', 'yearly')),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    alert_at DECIMAL(3,2) DEFAULT 0.80 CHECK (alert_at >= 0 AND alert_at <= 1),
    status VARCHAR(20) DEFAULT 'on_track' CHECK (status IN ('on_track', 'warning', 'exceeded')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de notificaciones de presupuesto
CREATE TABLE IF NOT EXISTS budget_notifications (
    id VARCHAR(50) PRIMARY KEY,
    budget_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL, -- Referencia externa al microservicio de usuarios
    notification_type VARCHAR(30) NOT NULL CHECK (notification_type IN ('alert', 'exceeded', 'reset')),
    message TEXT NOT NULL,
    threshold_percentage DECIMAL(5,2),
    exceeded_amount DECIMAL(15,2),
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- TRANSACCIONES RECURRENTES
-- =====================================================

-- Tabla principal de transacciones recurrentes
CREATE TABLE IF NOT EXISTS recurring_transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL, -- Referencia externa al microservicio de usuarios
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NOT NULL,
    category_id VARCHAR(36),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
    next_date DATE NOT NULL,
    last_executed TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    auto_create BOOLEAN DEFAULT TRUE,
    notify_before INTEGER DEFAULT 1 CHECK (notify_before >= 0),
    end_date DATE,
    execution_count INTEGER DEFAULT 0 CHECK (execution_count >= 0),
    max_executions INTEGER CHECK (max_executions > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de ejecuciones de transacciones recurrentes (audit trail)
CREATE TABLE IF NOT EXISTS recurring_transaction_executions (
    id VARCHAR(36) PRIMARY KEY,
    recurring_transaction_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL, -- Referencia externa al microservicio de usuarios
    amount DECIMAL(15,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN NOT NULL,
    created_transaction_id VARCHAR(36),
    error_message TEXT,
    execution_method VARCHAR(20) DEFAULT 'automatic'
);

-- Tabla de notificaciones de transacciones recurrentes
CREATE TABLE IF NOT EXISTS recurring_transaction_notifications (
    id VARCHAR(36) PRIMARY KEY,
    recurring_transaction_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL, -- Referencia externa al microservicio de usuarios
    notification_type VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivery_status VARCHAR(20) DEFAULT 'pending',
    delivery_method VARCHAR(20) DEFAULT 'system',
    error_message TEXT
);

-- =====================================================
-- METAS DE AHORRO
-- =====================================================

-- Tabla principal de metas de ahorro
CREATE TABLE IF NOT EXISTS savings_goals (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL, -- Referencia externa al microservicio de usuarios
    name VARCHAR(255) NOT NULL,
    description TEXT,
    target_amount DECIMAL(15,2) NOT NULL CHECK (target_amount > 0),
    current_amount DECIMAL(15,2) DEFAULT 0.00 CHECK (current_amount >= 0),
    category VARCHAR(20) NOT NULL CHECK (category IN ('vacation', 'emergency', 'house', 'car', 'education', 'retirement', 'investment', 'other')),
    priority VARCHAR(10) NOT NULL DEFAULT 'medium' CHECK (priority IN ('high', 'medium', 'low')),
    target_date DATE NOT NULL,
    status VARCHAR(15) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'achieved', 'paused', 'cancelled')),
    monthly_target DECIMAL(15,2) DEFAULT 0.00,
    weekly_target DECIMAL(15,2) DEFAULT 0.00,
    daily_target DECIMAL(15,2) DEFAULT 0.00,
    progress DECIMAL(5,4) DEFAULT 0.0000 CHECK (progress >= 0 AND progress <= 1),
    remaining_amount DECIMAL(15,2) DEFAULT 0.00,
    days_remaining INTEGER DEFAULT 0,
    is_auto_save BOOLEAN DEFAULT FALSE,
    auto_save_amount DECIMAL(15,2) DEFAULT 0.00,
    auto_save_frequency VARCHAR(10) CHECK (auto_save_frequency IN ('daily', 'weekly', 'monthly')),
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    achieved_at TIMESTAMP
);

-- Tabla de transacciones de ahorro (historial de depósitos/retiros)
CREATE TABLE IF NOT EXISTS savings_transactions (
    id VARCHAR(50) PRIMARY KEY,
    goal_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL, -- Referencia externa al microservicio de usuarios
    amount DECIMAL(15,2) NOT NULL CHECK (amount != 0),
    type VARCHAR(10) NOT NULL CHECK (type IN ('deposit', 'withdrawal')),
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (goal_id) REFERENCES savings_goals(id) ON DELETE CASCADE
);

-- =====================================================
-- ÍNDICES PARA MEJORAR EL RENDIMIENTO
-- =====================================================

-- ✅ Índices de usuarios removidos - manejados por users-service

-- Índices para gamificación
CREATE INDEX IF NOT EXISTS idx_user_gamification_user_id ON user_gamification(user_id);
CREATE INDEX IF NOT EXISTS idx_user_gamification_total_xp ON user_gamification(total_xp DESC);
CREATE INDEX IF NOT EXISTS idx_user_gamification_level ON user_gamification(current_level DESC);
CREATE INDEX IF NOT EXISTS idx_user_gamification_last_activity ON user_gamification(last_activity DESC);

CREATE INDEX IF NOT EXISTS idx_achievements_user_id ON achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_achievements_type ON achievements(type);
CREATE INDEX IF NOT EXISTS idx_achievements_completed ON achievements(completed);
CREATE INDEX IF NOT EXISTS idx_achievements_user_type ON achievements(user_id, type);

CREATE INDEX IF NOT EXISTS idx_user_actions_user_id ON user_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_action_type ON user_actions(action_type);
CREATE INDEX IF NOT EXISTS idx_user_actions_entity_type ON user_actions(entity_type);
CREATE INDEX IF NOT EXISTS idx_user_actions_created_at ON user_actions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_actions_user_date ON user_actions(user_id, created_at DESC);

-- Índices para transacciones
CREATE INDEX IF NOT EXISTS idx_incomes_user_id ON incomes(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_due_date ON expenses(due_date);
CREATE INDEX IF NOT EXISTS idx_expenses_paid ON expenses(paid);
CREATE INDEX IF NOT EXISTS idx_expenses_amount_paid ON expenses(amount_paid);
CREATE INDEX IF NOT EXISTS idx_expenses_paid_status ON expenses(paid, amount_paid);
CREATE INDEX IF NOT EXISTS idx_incomes_category_id ON incomes(category_id);
CREATE INDEX IF NOT EXISTS idx_expenses_category_id ON expenses(category_id);

-- Índices para presupuestos
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);
CREATE INDEX IF NOT EXISTS idx_budgets_category_id ON budgets(category_id);
CREATE INDEX IF NOT EXISTS idx_budgets_period ON budgets(period);
CREATE INDEX IF NOT EXISTS idx_budgets_status ON budgets(status);
CREATE INDEX IF NOT EXISTS idx_budgets_active ON budgets(is_active);
CREATE INDEX IF NOT EXISTS idx_budgets_period_dates ON budgets(period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_budgets_user_active ON budgets(user_id, is_active);

-- Índices para notificaciones
CREATE INDEX IF NOT EXISTS idx_budget_notifications_user_id ON budget_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_budget_notifications_budget_id ON budget_notifications(budget_id);
CREATE INDEX IF NOT EXISTS idx_budget_notifications_type ON budget_notifications(notification_type);
CREATE INDEX IF NOT EXISTS idx_budget_notifications_read ON budget_notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_budget_notifications_created ON budget_notifications(created_at);

-- Índices para transacciones recurrentes
CREATE INDEX IF NOT EXISTS idx_recurring_user_id ON recurring_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_recurring_next_date ON recurring_transactions(next_date);
CREATE INDEX IF NOT EXISTS idx_recurring_active ON recurring_transactions(is_active);
CREATE INDEX IF NOT EXISTS idx_recurring_type ON recurring_transactions(type);
CREATE INDEX IF NOT EXISTS idx_recurring_frequency ON recurring_transactions(frequency);
CREATE INDEX IF NOT EXISTS idx_recurring_category ON recurring_transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_recurring_user_active ON recurring_transactions(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_recurring_pending_execution ON recurring_transactions(is_active, next_date, auto_create);

-- Índices para ejecuciones
CREATE INDEX IF NOT EXISTS idx_execution_recurring_id ON recurring_transaction_executions(recurring_transaction_id);
CREATE INDEX IF NOT EXISTS idx_execution_user_id ON recurring_transaction_executions(user_id);
CREATE INDEX IF NOT EXISTS idx_execution_date ON recurring_transaction_executions(executed_at);
CREATE INDEX IF NOT EXISTS idx_execution_success ON recurring_transaction_executions(success);
CREATE INDEX IF NOT EXISTS idx_execution_created_transaction ON recurring_transaction_executions(created_transaction_id);

-- Índices para notificaciones recurrentes
CREATE INDEX IF NOT EXISTS idx_notification_recurring_id ON recurring_transaction_notifications(recurring_transaction_id);
CREATE INDEX IF NOT EXISTS idx_notification_user_id ON recurring_transaction_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_type ON recurring_transaction_notifications(notification_type);
CREATE INDEX IF NOT EXISTS idx_notification_status ON recurring_transaction_notifications(delivery_status);
CREATE INDEX IF NOT EXISTS idx_notification_date ON recurring_transaction_notifications(sent_at);

-- Índices para metas de ahorro
CREATE INDEX IF NOT EXISTS idx_savings_goals_user_id ON savings_goals(user_id);
CREATE INDEX IF NOT EXISTS idx_savings_goals_status ON savings_goals(status);
CREATE INDEX IF NOT EXISTS idx_savings_goals_category ON savings_goals(category);
CREATE INDEX IF NOT EXISTS idx_savings_goals_priority ON savings_goals(priority);
CREATE INDEX IF NOT EXISTS idx_savings_goals_target_date ON savings_goals(target_date);
CREATE INDEX IF NOT EXISTS idx_savings_goals_user_status ON savings_goals(user_id, status);
CREATE INDEX IF NOT EXISTS idx_savings_goals_user_category ON savings_goals(user_id, category);
CREATE INDEX IF NOT EXISTS idx_savings_goals_auto_save ON savings_goals(is_auto_save, auto_save_frequency);

-- Índices para transacciones de ahorro
CREATE INDEX IF NOT EXISTS idx_savings_transactions_goal_id ON savings_transactions(goal_id);
CREATE INDEX IF NOT EXISTS idx_savings_transactions_user_id ON savings_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_savings_transactions_type ON savings_transactions(type);
CREATE INDEX IF NOT EXISTS idx_savings_transactions_created ON savings_transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_savings_transactions_goal_created ON savings_transactions(goal_id, created_at);

-- =====================================================
-- FUNCIONES Y TRIGGERS
-- =====================================================

-- Función para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

-- Triggers para actualizar updated_at
-- ✅ Trigger de users removido - manejado por users-service

CREATE TRIGGER update_user_gamification_updated_at 
    BEFORE UPDATE ON user_gamification 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_achievements_updated_at 
    BEFORE UPDATE ON achievements 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_budgets_updated_at 
    BEFORE UPDATE ON budgets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recurring_transactions_updated_at 
    BEFORE UPDATE ON recurring_transactions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_savings_goals_updated_at 
    BEFORE UPDATE ON savings_goals 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- ESTRUCTURA COMPLETA LISTA
-- =====================================================

-- Mostrar mensaje de completado
COMMIT;
