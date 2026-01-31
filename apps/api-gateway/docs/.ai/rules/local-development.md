# 🔧 Reglas de Desarrollo Local

## Configuración Local

### Prerequisitos
- **Go 1.23+**: Versión compatible instalada
- **Docker**: Para ejecutar servicios dependientes
- **Make**: Para comandos de automatización

### Configuración Inicial
```bash
# Clonar dependencias
go mod download

# Ejecutar con Docker Compose (recomendado)
docker-compose up -d

# O ejecutar directamente
make run
```

### Variables de Entorno
Crear archivo `.env` en la raíz del proyecto:
```bash
# Ejemplo de configuración local
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=financial_resume_dev
DB_USER=postgres
DB_PASSWORD=postgres
LOG_LEVEL=debug
```

### Cargar Environment
Seguir las instrucciones del README para cargar el archivo `.env` en tu IDE de preferencia.

## Desarrollo Local vs Testing

### Desarrollo Local
```bash
# Para desarrollo y debugging básico
go test ./... -count=1

# Verificar linters
golangci-lint run
go vet ./...
go fmt ./...

# Ejecutar con hot reload (opcional)
# Instalar air: go install github.com/cosmtrek/air@latest
air
```

### Testing Local
- **Unit tests**: Para validar lógica de negocio
- **Integration tests**: Para validar conectividad con servicios
- **API tests**: Usando Postman collection incluida

## Configuración por IDE

### VSCode/Cursor
```json
// .vscode/launch.json
{
   "version": "0.2.0",
   "configurations": [
       {
           "name": "Launch Financial Resume Engine",
           "type": "go",
           "request": "launch",
           "mode": "auto",
           "program": "${workspaceFolder}/cmd/api",
           "envFile": "${workspaceFolder}/.env"
       }
   ]
}
```

### GoLand/IntelliJ
1. **Run/Debug Configurations**
2. **Go Build**
3. **Package path**: `./cmd/api`
4. **Environment**: Cargar desde archivo `.env`
5. **Working directory**: Raíz del proyecto

### Configuración de Environment Variables
```bash
# Variables principales para desarrollo local
export PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=financial_resume_dev
export LOG_LEVEL=debug
```

## Base de Datos Local

### Configuración con Docker
```bash
# Iniciar PostgreSQL con Docker Compose
docker-compose up -d postgres

# Ejecutar migraciones
make migrate

# Cargar datos de prueba (opcional)
make seed
```

### Configuración Manual
```bash
# Crear base de datos
createdb financial_resume_dev

# Ejecutar script de inicialización
psql -d financial_resume_dev -f scripts/init.sql
```

## Monitoreo Local

### Logs
```bash
# Verificar logs de la aplicación
tail -f logs/app.log

# Ver logs estructurados (si usa JSON)
tail -f logs/app.log | jq .
```

### Health Check
```bash
# Verificar que el servicio está corriendo
curl http://localhost:8080/health

# Verificar endpoints principales
curl http://localhost:8080/api/v1/expenses
curl http://localhost:8080/api/v1/incomes
curl http://localhost:8080/api/v1/dashboard
```

### Swagger UI
```bash
# Acceder a documentación interactiva
open http://localhost:8080/swagger/index.html
```

## Debugging

### Herramientas
- **Delve**: Para debugging de Go
- **IDE Debugger**: Usar breakpoints en IDE
- **Logging**: Agregar logs temporales para debugging

### Casos Comunes
```go
// ✅ Debug logging temporal
logger.Debug("Processing request", 
    zap.Int64("user_id", userID),
    zap.String("step", "validation"),
    zap.Any("input_data", input))
```

## Comandos Útiles

### Makefile
```bash
# Comandos disponibles
make help          # Ver todos los comandos
make build         # Compilar aplicación
make run           # Ejecutar aplicación
make test          # Ejecutar tests
make lint          # Ejecutar linters
make migrate       # Ejecutar migraciones
make docker-build  # Construir imagen Docker
make docker-run    # Ejecutar con Docker
```

## Troubleshooting

### Problemas Comunes
- **Puerto ocupado**: Cambiar puerto en `.env`
- **Base de datos no conecta**: Verificar Docker containers
- **Dependencias**: Ejecutar `go mod download`
- **Permisos**: Verificar permisos de archivos y directorios

```bash
# Verificar estado de servicios
docker-compose ps

# Ver logs de servicios
docker-compose logs postgres
docker-compose logs app

# Reiniciar servicios
docker-compose restart
``` 