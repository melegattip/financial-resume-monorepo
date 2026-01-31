package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims estructura para decodificar tokens del microservicio de usuarios
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware valida tokens JWT del microservicio de usuarios
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("🔍 [JWTAuth] Procesando request: %s %s", c.Request.Method, path)

		// Paths que no requieren autenticación
		publicPaths := []string{
			"/health",
			"/favicon.ico",
			"/robots.txt",
			"/manifest.json",
			"/swagger/",
			"/docs/",
			"/api/v1/gamification/action-types",
			"/api/v1/gamification/levels",
		}

		// Verificar si es un path público
		for _, publicPath := range publicPaths {
			if strings.HasPrefix(path, publicPath) {
				log.Printf("✅ [JWTAuth] Path público detectado: %s", path)
				c.Next()
				return
			}
		}

		// Extraer token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("❌ [JWTAuth] Authorization header faltante para: %s", path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		log.Printf("🔍 [JWTAuth] Authorization header encontrado para: %s", path)

		// Verificar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("❌ [JWTAuth] Formato de Authorization header inválido para: %s", path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		log.Printf("🔍 [JWTAuth] Token extraído, validando...")

		// Validar token JWT del microservicio de usuarios
		claims, err := validateUserServiceToken(token, jwtSecret)
		if err != nil {
			log.Printf("❌ [JWTAuth] Error validando token para %s: %v", path, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Printf("✅ [JWTAuth] Token válido para user_id: %d", claims.UserID)

		// Establecer información del usuario en el contexto
		userIDStr := strconv.FormatUint(uint64(claims.UserID), 10)
		c.Set("user_id", userIDStr)
		c.Set("user_email", claims.Email)
		c.Set("x_caller_id", userIDStr)

		log.Printf("✅ [JWTAuth] Contexto establecido - user_id: %s, email: %s", userIDStr, claims.Email)

		c.Next()
	}
}

// validateUserServiceToken valida un token JWT del microservicio de usuarios
func validateUserServiceToken(tokenString, jwtSecret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// OptionalJWTAuthMiddleware es un middleware opcional que no bloquea si no hay token
func OptionalJWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]

				if claims, err := validateUserServiceToken(token, jwtSecret); err == nil {
					userIDStr := strconv.FormatUint(uint64(claims.UserID), 10)
					c.Set("user_id", userIDStr)
					c.Set("user_email", claims.Email)
					c.Set("x_caller_id", userIDStr)
				}
			}
		}

		c.Next()
	}
}
