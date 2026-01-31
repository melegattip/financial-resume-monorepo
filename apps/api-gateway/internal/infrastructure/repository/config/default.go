package config

// DefaultPath contiene la ruta al archivo de configuración por defecto
const DefaultPath = "internal/infrastructure/repository/config/test.properties"

// Default contiene la configuración por defecto para pruebas
var Default = map[string]interface{}{
	"database": map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"user":     "postgres",
		"password": "postgres",
		"dbname":   "financial_resume",
	},
}
