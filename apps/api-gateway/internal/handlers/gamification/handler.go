package gamification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/http/helpers"
)

// ActionType información sobre tipos de acciones
type ActionType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	BaseXP      int    `json:"base_xp"`
}

// LevelInfo información sobre niveles
type LevelInfo struct {
	Level      int    `json:"level"`
	Name       string `json:"name"`
	XPRequired int    `json:"xp_required"`
	XPToNext   int    `json:"xp_to_next,omitempty"`
}

// Handler maneja las peticiones HTTP de gamificación
type Handler struct {
	gamificationUseCase usecases.GamificationUseCase
}

// NewHandler crea una nueva instancia del handler de gamificación
func NewHandler(gamificationUseCase usecases.GamificationUseCase) *Handler {
	return &Handler{
		gamificationUseCase: gamificationUseCase,
	}
}

// GetUserGamification obtiene el estado de gamificación del usuario
// @Summary Obtener estado de gamificación del usuario
// @Description Obtiene las estadísticas de gamificación del usuario autenticado
// @Tags gamification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.UserGamification
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/gamification/profile [get]
func (h *Handler) GetUserGamification(c *gin.Context) {
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	gamification, err := h.gamificationUseCase.GetUserGamification(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo datos de gamificación"})
		return
	}

	c.JSON(http.StatusOK, gamification)
}

// GetUserStats obtiene estadísticas detalladas de gamificación
// @Summary Obtener estadísticas de gamificación
// @Description Obtiene estadísticas detalladas incluyendo progreso a siguiente nivel
// @Tags gamification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.GamificationStats
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/gamification/stats [get]
func (h *Handler) GetUserStats(c *gin.Context) {
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.gamificationUseCase.GetGamificationStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo estadísticas"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetUserAchievements obtiene todos los achievements del usuario
// @Summary Obtener achievements del usuario
// @Description Obtiene todos los logros del usuario, completados y en progreso
// @Tags gamification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Achievement
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/gamification/achievements [get]
func (h *Handler) GetUserAchievements(c *gin.Context) {
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	achievements, err := h.gamificationUseCase.GetUserAchievements(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
		"total":        len(achievements),
	})
}

// RecordActionRequest estructura para registrar una acción
type RecordActionRequest struct {
	ActionType  string `json:"action_type" binding:"required"`
	EntityType  string `json:"entity_type" binding:"required"`
	EntityID    string `json:"entity_id" binding:"required"`
	Description string `json:"description"`
}

// RecordUserAction registra una acción del usuario y otorga XP
// @Summary Registrar acción del usuario
// @Description Registra una acción del usuario (ver insight, completar acción, etc.) y otorga XP
// @Tags gamification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param action body RecordActionRequest true "Datos de la acción"
// @Success 200 {object} usecases.ActionResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/gamification/actions [post]
func (h *Handler) RecordUserAction(c *gin.Context) {
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request RecordActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	params := usecases.RecordActionParams{
		UserID:      userID,
		ActionType:  request.ActionType,
		EntityType:  request.EntityType,
		EntityID:    request.EntityID,
		Description: request.Description,
	}

	result, err := h.gamificationUseCase.RecordUserAction(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registrando acción: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetActionTypes obtiene los tipos de acciones disponibles
// @Summary Obtener tipos de acciones
// @Description Obtiene la lista de tipos de acciones disponibles para gamificación
// @Tags gamification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ActionType
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/gamification/action-types [get]
func (h *Handler) GetActionTypes(c *gin.Context) {
	_, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	actionTypes := []ActionType{
		{
			Type:        "view_insight",
			Description: "Ver un insight de IA",
			BaseXP:      1,
		},
		{
			Type:        "understand_insight",
			Description: "Marcar insight como entendido",
			BaseXP:      3,
		},
		{
			Type:        "complete_action",
			Description: "Completar una acción recomendada",
			BaseXP:      10,
		},
		{
			Type:        "view_pattern",
			Description: "Ver un patrón de gastos",
			BaseXP:      2,
		},
		{
			Type:        "use_suggestion",
			Description: "Usar una sugerencia del sistema",
			BaseXP:      5,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"action_types": actionTypes,
		"total":        len(actionTypes),
	})
}

// GetLevels obtiene información sobre los niveles disponibles
// @Summary Obtener información de niveles
// @Description Obtiene la lista de niveles y XP requerido para cada uno
// @Tags gamification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} LevelInfo
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/gamification/levels [get]
func (h *Handler) GetLevels(c *gin.Context) {
	_, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	levels := []LevelInfo{
		{Level: 0, Name: "Novato", XPRequired: 0, XPToNext: 100},
		{Level: 1, Name: "Aprendiz", XPRequired: 100, XPToNext: 150},
		{Level: 2, Name: "Explorador", XPRequired: 250, XPToNext: 250},
		{Level: 3, Name: "Analista", XPRequired: 500, XPToNext: 500},
		{Level: 4, Name: "Estratega", XPRequired: 1000, XPToNext: 1000},
		{Level: 5, Name: "Experto", XPRequired: 2000, XPToNext: 2000},
		{Level: 6, Name: "Maestro", XPRequired: 4000, XPToNext: 4000},
		{Level: 7, Name: "Gurú", XPRequired: 8000, XPToNext: 8000},
		{Level: 8, Name: "Leyenda", XPRequired: 16000, XPToNext: 16000},
		{Level: 9, Name: "Magnate", XPRequired: 32000, XPToNext: 0},
	}

	c.JSON(http.StatusOK, gin.H{
		"levels": levels,
		"total":  len(levels),
	})
}
