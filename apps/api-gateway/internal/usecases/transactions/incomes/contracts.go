package incomes

// CreateIncomeRequest representa el request para crear un ingreso
type CreateIncomeRequest struct {
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id,omitempty"`
	Source      string  `json:"source"`
}

// CreateIncomeResponse representa la respuesta al crear un ingreso
type CreateIncomeResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// GetIncomeResponse representa la respuesta al obtener un ingreso
type GetIncomeResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ListIncomesResponse representa la respuesta al listar ingresos
type ListIncomesResponse struct {
	Incomes []GetIncomeResponse `json:"incomes"`
}

// UpdateIncomeRequest representa el request para actualizar un ingreso
type UpdateIncomeRequest struct {
	Amount      float64 `json:"amount,omitempty"`
	Description string  `json:"description,omitempty"`
	CategoryID  string  `json:"category_id,omitempty"`
	Source      string  `json:"source,omitempty"`
}

// UpdateIncomeResponse representa la respuesta al actualizar un ingreso
type UpdateIncomeResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
