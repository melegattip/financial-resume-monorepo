# 🏗️ Reglas de Arquitectura

## Proyecto
**Financial Resume Engine** desarrollado en **Go 1.23** implementando **Clean Architecture** con separación clara de capas.

## Estructura de Directorios

### 📁 Estructura Actual del Proyecto
```
financial-resume-engine/
├── cmd/                     # Puntos de entrada
│   └── api/                # Servidor HTTP principal
│       └── main.go
│
├── internal/               # Código privado
│   ├── core/              # 🧠 DOMINIO Y LÓGICA DE NEGOCIO
│   │   ├── domain/        # Entidades y value objects
│   │   ├── usecases/      # Casos de uso
│   │   ├── repository/    # Interfaces/contratos
│   │   ├── errors/        # Errores de dominio
│   │   ├── logs/          # Logging del core
│   │   ├── services/      # Servicios de dominio
│   │   └── web/           # Decodificadores web
│   │
│   ├── infrastructure/    # 🔌 ADAPTADORES EXTERNOS
│   │   ├── calculators/   # Calculadoras de analytics
│   │   ├── context/       # Contexto de infraestructura
│   │   ├── handler/       # Handlers de infraestructura
│   │   ├── http/          # Middlewares y respuestas HTTP
│   │   ├── logger/        # Sistema de logging
│   │   ├── repository/    # Acceso a datos
│   │   ├── router/        # Configuración de rutas
│   │   └── services/      # Servicios de infraestructura
│   │
│   ├── handlers/          # 🎯 CONTROLADORES HTTP
│   │   ├── analytics/     # Analytics de gastos
│   │   ├── categories/    # Gestión de categorías
│   │   ├── dashboard/     # Dashboard principal
│   │   ├── expenses/      # Gestión de gastos
│   │   └── incomes/       # Gestión de ingresos
│   │
│   ├── usecases/          # 📋 CASOS DE USO
│   │   ├── analytics/     # Analytics de transacciones
│   │   ├── categories/    # Casos de uso de categorías
│   │   ├── dashboard/     # Dashboard y reportes
│   │   ├── reports/       # Generación de reportes
│   │   └── transactions/  # Transacciones (expenses/incomes)
│   │
│   ├── adapters/          # 🔌 ADAPTADORES HTTP
│   │   └── http/          # Handlers HTTP adicionales
│   │
│   └── config/            # ⚙️ CONFIGURACIÓN
│       ├── configuration/ # Configuración general
│       ├── database.go    # Configuración de BD
│       └── environment/   # Variables de entorno
│
├── pkg/                   # Código público/compartido
│   ├── config/           # Configuración compartida
│   └── db/               # Conexión a base de datos
│
├── docs/                 # Documentación
├── .ai/                  # Reglas para IA
│   └── rules/           # Reglas de arquitectura
│
├── go.mod               # Dependencias
├── go.sum               # Checksums
├── Dockerfile           # Containerización
├── docker-compose.yml   # Orquestación local
└── README.md           # Documentación principal
```

## Principios Arquitectónicos

### Clean Architecture
- **Separación clara** entre core, infrastructure y handlers
- **Dependencias hacia adentro**: Las capas externas dependen de las internas
- **Independencia de frameworks**: El core no debe depender de herramientas externas

### Patrones de Diseño
- **Dependency Injection**: Desacoplar componentes
- **Repository Pattern**: Para acceso a datos externos
- **Interface Segregation**: Interfaces pequeñas y específicas

## Dependency Management

### Go Modules
- Mantener `go.mod` actualizado
- Usar versiones específicas de dependencias externas
- No agregar dependencias innecesarias

### Dependencias Principales
```go
// Dependencias principales del proyecto
github.com/gin-gonic/gin           // Framework HTTP
github.com/stretchr/testify        // Testing y mocks
github.com/swaggo/swag             // Documentación Swagger
go.uber.org/zap                    // Logging estructurado
``` 