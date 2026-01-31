package categories

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/handlers/categories/create"
	"github.com/melegattip/financial-resume-engine/internal/handlers/categories/delete"
	"github.com/melegattip/financial-resume-engine/internal/handlers/categories/get"
	"github.com/melegattip/financial-resume-engine/internal/handlers/categories/list"
	"github.com/melegattip/financial-resume-engine/internal/handlers/categories/update"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
)

// Handler maneja las peticiones HTTP relacionadas con categorías
type Handler struct {
	createHandler      *create.Handler
	getHandler         *get.Handler
	listHandler        *list.Handler
	updateHandler      *update.Handler
	deleteHandler      *delete.Handler
	service            categories.CategoryService
	gamificationHelper *services.GamificationHelper
}

// NewHandler crea una nueva instancia de Handler
func NewHandler(service categories.CategoryService, gamificationHelper *services.GamificationHelper) *Handler {
	return &Handler{
		createHandler:      create.NewHandler(service),
		getHandler:         get.NewHandler(service),
		listHandler:        list.NewHandler(service),
		updateHandler:      update.NewHandler(service),
		deleteHandler:      delete.NewHandler(service),
		service:            service,
		gamificationHelper: gamificationHelper,
	}
}

// CreateCategoryRequestBody representa solo los datos necesarios en el body
type CreateCategoryRequestBody struct {
	Name string `json:"name" binding:"required"`
}

// getUserIDFromContext extrae el user_id del contexto JWT
func (h *Handler) getUserIDFromContext(c *gin.Context) (string, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return "", coreErrors.NewUnauthorizedRequest("Usuario no autenticado")
	}

	// Convertir a string según el tipo
	switch v := userIDInterface.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case string:
		return v, nil
	default:
		return "", coreErrors.NewBadRequest("Formato de user_id inválido")
	}
}

// CreateCategory maneja la creación de una nueva categoría con integración de gamificación
func (h *Handler) CreateCategory(c *gin.Context) {
	log.Printf("🔍 [Categories] Iniciando creación de categoría")

	var requestBody CreateCategoryRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("❌ [Categories] Error binding JSON: %v", err)
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Invalid request body"))
		return
	}

	log.Printf("🔍 [Categories] Request body recibido: %+v", requestBody)

	// Obtener userID desde el contexto JWT
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		log.Printf("❌ [Categories] Error obteniendo user_id: %v", err)
		httpUtil.HandleError(c, err)
		return
	}

	log.Printf("🔍 [Categories] UserID extraído: %s", userID)

	if requestBody.Name == "" {
		log.Printf("❌ [Categories] Nombre de categoría vacío")
		httpUtil.BadRequest(c, coreErrors.NewBadRequest("Category name is required"))
		return
	}

	// Crear la categoría usando el servicio
	request := categories.CreateCategoryRequest{
		Name:   requestBody.Name,
		UserID: userID,
	}

	log.Printf("🔍 [Categories] Creando categoría con request: %+v", request)

	category, err := h.service.Create(request)
	if err != nil {
		log.Printf("❌ [Categories] Error en service.Create: %v", err)
		httpUtil.HandleError(c, err)
		return
	}

	log.Printf("✅ [Categories] Categoría creada exitosamente: %+v", category)

	// 🎮 REGISTRAR GAMIFICACIÓN después de la creación exitosa
	if h.gamificationHelper != nil {
		log.Printf("🎮 [Categories] Registrando creación de categoría para gamificación: %s (ID: %s)", category.Name, category.ID)

		err = h.gamificationHelper.RecordAction(
			userID,
			"create_category",
			"category",
			category.ID,
			"Created category: "+category.Name,
		)
		if err != nil {
			// Log el error pero no fallar la operación principal
			log.Printf("⚠️ [Categories] Error registrando acción de gamificación: %v", err)
		} else {
			log.Printf("✅ [Categories] Acción de gamificación registrada exitosamente para categoría: %s", category.Name)
		}
	} else {
		log.Printf("⚠️ [Categories] gamificationHelper es nil, no se puede registrar acción")
	}

	log.Printf("✅ [Categories] Enviando respuesta exitosa")
	httpUtil.SendSuccess(c, http.StatusCreated, category, "Category created successfully")
}

// ListCategories maneja la lista de categorías
func (h *Handler) ListCategories(c *gin.Context) {
	h.listHandler.Handle(c)
}

// GetCategory maneja la obtención de una categoría específica
func (h *Handler) GetCategory(c *gin.Context) {
	h.getHandler.Handle(c)
}

// UpdateCategory maneja la actualización de una categoría
func (h *Handler) UpdateCategory(c *gin.Context) {
	h.updateHandler.Handle(c)
}

// DeleteCategory maneja la eliminación de una categoría
func (h *Handler) DeleteCategory(c *gin.Context) {
	h.deleteHandler.Handle(c)
}
