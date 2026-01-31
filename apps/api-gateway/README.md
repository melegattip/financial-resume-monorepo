# Financial Resume Engine (Core + API Gateway)

Este servicio es el núcleo financiero y el API Gateway del ecosistema. Orquesta los módulos de dominio (transacciones, categorías, presupuestos, metas, recurrentes) y actúa como proxy hacia microservicios independientes (gamificación, IA, usuarios).

## Estado actual

- Arquitectura: Clean Architecture con capas `handlers → usecases → repository/domain`.
- Router: `internal/infrastructure/router/router.go` registra todas las rutas documentadas en `docs/00_ESTADO_ACTUAL_CONSOLIDADO.md`.
- Seguridad: JWT mediante proxy a users-service; CORS configurable; rate limiting y métricas in-memory.
- DB: PostgreSQL principal (`scripts/init.sql`).
- Desac acoplamiento: integración con gamificación/IA/usuarios vía HTTP proxy; la lógica financiera no depende de respuestas externas.

## Endpoints principales

Ver lista completa en `docs/00_ESTADO_ACTUAL_CONSOLIDADO.md` (sección API REST). Coincide 1:1 con `router.go` y con `frontend/src/services/api.js`.

## Cambios recientes relevantes

- Metas de ahorro: historial de movimientos (`GET /savings-goals/:id/transactions`), payloads parciales en `PUT`, validación de fecha relajada en updates.
- Gamificación (vía proxy): idempotencia diaria de `view_dashboard`; `view_insight` sin XP.
- Frontend: barras de progreso y sincronización detalle↔movimientos.

## Cómo ejecutar

```bash
go build ./...
go run cmd/api/main.go
```

Variables críticas: ver `env.example`.

## Tests

- `internal/core/usecases/savings_goal_scenarios_test.go`: flujo e2e en memoria para metas.

## Roadmap corto

- Event-driven para gamificación (reemplazar helper HTTP por cola de eventos).
- Tests de contrato router↔frontend (generación automática de listado de rutas y verificación en CI).

# Financial Resume Engine

Motor de resumen financiero que permite gestionar transacciones, categorías y generar reportes.

## Requisitos Previos

- Go 1.23 o superior
- PostgreSQL 15 o superior
- Sistema operativo compatible:
  - Windows 10/11 Pro, Enterprise o Education (64-bit)
  - Linux (Ubuntu 20.04 LTS o superior)
  - macOS 10.15 o superior

## Instalación

1. Clonar el repositorio:
```bash
git clone https://github.com/melegattip/financial-resume-engine.git
cd financial-resume-engine
```

2. Configurar variables de entorno:
```bash
# Windows (PowerShell)
Copy-Item .env.example .env

# Linux/macOS
cp .env.example .env

# Editar .env con tus credenciales seguras
```

## Ejecución con Docker Compose

El proyecto incluye una configuración de Docker Compose que facilita la ejecución tanto de la API como de la base de datos PostgreSQL:

```bash
# Iniciar todos los servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener todos los servicios
docker-compose down

# Detener y eliminar volúmenes
docker-compose down -v
```

La configuración incluye:
- PostgreSQL 15 en el puerto 5432
- Base de datos: financial_resume
- Usuario: postgres
- Contraseña: postgres
- Volumen persistente para los datos
- Healthcheck para asegurar que la base de datos esté lista antes de iniciar la API

## Ejecución

1. Asegúrate de tener PostgreSQL instalado y corriendo:

### Windows
- Descarga e instala PostgreSQL desde [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)
- Inicia el servicio de PostgreSQL desde Servicios de Windows

### Linux (Ubuntu)
```bash
# Instalar PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Iniciar el servicio
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### macOS
```bash
# Instalar PostgreSQL usando Homebrew
brew install postgresql@15

# Iniciar el servicio
brew services start postgresql@15
```

2. Crear la base de datos y el usuario:
```sql
-- Conectarse a PostgreSQL
psql -U postgres

-- Crear la base de datos y el usuario (ajusta según tu .env)
CREATE DATABASE financial_resume;
CREATE USER financial_user WITH PASSWORD 'tu_contraseña';
GRANT ALL PRIVILEGES ON DATABASE financial_resume TO financial_user;
```

3. Ejecutar la API:
```bash
# Instalar dependencias
go mod download

# Ejecutar la aplicación
go run cmd/app/main.go
```

## Documentación API (Swagger)

La documentación completa de la API está disponible a través de Swagger UI:

- URL: http://localhost:8080/swagger/index.html

En la interfaz de Swagger podrás:
1. Ver todas las operaciones disponibles organizadas por categorías
2. Probar los endpoints directamente desde la interfaz
3. Ver los modelos de datos y parámetros requeridos
4. Autenticarte usando el botón "Authorize" para probar endpoints protegidos

Parámetros comunes para todas las operaciones:
- Header `x-caller-id`: Identificador del usuario que realiza la llamada
- Header `Authorization`: Token Bearer para autenticación (cuando sea requerido)

## Endpoints Disponibles

### Transacciones

#### Crear Transacción
```bash
# Windows (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/transactions" `
  -Method Post `
  -Headers @{
    "Content-Type" = "application/json"
    "x-caller-id" = "test-user"
  } `
  -Body '{
    "description": "Test Transaction",
    "amount": 100.00,
    "paid": false,
    "expiry_date": "2024-12-31",
    "category": "Test"
  }'

# Linux/macOS
curl --location 'http://localhost:8080/api/v1/transactions' \
--header 'Content-Type: application/json' \
--header 'x-caller-id: test-user' \
--data '{
    "description": "Test Transaction",
    "amount": 100.00,
    "paid": false,
    "expiry_date": "2024-12-31",
    "category": "Test"
}'
```

#### Listar Transacciones
```bash
# Windows (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/transactions" `
  -Method Get `
  -Headers @{
    "x-caller-id" = "test-user"
  }

# Linux/macOS
curl --location 'http://localhost:8080/api/v1/transactions' \
--header 'x-caller-id: test-user'
```

### Categorías

#### Crear Categoría
```bash
# Windows (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/categories" `
  -Method Post `
  -Headers @{
    "Content-Type" = "application/json"
    "x-caller-id" = "test-user"
  } `
  -Body '{
    "name": "Test Category",
    "description": "Test Description"
  }'

# Linux/macOS
curl --location 'http://localhost:8080/api/v1/categories' \
--header 'Content-Type: application/json' \
--header 'x-caller-id: test-user' \
--data '{
    "name": "Test Category",
    "description": "Test Description"
}'
```

#### Listar Categorías
```bash
# Windows (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/categories" `
  -Method Get `
  -Headers @{
    "x-caller-id" = "test-user"
  }

# Linux/macOS
curl --location 'http://localhost:8080/api/v1/categories' \
--header 'x-caller-id: test-user'
```

### Reportes

#### Obtener Resumen por Categoría
```bash
# Windows (PowerShell)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/reports/category-summary" `
  -Method Get `
  -Headers @{
    "x-caller-id" = "test-user"
  }

# Linux/macOS
curl --location 'http://localhost:8080/api/v1/reports/category-summary' \
--header 'x-caller-id: test-user'
```

## Desarrollo

### Configuración de VS Code

El proyecto incluye configuración para debugging en VS Code:

1. Abrir el proyecto en VS Code
2. Presionar F5 para iniciar el debugging
3. Seleccionar "Debug Go"

### Estructura del Proyecto

```
.
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── categories/
│   ├── config/
│   ├── models/
│   ├── reports/
│   └── transactions/
├── .env.example
└── go.mod
```

## Notas Importantes

- El header `x-caller-id` es requerido en todas las peticiones
- Las transacciones están asociadas a un usuario específico
- Los reportes se generan por usuario
- Nunca compartas o comitees el archivo `.env` con tus credenciales reales

## Solución de Problemas

### Base de Datos

1. Si no puedes conectarte a la base de datos:
```bash
# Verificar que PostgreSQL está corriendo
# Windows
sc query postgresql

# Linux
sudo systemctl status postgresql

# macOS
brew services list | grep postgresql
```

2. Si necesitas reiniciar la base de datos:
```bash
# Windows
net stop postgresql
net start postgresql

# Linux
sudo systemctl restart postgresql

# macOS
brew services restart postgresql@15
```

3. Si necesitas verificar los logs de PostgreSQL:
```bash
# Windows
# Ver en el Visor de eventos de Windows

# Linux
sudo tail -f /var/log/postgresql/postgresql-15-main.log

# macOS
tail -f /usr/local/var/log/postgresql@15.log
```

# Modelos de datos

- **paid**: booleano. Indica si el gasto está pagado o pendiente. Este campo es relevante para egresos y se utiliza en la gestión y reporte de gastos.

## Ejemplo de respuesta de reporte financiero

```json
{
  "start_date": "2024-06-01",
  "end_date": "2024-06-30",
  "transactions": [
    {
      "id": "...",
      "user_id": "...",
      "amount": -100.00,
      "description": "Supermercado",
      "category_id": "cat1",
      "paid": true,
      "created_at": "2024-06-01T00:00:00Z",
      "updated_at": "2024-06-01T00:00:00Z",
      "percentage": -3.45
    }
    // ... más transacciones ...
  ],
  "total_income": 2900.00,
  "total_expenses": 2863.22,
  "net_balance": 36.78,
  "category_summary": [
    {
      "category_id": "cat1",
      "category_name": "Supermercado",
      "total_amount": 323.00,
      "percentage": 11.14
    },
    {
      "category_id": "cat2",
      "category_name": "Servicios",
      "total_amount": 880.00,
      "percentage": 30.68
    }
    // ... más categorías ...
  ]
}
```

- **percentage**: En cada transacción, indica el porcentaje que representa ese monto respecto al total de ingresos del período.

- **category_summary**: Array con el total y porcentaje de egresos agrupados por categoría. Útil para gráficos y análisis como en la planilla de ejemplo.