package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
)

// ParseIntQuery parsea un parámetro de query como entero
func ParseIntQuery(c *gin.Context, key string, required bool) (*int, error) {
	valueStr := c.Query(key)
	if valueStr == "" {
		if required {
			return nil, errors.NewBadRequest(key + " es requerido")
		}
		return nil, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		if key == "year" {
			return nil, errors.NewBadRequest("Año inválido")
		} else if key == "month" {
			return nil, errors.NewBadRequest("Mes inválido")
		}
		return nil, errors.NewBadRequest(key + " debe ser un número válido")
	}

	return &value, nil
}

// ParseRequiredParam parsea un parámetro de ruta requerido
func ParseRequiredParam(c *gin.Context, key string) (string, error) {
	value := c.Param(key)
	if value == "" {
		return "", errors.NewBadRequest(key + " es requerido")
	}
	return value, nil
}

// ValidateYearMonth valida parámetros de año y mes
func ValidateYearMonth(year, month *int) error {
	if year != nil && (*year < 2020 || *year > 2030) {
		return errors.NewBadRequest("Año debe estar entre 2020 y 2030")
	}

	if month != nil && (*month < 1 || *month > 12) {
		return errors.NewBadRequest("Mes debe estar entre 1 y 12")
	}

	return nil
}
