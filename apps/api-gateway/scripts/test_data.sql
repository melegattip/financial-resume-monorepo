-- =====================================================
-- FINANCIAL RESUME ENGINE - DATOS DE PRUEBA
-- Script con usuarios, transacciones y datos de testing
-- =====================================================
-- 
-- CREDENCIALES DE USUARIOS PARA TESTING:
-- 📧 nivel1@test.com     🔑 password123    🎮 Nivel 1 (0 XP)
-- 📧 nivel3@test.com     🔑 password123    🎮 Nivel 3 (200 XP) 
-- 📧 nivel5@test.com     🔑 password123    🎮 Nivel 5 (700 XP)
-- 📧 nivel10@test.com    🔑 password123    🎮 Nivel 10 (5500 XP)
-- 📧 pablo@niloft.com    🔑 password123    🎮 Usuario principal
-- =====================================================

-- =====================================================
-- USUARIOS DE PRUEBA CON CREDENCIALES REALES
-- =====================================================
-- Contraseña para todos: "password123"
-- Hashes generados con bcrypt.DefaultCost

-- USUARIO NIVEL 1 - PRINCIPIANTE (0 XP)
INSERT INTO users (id, email, password, first_name, last_name, is_active, is_verified, created_at, updated_at) VALUES 
(1, 'nivel1@test.com', '$2a$10$Z8bsd0n1VIb8OPZrztu90u1Q.G6sIOhWEJZSBy.SGf382MNCSRPJy', 'Usuario', 'Nivel1', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- USUARIO NIVEL 3 - SMART SAVER (200 XP)
INSERT INTO users (id, email, password, first_name, last_name, is_active, is_verified, created_at, updated_at) VALUES 
(2, 'nivel3@test.com', '$2a$10$RGmmn3yvBEqrslPi7GM4MuNriFZ07XFatUc3.pJn.AIpXZow197VG', 'Usuario', 'Nivel3', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- USUARIO NIVEL 5 - FINANCIAL PLANNER (700 XP)
INSERT INTO users (id, email, password, first_name, last_name, is_active, is_verified, created_at, updated_at) VALUES 
(3, 'nivel5@test.com', '$2a$10$vdNVzoJUrTHIeZtkb2g5l..1wWE4mZRtYhkZx3zyrFUQokgTjjEdi', 'Usuario', 'Nivel5', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- USUARIO NIVEL 10 - FINANCIAL MAGNATE (5500 XP)
INSERT INTO users (id, email, password, first_name, last_name, is_active, is_verified, created_at, updated_at) VALUES 
(4, 'nivel10@test.com', '$2a$10$y/GDFZFtwIiYu5j0YjobJO/Myf2AqER0JfUbto3lX9N4KmGPrLLdS', 'Usuario', 'Nivel10', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- USUARIO PRINCIPAL - PABLO MELEGATTI (para datos existentes)
INSERT INTO users (id, email, password, first_name, last_name, is_active, is_verified, created_at, updated_at) VALUES 
(5, 'pablo@niloft.com', '$2a$10$Z8bsd0n1VIb8OPZrztu90u1Q.G6sIOhWEJZSBy.SGf382MNCSRPJy', 'Pablo', 'Melegatti', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Reiniciar secuencia para IDs automáticos
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- =====================================================
-- DATOS DE GAMIFICACIÓN PARA USUARIOS DE PRUEBA  
-- =====================================================

-- USUARIO NIVEL 1 - PRINCIPIANTE (0 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_1', '1', 0, 1, 0, 0, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- USUARIO NIVEL 3 - SMART SAVER (200 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_3', '2', 200, 3, 15, 25, 2, 5, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP);

-- USUARIO NIVEL 5 - FINANCIAL PLANNER (700 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_5', '3', 700, 5, 45, 87, 4, 12, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP);

-- USUARIO NIVEL 10 - FINANCIAL MAGNATE (5500 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_nivel_10', '4', 5500, 10, 120, 275, 8, 30, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP);

-- USUARIO PRINCIPAL PABLO (NIVEL 7 - FINANCIAL EXPERT) (2000 XP)
INSERT INTO user_gamification (id, user_id, total_xp, current_level, insights_viewed, actions_completed, achievements_count, current_streak, last_activity, created_at, updated_at) VALUES 
('gam_user_pablo', '5', 2000, 7, 85, 150, 6, 20, CURRENT_TIMESTAMP - INTERVAL '15 minutes', CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP);

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

-- =====================================================
-- TRANSACCIONES PARA USUARIOS DE DIFERENTES NIVELES
-- =====================================================

-- USUARIO NIVEL 1 (PRINCIPIANTE) - Transacciones básicas
INSERT INTO incomes (id, user_id, amount, description, created_at) VALUES 
('inc_nivel1_01', '1', 120000.00, 'Sueldo básico', '2025-06-01 08:00:00');

INSERT INTO expenses (id, user_id, amount, description, category_id, paid, due_date, created_at) VALUES 
('exp_nivel1_01', '1', 45000.00, 'Supermercado', 'cat_nivel1_aliment', true, '2025-06-05', '2025-06-05 18:30:00'),
('exp_nivel1_02', '1', 15000.00, 'Transporte público', 'cat_nivel1_transp', true, '2025-06-10', '2025-06-10 08:15:00');

-- USUARIO NIVEL 3 (SMART SAVER) - Más variedad
INSERT INTO incomes (id, user_id, amount, description, created_at) VALUES 
('inc_nivel3_01', '2', 150000.00, 'Sueldo mensual', '2025-06-01 08:00:00');

INSERT INTO expenses (id, user_id, amount, description, category_id, paid, due_date, created_at) VALUES 
('exp_nivel3_01', '2', 35000.00, 'Compras inteligentes', 'cat_nivel3_aliment', true, '2025-06-03', '2025-06-03 16:20:00'),
('exp_nivel3_02', '2', 12000.00, 'Netflix', 'cat_nivel3_entret', true, '2025-06-15', '2025-06-15 20:30:00'),
('exp_nivel3_03', '2', 18000.00, 'Uber', 'cat_nivel3_transp', true, '2025-06-08', '2025-06-08 19:00:00');

-- USUARIO NIVEL 5 (FINANCIAL PLANNER) - Diversificado
INSERT INTO incomes (id, user_id, amount, description, created_at) VALUES 
('inc_nivel5_01', '3', 200000.00, 'Sueldo senior', '2025-06-01 08:00:00'),
('inc_nivel5_02', '3', 75000.00, 'Freelance', '2025-06-15 14:30:00');

INSERT INTO expenses (id, user_id, amount, description, category_id, paid, due_date, created_at) VALUES 
('exp_nivel5_01', '3', 40000.00, 'Alimentación planificada', 'cat_nivel5_aliment', true, '2025-06-02', '2025-06-02 17:45:00'),
('exp_nivel5_02', '3', 18000.00, 'Seguro médico', 'cat_nivel5_salud', true, '2025-06-12', '2025-06-12 10:30:00'),
('exp_nivel5_03', '3', 25000.00, 'Entretenimiento', 'cat_nivel5_entret', true, '2025-06-20', '2025-06-20 21:00:00');

-- USUARIO NIVEL 10 (FINANCIAL MAGNATE) - Complejo
INSERT INTO incomes (id, user_id, amount, description, created_at) VALUES 
('inc_nivel10_01', '4', 350000.00, 'Sueldo ejecutivo', '2025-06-01 08:00:00'),
('inc_nivel10_02', '4', 150000.00, 'Dividendos', '2025-06-05 11:30:00'),
('inc_nivel10_03', '4', 200000.00, 'Consultoría', '2025-06-20 16:00:00');

INSERT INTO expenses (id, user_id, amount, description, category_id, paid, due_date, created_at) VALUES 
('exp_nivel10_01', '4', 60000.00, 'Alimentación premium', 'cat_nivel10_aliment', true, '2025-06-03', '2025-06-03 19:30:00'),
('exp_nivel10_02', '4', 85000.00, 'Inversiones', 'cat_nivel10_invest', true, '2025-06-08', '2025-06-08 14:00:00'),
('exp_nivel10_03', '4', 45000.00, 'Entretenimiento VIP', 'cat_nivel10_entret', true, '2025-06-18', '2025-06-18 21:45:00'),
('exp_nivel10_04', '4', 30000.00, 'Salud premium', 'cat_nivel10_salud', true, '2025-06-25', '2025-06-25 11:00:00');

-- =====================================================
-- TRANSACCIONES PARA PABLO (ID 5) - USUARIO PRINCIPAL
-- =====================================================

-- Ingresos de abril 2025 para Pablo (ID 5)
INSERT INTO incomes (id, user_id, amount, description, created_at) VALUES 
('inc_abril_01', '5', 180000.00, 'Sueldo Abril', '2025-04-01 08:00:00'),
('inc_abril_02', '5', 35000.00, 'Freelance Abril', '2025-04-15 14:30:00'),
('inc_abril_03', '5', 15000.00, 'Venta productos', '2025-04-28 10:15:00');

-- Gastos de abril 2025 para Pablo (ID 5)
INSERT INTO expenses (id, user_id, amount, description, category_id, paid, due_date, created_at) VALUES 
('exp_abril_01', '5', 45000.00, 'Supermercado semanal', 'cat_aliment', true, '2025-04-05', '2025-04-02 09:30:00'),
('exp_abril_02', '5', 12000.00, 'Combustible', 'cat_transp', true, '2025-04-08', '2025-04-07 16:45:00'),
('exp_abril_03', '5', 8500.00, 'Netflix + Spotify', 'cat_entret', true, '2025-04-10', '2025-04-09 11:20:00'),
('exp_abril_04', '5', 25000.00, 'Farmacia', 'cat_salud', false, '2025-04-30', '2025-04-14 13:15:00'),
('exp_abril_05', '5', 18000.00, 'Uber/Taxi', 'cat_transp', true, '2025-04-18', '2025-04-17 19:30:00'),
('exp_abril_06', '5', 32000.00, 'Compras varias', 'cat_aliment', true, '2025-04-22', '2025-04-21 12:00:00'),
('exp_abril_07', '5', 6500.00, 'Café y snacks', 'cat_entret', true, '2025-04-25', '2025-04-24 15:45:00');

-- Ingresos de junio 2025
INSERT INTO incomes (id, user_id, amount, description, created_at) VALUES 
('inc_junio_01', '5', 185000.00, 'Sueldo Junio', '2025-06-01 08:00:00'),
('inc_junio_02', '5', 42000.00, 'Freelance Junio', '2025-06-16 14:30:00'),
('inc_junio_03', '5', 18000.00, 'Comisión ventas', '2025-06-27 10:15:00');

-- Gastos de junio 2025
INSERT INTO expenses (id, user_id, amount, description, category_id, paid, due_date, created_at) VALUES 
('exp_junio_01', '5', 48000.00, 'Supermercado junio', 'cat_aliment', true, '2025-06-03', '2025-06-01 10:30:00'),
('exp_junio_02', '5', 15000.00, 'Combustible junio', 'cat_transp', true, '2025-06-05', '2025-06-04 16:45:00'),
('exp_junio_03', '5', 8500.00, 'Suscripciones', 'cat_entret', true, '2025-06-10', '2025-06-08 11:20:00'),
('exp_junio_04', '5', 28000.00, 'Consulta médica', 'cat_salud', true, '2025-06-12', '2025-06-11 14:15:00'),
('exp_junio_05', '5', 22000.00, 'Transporte público', 'cat_transp', true, '2025-06-18', '2025-06-15 08:30:00'),
('exp_junio_06', '5', 35000.00, 'Cena restaurante', 'cat_entret', true, '2025-06-20', '2025-06-19 20:00:00'),
('exp_junio_07', '5', 12000.00, 'Compras online', 'cat_aliment', false, '2025-06-30', '2025-06-23 16:45:00'),
('exp_junio_08', '5', 7500.00, 'Cinema + palomitas', 'cat_entret', true, '2025-06-28', '2025-06-26 19:30:00');

-- =====================================================
-- PRESUPUESTOS DE EJEMPLO
-- =====================================================

-- Presupuestos de ejemplo para abril 2025
INSERT INTO budgets (id, user_id, category_id, amount, period, period_start, period_end) VALUES 
('budget_aliment_04', '5', 'cat_aliment', 80000.00, 'monthly', '2025-04-01 00:00:00', '2025-04-30 23:59:59'),
('budget_transp_04', '5', 'cat_transp', 35000.00, 'monthly', '2025-04-01 00:00:00', '2025-04-30 23:59:59'),
('budget_entret_04', '5', 'cat_entret', 20000.00, 'monthly', '2025-04-01 00:00:00', '2025-04-30 23:59:59')
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- TRANSACCIONES RECURRENTES DE EJEMPLO
-- =====================================================

-- Transacciones recurrentes de ejemplo
INSERT INTO recurring_transactions (id, user_id, amount, description, category_id, type, frequency, next_date) VALUES 
('rec_sueldo_01', '5', 180000.00, 'Sueldo mensual', NULL, 'income', 'monthly', '2025-05-01'),
('rec_netflix_01', '5', 8500.00, 'Netflix + Spotify', 'cat_entret', 'expense', 'monthly', '2025-05-10'),
('rec_farmacia_01', '5', 15000.00, 'Medicamentos recurrentes', 'cat_salud', 'expense', 'monthly', '2025-05-15')
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- METAS DE AHORRO DE EJEMPLO
-- =====================================================

-- Metas de ahorro de ejemplo
INSERT INTO savings_goals (id, user_id, name, description, target_amount, current_amount, category, priority, target_date, monthly_target, weekly_target, daily_target) VALUES 
('goal_vacation_01', '5', 'Vacaciones Europa', 'Viaje a Europa para 2 personas por 15 días', 500000.00, 125000.00, 'vacation', 'high', '2025-12-01', 62500.00, 14423.08, 2060.44),
('goal_emergency_01', '5', 'Fondo de Emergencia', 'Fondo de emergencia para 6 meses de gastos', 1080000.00, 320000.00, 'emergency', 'high', '2025-09-01', 253333.33, 58461.54, 8351.65),
('goal_car_01', '5', 'Auto Nuevo', 'Ahorro para comprar un auto 0km', 800000.00, 180000.00, 'car', 'medium', '2026-03-01', 77500.00, 17857.14, 2551.02),
('goal_house_01', '5', 'Entrada Casa', 'Ahorro para entrada de casa propia', 2000000.00, 450000.00, 'house', 'high', '2026-12-01', 86111.11, 19871.79, 2838.83)
ON CONFLICT (id) DO NOTHING;

-- Transacciones de ahorro de ejemplo
INSERT INTO savings_transactions (id, goal_id, user_id, amount, type, description) VALUES 
('trans_vac_01', 'goal_vacation_01', '5', 50000.00, 'deposit', 'Depósito inicial vacaciones'),
('trans_vac_02', 'goal_vacation_01', '5', 25000.00, 'deposit', 'Ahorro mensual abril'),
('trans_vac_03', 'goal_vacation_01', '5', 30000.00, 'deposit', 'Ahorro mensual mayo'),
('trans_vac_04', 'goal_vacation_01', '5', 20000.00, 'deposit', 'Ahorro mensual junio'),
('trans_emer_01', 'goal_emergency_01', '5', 100000.00, 'deposit', 'Depósito inicial emergencia'),
('trans_emer_02', 'goal_emergency_01', '5', 80000.00, 'deposit', 'Ahorro abril'),
('trans_emer_03', 'goal_emergency_01', '5', 70000.00, 'deposit', 'Ahorro mayo'),
('trans_emer_04', 'goal_emergency_01', '5', 70000.00, 'deposit', 'Ahorro junio'),
('trans_car_01', 'goal_car_01', '5', 80000.00, 'deposit', 'Depósito inicial auto'),
('trans_car_02', 'goal_car_01', '5', 50000.00, 'deposit', 'Ahorro mensual abril'),
('trans_car_03', 'goal_car_01', '5', 50000.00, 'deposit', 'Ahorro mensual mayo'),
('trans_house_01', 'goal_house_01', '5', 200000.00, 'deposit', 'Depósito inicial casa'),
('trans_house_02', 'goal_house_01', '5', 100000.00, 'deposit', 'Ahorro abril'),
('trans_house_03', 'goal_house_01', '5', 75000.00, 'deposit', 'Ahorro mayo'),
('trans_house_04', 'goal_house_01', '5', 75000.00, 'deposit', 'Ahorro junio')
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- ACHIEVEMENTS PARA USUARIOS DE DIFERENTES NIVELES
-- =====================================================

-- Achievements para Usuario Nivel 1 (Principiante - sin logros completados)
INSERT INTO achievements (id, user_id, type, name, description, points, progress, target, completed, created_at, updated_at) VALUES 
-- Básicos en progreso
('ach_l1_trans_start', '1', 'transaction_starter', 'Primer Paso', 'Registra tu primera transacción', 25, 0, 1, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ach_l1_cat_creator', '1', 'category_creator', 'Creador de Categorías', 'Crea 5 categorías personalizadas', 75, 0, 5, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ach_l1_org_expert', '1', 'organization_expert', 'Experto en Organización', 'Categoriza 50 transacciones', 100, 0, 50, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Achievements para Usuario Nivel 3 (Algunos completados)
INSERT INTO achievements (id, user_id, type, name, description, points, progress, target, completed, unlocked_at, created_at, updated_at) VALUES 
-- Completados
('ach_l3_trans_start', '2', 'transaction_starter', 'Primer Paso', 'Registra tu primera transacción', 25, 1, 1, true, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '5 days'),
('ach_l3_weekly_warr', '2', 'weekly_warrior', 'Guerrero Semanal', 'Mantén una racha de 7 días', 100, 7, 7, true, CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '3 days'),
-- En progreso
('ach_l3_trans_app', '2', 'transaction_apprentice', 'Aprendiz Financiero', 'Registra 10 transacciones', 50, 8, 10, false, NULL, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP),
('ach_l3_cat_creator', '2', 'category_creator', 'Creador de Categorías', 'Crea 5 categorías personalizadas', 75, 3, 5, false, NULL, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Achievements para Usuario Nivel 5 (Varios completados)
INSERT INTO achievements (id, user_id, type, name, description, points, progress, target, completed, unlocked_at, created_at, updated_at) VALUES 
-- Completados
('ach_l5_trans_start', '3', 'transaction_starter', 'Primer Paso', 'Registra tu primera transacción', 25, 1, 1, true, CURRENT_TIMESTAMP - INTERVAL '14 days', CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '14 days'),
('ach_l5_trans_app', '3', 'transaction_apprentice', 'Aprendiz Financiero', 'Registra 10 transacciones', 50, 10, 10, true, CURRENT_TIMESTAMP - INTERVAL '12 days', CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '12 days'),
('ach_l5_cat_creator', '3', 'category_creator', 'Creador de Categorías', 'Crea 5 categorías personalizadas', 75, 5, 5, true, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '10 days'),
('ach_l5_weekly_warr', '3', 'weekly_warrior', 'Guerrero Semanal', 'Mantén una racha de 7 días', 100, 12, 7, true, CURRENT_TIMESTAMP - INTERVAL '8 days', CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '8 days'),
-- En progreso avanzado
('ach_l5_trans_master', '3', 'transaction_master', 'Maestro de Transacciones', 'Registra 100 transacciones', 200, 67, 100, false, NULL, CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP),
('ach_l5_org_expert', '3', 'organization_expert', 'Experto en Organización', 'Categoriza 50 transacciones', 100, 43, 50, false, NULL, CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Achievements para Usuario Nivel 10 (Máximo - Muchos completados)
INSERT INTO achievements (id, user_id, type, name, description, points, progress, target, completed, unlocked_at, created_at, updated_at) VALUES 
-- Todos los básicos completados
('ach_l10_trans_start', '4', 'transaction_starter', 'Primer Paso', 'Registra tu primera transacción', 25, 1, 1, true, CURRENT_TIMESTAMP - INTERVAL '44 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '44 days'),
('ach_l10_trans_app', '4', 'transaction_apprentice', 'Aprendiz Financiero', 'Registra 10 transacciones', 50, 10, 10, true, CURRENT_TIMESTAMP - INTERVAL '43 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '43 days'),
('ach_l10_trans_master', '4', 'transaction_master', 'Maestro de Transacciones', 'Registra 100 transacciones', 200, 100, 100, true, CURRENT_TIMESTAMP - INTERVAL '30 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '30 days'),
('ach_l10_cat_creator', '4', 'category_creator', 'Creador de Categorías', 'Crea 5 categorías personalizadas', 75, 5, 5, true, CURRENT_TIMESTAMP - INTERVAL '42 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '42 days'),
('ach_l10_org_expert', '4', 'organization_expert', 'Experto en Organización', 'Categoriza 50 transacciones', 100, 50, 50, true, CURRENT_TIMESTAMP - INTERVAL '35 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '35 days'),
('ach_l10_weekly_warr', '4', 'weekly_warrior', 'Guerrero Semanal', 'Mantén una racha de 7 días', 100, 30, 7, true, CURRENT_TIMESTAMP - INTERVAL '37 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '37 days'),
('ach_l10_data_exp', '4', 'data_explorer', 'Explorador de Datos', 'Revisa analytics 5 días seguidos', 150, 5, 5, true, CURRENT_TIMESTAMP - INTERVAL '39 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '39 days'),
('ach_l10_ai_partner', '4', 'ai_partner', 'Socio de IA', 'Utiliza 100 insights de IA', 500, 100, 100, true, CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '45 days', CURRENT_TIMESTAMP - INTERVAL '20 days')
ON CONFLICT (id) DO NOTHING;

-- =====================================================
-- ACCIONES DE USUARIO PARA DIFERENTES NIVELES
-- =====================================================

-- Acciones para Usuario Nivel 3 (Progreso moderado)
INSERT INTO user_actions (id, user_id, action_type, entity_type, entity_id, xp_earned, description, created_at) VALUES 
('action_l3_01', '2', 'create_expense', 'transaction', 'exp_1', 8, 'Creó primer gasto', CURRENT_TIMESTAMP - INTERVAL '5 days'),
('action_l3_02', '2', 'create_income', 'transaction', 'inc_1', 8, 'Registró primer ingreso', CURRENT_TIMESTAMP - INTERVAL '5 days'),
('action_l3_03', '2', 'view_dashboard', 'page', 'dashboard', 2, 'Visitó dashboard', CURRENT_TIMESTAMP - INTERVAL '4 days'),
('action_l3_04', '2', 'create_category', 'category', 'cat_1', 10, 'Creó categoría Alimentación', CURRENT_TIMESTAMP - INTERVAL '4 days'),
('action_l3_05', '2', 'view_analytics', 'page', 'analytics', 3, 'Revisó reportes', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('action_l3_06', '2', 'daily_login', 'user', '2', 5, 'Login diario - día 1', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('action_l3_07', '2', 'daily_login', 'user', '2', 5, 'Login diario - día 2', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('action_l3_08', '2', 'daily_login', 'user', '2', 5, 'Login diario - día 3', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('action_l3_09', '2', 'create_expense', 'transaction', 'exp_2', 8, 'Segundo gasto registrado', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('action_l3_10', '2', 'update_expense', 'transaction', 'exp_1', 5, 'Actualizó gasto', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Acciones para Usuario Nivel 5 (Progreso avanzado)
INSERT INTO user_actions (id, user_id, action_type, entity_type, entity_id, xp_earned, description, created_at) VALUES 
('action_l5_01', '3', 'complete_profile', 'user', '3', 50, 'Completó perfil', CURRENT_TIMESTAMP - INTERVAL '15 days'),
('action_l5_02', '3', 'create_category', 'category', 'cat_1', 10, 'Creó Alimentación', CURRENT_TIMESTAMP - INTERVAL '14 days'),
('action_l5_03', '3', 'create_category', 'category', 'cat_2', 10, 'Creó Transporte', CURRENT_TIMESTAMP - INTERVAL '14 days'),
('action_l5_04', '3', 'create_savings_goal', 'savings', 'goal_1', 15, 'Primera meta de ahorro', CURRENT_TIMESTAMP - INTERVAL '13 days'),
('action_l5_05', '3', 'weekly_streak', 'streak', 'week_1', 25, 'Primera semana completada', CURRENT_TIMESTAMP - INTERVAL '8 days'),
('action_l5_06', '3', 'create_budget', 'budget', 'budget_1', 20, 'Primer presupuesto', CURRENT_TIMESTAMP - INTERVAL '7 days'),
('action_l5_07', '3', 'daily_challenge_complete', 'challenge', 'daily_1', 20, 'Challenge diario completado', CURRENT_TIMESTAMP - INTERVAL '5 days'),
('action_l5_08', '3', 'view_monthly_report', 'report', 'jan_2025', 5, 'Reporte mensual revisado', CURRENT_TIMESTAMP - INTERVAL '4 days'),
('action_l5_09', '3', 'deposit_savings', 'savings', 'goal_1', 8, 'Depósito en meta', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('action_l5_10', '3', 'stay_within_budget', 'budget', 'budget_1', 15, 'Mantuvo presupuesto', CURRENT_TIMESTAMP - INTERVAL '2 days')
ON CONFLICT (id) DO NOTHING;

-- Acciones para Usuario Nivel 10 (Máximo nivel - muchas acciones)
INSERT INTO user_actions (id, user_id, action_type, entity_type, entity_id, xp_earned, description, created_at) VALUES 
('action_l10_01', '4', 'monthly_streak', 'streak', 'month_1', 100, 'Racha de 30 días', CURRENT_TIMESTAMP - INTERVAL '15 days'),
('action_l10_02', '4', 'achieve_savings_goal', 'savings', 'goal_emergency', 100, 'Completó fondo emergencia', CURRENT_TIMESTAMP - INTERVAL '30 days'),
('action_l10_03', '4', 'weekly_challenge_complete', 'challenge', 'week_12', 75, 'Challenge semanal completado', CURRENT_TIMESTAMP - INTERVAL '25 days'),
('action_l10_04', '4', 'use_ai_analysis', 'ai', 'analysis_1', 10, 'Usó análisis de IA', CURRENT_TIMESTAMP - INTERVAL '23 days'),
('action_l10_05', '4', 'apply_ai_suggestion', 'ai', 'suggestion_1', 25, 'Aplicó sugerencia IA', CURRENT_TIMESTAMP - INTERVAL '22 days'),
('action_l10_06', '4', 'view_insight', 'insight', 'insight_1', 1, 'Vio insight de IA', CURRENT_TIMESTAMP - INTERVAL '20 days'),
('action_l10_07', '4', 'understand_insight', 'insight', 'insight_1', 15, 'Entendió insight', CURRENT_TIMESTAMP - INTERVAL '20 days'),
('action_l10_08', '4', 'export_data', 'data', 'export_2024', 10, 'Exportó datos anuales', CURRENT_TIMESTAMP - INTERVAL '17 days'),
('action_l10_09', '4', 'view_pattern', 'pattern', 'spending_pattern', 2, 'Revisó patrón gastos', CURRENT_TIMESTAMP - INTERVAL '15 days'),
('action_l10_10', '4', 'use_suggestion', 'suggestion', 'budget_opt', 5, 'Usó sugerencia presupuesto', CURRENT_TIMESTAMP - INTERVAL '1 day')
ON CONFLICT (id) DO NOTHING; 