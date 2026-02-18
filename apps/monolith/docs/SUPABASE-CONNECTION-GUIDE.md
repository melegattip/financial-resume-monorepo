# Guía: Obtener Connection String de Supabase

## Información que ya tenemos:
- **Project REF**: njgddjhjqzhhruklxrzg
- **Project URL**: https://njgddjhjqzhhruklxrzg.supabase.co

## Pasos para obtener la Connection String:

### 1. Accede a Supabase Dashboard
- Ve a: https://app.supabase.com/project/njgddjhjqzhhruklxrzg
- O desde https://app.supabase.com y selecciona tu proyecto

### 2. Obtén la Connection String de PostgreSQL

#### Opción A: Desde Settings → Database
1. En el menú lateral izquierdo, haz clic en **Settings** (⚙️)
2. Haz clic en **Database**
3. Busca la sección **"Connection String"**
4. Cambia el dropdown a **"URI"** (no "Transaction" ni "Session")
5. Verás algo como:
   ```
   postgresql://postgres.njgddjhjqzhhruklxrzg:[YOUR-PASSWORD]@aws-0-[region].pooler.supabase.com:6543/postgres
   ```
6. Copia esta cadena completa

#### Opción B: Construir manualmente
Si conoces la contraseña, la connection string sigue este formato:
```
postgresql://postgres.njgddjhjqzhhruklxrzg:[PASSWORD]@aws-0-[REGION].pooler.supabase.com:6543/postgres?sslmode=require
```

Donde:
- `njgddjhjqzhhruklxrzg` = tu Project REF (ya lo tenemos)
- `[PASSWORD]` = la contraseña de la base de datos (la que guardaste al crear el proyecto)
- `[REGION]` = la región de AWS donde se creó (ejemplo: `sa-east-1`, `us-east-1`, etc.)

### 3. Verificar la región

En **Settings → General** verás:
- **Region**: La región donde está deployado tu proyecto
- Ejemplos comunes:
  - South America (São Paulo) = `sa-east-1`
  - US East (N. Virginia) = `us-east-1`
  - US West (Oregon) = `us-west-2`

### 4. Obtener/Resetear la contraseña

Si no guardaste la contraseña cuando creaste el proyecto:
1. Ve a **Settings → Database**
2. Busca **"Database Password"**
3. Haz clic en **"Reset Database Password"**
4. Se generará una nueva contraseña
5. **¡Guárdala inmediatamente en un lugar seguro!**
6. Esta acción NO afecta datos, solo resetea la contraseña de conexión

### 5. Información que necesito

Una vez que tengas todo, completa:

```bash
# Connection String completa (REEMPLAZA la contraseña con [PASSWORD] al compartirla)
DATABASE_URL=postgresql://postgres.njgddjhjqzhhruklxrzg:[PASSWORD]@aws-0-XXXX.pooler.supabase.com:6543/postgres?sslmode=require

# Región (ejemplo: sa-east-1)
REGION=XXXX
```

## Próximo Paso

Una vez que tengas esta información, podremos:
1. Crear el archivo `.env.production` con todas las credenciales
2. Verificar la conexión a la nueva base de datos
3. Preparar la migración de datos desde las bases legacy
