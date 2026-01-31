# 🔍 **GUÍA DE DIAGNÓSTICO - ERROR 500 EN CREACIÓN DE CATEGORÍAS**

## 📋 **RESUMEN DEL PROBLEMA**

**Error**: `POST https://financial-resume-engine.onrender.com/api/v1/categories 500 (Internal Server Error)`

**Síntomas**:
- ✅ Funciona correctamente en desarrollo local (puerto 8080)
- ❌ Error 500 en producción (Render.com)
- ❌ Error al crear categorías desde el frontend
- ✅ Crear incomes funciona correctamente

---

## 🔧 **SOLUCIONES IMPLEMENTADAS**

### **1. Logging Detallado Agregado**

Se ha agregado logging extensivo en los siguientes componentes:

#### **A. Handler de Categorías** (`internal/handlers/categories/handler.go`)
```go
log.Printf("🔍 [Categories] Iniciando creación de categoría")
log.Printf("🔍 [Categories] Request body recibido: %+v", requestBody)
log.Printf("🔍 [Categories] UserID extraído: %s", userID)
log.Printf("🔍 [Categories] Creando categoría con request: %+v", request)
log.Printf("✅ [Categories] Categoría creada exitosamente: %+v", category)
```

#### **B. Middleware JWT** (`internal/infrastructure/http/middleware/jwt_auth.go`)
```go
log.Printf("🔍 [JWTAuth] Procesando request: %s %s", c.Request.Method, path)
log.Printf("✅ [JWTAuth] Token válido para user_id: %d", claims.UserID)
log.Printf("✅ [JWTAuth] Contexto establecido - user_id: %s, email: %s", userIDStr, claims.Email)
```

#### **C. Repositorio de Categorías** (`internal/infrastructure/repository/category.go`)
```go
log.Printf("🔍 [CategoryRepo] Iniciando creación de categoría: %+v", category)
log.Printf("✅ [CategoryRepo] Categoría creada exitosamente: %s", category.ID)
```

### **2. Endpoint de Diagnóstico**

Nuevo endpoint público para verificar la configuración del sistema:

```bash
# Verificar configuración
curl https://financial-resume-engine.onrender.com/api/v1/config

# Verificar diagnóstico completo
curl https://financial-resume-engine.onrender.com/api/v1/diagnostics
```

### **3. Script de Prueba**

Script automatizado para verificar la funcionalidad:

```bash
./scripts/test_categories.sh
```

---

## 🚀 **PASOS PARA DIAGNOSTICAR**

### **Paso 1: Verificar Logs de Producción**

1. **Acceder a Render Dashboard**
   - Ir a: https://dashboard.render.com
   - Seleccionar el servicio `financial-resume-engine`
   - Ir a la pestaña "Logs"

2. **Buscar logs específicos**:
   ```
   [Categories] Iniciando creación de categoría
   [JWTAuth] Procesando request: POST /api/v1/categories
   [CategoryRepo] Iniciando creación de categoría
   ```

### **Paso 2: Ejecutar Diagnóstico**

```bash
# Desde tu máquina local
curl -s https://financial-resume-engine.onrender.com/api/v1/diagnostics | jq '.'

# O usar el script
./scripts/test_categories.sh
```

### **Paso 3: Verificar Variables de Entorno**

El endpoint de diagnóstico mostrará:
- ✅ Variables de entorno críticas
- ✅ Estado de la base de datos
- ✅ Configuración de servicios externos

---

## 🎯 **POSIBLES CAUSAS Y SOLUCIONES**

### **Causa 1: Problema de Autenticación JWT**

**Síntomas**:
- Logs muestran "Error obteniendo user_id"
- Token inválido o mal formado

**Solución**:
1. Verificar que `JWT_SECRET` esté configurado correctamente en Render
2. Verificar que el token del frontend sea válido
3. Comprobar que el middleware JWT esté funcionando

### **Causa 2: Problema de Base de Datos**

**Síntomas**:
- Logs muestran "Error creando categoría"
- Error de conexión a PostgreSQL

**Solución**:
1. Verificar variables de entorno de base de datos:
   ```env
   DATABASE_URL=postgresql://user:pass@host:port/db
   DB_HOST=your-host
   DB_USER=your-user
   DB_PASSWORD=your-password
   DB_NAME=financial_resume
   DB_SSLMODE=require
   ```

2. Verificar que la tabla `categories` exista:
   ```sql
   \dt categories;
   ```

### **Causa 3: Problema de Gamificación**

**Síntomas**:
- Error después de crear la categoría exitosamente
- Logs muestran "Error registrando acción de gamificación"

**Solución**:
1. Verificar que `GAMIFICATION_SERVICE_URL` esté configurado
2. Verificar que el servicio de gamificación esté funcionando
3. El error de gamificación no debería fallar la operación principal

### **Causa 4: Problema de Configuración**

**Síntomas**:
- Variables de entorno faltantes
- URLs de servicios incorrectas

**Solución**:
1. Verificar todas las variables de entorno en Render
2. Verificar que los servicios externos estén accesibles

---

## 🔄 **PROCESO DE DEPLOYMENT**

### **1. Hacer Commit y Push**

```bash
git add .
git commit -m "🔧 Add detailed logging for categories debugging"
git push origin main
```

### **2. Verificar Deployment en Render**

1. Ir a Render Dashboard
2. Verificar que el deployment se completó exitosamente
3. Revisar logs del deployment

### **3. Probar en Producción**

1. **Ejecutar diagnóstico**:
   ```bash
   curl https://financial-resume-engine.onrender.com/api/v1/diagnostics
   ```

2. **Intentar crear categoría** desde el frontend

3. **Revisar logs** en tiempo real en Render Dashboard

---

## 📊 **MÉTRICAS DE MONITOREO**

### **Logs a Monitorear**

```bash
# Logs de éxito
✅ [Categories] Categoría creada exitosamente
✅ [JWTAuth] Token válido para user_id
✅ [CategoryRepo] Categoría creada exitosamente

# Logs de error
❌ [Categories] Error obteniendo user_id
❌ [Categories] Error en service.Create
❌ [CategoryRepo] Error creando categoría
❌ [JWTAuth] Error validando token
```

### **Endpoints de Health Check**

```bash
# Health general
curl https://financial-resume-engine.onrender.com/health

# Configuración
curl https://financial-resume-engine.onrender.com/api/v1/config

# Diagnóstico completo
curl https://financial-resume-engine.onrender.com/api/v1/diagnostics
```

---

## 🎯 **PRÓXIMOS PASOS**

1. **Deploy los cambios** con logging detallado
2. **Ejecutar el script de diagnóstico**
3. **Intentar crear una categoría** desde el frontend
4. **Revisar logs** en Render Dashboard
5. **Identificar el punto exacto** donde falla
6. **Aplicar la solución específica** según el diagnóstico

---

## 📞 **CONTACTO**

Si el problema persiste después de seguir esta guía:

1. **Compartir logs** de producción (sin información sensible)
2. **Compartir resultado** del endpoint de diagnóstico
3. **Describir pasos exactos** que reproducen el error

---

*Documento creado: Enero 2025*  
*Versión: 1.0* 