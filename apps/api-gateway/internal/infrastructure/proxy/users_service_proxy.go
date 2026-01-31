package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UsersServiceProxy maneja el proxy hacia el microservicio de usuarios
type UsersServiceProxy struct {
	usersServiceURL string
	httpClient      *http.Client
}

// NewUsersServiceProxy crea una nueva instancia del proxy
func NewUsersServiceProxy(serviceURL string) *UsersServiceProxy {
	log.Printf("👤 [UsersServiceProxy] Inicializando proxy con URL: %s", serviceURL)

	proxy := &UsersServiceProxy{
		usersServiceURL: serviceURL,
		httpClient:      &http.Client{},
	}

	// Verificar conectividad al inicializar
	go func() {
		if err := proxy.HealthCheck(); err != nil {
			log.Printf("⚠️ [UsersServiceProxy] Servicio de usuarios no disponible al inicializar: %v", err)
		} else {
			log.Printf("✅ [UsersServiceProxy] Servicio de usuarios disponible en: %s", serviceURL)
		}
	}()

	return proxy
}

// ProxyRequest redirige la petición al microservicio de usuarios (requiere autenticación)
func (p *UsersServiceProxy) ProxyRequest(c *gin.Context) {
	// Mapear rutas del FRE a rutas del servicio de usuarios
	targetURL := p.mapRoute(c.Request.URL.Path)
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Logging para diagnóstico
	log.Printf("🔧 [UsersServiceProxy] ProxyRequest - Target URL: %s", targetURL)
	log.Printf("🔧 [UsersServiceProxy] Service URL configurada: %s", p.usersServiceURL)

	// Leer el body de la petición original
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Crear nueva petición
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("❌ [UsersServiceProxy] Error creando petición: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando petición"})
		return
	}

	// Copiar headers importantes
	p.copyHeaders(c.Request, req)

	// Extraer y agregar user_id al JWT token
	if err := p.enrichJWTWithUserID(c, req); err != nil {
		log.Printf("❌ [UsersServiceProxy] Error procesando autenticación: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error procesando autenticación"})
		return
	}

	// Ejecutar la petición
	resp, err := p.httpClient.Do(req)
	if err != nil {
		// Mejor logging y diagnóstico
		log.Printf("❌ [UsersServiceProxy] Error conectando con servicio de usuarios: %v", err)
		log.Printf("🔧 [UsersServiceProxy] Target URL: %s", targetURL)
		log.Printf("🔧 [UsersServiceProxy] Service URL configurada: %s", p.usersServiceURL)

		c.JSON(http.StatusBadGateway, gin.H{
			"error":       "Error conectando con servicio de usuarios",
			"details":     err.Error(),
			"target_url":  targetURL,
			"service_url": p.usersServiceURL,
		})
		return
	}
	defer resp.Body.Close()

	// Logging del código de respuesta
	log.Printf("🔧 [UsersServiceProxy] Respuesta del servicio: %d", resp.StatusCode)

	// Leer respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [UsersServiceProxy] Error leyendo respuesta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo respuesta"})
		return
	}

	// Copiar headers de respuesta
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Enviar respuesta
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
}

// mapRoute mapea las rutas del FRE a las rutas del servicio de usuarios
func (p *UsersServiceProxy) mapRoute(path string) string {
	// Mapear /api/v1/auth/* a /api/v1/auth/*
	if strings.HasPrefix(path, "/api/v1/auth/") {
		return fmt.Sprintf("%s%s", p.usersServiceURL, path)
	}

	// Mapear /api/v1/users/* a /api/v1/users/*
	if strings.HasPrefix(path, "/api/v1/users/") {
		return fmt.Sprintf("%s%s", p.usersServiceURL, path)
	}

	// Para rutas sin prefijo /api/v1 (legacy), agregar /api/v1/users/
	if strings.HasPrefix(path, "/users/") {
		userPath := strings.TrimPrefix(path, "/users/")
		return fmt.Sprintf("%s/api/v1/users/%s", p.usersServiceURL, userPath)
	}

	// URL por defecto
	return fmt.Sprintf("%s%s", p.usersServiceURL, path)
}

// ProxyPublicRequest redirige la petición al microservicio de usuarios (sin autenticación)
func (p *UsersServiceProxy) ProxyPublicRequest(c *gin.Context) {
	// Mapear rutas del FRE a rutas del servicio de usuarios
	targetURL := p.mapRoute(c.Request.URL.Path)
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Logging para diagnóstico
	log.Printf("🔧 [UsersServiceProxy] ProxyPublicRequest - Target URL: %s", targetURL)
	log.Printf("🔧 [UsersServiceProxy] Service URL configurada: %s", p.usersServiceURL)

	// Leer el body de la petición original
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Crear nueva petición
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("❌ [UsersServiceProxy] Error creando petición pública: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando petición"})
		return
	}

	// Copiar headers importantes (sin autenticación)
	p.copyHeaders(c.Request, req)

	// Ejecutar la petición
	resp, err := p.httpClient.Do(req)
	if err != nil {
		// Mejor logging y diagnóstico
		log.Printf("❌ [UsersServiceProxy] Error conectando con servicio de usuarios (público): %v", err)
		log.Printf("🔧 [UsersServiceProxy] Target URL: %s", targetURL)
		log.Printf("🔧 [UsersServiceProxy] Service URL configurada: %s", p.usersServiceURL)

		c.JSON(http.StatusBadGateway, gin.H{
			"error":       "Error conectando con servicio de usuarios",
			"details":     err.Error(),
			"target_url":  targetURL,
			"service_url": p.usersServiceURL,
		})
		return
	}
	defer resp.Body.Close()

	// Logging del código de respuesta
	log.Printf("🔧 [UsersServiceProxy] Respuesta del servicio (público): %d", resp.StatusCode)

	// Leer respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [UsersServiceProxy] Error leyendo respuesta pública: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo respuesta"})
		return
	}

	// Copiar headers de respuesta
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Enviar respuesta
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), responseBody)
}

// copyHeaders copia headers importantes de la petición original
func (p *UsersServiceProxy) copyHeaders(source *http.Request, target *http.Request) {
	headersToProxy := []string{
		"Content-Type",
		"Accept",
		"User-Agent",
		"X-Forwarded-For",
		"X-Real-IP",
	}

	for _, header := range headersToProxy {
		if value := source.Header.Get(header); value != "" {
			target.Header.Set(header, value)
		}
	}
}

// enrichJWTWithUserID extrae el user_id del contexto de Gin y lo agrega al JWT para el microservicio
func (p *UsersServiceProxy) enrichJWTWithUserID(c *gin.Context, req *http.Request) error {
	// Obtener el token JWT original
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("missing authorization header")
	}

	// Obtener user_id del contexto de Gin (ya extraído por el middleware JWT del FRE)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return fmt.Errorf("user_id not found in context")
	}

	// Convertir user_id a string
	var userID string
	switch v := userIDInterface.(type) {
	case uint:
		userID = strconv.FormatUint(uint64(v), 10)
	case string:
		userID = v
	default:
		return fmt.Errorf("invalid user_id format")
	}

	// Pasar el user_id en un header personalizado
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("Authorization", authHeader)

	return nil
}

// RegisterRoutes registra las rutas del proxy en el router
func (p *UsersServiceProxy) RegisterRoutes(router *gin.Engine) {
	// Rutas públicas de usuarios (no requieren autenticación)
	publicUsersGroup := router.Group("/api/v1/users")
	{
		publicUsersGroup.POST("/register", p.ProxyPublicRequest)
		publicUsersGroup.POST("/login", p.ProxyPublicRequest)
		publicUsersGroup.POST("/refresh", p.ProxyPublicRequest)
		publicUsersGroup.GET("/verify-email/:token", p.ProxyPublicRequest)
		publicUsersGroup.POST("/request-password-reset", p.ProxyPublicRequest)
		publicUsersGroup.POST("/reset-password", p.ProxyPublicRequest)
	}

	// Rutas protegidas de usuarios (requieren autenticación)
	// Estas se manejarán en el router principal con middleware JWT
}

// HealthCheck verifica que el microservicio de usuarios esté disponible
func (p *UsersServiceProxy) HealthCheck() error {
	resp, err := p.httpClient.Get(p.usersServiceURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users service unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
