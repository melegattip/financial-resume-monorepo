# Script completo para setup y migración
# Uso: .\scripts\setup-migration.ps1 -LegacyUsersDB "url" -LegacyGamificationDB "url"

param(
    [Parameter(Mandatory=$true)]
    [string]$LegacyUsersDB,
    
    [Parameter(Mandatory=$true)]
    [string]$LegacyGamificationDB,
    
    [string]$NewDBName = "financial_resume_v2"
)

Write-Host "=== Setup de Migración - Fase 3 ===" -ForegroundColor Cyan

# Configuración PostgreSQL local
$LOCAL_DB_HOST = "localhost"
$LOCAL_DB_PORT = "5432"
$LOCAL_DB_USER = "postgres"
$LOCAL_DB_PASSWORD = "postgres"
$NEW_DB_URL = "postgresql://${LOCAL_DB_USER}:${LOCAL_DB_PASSWORD}@${LOCAL_DB_HOST}:${LOCAL_DB_PORT}/${NewDBName}"

# Paso 1: Levantar PostgreSQL en Docker
Write-Host "`n[1/5] Verificando PostgreSQL en Docker..." -ForegroundColor Yellow
$postgresRunning = docker ps --filter "name=fr_postgres_main" --format "{{.Names}}"

if (-not $postgresRunning) {
    Write-Host "  Levantando contenedor PostgreSQL..." -ForegroundColor Yellow
    $dockerComposePath = "infrastructure\docker\docker-compose.yml"
    
    if (Test-Path $dockerComposePath) {
        Push-Location "infrastructure\docker"
        docker-compose up -d postgres-main
        Pop-Location
        
        # Esperar a que PostgreSQL esté listo
        Write-Host "  Esperando a que PostgreSQL esté listo..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
        
        $maxRetries = 30
        $retry = 0
        while ($retry -lt $maxRetries) {
            $ready = docker exec fr_postgres_main pg_isready -U postgres 2>&1
            if ($ready -match "accepting connections") {
                Write-Host "  ✅ PostgreSQL está listo" -ForegroundColor Green
                break
            }
            $retry++
            Start-Sleep -Seconds 1
        }
    } else {
        Write-Host "  ❌ No se encontró docker-compose.yml" -ForegroundColor Red
        Write-Host "  Por favor, levanta PostgreSQL manualmente:" -ForegroundColor Yellow
        Write-Host "  docker run -d --name fr_postgres_main -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:15" -ForegroundColor Gray
        exit 1
    }
} else {
    Write-Host "  ✅ PostgreSQL ya está corriendo" -ForegroundColor Green
}

# Paso 2: Crear nueva base de datos
Write-Host "`n[2/5] Creando nueva base de datos: $NewDBName" -ForegroundColor Yellow

$createDbSQL = "CREATE DATABASE $NewDBName;"
$checkDbSQL = "SELECT 1 FROM pg_database WHERE datname = '$NewDBName';"

# Verificar si ya existe
$exists = docker exec fr_postgres_main psql -U $LOCAL_DB_USER -t -c $checkDbSQL 2>&1

if ($exists -match "1") {
    Write-Host "  ⚠️  La base de datos '$NewDBName' ya existe" -ForegroundColor Yellow
    $response = Read-Host "  ¿Deseas continuar usando la base existente? (S/N)"
    if ($response -ne "S" -and $response -ne "s") {
        exit 0
    }
} else {
    $result = docker exec fr_postgres_main psql -U $LOCAL_DB_USER -c $createDbSQL 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ Base de datos '$NewDBName' creada" -ForegroundColor Green
    } else {
        Write-Host "  ❌ Error al crear la base de datos:" -ForegroundColor Red
        Write-Host $result -ForegroundColor Red
        exit 1
    }
}

# Paso 3: Crear archivo .env
Write-Host "`n[3/5] Creando archivo .env con configuración..." -ForegroundColor Yellow

$envContent = @"
# ============================================
# NUEVA BASE DE DATOS (Objetivo - Monolith)
# ============================================
DATABASE_URL=$NEW_DB_URL

# ============================================
# BASES LEGACY (Solo Lectura - Productivas)
# ============================================
USERS_DB_URL=$LegacyUsersDB
GAMIFICATION_DB_URL=$LegacyGamificationDB

# ============================================
# Configuración Adicional
# ============================================
LOG_LEVEL=info
APP_ENV=development
PORT=8080
"@

$envPath = "apps\monolith\.env"
$envContent | Out-File -FilePath $envPath -Encoding UTF8
Write-Host "  ✅ Archivo .env creado en: $envPath" -ForegroundColor Green

# Paso 4: Compilar comando de migración
Write-Host "`n[4/5] Compilando comando de migración..." -ForegroundColor Yellow
Push-Location "apps\monolith"
go build -o bin\migrate.exe .\cmd\migrate
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✅ Comando compilado exitosamente" -ForegroundColor Green
} else {
    Write-Host "  ❌ Error al compilar" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location

# Paso 5: Ejecutar auditoría
Write-Host "`n[5/5] Ejecutando auditoría pre-migración..." -ForegroundColor Yellow
Push-Location "apps\monolith"
.\bin\migrate.exe audit
Pop-Location

Write-Host "`n=== Setup Completado ===" -ForegroundColor Green
Write-Host "`nPróximos pasos:" -ForegroundColor Cyan
Write-Host "1. Revisa el reporte de auditoría arriba" -ForegroundColor White
Write-Host "2. Ejecuta dry-run: cd apps\monolith; .\bin\migrate.exe migrate --dry-run" -ForegroundColor White
Write-Host "3. Si todo se ve bien, ejecuta: .\bin\migrate.exe migrate" -ForegroundColor White
