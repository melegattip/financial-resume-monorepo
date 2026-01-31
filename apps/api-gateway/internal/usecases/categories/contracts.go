package categories

// CreateCategoryRequest representa el request para crear una categoría
type CreateCategoryRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name" binding:"required"`
}

// CreateCategoryResponse representa la respuesta al crear una categoría
type CreateCategoryResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// ListCategoriesResponse representa la respuesta al listar categorías
type ListCategoriesResponse struct {
	Categories []GetCategoryResponse `json:"categories"`
}

// GetCategoryResponse representa la respuesta al obtener una categoría
type GetCategoryResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// UpdateCategoryRequest representa el request para actualizar una categoría
type UpdateCategoryRequest struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name" binding:"required"`     // Nombre actual de la categoría
	NewName string `json:"new_name" binding:"required"` // Nuevo nombre para la categoría
}

// UpdateCategoryResponse representa la respuesta al actualizar una categoría
type UpdateCategoryResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// DeleteCategoryRequest representa el request para eliminar una categoría
type DeleteCategoryRequest struct {
	UserID string `json:"user_id"`
	ID     string `json:"id" binding:"required"`
}

// DeleteCategoryResponse representa la respuesta al eliminar una categoría
type DeleteCategoryResponse struct {
	Message string `json:"message"`
}
