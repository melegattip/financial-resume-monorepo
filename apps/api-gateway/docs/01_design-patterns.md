# Patrones de Diseño en Financial Resume Engine

## Índice
1. [Introducción](#introducción)
2. [Patrones Estructurales](#patrones-estructurales)
3. [Patrones Creacionales](#patrones-creacionales)
4. [Patrones de Comportamiento](#patrones-de-comportamiento)
5. [Arquitectura](#arquitectura)

## Introducción

Financial Resume Engine implementa varios patrones de diseño que siguen las mejores prácticas de desarrollo de software. Estos patrones ayudan a mantener el código organizado, escalable y fácil de mantener.

## Patrones Estructurales

### 1. Patrón Adapter

El patrón Adapter permite que interfaces incompatibles trabajen juntas. En nuestro caso, se utiliza para adaptar diferentes tipos de transacciones a una interfaz común.

```go
// Ejemplo de implementación
func toTransaction(response interface{}) Transaction {
    switch v := response.(type) {
    case *incomes.CreateIncomeResponse:
        return &TransactionModel{
            ID:          v.ID,
            UserID:      v.UserID,
            Amount:      v.Amount,
            Description: v.Description,
        }
    // ... otros casos
    }
}
```

Diagrama:
```
+----------------+     +----------------+     +----------------+
| IncomeResponse | --> |    Adapter     | --> |  Transaction   |
+----------------+     +----------------+     +----------------+
```

### 2. Patrón Repository

El patrón Repository abstrae el acceso a datos, permitiendo que la lógica de negocio no dependa directamente de la capa de persistencia.

```go
type TransactionRepository interface {
    Save(ctx context.Context, transaction Transaction) error
    FindByID(ctx context.Context, id string) (Transaction, error)
    // ... otros métodos
}
```

Diagrama:
```
+----------------+     +----------------+     +----------------+
|    Service     | --> |  Repository   | --> |  Database      |
+----------------+     +----------------+     +----------------+
```

## Patrones Creacionales

### 1. Patrón Factory

El patrón Factory se utiliza para crear diferentes tipos de transacciones sin exponer la lógica de creación al cliente.

```go
type TransactionFactory interface {
    CreateTransaction(ctx context.Context, userID string, amount float64, 
                     description, categoryID string, dueDate *time.Time) (Transaction, error)
}

type TransactionFactoryImpl struct {
    incomeService  IncomeService
    expenseService ExpenseService
}
```

Diagrama:
```
+----------------+     +----------------+     +----------------+
|    Client      | --> |    Factory    | --> |  Concrete      |
|                |     |               |     |  Products      |
+----------------+     +----------------+     +----------------+
```

### 2. Patrón Builder

El patrón Builder se utiliza para construir objetos complejos paso a paso, permitiendo diferentes representaciones del mismo objeto.

```go
// Interfaz del Builder
type TransactionBuilder interface {
    SetID(id string) TransactionBuilder
    SetUserID(userID string) TransactionBuilder
    SetAmount(amount float64) TransactionBuilder
    SetDescription(description string) TransactionBuilder
    SetCategoryID(categoryID string) TransactionBuilder
    SetDueDate(dueDate time.Time) TransactionBuilder
    Build() Transaction
}

// Implementación concreta del Builder
type transactionBuilder struct {
    transaction Transaction
}

func NewTransactionBuilder() TransactionBuilder {
    return &transactionBuilder{
        transaction: &BaseTransaction{},
    }
}

func (b *transactionBuilder) SetID(id string) TransactionBuilder {
    b.transaction.ID = id
    return b
}

func (b *transactionBuilder) SetUserID(userID string) TransactionBuilder {
    b.transaction.UserID = userID
    return b
}

func (b *transactionBuilder) SetAmount(amount float64) TransactionBuilder {
    b.transaction.Amount = amount
    return b
}

func (b *transactionBuilder) SetDescription(description string) TransactionBuilder {
    b.transaction.Description = description
    return b
}

func (b *transactionBuilder) SetCategoryID(categoryID string) TransactionBuilder {
    b.transaction.CategoryID = categoryID
    return b
}

func (b *transactionBuilder) SetDueDate(dueDate time.Time) TransactionBuilder {
    b.transaction.DueDate = dueDate
    return b
}

func (b *transactionBuilder) Build() Transaction {
    return b.transaction
}

// Uso del Builder
transaction := NewTransactionBuilder().
    SetID(uuid.New().String()).
    SetUserID(userID).
    SetAmount(100.50).
    SetDescription("Salario mensual").
    SetCategoryID("salary").
    SetDueDate(time.Now()).
    Build()
```

Diagrama:
```
+----------------+     +----------------+     +----------------+
|    Director    | --> |    Builder     | --> |   Producto     |
|                |     |   Interface    |     |   (Transaction)|
+----------------+     +----------------+     +----------------+
```

Este ejemplo muestra cómo:
1. El Builder permite construir el objeto paso a paso
2. Cada método retorna el Builder para permitir encadenamiento de métodos
3. El método `Build()` finaliza la construcción y retorna el objeto completo
4. Se puede crear diferentes tipos de transacciones usando el mismo Builder

## Patrones de Comportamiento

### 1. Patrón Strategy

El patrón Strategy permite definir diferentes algoritmos o comportamientos que pueden ser intercambiables según el contexto. En nuestro caso, lo usamos para manejar diferentes tipos de transacciones (ingresos y gastos) que tienen comportamientos específicos.

```go
// Interfaz que define la estrategia
type TransactionStrategy interface {
    // Procesa la transacción según su tipo específico
    Process(ctx context.Context, transaction Transaction) error
    // Valida la transacción según sus reglas específicas
    Validate(transaction Transaction) error
    // Calcula el impacto en el balance según el tipo
    CalculateBalanceImpact(transaction Transaction) float64
}

// Estrategia para ingresos
type IncomeStrategy struct {
    balanceService BalanceService
}

func (s *IncomeStrategy) Process(ctx context.Context, transaction Transaction) error {
    // Lógica específica para procesar ingresos
    if err := s.Validate(transaction); err != nil {
        return err
    }
    
    // Actualizar el balance con el ingreso
    impact := s.CalculateBalanceImpact(transaction)
    return s.balanceService.AddToBalance(ctx, transaction.GetUserID(), impact)
}

func (s *IncomeStrategy) Validate(transaction Transaction) error {
    // Validaciones específicas para ingresos
    if transaction.GetAmount() <= 0 {
        return errors.New("el monto del ingreso debe ser positivo")
    }
    return nil
}

func (s *IncomeStrategy) CalculateBalanceImpact(transaction Transaction) float64 {
    // Los ingresos aumentan el balance
    return transaction.GetAmount()
}

// Estrategia para gastos
type ExpenseStrategy struct {
    balanceService BalanceService
    categoryService CategoryService
}

func (s *ExpenseStrategy) Process(ctx context.Context, transaction Transaction) error {
    // Lógica específica para procesar gastos
    if err := s.Validate(transaction); err != nil {
        return err
    }
    
    // Verificar categoría
    if err := s.categoryService.ValidateCategory(ctx, transaction.GetCategoryID()); err != nil {
        return err
    }
    
    // Actualizar el balance con el gasto
    impact := s.CalculateBalanceImpact(transaction)
    return s.balanceService.SubtractFromBalance(ctx, transaction.GetUserID(), impact)
}

func (s *ExpenseStrategy) Validate(transaction Transaction) error {
    // Validaciones específicas para gastos
    if transaction.GetAmount() <= 0 {
        return errors.New("el monto del gasto debe ser positivo")
    }
    if transaction.GetCategoryID() == "" {
        return errors.New("los gastos deben tener una categoría")
    }
    return nil
}

func (s *ExpenseStrategy) CalculateBalanceImpact(transaction Transaction) float64 {
    // Los gastos disminuyen el balance
    return -transaction.GetAmount()
}

// Contexto que usa la estrategia
type TransactionProcessor struct {
    strategy TransactionStrategy
}

func (p *TransactionProcessor) SetStrategy(strategy TransactionStrategy) {
    p.strategy = strategy
}

func (p *TransactionProcessor) ProcessTransaction(ctx context.Context, transaction Transaction) error {
    return p.strategy.Process(ctx, transaction)
}

// Uso del patrón Strategy
func main() {
    // Crear el procesador
    processor := &TransactionProcessor{}
    
    // Crear una transacción de ingreso
    income := &IncomeTransaction{
        Amount: 1000.00,
        Description: "Salario",
    }
    
    // Configurar estrategia para ingresos
    processor.SetStrategy(&IncomeStrategy{})
    if err := processor.ProcessTransaction(ctx, income); err != nil {
        log.Fatal(err)
    }
    
    // Crear una transacción de gasto
    expense := &ExpenseTransaction{
        Amount: 500.00,
        Description: "Compras",
        CategoryID: "groceries",
    }
    
    // Configurar estrategia para gastos
    processor.SetStrategy(&ExpenseStrategy{})
    if err := processor.ProcessTransaction(ctx, expense); err != nil {
        log.Fatal(err)
    }
}

Este ejemplo muestra cómo:
1. Cada tipo de transacción (ingreso/gasto) tiene su propia estrategia
2. Las estrategias implementan la misma interfaz pero con comportamientos diferentes
3. El contexto (`TransactionProcessor`) puede cambiar de estrategia según el tipo de transacción
4. Cada estrategia tiene sus propias reglas de validación y cálculo de impacto

Beneficios del patrón Strategy:
- Separación clara de responsabilidades
- Fácil adición de nuevos tipos de transacciones
- Código más mantenible y testeable
- Flexibilidad para cambiar comportamientos en tiempo de ejecución.
```

### 2. Patrón Observer

El patrón Observer permite que un objeto (sujeto) notifique a otros objetos (observadores) sobre cambios en su estado. En nuestra API, lo usamos para manejar eventos de transacciones y notificar a diferentes servicios.

```go
// Interfaz para los observadores
type TransactionObserver interface {
    OnTransactionCreated(ctx context.Context, transaction Transaction) error
    OnTransactionUpdated(ctx context.Context, transaction Transaction) error
    OnTransactionDeleted(ctx context.Context, transactionID string) error
}

// Sujeto que mantiene la lista de observadores
type TransactionSubject struct {
    observers []TransactionObserver
    mu        sync.RWMutex
}

func NewTransactionSubject() *TransactionSubject {
    return &TransactionSubject{
        observers: make([]TransactionObserver, 0),
    }
}

// Métodos para gestionar observadores
func (s *TransactionSubject) Attach(observer TransactionObserver) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.observers = append(s.observers, observer)
}

func (s *TransactionSubject) Detach(observer TransactionObserver) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for i, o := range s.observers {
        if o == observer {
            s.observers = append(s.observers[:i], s.observers[i+1:]...)
            break
        }
    }
}

// Métodos para notificar a los observadores
func (s *TransactionSubject) NotifyTransactionCreated(ctx context.Context, transaction Transaction) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var errs []error
    for _, observer := range s.observers {
        if err := observer.OnTransactionCreated(ctx, transaction); err != nil {
            errs = append(errs, err)
        }
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("errores al notificar creación: %v", errs)
    }
    return nil
}

// Implementaciones concretas de observadores
type BalanceObserver struct {
    balanceService BalanceService
}

func (o *BalanceObserver) OnTransactionCreated(ctx context.Context, transaction Transaction) error {
    // Actualizar balance cuando se crea una transacción
    return o.balanceService.UpdateBalance(ctx, transaction.GetUserID(), transaction.GetAmount())
}

func (o *BalanceObserver) OnTransactionUpdated(ctx context.Context, transaction Transaction) error {
    // Actualizar balance cuando se modifica una transacción
    return o.balanceService.UpdateBalance(ctx, transaction.GetUserID(), transaction.GetAmount())
}

func (o *BalanceObserver) OnTransactionDeleted(ctx context.Context, transactionID string) error {
    // Manejar la eliminación de transacción
    return nil
}

type NotificationObserver struct {
    notificationService NotificationService
}

func (o *NotificationObserver) OnTransactionCreated(ctx context.Context, transaction Transaction) error {
    // Enviar notificación al usuario
    return o.notificationService.SendTransactionNotification(ctx, transaction.GetUserID(), "Nueva transacción creada")
}

// Uso del patrón Observer en el handler
type TransactionHandler struct {
    subject *TransactionSubject
}

func NewTransactionHandler() *TransactionHandler {
    handler := &TransactionHandler{
        subject: NewTransactionSubject(),
    }
    
    // Registrar observadores
    handler.subject.Attach(&BalanceObserver{})
    handler.subject.Attach(&NotificationObserver{})
    
    return handler
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    var request CreateTransactionRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Crear la transacción
    transaction, err := h.createTransaction(c.Request.Context(), request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Notificar a los observadores
    if err := h.subject.NotifyTransactionCreated(c.Request.Context(), transaction); err != nil {
        // Loggear el error pero no fallar la operación
        log.Printf("Error al notificar observadores: %v", err)
    }
    
    c.JSON(http.StatusCreated, transaction)
}

Este ejemplo muestra cómo:
1. El `TransactionSubject` mantiene una lista de observadores
2. Los observadores (`BalanceObserver`, `NotificationObserver`) implementan la interfaz `TransactionObserver`
3. Cuando ocurre un evento (creación, actualización, eliminación), se notifica a todos los observadores
4. Cada observador puede reaccionar al evento de manera diferente

Beneficios del patrón Observer:
- Desacoplamiento entre el sujeto y los observadores
- Fácil adición de nuevos observadores sin modificar el código existente
- Notificaciones asíncronas y no bloqueantes
- Flexibilidad para manejar diferentes tipos de eventos
```

## Arquitectura

### Clean Architecture

La aplicación sigue los principios de Clean Architecture, con una clara separación de responsabilidades:

```
+----------------+
|   Handlers     |  <- Capa más externa
+----------------+
        ↓
+----------------+
|   Use Cases    |  <- Lógica de negocio
+----------------+
        ↓
+----------------+
|    Domain      |  <- Entidades y reglas de negocio
+----------------+
        ↓
+----------------+
| Infrastructure |  <- Implementaciones técnicas
+----------------+
```

### Principios SOLID

1. **Single Responsibility Principle (SRP)**
   - Cada componente tiene una única responsabilidad
   - Ejemplo: Separación de handlers, servicios y repositorios

2. **Open/Closed Principle (OCP)**
   - Las entidades están abiertas para extensión pero cerradas para modificación
   - Ejemplo: Interfaz Transaction que puede extenderse con nuevos tipos

3. **Liskov Substitution Principle (LSP)**
   - Los subtipos deben ser sustituibles por sus tipos base
   - Ejemplo: Income y Expense pueden sustituir a Transaction

4. **Interface Segregation Principle (ISP)**
   - Las interfaces deben ser específicas del cliente
   - Ejemplo: Separación de IncomeService y ExpenseService

5. **Dependency Inversion Principle (DIP)**
   - Las dependencias deben apuntar a abstracciones
   - Ejemplo: Inyección de dependencias en constructores

## Ejemplos de Implementación

### Factory Pattern en Acción

```go
// Creación de una transacción usando el factory
transaction, err := factory.CreateTransaction(ctx, userID, amount, description, categoryID, dueDate)
if err != nil {
    return nil, err
}
```

### Strategy Pattern en Acción

```go
// Selección de estrategia basada en el tipo de transacción
var strategy TransactionStrategy
if transactionType == IncomeType {
    strategy = &IncomeStrategy{}
} else {
    strategy = &ExpenseStrategy{}
}

// Uso de la estrategia
err := strategy.Process(ctx, transaction)
```

## Beneficios de los Patrones Implementados

1. **Mantenibilidad**
   - Código organizado y fácil de entender
   - Separación clara de responsabilidades

2. **Escalabilidad**
   - Fácil adición de nuevos tipos de transacciones
   - Extensibilidad del sistema

3. **Testabilidad**
   - Componentes aislados y fáciles de probar
   - Uso de interfaces para mocking

4. **Flexibilidad**
   - Fácil adaptación a nuevos requisitos
   - Intercambio de implementaciones

5. **Reusabilidad**
   - Componentes modulares
   - Código compartido entre diferentes partes del sistema 