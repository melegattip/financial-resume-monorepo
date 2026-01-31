-- =====================================================
-- SCRIPT DE SINCRONIZACIÓN DE DATOS DE USUARIOS
-- Sincroniza datos desde financial_resume a gamification_db
-- =====================================================

-- Primero, eliminar datos existentes en gamification_db para resincronizar
DELETE FROM user_actions;
DELETE FROM achievements;
DELETE FROM user_gamification;

-- Copiar datos de gamificación desde financial_resume
-- Nota: Ejecutar este INSERT desde la BD gamification_db conectándose a financial_resume via foreign data wrapper o manual

-- USUARIOS DE PRUEBA CON GAMIFICACIÓN CORRECTA
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 

-- USUARIO NIVEL 1 - PRINCIPIANTE (0 XP)
('gam_user_nivel_1', '1', 0, 1, 0, 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- USUARIO NIVEL 3 - SMART SAVER (200 XP)
('gam_user_nivel_3', '2', 200, 3, 15, 25, 2, 5, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP),

-- USUARIO NIVEL 5 - FINANCIAL PLANNER (700 XP)
('gam_user_nivel_5', '3', 700, 5, 45, 87, 4, 12, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP),

-- USUARIO NIVEL 10 - FINANCIAL MAGNATE (5500 XP)
('gam_user_nivel_10', '4', 5500, 10, 120, 275, 8, 30, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP),

-- USUARIO PRINCIPAL PABLO (NIVEL 7 - FINANCIAL EXPERT) (2000 XP)
('gam_user_pablo', '5', 2000, 7, 85, 150, 6, 20, CURRENT_TIMESTAMP - INTERVAL '15 minutes', CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP);

-- ✅ ACHIEVEMENTS BÁSICOS PARA TODOS LOS USUARIOS
INSERT INTO achievements (id, user_id, type, name, description, points, progress, target, completed, unlocked_at, created_at, updated_at) VALUES

-- ACHIEVEMENTS USUARIO NIVEL 1
('ach_nivel1_explorer', '1', 'first_steps', '🚀 Financial Explorer', 'Has dado tus primeros pasos en el mundo financiero', 50, 1, 1, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- ACHIEVEMENTS USUARIO NIVEL 3  
('ach_nivel3_saver', '2', 'smart_saver', '💰 Smart Saver', 'Has demostrado ser un ahorrador inteligente', 100, 200, 200, true, CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP),
('ach_nivel3_analyst', '2', 'data_analyst', '📊 Data Analyst', 'Has revisado insights financieros 15 veces', 75, 15, 25, false, null, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP),

-- ACHIEVEMENTS USUARIO NIVEL 5
('ach_nivel5_planner', '3', 'financial_planner', '📈 Financial Planner', 'Has creado y gestionado presupuestos como un experto', 200, 700, 700, true, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP),
('ach_nivel5_organizer', '3', 'category_master', '🏷️ Category Master', 'Has organizado tus gastos en categorías detalladas', 150, 45, 50, false, null, CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP),
('ach_nivel5_insights', '3', 'insight_seeker', '🔍 Insight Seeker', 'Has explorado análisis avanzados', 125, 45, 50, false, null, CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP),

-- ACHIEVEMENTS USUARIO NIVEL 10
('ach_nivel10_magnate', '4', 'financial_magnate', '👑 Financial Magnate', 'Has alcanzado el máximo nivel de maestría financiera', 1000, 5500, 5500, true, CURRENT_TIMESTAMP - INTERVAL '30 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP),
('ach_nivel10_master', '4', 'ai_master', '🤖 AI Master', 'Has utilizado la IA financiera como un verdadero experto', 500, 120, 100, true, CURRENT_TIMESTAMP - INTERVAL '35 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP),
('ach_nivel10_guru', '4', 'financial_guru', '💎 Financial Guru', 'Has completado más de 275 acciones financieras', 750, 275, 250, true, CURRENT_TIMESTAMP - INTERVAL '25 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP),
('ach_nivel10_streak', '4', 'streak_legend', '🔥 Streak Legend', 'Has mantenido una racha de 30 días consecutivos', 300, 30, 30, true, CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP),

-- ACHIEVEMENTS USUARIO PABLO
('ach_pablo_expert', '5', 'financial_expert', '🎓 Financial Expert', 'Has demostrado experiencia avanzada en finanzas', 500, 2000, 2000, true, CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP),
('ach_pablo_ai', '5', 'ai_partner', '🤖 AI Partner', 'Has establecido una alianza con la IA financiera', 300, 85, 75, true, CURRENT_TIMESTAMP - INTERVAL '50 days', CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP),
('ach_pablo_action', '5', 'action_taker', '⚡ Action Taker', 'Has completado 150 acciones financieras', 200, 150, 150, true, CURRENT_TIMESTAMP - INTERVAL '40 days', CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP);

-- ✅ ACCIONES DE EJEMPLO PARA TRACKING
INSERT INTO user_actions (id, user_id, action_type, entity_type, entity_id, xp_earned, description, created_at) VALUES

-- Acciones recientes del usuario nivel 10
('action_nivel10_1', '4', 'view_dashboard', 'dashboard', 'main_dashboard', 10, 'Usuario visualizó dashboard principal', CURRENT_TIMESTAMP - INTERVAL '30 minutes'),
('action_nivel10_2', '4', 'view_insights', 'insights', 'ai_financial_health', 25, 'Usuario revisó análisis de IA', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('action_nivel10_3', '4', 'create_budget', 'budget', 'budget_2025_01', 50, 'Usuario creó presupuesto mensual', CURRENT_TIMESTAMP - INTERVAL '2 days'),

-- Acciones recientes del usuario pablo
('action_pablo_1', '5', 'view_expenses', 'expenses', 'expense_list', 15, 'Usuario revisó lista de gastos', CURRENT_TIMESTAMP - INTERVAL '15 minutes'),
('action_pablo_2', '5', 'create_category', 'category', 'cat_invest', 20, 'Usuario creó categoría Inversiones', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('action_pablo_3', '5', 'use_ai_insight', 'insights', 'can_i_buy_analysis', 30, 'Usuario utilizó análisis ¿Puedo comprarlo?', CURRENT_TIMESTAMP - INTERVAL '3 days');

-- Mostrar estadísticas finales
SELECT 
    u.user_id,
    u.total_xp,
    u.current_level,
    u.achievements_count,
    COUNT(a.id) as total_achievements,
    COUNT(ac.id) as total_actions
FROM user_gamification u
LEFT JOIN achievements a ON u.user_id = a.user_id AND a.completed = true
LEFT JOIN user_actions ac ON u.user_id = ac.user_id
GROUP BY u.user_id, u.total_xp, u.current_level, u.achievements_count
ORDER BY u.current_level DESC; 