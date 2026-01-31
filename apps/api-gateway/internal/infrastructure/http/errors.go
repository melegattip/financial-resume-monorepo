package http

import (
	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
)

// ErrorResponse representa la estructura de respuesta de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Error message"`
	Message string `json:"message" example:"Bad Request"`
}

// HandleError maneja los errores de la aplicación y los convierte en respuestas HTTP apropiadas
func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.BadRequest:
		BadRequest(c, e)
	case *errors.UnauthorizedRequest:
		Unauthorized(c, e)
	case *errors.ResourceNotFound:
		NotFound(c, e)
	case *errors.TooManyRequests:
		TooManyRequests(c, e)
	case *errors.ResourceAlreadyExists:
		Conflict(c, e)
	default:
		InternalServerError(c, e)
	}
}
