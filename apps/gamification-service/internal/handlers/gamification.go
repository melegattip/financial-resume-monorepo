package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/melegattip/financial-gamification-service/internal/core/domain"
	"github.com/melegattip/financial-gamification-service/internal/core/usecases"
	"github.com/melegattip/financial-gamification-service/internal/infrastructure/http/middleware"
)

// GamificationHandlers maneja las peticiones HTTP de gamificación
type GamificationHandlers struct {
	gamificationUseCase usecases.GamificationUseCase
}

// NewGamificationHandlers crea una nueva instancia de los handlers
func NewGamificationHandlers(gamificationUseCase usecases.GamificationUseCase) *GamificationHandlers {
	return &GamificationHandlers{
		gamificationUseCase: gamificationUseCase,
	}
}

// RegisterRoutes registra todas las rutas de gamificación
func (h *GamificationHandlers) RegisterRoutes(router *mux.Router) {
	// Debug endpoint (temporal) - no requiere auth
	router.HandleFunc("/debug/user/{userID}/actions", h.DebugUserActions).Methods("GET")

	// Rutas protegidas con JWT
	protected := router.PathPrefix("/api/v1/gamification").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)

	// User gamification endpoints
	protected.HandleFunc("/profile", h.GetUserProfile).Methods("GET")
	protected.HandleFunc("/stats", h.GetUserStats).Methods("GET")
	protected.HandleFunc("/achievements", h.GetUserAchievements).Methods("GET")
	protected.HandleFunc("/actions", h.RecordUserAction).Methods("POST")

	// Feature Gates endpoints
	protected.HandleFunc("/features", h.GetUserFeatures).Methods("GET")
	protected.HandleFunc("/features/{featureKey}/access", h.CheckFeatureAccess).Methods("GET")

	// Challenge endpoints
	protected.HandleFunc("/challenges/daily", h.GetDailyChallenges).Methods("GET")
	protected.HandleFunc("/challenges/weekly", h.GetWeeklyChallenges).Methods("GET")
	protected.HandleFunc("/challenges/progress", h.ProcessChallengeProgress).Methods("POST")

	// Public endpoints (no requieren autenticación)
	router.HandleFunc("/api/v1/gamification/action-types", h.GetActionTypes).Methods("GET")
	router.HandleFunc("/api/v1/gamification/levels", h.GetLevels).Methods("GET")
}

// GetUserProfile obtiene el perfil de gamificación del usuario
// @Summary Obtener perfil de gamificación
// @Description Obtiene el estado completo de gamificación del usuario autenticado
// @Tags gamification
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.UserGamification
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/profile [get]
func (h *GamificationHandlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		log.Printf("❌ [GetUserProfile] User ID not found in context")
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	log.Printf("🔍 [GetUserProfile] Processing request for user: %s", userID)

	gamification, err := h.gamificationUseCase.GetUserGamification(r.Context(), userID)
	if err != nil {
		log.Printf("❌ [GetUserProfile] Error for user %s: %v", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [GetUserProfile] Successfully retrieved profile for user: %s", userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gamification); err != nil {
		log.Printf("❌ [GetUserProfile] Error encoding response for user %s: %v", userID, err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetUserStats obtiene las estadísticas de gamificación del usuario
// @Summary Obtener estadísticas de gamificación
// @Description Obtiene estadísticas detalladas de gamificación del usuario
// @Tags gamification
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.GamificationStats
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/stats [get]
func (h *GamificationHandlers) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	stats, err := h.gamificationUseCase.GetGamificationStats(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetUserAchievements obtiene los achievements del usuario
// @Summary Obtener achievements del usuario
// @Description Obtiene todos los logros del usuario con su progreso
// @Tags gamification
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Achievement
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/achievements [get]
func (h *GamificationHandlers) GetUserAchievements(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		log.Printf("❌ [GetUserAchievements] User ID not found in context")
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	log.Printf("🔍 [GetUserAchievements] Processing request for user: %s", userID)

	achievements, err := h.gamificationUseCase.GetUserAchievements(r.Context(), userID)
	if err != nil {
		log.Printf("❌ [GetUserAchievements] Error for user %s: %v", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Asegurar que siempre devolvemos un array, incluso si está vacío
	if achievements == nil {
		achievements = []domain.Achievement{}
	}

	log.Printf("✅ [GetUserAchievements] Successfully retrieved %d achievements for user: %s", len(achievements), userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(achievements); err != nil {
		log.Printf("❌ [GetUserAchievements] Error encoding response for user %s: %v", userID, err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// RecordUserAction registra una acción del usuario y otorga XP
// @Summary Registrar acción del usuario
// @Description Registra una acción del usuario y otorga XP correspondiente
// @Tags gamification
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param action body usecases.RecordActionParams true "Datos de la acción"
// @Success 200 {object} usecases.ActionResult
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/actions [post]
func (h *GamificationHandlers) RecordUserAction(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var params usecases.RecordActionParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Asegurar que el userID viene del token JWT
	params.UserID = userID

	result, err := h.gamificationUseCase.RecordUserAction(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetActionTypes obtiene los tipos de acciones disponibles
// @Summary Obtener tipos de acciones
// @Description Obtiene todos los tipos de acciones disponibles con su XP base
// @Tags gamification
// @Produce json
// @Success 200 {array} usecases.ActionType
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/action-types [get]
func (h *GamificationHandlers) GetActionTypes(w http.ResponseWriter, r *http.Request) {
	actionTypes, err := h.gamificationUseCase.GetActionTypes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(actionTypes); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetLevels obtiene información sobre los niveles
// @Summary Obtener información de niveles
// @Description Obtiene información sobre todos los niveles disponibles
// @Tags gamification
// @Produce json
// @Success 200 {array} usecases.LevelInfo
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/levels [get]
func (h *GamificationHandlers) GetLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := h.gamificationUseCase.GetLevels(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(levels); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetUserFeatures obtiene todas las features disponibles para el usuario
// @Summary Obtener features del usuario
// @Description Obtiene todas las features desbloqueadas y bloqueadas para el usuario autenticado
// @Tags gamification
// @Security BearerAuth
// @Produce json
// @Success 200 {object} usecases.UserFeaturesResult
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/features [get]
func (h *GamificationHandlers) GetUserFeatures(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		log.Printf("❌ [GetUserFeatures] User ID not found in context")
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	log.Printf("🔍 [GetUserFeatures] Processing request for user: %s", userID)

	features, err := h.gamificationUseCase.GetUserFeatures(r.Context(), userID)
	if err != nil {
		log.Printf("❌ [GetUserFeatures] Error for user %s: %v", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [GetUserFeatures] Successfully retrieved features for user: %s", userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(features); err != nil {
		log.Printf("❌ [GetUserFeatures] Error encoding response for user %s: %v", userID, err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// CheckFeatureAccess verifica acceso a una feature específica
// @Summary Verificar acceso a feature
// @Description Verifica si el usuario tiene acceso a una feature específica
// @Tags gamification
// @Security BearerAuth
// @Produce json
// @Param featureKey path string true "Clave de la feature (SAVINGS_GOALS, BUDGETS, AI_INSIGHTS)"
// @Success 200 {object} usecases.FeatureAccessResult
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/features/{featureKey}/access [get]
func (h *GamificationHandlers) CheckFeatureAccess(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	featureKey := vars["featureKey"]
	if featureKey == "" {
		http.Error(w, "Feature key is required", http.StatusBadRequest)
		return
	}

	access, err := h.gamificationUseCase.CheckFeatureAccess(r.Context(), userID, featureKey)
	if err != nil {
		if err.Error() == "feature key not found: "+featureKey {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(access); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// ====================================
// CHALLENGE HANDLERS
// ====================================

// GetDailyChallenges obtiene los challenges diarios del usuario
// @Summary Obtener challenges diarios
// @Description Obtiene los challenges diarios disponibles y su progreso
// @Tags gamification,challenges
// @Security BearerAuth
// @Produce json
// @Success 200 {object} []domain.ChallengeResult
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/challenges/daily [get]
func (h *GamificationHandlers) GetDailyChallenges(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	challenges, err := h.gamificationUseCase.GetDailyChallenges(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(challenges); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetWeeklyChallenges obtiene los challenges semanales del usuario
// @Summary Obtener challenges semanales
// @Description Obtiene los challenges semanales disponibles y su progreso
// @Tags gamification,challenges
// @Security BearerAuth
// @Produce json
// @Success 200 {object} []domain.ChallengeResult
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/challenges/weekly [get]
func (h *GamificationHandlers) GetWeeklyChallenges(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	challenges, err := h.gamificationUseCase.GetWeeklyChallenges(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(challenges); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// ProcessChallengeProgressRequest representa la request para procesar progreso de challenges
type ProcessChallengeProgressRequest struct {
	ActionType  string `json:"action_type" binding:"required"`
	EntityType  string `json:"entity_type,omitempty"`
	EntityID    string `json:"entity_id,omitempty"`
	XPEarned    int    `json:"xp_earned"`
	Description string `json:"description,omitempty"`
}

// ProcessChallengeProgress procesa el progreso de challenges para una acción
// @Summary Procesar progreso de challenges
// @Description Procesa una acción del usuario y actualiza el progreso de challenges
// @Tags gamification,challenges
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ProcessChallengeProgressRequest true "Datos de la acción"
// @Success 200 {object} []domain.ChallengeResult
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gamification/challenges/progress [post]
func (h *GamificationHandlers) ProcessChallengeProgress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req ProcessChallengeProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ActionType == "" {
		http.Error(w, "action_type is required", http.StatusBadRequest)
		return
	}

	results, err := h.gamificationUseCase.ProcessChallengeProgress(r.Context(), userID, req.ActionType, req.EntityType, req.EntityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// DebugUserActions endpoint temporal para debuggear acciones de usuario
func (h *GamificationHandlers) DebugUserActions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 [DEBUG ENDPOINT] Obteniendo acciones para userID: %s", userID)

	// Para debugging, vamos a usar el método getTransactionCount que ya implementamos
	// Primero intentamos obtener la gamificación del usuario
	userGamification, err := h.gamificationUseCase.GetUserGamification(r.Context(), userID)
	if err != nil {
		log.Printf("❌ [DEBUG ENDPOINT] Error getting user gamification: %v", err)
		http.Error(w, "Error getting user gamification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener achievements para ver el progreso actual
	achievements, err := h.gamificationUseCase.GetUserAchievements(r.Context(), userID)
	if err != nil {
		log.Printf("❌ [DEBUG ENDPOINT] Error getting achievements: %v", err)
		http.Error(w, "Error getting achievements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encontrar el achievement de primer paso
	var firstTransactionAchievement *domain.Achievement
	for i, achievement := range achievements {
		if achievement.Type == "transaction_starter" {
			firstTransactionAchievement = &achievements[i]
			break
		}
	}

	debugResponse := map[string]interface{}{
		"user_id":                       userID,
		"user_gamification":             userGamification,
		"total_achievements":            len(achievements),
		"first_transaction_achievement": firstTransactionAchievement,
		"achievements":                  achievements,
		"note":                          "Check logs for detailed transaction count calculation",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debugResponse)
}
