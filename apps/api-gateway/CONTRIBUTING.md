# Guía de Contribución y Flujo de Trabajo

## Flujo de Desarrollo de Features

### 1. Actualizar Develop
```bash
git checkout develop
git pull origin develop
```

### 2. Crear Rama Feature
```bash
git checkout -b feature/feature-description
```

### 3. Desarrollo y Commits
- Realiza tus cambios en el código
- Haz commits frecuentes y descriptivos
- Usa mensajes de commit claros y concisos
```bash
git add .
git commit -m "tipo(alcance): descripción"
```

### 4. Crear Pull Request
- Crear PR desde `feature/feature-description` hacia `develop`
- La descripción del PR debe ser clara y detallada
- Incluir información relevante sobre los cambios

### 5. Proceso de Revisión
- Esperar al menos 1 aprobación
- Resolver cualquier comentario o conflicto
- Una vez aprobado, realizar "Squash and Merge" a develop

### 6. Limpieza Post-Merge
```bash
git checkout develop
git pull origin develop
git branch -d feature/feature-description
```

## Flujo de Release

### 1. Preparación Inicial
```bash
git checkout develop
git pull origin develop
```

### 2. Crear Rama Release
```bash
git checkout -b release/dd-mm-yyyy
```

### 3. Sincronizar con Master
```bash
git checkout master
git pull origin master
git checkout release/dd-mm-yyyy
git merge master
git push origin release/dd-mm-yyyy
```

### 4. Proceso de Release
1. Crear Pull Request desde `release/dd-mm-yyyy` hacia `master`
2. Realizar el deploy en los scopes productivos (excepto scope job)
3. Una vez completado el deploy:
   - Mergear `release/dd-mm-yyyy` a `master`
   - Eliminar la rama `release/dd-mm-yyyy`
   - Realizar backport a `develop` y aprobar desde la sección "files"

## Convenciones de Commits

Usa el siguiente formato para los mensajes de commit:
```
tipo(alcance): descripción

[opcional] cuerpo del mensaje

[opcional] pie de página
```

Tipos de commits:
- feat: Nueva característica
- fix: Corrección de bug
- docs: Cambios en documentación
- style: Cambios de formato
- refactor: Refactorización de código
- test: Añadir o modificar tests
- chore: Tareas de mantenimiento

## Buenas Prácticas

1. Mantén las ramas feature pequeñas y enfocadas
2. Haz commits frecuentes y descriptivos
3. Siempre actualiza tu rama con develop antes de crear un PR
4. Revisa tus cambios antes de hacer push
5. Asegúrate de que todos los tests pasen antes de crear un PR
6. En el proceso de release, verifica que no haya procesos críticos ejecutándose antes del deploy
7. Mantén una comunicación clara con el equipo durante el proceso de release 

## Resolución de Conflictos

### Proceso de Resolución
1. Actualizar tu rama con los últimos cambios:
```bash
git checkout develop
git pull origin develop
git checkout feature/feature-description
git merge develop
```

2. Resolver conflictos en archivos:
   - Abrir archivos con conflictos
   - Buscar marcadores de conflicto (`<<<<<<<`, `=======`, `>>>>>>>`)
   - Elegir la versión correcta o combinar cambios
   - Eliminar marcadores de conflicto

3. Finalizar resolución:
```bash
git add .
git commit -m "fix: resolve merge conflicts"
git push origin feature/feature-description
```

## Comandos Git Útiles

### Ver Estado y Cambios
```bash
# Ver estado de cambios
git status

# Ver historial de commits
git log --oneline --graph --all

# Ver cambios en archivos
git diff
```

### Limpieza y Mantenimiento
```bash
# Eliminar ramas locales ya mergeadas
git branch --merged | grep -v "\*" | xargs -n 1 git branch -d

# Limpiar archivos no trackeados
git clean -fd

# Ver ramas remotas
git branch -r
```

## Checklist de Pull Request

### Antes de Crear el PR
- [ ] Código sigue las convenciones del proyecto
- [ ] Tests pasan localmente
- [ ] No hay conflictos con develop
- [ ] Commits siguen el formato convencional
- [ ] Documentación actualizada si es necesario

### En la Descripción del PR
- [ ] Descripción clara del cambio
- [ ] Screenshots si hay cambios UI
- [ ] Pasos para probar
- [ ] Referencia a issues relacionados
- [ ] Impacto en la base de datos (si aplica)

## Proceso de Rollback

### Rollback en Desarrollo
```bash
# Identificar el commit anterior
git log --oneline

# Revertir el último commit
git revert HEAD

# O volver a un commit específico
git revert <commit-hash>
```

### Rollback en Producción
1. Crear rama hotfix:
```bash
git checkout master
git pull origin master
git checkout -b hotfix/rollback-description
```

2. Revertir cambios:
```bash
git revert <commit-hash>
```

3. Deploy y merge:
- Deployar cambios
- Crear PR de hotfix a master
- Mergear y eliminar rama hotfix
- Backport a develop

## Buenas Prácticas Adicionales

8. Mantén un registro de cambios en el CHANGELOG.md
9. Documenta cualquier cambio en la configuración del entorno
10. Notifica al equipo sobre cambios que afecten el flujo de trabajo
11. Realiza code reviews constructivas
12. Mantén las dependencias actualizadas
13. Documenta decisiones técnicas importantes
14. Sigue el principio de "fail fast" en desarrollo
15. Mantén backups de la base de datos antes de migraciones

## Herramientas Recomendadas

### Desarrollo
- IDE: VSCode con extensiones recomendadas
- Linters: golangci-lint
- Formateadores: gofmt
- Testing: go test

### Git
- GitLens para VSCode
- GitKraken para visualización de ramas
- GitHub CLI para operaciones rápidas

### Monitoreo
- Logs de aplicación
- Métricas de rendimiento
- Alertas de errores 