package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey tipo para keys del contexto
type ContextKey string

const (
	// UserIDKey key para el user ID en el contexto
	UserIDKey ContextKey = "user_id"
)

// getJWTSecret obtiene el JWT secret de las variables de entorno
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback para desarrollo - debe coincidir con el del users-service
		secret = "financial_resume_secret_key_2024"
	}
	return secret
}

// JWTAuthMiddleware middleware para validar tokens JWT
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar si viene el user_id desde el proxy (header X-User-ID)
		proxyUserID := r.Header.Get("X-User-ID")
		if proxyUserID != "" {
			log.Printf("✅ [JWTMiddleware] Using proxy user ID: %s", proxyUserID)
			// Si viene del proxy del FRE, confiar en el X-User-ID
			// El FRE ya validó el JWT, solo necesitamos el user_id
			ctx := context.WithValue(r.Context(), UserIDKey, proxyUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		log.Printf("🔍 [JWTMiddleware] No proxy user ID found, validating JWT directly")

		// Validación JWT directa (para llamadas directas al microservicio)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("❌ [JWTMiddleware] Missing Authorization header")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Verificar formato Bearer
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Printf("❌ [JWTMiddleware] Invalid token format - Bearer prefix missing")
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		// Parsear y validar el token
		jwtSecret := getJWTSecret()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Retornar la clave secreta desde variable de entorno
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("❌ [JWTMiddleware] Token parsing failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extraer claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userID string

			// Manejar user_id como string o número
			switch v := claims["user_id"].(type) {
			case string:
				userID = v
			case float64:
				userID = fmt.Sprintf("%.0f", v)
			case int:
				userID = fmt.Sprintf("%d", v)
			default:
				log.Printf("❌ [JWTMiddleware] Invalid user_id claim type in token: %T", claims["user_id"])
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			log.Printf("✅ [JWTMiddleware] Successfully validated JWT for user: %s", userID)

			// Agregar user ID al contexto
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Printf("❌ [JWTMiddleware] Token validation failed")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}

// GetUserIDFromContext extrae el user ID del contexto
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
