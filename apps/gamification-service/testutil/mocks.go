package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/melegattip/financial-gamification-service/internal/core/domain"
	"github.com/melegattip/financial-gamification-service/internal/core/ports"
)

// MockGamificationRepository implementa GamificationRepository para tests
type MockGamificationRepository struct {
	users          map[string]*domain.UserGamification
	achievements   map[string][]domain.Achievement
	actions        map[string][]domain.UserAction
	userChallenges map[string][]domain.UserChallenge // Key: userID, Value: challenges for that user
}

func NewMockGamificationRepository() *MockGamificationRepository {
	return &MockGamificationRepository{
		users:          make(map[string]*domain.UserGamification),
		achievements:   make(map[string][]domain.Achievement),
		actions:        make(map[string][]domain.UserAction),
		userChallenges: make(map[string][]domain.UserChallenge),
	}
}

// SetupUser configura un usuario para tests
func (m *MockGamificationRepository) SetupUser(userID string, totalXP int, level int) {
	m.users[userID] = &domain.UserGamification{
		UserID:            userID,
		TotalXP:           totalXP,
		CurrentLevel:      level,
		InsightsViewed:    0,
		ActionsCompleted:  0,
		AchievementsCount: 0,
		CurrentStreak:     0,
		LastActivity:      time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Inicializar achievements básicos automáticamente
	if m.achievements[userID] == nil {
		m.initializeBasicAchievements(userID)
	}
}

// SetupUserWithXP configura un usuario con XP específico para tests de niveles
func (m *MockGamificationRepository) SetupUserWithXP(userID string, totalXP int) {
	gamification := &domain.UserGamification{
		UserID:            userID,
		TotalXP:           totalXP,
		InsightsViewed:    0,
		ActionsCompleted:  0,
		AchievementsCount: 0,
		CurrentStreak:     0,
		LastActivity:      time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Calcular nivel basado en XP
	gamification.CurrentLevel = gamification.CalculateLevel()

	m.users[userID] = gamification

	// Inicializar achievements básicos si no existen
	if m.achievements[userID] == nil {
		m.initializeBasicAchievements(userID)
	}
}

// initializeBasicAchievements crea los achievements básicos para un usuario
func (m *MockGamificationRepository) initializeBasicAchievements(userID string) {
	basicAchievements := []domain.Achievement{
		{
			ID:          "ach_ai_partner_" + userID,
			UserID:      userID,
			Type:        "ai_partner",
			Name:        "🤖 AI Explorer",
			Description: "Utiliza 10 insights de IA",
			Points:      100,
			Progress:    0,
			Target:      10,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "ach_action_taker_" + userID,
			UserID:      userID,
			Type:        "action_taker",
			Name:        "🎯 Action Taker",
			Description: "Completa 25 acciones",
			Points:      200,
			Progress:    0,
			Target:      25,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "ach_quick_learner_" + userID,
			UserID:      userID,
			Type:        "quick_learner",
			Name:        "⚡ Quick Learner",
			Description: "Marca 5 insights como entendidos",
			Points:      100,
			Progress:    0,
			Target:      5,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "ach_budget_master_" + userID,
			UserID:      userID,
			Type:        "budget_master",
			Name:        "💰 Budget Master",
			Description: "Crea 5 presupuestos",
			Points:      300,
			Progress:    0,
			Target:      5,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "ach_goal_setter_" + userID,
			UserID:      userID,
			Type:        "goal_setter",
			Name:        "🎯 Goal Setter",
			Description: "Crea 3 metas de ahorro",
			Points:      250,
			Progress:    0,
			Target:      3,
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	m.achievements[userID] = basicAchievements
}

// Implementación de GamificationRepository interface

func (m *MockGamificationRepository) Create(ctx context.Context, gamification *domain.UserGamification) error {
	if m.users[gamification.UserID] == nil {
		m.users[gamification.UserID] = gamification
		m.initializeBasicAchievements(gamification.UserID)
	}
	return nil
}

func (m *MockGamificationRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserGamification, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}

	// Auto-crear usuario si no existe (comportamiento para tests)
	newUser := &domain.UserGamification{
		UserID:            userID,
		TotalXP:           0,
		CurrentLevel:      0,
		InsightsViewed:    0,
		ActionsCompleted:  0,
		AchievementsCount: 0,
		CurrentStreak:     0,
		LastActivity:      time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	m.users[userID] = newUser
	m.initializeBasicAchievements(userID)

	return newUser, nil
}

func (m *MockGamificationRepository) Update(ctx context.Context, gamification *domain.UserGamification) error {
	m.users[gamification.UserID] = gamification
	return nil
}

func (m *MockGamificationRepository) Delete(ctx context.Context, userID string) error {
	delete(m.users, userID)
	delete(m.achievements, userID)
	delete(m.actions, userID)
	return nil
}

func (m *MockGamificationRepository) GetGamificationStats(ctx context.Context, userID string) (*domain.GamificationStats, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, ports.ErrGamificationNotFound
	}

	completedCount := 0
	for _, achievement := range m.achievements[userID] {
		if achievement.Completed {
			completedCount++
		}
	}

	return &domain.GamificationStats{
		UserID:                userID,
		TotalXP:               user.TotalXP,
		CurrentLevel:          user.CurrentLevel,
		XPToNextLevel:         user.XPToNextLevel(),
		ProgressPercent:       user.ProgressToNextLevel(),
		TotalAchievements:     len(m.achievements[userID]),
		CompletedAchievements: completedCount,
		CurrentStreak:         user.CurrentStreak,
		LastActivity:          user.LastActivity,
	}, nil
}

// Achievement methods

func (m *MockGamificationRepository) CreateAchievement(ctx context.Context, achievement *domain.Achievement) error {
	if m.achievements[achievement.UserID] == nil {
		m.achievements[achievement.UserID] = []domain.Achievement{}
	}
	m.achievements[achievement.UserID] = append(m.achievements[achievement.UserID], *achievement)
	return nil
}

func (m *MockGamificationRepository) GetAchievementsByUserID(ctx context.Context, userID string) ([]domain.Achievement, error) {
	if _, exists := m.achievements[userID]; exists {
		// Actualizar progreso automáticamente basado en acciones
		m.updateAchievementProgress(userID)
		return m.achievements[userID], nil
	}
	return []domain.Achievement{}, nil
}

func (m *MockGamificationRepository) GetAchievementByID(ctx context.Context, achievementID string) (*domain.Achievement, error) {
	for _, userAchievements := range m.achievements {
		for _, achievement := range userAchievements {
			if achievement.ID == achievementID {
				return &achievement, nil
			}
		}
	}
	return nil, ports.ErrAchievementNotFound
}

func (m *MockGamificationRepository) UpdateAchievement(ctx context.Context, achievement *domain.Achievement) error {
	userAchievements := m.achievements[achievement.UserID]
	for i, existingAchievement := range userAchievements {
		if existingAchievement.ID == achievement.ID {
			userAchievements[i] = *achievement
			return nil
		}
	}
	return ports.ErrAchievementNotFound
}

func (m *MockGamificationRepository) DeleteAchievement(ctx context.Context, achievementID string) error {
	for userID, userAchievements := range m.achievements {
		for i, achievement := range userAchievements {
			if achievement.ID == achievementID {
				m.achievements[userID] = append(userAchievements[:i], userAchievements[i+1:]...)
				return nil
			}
		}
	}
	return ports.ErrAchievementNotFound
}

// Action methods

func (m *MockGamificationRepository) CreateAction(ctx context.Context, action *domain.UserAction) error {
	if m.actions[action.UserID] == nil {
		m.actions[action.UserID] = []domain.UserAction{}
	}
	m.actions[action.UserID] = append(m.actions[action.UserID], *action)
	return nil
}

func (m *MockGamificationRepository) GetActionsByUserID(ctx context.Context, userID string) ([]domain.UserAction, error) {
	if actions, exists := m.actions[userID]; exists {
		return actions, nil
	}
	return []domain.UserAction{}, nil
}

func (m *MockGamificationRepository) GetActionsByUserIDAndPeriod(ctx context.Context, userID string, startDate, endDate string) ([]domain.UserAction, error) {
	return m.GetActionsByUserID(ctx, userID)
}

// updateAchievementProgress actualiza el progreso de achievements basado en acciones
func (m *MockGamificationRepository) updateAchievementProgress(userID string) {
	actions := m.actions[userID]
	if actions == nil {
		return
	}

	// Contar acciones por tipo
	actionCounts := make(map[string]int)
	for _, action := range actions {
		actionCounts[action.ActionType]++
	}

	// Actualizar progreso de cada achievement
	for i, achievement := range m.achievements[userID] {
		switch achievement.Type {
		case "ai_partner":
			// AI Explorer: ya no progresa con view_insight
			m.achievements[userID][i].Progress = 0
		case "action_taker":
			// Action Taker: contar complete_action
			m.achievements[userID][i].Progress = actionCounts["complete_action"]
		case "quick_learner":
			// Quick Learner: contar understand_insight
			m.achievements[userID][i].Progress = actionCounts["understand_insight"]
		case "budget_master":
			// Budget Master: contar create_budget
			m.achievements[userID][i].Progress = actionCounts["create_budget"]
		case "goal_setter":
			// Goal Setter: contar create_goal
			m.achievements[userID][i].Progress = actionCounts["create_goal"]
		}

		// Marcar como completado si alcanzó el target
		if m.achievements[userID][i].Progress >= achievement.Target && !achievement.Completed {
			m.achievements[userID][i].Completed = true
			now := time.Now()
			m.achievements[userID][i].UnlockedAt = &now
			m.achievements[userID][i].UpdatedAt = now
		}
	}
}

// Action delete method
func (m *MockGamificationRepository) DeleteAction(ctx context.Context, actionID string) error {
	for userID, userActions := range m.actions {
		for i, action := range userActions {
			if action.ID == actionID {
				m.actions[userID] = append(userActions[:i], userActions[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("action not found")
}

// Challenge methods (implementación mínima)

func (m *MockGamificationRepository) GetActiveChallenges(ctx context.Context, challengeType string) ([]domain.Challenge, error) {
	// Return realistic challenge data for testing
	challenges := []domain.Challenge{}

	if challengeType == "daily" {
		challenges = append(challenges,
			domain.Challenge{
				ID:                "daily_expense_1",
				ChallengeKey:      "daily_expense",
				Name:              "Daily Expense Tracker",
				Description:       "Create 3 expenses today",
				Icon:              "💰",
				ChallengeType:     "daily",
				RequirementType:   "count",
				RequirementTarget: 3,
				RequirementData: map[string]interface{}{
					"actions": []interface{}{"create_expense"},
				},
				XPReward:  15,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			domain.Challenge{
				ID:                "daily_login_1",
				ChallengeKey:      "daily_login",
				Name:              "Daily Login",
				Description:       "Login once today",
				Icon:              "🎯",
				ChallengeType:     "daily",
				RequirementType:   "count",
				RequirementTarget: 1,
				RequirementData: map[string]interface{}{
					"actions": []interface{}{"login", "view_dashboard"},
				},
				XPReward:  5,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		)
	}

	if challengeType == "weekly" {
		challenges = append(challenges,
			domain.Challenge{
				ID:                "weekly_insights_1",
				ChallengeKey:      "weekly_insights",
				Name:              "Weekly Insight Explorer",
				Description:       "View 10 insights this week",
				Icon:              "📊",
				ChallengeType:     "weekly",
				RequirementType:   "count",
				RequirementTarget: 10,
				RequirementData: map[string]interface{}{
					"actions": []interface{}{"view_insight", "understand_insight"},
				},
				XPReward:  50,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			domain.Challenge{
				ID:                "weekly_budget_1",
				ChallengeKey:      "weekly_budget",
				Name:              "Weekly Budget Master",
				Description:       "Create 5 budget entries this week",
				Icon:              "📝",
				ChallengeType:     "weekly",
				RequirementType:   "count",
				RequirementTarget: 5,
				RequirementData: map[string]interface{}{
					"actions": []interface{}{"create_budget_entry", "update_budget"},
				},
				XPReward:  75,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		)
	}

	return challenges, nil
}

func (m *MockGamificationRepository) GetChallengeByKey(ctx context.Context, challengeKey string) (*domain.Challenge, error) {
	return nil, fmt.Errorf("challenge not found")
}

func (m *MockGamificationRepository) GetUserChallengesForDate(ctx context.Context, userID string, challengeDate time.Time, challengeType string) ([]domain.UserChallenge, error) {
	// Return user challenges for this user
	userChallenges, exists := m.userChallenges[userID]
	if !exists {
		// Initialize with some sample challenges with partial progress
		today := time.Now()
		m.userChallenges[userID] = []domain.UserChallenge{
			{
				ID:            "user_challenge_1",
				UserID:        userID,
				ChallengeID:   "daily_expense_1",
				ChallengeDate: today,
				Progress:      1, // Already made 1 expense out of 3
				Target:        3,
				Completed:     false,
				CreatedAt:     today,
				UpdatedAt:     today,
			},
			{
				ID:            "user_challenge_2",
				UserID:        userID,
				ChallengeID:   "weekly_insights_1",
				ChallengeDate: today,
				Progress:      3, // Already viewed 3 insights out of 10
				Target:        10,
				Completed:     false,
				CreatedAt:     today,
				UpdatedAt:     today,
			},
		}
		userChallenges = m.userChallenges[userID]
	}

	// Filter by challenge type if provided
	var filtered []domain.UserChallenge
	for _, uc := range userChallenges {
		// Simple date filtering (same day)
		if challengeDate.Day() == uc.ChallengeDate.Day() {
			filtered = append(filtered, uc)
		}
	}

	return filtered, nil
}

func (m *MockGamificationRepository) CreateOrUpdateUserChallenge(ctx context.Context, userChallenge *domain.UserChallenge) error {
	// Update or create the user challenge
	userChallenges := m.userChallenges[userChallenge.UserID]

	// Find existing challenge by ID
	found := false
	for i, uc := range userChallenges {
		if uc.ID == userChallenge.ID {
			userChallenges[i] = *userChallenge
			found = true
			break
		}
	}

	// If not found, add it
	if !found {
		userChallenges = append(userChallenges, *userChallenge)
	}

	m.userChallenges[userChallenge.UserID] = userChallenges
	return nil
}

func (m *MockGamificationRepository) GetUserChallengeByID(ctx context.Context, userChallengeID string) (*domain.UserChallenge, error) {
	return nil, fmt.Errorf("user challenge not found")
}

func (m *MockGamificationRepository) UpdateChallengeProgressTracking(ctx context.Context, tracking *domain.ChallengeProgressTracking) error {
	return nil // Implementación mínima para tests
}

func (m *MockGamificationRepository) GetChallengeProgressTracking(ctx context.Context, userID string, challengeDate time.Time) ([]domain.ChallengeProgressTracking, error) {
	return []domain.ChallengeProgressTracking{}, nil
}

// AchievementTestData contiene datos helper para tests de achievements
type AchievementTestData struct {
	ActionType string
	EntityType string
	Count      int
}

// GetAchievementName retorna el nombre de un achievement por tipo
func GetAchievementName(achievementType string) string {
	names := map[string]string{
		"ai_partner":    "🤖 AI Explorer",
		"action_taker":  "🎯 Action Taker",
		"quick_learner": "⚡ Quick Learner",
		"data_explorer": "📊 Data Analyst",
		"budget_master": "💰 Budget Master",
		"goal_setter":   "🎯 Goal Setter",
	}
	return names[achievementType]
}

// GetAchievementPoints retorna los puntos de un achievement por tipo
func GetAchievementPoints(achievementType string) int {
	points := map[string]int{
		"ai_partner":    100,
		"action_taker":  200,
		"quick_learner": 100,
		"data_explorer": 150,
		"budget_master": 300,
		"goal_setter":   250,
	}
	return points[achievementType]
}

// GetAchievementTarget retorna el target de un achievement por tipo
func GetAchievementTarget(achievementType string) int {
	targets := map[string]int{
		"ai_partner":    10, // 10 insights utilizados
		"action_taker":  25, // 25 acciones completadas
		"quick_learner": 5,  // 5 insights entendidos
		"data_explorer": 50, // 50 views de analytics
		"budget_master": 5,  // 5 presupuestos creados
		"goal_setter":   3,  // 3 metas de ahorro creadas
	}
	return targets[achievementType]
}

// FindAchievementByType busca un achievement por tipo
func FindAchievementByType(achievements []domain.Achievement, achievementType string) *domain.Achievement {
	for _, ach := range achievements {
		if ach.Type == achievementType {
			return &ach
		}
	}
	return nil
}
