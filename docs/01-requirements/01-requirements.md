# Requerimientos - Financial Resume

**Documento**: Requerimientos Funcionales y No Funcionales
**Versión**: 1.0
**Fecha**: 2026-02-09
**Estado**: Documento Vivo
**Basado en**: Sistema actual en producción (MVP Fase 1)

---

## Tabla de Contenidos

1. [Requerimientos Funcionales](#1-requerimientos-funcionales)
2. [Reglas de Negocio](#2-reglas-de-negocio)
3. [Requerimientos No Funcionales](#3-requerimientos-no-funcionales)
4. [Restricciones](#4-restricciones)
5. [Dependencias](#5-dependencias)
6. [Matriz de Trazabilidad](#6-matriz-de-trazabilidad)

---

## 1. Requerimientos Funcionales

### 1.1 Gestión de Transacciones - Ingresos

**RF-001**: El sistema debe permitir crear ingresos con descripción, monto y categoría opcional, fecha de transaccion (defaul fecha de creacion), dado a que puede crear transacciones pasadas o futuras
**Prioridad**: Alta
**Estado**: 

**RF-002**: El sistema debe permitir editar ingresos existentes modificando todos sus campos
**Prioridad**: Alta
**Estado**: 

**RF-003**: El sistema debe permitir eliminar ingresos con confirmación del usuario
**Prioridad**: Alta
**Estado**: 

**RF-004**: El sistema debe calcular automáticamente el porcentaje de cada ingreso respecto al total de ingresos
**Prioridad**: Media
**Estado**: 

**RF-005**: El sistema debe mostrar los ingresos ordenados por fecha de transaccion, monto categoría
**Prioridad**: Media
**Estado**: 

**RF-006**: El sistema debe permitir buscar ingresos por descripción, fecha
**Prioridad**: Media
**Estado**: 

**RF-007**: El sistema debe permitir filtrar ingresos por mes y año
**Prioridad**: Alta
**Estado**: 

**RF-008**: El sistema debe registrar automáticamente la fecha de creación de cada ingreso
**Prioridad**: Alta
**Estado**: 

---

### 1.2 Gestión de Transacciones - Gastos

**RF-009**: El sistema debe permitir crear gastos con descripción, monto, categoría opcional, fecha de vencimiento estado de pago, prioridad y fecha de transaccion
**Prioridad**: Alta
**Estado**: 

**RF-010**: El sistema debe permitir editar gastos existentes modificando todos sus campos
**Prioridad**: Alta
**Estado**: 

**RF-011**: El sistema debe permitir eliminar gastos con confirmación del usuario
**Prioridad**: Alta
**Estado**: 

**RF-012**: El sistema debe permitir marcar gastos como pagados o pendientes
**Prioridad**: Alta
**Estado**: 

**RF-013**: El sistema debe soportar pagos parciales de gastos
**Prioridad**: Alta
**Estado**: 

**RF-014**: El sistema debe registrar el monto pagado y el monto pendiente de cada gasto
**Prioridad**: Alta
**Estado**: 

**RF-015**: El sistema debe calcular automáticamente el porcentaje de cada gasto respecto al total de ingresos
**Prioridad**: Alta
**Estado**: 

**RF-016**: El sistema debe permitir editar campos individuales de gastos inline (modo Excel)
**Prioridad**: Media
**Estado**: 

**RF-017**: El sistema debe mostrar gastos ordenados por fecha de transaccion, prioridad, fecha de vencimiento, monto o categoría
**Prioridad**: Alta
**Estado**: 

**RF-018**: El sistema debe permitir filtrar gastos por estado de pago (todos, pagados, pendientes)
**Prioridad**: Alta
**Estado**: 

**RF-019**: El sistema debe permitir buscar gastos por descripción
**Prioridad**: Media
**Estado**: 

**RF-020**: El sistema debe permitir filtrar gastos por mes y año
**Prioridad**: Alta
**Estado**: 

**RF-021**: El sistema debe soportar fechas de transacción personalizadas para gastos futuros
**Prioridad**: Alta
**Estado**: 

**RF-022**: El sistema debe detectar y manejar sobrepagos en gastos parciales
**Prioridad**: Media
**Estado**: 

---

### 1.3 Gestión de Categorías

**RF-023**: El sistema debe permitir crear categorías personalizadas con nombre único por usuario
**Prioridad**: Alta
**Estado**: 

**RF-024**: El sistema debe permitir editar el nombre de categorías existentes
**Prioridad**: Alta
**Estado**: 

**RF-025**: El sistema debe permitir eliminar categorías no utilizadas
**Prioridad**: Media
**Estado**: 

**RF-026**: El sistema debe asignar iconos y colores por defecto a las categorías
**Prioridad**: Baja
**Estado**: 

**RF-027**: El sistema debe mostrar colores consistentes para cada categoría en toda la aplicación
**Prioridad**: Media
**Estado**: 

**RF-028**: El sistema debe permitir buscar categorías por nombre
**Prioridad**: Baja
**Estado**: 

**RF-029**: El sistema debe prevenir la creación de categorías duplicadas por usuario
**Prioridad**: Alta
**Estado**: 

---

### 1.4 Gestión de Presupuestos

**RF-030**: El sistema debe permitir crear presupuestos asociados a categorías con monto límite
**Prioridad**: Alta
**Estado**: 

**RF-031**: El sistema debe soportar presupuestos por períodos: semanal, mensual, trimestral, anual
**Prioridad**: Alta
**Estado**: 

**RF-032**: El sistema debe calcular automáticamente el monto gastado vs presupuestado
**Prioridad**: Alta
**Estado**: 

**RF-033**: El sistema debe calcular el porcentaje de uso del presupuesto
**Prioridad**: Alta
**Estado**: 

**RF-034**: El sistema debe permitir configurar umbrales de alerta personalizables (por defecto 80%)
**Prioridad**: Media
**Estado**: 

**RF-035**: El sistema debe cambiar automáticamente el estado del presupuesto según el uso: on_track (<80%), warning (80-99%), exceeded (≥100%)
**Prioridad**: Alta
**Estado**: 

**RF-036**: El sistema debe permitir editar presupuestos existentes (monto, período, umbral)
**Prioridad**: Alta
**Estado**: 

**RF-037**: El sistema debe permitir eliminar presupuestos con confirmación
**Prioridad**: Media
**Estado**: 

**RF-038**: El sistema debe permitir filtrar presupuestos por categoría, período y estado
**Prioridad**: Media
**Estado**: 

**RF-039**: El sistema debe permitir ordenar presupuestos por fecha, nombre, monto o porcentaje usado
**Prioridad**: Baja
**Estado**: 

**RF-040**: El sistema debe mostrar un dashboard resumen con total de presupuestos, cantidad en meta, en alerta y excedidos
**Prioridad**: Alta
**Estado**: 

**RF-041**: El sistema debe generar notificaciones cuando un presupuesto alcanza el umbral de alerta
**Prioridad**: Media
**Estado**:  (estructura DB)

**RF-042**: El sistema debe generar notificaciones cuando un presupuesto es excedido
**Prioridad**: Media
**Estado**:  (estructura DB)

---

### 1.5 Gestión de Metas de Ahorro

**RF-043**: El sistema debe permitir crear metas de ahorro con nombre, monto objetivo, moneda custom y fecha límite
**Prioridad**: Alta
**Estado**: 

**RF-044**: El sistema debe soportar categorías predefinidas para metas: vacation, emergency, house, car, education, retirement, investment, custom
**Prioridad**: Media
**Estado**: 

**RF-045**: El sistema debe permitir asignar íconos personalizados a las metas de ahorro
**Prioridad**: Baja
**Estado**: 

**RF-046**: El sistema debe permitir asignar niveles de prioridad: high, medium, low
**Prioridad**: Media
**Estado**: 

**RF-047**: El sistema debe permitir realizar depósitos en metas de ahorro con descripción opcional
**Prioridad**: Alta
**Estado**: 

**RF-048**: El sistema debe permitir realizar retiros de metas de ahorro con descripción opcional
**Prioridad**: Alta
**Estado**: 

**RF-049**: El sistema debe calcular automáticamente el progreso de la meta como porcentaje
**Prioridad**: Alta
**Estado**: 

**RF-050**: El sistema debe calcular automáticamente el monto restante para alcanzar la meta
**Prioridad**: Alta
**Estado**: 

**RF-051**: El sistema debe calcular automáticamente los días restantes hasta la fecha objetivo
**Prioridad**: Alta
**Estado**: 

**RF-052**: El sistema debe calcular targets de ahorro diario, semanal y mensual basados en el tiempo restante
**Prioridad**: Media
**Estado**: 

**RF-053**: El sistema debe cambiar automáticamente el estado de la meta a "achieved" cuando se alcanza el monto objetivo
**Prioridad**: Alta
**Estado**: 

**RF-054**: El sistema debe soportar ahorro automático programado con frecuencia configurable
**Prioridad**: Media
**Estado**:  (estructura DB, falta cron)

**RF-055**: El sistema debe mostrar el historial completo de transacciones (depósitos/retiros) por meta
**Prioridad**: Alta
**Estado**: 

**RF-056**: El sistema debe permitir editar metas existentes (nombre, monto, fecha, prioridad)
**Prioridad**: Alta
**Estado**: 

**RF-057**: El sistema debe permitir eliminar metas con confirmación, incluyendo su historial
**Prioridad**: Media
**Estado**: 

**RF-058**: El sistema debe permitir pausar y reanudar metas de ahorro
**Prioridad**: Baja
**Estado**:  (estructura DB)

**RF-059**: El sistema debe mostrar un dashboard resumen con total ahorrado, cantidad de metas activas y monto objetivo total
**Prioridad**: Alta
**Estado**: 

**RF-060**: El sistema debe permitir filtrar metas por estado, categoría y prioridad
**Prioridad**: Media
**Estado**: 

---

### 1.6 Gestión de Transacciones Recurrentes

**RF-061**: El sistema debe permitir crear transacciones recurrentes para ingresos o gastos
**Prioridad**: Alta
**Estado**: 

**RF-062**: El sistema debe soportar frecuencias: weekly, monthly, yearly
**Prioridad**: Alta
**Estado**: 

**RF-063**: El sistema debe permitir configurar la próxima fecha de ejecución
**Prioridad**: Alta
**Estado**: 

**RF-064**: El sistema debe permitir configurar fecha de finalización opcional
**Prioridad**: Media
**Estado**: 

**RF-065**: El sistema debe permitir configurar número máximo de ejecuciones opcional
**Prioridad**: Baja
**Estado**: 

**RF-066**: El sistema debe permitir activar/desactivar la creación automática de transacciones
**Prioridad**: Alta
**Estado**: 

**RF-067**: El sistema debe permitir configurar notificaciones previas (días antes de ejecución)
**Prioridad**: Media
**Estado**: 

**RF-068**: El sistema debe crear automáticamente transacciones cuando llega la fecha programada
**Prioridad**: Alta
**Estado**: ⚠️ Requiere cron job

**RF-069**: El sistema debe actualizar automáticamente la próxima fecha de ejecución después de crear una transacción
**Prioridad**: Alta
**Estado**: ⚠️ Requiere cron job

**RF-070**: El sistema debe mantener un registro de ejecuciones con fecha, éxito/fallo y mensaje de error
**Prioridad**: Alta
**Estado**: 

**RF-071**: El sistema debe permitir ejecutar manualmente una transacción recurrente
**Prioridad**: Alta
**Estado**: 

**RF-072**: El sistema debe permitir ejecutar todas las transacciones pendientes en batch
**Prioridad**: Alta
**Estado**: 

**RF-073**: El sistema debe permitir pausar y reanudar transacciones recurrentes
**Prioridad**: Alta
**Estado**: 

**RF-074**: El sistema debe permitir editar transacciones recurrentes
**Prioridad**: Alta
**Estado**: 

**RF-075**: El sistema debe permitir eliminar transacciones recurrentes con confirmación
**Prioridad**: Media
**Estado**: 

**RF-076**: El sistema debe desactivar automáticamente transacciones cuando se alcanza max_executions o end_date
**Prioridad**: Media
**Estado**: ⚠️ Requiere cron job

**RF-077**: El sistema debe mostrar un dashboard resumen con total de transacciones activas, ingresos mensuales proyectados y gastos mensuales proyectados
**Prioridad**: Alta
**Estado**: 

**RF-078**: El sistema debe generar proyecciones de flujo de caja para N meses
**Prioridad**: Media
**Estado**: 

**RF-079**: El sistema debe permitir filtrar transacciones recurrentes por tipo, frecuencia, estado y categoría
**Prioridad**: Media
**Estado**: 

**RF-080**: El sistema debe enviar notificaciones antes de ejecutar transacciones recurrentes
**Prioridad**: Baja
**Estado**:  (estructura DB)

**RF-080.1**: El sistema debe permitir seleccionar recurrencia de las transacciones dentro del CRUD de ingresos y egresos
**Prioridad**: Alta
**Estado**: 

---

### 1.7 Dashboard y Reportes

**RF-081**: El sistema debe calcular y mostrar el balance total (ingresos - gastos)
**Prioridad**: Alta
**Estado**: 

**RF-082**: El sistema debe mostrar total de ingresos del período seleccionado
**Prioridad**: Alta
**Estado**: 

**RF-083**: El sistema debe mostrar total de gastos del período seleccionado
**Prioridad**: Alta
**Estado**: 

**RF-084**: El sistema debe mostrar cantidad de gastos pendientes de pago
**Prioridad**: Alta
**Estado**: 

**RF-085**: El sistema debe mostrar monto total de gastos pendientes
**Prioridad**: Alta
**Estado**: 

**RF-086**: El sistema debe generar gráfico de distribución de gastos por categoría (pie chart)
**Prioridad**: Alta
**Estado**: 

**RF-087**: El sistema debe calcular porcentaje de gasto por categoría respecto al total
**Prioridad**: Alta
**Estado**: 

**RF-088**: El sistema debe mostrar las transacciones más recientes con indicador de pago
**Prioridad**: Alta
**Estado**: 

**RF-089**: El sistema debe calcular total de transacciones del período
**Prioridad**: Media
**Estado**: 

**RF-090**: El sistema debe calcular promedio diario de gastos
**Prioridad**: Media
**Estado**: 

**RF-091**: El sistema debe identificar la categoría con mayor gasto
**Prioridad**: Media
**Estado**: 

**RF-092**: El sistema debe mostrar widgets resumen de presupuestos en el dashboard
**Prioridad**: Media
**Estado**: 

**RF-093**: El sistema debe mostrar widgets resumen de metas de ahorro en el dashboard
**Prioridad**: Media
**Estado**: 

**RF-094**: El sistema debe mostrar widgets resumen de transacciones recurrentes en el dashboard
**Prioridad**: Media
**Estado**: 

**RF-095**: El sistema debe permitir ocultar/mostrar balances (modo privacidad)
**Prioridad**: Media
**Estado**: 

**RF-096**: El sistema debe permitir filtrar datos por año y mes
**Prioridad**: Alta
**Estado**: 

**RF-097**: El sistema debe mostrar transacciones del período en dos columnas (gastos/ingresos)
**Prioridad**: Media
**Estado**: 

**RF-098**: El sistema debe permitir ordenar transacciones por fecha, monto o categoría
**Prioridad**: Media
**Estado**: 

**RF-099**: El sistema debe permitir limitar cantidad de filas mostradas (10, 25, 50, todas)
**Prioridad**: Baja
**Estado**: 

---

### 1.8 Sistema de Financial Health Score

**RF-100**: El sistema debe calcular un Financial Health Score de 1-1000 puntos combinando engagement (40%) y salud financiera (60%)
**Prioridad**: Alta
**Estado**:

**RF-101**: El sistema debe calcular el nivel del usuario basado en el Score acumulado
**Prioridad**: Alta
**Estado**:

**RF-102**: El sistema debe contribuir al componente de Engagement del Score por crear un gasto
**Prioridad**: Media
**Estado**:

**RF-103**: El sistema debe contribuir al componente de Engagement del Score por crear un ingreso
**Prioridad**: Media
**Estado**:

**RF-104**: El sistema debe contribuir al componente de Engagement del Score por crear un presupuesto
**Prioridad**: Media
**Estado**:

**RF-105**: El sistema debe contribuir al componente de Engagement del Score por crear una meta de ahorro
**Prioridad**: Media
**Estado**:

**RF-106**: El sistema debe contribuir al componente de Engagement del Score por ver el dashboard (máximo 1 vez al día)
**Prioridad**: Baja
**Estado**:

**RF-107**: El sistema debe contribuir al componente de Engagement del Score por ver analytics
**Prioridad**: Baja
**Estado**:

**RF-108**: El sistema debe mantener un registro de todas las acciones del usuario con fecha y contribución al Score
**Prioridad**: Media
**Estado**:

**RF-109**: El sistema debe rastrear rachas de días consecutivos con actividad (0-100 puntos del componente Engagement)
**Prioridad**: Alta
**Estado**:

**RF-110**: El sistema debe rastrear consumo de insights vistos por el usuario (0-100 puntos del componente Engagement)
**Prioridad**: Media
**Estado**:

**RF-111**: El sistema debe rastrear adopción de features (0-100 puntos del componente Engagement)
**Prioridad**: Media
**Estado**:

**RF-112**: El sistema debe rastrear challenges completados (0-100 puntos del componente Engagement)
**Prioridad**: Media
**Estado**:

**RF-113**: El sistema debe mantener fecha de última actividad del usuario
**Prioridad**: Media
**Estado**:

**RF-114**: El sistema debe soportar logros/badges con progreso y target
**Prioridad**: Media
**Estado**:

**RF-115**: El sistema debe desbloquear logros automáticamente al alcanzar el target
**Prioridad**: Media
**Estado**:

**RF-116**: El sistema debe registrar fecha de desbloqueo de logros
**Prioridad**: Baja
**Estado**:

**RF-117**: El sistema debe soportar challenges diarios y semanales
**Prioridad**: Media
**Estado**:

**RF-118**: El sistema debe rastrear progreso de challenges activos
**Prioridad**: Media
**Estado**:

**RF-119**: El sistema debe contribuir al Score al completar challenges
**Prioridad**: Media
**Estado**:

**RF-120**: El sistema debe resetear challenges diarios/semanales automáticamente
**Prioridad**: Media
**Estado**: ⚠️ Requiere cron job

**RF-121**: El sistema debe bloquear features por nivel de usuario basado en el Score
**Prioridad**: Alta
**Estado**:

**RF-122**: El sistema debe desbloquear Presupuestos en Nivel 2 (401 Score)
**Prioridad**: Alta
**Estado**:

**RF-123**: El sistema debe desbloquear Metas de Ahorro en Nivel 3 (501 Score)
**Prioridad**: Alta
**Estado**:

**RF-124**: El sistema debe desbloquear IA Financiera en Nivel 4 (601 Score)
**Prioridad**: Alta
**Estado**: 

---

### 1.9 Insights y Analytics con IA

**RF-125**: El sistema debe calcular el componente de Salud Financiera del Score (600 puntos máximo) basado en: Savings Rate (0-150), Budget Adherence (0-150), Goals Progress (0-150), Expense Management (0-150)
**Prioridad**: Alta
**Estado**:

**RF-126**: El sistema debe clasificar al usuario en niveles basados en el Score total: Nivel 1 Aprendiz (1-400), Nivel 2 Manager (401-500), Nivel 3 Guru (501-600), Nivel 4 Master (601-900), Nivel 5 Magnate (901-1000)
**Prioridad**: Alta
**Estado**:

**RF-127**: El sistema debe analizar patrones de gasto por categoría para el componente de Salud Financiera
**Prioridad**: Media
**Estado**:

**RF-127.1**: El sistema debe diferenciar categorías entre pasivos (egresos de dinero) y activos (ahorros, inversiones) para el cálculo del Score
**Prioridad**: Media
**Estado**:

**RF-127.2**: El sistema debe analizar las metas de ahorro para calcular el componente Goals Progress del Score
**Prioridad**: Alta
**Estado**:

**RF-128**: El sistema debe generar recomendaciones financieras personalizadas basadas en el Score y sus componentes
**Prioridad**: Media
**Estado**:

**RF-129**: El sistema debe detectar anomalías en gastos que afecten el componente Expense Management
**Prioridad**: Baja
**Estado**:

**RF-130**: El sistema debe generar pronósticos de flujo de caja considerando el comportamiento histórico del Score
**Prioridad**: Media
**Estado**:

**RF-131**: El sistema debe mostrar el Financial Health Score en el dashboard principal con desglose de componentes
**Prioridad**: Alta
**Estado**: 

---

### 1.10 Gestión de Usuarios y Autenticación

**RF-132**: El sistema debe permitir registro de usuarios con email y contraseña
**Prioridad**: Alta
**Estado**: 

**RF-133**: El sistema debe permitir login con email y contraseña
**Prioridad**: Alta
**Estado**: 

**RF-134**: El sistema debe generar tokens JWT para autenticación
**Prioridad**: Alta
**Estado**: 

**RF-135**: El sistema debe validar tokens en todas las peticiones protegidas
**Prioridad**: Alta
**Estado**: 

**RF-136**: El sistema debe cerrar sesión y eliminar tokens
**Prioridad**: Alta
**Estado**: 

**RF-137**: El sistema debe mantener sesión persistente (localStorage/sessionStorage)
**Prioridad**: Alta
**Estado**: 

**RF-138**: El sistema debe aislar datos por usuario (multi-tenancy)
**Prioridad**: Alta
**Estado**: 

**RF-139**: El sistema debe inicializar perfil de gamificación al crear usuario
**Prioridad**: Media
**Estado**: 

**RF-139.1**: El sistema debe permitir compartir la cuenta financiera en varios perfiles
**Prioridad**: Media
**Estado**: 
---

## 2. Reglas de Negocio

### RN-001: Cálculo de Porcentajes de Transacciones
- Los porcentajes de transacciones se calculan como: `(monto_transacción / total_ingresos) * 100`
- Si no hay ingresos registrados, el porcentaje es 0%
- Los porcentajes se calculan dinámicamente en cada request

### RN-002: Balance Financiero
- Balance = Total de Ingresos - Total de Gastos
- El balance se calcula sobre las transacciones del período seleccionado
- Los gastos parcialmente pagados cuentan por su monto total, no por el monto pagado

### RN-003: Estados de Pago de Gastos
- Un gasto está `paid = true` cuando `amount_paid >= amount`
- Un gasto está `paid = false` cuando `amount_paid < amount`
- `pending_amount = amount - amount_paid`

### RN-004: Pagos Parciales
- Se permite registrar pagos parciales que incrementan `amount_paid`
- Si un pago parcial hace que `amount_paid >= amount`, el gasto se marca automáticamente como `paid = true`
- Si el usuario intenta pagar más del monto pendiente, se ofrecen dos opciones:
  - Aumentar el monto total del gasto al monto del pago
  - Aplicar pago total con el monto original

### RN-005: Estados de Presupuestos
- `on_track`: `spent_amount / amount < alert_at` (por defecto < 80%)
- `warning`: `spent_amount / amount >= alert_at AND < 1.0` (80% - 99%)
- `exceeded`: `spent_amount / amount >= 1.0` (≥ 100%)
- El estado se calcula automáticamente en cada actualización

### RN-006: Cálculo de Progreso de Metas de Ahorro
- `progress = current_amount / target_amount` (valor decimal 0.0 - 1.0)
- `remaining_amount = target_amount - current_amount`
- `days_remaining = (target_date - today)`
- `daily_target = remaining_amount / days_remaining`
- `weekly_target = daily_target * 7`
- `monthly_target = remaining_amount / (days_remaining / 30)`

### RN-007: Estados de Metas de Ahorro
- `active`: Meta en progreso
- `achieved`: `current_amount >= target_amount`
- `paused`: Meta pausada manualmente
- `cancelled`: Meta cancelada manualmente
- Cuando una meta alcanza el 100%, cambia automáticamente a `achieved` y se registra `achieved_at`

### RN-008: Transacciones Recurrentes - Próxima Ejecución
- Después de cada ejecución exitosa, `next_date` se calcula según `frequency`:
  - `daily`: `next_date + 1 day`
  - `weekly`: `next_date + 7 days`
  - `monthly`: `next_date + 1 month` (mismo día)
  - `yearly`: `next_date + 1 year` (mismo día y mes)
- `execution_count` se incrementa después de cada ejecución
- Si `execution_count >= max_executions` o `next_date > end_date`, la transacción se desactiva automáticamente

### RN-009: Sistema de Financial Health Score y Niveles
- **Nivel 1: Aprendiz**: 1-400 Score 🌱
- **Nivel 2: Manager**: 401-500 Score 💼 (desbloquea Presupuestos)
- **Nivel 3: Guru**: 501-600 Score 🧠 (desbloquea Metas de Ahorro)
- **Nivel 4: Master**: 601-900 Score 👑 (desbloquea IA Financiera)
- **Nivel 5: Magnate**: 901-1000 Score 💎
- El nivel se calcula automáticamente basado en `financial_health_score`

### RN-010: Contribución al Financial Health Score
El Score se compone de dos componentes principales:

**Componente Engagement (400 puntos máximo - 40%)**:
- **Racha de días**: 0-100 puntos (basado en días consecutivos de actividad)
- **Adopción de features**: 0-100 puntos (uso de categorías, presupuestos, metas)
- **Consumo de insights**: 0-100 puntos (vistas de dashboard, analytics, reportes)
- **Challenges completados**: 0-100 puntos (challenges diarios y semanales)

**Componente Salud Financiera (600 puntos máximo - 60%)**:
- **Savings Rate**: 0-150 puntos (porcentaje de ingresos ahorrados)
- **Budget Adherence**: 0-150 puntos (adherencia a presupuestos establecidos)
- **Goals Progress**: 0-150 puntos (progreso en metas de ahorro)
- **Expense Management**: 0-150 puntos (control de gastos, gastos pendientes, recurrentes)

### RN-011: Validación de Categorías
- Los nombres de categoría deben ser únicos por usuario
- Longitud mínima: 2 caracteres
- Longitud máxima: 50 caracteres
- No se permiten caracteres especiales excepto espacios, guiones y guiones bajos

### RN-012: Validación de Montos
- Todos los montos deben ser números positivos mayores a 0
- Máximo 2 decimales
- Formato: DECIMAL(10,2) para transacciones
- Formato: DECIMAL(15,2) para presupuestos y metas

### RN-013: Validación de Descripciones
- Longitud mínima: 3 caracteres
- Longitud máxima: 255 caracteres
- No se permiten solo espacios en blanco

### RN-014: Filtrado de Datos por Período
- Cuando se selecciona un año: se filtran transacciones del 1 de enero al 31 de diciembre
- Cuando se selecciona un mes: se filtran transacciones del día 1 al último día del mes
- El filtrado se aplica sobre el campo `created_at` de las transacciones

### RN-015: Isolation de Datos por Usuario
- Todas las queries deben incluir filtro `WHERE user_id = :user_id`
- No se permite acceso cruzado entre datos de diferentes usuarios
- El `user_id` se extrae del token JWT autenticado

### RN-016: Soft Delete vs Hard Delete
- **Actualmente**: Todos los deletes son físicos (hard delete)
- Se pierde el historial al eliminar datos
- **Recomendación futura**: Implementar soft delete con campo `deleted_at`

### RN-017: Frecuencia de Actualización de Dashboard
- Los datos del dashboard se recargan automáticamente después de mutaciones
- El caché se invalida después de operaciones CRUD en transacciones, presupuestos o metas
- El sistema emite eventos de sincronización para refrescar vistas relacionadas

### RN-018: Cálculo del Financial Health Score
El Financial Health Score se calcula combinando dos componentes principales:

**Fórmula Total**:
```
Financial Health Score = Engagement Component + Health Component
Donde:
  - Engagement Component: 0-400 puntos (40%)
  - Health Component: 0-600 puntos (60%)
  - Total: 1-1000 puntos
```

**Engagement Component (400 pts)**:
- Racha de días (0-100): Basado en current_streak, con bonificación por rachas largas
- Adopción de features (0-100): Uso activo de categorías, presupuestos, metas, transacciones recurrentes
- Consumo de insights (0-100): Frecuencia de vistas de dashboard, analytics, reportes
- Challenges completados (0-100): Tasa de completación de challenges diarios/semanales

**Health Component (600 pts)**:
- Savings Rate (0-150): `(total_ahorrado / total_ingresos) * 150`
- Budget Adherence (0-150): Promedio de presupuestos en estado "on_track" vs "exceeded"
- Goals Progress (0-150): Promedio de progreso en metas activas ponderado por prioridad
- Expense Management (0-150): Basado en: ratio gastos/ingresos, gastos pendientes, control de recurrentes

Escala: 1-1000 puntos (nunca 0 para evitar desmotivación)

---

## 3. Requerimientos No Funcionales

### 3.1 Performance

**RNF-001**: Las operaciones CRUD de transacciones deben responder en menos de 100ms (percentil 95)
**Estado**: 

**RNF-002**: La carga del dashboard debe completarse en menos de 500ms con hasta 1000 transacciones
**Estado**: 

**RNF-003**: Los reportes complejos (analytics, proyecciones) deben generarse en menos de 2 segundos
**Estado**: 

**RNF-004**: El sistema debe soportar paginación con límite de 100 items por página
**Estado**: 

**RNF-005**: Las queries de base de datos deben tener índices optimizados para operaciones frecuentes
**Estado**: 

**RNF-006**: El sistema debe implementar caché para datos de solo lectura con TTL de 5 minutos
**Estado**: 

**RNF-007**: Las mutaciones deben invalidar el caché de forma inteligente para mantener consistencia
**Estado**: 

**RNF-008**: La UI debe implementar actualizaciones optimistas para mejorar percepción de velocidad
**Estado**:  

**RNF-009**: El sistema debe soportar edición inline de campos sin recargar la página completa
**Estado**: 

---

### 3.2 Seguridad

**RNF-010**: El sistema debe usar JWT para autenticación de usuarios
**Estado**: 

**RNF-011**: Los tokens JWT deben expirar después de 24 horas
**Estado**: 

**RNF-012**: Las contraseñas deben ser hasheadas con bcrypt (mínimo 10 rounds)
**Estado**: 

**RNF-013**: Todas las comunicaciones deben ser sobre HTTPS en producción
**Estado**: 

**RNF-014**: Las API keys y secrets deben ser almacenadas en variables de entorno
**Estado**: 

**RNF-015**: El sistema debe validar y sanitizar todas las entradas de usuario
**Estado**: 

**RNF-016**: El sistema debe prevenir inyección SQL usando prepared statements
**Estado**:  (ORM/query builder)

**RNF-017**: El sistema debe prevenir XSS escapando contenido HTML
**Estado**:  (React)

**RNF-018**: El sistema debe implementar rate limiting para prevenir abuso (100 req/min por usuario)
**Estado**: 

**RNF-019**: El sistema debe implementar CORS restrictivo para permitir solo orígenes autorizados
**Estado**: 

**RNF-020**: El sistema debe aislar datos entre usuarios (multi-tenancy a nivel de aplicación)
**Estado**: 

**RNF-021**: El sistema debe proteger rutas del frontend con guardias de autenticación
**Estado**: 

**RNF-022**: El sistema debe validar permisos en cada endpoint del backend
**Estado**: 

---

### 3.3 Escalabilidad

**RNF-023**: El sistema debe soportar hasta 100 usuarios concurrentes sin degradación
**Estado**: 

**RNF-024**: El sistema debe manejar hasta 10,000 transacciones por usuario sin problemas de performance
**Estado**: 

**RNF-025**: La base de datos debe soportar crecimiento horizontal mediante sharding por user_id
**Estado**: 

**RNF-026**: El backend debe ser stateless para permitir escalamiento horizontal
**Estado**: 

**RNF-027**: El sistema debe usar connection pooling para la base de datos
**Estado**: 

**RNF-028**: El sistema debe implementar timeouts en todas las operaciones de I/O (5s para DB, 10s para APIs externas)
**Estado**: 

---

### 3.4 Disponibilidad y Confiabilidad

**RNF-029**: El sistema debe tener un uptime objetivo de 99.5% (tolerancia de downtime: 3.65 horas/mes)
**Estado**: 

**RNF-030**: El sistema debe implementar health checks para monitoreo
**Estado**: 

**RNF-031**: El sistema debe registrar errores y excepciones en logs estructurados
**Estado**: 

**RNF-032**: El sistema debe tener backups automáticos diarios de la base de datos
**Estado**: 

**RNF-033**: El sistema debe recuperarse automáticamente de errores transitorios (retry con exponential backoff)
**Estado**: 

**RNF-034**: El sistema debe validar integridad de datos antes de mutaciones críticas
**Estado**: 

**RNF-035**: El sistema debe usar transacciones de base de datos para operaciones atómicas
**Estado**: 

---

### 3.5 Usabilidad

**RNF-036**: La interfaz debe ser responsive y funcionar en móviles, tablets y desktop
**Estado**: 

**RNF-037**: El sistema debe soportar modo oscuro (dark mode)
**Estado**: 

**RNF-038**: El sistema debe mostrar mensajes de error claros y accionables
**Estado**: 

**RNF-039**: El sistema debe mostrar feedback visual inmediato para acciones del usuario (toasts, spinners)
**Estado**: 

**RNF-040**: El sistema debe prevenir doble-submit de formularios
**Estado**: 

**RNF-041**: El sistema debe preservar el estado del scroll después de mutaciones
**Estado**:  

**RNF-042**: El sistema debe mostrar indicadores de progreso para operaciones largas (>2s)
**Estado**: 

**RNF-043**: El sistema debe usar confirmaciones para acciones destructivas (delete)
**Estado**: 

**RNF-044**: El sistema debe permitir deshacer acciones críticas
**Estado**: 

**RNF-045**: La interfaz debe usar colores consistentes para categorías en toda la aplicación
**Estado**: 

**RNF-046**: El sistema debe mostrar tooltips explicativos en acciones no obvias
**Estado**:  

**RNF-047**: El sistema debe mostrar placeholders informativos cuando no hay datos
**Estado**: 

---

### 3.6 Mantenibilidad

**RNF-048**: El código debe seguir principios SOLID y de Clean Architecture (Domain → Use Cases → Infrastructure)
**Estado**: 

**RNF-049**: El código debe tener separación clara entre backend (Go) y frontend (React)
**Estado**: 

**RNF-050**: El código debe usar nombres descriptivos y convenciones consistentes
**Estado**: 

**RNF-051**: El sistema debe documentar APIs con comentarios y ejemplos
**Estado**: 

**RNF-052**: El sistema debe mantener migraciones versionadas para cambios de esquema
**Estado**: 

**RNF-053**: El código debe tener cobertura de tests unitarios >80%
**Estado**: 

**RNF-054**: El código debe tener tests de integración para flujos críticos
**Estado**: 

**RNF-055**: El sistema debe usar linting y formateo automático (golangci-lint, prettier)
**Estado**: 

---

### 3.7 Portabilidad

**RNF-056**: El sistema debe funcionar en PostgreSQL 13+
**Estado**: 

**RNF-057**: El backend debe ser compilable para Linux, macOS y Windows
**Estado**: ✅ Go es multiplataforma

**RNF-058**: El frontend debe funcionar en Chrome, Firefox, Safari y Edge (últimas 2 versiones)
**Estado**: 

**RNF-059**: El sistema debe soportar despliegue en contenedores Docker
**Estado**: 

**RNF-060**: El sistema debe tener variables de entorno configurables para diferentes ambientes
**Estado**: 

---

### 3.8 Observabilidad

**RNF-061**: El sistema debe registrar logs con niveles: DEBUG, INFO, WARN, ERROR
**Estado**: 

**RNF-062**: Los logs deben incluir timestamp, nivel, mensaje y contexto (user_id, request_id)
**Estado**: 

**RNF-063**: El sistema debe exponer métricas de performance (latencias, throughput)
**Estado**: 

**RNF-064**: El sistema debe registrar eventos de auditoría para acciones críticas
**Estado**:  (user_actions, pero no completo)

**RNF-065**: El sistema debe tener dashboards de monitoreo para operación
**Estado**: 

---

## 4. Restricciones

### 4.1 Restricciones Técnicas

**REST-001**: Backend desarrollado en **Go** con arquitectura limpia (DDD)
**REST-002**: Frontend desarrollado en **React 18+** con Hooks
**REST-003**: Base de datos **PostgreSQL 13+**
**REST-004**: Autenticación mediante **JWT**
**REST-005**: Diseño responsive con **Tailwind CSS**
**REST-006**: Gráficos con **Recharts**
**REST-007**: Gestión de estado con **React Context API**
**REST-008**: No se usan frameworks de backend adicionales (monolito modular puro)
**REST-009**: No se usa Redux ni gestión de estado compleja
**REST-010**: No se usa ORM pesado (queries SQL directas con pgx)

---

### 4.2 Restricciones Operacionales

**REST-011**: Deployment en **Render.com** (PaaS)
**REST-012**: Base de datos hosteada en **Supabase** (PostgreSQL managed)
**REST-013**: Downtime tolerable: **1-2 horas** para mantenimiento planificado
**REST-014**: Usuarios actuales: **1-10** (fase beta/MVP)
**REST-015**: Presupuesto de infraestructura: **Mínimo** (free tier cuando sea posible)
**REST-016**: No se requiere CDN en esta fase
**REST-017**: No se requiere auto-scaling en esta fase
**REST-018**: Backups dependientes de proveedor de infraestructura

---

### 4.3 Restricciones de Negocio

**REST-019**: El sistema está en **fase MVP/beta** (no producción masiva)
**REST-020**: No se requiere soporte 24/7
**REST-021**: No se requiere SLA formal en esta fase
**REST-022**: No se requiere cumplimiento de regulaciones financieras (no es institución financiera)
**REST-023**: Los datos no se comparten con terceros
**REST-024**: No se implementa facturación/pagos en esta fase (modelo freemium futuro)
**REST-025**: No se requiere soporte multi-idioma en esta fase (solo español)
**REST-026**: No se requiere soporte multi-moneda en esta fase (moneda única)

---

### 4.4 Restricciones de Tiempo

**REST-027**: El MVP debe estar completo para **Q2 2026**
**REST-028**: Fase 2 (Inteligencia IA) planificada para **6 meses** 
**REST-029**: No hay deadline rígido para Fase 2
**REST-030**: El desarrollo es incremental e iterativo

---

## 5. Dependencias

### 5.1 Dependencias Externas

**DEP-001**: **PostgreSQL** - Base de datos principal
**DEP-002**: **Render.com** - Plataforma de hosting
**DEP-003**: **JWT** - Librería de autenticación
**DEP-004**: **React** - Framework de frontend
**DEP-005**: **Recharts** - Librería de gráficos
**DEP-006**: **Tailwind CSS** - Framework de estilos
**DEP-007**: **Go** - Lenguaje de backend
**DEP-008**: Integración con **servicios de IA** (OpenAI, Anthropic) en Fase 2
**DEP-009**: Posible integración con **servicios de notificaciones** (email, push) en Fase 2
**DEP-010**: Posible integración con **servicios de monitoreo** (Sentry, DataDog) en futuro
**DEP-011**:**Supabase** - proveedor de base de datos

---

### 5.2 Dependencias Internas

**DEP-011**: Los **presupuestos** dependen de la existencia de **categorías**
**DEP-012**: Los **gastos e ingresos** pueden asociarse opcionalmente a **categorías**
**DEP-013**: El cálculo de **porcentajes de gastos** depende del **total de ingresos**
**DEP-014**: Las **transacciones recurrentes** generan **ingresos o gastos** automáticamente
**DEP-015**: El **score de IA** depende de tener datos de **transacciones, presupuestos y metas**
**DEP-016**: La **gamificación** depende del registro de **acciones de usuario**
**DEP-017**: Los **niveles de usuario** controlan el acceso a **features bloqueadas**
**DEP-018**: El **dashboard** depende de datos de **todas las entidades** principales

---

### 5.3 Dependencias de Infraestructura

**DEP-019**: El **backend** requiere conexión a **PostgreSQL**
**DEP-020**: El **frontend** requiere acceso al **backend API**
**DEP-021**: Las **transacciones recurrentes** requieren un **cron job o scheduler**
**DEP-022**: Los **challenges diarios/semanales** requieren un **cron job** para reset
**DEP-023**: El **reinicio automático de presupuestos** requiere un **cron job**
**DEP-024**: Las **notificaciones** requieren un **servicio de email/push** (futuro)
**DEP-025**: El **backup automático** depende de **servicios de Supabase.com**

---

## 6. Matriz de Trazabilidad

| Feature | Requerimientos Funcionales | Reglas de Negocio | Casos de Uso | Endpoints API | Tablas DB |
|---------|---------------------------|-------------------|--------------|---------------|-----------|
| **Ingresos** | RF-001 a RF-008 | RN-001, RN-012, RN-013, RN-015 | UC-001 a UC-004 | `/api/incomes` | `incomes` |
| **Gastos** | RF-009 a RF-022 | RN-001, RN-003, RN-004, RN-012, RN-013, RN-015 | UC-005 a UC-012 | `/api/expenses` | `expenses` |
| **Categorías** | RF-023 a RF-029 | RN-011, RN-015 | UC-013 a UC-016 | `/api/categories` | `categories` |
| **Presupuestos** | RF-030 a RF-042 | RN-005, RN-012, RN-015 | UC-017 a UC-024 | `/api/budgets` | `budgets`, `budget_notifications` |
| **Metas de Ahorro** | RF-043 a RF-060 | RN-006, RN-007, RN-012, RN-015 | UC-025 a UC-036 | `/api/savings-goals` | `savings_goals`, `savings_transactions` |
| **Transacciones Recurrentes** | RF-061 a RF-080 | RN-008, RN-012, RN-015 | UC-037 a UC-048 | `/api/recurring-transactions` | `recurring_transactions`, `recurring_transaction_executions`, `recurring_transaction_notifications` |
| **Dashboard y Reportes** | RF-081 a RF-099 | RN-001, RN-002, RN-014, RN-017 | UC-049 a UC-054 | `/api/dashboard` | Múltiples tablas |
| **Financial Health Score** | RF-100 a RF-124, RF-125 a RF-131 | RN-009, RN-010, RN-018 | UC-055 a UC-070 | `/api/gamification`, `/api/health-score` | `user_gamification`, `achievements`, `user_actions`, `challenges`, `user_challenges` |
| **Autenticación** | RF-132 a RF-139 | RN-015 | UC-071 a UC-075 | `/api/auth` | `users` (DB separada) |

---

## 7. Identificación de Gaps y Mejoras

### 7.1 Gaps Funcionales Identificados

**GAP-001**: **Transacciones Recurrentes** - Cron job para ejecución automática
**Impacto**: Alto - Feature crítica parcialmente funcional
**Recomendación**: Implementar worker con cron job o servicio de scheduling

**GAP-002**: **Challenges Diarios/Semanales** - Cron job para reset automático
**Impacto**: Medio - Afecta experiencia de gamificación
**Recomendación**: Implementar worker para reset de challenges

**GAP-003**: **Presupuestos** - Falta reinicio automático al final del período
**Impacto**: Medio - Requiere acción manual del usuario
**Recomendación**: Implementar worker para rollover de presupuestos

**GAP-004**: **Ahorro Automático** - Estructura en DB pero no implementado
**Impacto**: Bajo - Feature nice-to-have
**Recomendación**: Implementar en Fase 2

**GAP-005**: **Notificaciones** - Estructura en DB pero no se envían
**Impacto**: Medio - Afecta engagement
**Recomendación**: Integrar servicio de email/push en Fase 2

**GAP-006**: **Detección de Anomalías** - Parcialmente implementado
**Impacto**: Bajo - Feature de IA avanzada
**Recomendación**: Mejorar algoritmos en Fase 2

**GAP-007**: **Pronóstico de Flujo de Caja** - Parcialmente implementado
**Impacación**: Bajo - Feature de IA avanzada
**Recomendación**: Mejorar modelos predictivos en Fase 2

---

### 7.2 Gaps No Funcionales Identificados

**GAP-NF-001**: **Rate Limiting** - No implementado
**Impacto**: Alto - Riesgo de abuso
**Recomendación**: Implementar middleware de rate limiting

**GAP-NF-002**: **Tests Automatizados** - Cobertura insuficiente
**Impacto**: Alto - Riesgo de regresiones
**Recomendación**: Establecer cobertura mínima de 80%

**GAP-NF-003**: **Auditoría Completa** - Solo parcial (user_actions)
**Impacto**: Medio - Dificulta debugging y compliance
**Recomendación**: Agregar audit log comprehensivo

**GAP-NF-004**: **Health Checks** - Parcialmente implementado
**Impacto**: Medio - Dificulta monitoreo
**Recomendación**: Implementar endpoints de health checks

**GAP-NF-005**: **Métricas y Observabilidad** - No implementadas
**Impacto**: Medio - Dificulta optimización
**Recomendación**: Integrar sistema de métricas (Prometheus, DataDog)

**GAP-NF-006**: **Soft Delete** - No implementado
**Impacto**: Medio - Pérdida irreversible de datos
**Recomendación**: Agregar campo `deleted_at` en tablas críticas

**GAP-NF-007**: **Row-Level Security** - No implementado
**Impacto**: Medio - Dependencia en capa de aplicación
**Recomendación**: Considerar PostgreSQL RLS para seguridad adicional

**GAP-NF-008**: **Containerización** - No configurado
**Impacto**: Bajo - Dificulta despliegues consistentes
**Recomendación**: Crear Dockerfile y docker-compose

---

### 7.3 Deuda Técnica Identificada

**TD-001**: **Duplicación de Datos de Gamificación** - Tablas duplicadas entre main-db y gamification-db
**Impacto**: Alto - Riesgo de inconsistencia
**Recomendación**: Consolidar en una sola DB en Backend v2

**TD-002**: **IDs Inconsistentes** - Mezcla de UUID y IDs con prefijos
**Impacto**: Medio - Confusión en contratos API
**Recomendación**: Estandarizar a UUID v4 o ULID

**TD-003**: **No hay Foreign Keys a Users** - `user_id` es VARCHAR sin FK
**Impacto**: Medio - Riesgo de datos huérfanos
**Recomendación**: Consolidar tabla users o implementar cleanup application-layer

**TD-004**: **Historial de Pagos Parciales** - Solo se guarda `amount_paid` actual
**Impacto**: Bajo - No se puede rastrear historial
**Recomendación**: Agregar tabla `expense_payments` en futuro

**TD-005**: **Transacciones DB no Atómicas** - Algunas operaciones no usan transactions
**Impacto**: Medio - Riesgo de inconsistencia
**Recomendación**: Envolver operaciones críticas en transactions

---

## 8. Conclusiones y Próximos Pasos

### 8.1 Estado Actual del Sistema

El sistema **Financial Resume** se encuentra en **Fase MVP completada** con las siguientes características:

✅ **Fortalezas**:
- Arquitectura limpia y bien estructurada (DDD)
- Features core completas (transacciones, presupuestos, metas, recurrentes)
- Gamificación funcional con sistema de niveles
- UI responsive y moderna con dark mode
- Performance aceptable (<100ms en operaciones CRUD)
- Seguridad básica implementada (JWT, bcrypt, validaciones)

⚠️ **Áreas de Mejora**:
- Falta automatización completa (cron jobs para recurrentes, challenges, presupuestos)
- Cobertura de tests insuficiente
- Observabilidad limitada
- Deuda técnica acumulada (duplicación de datos, inconsistencias)

---

### 8.2 Recomendaciones Prioritarias

**Prioridad Alta (Antes de Fase 2)**:
1. Implementar cron jobs para transacciones recurrentes
2. Consolidar datos de gamificación (eliminar duplicación)
3. Implementar rate limiting
4. Agregar cobertura de tests >80%
5. Implementar soft delete en tablas críticas

**Prioridad Media (Durante Fase 2)**:
6. Mejorar sistema de notificaciones
7. Implementar health checks y métricas
8. Estandarizar IDs (UUID)
9. Agregar auditoría completa
10. Implementar Row-Level Security

**Prioridad Baja (Post-Fase 2)**:
11. Containerización con Docker
12. Historial de pagos parciales
13. Mejoras de IA (anomalías, pronósticos)
14. Features sociales y colaborativas

---

### 8.3 Roadmap de Implementación

**Q1 2026** ✅ Completado:
- MVP Fase 1 con todas las features core

**Q2 2026** (En progreso):
- Implementar cron jobs faltantes
- Mejorar observabilidad
- Consolidar arquitectura de datos
- Agregar tests automatizados

**Q3 2026** (Planificado - Fase 2):
- Mejorar insights de IA
- Implementar notificaciones completas
- Optimizar performance para escala
- Preparar para beta pública

**Q4 2026** (Futuro):
- App móvil (iOS/Android)
- Features sociales
- Plan premium
- Integraciones externas

---

**Documento elaborado**: 2026-02-09
**Próxima revisión**: 2026-03-09
**Mantenido por**: Equipo de Producto

---

## Requisitos de Refactoring (RR)

> Requisitos técnicos para la migración de Distributed Monolith a Modular Monolith.

### RR-001: Consolidación a Binary Único
**Prioridad**: Alta
**Descripción**: Reemplazar los 4 binarios Go + nginx + supervisord por un único binary Go ejecutable.
**Criterio de Aceptación**: Un solo proceso en producción, sin supervisord, sin nginx interno.

### RR-002: Eliminación de HTTP Interno
**Prioridad**: Alta
**Descripción**: Reemplazar todas las llamadas HTTP entre servicios internos (gamification, ai, users) por llamadas directas a funciones/interfaces Go.
**Criterio de Aceptación**: Sin llamadas HTTP a localhost dentro del proceso. Latencia interna < 1ms.

### RR-003: Consolidación de Bases de Datos
**Prioridad**: Alta
**Descripción**: Unificar main-db y gamification-db en una sola base de datos PostgreSQL.
**Criterio de Aceptación**: Un solo connection string de DB. Sin duplicación de tablas. 0 pérdida de datos.

### RR-004: Estandarización de UserID
**Prioridad**: Media
**Descripción**: Estandarizar el tipo de UserID a string (UUID) en todos los modelos del dominio.
**Criterio de Aceptación**: UserID es string en 100% de los modelos. Sin conversiones de tipo en runtime.

### RR-005: Implementación de Soft Delete
**Prioridad**: Media
**Descripción**: Reemplazar borrado físico por soft delete (campo deleted_at) en todas las entidades principales.
**Criterio de Aceptación**: Expenses, Incomes, Budgets, Categories tienen campo deleted_at. Las queries filtran deleted_at IS NULL por defecto.

### RR-006: Event Bus Interno
**Prioridad**: Media
**Descripción**: Implementar un event bus in-memory para comunicación asíncrona entre módulos (especialmente para gamification).
**Criterio de Aceptación**: Módulo Gamification recibe eventos via bus, no via HTTP. Sin pérdida de eventos en operación normal.

### RR-007: Preservación de Datos durante Migración
**Prioridad**: Alta
**Descripción**: La migración debe preservar todos los datos de usuario existentes. Las bases de datos originales deben mantenerse intactas como backup natural (la nueva DB unificada se crea por separado).
**Criterio de Aceptación**: 0 pérdida de datos. Bases originales sin modificar hasta verificación completa de la nueva DB.

### RR-008: Compatibilidad de API
**Prioridad**: Alta
**Descripción**: Mantener 100% compatibilidad de la API REST pública durante y después de la migración.
**Criterio de Aceptación**: Todos los endpoints existentes responden con el mismo contrato (paths, métodos, request/response schemas).

### RR-009: Reducción de Complejidad Operacional
**Prioridad**: Media
**Descripción**: Simplificar el deployment eliminando nginx, supervisord y la gestión de múltiples procesos.
**Criterio de Aceptación**: Dockerfile con single-stage build. 1 proceso. Sin configuración nginx/supervisord.

### RR-010: Arquitectura Modular con Interfaces
**Prioridad**: Media
**Descripción**: Cada módulo del monolito debe exponer interfaces Go en lugar de HTTP endpoints para comunicación interna.
**Criterio de Aceptación**: Cada módulo tiene puertos/interfaces definidos. Sin imports cruzados directos entre módulos (solo via interfaces).
