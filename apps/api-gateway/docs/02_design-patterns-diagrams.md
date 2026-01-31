# Diagramas de Patrones de Diseño

## 1. Patrón Factory

```mermaid
classDiagram
    class TransactionFactory {
        <<interface>>
        +CreateTransaction()
        +GetTransaction()
        +ListTransactions()
        +UpdateTransaction()
        +DeleteTransaction()
    }
    
    class TransactionFactoryImpl {
        -incomeService IncomeService
        -expenseService ExpenseService
        +CreateTransaction()
        +GetTransaction()
        +ListTransactions()
        +UpdateTransaction()
        +DeleteTransaction()
    }
    
    class Transaction {
        <<interface>>
        +GetID()
        +GetUserID()
        +GetAmount()
        +GetDescription()
        +GetCategoryID()
        +GetCreatedAt()
        +GetUpdatedAt()
    }
    
    TransactionFactory <|.. TransactionFactoryImpl
    TransactionFactoryImpl --> Transaction
```

## 2. Patrón Strategy

```mermaid
classDiagram
    class TransactionStrategy {
        <<interface>>
        +Process(context.Context, Transaction) error
        +Validate(Transaction) error
        +CalculateBalanceImpact(Transaction) float64
    }
    
    class IncomeStrategy {
        -balanceService BalanceService
        +Process(context.Context, Transaction) error
        +Validate(Transaction) error
        +CalculateBalanceImpact(Transaction) float64
    }
    
    class ExpenseStrategy {
        -balanceService BalanceService
        -categoryService CategoryService
        +Process(context.Context, Transaction) error
        +Validate(Transaction) error
        +CalculateBalanceImpact(Transaction) float64
    }
    
    class TransactionProcessor {
        -strategy TransactionStrategy
        +SetStrategy(TransactionStrategy)
        +ProcessTransaction(context.Context, Transaction) error
    }
    
    class Transaction {
        <<interface>>
        +GetID() string
        +GetUserID() string
        +GetAmount() float64
        +GetDescription() string
        +GetCategoryID() string
    }
    
    class IncomeTransaction {
        +ID
        +UserID
        +Amount
        +Description
    }
    
    class ExpenseTransaction {
        +ID
        +UserID
        +Amount
        +Description
        +CategoryID
    }
    
    TransactionStrategy <|.. IncomeStrategy
    TransactionStrategy <|.. ExpenseStrategy
    TransactionProcessor --> TransactionStrategy
    Transaction <|.. IncomeTransaction
    Transaction <|.. ExpenseTransaction
    IncomeStrategy --> BalanceService
    ExpenseStrategy --> BalanceService
    ExpenseStrategy --> CategoryService
```

## 3. Clean Architecture

```mermaid
graph TD
    A[Handlers] --> B[Use Cases]
    B --> C[Domain]
    C --> D[Infrastructure]
    
    subgraph "Capa Externa"
    A
    end
    
    subgraph "Capa de Aplicación"
    B
    end
    
    subgraph "Capa de Dominio"
    C
    end
    
    subgraph "Capa de Infraestructura"
    D
    end
```

## 4. Patrón Repository

```mermaid
classDiagram
    class TransactionRepository {
        <<interface>>
        +Save()
        +FindByID()
        +FindAll()
        +Update()
        +Delete()
    }
    
    class TransactionService {
        -repository TransactionRepository
        +CreateTransaction()
        +GetTransaction()
        +UpdateTransaction()
        +DeleteTransaction()
    }
    
    class DatabaseRepository {
        +Save()
        +FindByID()
        +FindAll()
        +Update()
        +Delete()
    }
    
    TransactionRepository <|.. DatabaseRepository
    TransactionService --> TransactionRepository
```

## 5. Patrón Adapter

```mermaid
classDiagram
    class IncomeResponse {
        +ID
        +UserID
        +Amount
        +Description
    }
    
    class ExpenseResponse {
        +ID
        +UserID
        +Amount
        +Description
        +CategoryID
    }
    
    class TransactionAdapter {
        +toTransaction()
        +toTransactionSlice()
    }
    
    class Transaction {
        <<interface>>
        +GetID()
        +GetUserID()
        +GetAmount()
        +GetDescription()
        +GetCategoryID()
    }
    
    IncomeResponse --> TransactionAdapter
    ExpenseResponse --> TransactionAdapter
    TransactionAdapter --> Transaction
```

## 6. Flujo de una Transacción

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Factory
    participant Service
    participant Repository
    participant Database
    
    Client->>Handler: POST /api/v1/transactions
    Handler->>Factory: CreateTransaction()
    Factory->>Service: Process()
    Service->>Repository: Save()
    Repository->>Database: INSERT
    Database-->>Repository: OK
    Repository-->>Service: Transaction
    Service-->>Factory: Response
    Factory-->>Handler: Transaction
    Handler-->>Client: 201 Created
```

## 7. Principios SOLID

```mermaid
graph TD
    A[SOLID] --> B[Single Responsibility]
    A --> C[Open/Closed]
    A --> D[Liskov Substitution]
    A --> E[Interface Segregation]
    A --> F[Dependency Inversion]
    
    B --> G[Cada clase una responsabilidad]
    C --> H[Extensible sin modificar]
    D --> I[Subtipo sustituible]
    E --> J[Interfaces específicas]
    F --> K[Depender de abstracciones]
```

## 8. Estructura de Directorios

```mermaid
graph TD
    A[financial-resume-engine] --> B[cmd]
    A --> C[internal]
    A --> D[pkg]
    A --> E[docs]
    
    C --> F[handlers]
    C --> G[usecases]
    C --> H[core]
    C --> I[infrastructure]
    
    F --> J[incomes]
    F --> K[expenses]
    F --> L[categories]
    
    G --> M[transactions]
    G --> N[reports]
    
    H --> O[domain]
    H --> P[entities]
    
    I --> Q[persistence]
    I --> R[services]
```

## 9. Patrón Builder

```mermaid
classDiagram
    class TransactionBuilder {
        <<interface>>
        +SetID() TransactionBuilder
        +SetUserID() TransactionBuilder
        +SetAmount() TransactionBuilder
        +SetDescription() TransactionBuilder
        +SetCategoryID() TransactionBuilder
        +SetDueDate() TransactionBuilder
        +Build() Transaction
    }
    
    class transactionBuilder {
        -transaction Transaction
        +SetID() TransactionBuilder
        +SetUserID() TransactionBuilder
        +SetAmount() TransactionBuilder
        +SetDescription() TransactionBuilder
        +SetCategoryID() TransactionBuilder
        +SetDueDate() TransactionBuilder
        +Build() Transaction
    }
    
    class Transaction {
        <<interface>>
        +GetID()
        +GetUserID()
        +GetAmount()
        +GetDescription()
        +GetCategoryID()
        +GetDueDate()
    }
    
    class BaseTransaction {
        +ID
        +UserID
        +Amount
        +Description
        +CategoryID
        +DueDate
    }
    
    TransactionBuilder <|.. transactionBuilder
    transactionBuilder --> Transaction
    Transaction <|.. BaseTransaction
```

## 10. Patrón Observer

```mermaid
classDiagram
    class TransactionObserver {
        <<interface>>
        +OnTransactionCreated(context.Context, Transaction) error
        +OnTransactionUpdated(context.Context, Transaction) error
        +OnTransactionDeleted(context.Context, string) error
    }
    
    class TransactionSubject {
        -observers []TransactionObserver
        -mu sync.RWMutex
        +Attach(TransactionObserver)
        +Detach(TransactionObserver)
        +NotifyTransactionCreated(context.Context, Transaction) error
        +NotifyTransactionUpdated(context.Context, Transaction) error
        +NotifyTransactionDeleted(context.Context, string) error
    }
    
    class BalanceObserver {
        -balanceService BalanceService
        +OnTransactionCreated(context.Context, Transaction) error
        +OnTransactionUpdated(context.Context, Transaction) error
        +OnTransactionDeleted(context.Context, string) error
    }
    
    class NotificationObserver {
        -notificationService NotificationService
        +OnTransactionCreated(context.Context, Transaction) error
        +OnTransactionUpdated(context.Context, Transaction) error
        +OnTransactionDeleted(context.Context, string) error
    }
    
    class TransactionHandler {
        -subject *TransactionSubject
        +CreateTransaction(*gin.Context)
        +UpdateTransaction(*gin.Context)
        +DeleteTransaction(*gin.Context)
    }
    
    TransactionObserver <|.. BalanceObserver
    TransactionObserver <|.. NotificationObserver
    TransactionSubject --> TransactionObserver
    TransactionHandler --> TransactionSubject
    BalanceObserver --> BalanceService
    NotificationObserver --> NotificationService
``` 