# Financial Resume - Hub de Documentación

**Última Actualización**: 2026-02-13
**Estado**: 🚧 Migración en Progreso
**Estado Producción**: ⚠️ En vivo con 1-10 usuarios beta
**Fase Actual**: Fase 3 (Database Consolidation) - ✅ Código Completo, Listo para Ejecución

---

## 📋 Resumen General

Esta documentación captura la ingeniería inversa completa y re-arquitectura de la plataforma Financial Resume. El proyecto está en transición de una **arquitectura de microservicios** hacia un **monolito modular** con Clean Architecture y Domain-Driven Design.

---

## 🎯 Hoja de Ruta de Documentación

### ✅ Completado
- [00-vision/vision.md](./00-vision/vision.md) - Visión del producto, objetivos y dirección estratégica
- [06-data-models/01-current-state/main-db.md](./06-data-models/01-current-state/main-db.md) - Schema actual de base de datos en producción (main DB)
- [06-data-models/01-current-state/gamification-db.md](./06-data-models/01-current-state/gamification-db.md) - Schema actual de base de datos de gamificación
- [03-architecture/01-current-state.md](./03-architecture/01-current-state.md) - Arquitectura actual (Distributed Monolith)
- [03-architecture/02-target-state.md](./03-architecture/02-target-state.md) - Arquitectura objetivo (Modular Monolith)
- [03-architecture/03-migration-plan.md](./03-architecture/03-migration-plan.md) - Plan detallado de migración (7 fases)
- [03-architecture/04-implementation-roadmap.md](./03-architecture/04-implementation-roadmap.md) - Roadmap de implementación semana por semana
- [03-architecture/00-VALIDATION-REPORT.md](./03-architecture/00-VALIDATION-REPORT.md) - Reporte de validación de documentación
- [03-architecture/05-phase3-migration-guide.md](./03-architecture/05-phase3-migration-guide.md) - Guía de ejecución Fase 3
- [03-architecture/PHASE3-STATUS.md](./03-architecture/PHASE3-STATUS.md) - Estado de la Fase 3

### 🔄 En Progreso
- [ ] 06-data-models/01-current-state/domain-models.md - Modelos de dominio (Go structs)
- [ ] 04-api-contracts/current-api.yaml - Contratos API actuales (OpenAPI)
- [ ] 01-requirements.md - Requerimientos funcionales y no funcionales
- [ ] 02-user-stories.md - Historias de usuario y personas

### 📋 Planificado
- [ ] 03-architecture/adr/ - Registros de Decisiones de Arquitectura (ADRs)
- [ ] 06-data-models/02-target-state/ - Schema objetivo de base de datos y modelos de dominio
- [ ] 05-features/ - Especificaciones de features (usando spec-kit)

### 🚀 Estado de la Migración

#### ✅ Fase 1: Monolith Foundation Setup - COMPLETADA
- Estructura del monolith creada
- Event bus in-memory implementado
- Conexión a base de datos
- Servidor HTTP con Gin
- Health check endpoint

#### ✅ Fase 2: Auth Module Migration - COMPLETADA
- Módulo de autenticación completo
- Domain models, services, repository, handlers
- Tests implementados

#### ✅ Fase 3: Database Consolidation - CÓDIGO COMPLETO (Listo para Ejecución)
- ✅ Código de migración completo
- ✅ Comando CLI funcional (`cmd/migrate`)
- ✅ Auditoría pre/post migración
- ✅ Cambios de schema
- ✅ Migración de datos + deduplicación
- ✅ Validación completa
- ✅ Reportes detallados
- 📖 **Guía de ejecución**: [05-phase3-migration-guide.md](./03-architecture/05-phase3-migration-guide.md)
- 📊 **Estado detallado**: [PHASE3-STATUS.md](./03-architecture/PHASE3-STATUS.md)

#### ⏳ Fase 4: Transactions Module Migration - PENDIENTE
- Migrar módulo de transacciones (expenses, incomes, categories)

#### ⏳ Fase 5: Gamification Module + Event Bus - PARCIAL
- ✅ Event bus implementado
- ⏳ Módulo de gamificación pendiente

#### ⏳ Fase 6: Remaining Modules - PENDIENTE
- Budgets, Savings, Recurring, Analytics, AI

#### ⏳ Fase 7: Technical Debt Resolution - PENDIENTE
- Cleanup de código antiguo

---

## 🔍 Hallazgos Clave (Ingeniería Inversa)

### Problemas Críticos Descubiertos

#### 1. ⚠️ **DUPLICACIÓN DE DATOS** (Crítico)
**Problema**: Las tablas `user_gamification`, `achievements`, `user_actions` existen en **AMBAS** bases de datos (main-db y gamification-db).

**Impacto**:
- No hay una fuente única de verdad
- Problemas de sincronización entre bases de datos
- Riesgo de inconsistencia de datos

**Recomendación**: Consolidar en main-db, eliminar el microservicio de gamification (enfoque monolito modular).

---

#### 2. ⚠️ **Sin Eliminación Lógica (Soft Delete)** (Alta Prioridad)
**Problema**: Todas las eliminaciones son físicas (pérdida permanente de datos).

**Impacto**:
- No se pueden restaurar datos eliminados accidentalmente
- No hay auditoría de registros eliminados

**Recomendación**: Agregar columna `deleted_at TIMESTAMP` a tablas críticas.

---

#### 3. ⚠️ **Sin Claves Foráneas a Usuarios** (Prioridad Media)
**Problema**: `user_id` es VARCHAR sin restricción FK (usuarios están en `users_db` separada).

**Impacto**:
- Datos huérfanos si se eliminan usuarios
- Integridad referencial no garantizada a nivel de base de datos

**Recomendación**: Consolidar tabla de usuarios O implementar limpieza a nivel de aplicación.

---

#### 4. ⚠️ **Sin Workers de Fondo** (Prioridad Media)
**Problema**: No hay trabajos cron para transacciones recurrentes o reinicio de períodos de presupuesto.

**Impacto**:
- Se requiere intervención manual para transacciones recurrentes
- Los presupuestos no se reinician automáticamente al final del período

**Recomendación**: Implementar background worker (Go cron o scheduler externo).

---

#### 5. ⚠️ **Formatos de ID Inconsistentes** (Baja Prioridad)
**Problema**: Algunas tablas usan `VARCHAR(36)` (UUIDs), otras usan `VARCHAR(50)` (prefijos custom como `goal_`, `bud_`).

**Impacto**:
- Contratos API inconsistentes
- Confusión para desarrolladores

**Recomendación**: Estandarizar a UUID v4 o ULID en todas las tablas.

---

#### 6. ⚠️ **Sin Row-Level Security** (Seguridad)
**Problema**: No hay aislamiento a nivel de base de datos entre usuarios (depende de cláusulas WHERE en aplicación).

**Impacto**:
- Riesgo de filtración de datos por bugs
- Multi-tenancy no garantizada a nivel de DB

**Recomendación**: Implementar PostgreSQL Row-Level Security (RLS).

---

#### 7. ⚠️ **Sin Logging de Auditoría** (Compliance)
**Problema**: No hay tracking de quién cambió qué y cuándo.

**Impacto**:
- No se pueden debuggear problemas de datos
- No se puede cumplir con requerimientos de auditoría

**Recomendación**: Agregar tabla de audit log o usar event triggers de PostgreSQL.

---

### Resumen de Deuda Técnica

| Categoría | Problema | Prioridad | Esfuerzo | Impacto |
|----------|-------|----------|--------|--------|
| **Integridad de Datos** | Duplicación de datos (gamification) | 🔴 Crítico | Alto | Alto |
| **Integridad de Datos** | Sin soft delete | 🟠 Alto | Medio | Medio |
| **Integridad de Datos** | Sin FK a usuarios | 🟡 Medio | Bajo | Medio |
| **Automatización** | Sin background workers | 🟡 Medio | Alto | Medio |
| **Consistencia** | Formatos de ID inconsistentes | 🟢 Bajo | Bajo | Bajo |
| **Seguridad** | Sin Row-Level Security | 🟠 Alto | Medio | Alto |
| **Compliance** | Sin audit logging | 🟡 Medio | Medio | Medio |

---

## 🏗️ Arquitectura Actual (as-is)

### Microservicios
```
┌─────────────────┐
│  Frontend       │ (React + TailwindCSS)
│  Puerto: 3000   │
└────────┬────────┘
         │ HTTP
         ▼
┌─────────────────┐
│  Nginx Proxy    │
│  Puerto: 8080   │
└────────┬────────┘
         │
         ├──────────────────┬──────────────────┬──────────────────┐
         ▼                  ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│ api-gateway   │  │ users-service │  │ ai-service    │  │ gamification  │
│ Puerto: 8081  │  │ Puerto: 8083  │  │ Puerto: 8082  │  │ Puerto: 8084  │
│ (Finanzas)    │  │ (Auth/JWT/2FA)│  │ (Insights IA) │  │ (XP/niveles)  │
└───────┬───────┘  └───────┬───────┘  └───────┬───────┘  └───────┬───────┘
        │                  │                  │                  │
        ▼                  ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│ postgres-main │  │ users-db      │  │ redis         │  │ gamification  │
│ Puerto: 5432  │  │ Puerto: 5434  │  │ Puerto: 6379  │  │ -db           │
└───────────────┘  └───────────────┘  └───────────────┘  └───────────────┘
```

**Problemas con Arquitectura Actual**:
- ⚠️ **Acoplamiento HTTP síncrono** entre servicios (latencia, fallos en cascada)
- ⚠️ **Duplicación de datos** (tablas de gamification en 2 DBs)
- ⚠️ **Sin arquitectura event-driven** (acoplamiento fuerte)
- ⚠️ **Complejidad operacional** (4 servicios para deployar, monitorear, escalar)

**Ventajas**:
- ✅ Clean Architecture dentro de cada servicio
- ✅ Separación de responsabilidades (auth, finanzas, gamification)

---

## 🎯 Arquitectura Objetivo (to-be)

### Monolito Modular
```
┌─────────────────┐
│  Frontend       │ (React + TailwindCSS)
│  Puerto: 3000   │
└────────┬────────┘
         │ HTTP
         ▼
┌─────────────────────────────────────────────────────────┐
│              Backend v2 (Monolito Modular)               │
│              Go + Clean Architecture + DDD               │
│  Puerto: 8080                                            │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │          Capa HTTP (Handlers/Controllers)           │ │
│  └───────────────────────┬────────────────────────────┘ │
│                          │                              │
│  ┌────────────────────────────────────────────────────┐ │
│  │           Capa de Aplicación (Casos de Uso)         │ │
│  │                                                      │ │
│  │  ┌────────┐  ┌─────────┐  ┌────────┐  ┌──────────┐ │ │
│  │  │Finanzas│  │Usuarios │  │Gamif.  │  │   IA     │ │ │
│  │  │Módulo  │  │ Módulo  │  │ Módulo │  │  Módulo  │ │ │
│  │  └────────┘  └─────────┘  └────────┘  └──────────┘ │ │
│  └───────────────────────┬────────────────────────────┘ │
│                          │                              │
│  ┌────────────────────────────────────────────────────┐ │
│  │         Capa de Dominio (Lógica de Negocio)         │ │
│  │                                                      │ │
│  │  Transaction │ Budget │ SavingsGoal │ User │ etc.  │ │
│  └───────────────────────┬────────────────────────────┘ │
│                          │                              │
│  ┌────────────────────────────────────────────────────┐ │
│  │    Capa de Infraestructura (Repositorios)           │ │
│  └───────────────────────┬────────────────────────────┘ │
└──────────────────────────┼──────────────────────────────┘
                           │
                           ▼
               ┌───────────────────────┐
               │  PostgreSQL (Única)   │
               │  Puerto: 5432         │
               │                       │
               │  - Transacciones      │
               │  - Usuarios           │
               │  - Gamificación       │
               │  - Todos los datos    │
               └───────────────────────┘
```

**Beneficios**:
- ✅ **Deployment simplificado** (un solo servicio)
- ✅ **Sin duplicación de datos** (fuente única de verdad)
- ✅ **Transacciones atómicas** entre módulos
- ✅ **Menor latencia** (sin overhead de HTTP)
- ✅ **Debugging más fácil** (un solo codebase)
- ✅ **Sigue siendo modular** (límites claros entre dominios)

**Cuándo Extraer Microservicios** (si se necesita después):
- Cuando el módulo de IA necesite escalar independientemente
- Cuando gamification necesite ownership de un equipo separado
- Cuando el servicio de usuarios necesite garantías de SLA diferentes

---

## 🗂️ Estructura de Documentación

```
docs/
├── README.md                                    ← Estás aquí
├── 00-vision/
│   └── vision.md                                ✅ Visión del producto, objetivos
├── 01-requirements.md                           📋 Requerimientos funcionales + no funcionales
├── 02-user-stories.md                           📋 Personas, historias, journeys
├── 03-architecture/
│   ├── current-architecture.md                  📋 Microservicios as-is
│   ├── target-architecture.md                   📋 Monolito modular to-be
│   ├── deployment.md                            📋 Docker, CI/CD, infraestructura
│   ├── patterns.md                              📋 DDD, CQRS, patrones aplicados
│   └── adr/                                     📋 Registros de Decisiones de Arquitectura
│       ├── 001-migracion-monolito.md
│       ├── 002-consolidacion-datos.md
│       └── 003-gamificacion-event-driven.md
├── 04-api-contracts/
│   ├── current-api.yaml                         📋 Spec OpenAPI 3.0 (actual)
│   ├── target-api.yaml                          📋 Spec OpenAPI 3.0 (objetivo)
│   └── integration-patterns.md                  📋 Sync vs async, políticas de retry
├── 05-features/
│   ├── transactions/
│   │   ├── spec.md                              📋 Especificación de feature
│   │   ├── plan.md                              📋 Plan de implementación
│   │   └── tasks.md                             📋 Tareas descompuestas
│   ├── budgets/
│   ├── savings-goals/
│   ├── gamification/
│   └── ai-insights/
└── 06-data-models/
    ├── 01-current-state/                        ← ESTADO ACTUAL EN PRODUCCIÓN
    │   ├── main-db.md                           ✅ Schema DB principal (as-is)
    │   ├── gamification-db.md                   ✅ Schema DB gamification (as-is)
    │   ├── domain-models.md                     🔄 Modelos de dominio Go (as-is)
    │   └── data-flow.md                         📋 Cómo fluyen los datos entre servicios
    ├── 02-target-state/                         ← OBJETIVO PARA BACKEND V2
    │   ├── unified-schema.md                    📋 Schema DB consolidado
    │   ├── entity-models.md                     📋 Modelos de entidad DDD (objetivo)
    │   └── improvements.md                      📋 Qué estamos mejorando
    └── 03-migrations/
        ├── strategy.md                          📋 Estrategia de migración
        └── scripts/                             📋 Scripts SQL de migración
```

---

## 🚀 Próximos Pasos

### Fase 1: Completar Documentación (Esta Semana)
- [ ] Terminar de documentar estado actual (modelos de dominio, contratos API)
- [ ] Definir requerimientos desde features existentes
- [ ] Crear historias de usuario
- [ ] Documentar arquitectura actual (microservicios)

### Fase 2: Diseñar Estado Objetivo (Semana Próxima)
- [ ] Diseñar schema objetivo de base de datos (plan de consolidación)
- [ ] Diseñar modelos de dominio objetivo (entidades DDD, value objects, agregados)
- [ ] Diseñar contratos API objetivo (OpenAPI 3.0)
- [ ] Escribir ADRs para decisiones clave (migración a monolito, consolidación de datos, etc.)

### Fase 3: Planificar Migración (Semana 3)
- [ ] Crear documento de estrategia de migración
- [ ] Escribir scripts de migración (SQL) para consolidación de datos
- [ ] Planificar procedimientos de rollback
- [ ] Definir estrategia de testing

### Fase 4: Implementar Backend v2 (Semanas 4-8)
- [ ] Configurar nueva estructura de proyecto monolito modular
- [ ] Implementar capa de dominio (entidades, value objects, agregados)
- [ ] Implementar casos de uso (capa de aplicación)
- [ ] Implementar repositorios (capa de infraestructura)
- [ ] Implementar handlers HTTP (capa de presentación)
- [ ] Escribir tests (unitarios, integración, e2e)

### Fase 5: Migración & Deployment (Semana 9)
- [ ] Ejecutar scripts de migración en staging
- [ ] Testear backend v2 con frontend
- [ ] Programar ventana de mantenimiento (1-2 horas)
- [ ] Ejecutar migración en producción
- [ ] Verificar integridad de datos
- [ ] Monitorear problemas

### Fase 6: Refactor de Frontend (Semana 10)
- [ ] Actualizar llamadas API a nuevos endpoints (si cambiaron)
- [ ] Agregar nuevas features aprovechando capacidades de backend v2
- [ ] Mejorar UX basado en feedback

---

## 📚 Integración con Spec-Kit

Este proyecto usa [spec-kit](https://github.com/anthropics/spec-kit) para especificación y planificación automatizada de features.

**Comandos Disponibles**:
- `/specify` - Crear especificaciones de features
- `/clarify` - Identificar áreas poco especificadas
- `/plan` - Generar planes de implementación
- `/tasks` - Generar tareas accionables
- `/implement` - Ejecutar implementación
- `/analyze` - Analizar consistencia entre artefactos

**Uso**:
```bash
# Ejemplo: Crear spec para feature de transacciones
/specify "Gestión completa de transacciones con CRUD, filtrado y paginación"
```

---

## 🤝 Contribución

Esta documentación es un documento vivo. A medida que descubrimos nuevos insights durante la ingeniería inversa o tomamos decisiones arquitectónicas, actualiza las secciones relevantes.

**Estándares de Documentación**:
- Usar markdown para todos los documentos
- Incluir diagramas (ASCII o Mermaid) donde sea útil
- Mantener docs de estado actual separados de estado objetivo
- Actualizar README.md al agregar nuevas secciones
- Poner fecha a todos los cambios mayores

---

## 📞 Contacto

**Owner del Proyecto**: [Tu Nombre]
**Estado**: Beta (1-10 usuarios)
**Última Revisión**: 2026-02-09

---

**Leyenda**:
- ✅ Completado
- 🔄 En Progreso
- 📋 Planificado
- ⚠️ Problema Crítico
- 🔴 Alta Prioridad
- 🟠 Prioridad Media
- 🟢 Baja Prioridad
