# 🧪 Reglas de Testing

## Comandos de Testing

### Con Go Nativo
```bash
# Tests básicos
go test ./... -count=1

# Tests con race detection
go test ./... -count=1 -race

# Tests con verbose output
go test ./... -v
```

### Coverage
```bash
# Generar coverage de un package específico
go test -coverprofile=coverage.out [package_path]

# Ejemplo
go test -coverprofile=coverage.out github.com/financial-resume-engine/internal/core/usecases/expenses

# Visualizar coverage
go tool cover -html="coverage.out"

# Coverage para todo el proyecto
go test ./... -coverprofile=coverage.out
```

## Convenciones de Testing

### Estructura de Tests
Usar el patrón **Given-When-Then**:

```go
// ✅ Estructura recomendada
func TestFunctionName_Scenario_ExpectedBehavior(t *testing.T) {
    // Given - Preparar datos de prueba
    input := setupTestData()
    expectedResult := "expected_value"
    
    // When - Ejecutar función bajo prueba
    result, err := FunctionToTest(input)
    
    // Then - Verificar resultados
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}
```

### Convenciones de Nombrado de Variables en Test Cases

Para estructuras de casos de prueba, usar prefijos específicos para identificar claramente el propósito de cada variable:

#### Prefijo `given` - Datos de Entrada y Configuración
Variables que **nosotros establecemos** como datos de entrada, configuración de mocks, o condiciones iniciales:

```go
// ✅ Ejemplos de variables "given"
type testScenario struct {
    name                        string
    givenUserID                 int64                     // Datos de entrada
    givenCategoryID             string                    // Parámetros de configuración
    givenExpenseData            *domain.Expense           // Estados iniciales
    givenRepositoryError        error                     // Errores simulados en mocks
    givenDatabaseError          error                     // Comportamiento de dependencias
    givenCurrentTime            time.Time                // Respuestas mockeadas
}
```

#### Prefijo `expected` - Resultados Esperados
Variables que representan **los resultados que esperamos** obtener del test:

```go
// ✅ Ejemplos de variables "expected"  
type testScenario struct {
    expectedError               error  // Error esperado como resultado
    expectedRepositoryCalls     int    // Número de llamadas esperadas
    expectedDatabaseCalls       int    // Interacciones esperadas con mocks
    expectedNotificationCalls   int    // Validaciones de comportamiento
    expectedResult              *domain.Expense // Valores de retorno esperados
}
```

#### Ejemplo Completo
```go
type testScenario struct {
    name                        string
    // Given - Lo que establecemos
    givenUserID                 int64
    givenExpenseData            *domain.Expense
    givenRepositoryError        error
    givenCategoryResponse       *domain.Category
    
    // Expected - Lo que esperamos
    expectedError               error
    expectedResult              *domain.Expense
    expectedRepositoryCalls     int
    expectedCategoryCalls       int
}
```

#### ⚠️ Beneficios de esta Convención
- **Claridad**: Inmediatamente se identifica si una variable es entrada o resultado esperado
- **Mantenibilidad**: Facilita la lectura y modificación de tests complejos
- **Consistencia**: Patrón uniforme en todo el codebase para estructuras de test cases

### Naming Convention
```go
// Formato: TestFunctionName_Scenario_ExpectedBehavior
func TestCreateExpense_WithValidData_ReturnsExpense(t *testing.T) {}
func TestCreateExpense_WithInvalidUserID_ReturnsError(t *testing.T) {}
func TestListExpenses_WithEmptyCategory_ReturnsAllExpenses(t *testing.T) {}
```

## Convenciones de Mocks

### Ubicación y Nomenclatura
- **Ubicación**: Mismo package que lo que mockea
- **Nombre**: Igual al original + sufijo `Mock`
- **Archivo**: Termina en `_mock.go`

```go
// ✅ Ejemplo: internal/core/repository/expense_repository_mock.go
type ExpenseRepositoryMock struct {
    mock.Mock
}

func (m *ExpenseRepositoryMock) Create(ctx context.Context, expense *domain.Expense) error {
    args := m.Called(ctx, expense)
    return args.Error(0)
}

func (m *ExpenseRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Expense, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*domain.Expense), args.Error(1)
}
```

### Configuración de Mocks
```go
// ✅ Setup en tests - SER ESPECÍFICO en los parámetros
func TestExpenseService_CreateExpense_Success(t *testing.T) {
    // Given
    mockRepo := new(ExpenseRepositoryMock)
    service := NewExpenseService(mockRepo)
    
    expectedExpense := &domain.Expense{
        ID:     "123",
        Amount: 100.50,
        UserID: 456,
    }
    
    // ✅ Especificar valores exactos en lugar de mock.Anything
    mockRepo.On("Create", 
        context.Background(),
        expectedExpense,
    ).Return(nil)
    
    // When
    err := service.CreateExpense(context.Background(), expectedExpense)
    
    // Then
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### ⚠️ Reglas para mock.Anything
```go
// ❌ MAL - Usar mock.Anything indiscriminadamente
mockRepo.On("Create", mock.Anything, mock.Anything)

// ✅ BIEN - Solo cuando realmente no importa el valor específico
mockRepo.On("Create", 
    mock.AnythingOfType("*context.Context"), // Cuando el contexto no es relevante para el test
    mock.AnythingOfType("*domain.Expense"),  // Solo si los datos exactos no son relevantes
)

// 🎯 MEJOR - Ser específico siempre que sea posible
mockRepo.On("Create",
    context.Background(),
    &domain.Expense{
        Amount: 100.50,
        UserID: 456,
    },
).Return(nil)
```

### 🎯 Best Practices para Mocks

#### Principios Fundamentales
- **Ser específico** es mejor que usar `mock.Anything`
- **Validar comportamiento** no solo resultados
- **Un mock por responsabilidad** - no mockear todo
- **Limpiar mocks** entre tests

#### Cuándo usar mock.Anything
```go
// ✅ ACEPTABLE - Context que no afecta la lógica de negocio
mockService.On("ProcessExpense", mock.AnythingOfType("*context.Context"), specificExpense)

// ✅ ACEPTABLE - IDs generados dinámicamente que no se pueden predecir
mockRepo.On("Save", mock.AnythingOfType("*domain.Expense")).Return(nil)

// ❌ EVITAR - Parámetros que sí importan para la lógica
mockRepo.On("FindByCategory", mock.Anything) // ¿Qué categoría? ¡Debe ser específico!
```

#### Validaciones en Mocks
```go
// ✅ Validar que se llamó con parámetros correctos
mockRepo.On("UpdateStatus", "123", "approved").Return(nil).Once()

// ✅ Validar número de llamadas
mockRepo.AssertNumberOfCalls(t, "UpdateStatus", 1)

// ✅ Validar que NO se llamó un método
mockRepo.AssertNotCalled(t, "Delete")
```

## Estrategias de Testing

### Tipos de Tests
- **Unit Tests**: Para cada función/método público
- **Integration Tests**: Para flujos completos entre componentes
- **Repository Tests**: Para validar integración con base de datos

### Coverage Goals
- **Mínimo**: 80% de coverage
- **Objetivo**: 90%+ para código crítico
- **Foco**: Priorizar paths principales y casos de error

### Test Data Management

#### Para Tests Unitarios
```go
// ✅ Usar fechas fijas en tests
func setupTestData() *domain.Expense {
    testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
    
    return &domain.Expense{
        ID:          "test-123",
        Amount:      100.50,
        UserID:      456,
        CategoryID:  "food",
        Description: "Test expense",
        CreatedAt:   testDate,
    }
}
```

## Testing Local

### Configuración de Tests
- **Base de datos**: Usar base de datos de testing separada
- **Mocks**: Para dependencias externas
- **Fixtures**: Datos de prueba consistentes

### Ejecutar Tests
```bash
# Tests unitarios
go test ./internal/... -count=1

# Tests de integración (si existen)
go test ./tests/integration/... -count=1

# Tests con coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Checklist de Testing

- [ ] ¿Tests unitarios para nuevas funciones públicas?
- [ ] ¿Tests de error cases y edge cases?
- [ ] ¿Mocks siguen las convenciones de naming?
- [ ] ¿Coverage mínimo del 80%?
- [ ] ¿Tests de integración para flujos críticos?
- [ ] ¿Son los tests unitarios lo más atómicos posibles?
- [ ] ¿Se usan datos de test específicos en lugar de mock.Anything?
- [ ] ¿Se validan tanto resultados como comportamiento de mocks?  