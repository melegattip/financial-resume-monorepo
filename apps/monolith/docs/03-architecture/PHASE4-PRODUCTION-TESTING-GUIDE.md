# Guía de Pruebas en Producción - Fase 4
**Testing del Monolith Deployado**

**Fecha de Creación**: 2026-02-13  
**Ambiente**: Producción (Monolith en Render + Nueva DB en Supabase)  
**Tiempo Estimado**: 1-2 horas

---

## 🎯 Objetivo

Realizar pruebas exhaustivas del monolith deployado en producción **antes** de redirigir tráfico de usuarios reales. Esta fase asegura que:

- ✅ Todos los endpoints críticos funcionan correctamente
- ✅ La conexión a la nueva base de datos es estable
- ✅ Los datos migrados son accesibles y correctos
- ✅ No hay errores críticos en logs
- ✅ El rendimiento es aceptable
- ✅ La autenticación y autorización funcionan

---

## 📋 Pre-Requisitos

### Estado Actual (Completado en Fase 3)
- ✅ Base de datos migrada a Supabase (nueva DB)
- ✅ Validación de datos completada (PASS)
- ✅ Monolith deployado en Render
- ✅ Variables de entorno configuradas en producción
- ✅ Health check básico funcionando

### Herramientas Necesarias
- **Postman** o **Insomnia** (para testing de APIs)
- **curl** (línea de comandos)
- **Navegador web** (para testing del frontend)
- **Acceso a Render Dashboard** (para revisar logs)
- **Acceso a Supabase Dashboard** (para queries manuales)

### URLs de Referencia
```bash
# Monolith en producción (ejemplo)
MONOLITH_URL=https://financial-resume-monolith.onrender.com

# Frontend en producción (ejemplo)
FRONTEND_URL=https://financial-resume.vercel.app

# Legacy API Gateway (para comparación)
LEGACY_URL=https://api-gateway-legacy.onrender.com
```

---

## 🧪 FASE 4.1: Pruebas de Infraestructura (15 min)

### 4.1.1 Health Check Completo

```bash
# 1. Health check básico
curl https://financial-resume-monolith.onrender.com/health

# Respuesta esperada:
# {
#   "status": "ok",
#   "database": "connected",
#   "timestamp": "2026-02-13T19:21:30Z"
# }
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ `status: "ok"`
- ✅ `database: "connected"`

### 4.1.2 Verificar Conexión a Base de Datos

```bash
# Endpoint de diagnóstico de DB (si existe)
curl https://financial-resume-monolith.onrender.com/internal/db-status

# O verificar manualmente en Supabase
```

**En Supabase Dashboard**:
1. Ve a tu proyecto → Logs
2. Filtra por "Connections"
3. Verifica que hay conexiones activas desde el monolith

### 4.1.3 Revisar Logs de Startup

**En Render Dashboard**:
1. Ir a `financial-resume-monolith` → Logs
2. Filtrar logs de los últimos 15 minutos
3. Buscar mensajes de inicio:

```log
[INFO] Server starting on port 8080
[INFO] Connected to database: postgresql://...
[INFO] Migrations up to date
[INFO] Server ready to accept connections
```

**Verificar que NO hay**:
- ❌ `[ERROR]` durante el startup
- ❌ `connection refused`
- ❌ `panic` o `fatal`

---

## 🔐 FASE 4.2: Pruebas de Autenticación (30 min)

### 4.2.1 Registro de Usuario de Prueba

**Opción A: Usando curl**

```bash
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-phase4@example.com",
    "password": "TestPass123!",
    "name": "Phase 4 Test User"
  }'
```

**Respuesta esperada**:
```json
{
  "user": {
    "id": "uuid-generated",
    "email": "test-phase4@example.com",
    "name": "Phase 4 Test User",
    "created_at": "2026-02-13T19:21:30Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Verificar**:
- ✅ Status code: `201 Created`
- ✅ Usuario tiene un `id` válido (UUID)
- ✅ Token JWT retornado
- ✅ No hay errores en logs

**Opción B: Usando Postman**

1. Crear nueva request: `POST` → `{{MONOLITH_URL}}/api/v1/auth/register`
2. Headers: `Content-Type: application/json`
3. Body (raw JSON):
```json
{
  "email": "test-phase4@example.com",
  "password": "TestPass123!",
  "name": "Phase 4 Test User"
}
```
4. Send y guardar el token retornado

### 4.2.2 Login con Usuario de Prueba

```bash
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-phase4@example.com",
    "password": "TestPass123!"
  }'
```

**Respuesta esperada**:
```json
{
  "user": {
    "id": "uuid-generated",
    "email": "test-phase4@example.com",
    "name": "Phase 4 Test User"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ Token JWT retornado
- ✅ Datos de usuario correctos

**⚠️ GUARDAR EL TOKEN** para las siguientes pruebas:
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 4.2.3 Validación de Token JWT

```bash
# Obtener perfil del usuario autenticado
curl https://financial-resume-monolith.onrender.com/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta esperada**:
```json
{
  "id": "uuid-generated",
  "email": "test-phase4@example.com",
  "name": "Phase 4 Test User",
  "created_at": "2026-02-13T19:21:30Z",
  "updated_at": "2026-02-13T19:21:30Z"
}
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ Datos de usuario retornados
- ✅ No hay error de autenticación

### 4.2.4 Pruebas de Casos de Error

**Test 1: Login con credenciales inválidas**
```bash
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-phase4@example.com",
    "password": "WrongPassword123!"
  }'
```

**Verificar**:
- ✅ Status code: `401 Unauthorized`
- ✅ Mensaje de error apropiado: `"invalid credentials"`

**Test 2: Acceso sin token**
```bash
curl https://financial-resume-monolith.onrender.com/api/v1/users/me
```

**Verificar**:
- ✅ Status code: `401 Unauthorized`
- ✅ Mensaje de error: `"missing or invalid token"`

**Test 3: Token inválido**
```bash
curl https://financial-resume-monolith.onrender.com/api/v1/users/me \
  -H "Authorization: Bearer invalid-token-12345"
```

**Verificar**:
- ✅ Status code: `401 Unauthorized`
- ✅ Mensaje de error: `"invalid token"`

---

## 🎮 FASE 4.3: Pruebas de Datos de Gamificación (30 min)

### 4.3.1 Obtener Perfil de Gamificación

```bash
curl https://financial-resume-monolith.onrender.com/api/v1/gamification/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta esperada (usuario nuevo)**:
```json
{
  "user_id": "uuid-generated",
  "level": 1,
  "experience_points": 0,
  "total_quests_completed": 0,
  "total_achievements_unlocked": 0,
  "rank": "Novice",
  "streak_days": 0
}
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ Perfil de gamificación creado automáticamente
- ✅ Valores iniciales correctos

### 4.3.2 Obtener Achievements

```bash
curl https://financial-resume-monolith.onrender.com/api/v1/gamification/achievements \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta esperada**:
```json
{
  "achievements": [],
  "total": 0
}
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ Lista vacía (usuario nuevo sin achievements)

### 4.3.3 Prueba con Usuario Migrado (Datos Legacy)

**IMPORTANTE**: Usar credenciales de un usuario real que existía en la base legacy

```bash
# Login con usuario migrado
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario-real@example.com",
    "password": "PasswordReal123!"
  }'

# Guardar token migrado
export TOKEN_MIGRATED="..."

# Obtener perfil de gamificación
curl https://financial-resume-monolith.onrender.com/api/v1/gamification/profile \
  -H "Authorization: Bearer $TOKEN_MIGRATED"
```

**Respuesta esperada (usuario migrado)**:
```json
{
  "user_id": "uuid-migrated",
  "level": 5,
  "experience_points": 1250,
  "total_quests_completed": 12,
  "total_achievements_unlocked": 3,
  "rank": "Adventurer",
  "streak_days": 7
}
```

**Verificar que los datos migrados son correctos**:
- ✅ `level` corresponde a datos legacy
- ✅ `experience_points` corresponde a datos legacy
- ✅ `total_quests_completed` correcto
- ✅ No hay pérdida de datos

### 4.3.4 Comparar con Legacy API (Verificación Cruzada)

**Obtener datos del mismo usuario desde el legacy** (para comparación):

```bash
curl https://api-gateway-legacy.onrender.com/api/gamification/profile \
  -H "Authorization: Bearer $TOKEN_LEGACY"
```

**Comparar**:
- ✅ Los valores deben ser idénticos entre monolith y legacy
- ✅ No hay discrepancias en datos críticos

---

## 📊 FASE 4.4: Pruebas de Endpoints Críticos (20 min)

### 4.4.1 CRUD de Productos (Si aplica)

**Crear producto**:
```bash
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product Phase 4",
    "description": "Testing product creation",
    "price": 99.99,
    "stock": 100
  }'
```

**Verificar**:
- ✅ Status code: `201 Created`
- ✅ Producto creado con ID válido

**Listar productos**:
```bash
curl https://financial-resume-monolith.onrender.com/api/v1/products \
  -H "Authorization: Bearer $TOKEN"
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ Lista de productos retornada

**Actualizar producto**:
```bash
curl -X PUT https://financial-resume-monolith.onrender.com/api/v1/products/{id} \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product Name"
  }'
```

**Verificar**:
- ✅ Status code: `200 OK`
- ✅ Producto actualizado correctamente

**Eliminar producto**:
```bash
curl -X DELETE https://financial-resume-monolith.onrender.com/api/v1/products/{id} \
  -H "Authorization: Bearer $TOKEN"
```

**Verificar**:
- ✅ Status code: `204 No Content` o `200 OK`
- ✅ Producto eliminado

### 4.4.2 Otros Endpoints Críticos

**Adapta según tu aplicación**:
- `GET /api/v1/transactions`
- `POST /api/v1/transactions`
- `GET /api/v1/dashboards/overview`
- `GET /api/v1/reports/summary`
- etc.

---

## 🌐 FASE 4.5: Pruebas de Frontend Integrado (20 min)

### 4.5.1 Configurar Frontend Temporal

**Opción A: Modificar variable de entorno localmente**

```bash
# En tu proyecto de frontend
export VITE_API_URL=https://financial-resume-monolith.onrender.com

# Iniciar frontend
npm run dev
```

**Opción B: Deployment temporal de frontend**

Deployar una versión del frontend que apunte al monolith (rama temporal):
```bash
# En Vercel
vercel --env VITE_API_URL=https://financial-resume-monolith.onrender.com
```

### 4.5.2 Pruebas de Flujo Completo (User Journey)

**Test 1: Registro → Login → Dashboard**

1. Abrir el frontend en el navegador
2. Ir a `/register`
3. Registrar nuevo usuario: `frontend-test@example.com`
4. Verificar redirección automática al dashboard
5. Verificar que se muestra el perfil del usuario
6. Verificar que se muestran datos de gamificación (iniciales)

**Verificar**:
- ✅ Registro exitoso
- ✅ Login automático
- ✅ Dashboard carga sin errores
- ✅ No hay errores en consola del navegador
- ✅ No hay errores en Network tab

**Test 2: Login con Usuario Migrado**

1. Cerrar sesión
2. Ir a `/login`
3. Login con usuario migrado: `usuario-real@example.com`
4. Verificar que se muestra el perfil correcto
5. Verificar que se muestran datos históricos (achievements, quests, etc.)
6. Navegar por todas las secciones críticas

**Verificar**:
- ✅ Login exitoso
- ✅ Datos migrados visibles
- ✅ No hay pérdida de información
- ✅ UI renderiza correctamente

**Test 3: Crear/Editar Datos**

1. Crear una nueva transacción/producto/etc.
2. Verificar que se guarda correctamente
3. Editar el registro creado
4. Verificar que se actualiza correctamente
5. Eliminar el registro
6. Verificar que se elimina correctamente

**Verificar**:
- ✅ CRUD completo funciona
- ✅ Feedback visual correcto
- ✅ No hay errores de validación inesperados

### 4.5.3 Pruebas de Performance en Frontend

**En Chrome DevTools**:

1. Abrir DevTools → Network tab
2. Refresh la página del dashboard
3. Observar tiempos de carga

**Verificar**:
- ✅ Tiempo de carga inicial: `< 2 segundos`
- ✅ Tiempo de respuesta de API: `< 500ms` (por request)
- ✅ No hay requests fallidos (status 500, 404, etc.)

**En Chrome DevTools → Lighthouse**:

1. Ejecutar audit de Performance
2. Verificar score

**Objetivo**:
- ✅ Performance score: `> 80`
- ✅ No hay errores críticos en consola

---

## 🔍 FASE 4.6: Auditoría de Logs y Monitoreo (15 min)

### 4.6.1 Revisar Logs del Monolith

**En Render Dashboard**:

1. Ir a `financial-resume-monolith` → Logs
2. Filtrar logs de los últimos 30 minutos (periodo de testing)
3. Buscar patrones de error

**Buscar**:
- ❌ `[ERROR]` - Errores críticos
- ❌ `[WARN]` - Warnings que podrían ser problemas
- ❌ `panic` - Crashes de la aplicación
- ❌ `connection pool exhausted` - Problemas de DB
- ❌ `timeout` - Problemas de latencia

**Si hay errores**:
- **Documentar** cada error encontrado
- **Clasificar** por severidad (crítico, medio, bajo)
- **Priorizar** correcciones

### 4.6.2 Monitorear Uso de Recursos

**En Render Dashboard → Metrics**:

Verificar:
- **CPU Usage**: `< 60%` promedio
- **Memory Usage**: `< 70%` promedio
- **Response Time**: `< 500ms` P95

**Si los recursos están altos**:
- Considerar upgrade del plan de Render
- Optimizar queries de DB
- Implementar caching

### 4.6.3 Verificar Estado de Base de Datos

**En Supabase Dashboard**:

1. Ir a Database → Logs
2. Filtrar por tipo: `errors`

**Verificar que NO hay**:
- ❌ `deadlock detected`
- ❌ `too many connections`
- ❌ `out of memory`
- ❌ `slow query` (> 1 segundo)

**Queries lentas** (si las hay):
```sql
-- En Supabase SQL Editor
SELECT 
  query,
  mean_exec_time,
  calls
FROM pg_stat_statements
WHERE mean_exec_time > 1000  -- Más de 1 segundo
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Si hay queries lentas**:
- Documentar queries problemáticas
- Verificar que los índices necesarios existen
- Considerar optimización de queries

---

## ✅ FASE 4.7: Validación Final y Reporte (10 min)

### 4.7.1 Checklist de Validación

**Autenticación**:
- [ ] Registro de usuario funciona
- [ ] Login funciona
- [ ] Tokens JWT válidos
- [ ] Autorización funciona (endpoints protegidos)
- [ ] Manejo de errores correcto (401, 403)

**Datos Migrados**:
- [ ] Usuarios legacy pueden hacer login
- [ ] Datos de gamificación accesibles
- [ ] No hay pérdida de datos
- [ ] Achievements/quests visibles

**Endpoints Críticos**:
- [ ] Todos los endpoints críticos responden 200 OK
- [ ] CRUD completo funciona
- [ ] Validaciones funcionan correctamente
- [ ] No hay errores 500 inesperados

**Frontend**:
- [ ] Login/registro desde UI funciona
- [ ] Dashboard carga correctamente
- [ ] Datos legacy visibles en UI
- [ ] CRUD desde UI funciona
- [ ] No hay errores en consola

**Performance y Estabilidad**:
- [ ] Tiempos de respuesta aceptables (< 500ms)
- [ ] No hay memory leaks
- [ ] No hay errores en logs
- [ ] Uso de recursos aceptable

### 4.7.2 Crear Reporte de Pruebas

**Crear archivo**: `PHASE4-TESTING-REPORT.md`

```markdown
# Reporte de Pruebas - Fase 4

**Fecha**: 2026-02-13
**Ejecutor**: [Tu nombre]
**Duración**: [tiempo total]

## Resumen Ejecutivo
- ✅ Todas las pruebas pasadas
- ⚠️ [N] warnings detectados
- ❌ [N] errores críticos detectados

## Resultados por Categoría

### Autenticación
- Estado: ✅ PASS
- Pruebas ejecutadas: 7
- Pruebas pasadas: 7
- Notas: Ninguna

### Datos Migrados
- Estado: ✅ PASS
- Usuarios probados: 3
- Discrepancias: 0
- Notas: Datos legacy completamente accesibles

### Performance
- Tiempo promedio de respuesta: 234ms
- P95: 450ms
- P99: 680ms
- Estado: ✅ PASS

### Errores Encontrados
1. [ERROR-001]: [Descripción]
   - Severidad: Crítico/Medio/Bajo
   - Acción requerida: [...]

## Recomendaciones
1. [Recomendación 1]
2. [Recomendación 2]

## Conclusión
✅ El monolith está listo para recibir tráfico de producción

o

⚠️ Se requieren correcciones antes de proceder con Fase 5
```

### 4.7.3 Decisión Go/No-Go

**Criterios para proceder a Fase 5** (Migración de Tráfico):

**TODOS deben ser ✅**:
- [ ] No hay errores críticos (status 500, crashes, data loss)
- [ ] Autenticación 100% funcional
- [ ] Datos migrados 100% accesibles
- [ ] Performance aceptable (< 500ms P95)
- [ ] Frontend funciona sin errores
- [ ] Logs limpios (sin errores recurrentes)

**Si alguno es ❌**:
→ **NO proceder a Fase 5**
→ **Corregir problemas identificados**
→ **Re-ejecutar pruebas de Fase 4**

---

## 🚨 Plan de Acción si hay Problemas

### Errores Críticos Detectados

**Acción inmediata**:
1. **Documentar** el error con máximo detalle
2. **NO proceder** a Fase 5
3. **Rollback** si ya hay tráfico en el monolith
4. **Investigar** root cause en desarrollo
5. **Corregir** y re-deployar
6. **Re-ejecutar** Fase 4 completa

### Errores Menores/Warnings

**Acción**:
1. **Documentar** para corrección posterior
2. **Evaluar** si bloquean Fase 5
3. **Si no bloquean**: Proceder con Fase 5 pero monitorear de cerca
4. **Si bloquean**: Corregir primero

### Performance Degradada

**Si latencia > 1 segundo promedio**:
1. **Investigar** queries lentas en Supabase
2. **Optimizar** con índices
3. **Considerar** caching (Redis)
4. **Upgrade** recursos de Render si es necesario

---

## 📊 Métricas de Éxito de Fase 4

**La Fase 4 es exitosa si**:
- ✅ 100% de pruebas de autenticación pasan
- ✅ 100% de usuarios migrados pueden acceder sus datos
- ✅ 0 errores críticos (status 500, crashes)
- ✅ Latencia P95 < 500ms
- ✅ Frontend funciona sin errores
- ✅ Logs limpios sin errores recurrentes
- ✅ Equipo confiado para proceder a Fase 5

---

## 📚 Referencias

- **Fase 3 - Migración**: `PHASE3-PRODUCTION-MIGRATION-GUIDE.md`
- **Estado de Base de Datos**: Supabase Dashboard
- **Logs de Aplicación**: Render Dashboard
- **Next Step**: `PHASE5-TRAFFIC-MIGRATION-GUIDE.md` (a crear)

---

## 📝 Anexo: Scripts de Automatización

### Script de Testing Automatizado (Bash)

Guardar como `scripts/phase4-test.sh`:

```bash
#!/bin/bash

# Configuración
MONOLITH_URL="https://financial-resume-monolith.onrender.com"
TEST_EMAIL="test-phase4@example.com"
TEST_PASSWORD="TestPass123!"

echo "========================================="
echo "PHASE 4: Automated Testing"
echo "========================================="

# 1. Health Check
echo ""
echo "1. Testing health endpoint..."
HEALTH=$(curl -s -w "\n%{http_code}" $MONOLITH_URL/health)
HTTP_CODE=$(echo "$HEALTH" | tail -n 1)

if [ "$HTTP_CODE" == "200" ]; then
  echo "✅ Health check PASS"
else
  echo "❌ Health check FAIL (HTTP $HTTP_CODE)"
  exit 1
fi

# 2. Register User
echo ""
echo "2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $MONOLITH_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"name\":\"Test User\"}")

HTTP_CODE=$(echo "$REGISTER_RESPONSE" | tail -n 1)
if [ "$HTTP_CODE" == "201" ]; then
  echo "✅ Registration PASS"
  TOKEN=$(echo "$REGISTER_RESPONSE" | head -n -1 | jq -r '.token')
else
  echo "❌ Registration FAIL (HTTP $HTTP_CODE)"
  exit 1
fi

# 3. Login
echo ""
echo "3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $MONOLITH_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n 1)
if [ "$HTTP_CODE" == "200" ]; then
  echo "✅ Login PASS"
  TOKEN=$(echo "$LOGIN_RESPONSE" | head -n -1 | jq -r '.token')
else
  echo "❌ Login FAIL (HTTP $HTTP_CODE)"
  exit 1
fi

# 4. Get User Profile
echo ""
echo "4. Testing authenticated endpoint (user profile)..."
PROFILE_RESPONSE=$(curl -s -w "\n%{http_code}" $MONOLITH_URL/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$PROFILE_RESPONSE" | tail -n 1)
if [ "$HTTP_CODE" == "200" ]; then
  echo "✅ User profile PASS"
else
  echo "❌ User profile FAIL (HTTP $HTTP_CODE)"
  exit 1
fi

# 5. Get Gamification Profile
echo ""
echo "5. Testing gamification endpoint..."
GAMIFICATION_RESPONSE=$(curl -s -w "\n%{http_code}" $MONOLITH_URL/api/v1/gamification/profile \
  -H "Authorization: Bearer $TOKEN")

HTTP_CODE=$(echo "$GAMIFICATION_RESPONSE" | tail -n 1)
if [ "$HTTP_CODE" == "200" ]; then
  echo "✅ Gamification profile PASS"
else
  echo "❌ Gamification profile FAIL (HTTP $HTTP_CODE)"
  exit 1
fi

echo ""
echo "========================================="
echo "✅ ALL TESTS PASSED"
echo "========================================="
```

**Ejecutar**:
```bash
chmod +x scripts/phase4-test.sh
./scripts/phase4-test.sh
```

---

**Última actualización**: 2026-02-13  
**Autor**: AI Assistant  
**Versión**: 1.0
