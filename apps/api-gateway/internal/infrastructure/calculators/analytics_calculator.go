package calculators

import (
	"hash/fnv"

	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// AnalyticsCalculatorImpl implementa la interfaz AnalyticsCalculator
type AnalyticsCalculatorImpl struct{}

// NewAnalyticsCalculator crea una nueva instancia del calculador analítico
func NewAnalyticsCalculator() usecases.AnalyticsCalculator {
	return &AnalyticsCalculatorImpl{}
}

// CalculatePercentages calcula el porcentaje de amount respecto al total
func (a *AnalyticsCalculatorImpl) CalculatePercentages(amount, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (amount / total) * 100
}

// CalculateAverage calcula el promedio
func (a *AnalyticsCalculatorImpl) CalculateAverage(total float64, count int) float64 {
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// GenerateColorSeed genera una semilla numérica para colores basada en un identificador
func (a *AnalyticsCalculatorImpl) GenerateColorSeed(identifier string) int {
	hash := fnv.New32a()
	hash.Write([]byte(identifier))
	return int(hash.Sum32())
}
