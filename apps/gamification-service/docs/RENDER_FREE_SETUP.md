# 🚀 Configuración Simplificada para Render Free

## 📋 Resumen

✅ **Una sola base de datos PostgreSQL** para ambos servicios  
✅ **Compatible con Render Free** (1 BD gratuita)  
✅ **Sin sincronización** necesaria  
✅ **Configuración simple** para desarrollo y producción  

## 🏗️ Arquitectura Simplificada

```
🌐 Internet
    ↓
📱 Frontend (React)
    ↓
┌─────────────────────────────────────────┐
│ 🚀 financial-resume-engine:8080        │ ──┐
│   (API principal, proxy gateway)       │   │
└─────────────────────────────────────────┘   │
    ↓ (proxy HTTP)                            ├── 🗄️ financial_resume DB
┌─────────────────────────────────────────┐   │     (UNA SOLA BASE DE DATOS)
│ 🎮 financial-gamification-service:8081 │ ──┘
│   (Microservicio de gamificación)      │
└─────────────────────────────────────────┘
```

## ☁️ Configuración en Render.com

### **1. Crear una sola PostgreSQL Database:**

En Render Dashboard:
1. **New PostgreSQL** → `financial-resume-db`
   - Name: `financial-resume-db`
   - Database: `financial_resume`
   - Plan: **Free** ✅

### **2. Variables de entorno para ambos servicios:**

#### **financial-resume-engine:**
```env
DATABASE_URL=postgresql://user:pass@host:port/financial_resume
GAMIFICATION_SERVICE_URL=https://financial-gamification-service.onrender.com
DB_HOST=your-db-host
DB_USER=your-db-user  
DB_PASSWORD=your-db-password
DB_NAME=financial_resume
```

#### **financial-gamification-service:**
```env
# ✅ MISMA BD QUE FINANCIAL-RESUME-ENGINE
DATABASE_URL=postgresql://user:pass@host:port/financial_resume
DB_HOST=your-db-host
DB_USER=your-db-user
DB_PASSWORD=your-db-password  
DB_NAME=financial_resume
JWT_SECRET=financial_resume_secret_key_2024
```

## 🔧 Desarrollo Local

### **Configuración actual (ya implementada):**

**financial-resume-engine:**
- BD: `localhost:5432/financial_resume`

**financial-gamification-service:**  
- BD: `host.docker.internal:5432/financial_resume` (misma BD)

### **Comandos para desarrollo:**

```bash
# Levantar BD principal
cd financial-resume-engine
docker-compose up -d postgres

# Levantar servicios
docker-compose up -d  # financial-resume-engine
cd ../financial-gamification-service  
docker-compose up -d  # gamification-service (sin BD propia)
```

## 🚀 Deployment en Render

### **Paso a paso:**

1. **Crear BD en Render:**
   - New PostgreSQL → `financial-resume-db`
   - Copiar DATABASE_URL

2. **Deploy financial-resume-engine:**
   - New Web Service
   - Variables: `DATABASE_URL`, `GAMIFICATION_SERVICE_URL`
   - Build: `go build cmd/api/main.go`
   - Start: `./main`

3. **Deploy financial-gamification-service:**
   - New Web Service  
   - Variables: `DATABASE_URL` (misma que step 1)
   - Build: `go build cmd/api/main.go`
   - Start: `./main`

4. **Inicializar BD:**
   ```bash
   # Solo UNA VEZ al inicio
   psql $DATABASE_URL -f financial-resume-engine/scripts/init.sql
   psql $DATABASE_URL -f financial-resume-engine/scripts/test_data.sql
   ```

## 📊 Verificación

### **Health checks:**
```bash
# API Principal
curl https://financial-resume-engine.onrender.com/health

# Gamificación  
curl https://financial-gamification-service.onrender.com/health
```

### **Verificar datos de gamificación:**
```bash
# Nivel 1
curl "https://financial-resume-engine.onrender.com/api/v1/gamification/profile" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-User-ID: 1"

# Nivel 10  
curl "https://financial-resume-engine.onrender.com/api/v1/gamification/profile" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-User-ID: 4"
```

## 💡 Usuarios de Prueba

| Email | Password | Nivel | XP | Descripción |
|-------|----------|-------|----|-----------| 
| `nivel1@test.com` | `password123` | **1** | 0 | Financial Newbie |
| `nivel3@test.com` | `password123` | **3** | 200 | Smart Saver |
| `nivel5@test.com` | `password123` | **5** | 700 | Financial Planner |
| `nivel10@test.com` | `password123` | **10** | 5500 | Financial Magnate |
| `pablo@niloft.com` | `password123` | **7** | 2000 | Financial Expert |

## 🔒 Seguridad

### **Variables de entorno en Render:**
```bash
DATABASE_URL=postgresql://...        # Automática
JWT_SECRET=your_secret_key_here      # Manual
GAMIFICATION_SERVICE_URL=https://... # Manual
```

### **SSL en producción:**
```bash
# Render automáticamente configura SSL
# No requiere configuración adicional
```

## 🚨 Troubleshooting

### **Problema: Usuario aparece con nivel incorrecto**
```bash
# Verificar BD directamente:
psql $DATABASE_URL -c "SELECT user_id, current_level, total_xp FROM user_gamification ORDER BY user_id;"

# Si falta algún usuario, re-ejecutar:
psql $DATABASE_URL -f financial-resume-engine/scripts/test_data.sql
```

### **Problema: Error de conexión a BD**
```bash
# Verificar DATABASE_URL en ambos servicios
# Debe ser exactamente igual en ambos

# En Render Dashboard:
# Settings → Environment → DATABASE_URL
```

### **Problema: Gamification service no responde**
```bash
# Verificar que ambos servicios usen misma BD
# Logs en Render Dashboard → Service → Logs
```

## 📋 Checklist de Deployment

### **Antes del deploy:**
- [ ] Una sola BD creada en Render
- [ ] DATABASE_URL copiada y configurada en ambos servicios
- [ ] Variables de entorno configuradas
- [ ] Scripts SQL aplicados

### **Después del deploy:**
- [ ] Health checks responden OK
- [ ] Login funciona con usuarios de prueba  
- [ ] Niveles de gamificación aparecen correctos
- [ ] Endpoints de proxy funcionan

## 💰 Costo en Render Free

✅ **Free Tier incluye:**
- 1 PostgreSQL Database (500MB) - GRATIS
- 2 Web Services - GRATIS  
- SSL automático - GRATIS
- Custom domains - GRATIS

**Total: $0/mes** 🎉

## 🔄 Migración Futura

Si en el futuro quieres volver a **2 BDs separadas**:

1. Crear segunda BD en plan pago
2. Ejecutar script de sincronización
3. Cambiar variables de entorno
4. ¡Sin cambios de código!

**La aplicación está preparada para ambas arquitecturas.** 