# Módulo AI — Niloft Financial

> Documento de referencia para implementar features, corregir bugs y entender la arquitectura del módulo de IA financiera.
> Última actualización: 2026-04-27

---

## Propósito

El módulo AI es el **asesor financiero inteligente** de Niloft. Su función es tomar datos financieros reales del usuario (ingresos, gastos, presupuestos, metas) y producir análisis personalizados, recomendaciones accionables y coaching mensual en español rioplatense.

No es un chatbot — es un sistema de análisis estructurado que siempre devuelve JSON con un esquema fijo, usando OpenAI como motor de razonamiento.

---

## Ubicación en el codebase

```
apps/monolith/internal/modules/ai/
├── module.go                    ← wiring, rutas, dependencias
├── handlers/
│   └── ai_handler.go            ← todos los HTTP handlers
├── service/
│   ├── openai_client.go         ← cliente OpenAI + mock mode
│   ├── analysis_service.go      ← health analysis, insights, coaching, education
│   ├── purchase_service.go      ← ¿puedo comprarlo? + alternativas
│   └── credit_service.go        ← plan crediticio + score
└── domain/
    ├── types.go                 ← FinancialAnalysisData y todos los tipos
    ├── coaching.go              ← MonthlyCoachingReport
    └── education.go             ← EducationCard

apps/frontend/src/
├── components/AIInsights.jsx    ← componente principal (orquesta todo)
├── pages/insights/tabs/
│   ├── MonthlyCoachingTab.jsx   ← tab "Reporte del Mes"
│   └── EducationTab.jsx         ← tab "Educación"
└── services/api.js              ← aiAPI (llamadas HTTP al backend)
```

---

## Endpoints

Todos bajo `/api/v1/ai`. **Sin autenticación propia** — el módulo AI no usa `authMW`/`permMW`. El `user_id` se extrae del JWT en el handler si no viene en el body.

| Método | Ruta | Handler | Descripción |
|--------|------|---------|-------------|
| POST | `/ai/health-analysis` | `AnalyzeFinancialHealth` | Análisis de salud financiera (score + insights) |
| POST | `/ai/insights` | `GenerateInsights` | 3 insights accionables del mes |
| POST | `/ai/can-i-buy` | `CanIBuy` | Decisión de compra |
| POST | `/ai/alternatives` | `SuggestAlternatives` | Alternativas más baratas a una compra |
| POST | `/ai/credit-plan` | `GenerateCreditPlan` | Plan de mejora crediticia |
| POST | `/ai/credit-score` | `CalculateCreditScore` | Score crediticio estimado |
| POST | `/ai/monthly-coaching` | `HandleMonthlyCoaching` | Reporte mensual de coaching |
| POST | `/ai/education-cards` | `HandleEducationCards` | 3 tarjetas educativas personalizadas |

> **Nota**: `GET /insights/financial-health` es del módulo **analytics**, no del AI. Calcula el score algorítmicamente sin OpenAI.

---

## Inicialización del módulo

```go
// module.go — firma de New()
func New(
    db *gorm.DB,
    logger zerolog.Logger,
    cfg *config.AppConfig,
    eventBus sharedports.EventBus,
    emailService sharedemail.EmailService,
) *Module
```

Diferencias respecto al patrón estándar de módulos:
- **No recibe** `authMW` ni `permMW` — sus rutas son públicas dentro del API
- **Recibe** `emailService` en lugar de los middlewares — lo usa para enviar el reporte mensual por email
- Si `OPENAI_API_KEY` está vacía → activa **mock mode** automáticamente (log warning al iniciar)

---

## OpenAI Client

**Archivo**: `service/openai_client.go`

| Config | Valor |
|--------|-------|
| Modelo | `gpt-4.1` |
| Temperature | `0.3` (respuestas determinísticas) |
| Max Tokens | `2000` |
| Timeout HTTP | `30 segundos` |

### Mock mode

Cuando `OPENAI_API_KEY` está vacía, `useMock = true`. El método `getMockResponse(userPrompt)` detecta keywords en el prompt para retornar JSON hardcodeado realista:

| Keyword detectada | Mock devuelto |
|-------------------|---------------|
| `"tarjetas educativas"` | 3 education cards |
| `"reporte de coaching"` | coaching report (positivo, 2 wins, 3 actions) |
| `"score crediticio"` | `{"score": 720}` |
| `"plan de mejora crediticia"` | credit plan 12 meses |
| `"can_buy"` o `"compra"` | purchase decision |
| `"insights"` | 3 insights |
| default | health analysis score 750 |

---

## Tipo central: FinancialAnalysisData

Todos los endpoints reciben este struct (o una variante). **El frontend es responsable de armarlo** antes de llamar al backend.

```go
type FinancialAnalysisData struct {
    UserID             string                  `json:"user_id"`
    TotalIncome        float64                 `json:"total_income"`
    TotalExpenses      float64                 `json:"total_expenses"`
    SavingsRate        float64                 `json:"savings_rate"`        // 0-1
    ExpensesByCategory map[string]float64      `json:"expenses_by_category"`
    IncomeStability    float64                 `json:"income_stability"`    // 0-1
    FinancialScore     int                     `json:"financial_score"`     // 0-1000
    Period             string                  `json:"period"`              // "YYYY-MM"
    SavingsGoals       []SavingsGoalInfo       `json:"savings_goals,omitempty"`
    BudgetsSummary     *BudgetsSummaryInfo     `json:"budgets_summary,omitempty"`
    BehaviorProfile    *BehaviorProfileContext `json:"behavior_profile,omitempty"`
}
```

### Cómo lo arma el frontend

El frontend hace llamadas paralelas a otros endpoints y construye `financialData`:

```javascript
const [catRes, incomeRes, budgetRes, goalsRes, profileRes] = await Promise.allSettled([
  analyticsAPI.categories(periodParam),   // GET /analytics/categories?from=...&to=...
  analyticsAPI.incomes(periodParam),      // GET /analytics/incomes?from=...&to=...
  budgetsAPI.getDashboard(),              // GET /budgets/dashboard (sin filtro de período)
  savingsGoalsAPI.list({ status: 'active' }),  // GET /savings-goals?status=active
  gamAPI.getBehaviorProfile(),            // GET /gamification/behavior-profile
]);
```

**Paths de respuesta críticos** (errores frecuentes):

| API | Respuesta JSON | Path correcto en JS |
|-----|---------------|---------------------|
| `analytics/categories` | `{ data: [...], total: N }` | `res.value.data?.data` |
| `analytics/incomes` | `{ total_amount, count, ... }` | `res.value.data?.total_amount` |
| `budgets/dashboard` | `{ summary: {...}, budgets: [...] }` | `res.value.data?.summary` |
| `savings-goals` | `{ data: [...], total: N }` | `res.value.data?.data` ← **NO `.goals`** |
| `gamification/behavior-profile` | BehaviorProfile directo | `res.value` |

---

## Caches en handlers

Los caches son **in-memory** (maps en el handler). Se pierden al reiniciar el servidor.

| Feature | Clave de cache | TTL efectivo |
|---------|---------------|--------------|
| Insights | `{userID}_{YYYY-MM-DD}` | Diario |
| Monthly Coaching | `{userID}_{YYYY-MM}` | Mensual |
| Education Cards | `{userID}_{YYYY-WW}` (ISO week) | Semanal |
| Email enviado | `{userID}_{YYYY-MM}` | Mensual |

El health analysis **no tiene cache en el módulo AI** — se recalcula siempre. El endpoint `GET /insights/financial-health` (analytics) tiene su propia lógica sin cache.

---

## Análisis de salud financiera (AnalyzeFinancialHealth)

**Endpoint**: `POST /ai/health-analysis`

**Request**: `FinancialAnalysisData`

**Response**:
```json
{
  "success": true,
  "data": {
    "score": 750,
    "level": "Bueno",
    "message": "Tu situación financiera...",
    "insights": [ ...3 AIInsight... ],
    "generated_at": "..."
  }
}
```

**Niveles**: `"Excelente"` | `"Bueno"` | `"Regular"` | `"Mejorable"`

**Fallback** en parse error: `{score: 500, level: "Regular", message: "Análisis completado..."}`

---

## Insights (GenerateInsights)

**Endpoint**: `POST /ai/insights`

**Response**:
```json
{
  "success": true,
  "data": {
    "insights": [ ...exactamente 3 AIInsight... ],
    "generated_at": "..."
  }
}
```

**Tipo AIInsight**:
```go
type AIInsight struct {
    Title       string // max 60 chars
    Description string // max 200 chars, con $montos y %
    Impact      string // "Positivo" | "Negativo" | "Neutro"
    Score       int    // 0-100
    ActionType  string // "maintain" | "improve" | "optimize" | "alert" | "invest"
    Category    string // "savings" | "expenses" | "income" | "debt" | "investment" | "budget" | "goals"
    NextAction  string // tarea concreta para esta semana, max 120 chars
}
```

**Personalización por nivel de sofisticación** (`detectSophistication`):

| Nivel | Condición | Tono del prompt |
|-------|-----------|-----------------|
| `"BÁSICO"` | discipline_score < 30 y sin presupuestos | Enfocado en crear primer presupuesto |
| `"AVANZADO"` | discipline_score >= 70 | Omite lo básico, optimización |
| `"EJECUTOR"` | ai_recommendations_applied >= 3 | Refuerza el hábito de acción |
| `"INTERMEDIO"` | default | Balance educación + acción |

---

## Coaching mensual (HandleMonthlyCoaching)

**Endpoint**: `POST /ai/monthly-coaching`

**Request**:
```json
{
  "financial_data": { ...FinancialAnalysisData... },
  "previous_month": "2026-04"
}
```

**Response**:
```json
{
  "report": {
    "month": "2026-04",
    "sentiment": "neutral",
    "summary": "...",
    "wins": [ {"title": "...", "description": "..."} ],
    "improvements": [ {"title": "...", "description": "..."} ],
    "actions": [ {"title": "...", "detail": "...", "deep_link": "/budgets"} ],
    "behavior_note": "...",
    "generated_at": "..."
  },
  "cached": false
}
```

**Restricciones del AI**: `wins`: 2-3 | `improvements`: 2-3 | `actions`: exactamente 3

**Deep links disponibles** para `actions.deep_link`:
`/dashboard`, `/expenses`, `/incomes`, `/budgets`, `/savings-goals`, `/recurring-transactions`, `/categories`, `/reports`, `/insights`

### Lógica clave en buildMonthlyCoachingPrompt

El prompt separa los egresos en dos grupos antes de enviárselos al AI:

```go
// isProductiveCoachingCategory detecta categorías de inversión/ahorro
keywords: "invers", "ahorro", "seguro", "educac", "retiro", "pension", "fondo",
          "activo", "propiedad", "inmueble", "capital", "emerg", "patrimonio",
          "cripto", "bitcoin", "etf", "accion", "bono", "plazo fijo"

consumptionTotal = totalExpenses - productiveTotal
netSavingsRate   = (income - consumptionTotal) / income * 100

// El balance se etiqueta explícitamente para evitar alucinaciones del AI:
balanceLabel = "SUPERÁVIT de $X" o "DÉFICIT de $X"
```

**Reglas del system prompt** (críticas para la coherencia):
- Principio 7: Los egresos en inversión/activos son **wins**, no problemas
- Principio 8: Las secciones `CUMPLIMIENTO DE PRESUPUESTOS` y `METAS DE AHORRO ACTIVAS` tienen datos REALES — prioridad sobre contadores del perfil conductual
- Principio 9: El campo `Balance neto del mes` es definitivo — el AI no puede invertirlo ni recalcularlo

### Email post-coaching

Al generar un reporte por primera vez (no cacheado), el handler dispara un envío de email **async** (goroutine separada). No bloquea la respuesta. Fallo en email = no crash.

---

## Tarjetas educativas (HandleEducationCards)

**Endpoint**: `POST /ai/education-cards`

**Request**: `{ "financial_data": { ...FinancialAnalysisData... } }`

**Response**:
```json
{
  "cards": [ ...exactamente 3 EducationCard... ],
  "generated_at": "...",
  "cached": false
}
```

**Tipo EducationCard**:
```go
type EducationCard struct {
    Topic      string // "emergencia"|"presupuesto"|"deuda"|"ahorro"|"inversión"|"impuestos"
    Title      string // max 60 chars
    Summary    string // 2-3 oraciones personalizadas
    KeyConcept string // frase memorable, max 80 chars
    CTA        string // label del botón, max 35 chars
    DeepLink   string // e.g., "/savings-goals"
    Difficulty string // "básico"|"intermedio"|"avanzado"
}
```

**Dificultad según sofisticación del usuario**:
- BÁSICO → básico/intermedio
- AVANZADO/EJECUTOR → avanzado

---

## ¿Puedo comprarlo? (CanIBuy / SuggestAlternatives)

**Endpoint**: `POST /ai/can-i-buy`

**Request**: `PurchaseAnalysisRequest`
```go
type PurchaseAnalysisRequest struct {
    UserID               string
    ItemName             string
    Amount               float64
    Description          string
    PaymentTypes         []string  // ["contado", "cuotas", "ahorro"]
    IsNecessary          bool
    UserFinancialProfile UserFinancialProfile
    SavingsGoal          *SavingsGoalInfo
}
```

**Response**: `PurchaseDecision`
```go
type PurchaseDecision struct {
    CanBuy       bool
    Confidence   float64   // 0-1
    Reasoning    string
    Alternatives []string
    ImpactScore  int       // 1-100
    GeneratedAt  time.Time
}
```

**Fallback** en parse error: `{can_buy: false, confidence: 0.5, impact_score: 50}`

**Endpoint**: `POST /ai/alternatives` — mismo request, devuelve `[]Alternative`:
```go
type Alternative struct {
    Name        string
    Description string
    Savings     float64
    Feasibility string // "alta"|"media"|"baja"
}
```

---

## Score y plan crediticio

**Endpoint**: `POST /ai/credit-score`
- Request: `FinancialAnalysisData`
- Response: `{score: int}` (1-1000)
- Fallback algorítmico si OpenAI falla: base 500 + bonus por savings_rate + estabilidad + capacidad de pago

**Endpoint**: `POST /ai/credit-plan`
- Response: `CreditPlan { current_score, target_score, timeline_months (3-24), actions, key_metrics }`
- Actions: `{ title, description, priority (alta|media|baja), timeline, impact (10-50), difficulty (fácil|media|difícil) }`

---

## Frontend: componente AIInsights.jsx

**Ruta**: `apps/frontend/src/components/AIInsights.jsx`

Es el componente orquestador principal. Renderiza:
1. **Score de salud financiera** (con desglose expandible: consumo vs inversión, 4 dimensiones)
2. **Tabs**: Reporte del Mes | Educación | ¿Puedo comprarlo?

### Flujo de carga

```
mount
  └── loadDashboardData()   → GET /dashboard (balance, income, expenses del mes)
  └── loadSavingsGoals()    → GET /savings-goals?status=active
  └── loadAIInsights()
        ├── analyticsAPI.categories()   → expenses_by_category + total_expenses
        ├── analyticsAPI.incomes()      → total_income + income_stability
        ├── budgetsAPI.getDashboard()   → budgets_summary
        ├── savingsGoalsAPI.list()      → savings_goals
        └── gamificationAPI.getBehaviorProfile() → behavior_profile
            └── aiAPI.getInsights(financialData)
                └── loadHealthScore(behaviorProfile)  → GET /insights/financial-health
```

### Health Score Display

Muestra datos del endpoint `GET /insights/financial-health` (módulo analytics, no AI):
- `score`: 0-1000
- `current_month_incomes`, `current_month_expenses`, `current_month_balance`
- `productive_expenses`, `consumption_expenses`
- `cash_flow_score` (40%), `planning_score` (30%), `consistency_score` (20%), `engagement_score` (10%)
- `savings_rate`: tasa real = (balance + inversiones) / ingresos

### MonthlyCoachingTab

- Usa `usePeriod()` para saber qué mes analizar (no hardcodea "mes anterior")
- Fallback: `subMonths(new Date(), 1)` si no hay mes seleccionado
- Llama `analyticsAPI.categories(periodParam)` y `analyticsAPI.incomes(periodParam)` con el período seleccionado
- Los botones de "Ir →" en las actions usan `navigate(action.deep_link)` (React Router)

### Purchase Analysis

- `paymentTypes` es **array** (multi-selección): `["contado"]`, `["cuotas", "ahorro"]`, etc.
- Filtra savings goals relevantes por nombre/descripción antes de enviar al AI
- El form valida: item_name, amount, y al menos 1 paymentType

---

## Integración con Gamificación

El módulo AI registra acciones de gamificación en el frontend:

| Acción | Cuándo |
|--------|--------|
| `'read_education_card'` | Click en CTA de tarjeta educativa |
| `'complete_monthly_review'` | Al generar el reporte mensual |
| `'view_ai_insight'` | Al ver un insight (deduplicado por sesión en sessionStorage) |
| `'understand_ai_insight'` | Al marcar insight como entendido (backend verifica duplicados) |
| `'apply_ai_recommendation'` | Al aplicar una recomendación |

**Importante**: La deduplicación real vive en el backend. El frontend usa sessionStorage solo para UX (evitar clicks dobles).

---

## Distinción producivo vs consumo

Este es el concepto más crítico del módulo. **Siempre distinguir** estos dos tipos de egresos:

**Consumo** (penaliza el score): gastos corrientes, viajes, servicios, alimentación, entretenimiento

**Construcción patrimonial** (es un WIN): inversiones, activos, ahorro, seguros, educación, retiro, fondo de emergencia, cripto, ETFs, bonos, propiedades

Keywords para detectar categorías productivas (sincronizadas entre `analytics/handlers` y `ai/service`):
```
"invers", "ahorro", "seguro", "educac", "retiro", "pension", "fondo",
"activo", "propiedad", "inmueble", "capital", "emerg", "patrimonio",
"cripto", "bitcoin", "etf", "accion", "bono", "plazo fijo"
```

La **tasa de ahorro neto** = `(ingresos - consumo) / ingresos` — excluye inversiones.
La **tasa de ahorro real** = `(balance + inversiones) / ingresos` — incluye todo.

---

## Idioma y estilo

- Todo el AI responde en **español rioplatense** ("vos", "hacé", "usá")
- Moneda: pesos argentinos (ARS), formato `$1.200.000,00`
- Los prompts incluyen contexto latinoamericano explícito
- Tono del coaching: **cálido y empático** (mentor de confianza, no banco)
- Tono de los insights: **directo y accionable** (con $montos y % reales)

---

## Bugs conocidos y workarounds activos

| Bug | Descripción | Estado / Workaround |
|-----|-------------|---------------------|
| BUG-007 | AI no distinguía inversiones de consumo en savings_rate | **Corregido** — `buildMonthlyCoachingPrompt` separa productivo/consumo |
| BUG-001 | Contadores de gamificación (budgets_created, etc.) siempre en 0 | **Sin resolver** — el system prompt tiene Principio 8 para que el AI priorice datos reales de API sobre contadores |
| — | Savings goals path incorrecto (`data?.data?.goals`) | **Corregido** en MonthlyCoachingTab y AIInsights — usar `data?.data` |
| — | `financial_score = 0` hardcodeado en coaching | **Corregido** — campo omitido del payload |
| — | AI alucinaba déficit con balance positivo | **Corregido** — balance se etiqueta SUPERÁVIT/DÉFICIT explícitamente |

---

## Checklist para implementar un nuevo feature AI

1. **Definir el domain type** en `domain/types.go` (request + response structs)
2. **Crear el método de servicio** en el service correspondiente:
   - System prompt en español
   - User prompt construido desde `FinancialAnalysisData`
   - `cleanJSONResponse()` antes de `json.Unmarshal`
   - Fallback seguro si el parse falla
3. **Agregar el mock** en `getMockResponse()` del OpenAIClient (keyword detection)
4. **Crear el handler** en `ai_handler.go`:
   - Extraer `user_id` del JWT si no viene en body
   - Decidir si necesita cache (¿tiene sentido repetir el análisis?)
   - Si tiene cache: map con key `{userID}_{período}`
5. **Registrar la ruta** en `module.go`
6. **Agregar el método** en `aiAPI` de `api.js` (timeout 60s para llamadas AI)
7. **Construir el `financialData`** en el componente frontend siguiendo el patrón establecido
8. **Verificar** que `go build ./...` compila sin errores
