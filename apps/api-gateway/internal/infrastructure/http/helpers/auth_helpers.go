package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
)

// GetUserIDFromContext extrae el user_id del contexto JWT
// Esta función centraliza la lógica que estaba duplicada en múltiples handlers
func GetUserIDFromContext(c *gin.Context) (string, error) {
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

// GetUserIDFromContextAsUint extrae el user_id del contexto JWT como uint
// Para handlers que necesitan específicamente uint
func GetUserIDFromContextAsUint(c *gin.Context) (uint, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return 0, errors.NewUnauthorizedRequest("User ID not found in context")
	}

	// El user_id puede venir como string o uint dependiendo del middleware
	switch v := userIDStr.(type) {
	case uint:
		return v, nil
	case string:
		userID, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, errors.NewBadRequest("Invalid user ID format")
		}
		return uint(userID), nil
	default:
		return 0, errors.NewBadRequest("Invalid user ID type")
	}
}
