# 💻 Estándares de Código Go

## Herramientas Obligatorias

### Linters
```bash
# Linter principal
golangci-lint run

# Linter con configuración custom si existe
golangci-lint run -c .golangci.yml

# Formateo estándar
gofmt -w .
go vet ./...
```

### Formateo
```bash
# Herramientas estándar
gofmt
go vet
go mod tidy
```

## Buenas Prácticas de Código

### Manejo de Errores
- **Siempre explícito**: Nunca ignorar errores
- **Errores tipados**: Usar constructores específicos del package `internal/core/errors`
- **No usar fmt.Errorf**: Preferir errores tipados sobre error wrapping genérico
- **Validación temprana**: Validar parámetros al inicio de funciones

```go
// ✅ Ejemplo correcto - Usar errores tipados
func (r *Repository) GetExpenses(ctx context.Context, userID int64, categoryID string) ([]domain.Expense, error) {
    if userID <= 0 {
        return nil, errors.NewBadRequest("invalid userID")
    }
    
    expenses, err := r.fetchExpenses(ctx, userID, categoryID)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to fetch expenses")
    }
    
    return expenses, nil
}

// ✅ Tipos de errores disponibles
func ExampleErrorTypes() {
    // Errores de cliente
    errors.NewBadRequest("invalid input")
    errors.NewResourceNotFound("expense not found")
    errors.NewUnauthorizedRequest("authentication required")
    
    // Errores de servidor
    errors.NewInternalServerError("database connection failed")
    errors.NewTooManyRequests("rate limit exceeded")
    
    // Errores específicos del dominio
    errors.NewResourceParsingError("invalid date format")
    errors.NewResourceAlreadyExists("category already exists")
}
```

### Context Management
- **Usar context.Context**: Para todas las operaciones asíncronas
- **Propagación**: Pasar context a través de todas las capas
- **Timeouts**: Configurar timeouts apropiados

```go
// ✅ Ejemplo correcto
func (s *Service) ProcessExpense(ctx context.Context, req *domain.ExpenseRequest) (*domain.Expense, error) {
    // Usar context en todas las llamadas downstream
    category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
    if err != nil {
        return nil, errors.NewInternalServerError("failed to get category")
    }
    
    return s.createExpense(ctx, req, category)
}
```

### Interfaces
- **Pequeñas y específicas**: Preferir interfaces pequeñas
- **Definir donde se usan**: No en packages globales
- **Composición**: Combinar interfaces pequeñas para crear mayores

```go
// ✅ Ejemplo correcto
type ExpenseRepository interface {
    Create(ctx context.Context, expense *domain.Expense) error
    GetByID(ctx context.Context, id string) (*domain.Expense, error)
    ListByUserID(ctx context.Context, userID int64) ([]domain.Expense, error)
}
```

### Convenciones de Naming
- **Descriptivos**: Nombres que expliquen el propósito
- **Consistentes**: Seguir convenciones de Go
- **Acrónimos**: En mayúsculas (ID, HTTP, JSON)

```go
// ✅ Correcto
type UserID int64
type HTTPClient interface{}
type APIResponse struct{}

// ❌ Incorrecto  
type UserId int64
type HttpClient interface{}
type ApiResponse struct{}
```

### Comentarios en el Código
- **NO comentar código**: El código debe ser autoexplicativo
- **Nombres descriptivos**: Usar nombres que expliquen el propósito
- **Funciones pequeñas**: Dividir en funciones más pequeñas y descriptivas

```go
// ❌ MAL - Comentar funcionalidad obvia
// Esta función suma dos números
func sum(a, b int) int {
    return a + b
}

// ✅ BIEN - Nombre descriptivo sin comentario
func calculateTotalWithTax(amount, taxRate float64) float64 {
    return amount * (1 + taxRate)
}
```

## Estructura de Funciones

### Orden de Validación
1. **Validación de parámetros**
2. **Lógica principal**
3. **Manejo de respuestas**

```go
// ✅ Estructura recomendada
func ProcessTransaction(ctx context.Context, input *domain.TransactionRequest) (*domain.Transaction, error) {
    // 1. Validaciones
    if input == nil {
        return nil, errors.NewBadRequest("input cannot be nil")
    }
    if input.UserID <= 0 {
        return nil, errors.NewBadRequest("invalid user ID")
    }
    
    // 2. Lógica principal
    transaction, err := processBusinessLogic(ctx, input)
    if err != nil {
        return nil, errors.NewInternalServerError("business logic failed")
    }
    
    // 3. Preparar respuesta
    return &domain.Transaction{
        ID:        transaction.ID,
        Amount:    transaction.Amount,
        CreatedAt: time.Now(),
    }, nil
}
```

## Logging Estructurado

### Formato Estándar
```go
// ✅ Usar logs estructurados con zap o similar
logger.Error("Error getting data from database", 
    zap.String("operation", "get_expenses"),
    zap.Int64("user_id", userID),
    zap.Error(err))

// ✅ Logs de información
logger.Info("Processing request", 
    zap.Int64("user_id", userID),
    zap.String("category", categoryID))
```

## Manejo de Fechas

### Uso de time.Time Estándar
- **Producción**: Usar `time.Time` nativo de Go para operaciones de fecha
- **Tests**: Usar mocks cuando sea necesario para tiempo controlado

```go
// ✅ Correcto - Uso estándar de time
type TransactionService struct {
    repo TransactionRepository
}

// ✅ Correcto - Uso en código de producción
func (s *TransactionService) CreateTransaction(transaction *domain.Transaction) error {
    transaction.CreatedAt = time.Now().UTC()
    return s.repo.Save(context.Background(), transaction)
}

// ✅ Correcto - Parsing de fechas
func ParseDateString(dateString string) (time.Time, error) {
    return time.Parse("2006-01-02", dateString)
}

// ✅ Correcto - Formateo de fechas
func FormatDateToISO(date time.Time) string {
    return date.Format("2006-01-02T15:04:05Z")
}
```

### Funciones Principales para Fechas
- `time.Now()` - Tiempo actual
- `time.Parse()` - Parse string a time.Time
- `time.Format()` - Formato de fecha a string
- `time.Date()` - Crear fecha específica
- `.UTC()` - Convertir a UTC
- `.Format("2006-01-02")` - Formato fecha YYYY-MM-DD

### Manejo de Fechas en Testing
```go
// ✅ Correcto - Fechas fijas en tests
func TestCreateTransaction_SetsCreatedAt(t *testing.T) {
    fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
    
    // Mock time si es necesario para el test
    transaction := &domain.Transaction{
        Amount:    100.0,
        CreatedAt: fixedTime,
    }
    
    assert.Equal(t, fixedTime, transaction.CreatedAt)
}
```
