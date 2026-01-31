# 🤖 Directorio .ai

Este directorio contiene configuraciones y reglas universales para IDEs con capacidades de AI (Cursor, Windsurf, Codeium, GitHub Copilot, etc.).

## 📁 Estructura

```
.ai/
├── README.md              # Este archivo - Guía principal + índice de reglas
└── rules/                 # Reglas organizadas por temas
    ├── architecture.md    # Arquitectura y Clean Architecture
    ├── development.md     # Proceso de desarrollo y branching
    ├── go-standards.md    # Estándares específicos de Go
    ├── testing.md         # Testing, mocks y coverage
    ├── local-development.md # Desarrollo local y configuración
    └── validation.md      # Validaciones y seguridad
```

## 🎯 Propósito

- **Universalidad**: Compatible con cualquier IDE con AI
- **Consistencia**: Mismas reglas independientemente de la herramienta
- **Simplicidad**: Enfocado en lo esencial para el desarrollo

## 📖 Reglas por Temas

### 🏗️ [Arquitectura](./rules/architecture.md)
- Estructura del proyecto Financial Resume Engine
- Principios de Clean Architecture  
- Patrones de diseño
- Dependency management con Go modules

### 🔄 [Proceso de Desarrollo](./rules/development.md)
- Branching strategy (feature/FRE-xxx)
- Versionado semántico
- Proceso de deployment
- Estrategias de merge

### 💻 [Estándares de Código Go](./rules/go-standards.md)
- Herramientas obligatorias (golangci-lint, gofmt)
- Buenas prácticas de código
- Manejo de errores y context
- Logging estructurado

### 🧪 [Testing](./rules/testing.md)
- Comandos de testing (Go nativo)
- Convenciones de mocks
- Estrategias de testing y coverage
- Testing local

### 🔧 [Desarrollo Local](./rules/local-development.md)
- Configuración de desarrollo local
- Variables de entorno
- Configuración por IDE
- Debugging y troubleshooting

### 🚨 [Validaciones](./rules/validation.md)
- Input validation y sanitización
- Security validations
- Business logic validation
- Error handling y logging

## 📖 Uso

### Para IDEs con AI
Los IDEs pueden leer automáticamente las reglas de este directorio para proporcionar asistencia consistente.

### Para Desarrolladores
Consulta los archivos específicos arriba para revisar las normas y buenas prácticas del proyecto Financial Resume Engine.

## 📝 Contexto del Proyecto

Para información detallada sobre arquitectura, APIs y documentación técnica, consulta la carpeta `/docs`.

## 🔧 Mantenimiento

- Actualiza las reglas cuando cambien los estándares del equipo
- Mantén las reglas simples y enfocadas en desarrollo
- El contexto específico del proyecto debe estar en `/docs`

## 🔗 Referencias Externas

- **Contexto del proyecto**: Ver `/docs` para arquitectura y APIs
- **Docker**: Para desarrollo local con contenedores
- **Swagger**: Para documentación de APIs (ver `/docs/swagger.yaml`) 