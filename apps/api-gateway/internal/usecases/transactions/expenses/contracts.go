package expenses

// CreateExpenseRequest representa el request para crear un gasto
type CreateExpenseRequest struct {
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	CategoryID      string  `json:"category_id,omitempty"`
	Paid            bool    `json:"paid"`
	DueDate         string  `json:"due_date,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty"` // NEW: When expense occurred
}

// CreateExpenseResponse representa la respuesta al crear un gasto
type CreateExpenseResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`         // Monto total original
	AmountPaid      float64 `json:"amount_paid"`    // Monto ya pagado
	PendingAmount   float64 `json:"pending_amount"` // Monto pendiente
	Description     string  `json:"description"`
	CategoryID      string  `json:"category_id,omitempty"`
	Paid            bool    `json:"paid"`
	DueDate         string  `json:"due_date,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty"` // NEW: When expense occurred
	Percentage      float64 `json:"percentage,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// GetExpenseResponse representa la respuesta al obtener un gasto
type GetExpenseResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`         // Monto total original
	AmountPaid      float64 `json:"amount_paid"`    // Monto ya pagado
	PendingAmount   float64 `json:"pending_amount"` // Monto pendiente
	Description     string  `json:"description"`
	CategoryID      string  `json:"category_id,omitempty"`
	Paid            bool    `json:"paid"`
	DueDate         string  `json:"due_date,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty"` // NEW: When expense occurred
	Percentage      float64 `json:"percentage,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// ListExpensesResponse representa la respuesta al listar gastos
type ListExpensesResponse struct {
	Expenses []GetExpenseResponse `json:"expenses"`
}

// UpdateExpenseRequest representa el request para actualizar un gasto
type UpdateExpenseRequest struct {
	Amount          float64 `json:"amount,omitempty"`         // Monto total (solo para editar)
	PaymentAmount   float64 `json:"payment_amount,omitempty"` // Monto de pago (para pagos parciales)
	Description     string  `json:"description,omitempty"`
	CategoryID      string  `json:"category_id,omitempty"`
	Paid            bool    `json:"paid,omitempty"`
	DueDate         string  `json:"due_date,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty"` // NEW: When expense occurred
}

// UpdateExpenseResponse representa la respuesta al actualizar un gasto
type UpdateExpenseResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`         // Monto total original
	AmountPaid      float64 `json:"amount_paid"`    // Monto ya pagado
	PendingAmount   float64 `json:"pending_amount"` // Monto pendiente
	Description     string  `json:"description"`
	CategoryID      string  `json:"category_id,omitempty"`
	Paid            bool    `json:"paid"`
	DueDate         string  `json:"due_date,omitempty"`
	TransactionDate string  `json:"transaction_date,omitempty"` // NEW: When expense occurred
	Percentage      float64 `json:"percentage,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
