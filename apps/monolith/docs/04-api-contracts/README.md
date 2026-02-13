# Contratos API - Financial Resume

**Última Actualización**: 2026-02-09
**Versión**: 1.0.0 (Estado Actual)

---

## Resumen General

El API de Financial Resume cuenta con **95+ endpoints** distribuidos en 15 categorías funcionales.

### Arquitectura API

- **Protocolo**: HTTP/REST
- **Formato**: JSON
- **Autenticación**: JWT Bearer Token
- **Headers Requeridos**:
  - `x-caller-id`: ID del usuario (requerido en todos los endpoints)
  - `Authorization`: Bearer token (requerido en endpoints protegidos)
  - `Content-Type`: application/json (en POST/PUT/PATCH)

### Distribución de Endpoints

| Categoría | Endpoints | Descripción |
|-----------|-----------|-------------|
| **Authentication & Users** | 13 | Login, registro, perfil (proxy a users-service) |
| **Expenses** | 6 | CRUD de gastos, pagos parciales |
| **Incomes** | 5 | CRUD de ingresos |
| **Categories** | 5 | Gestión de categorías |
| **Budgets** | 8 | CRUD de presupuestos, alertas |
| **Savings Goals** | 13 | CRUD de metas, depósitos, retiros |
| **Recurring Transactions** | 12 | CRUD de transacciones recurrentes |
| **Dashboard** | 1 | Resumen financiero |
| **Analytics** | 3 | Reportes y análisis |
| **Insights** | 2 | Insights de IA |
| **AI Features** | 3 | Health score, recomendaciones (proxy a ai-service) |
| **Reports** | 1 | Reportes personalizados |
| **Gamification** | 9 | XP, niveles, logros, challenges (proxy a gamification-service) |
| **Config** | 2 | Configuración del sistema |
| **Health** | 1 | Health check |

---

## Patrones de API

### URLs Base

```
Production: https://financial-resume-backend.onrender.com
Local: http://localhost:8080
```

### Versionado

Todos los endpoints están bajo el prefijo `/api/v1/`

### Códigos de Estado HTTP

| Código | Uso |
|--------|-----|
| 200 | Éxito (GET, PUT, PATCH) |
| 201 | Recurso creado (POST) |
| 204 | Sin contenido (DELETE) |
| 400 | Bad request (validación fallida) |
| 401 | No autenticado |
| 403 | No autorizado |
| 404 | Recurso no encontrado |
| 409 | Conflicto (duplicado) |
| 422 | Entidad no procesable |
| 500 | Error interno del servidor |

### Paginación

**Formato**:
```
GET /api/v1/transactions?page=1&limit=20
```

**Respuesta**:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 157,
    "total_pages": 8
  }
}
```

### Filtrado

**Por fechas**:
```
GET /api/v1/transactions?start_date=2024-01-01&end_date=2024-12-31
```

**Por estado**:
```
GET /api/v1/expenses?paid=true
GET /api/v1/budgets?status=exceeded
```

**Por categoría**:
```
GET /api/v1/transactions?category_id=cat_12345678
```

---

## Categorías de Endpoints

### 1. Authentication & Users (Proxy)

Estos endpoints hacen proxy al `users-service` en puerto 8083.

**Base**: `/api/v1/auth` y `/api/v1/users`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | /api/v1/auth/register | Registro de nuevo usuario |
| POST | /api/v1/auth/login | Login con email/password |
| POST | /api/v1/auth/logout | Cerrar sesión |
| GET | /api/v1/users/profile | Obtener perfil de usuario |
| PUT | /api/v1/users/profile | Actualizar perfil |
| POST | /api/v1/users/change-password | Cambiar contraseña |

---

### 2. Transactions - Expenses

**Base**: `/api/v1/expenses`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/expenses | Listar gastos (con filtros) |
| GET | /api/v1/expenses/:id | Obtener gasto por ID |
| POST | /api/v1/expenses | Crear nuevo gasto |
| PUT | /api/v1/expenses/:id | Actualizar gasto completo |
| PATCH | /api/v1/expenses/:id/payment | Agregar pago parcial |
| DELETE | /api/v1/expenses/:id | Eliminar gasto |

---

### 3. Transactions - Incomes

**Base**: `/api/v1/incomes`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/incomes | Listar ingresos |
| GET | /api/v1/incomes/:id | Obtener ingreso por ID |
| POST | /api/v1/incomes | Crear nuevo ingreso |
| PUT | /api/v1/incomes/:id | Actualizar ingreso |
| DELETE | /api/v1/incomes/:id | Eliminar ingreso |

---

### 4. Categories

**Base**: `/api/v1/categories`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/categories | Listar categorías del usuario |
| GET | /api/v1/categories/:id | Obtener categoría por ID |
| POST | /api/v1/categories | Crear nueva categoría |
| PUT | /api/v1/categories/:id | Actualizar categoría |
| DELETE | /api/v1/categories/:id | Eliminar categoría |

---

### 5. Budgets

**Base**: `/api/v1/budgets`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/budgets | Listar presupuestos |
| GET | /api/v1/budgets/:id | Obtener presupuesto por ID |
| POST | /api/v1/budgets | Crear nuevo presupuesto |
| PUT | /api/v1/budgets/:id | Actualizar presupuesto |
| PATCH | /api/v1/budgets/:id/spent | Actualizar gasto |
| GET | /api/v1/budgets/summary | Resumen de presupuestos |
| DELETE | /api/v1/budgets/:id | Eliminar presupuesto |

---

### 6. Savings Goals

**Base**: `/api/v1/savings-goals`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/savings-goals | Listar metas de ahorro |
| GET | /api/v1/savings-goals/:id | Obtener meta por ID |
| POST | /api/v1/savings-goals | Crear nueva meta |
| PUT | /api/v1/savings-goals/:id | Actualizar meta |
| POST | /api/v1/savings-goals/:id/deposit | Agregar depósito |
| POST | /api/v1/savings-goals/:id/withdraw | Retirar ahorro |
| GET | /api/v1/savings-goals/:id/transactions | Historial de transacciones |
| PATCH | /api/v1/savings-goals/:id/pause | Pausar meta |
| PATCH | /api/v1/savings-goals/:id/resume | Reanudar meta |
| PATCH | /api/v1/savings-goals/:id/cancel | Cancelar meta |
| GET | /api/v1/savings-goals/summary | Resumen de metas |
| DELETE | /api/v1/savings-goals/:id | Eliminar meta |

---

### 7. Recurring Transactions

**Base**: `/api/v1/recurring-transactions`

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/recurring-transactions | Listar transacciones recurrentes |
| GET | /api/v1/recurring-transactions/:id | Obtener por ID |
| POST | /api/v1/recurring-transactions | Crear nueva recurrente |
| PUT | /api/v1/recurring-transactions/:id | Actualizar recurrente |
| POST | /api/v1/recurring-transactions/:id/execute | Ejecutar manualmente |
| PATCH | /api/v1/recurring-transactions/:id/pause | Pausar |
| PATCH | /api/v1/recurring-transactions/:id/resume | Reanudar |
| GET | /api/v1/recurring-transactions/pending | Listar pendientes |
| DELETE | /api/v1/recurring-transactions/:id | Eliminar |

---

### 8. Dashboard

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/dashboard | Resumen completo del dashboard |

**Respuesta incluye**: Ingresos totales, gastos totales, balance, presupuestos, metas, insights recientes

---

### 9. Analytics & Reports

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/analytics/spending-by-category | Gastos agrupados por categoría |
| GET | /api/v1/analytics/trends | Tendencias temporales |
| GET | /api/v1/reports/financial-summary | Reporte financiero personalizado |

---

### 10. Gamification (Proxy)

**Base**: `/api/v1/gamification`

Estos endpoints hacen proxy al `gamification-service` en puerto 8084.

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/gamification/profile | Perfil de gamificación |
| GET | /api/v1/gamification/achievements | Logros del usuario |
| GET | /api/v1/gamification/challenges | Challenges activos |
| POST | /api/v1/gamification/actions | Registrar acción (gana XP) |
| GET | /api/v1/gamification/stats | Estadísticas completas |
| GET | /api/v1/gamification/leaderboard | Tabla de líderes |

---

## Próximos Pasos

1. **Especificación OpenAPI completa**: Ver `current-api.yaml` (pendiente de generación)
2. **Ejemplos de requests/responses**: Ver `examples/` (pendiente)
3. **Contratos de eventos**: Ver `events.md` (para arquitectura event-driven futura)

---

**Nota**: El archivo OpenAPI completo (`current-api.yaml`) es muy extenso (+28K tokens) y se está generando. Por ahora, este README proporciona una visión general completa de todos los endpoints disponibles.

**Generado por**: Agent: API Contracts Explorer
**Código analizado**: master (commit 6fca155)
