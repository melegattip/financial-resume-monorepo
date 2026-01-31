package usecases

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/melegattip/financial-gamification-service/internal/core/domain"
	"github.com/melegattip/financial-gamification-service/internal/core/ports"
)

// GamificationUseCase define los casos de uso de gamificación
type GamificationUseCase interface {
	// User gamification management
	GetUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error)
	InitializeUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error)

	// Action tracking
	RecordUserAction(ctx context.Context, params RecordActionParams) (*ActionResult, error)

	// Achievement management
	GetUserAchievements(ctx context.Context, userID string) ([]domain.Achievement, error)
	CheckAndUpdateAchievements(ctx context.Context, userID string) ([]domain.Achievement, []domain.Achievement, error)

	// Statistics
	GetGamificationStats(ctx context.Context, userID string) (*domain.GamificationStats, error)
	GetActionTypes(ctx context.Context) ([]ActionType, error)
	GetLevels(ctx context.Context) ([]LevelInfo, error)

	// Feature Gates
	CheckFeatureAccess(ctx context.Context, userID string, featureKey string) (*FeatureAccessResult, error)
	GetUserFeatures(ctx context.Context, userID string) (*UserFeaturesResult, error)

	// Daily Challenges
	GetDailyChallenges(ctx context.Context, userID string) ([]UserChallengeResult, error)
	GetWeeklyChallenges(ctx context.Context, userID string) ([]UserChallengeResult, error)
	ProcessChallengeProgress(ctx context.Context, userID string, actionType string, entityType string, entityID string) (*ChallengeProgressResult, error)
}

// RecordActionParams parámetros para registrar una acción del usuario
type RecordActionParams struct {
	UserID      string `json:"user_id"`
	ActionType  string `json:"action_type"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Description string `json:"description"`
}

// ActionResult resultado de registrar una acción
type ActionResult struct {
	XPEarned            int                  `json:"xp_earned"`
	NewLevel            int                  `json:"new_level"`
	LevelUp             bool                 `json:"level_up"`
	NewAchievements     []domain.Achievement `json:"new_achievements"`
	UpdatedAchievements []domain.Achievement `json:"updated_achievements"`
	TotalXP             int                  `json:"total_xp"`
}

// ActionType información sobre tipos de acciones
type ActionType struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	BaseXP      int     `json:"base_xp"`
	Multiplier  float64 `json:"multiplier,omitempty"`
}

// LevelInfo información sobre niveles
type LevelInfo struct {
	Level      int    `json:"level"`
	Name       string `json:"name"`
	XPRequired int    `json:"xp_required"`
	XPToNext   int    `json:"xp_to_next,omitempty"`
}

// FeatureAccessResult resultado de verificación de acceso a feature
type FeatureAccessResult struct {
	FeatureKey    string `json:"feature_key"`
	HasAccess     bool   `json:"has_access"`
	RequiredLevel int    `json:"required_level"`
	CurrentLevel  int    `json:"current_level"`
	XPNeeded      int    `json:"xp_needed"`
	FeatureName   string `json:"feature_name"`
	FeatureIcon   string `json:"feature_icon"`
	Description   string `json:"description"`
	// Trial info (si el acceso está habilitado temporalmente para nuevos usuarios)
	TrialActive bool       `json:"trial_active"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
}

// UserFeaturesResult resultado con todas las features del usuario
type UserFeaturesResult struct {
	UserLevel        int                   `json:"user_level"`
	TotalXP          int                   `json:"total_xp"`
	UnlockedFeatures []string              `json:"unlocked_features"`
	LockedFeatures   []FeatureAccessResult `json:"locked_features"`
}

// UserChallengeResult resultado de challenge de usuario con progreso
type UserChallengeResult struct {
	ID              string `json:"id"`
	ChallengeKey    string `json:"challenge_key"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Icon            string `json:"icon"`
	XPReward        int    `json:"xp_reward"`
	Progress        int    `json:"progress"`
	Target          int    `json:"target"`
	Completed       bool   `json:"completed"`
	ProgressPercent int    `json:"progress_percent"`
	TimeRemaining   string `json:"time_remaining"`
}

// ChallengeProgressResult resultado de actualizar progreso de challenges
type ChallengeProgressResult struct {
	UpdatedChallenges   []UserChallengeResult `json:"updated_challenges"`
	CompletedChallenges []UserChallengeResult `json:"completed_challenges"`
	TotalXPEarned       int                   `json:"total_xp_earned"`
}

// gamificationService implementa GamificationUseCase
type gamificationService struct {
	gamificationRepo ports.GamificationRepository
}

// NewGamificationUseCase crea una nueva instancia del servicio de gamificación
func NewGamificationUseCase(gamificationRepo ports.GamificationRepository) GamificationUseCase {
	return &gamificationService{
		gamificationRepo: gamificationRepo,
	}
}

// GetUserGamification obtiene el estado de gamificación del usuario
func (s *gamificationService) GetUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error) {
	gamification, err := s.gamificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Si no existe, inicializar automáticamente
		if err == ports.ErrGamificationNotFound {
			return s.InitializeUserGamification(ctx, userID)
		}
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Actualizar nivel calculado
	gamification.CurrentLevel = gamification.CalculateLevel()

	return gamification, nil
}

// InitializeUserGamification inicializa la gamificación para un nuevo usuario
func (s *gamificationService) InitializeUserGamification(ctx context.Context, userID string) (*domain.UserGamification, error) {
	gamification := &domain.UserGamification{
		ID:                domain.NewID(),
		UserID:            userID,
		TotalXP:           0,
		CurrentLevel:      0, // Se calculará después
		InsightsViewed:    0,
		ActionsCompleted:  0,
		AchievementsCount: 0,
		CurrentStreak:     0,
		LastActivity:      time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Calcular nivel inicial correctamente (0 XP = Nivel 1)
	gamification.CurrentLevel = gamification.CalculateLevel()

	err := s.gamificationRepo.Create(ctx, gamification)
	if err != nil {
		return nil, fmt.Errorf("error creating user gamification: %w", err)
	}

	// Inicializar achievements básicos
	err = s.initializeBasicAchievements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error initializing achievements: %w", err)
	}

	return gamification, nil
}

// RecordUserAction registra una acción del usuario y otorga XP
func (s *gamificationService) RecordUserAction(ctx context.Context, params RecordActionParams) (*ActionResult, error) {
	// Obtener gamificación actual del usuario
	gamification, err := s.GetUserGamification(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Bloquear completamente efectos por vistas de insights (IA Financiera)
	// - No XP
	// - No registro de acción
	// - No incremento de contadores ni recálculo de achievements
	if params.ActionType == domain.ActionTypeViewInsight {
		return &ActionResult{
			XPEarned:            0,
			NewLevel:            gamification.CurrentLevel,
			LevelUp:             false,
			NewAchievements:     []domain.Achievement{},
			UpdatedAchievements: []domain.Achievement{},
			TotalXP:             gamification.TotalXP,
		}, nil
	}

	// 🛡️ ANTI-FARMING: Verificación de duplicados para understand_insight por entity_id específico
	// Esta es la protección REAL contra farming - verifica en base de datos
	if params.ActionType == domain.ActionTypeUnderstandInsight && params.EntityID != "" {
		// Buscar si ya existe una acción con el mismo ActionType y EntityID para este usuario
		existingActions, err := s.gamificationRepo.GetActionsByUserID(ctx, params.UserID)
		if err == nil {
			for _, action := range existingActions {
				if action.ActionType == params.ActionType && action.EntityID == params.EntityID {
					log.Printf("🛡️ [Anti-Farming] BLOCKED: User %s already understood insight %s", params.UserID, params.EntityID)
					// Ya existe: no otorgar XP ni registrar de nuevo
					return &ActionResult{
						XPEarned:            0,
						NewLevel:            gamification.CurrentLevel,
						LevelUp:             false,
						NewAchievements:     []domain.Achievement{},
						UpdatedAchievements: []domain.Achievement{},
						TotalXP:             gamification.TotalXP,
					}, nil
				}
			}
		}
	}

	// Idempotencia diaria: evitar sumar XP por acciones de vista que pueden ocurrir en cada refresh
	// Actualmente aplicamos a: view_dashboard (solo una vez por día)
	if params.ActionType == "view_dashboard" {
		// Calcular ventana de hoy [00:00:00, 23:59:59]
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

		actionsToday, err := s.gamificationRepo.GetActionsByUserIDAndPeriod(ctx, params.UserID, startOfDay.Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))
		if err == nil {
			for _, a := range actionsToday {
				if a.ActionType == params.ActionType {
					// Ya registrada hoy: no otorgar XP ni registrar de nuevo
					return &ActionResult{
						XPEarned:            0,
						NewLevel:            gamification.CurrentLevel,
						LevelUp:             false,
						NewAchievements:     []domain.Achievement{},
						UpdatedAchievements: []domain.Achievement{},
						TotalXP:             gamification.TotalXP,
					}, nil
				}
			}
		}
	}

	xpEarned := s.calculateXPForAction(params.ActionType, params.EntityType)
	if params.ActionType == domain.ActionTypeViewInsight {
		xpEarned = 0
	}

	// Si es daily_login, actualizar racha antes de persistir (para reflejar progreso inmediato)
	if params.ActionType == "daily_login" {
		s.updateStreakOnDailyLogin(gamification)
	}

	// Registrar la acción
	action := &domain.UserAction{
		ID:          domain.NewID(),
		UserID:      params.UserID,
		ActionType:  params.ActionType,
		EntityType:  params.EntityType,
		EntityID:    params.EntityID,
		XPEarned:    xpEarned,
		Description: params.Description,
		CreatedAt:   time.Now(),
	}

	err = s.gamificationRepo.CreateAction(ctx, action)
	if err != nil {
		return nil, fmt.Errorf("error creating user action: %w", err)
	}

	// Actualizar tracking de challenges (diario y semanal) para acciones relevantes
	s.trackChallengeProgress(ctx, params.UserID, params.ActionType)

	// Actualizar estadísticas del usuario
	oldLevel := gamification.CurrentLevel
	gamification.TotalXP += xpEarned
	gamification.CurrentLevel = gamification.CalculateLevel()
	gamification.LastActivity = time.Now()
	gamification.UpdatedAt = time.Now()

	// Actualizar contadores específicos
	switch params.ActionType {
	case domain.ActionTypeViewInsight:
		// No contar vistas de insights para evitar progreso/XP por refresh o entrada
	case domain.ActionTypeCompleteAction:
		gamification.ActionsCompleted++
	}

	err = s.gamificationRepo.Update(ctx, gamification)
	if err != nil {
		return nil, fmt.Errorf("error updating user gamification: %w", err)
	}

	// Verificar nuevos achievements
	newAchievements, updatedAchievements, err := s.CheckAndUpdateAchievements(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("error checking achievements: %w", err)
	}

	// 🏆 CALCULAR XP BONUS POR ACHIEVEMENTS COMPLETADOS
	achievementXP := 0
	for _, achievement := range newAchievements {
		achievementXP += achievement.Points
		log.Printf("🏆 [DEBUG] Sumando %d XP por completar achievement: %s", achievement.Points, achievement.Name)
	}

	// Si hay XP de achievements, actualizar gamification
	if achievementXP > 0 {
		oldLevel = gamification.CurrentLevel // Actualizar oldLevel por si sube de nivel con achievement XP
		gamification.TotalXP += achievementXP
		gamification.CurrentLevel = gamification.CalculateLevel()
		gamification.UpdatedAt = time.Now()

		err = s.gamificationRepo.Update(ctx, gamification)
		if err != nil {
			return nil, fmt.Errorf("error updating user gamification with achievement XP: %w", err)
		}

		log.Printf("🎉 [DEBUG] Otorgados %d XP adicionales por achievements completados. Total XP: %d", achievementXP, gamification.TotalXP)
	}

	result := &ActionResult{
		XPEarned:            xpEarned + achievementXP, // ✅ Incluir XP de achievements
		NewLevel:            gamification.CurrentLevel,
		LevelUp:             gamification.CurrentLevel > oldLevel,
		NewAchievements:     newAchievements,
		UpdatedAchievements: updatedAchievements,
		TotalXP:             gamification.TotalXP,
	}

	return result, nil
}

// GetUserAchievements obtiene todos los achievements del usuario
func (s *gamificationService) GetUserAchievements(ctx context.Context, userID string) ([]domain.Achievement, error) {
	achievements, err := s.gamificationRepo.GetAchievementsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user achievements: %w", err)
	}

	// Asegurar que siempre devolvemos un array, incluso si está vacío
	if achievements == nil {
		achievements = []domain.Achievement{}
	}

	return achievements, nil
}

// CheckAndUpdateAchievements verifica y actualiza achievements del usuario
func (s *gamificationService) CheckAndUpdateAchievements(ctx context.Context, userID string) ([]domain.Achievement, []domain.Achievement, error) {
	log.Printf("🔍 [DEBUG] Verificando achievements para userID: %s", userID)

	// Obtener gamificación actual
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Obtener achievements actuales
	achievements, err := s.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting achievements: %w", err)
	}

	log.Printf("🔍 [DEBUG] Achievements encontrados: %d", len(achievements))

	var newAchievements []domain.Achievement
	var updatedAchievements []domain.Achievement

	// Verificar cada tipo de achievement
	for _, achievement := range achievements {
		log.Printf("🔍 [DEBUG] Verificando achievement: %s (Type: %s, Progress: %d/%d)", achievement.Name, achievement.Type, achievement.Progress, achievement.Target)

		oldProgress := achievement.Progress
		newProgress := s.calculateAchievementProgress(achievement.Type, gamification)

		log.Printf("🔍 [DEBUG] Achievement %s: oldProgress=%d, newProgress=%d", achievement.Name, oldProgress, newProgress)

		if newProgress > oldProgress {
			log.Printf("✅ [DEBUG] Actualizando achievement %s: %d -> %d", achievement.Name, oldProgress, newProgress)
			achievement.UpdateProgress(newProgress)

			// Si se completó por primera vez, es un nuevo achievement
			if achievement.Completed && oldProgress < achievement.Target {
				log.Printf("🎉 [DEBUG] Achievement COMPLETADO por primera vez: %s", achievement.Name)
				newAchievements = append(newAchievements, achievement)
				gamification.AchievementsCount++
			} else {
				log.Printf("🔄 [DEBUG] Achievement actualizado (no completado): %s", achievement.Name)
				updatedAchievements = append(updatedAchievements, achievement)
			}

			// Actualizar en base de datos
			err = s.gamificationRepo.UpdateAchievement(ctx, &achievement)
			if err != nil {
				return nil, nil, fmt.Errorf("error updating achievement: %w", err)
			}
		} else {
			log.Printf("⚪ [DEBUG] No hay cambios en achievement: %s", achievement.Name)
		}
	}

	// Actualizar contador de achievements si hay nuevos
	if len(newAchievements) > 0 {
		err = s.gamificationRepo.Update(ctx, gamification)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating gamification: %w", err)
		}
	}

	return newAchievements, updatedAchievements, nil
}

// GetGamificationStats obtiene estadísticas de gamificación del usuario
func (s *gamificationService) GetGamificationStats(ctx context.Context, userID string) (*domain.GamificationStats, error) {
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	achievements, err := s.GetUserAchievements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting achievements: %w", err)
	}

	completedCount := 0
	for _, achievement := range achievements {
		if achievement.Completed {
			completedCount++
		}
	}

	stats := &domain.GamificationStats{
		UserID:                userID,
		TotalXP:               gamification.TotalXP,
		CurrentLevel:          gamification.CalculateLevel(),
		XPToNextLevel:         gamification.XPToNextLevel(),
		ProgressPercent:       gamification.ProgressToNextLevel(),
		TotalAchievements:     len(achievements),
		CompletedAchievements: completedCount,
		CurrentStreak:         gamification.CurrentStreak,
		LastActivity:          gamification.LastActivity,
	}

	return stats, nil
}

// GetActionTypes obtiene los tipos de acciones disponibles
func (s *gamificationService) GetActionTypes(ctx context.Context) ([]ActionType, error) {
	actionTypes := []ActionType{
		// 🏠 ACCIONES BÁSICAS (Disponibles desde Nivel 0)
		{Type: "view_dashboard", Description: "Ver dashboard principal", BaseXP: 2},
		{Type: "view_expenses", Description: "Ver lista de gastos", BaseXP: 1},
		{Type: "view_incomes", Description: "Ver lista de ingresos", BaseXP: 1},
		{Type: "view_categories", Description: "Ver categorías", BaseXP: 1},
		{Type: "view_analytics", Description: "Ver reportes básicos", BaseXP: 3},

		// 💰 TRANSACCIONES (Motor principal de XP)
		{Type: "create_expense", Description: "Registrar gasto", BaseXP: 8},
		{Type: "create_income", Description: "Registrar ingreso", BaseXP: 8},
		{Type: "update_expense", Description: "Actualizar gasto", BaseXP: 5},
		{Type: "update_income", Description: "Actualizar ingreso", BaseXP: 5},
		{Type: "delete_expense", Description: "Eliminar gasto", BaseXP: 3},
		{Type: "delete_income", Description: "Eliminar ingreso", BaseXP: 3},

		// 🏷️ ORGANIZACIÓN (Disponible desde Nivel 0)
		{Type: "create_category", Description: "Crear categoría personalizada", BaseXP: 10},
		{Type: "update_category", Description: "Actualizar categoría", BaseXP: 5},
		{Type: "assign_category", Description: "Categorizar transacción", BaseXP: 3},

		// 🎯 ENGAGEMENT Y STREAKS
		{Type: "daily_login", Description: "Login diario", BaseXP: 5},
		{Type: "weekly_streak", Description: "Racha de 7 días", BaseXP: 25},
		{Type: "monthly_streak", Description: "Racha de 30 días", BaseXP: 100},
		{Type: "complete_profile", Description: "Completar perfil", BaseXP: 50},

		// 🏆 CHALLENGES DIARIOS
		{Type: "daily_challenge_complete", Description: "Completar challenge diario", BaseXP: 20},
		{Type: "weekly_challenge_complete", Description: "Completar challenge semanal", BaseXP: 75},

		// 📊 ANÁLISIS Y REPORTES
		{Type: "view_monthly_report", Description: "Ver reporte mensual", BaseXP: 5},
		{Type: "view_category_breakdown", Description: "Ver desglose por categorías", BaseXP: 3},
		{Type: "export_data", Description: "Exportar datos", BaseXP: 10},

		// 🔓 FEATURES DESBLOQUEABLES
		// Metas de Ahorro (Nivel 3+)
		{Type: "create_savings_goal", Description: "Crear meta de ahorro", BaseXP: 15},
		{Type: "deposit_savings", Description: "Depositar en meta", BaseXP: 8},
		{Type: "achieve_savings_goal", Description: "Completar meta de ahorro", BaseXP: 100},

		// Presupuestos (Nivel 5+)
		{Type: "create_budget", Description: "Crear presupuesto", BaseXP: 20},
		{Type: "stay_within_budget", Description: "Mantener presupuesto", BaseXP: 15},

		// IA Financiera (Nivel 7+)
		{Type: "view_insight", Description: "Visualizar insight de IA", BaseXP: 0},
		{Type: "understand_insight", Description: "Marcar insight como entendido", BaseXP: 15},
		{Type: "use_ai_analysis", Description: "Usar análisis de IA", BaseXP: 10},
		{Type: "apply_ai_suggestion", Description: "Aplicar sugerencia de IA", BaseXP: 25},
	}

	return actionTypes, nil
}

// GetLevels obtiene información sobre los niveles
func (s *gamificationService) GetLevels(ctx context.Context) ([]LevelInfo, error) {
	levels := []LevelInfo{
		{Level: 1, Name: "Financial Newbie", XPRequired: 0},
		{Level: 2, Name: "Money Tracker", XPRequired: 75},
		{Level: 3, Name: "Smart Saver", XPRequired: 200}, // 🔓 METAS DE AHORRO
		{Level: 4, Name: "Budget Master", XPRequired: 400},
		{Level: 5, Name: "Financial Planner", XPRequired: 700}, // 🔓 PRESUPUESTOS
		{Level: 6, Name: "Investment Seeker", XPRequired: 1200},
		{Level: 7, Name: "Wealth Builder", XPRequired: 1800}, // 🔓 IA FINANCIERA
		{Level: 8, Name: "Financial Strategist", XPRequired: 2600},
		{Level: 9, Name: "Money Mentor", XPRequired: 3600},
		{Level: 10, Name: "Financial Magnate", XPRequired: 5500},
	}

	return levels, nil
}

// CheckFeatureAccess verifica si un usuario tiene acceso a una feature específica
func (s *gamificationService) CheckFeatureAccess(ctx context.Context, userID string, featureKey string) (*FeatureAccessResult, error) {
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	// Definir features gates
	featureGates := map[string]struct {
		requiredLevel int
		name          string
		icon          string
		description   string
	}{
		"SAVINGS_GOALS": {3, "Metas de Ahorro", "🎯", "Crea y gestiona objetivos de ahorro personalizados"},
		"BUDGETS":       {5, "Presupuestos", "📊", "Controla tus gastos con límites inteligentes por categoría"},
		"AI_INSIGHTS":   {7, "IA Financiera", "🧠", "Análisis inteligente con IA para decisiones financieras"},
	}

	feature, exists := featureGates[featureKey]
	if !exists {
		return nil, fmt.Errorf("feature key not found: %s", featureKey)
	}

	currentLevel := gamification.CurrentLevel
	hasAccess := currentLevel >= feature.requiredLevel

	// Trial de 10 días para nuevos usuarios en features clave
	trialDuration := 10 * 24 * time.Hour
	trialEnds := gamification.CreatedAt.Add(trialDuration)
	trialActive := time.Now().Before(trialEnds)
	if trialActive {
		switch featureKey {
		case "SAVINGS_GOALS", "BUDGETS", "AI_INSIGHTS":
			hasAccess = true
		}
	}

	// Calcular XP necesario si no tiene acceso
	xpNeeded := 0
	if !hasAccess {
		levels := []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}
		if feature.requiredLevel <= len(levels) {
			xpNeeded = levels[feature.requiredLevel-1] - gamification.TotalXP
		}
	}

	return &FeatureAccessResult{
		FeatureKey:    featureKey,
		HasAccess:     hasAccess,
		RequiredLevel: feature.requiredLevel,
		CurrentLevel:  currentLevel,
		XPNeeded:      xpNeeded,
		FeatureName:   feature.name,
		FeatureIcon:   feature.icon,
		Description:   feature.description,
		TrialActive:   trialActive,
		TrialEndsAt: func() *time.Time {
			if trialActive {
				return &trialEnds
			}
			return nil
		}(),
	}, nil
}

// GetUserFeatures obtiene todas las features disponibles para un usuario
func (s *gamificationService) GetUserFeatures(ctx context.Context, userID string) (*UserFeaturesResult, error) {
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}

	featureKeys := []string{"SAVINGS_GOALS", "BUDGETS", "AI_INSIGHTS"}
	var unlockedFeatures []string
	var lockedFeatures []FeatureAccessResult

	for _, featureKey := range featureKeys {
		access, err := s.CheckFeatureAccess(ctx, userID, featureKey)
		if err != nil {
			continue
		}

		if access.HasAccess {
			unlockedFeatures = append(unlockedFeatures, featureKey)
			// Si tiene acceso por trial, también agregar a locked features para mostrar banner
			if access.TrialActive {
				lockedFeatures = append(lockedFeatures, *access)
			}
		} else {
			lockedFeatures = append(lockedFeatures, *access)
		}
	}

	return &UserFeaturesResult{
		UserLevel:        gamification.CurrentLevel,
		TotalXP:          gamification.TotalXP,
		UnlockedFeatures: unlockedFeatures,
		LockedFeatures:   lockedFeatures,
	}, nil
}

// GetDailyChallenges obtiene los challenges diarios del usuario
func (s *gamificationService) GetDailyChallenges(ctx context.Context, userID string) ([]UserChallengeResult, error) {
	return s.getChallengesByType(ctx, userID, domain.ChallengeTypeDaily)
}

// GetWeeklyChallenges obtiene los challenges semanales del usuario
func (s *gamificationService) GetWeeklyChallenges(ctx context.Context, userID string) ([]UserChallengeResult, error) {
	return s.getChallengesByType(ctx, userID, domain.ChallengeTypeWeekly)
}

// ProcessChallengeProgress procesa el progreso de challenges basado en una acción
func (s *gamificationService) ProcessChallengeProgress(ctx context.Context, userID string, actionType string, entityType string, entityID string) (*ChallengeProgressResult, error) {
	now := time.Now()
	var updatedChallenges []UserChallengeResult
	var completedChallenges []UserChallengeResult
	totalXPEarned := 0

	// Procesar challenges diarios
	dailyResults, dailyXP, err := s.processChallengesForType(ctx, userID, actionType, entityType, entityID, domain.ChallengeTypeDaily, now)
	if err != nil {
		return nil, fmt.Errorf("error processing daily challenges: %w", err)
	}
	updatedChallenges = append(updatedChallenges, dailyResults.UpdatedChallenges...)
	completedChallenges = append(completedChallenges, dailyResults.CompletedChallenges...)
	totalXPEarned += dailyXP

	// Procesar challenges semanales
	weeklyResults, weeklyXP, err := s.processChallengesForType(ctx, userID, actionType, entityType, entityID, domain.ChallengeTypeWeekly, now)
	if err != nil {
		return nil, fmt.Errorf("error processing weekly challenges: %w", err)
	}
	updatedChallenges = append(updatedChallenges, weeklyResults.UpdatedChallenges...)
	completedChallenges = append(completedChallenges, weeklyResults.CompletedChallenges...)
	totalXPEarned += weeklyXP

	return &ChallengeProgressResult{
		UpdatedChallenges:   updatedChallenges,
		CompletedChallenges: completedChallenges,
		TotalXPEarned:       totalXPEarned,
	}, nil
}

// Métodos auxiliares privados

// getChallengesByType obtiene challenges por tipo para un usuario
func (s *gamificationService) getChallengesByType(ctx context.Context, userID string, challengeType string) ([]UserChallengeResult, error) {
	// Obtener challenges activos del tipo
	challenges, err := s.gamificationRepo.GetActiveChallenges(ctx, challengeType)
	if err != nil {
		return nil, fmt.Errorf("error getting active challenges: %w", err)
	}

	// ✅ NUEVO: Obtener nivel del usuario para filtrar
	gamification, err := s.GetUserGamification(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user gamification: %w", err)
	}
	userLevel := gamification.CurrentLevel

	now := time.Now()
	challengeDate := domain.GetChallengeDate(now, challengeType)

	// Obtener progreso del usuario para estos challenges
	userChallenges, err := s.gamificationRepo.GetUserChallengesForDate(ctx, userID, challengeDate, challengeType)
	if err != nil {
		return nil, fmt.Errorf("error getting user challenges: %w", err)
	}

	// Crear mapa de progreso por challenge ID
	progressMap := make(map[string]*domain.UserChallenge)
	for i := range userChallenges {
		progressMap[userChallenges[i].ChallengeID] = &userChallenges[i]
	}

	var results []UserChallengeResult
	for _, challenge := range challenges {
		// ✅ NUEVO: Filtrar challenges por nivel de usuario
		requiredLevel := s.getChallengeRequiredLevel(challenge.ChallengeKey)
		if userLevel < requiredLevel {
			continue // Skip challenges que requieren mayor nivel
		}

		userChallenge := progressMap[challenge.ID]

		progress := 0
		completed := false
		challengeID := ""

		if userChallenge != nil {
			progress = userChallenge.Progress
			completed = userChallenge.Completed
			challengeID = userChallenge.ID
		}

		progressPercent := 0
		if challenge.RequirementTarget > 0 {
			progressPercent = min(100, (progress*100)/challenge.RequirementTarget)
		}

		timeRemaining := s.calculateTimeRemaining(challengeType, now)

		results = append(results, UserChallengeResult{
			ID:              challengeID,
			ChallengeKey:    challenge.ChallengeKey,
			Name:            challenge.Name,
			Description:     challenge.Description,
			Icon:            challenge.Icon,
			XPReward:        challenge.XPReward,
			Progress:        progress,
			Target:          challenge.RequirementTarget,
			Completed:       completed,
			ProgressPercent: progressPercent,
			TimeRemaining:   timeRemaining,
		})
	}

	return results, nil
}

// ✅ NUEVO: Método auxiliar para determinar nivel requerido por challenge
func (s *gamificationService) getChallengeRequiredLevel(challengeKey string) int {
	switch challengeKey {
	// 🟢 CHALLENGES BÁSICOS (Nivel 1+) - Disponibles para todos
	case "transaction_master", "category_organizer", "analytics_explorer", "streak_keeper":
		return 1

	// 🟡 CHALLENGES INTERMEDIOS (Nivel 3+) - Requieren experiencia básica
	case "weekly_warrior", "engagement_hero", "data_explorer":
		return 3

	// 🟠 CHALLENGES AVANZADOS (Nivel 5+) - Requieren conocimiento de presupuestos
	case "budget_challenger", "savings_champion", "financial_planner":
		return 5

	// 🔴 CHALLENGES PREMIUM (Nivel 7+) - Requieren acceso a IA
	case "ai_pioneer", "insight_master", "ai_analyst", "financial_advisor":
		return 7

	// 💎 CHALLENGES ELITE (Nivel 9+) - Solo para usuarios expertos
	case "wealth_builder", "investment_guru", "financial_magnate":
		return 9

	default:
		// Por defecto, challenges básicos disponibles desde nivel 1
		return 1
	}
}

// processChallengesForType procesa challenges de un tipo específico
func (s *gamificationService) processChallengesForType(ctx context.Context, userID string, actionType string, entityType string, entityID string, challengeType string, now time.Time) (*ChallengeProgressResult, int, error) {
	challengeDate := domain.GetChallengeDate(now, challengeType)

	// Obtener challenges activos del tipo
	challenges, err := s.gamificationRepo.GetActiveChallenges(ctx, challengeType)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting active challenges: %w", err)
	}

	var updatedChallenges []UserChallengeResult
	var completedChallenges []UserChallengeResult
	totalXPEarned := 0

	for _, challenge := range challenges {
		// Verificar si la acción aplica para este challenge
		if !s.actionAppliesToChallenge(challenge, actionType, entityType) {
			continue
		}

		// Obtener o crear user challenge
		userChallenge, err := s.getOrCreateUserChallenge(ctx, userID, challenge.ID, challengeDate, challenge.RequirementTarget)
		if err != nil {
			return nil, 0, fmt.Errorf("error getting/creating user challenge: %w", err)
		}

		// Si ya está completado, skip
		if userChallenge.Completed {
			continue
		}

		// Calcular nuevo progreso
		newProgress := s.calculateNewProgress(ctx, userID, challenge, challengeDate, actionType, entityType, entityID)
		oldProgress := userChallenge.Progress
		userChallenge.UpdateProgress(newProgress)

		// Guardar el progreso actualizado
		err = s.gamificationRepo.CreateOrUpdateUserChallenge(ctx, userChallenge)
		if err != nil {
			return nil, 0, fmt.Errorf("error updating user challenge: %w", err)
		}

		// Crear resultado
		result := UserChallengeResult{
			ID:              userChallenge.ID,
			ChallengeKey:    challenge.ChallengeKey,
			Name:            challenge.Name,
			Description:     challenge.Description,
			Icon:            challenge.Icon,
			XPReward:        challenge.XPReward,
			Progress:        userChallenge.Progress,
			Target:          userChallenge.Target,
			Completed:       userChallenge.Completed,
			ProgressPercent: userChallenge.GetProgressPercentage(),
			TimeRemaining:   s.calculateTimeRemaining(challengeType, now),
		}

		// Solo agregar si hubo cambio en el progreso
		if newProgress > oldProgress {
			updatedChallenges = append(updatedChallenges, result)
		}

		// Si se completó recién, agregar XP
		if userChallenge.Completed && newProgress >= userChallenge.Target && oldProgress < userChallenge.Target {
			completedChallenges = append(completedChallenges, result)
			totalXPEarned += challenge.XPReward
		}
	}

	return &ChallengeProgressResult{
		UpdatedChallenges:   updatedChallenges,
		CompletedChallenges: completedChallenges,
		TotalXPEarned:       totalXPEarned,
	}, totalXPEarned, nil
}

// actionAppliesToChallenge verifica si una acción aplica para un challenge
func (s *gamificationService) actionAppliesToChallenge(challenge domain.Challenge, actionType string, entityType string) bool {
	if challenge.RequirementData == nil {
		return false
	}

	actions, ok := challenge.RequirementData["actions"].([]interface{})
	if !ok {
		return false
	}

	for _, action := range actions {
		if actionStr, ok := action.(string); ok && actionStr == actionType {
			return true
		}
	}

	return false
}

// getOrCreateUserChallenge obtiene o crea un user challenge
func (s *gamificationService) getOrCreateUserChallenge(ctx context.Context, userID string, challengeID string, challengeDate time.Time, target int) (*domain.UserChallenge, error) {
	// Intentar obtener challenges existentes para esta fecha
	userChallenges, err := s.gamificationRepo.GetUserChallengesForDate(ctx, userID, challengeDate, "")
	if err != nil {
		return nil, err
	}

	// Buscar el challenge específico
	for i := range userChallenges {
		if userChallenges[i].ChallengeID == challengeID {
			return &userChallenges[i], nil
		}
	}

	// Crear nuevo user challenge
	userChallenge := &domain.UserChallenge{
		ID:            domain.NewID(),
		UserID:        userID,
		ChallengeID:   challengeID,
		ChallengeDate: challengeDate,
		Progress:      0,
		Target:        target,
		Completed:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.gamificationRepo.CreateOrUpdateUserChallenge(ctx, userChallenge)
	if err != nil {
		return nil, err
	}

	return userChallenge, nil
}

// calculateNewProgress calcula el nuevo progreso basado en el tipo de challenge
func (s *gamificationService) calculateNewProgress(ctx context.Context, userID string, challenge domain.Challenge, challengeDate time.Time, actionType string, entityType string, entityID string) int {
	switch challenge.RequirementType {
	case domain.RequirementTypeTransactionCount:
		return s.calculateTransactionProgress(ctx, userID, challengeDate, challenge.RequirementData)
	case domain.RequirementTypeCategoryVariety:
		return s.calculateCategoryVarietyProgress(ctx, userID, challengeDate)
	case domain.RequirementTypeViewCombo:
		return s.calculateViewComboProgress(ctx, userID, challengeDate, challenge.RequirementData)
	case domain.RequirementTypeDailyLogin:
		return 1 // Daily login siempre es 1
	case domain.RequirementTypeDailyLoginCount:
		return s.calculateDailyLoginCountProgress(ctx, userID, challengeDate)
	default:
		return 0
	}
}

// ============================
// Rachas y tracking de challenges
// ============================

// updateStreakOnDailyLogin aplica la lógica de racha al efectuar un daily_login
func (s *gamificationService) updateStreakOnDailyLogin(g *domain.UserGamification) {
	now := time.Now()
	last := g.LastActivity

	// Normalizar fechas a día
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	lastDay := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, last.Location())

	daysDiff := int(today.Sub(lastDay).Hours() / 24)

	switch {
	case daysDiff == 0:
		// Ya hubo actividad hoy: no cambiar racha
	case daysDiff == 1:
		// Día consecutivo
		if g.CurrentStreak < 0 {
			g.CurrentStreak = 1
		} else {
			g.CurrentStreak++
		}
	default:
		// Salto de más de 1 día o última actividad en el futuro: reiniciar racha
		g.CurrentStreak = 1
	}
}

// trackChallengeProgress registra/actualiza el tracking para challenges por día y semana
func (s *gamificationService) trackChallengeProgress(ctx context.Context, userID string, actionType string) {
	now := time.Now()

	// Helper interno para upsert tracking
	upsert := func(challengeDate time.Time) {
		tracking := &domain.ChallengeProgressTracking{
			ID:            domain.NewID(),
			UserID:        userID,
			ChallengeDate: challengeDate,
			ActionType:    actionType,
			Count:         1,
			CreatedAt:     now,
		}
		if err := s.gamificationRepo.UpdateChallengeProgressTracking(ctx, tracking); err != nil {
			log.Printf("⚠️ [Gamification] Error updating challenge tracking (%s): %v", actionType, err)
		}
	}

	// Registrar para ventana diaria
	upsert(domain.GetChallengeDate(now, domain.ChallengeTypeDaily))
	// Registrar para ventana semanal (contabiliza días únicos con daily_login y otras acciones)
	upsert(domain.GetChallengeDate(now, domain.ChallengeTypeWeekly))
}

// calculateTimeRemaining calcula el tiempo restante para un tipo de challenge
func (s *gamificationService) calculateTimeRemaining(challengeType string, now time.Time) string {
	switch challengeType {
	case domain.ChallengeTypeDaily:
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		remaining := endOfDay.Sub(now)
		hours := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	case domain.ChallengeTypeWeekly:
		// Calcular hasta el domingo
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		daysUntilSunday := 7 - weekday
		endOfWeek := now.AddDate(0, 0, daysUntilSunday)
		endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 999999999, endOfWeek.Location())
		remaining := endOfWeek.Sub(now)
		days := int(remaining.Hours()) / 24
		hours := int(remaining.Hours()) % 24
		return fmt.Sprintf("%dd %dh", days, hours)
	default:
		return "Unknown"
	}
}

// Métodos auxiliares para cálculo de progreso (simplificados)
func (s *gamificationService) calculateTransactionProgress(ctx context.Context, userID string, challengeDate time.Time, requirementData map[string]interface{}) int {
	// Simplificado: contar acciones de tracking
	tracking, err := s.gamificationRepo.GetChallengeProgressTracking(ctx, userID, challengeDate)
	if err != nil {
		return 0
	}

	count := 0
	for _, t := range tracking {
		if t.ActionType == "create_expense" || t.ActionType == "create_income" {
			count += t.Count
		}
	}
	return count
}

func (s *gamificationService) calculateCategoryVarietyProgress(ctx context.Context, userID string, challengeDate time.Time) int {
	// Simplificado: contar categorías únicas usadas
	tracking, err := s.gamificationRepo.GetChallengeProgressTracking(ctx, userID, challengeDate)
	if err != nil {
		return 0
	}

	uniqueCategories := make(map[string]bool)
	for _, t := range tracking {
		if t.ActionType == "assign_category" && t.UniqueEntities != nil {
			if categories, ok := t.UniqueEntities["categories"].([]interface{}); ok {
				for _, cat := range categories {
					if catStr, ok := cat.(string); ok {
						uniqueCategories[catStr] = true
					}
				}
			}
		}
	}
	return len(uniqueCategories)
}

func (s *gamificationService) calculateViewComboProgress(ctx context.Context, userID string, challengeDate time.Time, requirementData map[string]interface{}) int {
	// Simplificado: verificar que se hayan hecho las vistas requeridas
	tracking, err := s.gamificationRepo.GetChallengeProgressTracking(ctx, userID, challengeDate)
	if err != nil {
		return 0
	}

	requiredActions, ok := requirementData["actions"].([]interface{})
	if !ok {
		return 0
	}

	actionMap := make(map[string]bool)
	for _, t := range tracking {
		actionMap[t.ActionType] = true
	}

	count := 0
	for _, action := range requiredActions {
		if actionStr, ok := action.(string); ok && actionMap[actionStr] {
			count++
		}
	}
	return count
}

func (s *gamificationService) calculateDailyLoginCountProgress(ctx context.Context, userID string, challengeDate time.Time) int {
	// Para challenges semanales: contar días únicos con daily_login
	tracking, err := s.gamificationRepo.GetChallengeProgressTracking(ctx, userID, challengeDate)
	if err != nil {
		return 0
	}

	uniqueDays := make(map[string]bool)
	for _, t := range tracking {
		if t.ActionType == "daily_login" {
			day := t.CreatedAt.Format("2006-01-02")
			uniqueDays[day] = true
		}
	}
	return len(uniqueDays)
}

// Métodos auxiliares privados

func (s *gamificationService) calculateXPForAction(actionType, entityType string) int {
	baseXP := s.getBaseXPForAction(actionType)
	multiplier := s.getMultiplierForEntity(entityType)
	return int(float64(baseXP) * multiplier)
}

func (s *gamificationService) getBaseXPForAction(actionType string) int {
	switch actionType {
	// 🏠 ACCIONES BÁSICAS (Disponibles desde Nivel 0)
	case "view_dashboard":
		return 2
	case "view_expenses", "view_incomes", "view_categories":
		return 1
	case "view_analytics":
		return 3

	// 💰 TRANSACCIONES (Motor principal de XP)
	case "create_expense", "create_income":
		return 8
	case "update_expense", "update_income":
		return 5
	case "delete_expense", "delete_income":
		return 3

	// 🏷️ ORGANIZACIÓN
	case "create_category":
		return 10
	case "update_category":
		return 5
	case "assign_category":
		return 3

	// 🎯 ENGAGEMENT Y STREAKS
	case "daily_login":
		return 5
	case "weekly_streak":
		return 25
	case "monthly_streak":
		return 100
	case "complete_profile":
		return 50

	// 🏆 CHALLENGES
	case "daily_challenge_complete":
		return 20
	case "weekly_challenge_complete":
		return 75

	// 📊 ANÁLISIS Y REPORTES
	case "view_monthly_report":
		return 5
	case "view_category_breakdown":
		return 3
	case "export_data":
		return 10

	// 🔓 FEATURES DESBLOQUEABLES
	// Metas de Ahorro (Nivel 3+)
	case "create_savings_goal":
		return 15
	case "deposit_savings":
		return 8
	case "achieve_savings_goal":
		return 100

	// Presupuestos (Nivel 5+)
	case "create_budget":
		return 20
	case "stay_within_budget":
		return 15

	// IA Financiera (Nivel 7+) — no XP por vistas de insights para evitar XP por refresh
	case domain.ActionTypeViewInsight:
		return 0
	case domain.ActionTypeUnderstandInsight:
		return 15
	case "use_ai_analysis":
		return 10
	case "apply_ai_suggestion":
		return 25
	case domain.ActionTypeCompleteAction:
		return 10
	case domain.ActionTypeViewPattern:
		return 2
	case domain.ActionTypeUseSuggestion:
		return 5
	default:
		return 1
	}
}

func (s *gamificationService) getMultiplierForEntity(entityType string) float64 {
	switch entityType {
	case domain.EntityTypeInsight:
		return 1.0
	case domain.EntityTypeSuggestion:
		return 1.2
	case domain.EntityTypePattern:
		return 1.1
	default:
		return 1.0
	}
}

func (s *gamificationService) calculateAchievementProgress(achievementType string, gamification *domain.UserGamification) int {
	log.Printf("🔍 [DEBUG] Calculando progreso para achievement: %s, userID: %s", achievementType, gamification.UserID)

	switch achievementType {
	// 💰 ACHIEVEMENTS DE TRANSACCIONES
	case "transaction_starter", "transaction_apprentice", "transaction_master":
		log.Printf("🎯 [DEBUG] Es un achievement de transacciones: %s", achievementType)
		// Obtener conteo real de transacciones creadas por el usuario
		count := s.getTransactionCount(gamification.UserID)
		log.Printf("🎯 [DEBUG] Retornando count: %d para achievement: %s", count, achievementType)
		return count

	// 🏷️ ACHIEVEMENTS DE ORGANIZACIÓN
	case "category_creator":
		// ✅ CORREGIDO: Contar acciones reales de create_category
		count := s.getCategoryCount(gamification.UserID)
		log.Printf("🎯 [DEBUG] Categorías reales contadas para user %s: %d", gamification.UserID, count)
		return count
	case "organization_expert":
		// ✅ CORREGIDO: Contar categorizaciones reales (assign_category)
		assignCount := s.getAssignCategoryCount(gamification.UserID)
		log.Printf("🎯 [DEBUG] Categorizaciones reales contadas para user %s: %d", gamification.UserID, assignCount)
		return assignCount

	// 🔥 ACHIEVEMENTS DE ENGAGEMENT Y STREAKS
	case "weekly_warrior":
		return gamification.CurrentStreak
	case "monthly_legend":
		return gamification.CurrentStreak

		// 📈 ACHIEVEMENTS DE PROGRESO Y ANÁLISIS
	case "data_explorer":
		// Simplificado: cada día de actividad = 1 vista de analytics
		return int(time.Since(gamification.CreatedAt).Hours() / 24)

	// 🤖 ACHIEVEMENTS DE IA (para usuarios que lleguen a nivel 7+)
	case domain.AchievementTypeAIPartner:
		return gamification.InsightsViewed
	case domain.AchievementTypeActionTaker:
		return gamification.ActionsCompleted
	case domain.AchievementTypeQuickLearner:
		return gamification.InsightsViewed / 10

	default:
		log.Printf("⚠️ [DEBUG] Achievement type desconocido: %s", achievementType)
		return 0
	}
}

func (s *gamificationService) initializeBasicAchievements(ctx context.Context, userID string) error {
	basicAchievements := []domain.Achievement{
		// 💰 ACHIEVEMENTS DE TRANSACCIONES (Base de progresión)
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "transaction_starter",
			Name:        "🌱 Primer Paso",
			Description: "Registra tu primera transacción",
			Points:      25,
			Progress:    0,
			Target:      1,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "transaction_apprentice",
			Name:        "📝 Aprendiz Financiero",
			Description: "Registra 10 transacciones",
			Points:      50,
			Progress:    0,
			Target:      10,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "transaction_master",
			Name:        "💎 Maestro de Transacciones",
			Description: "Registra 100 transacciones",
			Points:      200,
			Progress:    0,
			Target:      100,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// 🏷️ ACHIEVEMENTS DE ORGANIZACIÓN
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "category_creator",
			Name:        "🎨 Creador de Categorías",
			Description: "Crea 5 categorías personalizadas",
			Points:      75,
			Progress:    0,
			Target:      5,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "organization_expert",
			Name:        "📊 Expert en Organización",
			Description: "Categoriza 50 transacciones",
			Points:      100,
			Progress:    0,
			Target:      50,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// 🔥 ACHIEVEMENTS DE ENGAGEMENT Y STREAKS
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "weekly_warrior",
			Name:        "⚡ Guerrero Semanal",
			Description: "Mantén una racha de 7 días",
			Points:      100,
			Progress:    0,
			Target:      7,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "monthly_legend",
			Name:        "👑 Leyenda Mensual",
			Description: "Mantén una racha de 30 días",
			Points:      500,
			Progress:    0,
			Target:      30,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// 📈 ACHIEVEMENTS DE PROGRESO Y ANÁLISIS
		{
			ID:          domain.NewID(),
			UserID:      userID,
			Type:        "data_explorer",
			Name:        "🔍 Explorador de Datos",
			Description: "Revisa analytics 25 veces",
			Points:      75,
			Progress:    0,
			Target:      25,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, achievement := range basicAchievements {
		err := s.gamificationRepo.CreateAchievement(ctx, &achievement)
		if err != nil {
			return fmt.Errorf("error creating achievement %s: %w", achievement.Name, err)
		}
	}

	return nil
}

// getTransactionCount cuenta las transacciones creadas por el usuario
func (s *gamificationService) getTransactionCount(userID string) int {
	// Contar acciones de tipo create_expense y create_income
	ctx := context.Background()

	log.Printf("🔍 [DEBUG] Contando transacciones para userID: %s", userID)

	// Intentar obtener las acciones del usuario desde el repositorio
	// Si hay error, devolvemos 0 para no romper el flujo
	actions, err := s.gamificationRepo.GetActionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("❌ [DEBUG] Error obteniendo acciones para user %s: %v", userID, err)
		return 0
	}

	log.Printf("🔍 [DEBUG] Total acciones encontradas para user %s: %d", userID, len(actions))

	transactionCount := 0
	for _, action := range actions {
		log.Printf("🔍 [DEBUG] Acción: %s, Tipo: %s, Entity: %s", action.ID, action.ActionType, action.EntityType)
		if action.ActionType == "create_expense" || action.ActionType == "create_income" {
			transactionCount++
			log.Printf("✅ [DEBUG] Transacción contada: %s", action.ActionType)
		}
	}

	log.Printf("🎯 [DEBUG] Total transacciones contadas para user %s: %d", userID, transactionCount)
	return transactionCount
}

// getCategoryCount cuenta las categorías creadas por el usuario
func (s *gamificationService) getCategoryCount(userID string) int {
	// Contar acciones de tipo create_category
	ctx := context.Background()

	log.Printf("🔍 [DEBUG] Contando categorías para userID: %s", userID)

	// Intentar obtener las acciones del usuario desde el repositorio
	// Si hay error, devolvemos 0 para no romper el flujo
	actions, err := s.gamificationRepo.GetActionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("❌ [DEBUG] Error obteniendo acciones para user %s: %v", userID, err)
		return 0
	}

	log.Printf("🔍 [DEBUG] Total acciones encontradas para user %s: %d", userID, len(actions))

	categoryCount := 0
	for _, action := range actions {
		log.Printf("🔍 [DEBUG] Acción: %s, Tipo: %s, Entity: %s", action.ID, action.ActionType, action.EntityType)
		if action.ActionType == "create_category" {
			categoryCount++
			log.Printf("✅ [DEBUG] Categoría contada: %s", action.ActionType)
		}
	}

	log.Printf("🎯 [DEBUG] Total categorías contadas para user %s: %d", userID, categoryCount)
	return categoryCount
}

// getAssignCategoryCount cuenta cuántas veces el usuario categorizó transacciones
func (s *gamificationService) getAssignCategoryCount(userID string) int {
	ctx := context.Background()

	log.Printf("🔍 [DEBUG] Contando asignaciones de categoría para userID: %s", userID)

	actions, err := s.gamificationRepo.GetActionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("❌ [DEBUG] Error obteniendo acciones para user %s: %v", userID, err)
		return 0
	}

	count := 0
	for _, action := range actions {
		if action.ActionType == "assign_category" {
			count++
		}
	}

	log.Printf("🎯 [DEBUG] Total asignaciones de categoría contadas para user %s: %d", userID, count)
	return count
}
