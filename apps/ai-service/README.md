# 🤖 Financial AI Service

Microservicio especializado en análisis financiero con Inteligencia Artificial, extraído del servicio principal como parte de la **Fase 2** del refactor arquitectónico.

## 🎯 Propósito

Este microservicio maneja todas las funcionalidades relacionadas con IA:
- **Análisis de salud financiera** con insights personalizados
- **Decisiones de compra inteligentes** con análisis de viabilidad
- **Planes de mejora crediticia** con acciones específicas y priorizadas

## 🏗️ Arquitectura

Sigue los principios de **Clean Architecture** y **SOLID**:

```
financial-ai-service/
├── cmd/api/                    # Punto de entrada
├── internal/
│   ├── core/                   # Lógica de negocio
│   │   ├── ports/              # Interfaces (puertos)
│   │   └── usecases/           # Casos de uso
│   ├── adapters/               # Adaptadores
│   │   ├── openai/             # Cliente OpenAI
│   │   ├── cache/              # Cliente Redis
│   │   └── http/               # Handlers HTTP
│   └── infrastructure/         # Configuración
├── pkg/                        # Utilidades compartidas
└── Dockerfile                  # Containerización
```

## 🚀 Endpoints

### Análisis Financiero
- `POST /api/v1/ai/health-analysis` - Análisis de salud financiera
- `POST /api/v1/ai/insights` - Generación de insights personalizados

### Decisiones de Compra
- `POST /api/v1/ai/can-i-buy` - Análisis de viabilidad de compra
- `POST /api/v1/ai/alternatives` - Sugerencias de alternativas

### Análisis Crediticio
- `POST /api/v1/ai/credit-plan` - Plan de mejora crediticia
- `POST /api/v1/ai/credit-score` - Cálculo de score crediticio

### Health Check
- `GET /health` - Estado del servicio

## ⚙️ Configuración

### Variables de Entorno

```bash
# Servidor
PORT=8082
HOST=localhost

# OpenAI
OPENAI_API_KEY=sk-...
USE_AI_MOCK=true

# Redis (Cache)
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Cache TTL
CACHE_DEFAULT_TTL_MINUTES=30
CACHE_INSIGHTS_TTL_HOURS=20
```

## 🛠️ Desarrollo

### Prerrequisitos
- Go 1.21+
- Redis (opcional, usa mock por defecto)
- OpenAI API Key (opcional para desarrollo)

### Instalación

```bash
# Clonar y navegar al directorio
cd financial-ai-service

# Instalar dependencias
go mod tidy

# Ejecutar en modo desarrollo
USE_AI_MOCK=true go run cmd/api/main.go
```

### Docker

```bash
# Construir imagen
docker build -t financial-ai-service .

# Ejecutar contenedor
docker run -p 8082:8082 \
  -e USE_AI_MOCK=true \
  -e PORT=8082 \
  financial-ai-service
```

## 🧪 Pruebas

### Health Check
```bash
curl http://localhost:8082/health
```

### Análisis de Salud Financiera
```bash
curl -X POST http://localhost:8082/api/v1/ai/health-analysis \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test-user",
    "total_income": 5000000,
    "total_expenses": 3500000,
    "savings_rate": 0.3,
    "income_stability": 0.8,
    "financial_score": 750,
    "period": "monthly",
    "expenses_by_category": {
      "Alimentación": 1200000,
      "Transporte": 800000,
      "Vivienda": 1500000
    }
  }'
```

### Decisión de Compra
```bash
curl -X POST http://localhost:8082/api/v1/ai/can-i-buy \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test-user",
    "item_name": "MacBook Pro",
    "amount": 8000000,
    "description": "Laptop para trabajo",
    "is_necessary": true,
    "payment_types": ["contado"],
    "user_financial_profile": {
      "current_balance": 10000000,
      "monthly_income": 5000000,
      "monthly_expenses": 3500000,
      "savings_rate": 0.3,
      "income_stability": 0.8,
      "financial_discipline": 750
    }
  }'
```

## 📊 Características

### ✅ Implementado
- **Clean Architecture** con separación clara de responsabilidades
- **Cache inteligente** con Redis (mock incluido)
- **Fallbacks robustos** para manejo de errores
- **Logging estructurado** para debugging
- **Configuración flexible** vía variables de entorno
- **Health checks** para monitoreo
- **Dockerización** para deployment

### 🔄 Próximas Mejoras
- **Métricas y monitoring** con Prometheus
- **Circuit breakers** para resilencia
- **Rate limiting** para protección
- **Autenticación** JWT
- **Documentación Swagger** automática

## 🔗 Integración

Este microservicio está diseñado para integrarse con:
- **Financial Resume Engine** (servicio principal)
- **Financial Gamification Service** (sistema de gamificación)
- **Redis** (cache distribuido)
- **OpenAI API** (análisis con IA)

## 📈 Métricas Esperadas

- **Latencia**: < 2s para análisis complejos
- **Throughput**: 100+ requests/segundo
- **Cache Hit Rate**: > 70%
- **Disponibilidad**: 99.9%

## 🚨 Troubleshooting

### Errores Comunes

1. **Error de conexión OpenAI**
   ```bash
   export USE_AI_MOCK=true
   ```

2. **Puerto ocupado**
   ```bash
   export PORT=8083
   ```

3. **Cache no disponible**
   - El servicio funciona sin Redis usando mock

### Logs Importantes
```bash
# Verificar inicialización
🤖 Starting Financial AI Service...
✅ OpenAI client initialized with API key: sk-1234567890...
🎭 Redis client initialized in MOCK mode

# Verificar requests
🧠 Analyzing financial health for user: test-user
✅ Health analysis completed for user: test-user (Score: 750)
```

---

**Desarrollado como parte del refactor arquitectónico Fase 2**  
**Migrado desde financial-resume-engine** ✨ 