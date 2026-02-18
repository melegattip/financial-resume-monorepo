# Guía de Migración en Producción - Fase 3
**Base de Datos: Consolidación a Monolith**

**Fecha de Creación**: 2026-02-13  
**Ambiente**: Producción (Supabase)  
**Tiempo Estimado**: 2-3 horas

---

## 🎯 Objetivo

Migrar las bases de datos legacy de producción (Supabase) a una nueva base de datos consolidada para el monolith, **sin downtime** y con capacidad de rollback inmediato.

---

## 📋 Pre-Requisitos

### 1. Infraestructura Actual (Legacy)
- ✅ `users_db` (Supabase) - Base de usuarios actual
- ✅ `financial_resume` (Supabase) - Base de datos financieros/gamificación
- ✅ Servicios corriendo en Render/producción

### 2. Infraestructura Nueva
- ⚠️ Nueva base de datos en Supabase (a crear)
- ⚠️ Comando de migración compilado (`migrate.exe` o `migrate` para Linux)
- ⚠️ Acceso de red desde donde ejecutarás la migración a Supabase

### 3. Accesos Necesarios
- ✅ Credenciales de admin de Supabase
- ✅ Variables de entorno de producción
- ✅ Acceso SSH/terminal donde se ejecutará la migración

---

## 🚨 IMPORTANTE: Estrategia de Migración

### **Opción Recomendada: Blue-Green Deployment** 

Esta estrategia permite:
- ✅ **Cero downtime** para los usuarios
- ✅ **Rollback inmediato** si algo falla
- ✅ **Validación completa** antes de cambiar tráfico
- ✅ **Coexistencia** de ambas versiones durante la transición

```
┌─────────────────────────────────────────────────────────────┐
│                    FASE 1: PREPARACIÓN                      │
│  - Crear nueva DB en Supabase (financial_resume_v2)        │
│  - Configurar pero NO activar                               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                  FASE 2: MIGRACIÓN (OFFLINE)                │
│  - Ejecutar migración desde tu máquina local                │
│  - Copiar datos de legacy → nueva DB                        │
│  - Validar integridad                                       │
│  - Servicios legacy siguen funcionando normalmente          │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                FASE 3: DEPLOYAR MONOLITH (NUEVO)            │
│  - Deployar monolith apuntando a nueva DB                   │
│  - Ejecutar en paralelo con servicios legacy                │
│  - NO recibe tráfico todavía                                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              FASE 4: PRUEBAS EN PRODUCCIÓN                  │
│  - Acceder directamente al monolith (URL interna)           │
│  - Probar funcionalidad crítica                             │
│  - Validar que todo funciona                                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│             FASE 5: CAMBIO DE TRÁFICO GRADUAL               │
│  - 10% usuarios → monolith                                  │
│  - Monitorear 1 hora                                        │
│  - 50% usuarios → monolith                                  │
│  - Monitorear 1 hora                                        │
│  - 100% usuarios → monolith                                 │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                  FASE 6: DESACTIVAR LEGACY                  │
│  - Mantener legacy activo por 1-2 semanas                   │
│  - Monitorear estabilidad                                   │
│  - Desactivar servicios legacy gradualmente                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 📝 FASE 1: Preparación (1 hora)

### 1.1 Crear Nueva Base de Datos en Supabase

**Paso 1**: Ve al Dashboard de Supabase → https://app.supabase.com

**Paso 2**: Crear nuevo proyecto (o nueva base de datos según tu plan)
- **Nombre**: `financial-resume-monolith` o `financial-resume-v2`
- **Región**: **LA MISMA** que tus bases legacy (para menor latencia)
- **Plan**: El mismo o superior al actual
- **Contraseña**: Guardar en lugar seguro (1Password, etc.)

**Paso 3**: Esperar a que la base esté lista (5-10 minutos)

**Paso 4**: Obtener la Connection String
- Settings → Database → Connection String
- Copiar la URI format: `postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres`

**Paso 5**: Guardar las credenciales

### 1.2 Crear Backups de las Bases Legacy

**CRÍTICO**: Antes de continuar, crear backups completos.

**Opción A: Desde Supabase Dashboard** (Recomendado)
- Dashboard → Database → Backups → "Backup Now"
- Hacer backup de cada proyecto

**Opción B: Usando pg_dump** (Backup local adicional)
```bash
# Backup de users_db
pg_dump "postgresql://postgres.akngrdpnwboujagnziqb:[PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres" > backup_users_db_$(date +%Y%m%d_%H%M%S).sql

# Backup de financial_resume (gamification)
pg_dump "postgresql://postgres.gtzkqlbkqgnaittehfey:[PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres" > backup_gamification_db_$(date +%Y%m%d_%H%M%S).sql
```

**Verificar backups**:
```bash
# Verificar que los archivos existen y tienen tamaño > 0
ls -lh backup_*.sql
```

### 1.3 Preparar Variables de Entorno

Crear archivo `.env.production` en `apps/monolith/`:

```bash
# ============================================
# NUEVA BASE DE DATOS (Objetivo - Producción)
# ============================================
DATABASE_URL=postgresql://postgres.[NEW-PROJECT-REF]:[NEW-PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?sslmode=require

# ============================================
# BASES LEGACY (Solo Lectura)
# ============================================
# IMPORTANTE: URL-encode caracteres especiales en la contraseña
# $ → %24, @ → %40, etc.
USERS_DB_URL=postgresql://postgres.akngrdpnwboujagnziqb:[PASSWORD-URL-ENCODED]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?sslmode=require
GAMIFICATION_DB_URL=postgresql://postgres.gtzkqlbkqgnaittehfey:[PASSWORD-URL-ENCODED]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?sslmode=require

# ============================================
# Configuración
# ============================================
LOG_LEVEL=info
APP_ENV=production
```

### 1.4 Compilar Comando de Migración

```bash
cd apps/monolith
go build -o bin/migrate ./cmd/migrate
```

**Verificar compilación**:
```bash
./bin/migrate --help
```

### 1.5 Ventana de Mantenimiento (OPCIONAL)

**Decisión**: ¿Quieres programar una ventana de mantenimiento?

**Opción A: Sin ventana** (Recomendado)
- Migración no afecta servicios legacy
- Usuarios no se enteran
- Migración offline

**Opción B: Con ventana** (Más conservador)
- Anunciar a usuarios 24-48 hrs antes
- Ventana de 2-3 horas
- Poner banner de "Mantenimiento programado"

---

## 🚀 FASE 2: Ejecución de la Migración (30-45 min)

### 2.1 Verificar Conexiones (Auditoría)

```bash
cd apps/monolith

# Cargar variables de entorno
export $(cat .env.production | xargs)

# Ejecutar auditoría
./bin/migrate audit
```

**Verificar output:**
- ✅ Conexión exitosa a nueva DB (debe estar vacía)
- ✅ Conexión exitosa a users_db (debe mostrar N usuarios)
- ✅ Conexión exitosa a gamification_db (debe mostrar datos)

**Si hay errores**:
- Verificar URLs de conexión
- Verificar que las contraseñas estén URL-encoded
- Verificar conexión de red a Supabase

### 2.2 Dry-Run (Simulación)

**IMPORTANTE**: SIEMPRE ejecutar dry-run primero

```bash
./bin/migrate migrate --dry-run
```

**Revisar cuidadosamente**:
- ✅ Conteo de registros a copiar es correcto
- ✅ No hay errores en las queries SQL
- ✅ Las transformaciones se ven correctas

**Guardar output**:
```bash
./bin/migrate migrate --dry-run > migration_dryrun_$(date +%Y%m%d_%H%M%S).log
```

### 2.3 Ejecutar Migración Real

**Hora de ejecución recomendada**: Durante horas de bajo tráfico (madrugada, fin de semana)

```bash
# Ejecutar migración completa
./bin/migrate migrate

# O por fases (más control):
./bin/migrate migrate --phase=1  # Auditoría
./bin/migrate migrate --phase=2  # Schema
./bin/migrate migrate --phase=3  # Datos
./bin/migrate migrate --phase=4  # Validación
```

**Guardar output completo**:
```bash
./bin/migrate migrate 2>&1 | tee migration_production_$(date +%Y%m%d_%H%M%S).log
```

**Tiempo estimado**: 15-30 minutos (depende del volumen de datos)

### 2.4 Validación Post-Migración

```bash
./bin/migrate validate
```

**Verificar que TODOS los checks pasen**:
- ✅ `overall: "PASS"`
- ✅ Conteo de usuarios correcto
- ✅ No hay referencias huérfanas
- ✅ Tipos de datos correctos
- ✅ Índices creados

**Si algún check falla**:
- Revisar el reporte de validación
- NO proceder a la siguiente fase
- Investigar y corregir el problema
- Posiblemente re-ejecutar la migración

### 2.5 Verificación Manual (Opcional pero Recomendado)

Conectarse a la nueva base y verificar datos:

```bash
psql "postgresql://postgres.[NEW-PROJECT-REF]:[PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres"
```

```sql
-- Verificar conteo de usuarios
SELECT COUNT(*) FROM users;

-- Verificar que hay datos de gamificación
SELECT COUNT(*) FROM user_gamification;

-- Verificar que no hay duplicados
SELECT user_id, COUNT(*) 
FROM user_gamification 
GROUP BY user_id 
HAVING COUNT(*) > 1;

-- Verificar tipos de columnas
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' AND column_name = 'id';
```

---

## 🎨 FASE 3: Deployar el Monolith (30 min)

### 3.1 Preparar Deployment del Monolith

**En Render** (o tu plataforma de deployment):

**Paso 1**: Crear nuevo Web Service
- **Nombre**: `financial-resume-monolith`
- **Repository**: Tu repo actual
- **Branch**: `main` o `production`
- **Root Directory**: `apps/monolith`
- **Build Command**: `go build -o bin/server ./cmd/server`
- **Start Command**: `./bin/server`

**Paso 2**: Configurar Variables de Entorno
```
DATABASE_URL=postgresql://postgres.[NEW-PROJECT-REF]:[PASSWORD]@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?sslmode=require
JWT_SECRET=[TU_JWT_SECRET_ACTUAL]
PORT=8080
APP_ENV=production
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=https://tu-frontend.vercel.app
```

**Paso 3**: **NO ACTIVAR TODAVÍA** el tráfico público

### 3.2 Deployar el Monolith

```bash
# Si usas Render, el deployment se hace automáticamente al hacer push
git push origin main

# Esperar a que el deployment complete (5-10 min)
```

### 3.3 Verificar Health del Monolith

```bash
# Obtener la URL interna del monolith desde Render
curl https://financial-resume-monolith.onrender.com/health

# Deberías ver:
# {"status": "ok", "database": "connected"}
```

---

## 🧪 FASE 4: Pruebas en Producción (30 min)

### 4.1 Pruebas de Endpoints Críticos

**Usar Postman, curl o alguna herramienta de testing:**

```bash
# 1. Registro de usuario (usa un email de prueba)
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-migration@example.com",
    "password": "Test123!",
    "name": "Migration Test User"
  }'

# 2. Login
curl -X POST https://financial-resume-monolith.onrender.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test-migration@example.com",
    "password": "Test123!"
  }'

# Guardar el token JWT que retorna

# 3. Obtener datos de usuario
curl https://financial-resume-monolith.onrender.com/api/v1/users/me \
  -H "Authorization: Bearer [TOKEN]"

# 4. Obtener datos de gamificación
curl https://financial-resume-monolith.onrender.com/api/v1/gamification/profile \
  -H "Authorization: Bearer [TOKEN]"
```

**Verificar**:
- ✅ Todos los endpoints responden 200 OK
- ✅ Los datos retornados son correctos
- ✅ No hay errores en los logs

### 4.2 Pruebas con Usuario Real

**IMPORTANTE**: Usar un usuario de prueba o tu propio usuario

1. Hacer login en el frontend apuntando temporalmente al monolith
2. Verificar que puedes:
   - ✅ Ver tu perfil
   - ✅ Ver tus datos de gamificación
   - ✅ Crear/editar datos (si aplica)
   - ✅ Ver achievements

### 4.3 Revisar Logs

**En Render**:
- Ir a Logs del servicio monolith
- Verificar que no hay errores críticos
- Verificar que las queries a la DB funcionan

---

## 🔄 FASE 5: Migración Gradual de Tráfico (2-3 horas)

### Estrategia: Canary Deployment

**Herramientas de routing progresivo**:
- **Cloudflare**: Workers con % routing
- **Vercel**: Edge Functions con A/B testing
- **Nginx**: Weighted routing
- **API Gateway**: Route weighting

### 5.1 Configurar Routing Progresivo

**Ejemplo con Cloudflare Workers** (pseudocódigo):

```javascript
// cloudflare-worker.js
async function handleRequest(request) {
  const random = Math.random();
  
  // 10% de tráfico al monolith nuevo
  if (random < 0.10) {
    return fetch('https://financial-resume-monolith.onrender.com' + request.url.pathname, request);
  } else {
    // 90% al API Gateway viejo
    return fetch('https://api-gateway-legacy.onrender.com' + request.url.pathname, request);
  }
}
```

### 5.2 Incremento Gradual

**Semana 1**:
- **Día 1**: 10% tráfico → monolith (monitorear 24h)
- **Día 2**: Si todo OK → 25% tráfico
- **Día 3**: Si todo OK → 50% tráfico
- **Día 4**: Si todo OK → 75% tráfico
- **Día 5**: Si todo OK → 100% tráfico

**En cada incremento, monitorear**:
- ✅ Tasa de errores (< 0.1%)
- ✅ Latencia promedio (< 500ms)
- ✅ Tasa de éxito de autenticación (> 99%)
- ✅ No hay quejas de usuarios

### 5.3 Monitoreo Continuo

**Herramientas recomendadas**:
- **Sentry**: Para tracking de errores
- **DataDog** o **New Relic**: Para métricas de performance
- **Logs de Render**: Para debugging

**Métricas clave**:
```
- Requests per minute (RPM)
- Error rate (%)
- P50, P95, P99 latency
- Database connection pool utilization
- Memory usage
- CPU usage
```

---

## 🚨 FASE 6: Plan de Rollback

### ⚠️ IMPORTANTE: Tener plan de rollback listo en todo momento

### 6.1 Rollback Inmediato (Si algo falla)

**Si detectas problemas graves**:

1. **Revertir Tráfico** (< 5 minutos)
   ```javascript
   // En Cloudflare Worker (o tu router)
   // Cambiar a 100% tráfico al legacy
   return fetch('https://api-gateway-legacy.onrender.com' + request.url.pathname, request);
   ```

2. **Verificar que legacy funciona**
   ```bash
   curl https://api-gateway-legacy.onrender.com/health
   ```

3. **Investigar logs del monolith**
   - Identificar el error
   - Corregir en desarrollo
   - Re-deployar cuando esté corregido

### 6.2 Rollback Parcial (Si algunos endpoints fallan)

**Routing selectivo**:
```javascript
// Solo rutas específicas al legacy
if (request.url.pathname.startsWith('/api/v1/problematic-endpoint')) {
  return fetch('https://api-gateway-legacy.onrender.com' + request.url.pathname);
} else {
  return fetch('https://financial-resume-monolith.onrender.com' + request.url.pathname);
}
```

---

## 📊 FASE 7: Post-Migración y Limpieza (1-2 semanas después)

### 7.1 Periodo de Observación

**Mantener ambas infraestructuras por 1-2 semanas**:
- Monolith recibe 100% del tráfico
- Legacy sigue activo pero en standby
- Monitorear estabilidad del monolith

### 7.2 Desactivar Servicios Legacy

**Cuando estés 100% seguro** (después de 1-2 semanas estables):

1. **Desactivar tráfico a API Gateway legacy**
2. **Pausar servicios en Render** (no eliminar todavía):
   - api-gateway (legacy)
   - users-service (legacy)
   - gamification-service (legacy)

3. **Mantener bases de datos legacy** por 1 mes más
   - Por si necesitas consultar datos históricos
   - Como backup de emergencia

### 7.3 Limpieza Final (1 mes después)

**Si todo ha sido estable por 1 mes**:

1. **Eliminar servicios legacy de Render**
2. **Pausar proyectos de Supabase legacy** (no eliminar, solo pausar)
3. **Crear backup final de bases legacy**
4. **Archivar documentación de arquitectura legacy**

---

## ✅ Checklist de Migración Completa

### Pre-Migración
- [ ] Nueva base de datos creada en Supabase
- [ ] Backups completos de bases legacy creados
- [ ] Variables de entorno preparadas (`.env.production`)
- [ ] Comando de migración compilado
- [ ] Dry-run ejecutado y revisado
- [ ] Ventana de mantenimiento comunicada (si aplica)

### Migración
- [ ] Auditoría ejecutada sin errores
- [ ] Migración real ejecutada
- [ ] Validación completa PASS
- [ ] Verificación manual de datos correcta
- [ ] Logs de migración guardados

### Deployment
- [ ] Monolith deployado en Render
- [ ] Health check del monolith OK
- [ ] Variables de entorno de producción configuradas
- [ ] Pruebas de endpoints críticos OK
- [ ] Pruebas con usuario real OK

### Migración de Tráfico
- [ ] Router configurado (Cloudflare/Vercel/etc.)
- [ ] 10% tráfico → monitoreo OK
- [ ] 25% tráfico → monitoreo OK
- [ ] 50% tráfico → monitoreo OK
- [ ] 75% tráfico → monitoreo OK
- [ ] 100% tráfico → monitoreo OK

### Post-Migración
- [ ] Monitoreo continuo por 1-2 semanas
- [ ] No hay incidentes reportados
- [ ] Servicios legacy pausados
- [ ] Bases legacy archivadas
- [ ] Documentación actualizada

---

## 🆘 Solución de Problemas

### Problema: "connection refused" durante migración

**Solución**:
- Verificar que las URLs de Supabase son correctas
- Verificar que las contraseñas están URL-encoded
- Verificar conexión de red a Supabase

### Problema: "orphaned user_ids" en validación

**Solución**:
- Hay user_ids en tablas secundarias que no existen en `users`
- Opciones:
  1. Crear usuarios placeholder para esos user_ids
  2. Eliminar esos registros huérfanos
  3. Investigar por qué existen esos user_ids

### Problema: Alta latencia en el monolith

**Solución**:
- Verificar índices en la base de datos
- Verificar que el connection pool está bien configurado
- Considerar aumentar recursos del servicio en Render
- Usar caching (Redis) para queries frecuentes

### Problema: Errores de autenticación después de migración

**Solución**:
- Verificar que `JWT_SECRET` es el MISMO que en legacy
- Verificar que el formato de los tokens es compatible
- Verificar que las sesiones de usuario se mantienen

---

## 📞 Contactos de Emergencia

**Durante la migración, tener a mano**:
- [ ] Acceso a dashboard de Supabase
- [ ] Acceso a dashboard de Render
- [ ] Logs de aplicación
- [ ] Herramienta de monitoreo (Sentry/DataDog)
- [ ] Plan de rollback impreso/accesible

---

## 📈 Métricas de Éxito

La migración es exitosa si:
- ✅ 100% de usuarios migrados sin pérdida de datos
- ✅ < 0.1% de error rate en producción
- ✅ Latencia promedio < 500ms
- ✅ No hay quejas de usuarios
- ✅ Monolith estable por 2 semanas consecutivas
- ✅ Servicios legacy desactivados sin problemas

---

## 📚 Referencias

- **Guía de Migración Local**: `05-phase3-migration-guide-NEW-DB.md`
- **Reporte de Verificación**: `PHASE3-VERIFICATION-REPORT.md`
- **Estado de Fase 3**: `PHASE3-STATUS.md`
- **Plan de Migración**: `03-migration-plan.md`

---

**Última actualización**: 2026-02-13  
**Autor**: AI Assistant  
**Versión**: 1.0
