-- =====================================================
-- USER SERVICE - ESTRUCTURA DE BASE DE DATOS
-- Script para microservicio de gestión de usuarios
-- =====================================================

-- Eliminar tablas existentes si existen
DROP TABLE IF EXISTS user_two_fa CASCADE;
DROP TABLE IF EXISTS user_notification_settings CASCADE;
DROP TABLE IF EXISTS user_preferences CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- =====================================================
-- TABLA PRINCIPAL DE USUARIOS
-- =====================================================

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Hash bcrypt
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- TABLA DE PREFERENCIAS DE USUARIO
-- =====================================================

CREATE TABLE IF NOT EXISTS user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'ARS',
    language VARCHAR(2) DEFAULT 'es',
    theme VARCHAR(10) DEFAULT 'dark',
    date_format VARCHAR(20) DEFAULT 'DD/MM/YYYY',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

-- =====================================================
-- TABLA DE CONFIGURACIÓN DE NOTIFICACIONES
-- =====================================================

CREATE TABLE IF NOT EXISTS user_notification_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT FALSE,
    weekly_reports BOOLEAN DEFAULT TRUE,
    expense_alerts BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

-- =====================================================
-- TABLA DE AUTENTICACIÓN DE DOS FACTORES (2FA)
-- =====================================================

CREATE TABLE IF NOT EXISTS user_two_fa (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    secret VARCHAR(255) NOT NULL,
    enabled BOOLEAN DEFAULT FALSE,
    backup_codes TEXT[], -- Array de códigos de respaldo
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

-- =====================================================
-- ÍNDICES PARA MEJORAR EL RENDIMIENTO
-- =====================================================

-- Índices para usuarios
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);

-- Índices para preferencias
CREATE INDEX IF NOT EXISTS idx_preferences_user_id ON user_preferences(user_id);

-- Índices para notificaciones
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON user_notification_settings(user_id);

-- Índices para 2FA
CREATE INDEX IF NOT EXISTS idx_two_fa_user_id ON user_two_fa(user_id);
CREATE INDEX IF NOT EXISTS idx_two_fa_enabled ON user_two_fa(enabled);

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
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_preferences_updated_at 
    BEFORE UPDATE ON user_preferences 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_notification_settings_updated_at 
    BEFORE UPDATE ON user_notification_settings 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_two_fa_updated_at 
    BEFORE UPDATE ON user_two_fa 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- DATOS DE CONFIGURACIÓN INICIAL
-- =====================================================

-- Insertar preferencias por defecto para usuarios existentes
INSERT INTO user_preferences (user_id, currency, language, theme, date_format)
SELECT id, 'ARS', 'es', 'dark', 'DD/MM/YYYY' 
FROM users 
WHERE id NOT IN (SELECT user_id FROM user_preferences);

-- Insertar configuración de notificaciones por defecto para usuarios existentes
INSERT INTO user_notification_settings (user_id, email_notifications, push_notifications, weekly_reports, expense_alerts)
SELECT id, TRUE, FALSE, TRUE, TRUE 
FROM users 
WHERE id NOT IN (SELECT user_id FROM user_notification_settings);

COMMIT; 