package services

import (
	"fmt"
	"strings"
)

// === UTILIDADES COMPARTIDAS PARA SERVICIOS IA ===

// GetTopExpenseCategories obtiene las top N categorías de gastos
func GetTopExpenseCategories(categories map[string]float64, n int) map[string]float64 {
	if len(categories) == 0 {
		return make(map[string]float64)
	}

	// Convertir a slice para ordenar
	type categoryAmount struct {
		name   string
		amount float64
	}

	var items []categoryAmount
	for name, amount := range categories {
		items = append(items, categoryAmount{name: name, amount: amount})
	}

	// Ordenar por monto descendente
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].amount < items[j].amount {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// Tomar top N
	result := make(map[string]float64)
	limit := n
	if len(items) < n {
		limit = len(items)
	}

	for i := 0; i < limit; i++ {
		result[items[i].name] = items[i].amount
	}

	return result
}

// FormatTopCategories formatea las categorías para el prompt
func FormatTopCategories(categories map[string]float64) string {
	if len(categories) == 0 {
		return "Sin datos de categorías"
	}

	var parts []string
	for name, amount := range categories {
		parts = append(parts, fmt.Sprintf("%s: $%.0f", name, amount))
	}

	return strings.Join(parts, ", ")
}

// GetImpactLevel determina el nivel de impacto basado en porcentaje
func GetImpactLevel(percentage float64) string {
	// ✅ Manejar valores negativos como impacto alto
	if percentage < 0 {
		return "high"
	}
	if percentage > 40 {
		return "high"
	} else if percentage > 25 {
		return "medium"
	}
	return "low"
}

// CalculateCategoryScore calcula score basado en porcentaje de categoría
func CalculateCategoryScore(percentage float64) int {
	if percentage > 50 {
		return 300
	} else if percentage > 40 {
		return 400
	} else if percentage > 30 {
		return 600
	} else if percentage > 20 {
		return 750
	}
	return 850
}
