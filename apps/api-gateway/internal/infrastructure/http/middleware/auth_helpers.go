package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
)

// AuthHelpers contiene funciones auxiliares para autenticación
type AuthHelpers struct{}

// NewAuthHelpers crea una nueva instancia de AuthHelpers
func NewAuthHelpers() *AuthHelpers {
	return &AuthHelpers{}
}

// GetUserIDFromContext extrae el user_id del contexto JWT
// Esta función centraliza la lógica que se repite en múltiples handlers
func (h *AuthHelpers) GetUserIDFromContext(c *gin.Context) (string, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", errors.NewUnauthorizedRequest("Usuario no autenticado")
	}

	// Convertir a string según el tipo
	switch v := userIDInterface.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case string:
		return v, nil
	default:
		return "", errors.NewBadRequest("Formato de user_id inválido")
	}
}

// GetUserIDFromContextOrEmpty extrae el user_id del contexto sin error si no existe
func (h *AuthHelpers) GetUserIDFromContextOrEmpty(c *gin.Context) string {
	userID, _ := h.GetUserIDFromContext(c)
	return userID
}

// RequireAuthentication middleware que requiere autenticación
func (h *AuthHelpers) RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := h.GetUserIDFromContext(c); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
