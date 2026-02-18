# Script para crear la nueva base de datos en PostgreSQL local (Docker)
# Ejecuta: .\scripts\create-new-db.ps1

Write-Host "=== Creando Nueva Base de Datos para Monolith ===" -ForegroundColor Cyan

# Configuración
$DB_HOST = "localhost"
$DB_PORT = "5432"
$DB_USER = "postgres"
$DB_PASSWORD = "postgres"
$NEW_DB_NAME = "financial_resume_v2"

# Verificar si PostgreSQL está corriendo en Docker
Write-Host "`nVerificando PostgreSQL en Docker..." -ForegroundColor Yellow
$postgresContainer = docker ps --filter "name=postgres" --format "{{.Names}}" | Select-Object -First 1

if (-not $postgresContainer) {
    Write-Host "No se encontró contenedor de PostgreSQL corriendo." -ForegroundColor Red
    Write-Host "Opciones:" -ForegroundColor Yellow
    Write-Host "1. Levantar PostgreSQL con docker-compose:" -ForegroundColor White
    Write-Host "   cd infrastructure/docker" -ForegroundColor Gray
    Write-Host "   docker-compose up -d postgres-main" -ForegroundColor Gray
    Write-Host "`n2. O usar PostgreSQL local si está instalado" -ForegroundColor White
    exit 1
}

Write-Host "PostgreSQL encontrado: $postgresContainer" -ForegroundColor Green

# Crear la nueva base de datos
Write-Host "`nCreando base de datos: $NEW_DB_NAME" -ForegroundColor Yellow

# Usar docker exec para crear la base de datos
$createDbCommand = "CREATE DATABASE $NEW_DB_NAME;"

# Intentar con el contenedor encontrado primero
$result = docker exec -i $postgresContainer psql -U $DB_USER -c $createDbCommand 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Base de datos '$NEW_DB_NAME' creada exitosamente" -ForegroundColor Green
    Write-Host "`nConnection String para .env:" -ForegroundColor Cyan
    Write-Host "DATABASE_URL=postgresql://$DB_USER`:$DB_PASSWORD@$DB_HOST`:$DB_PORT/$NEW_DB_NAME" -ForegroundColor White
} else {
    # Si falla, intentar con postgres-main directamente
    Write-Host "Intentando con contenedor postgres-main..." -ForegroundColor Yellow
    $result = docker exec -i fr_postgres_main psql -U $DB_USER -c $createDbCommand 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Base de datos '$NEW_DB_NAME' creada exitosamente" -ForegroundColor Green
        Write-Host "`nConnection String para .env:" -ForegroundColor Cyan
        Write-Host "DATABASE_URL=postgresql://$DB_USER`:$DB_PASSWORD@$DB_HOST`:$DB_PORT/$NEW_DB_NAME" -ForegroundColor White
    } else {
        # Verificar si la base ya existe
        $checkDbCommand = "SELECT 1 FROM pg_database WHERE datname = '$NEW_DB_NAME';"
        $exists = docker exec -i fr_postgres_main psql -U $DB_USER -t -c $checkDbCommand 2>&1
        
        if ($exists -match "1") {
            Write-Host "⚠️  La base de datos '$NEW_DB_NAME' ya existe" -ForegroundColor Yellow
            Write-Host "`nConnection String para .env:" -ForegroundColor Cyan
            Write-Host "DATABASE_URL=postgresql://$DB_USER`:$DB_PASSWORD@$DB_HOST`:$DB_PORT/$NEW_DB_NAME" -ForegroundColor White
        } else {
            Write-Host "❌ Error al crear la base de datos:" -ForegroundColor Red
            Write-Host $result -ForegroundColor Red
            exit 1
        }
    }
}

Write-Host "`n=== Base de datos lista para migración ===" -ForegroundColor Green
