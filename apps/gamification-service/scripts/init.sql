-- Financial Gamification Service Database Schema
-- Script para crear tablas en la base de datos existente
-- (La base de datos ya existe como financial_gamification)

-- Tabla principal de gamificación de usuarios
CREATE TABLE IF NOT EXISTS user_gamification (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    total_xp INTEGER DEFAULT 0,
    current_level INTEGER DEFAULT 0,
    insights_viewed INTEGER DEFAULT 0,
    actions_completed INTEGER DEFAULT 0,
    achievements_count INTEGER DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de logros/achievements
CREATE TABLE IF NOT EXISTS achievements (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    points INTEGER DEFAULT 0,
    progress INTEGER DEFAULT 0,
    target INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    unlocked_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);

-- Tabla de acciones de usuario para tracking de XP
CREATE TABLE IF NOT EXISTS user_actions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255),
    xp_earned INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);

-- Tabla de challenges disponibles
CREATE TABLE IF NOT EXISTS challenges (
    id VARCHAR(255) PRIMARY KEY,
    challenge_key VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    challenge_type VARCHAR(50) NOT NULL, -- 'daily', 'weekly', 'monthly'
    icon VARCHAR(20) DEFAULT '🎯',
    xp_reward INTEGER DEFAULT 0,
    requirement_type VARCHAR(100) NOT NULL, -- 'transaction_count', 'category_variety', 'view_combo', etc.
    requirement_target INTEGER NOT NULL,
    requirement_data JSONB, -- Para requisitos complejos
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de progreso de challenges por usuario
CREATE TABLE IF NOT EXISTS user_challenges (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    challenge_id VARCHAR(255) NOT NULL,
    challenge_date DATE NOT NULL, -- Para challenges diarios/semanales
    progress INTEGER DEFAULT 0,
    target INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE,
    FOREIGN KEY (challenge_id) REFERENCES challenges(id) ON DELETE CASCADE,
    UNIQUE(user_id, challenge_id, challenge_date)
);

-- Tabla de tracking de acciones para challenges
CREATE TABLE IF NOT EXISTS challenge_progress_tracking (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    challenge_date DATE NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    count INTEGER DEFAULT 1,
    unique_entities JSONB, -- Para tracking de entities únicas (categorías, etc.)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE,
    UNIQUE(user_id, challenge_date, action_type, entity_type)
);

-- Índices para optimizar performance
CREATE INDEX IF NOT EXISTS idx_user_gamification_user_id ON user_gamification(user_id);
CREATE INDEX IF NOT EXISTS idx_user_gamification_total_xp ON user_gamification(total_xp DESC);
CREATE INDEX IF NOT EXISTS idx_user_gamification_level ON user_gamification(current_level);

CREATE INDEX IF NOT EXISTS idx_achievements_user_id ON achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_achievements_type ON achievements(type);
CREATE INDEX IF NOT EXISTS idx_achievements_completed ON achievements(completed);

CREATE INDEX IF NOT EXISTS idx_user_actions_user_id ON user_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_type ON user_actions(action_type);
CREATE INDEX IF NOT EXISTS idx_user_actions_created_at ON user_actions(created_at DESC);

-- Índices para challenges
CREATE INDEX IF NOT EXISTS idx_challenges_type ON challenges(challenge_type);
CREATE INDEX IF NOT EXISTS idx_challenges_active ON challenges(active);
CREATE INDEX IF NOT EXISTS idx_challenges_key ON challenges(challenge_key);

CREATE INDEX IF NOT EXISTS idx_user_challenges_user_date ON user_challenges(user_id, challenge_date);
CREATE INDEX IF NOT EXISTS idx_user_challenges_completed ON user_challenges(completed);
CREATE INDEX IF NOT EXISTS idx_user_challenges_user_id ON user_challenges(user_id);

CREATE INDEX IF NOT EXISTS idx_challenge_tracking_user_date ON challenge_progress_tracking(user_id, challenge_date);
CREATE INDEX IF NOT EXISTS idx_challenge_tracking_action ON challenge_progress_tracking(action_type);

-- Función para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

-- Triggers para actualizar updated_at
CREATE TRIGGER update_user_gamification_updated_at 
    BEFORE UPDATE ON user_gamification 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_achievements_updated_at 
    BEFORE UPDATE ON achievements 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_challenges_updated_at 
    BEFORE UPDATE ON challenges 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_challenges_updated_at 
    BEFORE UPDATE ON user_challenges 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insertar challenges diarios por defecto
INSERT INTO challenges (id, challenge_key, name, description, challenge_type, icon, xp_reward, requirement_type, requirement_target, requirement_data, active) VALUES
('ch-daily-tx-master', 'transaction_master', 'Maestro de Transacciones', 'Registra 3 transacciones hoy', 'daily', '💰', 20, 'transaction_count', 3, '{"actions": ["create_expense", "create_income"]}', TRUE),
('ch-daily-cat-org', 'category_organizer', 'Organizador Expert', 'Usa 3 categorías diferentes', 'daily', '🏷️', 15, 'category_variety', 3, '{"actions": ["assign_category"]}', TRUE),
('ch-daily-analytics', 'analytics_explorer', 'Explorador de Datos', 'Revisa dashboard y analytics', 'daily', '📊', 10, 'view_combo', 2, '{"actions": ["view_dashboard", "view_analytics"]}', TRUE),
('ch-daily-streak', 'streak_keeper', 'Constancia', 'Mantén tu racha diaria', 'daily', '🔥', 5, 'daily_login', 1, '{"actions": ["daily_login"]}', TRUE)
ON CONFLICT (challenge_key) DO NOTHING;

-- Insertar challenges semanales por defecto
INSERT INTO challenges (id, challenge_key, name, description, challenge_type, icon, xp_reward, requirement_type, requirement_target, requirement_data, active) VALUES
('ch-weekly-tx-champ', 'transaction_champion', 'Campeón Financiero', 'Registra 15 transacciones esta semana', 'weekly', '🏆', 75, 'transaction_count', 15, '{"actions": ["create_expense", "create_income", "update_expense", "update_income"]}', TRUE),
('ch-weekly-cat-master', 'category_master', 'Maestro de Categorías', 'Usa al menos 5 categorías diferentes', 'weekly', '🎯', 50, 'category_variety', 5, '{"actions": ["assign_category"]}', TRUE),
('ch-weekly-engagement', 'engagement_hero', 'Héroe del Engagement', 'Inicia sesión 5 días esta semana', 'weekly', '⭐', 60, 'daily_login_count', 5, '{"actions": ["daily_login"]}', TRUE)
ON CONFLICT (challenge_key) DO NOTHING;

-- Insertar datos de ejemplo (opcional)
-- INSERT INTO user_gamification (id, user_id, total_xp, current_level) 
-- VALUES ('example-1', 'user-123', 150, 1)
-- ON CONFLICT (user_id) DO NOTHING;

COMMIT;

-- Mostrar tablas creadas
\dt 