# 🔄 Reglas de Proceso de Desarrollo

## Branching Strategy

### Formato de Ramas
```bash
# Formato recomendado
feature/FRE-[número_tarea]-[título-descriptivo]

# Ejemplos
feature/FRE-001-analytics-dashboard
feature/FRE-002-expense-categories
feature/FRE-003-income-tracking
```

### Flujo de Branches
- **Base Branch**: Todas las ramas deben generarse desde `main`
- **Target Branch**: PRs deben apuntar a `main`
- **Approvals**: Mínimo 1 approval para merge
- **Merge Strategy**: Usar "Squash and Merge" para features
- **Epic Branch**: Para features grandes que requieren múltiples ramas paralelas:
  - Crear una rama `feature/FRE-epic-[nombre]` desde `main`
  - Las ramas individuales se crean desde esta epic branch
  - Los PRs se apuntan a la epic branch
  - Una vez completadas todas las features, la epic branch se mergea a `main`

## Versionado Semántico

Seguir el formato `MAJOR.MINOR.PATCH`:

- **MAJOR**: Cambios incompatibles en API (rompe compatibilidad)
- **MINOR**: Nuevas funcionalidades compatibles hacia atrás
- **PATCH**: Correcciones de errores compatibles

### Ejemplos
```
1.0.0 → 1.0.1  (patch: bug fix)
1.0.1 → 1.1.0  (minor: nueva feature)
1.1.0 → 2.0.0  (major: breaking change)
```

## Deployment Strategy

### Estrategias de Merge por Branch
- **feature → main**: SQUASH
- **hotfix → main**: MERGE
- **release → main**: MERGE

### Proceso de Release
1. Crear release branch desde main
2. PR a main con formato: `Release - X.Y.Z - dd-mm-yyyy`
3. Deploy usando Docker
4. Verificar en logs de aplicación
5. Merge a main con "Create a merge commit"

## Pre-commit Checklist

- [ ] ¿La rama sigue el formato `feature/FRE-xxx`?
- [ ] ¿El PR apunta a `main`?
- [ ] ¿Tiene reviewer asignado?
- [ ] ¿Sigue las reglas arquitectónicas?
- [ ] ¿Pasan todos los linters?
- [ ] ¿Tiene tests unitarios?
- [ ] ¿Se actualiza la documentación Swagger si hay cambios en API? 