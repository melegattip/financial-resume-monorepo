package update

// UpdateCategoryRequest representa el request para actualizar una categoría
type UpdateCategoryRequest struct {
	UserID  string `json:"user_id"`
	ID      string `json:"id"`                          // ID de la categoría
	NewName string `json:"new_name" binding:"required"` // Nuevo nombre para la categoría
}

// UpdateCategoryResponse representa la respuesta al actualizar una categoría
type UpdateCategoryResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}
