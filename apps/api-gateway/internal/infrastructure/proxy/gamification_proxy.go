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

// GamificationProxy maneja el proxy hacia el microservicio de gamificación
type GamificationProxy struct {
	gamificationServiceURL string
	httpClient             *http.Client
}

// NewGamificationProxy crea una nueva instancia del proxy
func NewGamificationProxy(serviceURL string) *GamificationProxy {
	log.Printf("🎮 [GamificationProxy] Inicializando proxy con URL: %s", serviceURL)

	proxy := &GamificationProxy{
		gamificationServiceURL: serviceURL,
		httpClient:             &http.Client{},
	}

	// Verificar conectividad al inicializar
	go func() {
		if err := proxy.HealthCheck(); err != nil {
			log.Printf("⚠️ [GamificationProxy] Servicio de gamificación no disponible al inicializar: %v", err)
		} else {
			log.Printf("✅ [GamificationProxy] Servicio de gamificación disponible en: %s", serviceURL)
		}
	}()

	log.Printf("🎮 [GamificationProxy] Proxy inicializado exitosamente")
	return proxy
}

// ProxyRequest redirige la petición al microservicio de gamificación (requiere autenticación)
func (p *GamificationProxy) ProxyRequest(c *gin.Context) {
	// Remover el prefijo /api/v1 para evitar duplicación
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/v1")

	// Construir URL del microservicio
	targetURL := p.gamificationServiceURL + path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Logging para diagnóstico
	log.Printf("🔧 [GamificationProxy] ProxyRequest - Original Path: %s", c.Request.URL.Path)
	log.Printf("🔧 [GamificationProxy] ProxyRequest - Processed Path: %s", path)
	log.Printf("🔧 [GamificationProxy] ProxyRequest - Target URL: %s", targetURL)
	log.Printf("🔧 [GamificationProxy] Service URL configurada: %s", p.gamificationServiceURL)

	// Leer el body de la petición original
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Crear nueva petición
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("❌ [GamificationProxy] Error creando petición: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando petición"})
		return
	}

	// Copiar headers importantes
	p.copyHeaders(c.Request, req)

	// Extraer y agregar user_id al JWT token
	if err := p.enrichJWTWithUserID(c, req); err != nil {
		log.Printf("❌ [GamificationProxy] Error procesando autenticación: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error procesando autenticación"})
		return
	}

	// Ejecutar la petición
	resp, err := p.httpClient.Do(req)
	if err != nil {
		// Mejor logging y diagnóstico
		log.Printf("❌ [GamificationProxy] Error conectando con servicio de gamificación: %v", err)
		log.Printf("🔧 [GamificationProxy] Target URL: %s", targetURL)
		log.Printf("🔧 [GamificationProxy] Service URL configurada: %s", p.gamificationServiceURL)

		c.JSON(http.StatusBadGateway, gin.H{
			"error":       "Error conectando con servicio de gamificación",
			"details":     err.Error(),
			"target_url":  targetURL,
			"service_url": p.gamificationServiceURL,
		})
		return
	}
	defer resp.Body.Close()

	// Logging del código de respuesta
	log.Printf("🔧 [GamificationProxy] Respuesta del servicio: %d", resp.StatusCode)

	// Leer respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [GamificationProxy] Error leyendo respuesta: %v", err)
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

// ProxyPublicRequest redirige la petición al microservicio de gamificación (sin autenticación)
func (p *GamificationProxy) ProxyPublicRequest(c *gin.Context) {
	// Remover el prefijo /api/v1 de la ruta para evitar duplicación
	path := strings.TrimPrefix(c.Request.URL.Path, "/api/v1")

	// Construir URL del microservicio
	targetURL := p.gamificationServiceURL + path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Logging para diagnóstico
	log.Printf("🔧 [GamificationProxy] ProxyPublicRequest - Original Path: %s", c.Request.URL.Path)
	log.Printf("🔧 [GamificationProxy] ProxyPublicRequest - Processed Path: %s", path)
	log.Printf("🔧 [GamificationProxy] ProxyPublicRequest - Target URL: %s", targetURL)
	log.Printf("🔧 [GamificationProxy] Service URL configurada: %s", p.gamificationServiceURL)
	log.Printf("🔧 [GamificationProxy] ProxyPublicRequest - Method: %s", c.Request.Method)

	// Leer el body de la petición original
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Crear nueva petición
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("❌ [GamificationProxy] Error creando petición pública: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando petición"})
		return
	}

	// Copiar headers importantes (sin autenticación)
	p.copyHeaders(c.Request, req)

	// Ejecutar la petición
	resp, err := p.httpClient.Do(req)
	if err != nil {
		// Mejor logging y diagnóstico
		log.Printf("❌ [GamificationProxy] Error conectando con servicio de gamificación (público): %v", err)
		log.Printf("🔧 [GamificationProxy] Target URL: %s", targetURL)
		log.Printf("🔧 [GamificationProxy] Service URL configurada: %s", p.gamificationServiceURL)

		c.JSON(http.StatusBadGateway, gin.H{
			"error":       "Error conectando con servicio de gamificación",
			"details":     err.Error(),
			"target_url":  targetURL,
			"service_url": p.gamificationServiceURL,
		})
		return
	}
	defer resp.Body.Close()

	// Logging del código de respuesta
	log.Printf("🔧 [GamificationProxy] Respuesta del servicio (público): %d", resp.StatusCode)

	// Leer respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [GamificationProxy] Error leyendo respuesta pública: %v", err)
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
func (p *GamificationProxy) copyHeaders(source *http.Request, target *http.Request) {
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
func (p *GamificationProxy) enrichJWTWithUserID(c *gin.Context, req *http.Request) error {
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

	// Para simplificar, vamos a crear un nuevo JWT que incluya el user_id como string
	// En una implementación más robusta, podrías modificar el token existente
	// Por ahora, vamos a pasar el user_id en un header personalizado y modificar el microservicio

	// Opción 1: Pasar el user_id en un header personalizado (más simple)
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("Authorization", authHeader)

	return nil
}

// RegisterRoutes registra las rutas del proxy en el router
func (p *GamificationProxy) RegisterRoutes(api *gin.RouterGroup) {
	// Todas las rutas de gamificación van al proxy
	gamificationRoutes := api.Group("/gamification")
	{
		// Usar el proxy para todas las rutas
		gamificationRoutes.Any("/*path", p.ProxyRequest)
		gamificationRoutes.GET("", p.ProxyRequest)  // Para /gamification sin path
		gamificationRoutes.POST("", p.ProxyRequest) // Para /gamification sin path
	}
}

// HealthCheck verifica que el microservicio de gamificación esté disponible
func (p *GamificationProxy) HealthCheck() error {
	// Construir URL de health check correctamente
	healthURL := strings.TrimSuffix(p.gamificationServiceURL, "/api/v1") + "/health"
	resp, err := p.httpClient.Get(healthURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gamification service unhealthy: status %d", resp.StatusCode)
	}

	return nil
}
