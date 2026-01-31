# 🚨 Reglas de Validación

## Input Validation

### Validaciones Obligatorias
- **Parámetros nulos**: Verificar que no sean nil/null
- **Rangos numéricos**: Validar que estén dentro de rangos esperados
- **Longitud de strings**: Verificar límites máximos/mínimos
- **Formatos**: Validar emails, fechas, IDs, etc.

### Ejemplos de Validación
```go
// ✅ Validación de parámetros de entrada
func ValidateExpenseRequest(userID int64, amount float64, categoryID string) error {
    if userID <= 0 {
        return errors.NewBadRequest("userID must be positive")
    }
    
    if amount <= 0 {
        return errors.NewBadRequest("amount must be positive")
    }
    
    if len(categoryID) == 0 {
        return errors.NewBadRequest("categoryID cannot be empty")
    }
    
    if len(categoryID) > 50 {
        return errors.NewBadRequest("categoryID too long")
    }
    
    return nil
}

// ✅ Validación en funciones
func (r *ExpenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
    // Validación temprana
    if err := ValidateExpenseRequest(expense.UserID, expense.Amount, expense.CategoryID); err != nil {
        return err // Error ya tipado, no necesita wrapping
    }
    
    // Lógica principal...
}
```

## Security Validations

### Sanitización de Datos
```go
// ✅ Limpiar inputs antes de procesarlos
func SanitizeUserInput(input string) string {
    // Remover caracteres peligrosos
    cleaned := strings.TrimSpace(input)
    cleaned = strings.ReplaceAll(cleaned, "<", "")
    cleaned = strings.ReplaceAll(cleaned, ">", "")
    return cleaned
}

// ✅ Validar formato de IDs
func ValidateUserID(userID string) error {
    if _, err := strconv.ParseInt(userID, 10, 64); err != nil {
        return errors.NewBadRequest("invalid user ID format")
    }
    return nil
}
```

### Authentication & Authorization
```go
// ✅ Validar headers de autenticación
func ValidateAuthHeaders(headers map[string]string) error {
    userID, exists := headers["X-User-Id"]
    if !exists || userID == "" {
        return errors.NewBadRequest("X-User-Id header is required")
    }
    
    if _, err := strconv.ParseInt(userID, 10, 64); err != nil {
        return errors.NewBadRequest("invalid X-User-Id format")
    }
    
    return nil
}
```

## Business Logic Validation

### Validaciones de Dominio
```go
// ✅ Validar reglas de negocio
func ValidateTransactionRequest(req *domain.TransactionRequest) error {
    if req.UserID <= 0 {
        return errors.NewBadRequest("invalid user ID")
    }
    
    // Validar categorías disponibles
    validCategories := map[string]bool{
        "food":          true,
        "transportation": true,
        "entertainment": true,
        "utilities":     true,
        "healthcare":    true,
    }
    
    if !validCategories[req.CategoryID] {
        return errors.NewBadRequest("invalid category: " + req.CategoryID)
    }
    
    // Validar rangos de montos
    if req.Amount <= 0 {
        return errors.NewBadRequest("amount must be positive")
    }
    
    if req.Amount > 10000 {
        return errors.NewBadRequest("amount exceeds maximum limit")
    }
    
    return nil
}
```

### Estados y Transiciones
```go
// ✅ Validar transiciones de estado
func ValidateStatusTransition(from, to string) error {
    validTransitions := map[string][]string{
        "pending":   {"approved", "rejected"},
        "approved":  {"completed", "cancelled"},
        "completed": {"archived"},
        "rejected":  {"pending"},
    }
    
    validToStates, exists := validTransitions[from]
    if !exists {
        return errors.NewBadRequest("invalid from status: " + from)
    }
    
    for _, validTo := range validToStates {
        if validTo == to {
            return nil
        }
    }
    
    return errors.NewBadRequest("invalid transition from " + from + " to " + to)
}
```

## Error Handling

### Tipos de Errores Disponibles
```go
// ✅ Errores de cliente (4xx)
errors.NewBadRequest("mensaje")           // 400 - Input inválido
errors.NewUnauthorizedRequest("mensaje")  // 401 - Sin autenticación  
errors.NewResourceNotFound("mensaje")     // 404 - Recurso no encontrado

// ✅ Errores de servidor (5xx)
errors.NewInternalServerError("mensaje")  // 500 - Error interno
errors.NewTooManyRequests("mensaje")      // 429 - Rate limit excedido

// ✅ Errores específicos del dominio
errors.NewResourceParsingError("mensaje") // Error parseando datos
errors.NewResourceAlreadyExists("mensaje") // Recurso ya existe
```

### Ejemplos de Uso por Contexto
```go
// ✅ Funciones de error específicas por tipo de validación
func NewInvalidUserIDError() error {
    return errors.NewBadRequest("invalid user ID")
}

func NewInvalidCategoryError() error {
    return errors.NewBadRequest("invalid category")
}

func NewInvalidAmountError() error {
    return errors.NewBadRequest("invalid amount")
}

func NewMissingAuthHeaderError() error {
    return errors.NewBadRequest("missing authentication header")
}

// ✅ Manejo de errores sin wrapping genérico
func ValidateAndProcess(req *domain.TransactionRequest) error {
    if err := ValidateTransactionRequest(req); err != nil {
        return err // Error ya tipado, no necesita wrapping
    }
    
    if err := ProcessTransaction(req); err != nil {
        return errors.NewInternalServerError("failed to process transaction")
    }
    
    return nil
}
```

## API Response Validation

### Validar Responses Externas
```go
// ✅ Validar respuestas de APIs externas
func ValidateAPIResponse(response *APIResponse) error {
    if response == nil {
        return errors.NewBadRequest("response cannot be nil")
    }
    
    if response.StatusCode < 200 || response.StatusCode >= 300 {
        return errors.NewInternalServerError("API error: status " + strconv.Itoa(response.StatusCode))
    }
    
    if len(response.Data) == 0 {
        return errors.NewResourceNotFound("empty response data")
    }
    
    return nil
}
```

### Validar Formato de Datos
```go
// ✅ Validar estructura de datos
func ValidateExpenseResults(results []domain.Expense) error {
    if len(results) == 0 {
        return nil // Empty results are valid
    }
    
    for i, expense := range results {
        if expense.ID == "" {
            return errors.NewResourceParsingError("missing ID in expense " + strconv.Itoa(i))
        }
        
        if expense.Amount <= 0 {
            return errors.NewResourceParsingError("invalid amount in expense " + strconv.Itoa(i))
        }
        
        if expense.CategoryID == "" {
            return errors.NewResourceParsingError("missing category_id in expense " + strconv.Itoa(i))
        }
        
        if expense.CreatedAt.IsZero() {
            return errors.NewResourceParsingError("missing created_at in expense " + strconv.Itoa(i))
        }
    }
    
    return nil
}
```

## Logging de Validaciones

### Log de Errores de Validación
```go
// ✅ Log structured para validaciones
func LogValidationError(ctx context.Context, err error, operation string, userID int64) {
    logger.Error("Validation failed", 
        zap.String("operation", operation),
        zap.Int64("user_id", userID),
        zap.String("error_type", "validation"),
        zap.Error(err))
}

// ✅ Uso en funciones
func (s *ExpenseService) CreateExpense(ctx context.Context, req *domain.ExpenseRequest) error {
    if err := ValidateTransactionRequest(req); err != nil {
        LogValidationError(ctx, err, "create_expense", req.UserID)
        return err // Error ya tipado, no necesita wrapping
    }
    
    // Continuar procesamiento...
}
```

## Checklist de Validación

- [ ] ¿Se validan todos los parámetros de entrada?
- [ ] ¿Se manejan casos de inputs nulos/vacíos?
- [ ] ¿Se verifican rangos numéricos apropiados?
- [ ] ¿Se validan headers de autenticación?
- [ ] ¿Se sanitizan inputs antes de procesarlos?
- [ ] ¿Se validan reglas de negocio específicas?
- [ ] ¿Se logean errores de validación apropiadamente?
- [ ] ¿Se usan constructores de errores tipados ?
- [ ] ¿Se retornan errores descriptivos y tipados al cliente?
- [ ] ¿Se validan transiciones de estado cuando aplique?
- [ ] ¿Se verifican límites de montos y rangos financieros? 