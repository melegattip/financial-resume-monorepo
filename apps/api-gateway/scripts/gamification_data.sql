-- =====================================================
-- FINANCIAL RESUME ENGINE - DATOS DE GAMIFICACIÓN
-- Script con datos de gamificación para testing
-- =====================================================

-- =====================================================
-- DATOS DE GAMIFICACIÓN PARA USUARIOS DE PRUEBA  
-- =====================================================

-- USUARIO NIVEL 1 - PRINCIPIANTE (0 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_1', '1', 0, 1, 0, 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- USUARIO NIVEL 3 - SMART SAVER (200 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_3', '2', 200, 3, 15, 25, 2, 5, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- USUARIO NIVEL 5 - FINANCIAL PLANNER (700 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_5', '3', 700, 5, 45, 87, 4, 12, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- USUARIO NIVEL 10 - FINANCIAL MAGNATE (5500 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_10', '4', 5500, 10, 120, 275, 8, 30, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- USUARIO PRINCIPAL PABLO (NIVEL 7 - FINANCIAL EXPERT) (2000 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_pablo', '5', 2000, 7, 85, 150, 6, 20, CURRENT_TIMESTAMP - INTERVAL '15 minutes', CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- CATEGORÍAS PARA USUARIOS
-- =====================================================

-- Categorías de ejemplo para Pablo (ID 5)
INSERT INTO categories (id, user_id, name) VALUES 
('cat_aliment', '5', 'Alimentación'),
('cat_transp', '5', 'Transporte'),
('cat_entret', '5', 'Entretenimiento'),
('cat_salud', '5', 'Salud'),
('cat_educ', '5', 'Educación'),
('cat_invest', '5', 'Inversiones')
ON CONFLICT (user_id, name) DO NOTHING;

-- Categorías para usuarios de prueba
INSERT INTO categories (id, user_id, name) VALUES 
-- Usuario Nivel 1
('cat_nivel1_aliment', '1', 'Alimentación'),
('cat_nivel1_transp', '1', 'Transporte'),
-- Usuario Nivel 3  
('cat_nivel3_aliment', '2', 'Alimentación'),
('cat_nivel3_transp', '2', 'Transporte'),
('cat_nivel3_entret', '2', 'Entretenimiento'),
-- Usuario Nivel 5
('cat_nivel5_aliment', '3', 'Alimentación'),
('cat_nivel5_transp', '3', 'Transporte'),
('cat_nivel5_entret', '3', 'Entretenimiento'),
('cat_nivel5_salud', '3', 'Salud'),
('cat_nivel5_educ', '3', 'Educación'),
-- Usuario Nivel 10
('cat_nivel10_aliment', '4', 'Alimentación'),
('cat_nivel10_transp', '4', 'Transporte'),
('cat_nivel10_entret', '4', 'Entretenimiento'),
('cat_nivel10_salud', '4', 'Salud'),
('cat_nivel10_educ', '4', 'Educación'),
('cat_nivel10_invest', '4', 'Inversiones')
ON CONFLICT (user_id, name) DO NOTHING; 