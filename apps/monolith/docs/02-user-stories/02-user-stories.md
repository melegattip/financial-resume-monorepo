# Historias de Usuario - Financial Resume

**Fecha**: 2026-02-09
**Estado**: Documento Vivo
**Versión**: 1.0 (MVP Actual)

---

## Tabla de Contenidos

1. [Personas](#personas)
2. [Épicas y User Stories](#épicas-y-user-stories)
3. [Journey Maps](#journey-maps)
4. [Backlog Priorizado](#backlog-priorizado)
5. [Métricas de Éxito](#métricas-de-éxito)

---

## Personas

### 1. El Ahorrador Aspirante - Ana (28 años)

**Perfil**:
- Profesional joven, desarrolladora de software
- Ingreso estable pero variable (freelance ocasional)
- Vive sola en departamento alquilado
- Usa apps móviles diariamente, muy tech-savvy

**Objetivos**:
- Ahorrar para viaje a Europa (meta: $800,000 en 8 meses)
- Crear fondo de emergencia ($200,000)
- Reducir gastos impulsivos en delivery y suscripciones

**Frustraciones**:
- Las apps tradicionales son aburridas y las abandona después de 2 semanas
- No sabe exactamente en qué gasta su dinero
- Pierde motivación cuando no ve progreso
- Los presupuestos tradicionales se sienten restrictivos

**Comportamiento**:
- Revisa su teléfono al menos 50 veces al día
- Le gustan las apps con gamificación (Duolingo, Strava)
- Comparte logros en redes sociales
- Prefiere visualizaciones simples y coloridas

**Tech Savviness**: Alta (9/10)

---

### 2. El Padre/Madre Consciente del Presupuesto - Carlos (35 años)

**Perfil**:
- Ingeniero casado con 2 hijos (4 y 7 años)
- Ingreso familiar estable pero gastos variables
- Responsable de finanzas del hogar
- Usuario intermedio de tecnología

**Objetivos**:
- Controlar gastos del hogar (supermercado, colegio, servicios)
- Evitar sorpresas de fin de mes
- Ahorrar para educación de los hijos
- Mantener presupuesto mensual bajo control

**Frustraciones**:
- Gastos inesperados que desbalancean el presupuesto
- No tiene tiempo para hojas de cálculo complicadas
- Necesita visibilidad rápida del estado financiero
- Los gastos de su pareja no están centralizados (futuro)

**Comportamiento**:
- Revisa finanzas los domingos por la noche
- Prefiere notificaciones de alerta para no exceder presupuestos
- Busca simplicidad y rapidez
- Valora informes mensuales automáticos

**Tech Savviness**: Media (6/10)

---

### 3. El Optimizador Financiero - Laura (32 años)

**Perfil**:
- Contadora independiente
- Múltiples fuentes de ingreso (trabajo fijo + consultorías)
- Ya hace tracking financiero, busca mejor análisis
- Power user de hojas de cálculo y dashboards

**Objetivos**:
- Identificar patrones de gasto para optimizar
- Proyectar flujo de caja a 6 meses
- Automatizar gastos recurrentes (alquiler, servicios)
- Maximizar ahorro mediante análisis de datos

**Frustraciones**:
- El análisis manual consume mucho tiempo
- Las apps básicas no tienen suficiente profundidad
- Necesita insights accionables, no solo gráficos bonitos
- Quiere entender su comportamiento financiero

**Comportamiento**:
- Revisa dashboards diariamente
- Exporta datos para análisis avanzados
- Configura múltiples presupuestos y alertas
- Utiliza todas las funcionalidades disponibles

**Tech Savviness**: Muy Alta (10/10)

---

## Épicas y User Stories

### Épica 1: Gestión de Transacciones

#### US-001: Registrar Gasto

**Como** Ana (Ahorrador Aspirante)
**Quiero** registrar un gasto inline rápidamente desde el dashboard
**Para** mantener un registro actualizado de mis gastos sin fricción

**Criterios de Aceptación**:
- [x] Puedo ingresar monto, descripción, categoría y prioridad
- [x] Puedo seleccionar fecha del gasto (pasado o presente)
- [x] El sistema valida que el monto sea positivo
- [x] Veo confirmación visual después de guardar
- [x] El gasto aparece inmediatamente en mi lista
- [x] Puedo asignar una fecha de vencimiento opcional
- [x] Puedo marcar el gasto como pagado o pendiente o asignarle un pago parcial
- [x] El formulario tiene validación en tiempo real

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Gestión de Transacciones
**Estado**:

**Notas Técnicas**:
- Validación: monto > 0, descripción requerida (max 255 chars)
- Campos opcionales: category_id, due_date
- Actualización optimista de UI para feedback inmediato

---

#### US-002: Editar Gasto Inline

**Como** Laura (Optimizador Financiero)
**Quiero** editar gastos directamente en la lista sin abrir formularios
**Para** corregir errores rápidamente sin interrumpir mi flujo

**Criterios de Aceptación**:
- [x] Puedo hacer clic en la descripción para editarla inline
- [x] Puedo hacer clic en el monto para modificarlo inline
- [x] Puedo cambiar la categoría desde un dropdown inline
- [x] Puedo cambiar la fecha de vencimiento inline
- [x] Puedo cambiar la prioridad inline
- [x] Los cambios se guardan con Enter o al perder foco
- [x] Puedo cancelar con Escape
- [x] Veo feedback visual del guardado
- [x] Los cambios son instantáneos (actualización optimista)

**Prioridad**: Media
**Story Points**: 5
**Épica**: Gestión de Transacciones
**Estado**: 

**Pain Points Resueltos**:
- Elimina necesidad de abrir modal completo para ediciones menores
- Mejora velocidad para usuarios power como Laura

---

#### US-003: Gestionar Pagos Parciales

**Como** Carlos (Padre Consciente)
**Quiero** registrar pagos parciales de un gasto
**Para** llevar control de gastos que pago en cuotas

**Criterios de Aceptación**:
- [x] Puedo marcar un gasto como "Pago Parcial"
- [x] Puedo ingresar el monto pagado hasta el momento
- [x] El sistema calcula automáticamente el saldo pendiente
- [x] Veo visualmente cuánto he pagado vs. cuánto falta
- [x] Puedo hacer múltiples pagos parciales
- [x] Cuando el saldo llega a 0, el gasto se marca como "Pagado"
- [x] Puedo revertir un pago y volver a marcar como pendiente

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Gestión de Transacciones
**Estado**: 

**Ejemplo de Uso**:
- Carlos tiene un gasto de $50,000 en reparaciones del auto
- Paga $20,000 inicialmente → Sistema muestra "Pendiente: $30,000"
- Paga $15,000 más → Sistema muestra "Pendiente: $15,000"
- Paga $15,000 finales → Sistema marca como "Pagado Completamente"

---

#### US-004: Registrar Ingreso

**Como** Laura (Optimizador Financiero)
**Quiero** registrar ingresos con categorías personalizadas
**Para** tener visibilidad completa de mis múltiples fuentes de ingresos

**Criterios de Aceptación**:
- [x] Puedo crear un ingreso con descripción, monto y categoría
- [x] Puedo asignar categorías personalizadas (Sueldo, Freelance, Inversiones)
- [x] El sistema calcula automáticamente el porcentaje de cada ingreso sobre el total
- [x] Veo el ingreso inmediatamente en el dashboard
- [x] Puedo editar o eliminar ingresos
- [x] Los ingresos se consideran en el balance total

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Gestión de Transacciones
**Estado**: 

---

#### US-005: Gestión de Categorías

**Como** Ana (Ahorrador Aspirante)
**Quiero** crear categorías personalizadas con iconos y colores
**Para** organizar mis gastos de forma visual y significativa

**Criterios de Aceptación**:
- [x] Puedo crear categorías personalizadas
- [x] Puedo asignar nombre, icono emoji y color a cada categoría
- [x] Puedo editar o eliminar categorías existentes
- [x] El sistema valida que no haya nombres duplicados
- [x] Al eliminar una categoría, los gastos asociados no se pierden
- [x] Veo las categorías organizadas por color en reportes

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Gestión de Transacciones
**Estado**: 

---

#### US-006: Filtrado y Ordenamiento de Gastos

**Como** Carlos (Padre Consciente)
**Quiero** filtrar gastos por fecha, categoría y estado de pago
**Para** encontrar rápidamente transacciones específicas

**Criterios de Aceptación**:
- [x] Puedo filtrar gastos por mes, año o rango custom
- [x] Puedo filtrar por estado: Todos, Pagados, Pendientes
- [x] Puedo ordenar por fecha, monto, categoría, prioridad
- [x] Puedo cambiar el orden (ascendente/descendente)
- [x] Los filtros se aplican inmediatamente sin recargar página
- [x] Puedo combinar múltiples filtros
- [x] Los filtros persisten si navego a otra página y vuelvo

**Prioridad**: Media
**Story Points**: 3
**Épica**: Gestión de Transacciones
**Estado**: 

---

### Épica 2: Presupuestos Inteligentes

#### US-010: Crear Presupuesto por Categoría

**Como** Carlos (Padre Consciente)
**Quiero** establecer un límite mensual para cada categoría de gasto
**Para** evitar gastar de más en áreas específicas

**Criterios de Aceptación**:
- [x] Puedo crear un presupuesto seleccionando categoría, monto y período
- [x] Los períodos disponibles son: Semanal, Mensual, Trimestral, Anual
- [x] Puedo definir el umbral de alerta (default 80%)
- [x] El sistema calcula automáticamente cuánto he gastado vs. el límite
- [x] Veo una barra de progreso visual del presupuesto
- [x] El presupuesto se marca con colores: Verde (ok), Amarillo (alerta), Rojo (excedido)

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Presupuestos Inteligentes
**Estado**: 

**Niveles de Feature Gate**:
- Requerido: Nivel 2 (500 XP)
- Beneficio: Control inteligente de gastos

---

#### US-011: Alertas de Presupuesto

**Como** Carlos (Padre Consciente)
**Quiero** recibir alertas cuando esté cerca de exceder un presupuesto
**Para** ajustar mis gastos antes de que sea tarde

**Criterios de Aceptación**:
- [x] Recibo una notificación cuando alcanzo el 80% del presupuesto (configurable) y esta se guarda en la seccion notificaciones
- [x] Recibo una notificación cuando excedo el 100%
- [x] Las notificaciones aparecen en el dashboard
- [x] Puedo ver historial de notificaciones
- [x] Las notificaciones incluyen: categoría, monto gastado, monto límite
- [x] Puedo marcar notificaciones como leídas

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Presupuestos Inteligentes
**Estado**: 

**Nota**: Las notificaciones actualmente son in-app. Push notifications pendiente para v2.

---

#### US-012: Dashboard de Presupuestos

**Como** Laura (Optimizador Financiero)
**Quiero** ver todos mis presupuestos en un dashboard unificado
**Para** entender rápidamente mi situación financiera general

**Criterios de Aceptación**:
- [x] Veo un resumen con: Total presupuestos, En meta, Con alerta, Excedidos
- [x] Veo lista completa de presupuestos con estado actual
- [x] Puedo filtrar por categoría, período, estado
- [x] Puedo ordenar por monto gastado, porcentaje usado, fecha
- [x] Veo el período activo de cada presupuesto (fecha inicio - fecha fin)
- [x] Puedo acceder rápidamente a editar o eliminar presupuestos inline

**Prioridad**: Media
**Story Points**: 3
**Épica**: Presupuestos Inteligentes
**Estado**: 

---

### Épica 3: Ahorro con Objetivos

#### US-015: Crear Meta de Ahorro

**Como** Ana (Ahorrador Aspirante)
**Quiero** crear una meta de ahorro con nombre, monto objetivo y fecha límite
**Para** visualizar mi progreso hacia objetivos específicos

**Criterios de Aceptación**:
- [x] Puedo crear una meta con nombre descriptivo
- [x] Puedo elegir un icono/emoji y categoría (Vacaciones, Casa, Auto, etc.)
- [x] Puedo establecer monto objetivo
- [x] Puedo establecer fecha límite
- [x] El sistema calcula cuánto debo ahorrar diariamente/semanalmente/mensualmente y me lo sugiere en un pequeño lugar al pasar el cursor sobre un icono
- [x] Veo una barra de progreso visual
- [x] Puedo editar la meta en cualquier momento

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Ahorro con Objetivos
**Estado**: 

**Niveles de Feature Gate**:
- Requerido: Nivel 3 (1000 XP)
- Beneficio: Metas de ahorro personalizadas

---

#### US-016: Depositar en Meta de Ahorro

**Como** Ana (Ahorrador Aspirante)
**Quiero** registrar depósitos en mi meta de ahorro
**Para** ver cómo avanzo hacia mi objetivo

**Criterios de Aceptación**:
- [x] Puedo hacer un depósito con monto y descripción opcional
- [x] El saldo actual de la meta se actualiza inmediatamente
- [x] El progreso se recalcula automáticamente
- [x] Veo un historial completo de depósitos con fechas
- [x] Puedo hacer múltiples depósitos pequeños
- [x] El sistema celebra cuando completo la meta (100%)

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Ahorro con Objetivos
**Estado**: 

---

#### US-017: Retirar de Meta de Ahorro

**Como** Carlos (Padre Consciente)
**Quiero** retirar dinero de una meta de ahorro cuando lo necesite
**Para** usar los ahorros ante emergencias

**Criterios de Aceptación**:
- [x] Puedo registrar un retiro con monto y razón
- [x] El saldo actual de la meta se reduce
- [x] El progreso se recalcula automáticamente
- [x] Veo retiros en el historial diferenciados de depósitos
- [x] Puedo retirar hasta el monto disponible (no negativo)
- [x] El sistema muestra claramente cuánto queda disponible

**Prioridad**: Media
**Story Points**: 3
**Épica**: Ahorro con Objetivos
**Estado**: 

---

#### US-018: Dashboard de Metas

**Como** Ana (Ahorrador Aspirante)
**Quiero** ver todas mis metas en una vista principal
**Para** motivarme viendo mi progreso general

**Criterios de Aceptación**:
- [x] Veo total ahorrado entre todas las metas
- [x] Veo lista de metas con progreso visual (%)
- [x] Veo cuánto falta para cada meta
- [x] Puedo ver detalle de una meta con historial completo
- [x] Puedo crear nuevas metas desde el dashboard
- [x] La UI es motivacional y celebra el progreso

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Ahorro con Objetivos
**Estado**: 

**Nota de UX**: El diseño usa emojis grandes, barras de progreso gradiente y celebraciones visuales. (ver diseño nueva version de mercadopago)

---

### Épica 4: Transacciones Recurrentes

#### US-020: Crear Transacción Recurrente

**Como** Laura (Optimizador Financiero)
**Quiero** configurar ingresos y gastos que se repiten automáticamente
**Para** no tener que registrarlos manualmente cada mes

**Criterios de Aceptación**:
- [x] Puedo crear una transacción recurrente (ingreso o gasto)
- [x] Puedo elegir frecuencia: Semanal, Mensual, Anual
- [x] Puedo establecer fecha de próxima ejecución o ejecutar todas hasta el vencimiento
- [x] Puedo establecer fecha de fin opcional
- [x] Puedo activar/desactivar auto-creación de transacciones
- [x] Puedo configurar notificación previa (default 1 día antes)
- [x] El sistema valida que la fecha de fin sea posterior al inicio

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Transacciones Recurrentes
**Estado**: 

**Ejemplos de Uso**:
- Sueldo mensual: $150,000 cada 1 de mes
- Alquiler: $60,000 cada 10 de mes
- Netflix: $3,000 cada mes
- Gimnasio: $5,000 semanal

---

#### US-021: Pausar/Reanudar Transacción Recurrente

**Como** Carlos (Padre Consciente)
**Quiero** pausar temporalmente una transacción recurrente
**Para** cuando cancelo una suscripción temporalmente

**Criterios de Aceptación**:
- [x] Puedo pausar una transacción activa
- [x] La transacción pausada no se ejecuta automáticamente
- [x] Veo claramente el estado "Pausada" en la lista
- [x] Puedo reanudar una transacción pausada
- [x] Al reanudar, puedo actualizar la próxima fecha de ejecución
- [x] No pierdo el historial de ejecuciones previas
- [x] Al reanudar, no se hacen ejecuciones retroactivas, solo desde la fecha actual en adelante

**Prioridad**: Media
**Story Points**: 2
**Épica**: Transacciones Recurrentes
**Estado**: 

---

#### US-022: Ejecutar Transacción Vencida

**Como** Carlos (Padre Consciente)
**Quiero** ejecutar manualmente transacciones vencidas que no se procesaron
**Para** regularizar gastos atrasados

**Criterios de Aceptación**:
- [x] Veo claramente cuáles transacciones están vencidas (fecha pasada)
- [x] Puedo ejecutar una transacción vencida con un botón "⚡ Ejecutar"
- [x] La ejecución crea la transacción correspondiente (gasto o ingreso)
- [x] La próxima fecha se calcula automáticamente según frecuencia
- [x] Veo confirmación del resultado de la ejecución
- [x] Se registra en el historial de ejecuciones

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Transacciones Recurrentes
**Estado**: 

**Nota Técnica**: Botón "🔄 Ejecutar Pendientes" permite procesar todas las vencidas en lote.

---

#### US-023: Proyección de Flujo de Caja

**Como** Laura (Optimizador Financiero)
**Quiero** proyectar mi flujo de caja futuro basado en transacciones recurrentes
**Para** planificar mejor mis finanzas

**Criterios de Aceptación**:
- [x] Puedo generar una proyección a 3, 6 o 12 meses
- [x] Veo ingresos proyectados por mes
- [x] Veo gastos proyectados por mes
- [x] Veo balance neto mensual (ingresos - gastos)
- [x] Veo balance acumulado mes a mes
- [x] La proyección considera las frecuencias reales de cada transacción
- [x] Puedo exportar o imprimir la proyección

**Prioridad**: Media
**Story Points**: 5
**Épica**: Transacciones Recurrentes
**Estado**: 

**Diferenciador**: Calcula frecuencias exactas (diaria=x30, semanal=x4.33, anual=/12)

---

#### US-024: Dashboard de Transacciones Recurrentes

**Como** Laura (Optimizador Financiero)
**Quiero** ver resumen de todas mis transacciones recurrentes
**Para** entender mi flujo de caja mensual automático

**Criterios de Aceptación**:
- [x] Veo total de transacciones activas e inactivas
- [x] Veo total de ingresos recurrentes mensuales
- [x] Veo total de gastos recurrentes mensuales
- [x] Puedo filtrar por tipo (ingreso/gasto), frecuencia, estado, categoría
- [x] Puedo ordenar por próxima ejecución, monto, fecha de creación
- [x] Veo días restantes hasta próxima ejecución de cada transacción

**Prioridad**: Media
**Story Points**: 3
**Épica**: Transacciones Recurrentes
**Estado**: 

---

### Épica 5: Sistema de Financial Health Score

#### US-030: Financial Health Score y Niveles

**Como** Ana (Ahorrador Aspirante)
**Quiero** ver mi Financial Health Score que refleje tanto mi engagement como mi salud financiera real
**Para** tener una visión integral de mi progreso y mantenerme motivada

**Criterios de Aceptación**:
- [x] Veo mi Financial Health Score total (1-1000 puntos) en el dashboard
- [x] Veo el desglose del Score: Componente Engagement (400 pts) y Componente Salud Financiera (600 pts)
- [x] El componente Engagement refleja: Racha de días, Adopción de features, Consumo de insights, Challenges completados
- [x] El componente Salud Financiera refleja: Savings Rate, Budget Adherence, Goals Progress, Expense Management
- [x] Veo mi nivel actual basado en el Score total
- [x] Veo cuántos puntos necesito para el siguiente nivel
- [x] Veo una barra de progreso hacia el próximo nivel
- [x] Mis acciones financieras contribuyen al componente de Engagement
- [x] Mi comportamiento financiero real contribuye al componente de Salud Financiera
- [x] Recibo insights personalizados basados en ambos componentes del Score

**Prioridad**: Alta
**Story Points**: 8
**Épica**: Financial Health Score
**Estado**:

**Niveles y Rangos de Score**:
- Nivel 1: Aprendiz (1-400 Score) 🌱
- Nivel 2: Manager (401-500 Score) 💼 - Desbloquea Presupuestos
- Nivel 3: Guru (501-600 Score) 🧠 - Desbloquea Metas de Ahorro
- Nivel 4: Master (601-900 Score) 👑 - Desbloquea IA Financiera
- Nivel 5: Magnate (901-1000 Score) 💎

**Nota**: El Financial Health Score combina gamificación con análisis real de salud financiera para una métrica más significativa
---

#### US-031: Logros y Badges

**Como** Ana (Ahorrador Aspirante)
**Quiero** desbloquear logros por alcanzar hitos
**Para** tener objetivos claros y celebrar mis avances

**Criterios de Aceptación**:
- [x] Veo lista de logros disponibles y mi progreso en cada uno
- [x] Veo logros completados con fecha de desbloqueo
- [x] Cada logro contribuye a mejorar mi Score
- [x] Veo notificación cuando desbloqueo un logro
- [x] Los logros tienen categorías: AI Partner, Action Taker, Data Explorer, Financial Health
- [x] Puedo ver descripción y requerimientos de cada logro
- [x] Los logros contribuyen al componente de Engagement del Score

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Financial Health Score
**Estado**:

**Ejemplos de Logros**:
- "Primera Transacción": Registra tu primer gasto o ingreso (mejora Adopción de features)
- "Presupuesto Maestro": Crea tu primer presupuesto (mejora Adopción de features)
- "Ahorrador Iniciado": Completa tu primera meta de ahorro (mejora Goals Progress)
- "Racha de 7 Días": Usa la app 7 días consecutivos (mejora Racha de días)
- "Score 500": Alcanza 500 puntos de Score total (celebración de hito)

---

#### US-032: Challenges Diarios y Semanales

**Como** Ana (Ahorrador Aspirante)
**Quiero** completar desafíos diarios y semanales
**Para** mantener el hábito de usar la app y mejorar mi Score

**Criterios de Aceptación**:
- [x] Veo challenges disponibles del día y semana actual
- [x] Veo mi progreso en cada challenge en tiempo real
- [x] Los challenges se reinician automáticamente (diarios a las 00:00, semanales los lunes)
- [x] Completar challenges contribuye al componente "Challenges completados" del Score (0-100 pts)
- [x] Veo notificación cuando completo un challenge
- [x] Los challenges son variados: registrar transacciones, usar categorías, revisar analytics, mejorar salud financiera
- [x] Los challenges contribuyen directamente a mejorar mi Financial Health Score

**Prioridad**: Media
**Story Points**: 5
**Épica**: Financial Health Score
**Estado**:

**Ejemplos de Challenges**:
- Diario: Registra 3 transacciones hoy (mejora Adopción de features)
- Diario: Usa 3 categorías diferentes (mejora Adopción de features)
- Diario: Revisa dashboard y analytics (mejora Consumo de insights)
- Semanal: Registra 15 transacciones esta semana (mejora Adopción de features)
- Semanal: Inicia sesión 5 días esta semana (mejora Racha de días)
- Semanal: Mantén tus presupuestos en verde (mejora Budget Adherence)

---

#### US-033: Sistema de Rachas

**Como** Ana (Ahorrador Aspirante)
**Quiero** mantener una racha de días consecutivos usando la app
**Para** crear un hábito sostenible y mejorar mi Score

**Criterios de Aceptación**:
- [x] El sistema cuenta días consecutivos con al menos 1 acción
- [x] Veo mi racha actual en el perfil y su contribución al Score
- [x] La racha contribuye de 0-100 puntos al componente de Engagement
- [x] La racha se reinicia si paso un día sin actividad
- [x] Recibo notificación para no perder mi racha
- [x] Hay logros asociados a rachas (7, 30, 90 días) que mejoran el Score
- [x] Rachas más largas contribuyen proporcionalmente más al Score

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Financial Health Score
**Estado**:

**Nota**: Las rachas ayudan a formar hábitos según psicología del comportamiento y representan hasta 25% del componente de Engagement.

---

#### US-034: Feature Gates por Nivel de Score

**Como** Ana (Ahorrador Aspirante)
**Quiero** desbloquear funcionalidades avanzadas al mejorar mi Score
**Para** tener objetivos progresivos y aprender gradualmente

**Criterios de Aceptación**:
- [x] Nivel 1 Aprendiz (1-400 Score): Acceso a transacciones básicas y categorías
- [x] Nivel 2 Manager (401-500 Score): Desbloqueo de Presupuestos
- [x] Nivel 3 Guru (501-600 Score): Desbloqueo de Metas de Ahorro
- [x] Nivel 4 Master (601-900 Score): Desbloqueo de IA Financiera
- [x] Nivel 5 Magnate (901-1000 Score): Acceso completo + features premium futuras
- [x] Veo "widgets bloqueados" con indicación de qué Score necesito
- [x] Los widgets bloqueados muestran beneficios y cuántos puntos faltan para desbloqueo
- [x] Al desbloquear una funcionalidad, recibo celebración y tutorial
- [x] Las funcionalidades están todas desbloqueadas los primeros 35 días desde el registro del usuario
- [x] Veo claramente cómo mis acciones y comportamiento financiero contribuyen a alcanzar el Score necesario

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Financial Health Score
**Estado**:

**Diferenciador**: El sistema de Score impulsa el onboarding progresivo, reduce overwhelm y motiva mejora continua tanto en engagement como en salud financiera real.

---

### Épica 6: Insights y Reportes

#### US-040: Dashboard Principal

**Como** Carlos (Padre Consciente)
**Quiero** ver un resumen de mi situación financiera al abrir la app
**Para** entender rápidamente mi estado actual

**Criterios de Aceptación**:
- [x] Veo balance total (ingresos - gastos)
- [x] Veo total de ingresos y total de gastos
- [x] Veo cantidad de gastos pendientes de pago
- [x] Veo resumen de presupuestos (en meta, con alerta, excedidos)
- [x] Veo resumen de metas de ahorro (total ahorrado, metas activas)
- [x] Veo resumen de transacciones recurrentes (activas, ingresos/gastos mensuales)
- [x] Veo mi nivel de gamificación y XP
- [x] Puedo filtrar por mes y año

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Insights y Reportes
**Estado**: 

---

#### US-041: Gráfico de Gastos por Categoría

**Como** Ana (Ahorrador Aspirante)
**Quiero** ver un gráfico circular de mis gastos por categoría
**Para** identificar dónde gasto más

**Criterios de Aceptación**:
- [x] Veo un pie chart con distribución porcentual de gastos
- [x] Cada categoría tiene color distintivo
- [x] Al hover veo monto exacto y porcentaje
- [x] Veo leyenda con categorías ordenadas por monto (mayor a menor)
- [x] El gráfico se filtra por período seleccionado
- [x] Si no hay datos, veo mensaje educativo

**Prioridad**: Media
**Story Points**: 3
**Épica**: Insights y Reportes
**Estado**: 

---

#### US-042: Métricas Clave

**Como** Laura (Optimizador Financiero)
**Quiero** ver métricas calculadas automáticamente
**Para** entender patrones en mis finanzas

**Criterios de Aceptación**:
- [x] Veo total de transacciones del período
- [x] Veo promedio de gasto diario
- [x] Veo categoría con mayor gasto
- [x] Veo cantidad de gastos pendientes
- [x] Las métricas se actualizan en tiempo real
- [x] Puedo ver métricas para período personalizado (mes, año, custom)

**Prioridad**: Media
**Story Points**: 3
**Épica**: Insights y Reportes
**Estado**: 

---

#### US-043: Filtros Globales de Período

**Como** Laura (Optimizador Financiero)
**Quiero** filtrar todos los datos por mes o año desde un control global
**Para** comparar períodos sin cambiar filtros en cada vista

**Criterios de Aceptación**:
- [x] Hay un selector de período en el header/sidebar
- [x] Puedo seleccionar año específico o "Todos"
- [x] Puedo seleccionar mes específico dentro de año o "Todos"
- [x] Los filtros aplican a: dashboard, gastos, ingresos, presupuestos
- [x] El sistema recuerda mi selección entre sesiones
- [x] Veo indicación clara del período activo
- [x] Puedo resetear filtros fácilmente

**Prioridad**: Alta
**Story Points**: 5
**Épica**: Insights y Reportes
**Estado**: 

**Nota Técnica**: Implementado con React Context (PeriodContext)

---

#### US-044: Desglose del Financial Health Score (Consolidado en US-030)

**Nota**: Esta user story ha sido consolidada en US-030 como parte de la unificación del sistema de Score.

El Financial Health Score ya incluye:
- Score numérico de 1-1000 puntos
- Clasificación en niveles (Aprendiz, Manager, Guru, Master, Magnate)
- Desglose de componentes: Engagement (40%) y Salud Financiera (60%)
- Visualización detallada de cada factor que contribuye al Score
- Actualización en tiempo real basada en acciones y comportamiento financiero

**Ver**: US-030 para detalles completos del Financial Health Score

---

### Épica 7: Gestión de Cuenta

#### US-050: Registro de Usuario

**Como** nuevo usuario
**Quiero** crear una cuenta con email y contraseña
**Para** comenzar a usar la aplicación

**Criterios de Aceptación**:
- [x] Puedo registrarme con email y contraseña
- [x] El sistema valida formato de email
- [x] La contraseña debe cumplir requisitos mínimos de seguridad
- [x] Recibo confirmación de registro exitoso
- [x] Se crea automáticamente mi perfil de gamificación (Nivel 0, 0 Score)
- [x] Soy redirigido al dashboard después del registro
- [x] Tengo todas las funcionalides desbloqueadas por 35 dias *las premium muestran una leyenda con la cuenta regresiva

**Prioridad**: Alta
**Story Points**: 3
**Épica**: Gestión de Cuenta
**Estado**: 

---

#### US-051: Login y Autenticación

**Como** usuario existente
**Quiero** iniciar sesión con mis credenciales
**Para** acceder a mis datos financieros

**Criterios de Aceptación**:
- [x] Puedo hacer login con email y contraseña
- [x] El sistema valida credenciales correctamente
- [x] Si las credenciales son incorrectas, veo mensaje claro
- [x] La sesión persiste si cierro y reabro el navegador
- [x] Puedo cerrar sesión desde el menú
- [x] Al cerrar sesión, soy redirigido al login

**Prioridad**: Alta
**Story Points**: 2
**Épica**: Gestión de Cuenta
**Estado**: 

---

#### US-052: Modo Oscuro

**Como** usuario
**Quiero** alternar entre modo claro y oscuro
**Para** usar la app cómodamente en diferentes condiciones de luz

**Criterios de Aceptación**:
- [x] Hay un toggle para cambiar entre modo claro y oscuro
- [x] El modo se aplica inmediatamente a toda la app
- [x] El modo seleccionado se guarda y persiste
- [x] Todos los componentes se ven correctamente en ambos modos
- [x] Los gráficos y colores se adaptan al modo activo

**Prioridad**: Media
**Story Points**: 3
**Épica**: Gestión de Cuenta
**Estado**: 

---

#### US-053: Ocultar Balances

**Como** Carlos (Padre Consciente)
**Quiero** ocultar temporalmente mis balances y montos
**Para** revisar la app en público sin exponer mi información financiera

**Criterios de Aceptación**:
- [x] Hay un botón para ocultar balances rápidamente
- [x] Todos los montos se muestran como "••••••" cuando están ocultos
- [x] Los porcentajes y gráficos también se ocultan o se muestran de forma genérica
- [x] Puedo volver a mostrar los balances con el mismo botón
- [x] El estado se mantiene durante la sesión pero no persiste

**Prioridad**: Baja
**Story Points**: 2
**Épica**: Gestión de Cuenta
**Estado**: 

---

## Journey Maps

### Journey 1: Primer Registro de Gasto (Ana - Ahorrador Aspirante)

**Contexto**: Ana acaba de registrarse y quiere probar la app por primera vez.

#### Paso 1: Descubrimiento
- **Acción**: Ana abre la app después de registrarse
- **Pensamiento**: "Quiero ver qué puedo hacer aquí"
- **UI**: Dashboard vacío con widget destacado "Registrar Primer Gasto"
- **Emoción**: 😊 Curiosa

#### Paso 2: Exploración
- **Acción**: Hace clic en el botón "Nuevo Gasto"
- **Pensamiento**: "Parece simple, solo necesito completar estos campos"
- **UI**: Modal con formulario limpio: Descripción, Monto, Categoría, Fecha
- **Emoción**: 😌 Tranquila

#### Paso 3: Entrada de Datos
- **Acción**: Completa: "Almuerzo", $1,500, "Comida", hoy, Prioridad Media
- **Pensamiento**: "Me gusta que valide en tiempo real"
- **UI**: Validación inline con iconos ✅, sin errores
- **Emoción**: 😊 Satisfecha

#### Paso 4: Validación
- **Acción**: Hace clic en "Crear"
- **Pensamiento**: "¿Funcionará?"
- **UI**: Spinner breve, luego cierre de modal
- **Emoción**: 😐 Expectante

#### Paso 5: Confirmación
- **Acción**: Ve el gasto en el dashboard inmediatamente
- **Pensamiento**: "¡Está ahí! Y gané puntos 🎉"
- **UI**: Toast notification "¡Gasto creado! +5 Score", gasto visible en lista
- **Emoción**: 😃 Feliz

#### Paso 6: Celebración
- **Acción**: Ve su perfil y un badge "Primera Transacción" desbloqueado, su Score mejoró
- **Pensamiento**: "¡Conseguí mi primer logro y mi Score subió! Esto es divertido"
- **UI**: Badge brillante con animación + notificación de logro + Score mejorado en componente de Adopción de features
- **Emoción**: 😍 Emocionada

#### Paso 7: Retención
- **Acción**: Ve un challenge "Registra 3 transacciones hoy" y su progreso hacia el Nivel 2
- **Pensamiento**: "Voy por 2 más para completar el challenge y mejorar mi Score"
- **UI**: Challenge card con progreso 1/3, barra de Score muestra 185/400 para Nivel 2
- **Emoción**: 💪 Motivada

**Pain Points Identificados**:
- ❌ No hay: No se identifican pain points críticos en este flujo
- ⚠️ Posible confusión: ¿Qué es Score? → Solución: Tooltip educativo en primer login

**Oportunidades de Mejora**:
- 💡 Onboarding interactivo con tooltips en primer uso
- 💡 Sugerencia de categorías populares para nuevos usuarios
- 💡 Celebración visual más notoria al ganar primer XP

**Métricas de Éxito**:
- ✅ Tiempo de completar primer gasto < 2 minutos
- ✅ % usuarios que registran 2+ gastos en primera sesión > 60%
- ✅ % usuarios que vuelven al día siguiente > 40%

---

### Journey 2: Creación de Presupuesto (Carlos - Padre Consciente)

**Contexto**: Carlos quiere controlar sus gastos de supermercado que suelen excederse.

#### Paso 1: Reconocimiento de Necesidad
- **Acción**: Carlos revisa dashboard y ve que gastó $45,000 en "Supermercado" este mes
- **Pensamiento**: "¡Es demasiado! Necesito un límite"
- **UI**: Widget "Presupuestos" muestra "0 presupuestos activos"
- **Emoción**: 😟 Preocupado

#### Paso 2: Descubrimiento de Feature
- **Acción**: Hace clic en el widget de Presupuestos
- **Pensamiento**: "¿Está bloqueado? Necesito mejorar mi Score"
- **UI**: Widget bloqueado muestra "Nivel 2 Manager requerido (401 Score). Tu Score actual: 380. Te faltan 21 puntos"
- **Emoción**: 😐 Neutro pero motivado a desbloquear

#### Paso 3: Progreso hacia Desbloqueo
- **Acción**: Registra gastos pendientes, crea una categoría nueva, completa un challenge
- **Pensamiento**: "Ya casi llego al Score necesario. Veo que mis acciones y mi salud financiera mejoran el Score"
- **UI**: Barra de progreso de Score se llena progresivamente: 380 → 395 → 405
- **Emoción**: 😊 Optimista

#### Paso 4: Desbloqueo
- **Acción**: Alcanza 401 Score, sube a Nivel 2 Manager
- **Pensamiento**: "¡Lo logré! Ahora puedo crear presupuestos"
- **UI**: Animación de level up + notificación "🎉 Nivel 2 Manager Desbloqueado: Presupuestos disponibles"
- **Emoción**: 😃 Satisfecho

#### Paso 5: Creación de Presupuesto
- **Acción**: Crea presupuesto "Supermercado" con límite $35,000/mes
- **Pensamiento**: "Configuré alerta al 80%, me avisará cuando llegue a $28,000"
- **UI**: Formulario intuitivo, alerta_at=80% por defecto
- **Emoción**: 😌 Confiado

#### Paso 6: Monitoreo
- **Acción**: Durante el mes, registra gastos de supermercado
- **Pensamiento**: "Voy en $22,000, me quedan $13,000"
- **UI**: Barra de progreso verde, "63% usado"
- **Emoción**: 👍 Controlado

#### Paso 7: Alerta
- **Acción**: Llega a $28,500 (81% del presupuesto)
- **Pensamiento**: "¡Cuidado! Ya estoy en zona de alerta"
- **UI**: Notificación amarilla "⚠️ Alerta: Supermercado al 81%"
- **Emoción**: 😬 Alerta pero agradecido

#### Paso 8: Ajuste de Comportamiento
- **Acción**: Reduce compras de snacks, cocina más en casa
- **Pensamiento**: "No quiero pasarme, voy a cuidar mis gastos"
- **UI**: Barra se mantiene en amarillo pero no avanza a rojo
- **Emoción**: 💪 Empoderado

**Pain Points Identificados**:
- ⚠️ Frustración inicial: "¿Por qué está bloqueado?"
  - Solución: El widget bloqueado explica claramente el beneficio, el Score necesario y cuántos puntos faltan
- ⚠️ No saber cómo mejorar el Score rápido
  - Solución: Tooltip con sugerencias específicas "Completa un challenge (+10-20 pts al componente Engagement) o mejora tu tasa de ahorro (+puntos al componente Salud Financiera)"

**Oportunidades de Mejora**:
- 💡 Mostrar proyección: "A este ritmo, terminarás con $33,000 (dentro del presupuesto)"
- 💡 Sugerencias de ahorro cuando se activa alerta
- 💡 Celebración si termina el mes sin exceder presupuesto

**Métricas de Éxito**:
- ✅ % usuarios con presupuestos que reducen gastos en categoría > 40%
- ✅ % usuarios que mantienen presupuesto en verde > 50%
- ✅ Adherencia a presupuesto (avg % utilizado) < 90%

---

### Journey 3: Logro de Meta de Ahorro (Ana - Ahorrador Aspirante)

**Contexto**: Ana quiere ahorrar $800,000 para un viaje a Europa en 8 meses.

#### Paso 1: Definición de Objetivo
- **Acción**: Ana crea una meta "Viaje a Europa" con target $800,000 y fecha 8 meses adelante
- **Pensamiento**: "Necesito ahorrar $100,000 por mes"
- **UI**: Sistema calcula automáticamente: "Meta diaria: $3,333, Meta semanal: $23,333, Meta mensual: $100,000"
- **Emoción**: 😌 Claridad

#### Paso 2: Primer Depósito
- **Acción**: Deposita $50,000 que tenía guardados
- **Pensamiento**: "Buen comienzo, voy 6% de la meta"
- **UI**: Barra de progreso al 6%, confetti animation
- **Emoción**: 😊 Motivada

#### Paso 3: Rutina de Ahorro
- **Acción**: Cada semana deposita entre $20,000-$30,000
- **Pensamiento**: "Voy avanzando, a veces más, a veces menos"
- **UI**: Progreso va subiendo gradualmente: 15%, 23%, 31%...
- **Emoción**: 💪 Constante

#### Paso 4: Desmotivación Temporal
- **Acción**: En el mes 4, tuvo gastos inesperados y solo depositó $40,000
- **Pensamiento**: "Estoy atrasada, ¿lo lograré?"
- **UI**: El sistema recalcula: "Te faltan $430,000 en 4 meses → $107,500/mes"
- **Emoción**: 😟 Preocupada

#### Paso 5: Ajuste y Recuperación
- **Acción**: Reduce gastos en delivery y streaming, deposita $120,000 al mes siguiente
- **Pensamiento**: "¡Volví al ritmo! Todavía es posible"
- **UI**: Progreso salta de 46% a 61%
- **Emoción**: 💪 Determinada

#### Paso 6: Sprint Final
- **Acción**: En los últimos 2 meses, hace esfuerzo extra y deposita $130,000/mes
- **Pensamiento**: "Ya casi llego, veo la meta cerca"
- **UI**: 78%, 91%, 98%...
- **Emoción**: 😤 Enfocada

#### Paso 7: Logro de Meta
- **Acción**: Hace depósito final de $15,000, alcanza $800,000
- **Pensamiento**: "¡LO LOGRÉ! 🎉 Me voy a Europa"
- **UI**: Animación de celebración con fuegos artificiales, badge "Meta Completada", +50 pto de Score bonus
- **Emoción**: 😍 Eufórica

#### Paso 8: Compartir
- **Acción**: Toma screenshot del logro y lo comparte en Instagram
- **Pensamiento**: "Quiero que mis amigos vean que sí se puede ahorrar"
- **UI**: Silver Moment> Pantalla de logro compartible con estadísticas
- **Emoción**: 😊 Orgullosa

**Pain Points Identificados**:
- ⚠️ Desmotivación al atrasarse: "¿Abandono la meta?"
  - Solución actual: Recalculo automático de metas ajustadas
  - Mejora futura: Mensaje motivacional "No te rindas, ajusta tu ritmo"
- ⚠️ No poder retirar para emergencias
  - Solución actual: Ya existe funcionalidad de retiro
  - Mejora futura: Separar "emergencias" de "caprichos"

**Oportunidades de Mejora**:
- 💡 Recordatorios semanales: "Esta semana necesitas depositar $23,000"
- 💡 Gamificación de racha: "7 semanas seguidas depositando"
- 💡 Sugerencias automáticas: "Si reduces delivery $10,000, llegarás a tu meta 2 semanas antes"
- 💡 Compartir logros en redes sociales (link generado)

**Métricas de Éxito**:
- ✅ % metas alcanzadas vs. creadas > 40%
- ✅ Tiempo promedio para alcanzar meta < target_date + 10%
- ✅ % usuarios que completan 2+ metas > 25%

---

### Journey 4: Completar Challenge Diario (Ana - Ahorrador Aspirante)

**Contexto**: Ana abre la app por la mañana y ve challenges disponibles.

#### Paso 1: Awareness
- **Acción**: Ve en dashboard "Transaction Master: Registra 3 transacciones hoy (0/3) - Mejora tu Score"
- **Pensamiento**: "Puedo hacer eso fácil, ayer compré 3 cosas"
- **UI**: Challenge card destacada con progreso y contribución al Score
- **Emoción**: 💪 Motivada

#### Paso 2: Primera Transacción
- **Acción**: Registra "Café" $500
- **Pensamiento**: "Una menos"
- **UI**: Challenge card se actualiza en tiempo real: "1/3 completado"
- **Emoción**: 😊 Progreso

#### Paso 3: Segunda Transacción
- **Acción**: Registra "Almuerzo" $2,000
- **Pensamiento**: "Ya voy por la mitad"
- **UI**: "2/3 completado", barra de progreso al 66%
- **Emoción**: 👍 Constante

#### Paso 4: Tercera Transacción
- **Acción**: Registra "Transporte" $800
- **Pensamiento**: "¡Listo, completé el challenge!"
- **UI**: Animación de check ✅, notificación "🎉 Challenge Completado: Transaction Master - Score mejorado"
- **Emoción**: 😃 Satisfecha

#### Paso 5: Recompensa
- **Acción**: Ve su Score aumentar de 425 → 435 (componente de Challenges completados mejoró)
- **Pensamiento**: "Ya estoy más cerca del Nivel 2 Manager"
- **UI**: Barra de Score sube con animación, desglose muestra mejora en componente de Engagement
- **Emoción**: 😊 Motivada a seguir

#### Paso 6: Siguiente Challenge
- **Acción**: Ve otro challenge disponible "Category Organizer: Usa 3 categorías (ya completado)"
- **Pensamiento**: "¡Sin querer también completé este!"
- **UI**: Segundo check automático, Score mejora +8 puntos adicionales en componente de Adopción de features
- **Emoción**: 😍 Delighted

**Pain Points Identificados**:
- ❌ No identificados en este flujo

**Oportunidades de Mejora**:
- 💡 Sugerencias inteligentes: "Tip: Registra tu cena para completar challenge de 3 transacciones"
- 💡 Challenges personalizados basados en comportamiento del usuario
- 💡 Challenges sociales: "Reta a un amigo" (futuro)

**Métricas de Éxito**:
- ✅ % usuarios que completan al menos 1 challenge/día > 50%
- ✅ % challenges diarios completados vs. disponibles > 35%
- ✅ Retención D+1 de usuarios que completan challenges > 70%

---

## Backlog Priorizado

### Must Have (MVP Actual - Completado)

#### Gestión de Transacciones
- ✅ US-001: Registrar Gasto (Alta, 3 pts)
- ✅ US-002: Editar Gasto Inline (Media, 5 pts)
- ✅ US-003: Gestionar Pagos Parciales (Alta, 5 pts)
- ✅ US-004: Registrar Ingreso (Alta, 3 pts)
- ✅ US-005: Gestión de Categorías (Alta, 3 pts)
- ✅ US-006: Filtrado y Ordenamiento (Media, 3 pts)

#### Presupuestos
- ✅ US-010: Crear Presupuesto (Alta, 5 pts)
- ✅ US-011: Alertas de Presupuesto (Alta, 3 pts)
- ✅ US-012: Dashboard de Presupuestos (Media, 3 pts)

#### Metas de Ahorro
- ✅ US-015: Crear Meta de Ahorro (Alta, 5 pts)
- ✅ US-016: Depositar en Meta (Alta, 3 pts)
- ✅ US-017: Retirar de Meta (Media, 3 pts)
- ✅ US-018: Dashboard de Metas (Alta, 3 pts)

#### Transacciones Recurrentes
- ✅ US-020: Crear Transacción Recurrente (Alta, 5 pts)
- ✅ US-021: Pausar/Reanudar (Media, 2 pts)
- ✅ US-022: Ejecutar Vencidas (Alta, 3 pts)
- ✅ US-023: Proyección de Flujo (Media, 5 pts)
- ✅ US-024: Dashboard Recurrentes (Media, 3 pts)

#### Financial Health Score
- ✅ US-030: Financial Health Score y Niveles (Alta, 8 pts)
- ✅ US-031: Logros y Badges (Alta, 5 pts)
- ✅ US-032: Challenges Diarios/Semanales (Media, 5 pts)
- ✅ US-033: Sistema de Rachas (Alta, 3 pts)
- ✅ US-034: Feature Gates por Score (Alta, 5 pts)

#### Insights y Reportes
- ✅ US-040: Dashboard Principal (Alta, 5 pts)
- ✅ US-041: Gráfico por Categoría (Media, 3 pts)
- ✅ US-042: Métricas Clave (Media, 3 pts)
- ✅ US-043: Filtros Globales (Alta, 5 pts)
- ✅ US-044: Desglose del Score (Consolidado en US-030)

#### Gestión de Cuenta
- ✅ US-050: Registro (Alta, 3 pts)
- ✅ US-051: Login (Alta, 2 pts)
- ✅ US-052: Modo Oscuro (Media, 3 pts)
- ✅ US-053: Ocultar Balances (Baja, 2 pts)

---

### Should Have (Próxima Iteración - Q1 2026)

#### Mejoras de UX
- 🔜 US-060: Onboarding Interactivo (Alta, 5 pts)
  - Tutorial guiado para nuevos usuarios
  - Tooltips contextuales en primer uso
  - Checklist de primeros pasos

- 🔜 US-061: Exportar Datos (Media, 3 pts)
  - Exportar transacciones a CSV/Excel
  - Exportar reportes PDF
  - Backup completo de datos

- 🔜 US-062: Búsqueda Global (Media, 5 pts)
  - Buscar transacciones por descripción
  - Buscar por monto, categoría, fecha
  - Resultados instantáneos (typeahead)

#### Notificaciones
- 🔜 US-070: Push Notifications Web (Media, 8 pts)
  - Notificaciones de presupuestos excedidos
  - Recordatorios de metas de ahorro
  - Alertas de transacciones recurrentes pendientes

- 🔜 US-071: Preferencias de Notificaciones (Baja, 3 pts)
  - Configurar qué notificaciones recibir
  - Horarios preferidos para notificaciones
  - Canales (web, email)

#### Gamificación Avanzada
- 🔜 US-075: Leaderboards (Baja, 5 pts)
  - Top usuarios por XP
  - Top rachas
  - Sistema opt-in (privacidad)

- 🔜 US-076: Challenges Personalizados (Media, 8 pts)
  - Challenges basados en comportamiento del usuario
  - Dificultad adaptativa
  - Challenges de fin de semana

#### Insights Avanzados
- 🔜 US-080: Comparación de Períodos (Media, 5 pts)
  - "Este mes vs. mes pasado"
  - Gráficos de tendencias
  - % de variación por categoría

- 🔜 US-081: Análisis de Patrones (Alta, 8 pts)
  - Detección de gastos inusuales (anomalías)
  - Identificación de gastos hormiga
  - Sugerencias de optimización

- 🔜 US-082: Recomendaciones de IA (Alta, 13 pts)
  - "Podrías ahorrar $X reduciendo Y"
  - Predicción de flujo de caja
  - Alertas proactivas

---

### Could Have (Futuro - Q2-Q3 2026)

#### Colaboración
- 📋 US-090: Presupuestos Compartidos (Media, 13 pts)
  - Invitar pareja/familia a presupuesto
  - Permisos: Admin, Editor, Viewer
  - Sincronización en tiempo real

- 📋 US-091: Gastos Compartidos (Media, 8 pts)
  - Dividir gastos con otras personas
  - Tracking de "quién debe a quién"
  - Liquidación de cuentas

#### Integraciones
- 📋 US-095: Importar Transacciones (Alta, 13 pts)
  - Importar CSV de banco
  - Auto-categorización de importados
  - Detección de duplicados

- 📋 US-096: Sincronización Bancaria (Alta, 21 pts)
  - Conexión con bancos argentinos
  - Actualización automática de transacciones
  - Cumplimiento regulatorio y seguridad

#### Mobile
- 📋 US-100: App Móvil iOS (Alta, 34 pts)
  - Versión nativa iOS
  - Paridad de funcionalidades con web
  - Widgets de iOS

- 📋 US-101: App Móvil Android (Alta, 34 pts)
  - Versión nativa Android
  - Paridad de funcionalidades con web
  - Widgets de Android

#### Avanzado
- 📋 US-110: Inversiones (Baja, 21 pts)
  - Tracking de inversiones (acciones, bonos, crypto)
  - Gráficos de rentabilidad
  - Integración con brokers

- 📋 US-111: Planeación Fiscal (Baja, 13 pts)
  - Cálculo de impuestos (autónomos)
  - Exportar para contador
  - Deducciones sugeridas

---

### Won't Have (Out of Scope para MVP)

#### Fuera de Alcance Temporal
- ❌ US-120: Multi-moneda
  - Soporte para USD, EUR, etc.
  - Conversión automática
  - **Razón**: Complejidad vs. necesidad (usuarios argentinos = ARS)

- ❌ US-121: Tarjetas de Crédito Avanzadas
  - Tracking de límites de crédito
  - Fechas de cierre y vencimiento
  - Intereses y financiación
  - **Razón**: Scope muy grande, requiere dominio completo de TC

- ❌ US-122: Préstamos
  - Tracking de préstamos personales
  - Amortización e intereses
  - **Razón**: Nicho muy específico, no core

- ❌ US-123: White Label
  - Versión personalizable para empresas
  - Multi-tenancy corporativo
  - **Razón**: Fuera del mercado objetivo (B2C, no B2B)

---

## Métricas de Éxito

### North Star Metric
**Weekly Active Users (WAU)** - Usuarios que usan la app al menos 1 vez por semana

**Objetivo**:
- Mes 1: 5 WAU (50% de 10 usuarios beta)
- Mes 3: 8 WAU (80% de 10 usuarios beta)
- Mes 6: 100 WAU (con crecimiento orgánico)

---

### Métricas de Engagement

#### Daily Active Users (DAU)
- **Objetivo**: DAU/WAU > 0.4 (40% de usuarios semanales vuelven cada día)
- **Medición**: Usuarios únicos con al menos 1 acción por día

#### Duración de Sesión
- **Objetivo**: > 5 minutos promedio
- **Medición**: Tiempo entre login y última acción
- **Benchmark**: Apps de finanzas personales ~3-7 min

#### Retención de Racha
- **Objetivo**: 30% de usuarios mantienen racha de 7 días
- **Medición**: % usuarios con current_streak >= 7

#### Completion Rate de Challenges
- **Objetivo**: 40% de challenges diarios completados
- **Medición**: (Challenges completados / Challenges disponibles) * 100

---

### Métricas de Salud Financiera

#### Savings Rate
- **Objetivo**: Promedio de ahorro = 15% de ingresos
- **Medición**: (Total ahorrado en metas / Total ingresos) * 100
- **Nota**: Solo usuarios con metas de ahorro activas

#### Budget Adherence
- **Objetivo**: 60% de presupuestos terminan el período sin excederse
- **Medición**: % de presupuestos con status != 'exceeded' al final del período

#### Goals Achievement Rate
- **Objetivo**: 40% de metas se completan antes de target_date
- **Medición**: Metas con status='achieved' / Total metas creadas

---

### Métricas de Uso del Producto

#### Transacciones Registradas
- **Objetivo**: 15 transacciones/usuario/mes
- **Medición**: AVG(transacciones por usuario en último mes)
- **Segmentado por**: Tipo (ingreso/gasto), Categoría

#### Categorías Creadas
- **Objetivo**: 5-10 categorías/usuario (sweet spot)
- **Medición**: COUNT(categorías) WHERE user_id
- **Insight**: <3 = sub-utilización, >15 = complejidad excesiva

#### Insights Vistos
- **Objetivo**: 2 vistas de insights/semana/usuario
- **Medición**: user_gamification.insights_viewed / semanas activas
- **Tipos**: Dashboard, AI Health Score, Reportes

#### Challenges Completados
- **Objetivo**: 10 challenges/usuario/mes
- **Medición**: SUM(challenges completados) / usuarios activos

---

### Métricas de Retención

#### Day 1 Retention
- **Objetivo**: > 50%
- **Medición**: % usuarios que vuelven 1 día después del registro

#### Day 7 Retention
- **Objetivo**: > 60%
- **Medición**: % usuarios que vuelven 7 días después del registro

#### Day 30 Retention
- **Objetivo**: > 40%
- **Medición**: % usuarios que vuelven 30 días después del registro

#### Churn Rate
- **Objetivo**: < 5% mensual
- **Medición**: Usuarios sin actividad en 30 días / Total usuarios
- **Definición de churn**: Sin login en 30 días

---

### Métricas de Gamificación

#### Score Distribution
- **Objetivo**: Distribución balanceada entre componentes del Score
- **Medición**: % contribución por componente
- **Ideal**:
  - Engagement Component: 35-45% (target 40%)
  - Health Component: 55-65% (target 60%)
- **Sub-componentes Engagement**:
  - Racha de días: 20-30 pts
  - Adopción de features: 20-30 pts
  - Consumo de insights: 15-25 pts
  - Challenges completados: 15-25 pts

#### Level Progression
- **Objetivo**: Usuarios alcanzan Nivel 3 Guru (501 Score) en primer mes
- **Medición**: AVG(current_level) de usuarios con 30 días de antigüedad
- **Benchmark**: Nivel 3 = 501 Score = Engagement activo + salud financiera moderada

#### Achievement Completion
- **Objetivo**: 50% de logros disponibles completados por usuario activo
- **Medición**: AVG(achievements completados / total achievements) por usuario

---

### Métricas de Calidad

#### Error Rate
- **Objetivo**: < 1% de requests con error
- **Medición**: 5xx errors / total requests

#### Page Load Time
- **Objetivo**: < 2 segundos (P95)
- **Medición**: Time to Interactive (TTI)

#### Data Accuracy
- **Objetivo**: 100% de balances consistentes
- **Medición**: (Ingresos - Gastos) = Balance en DB
- **Validación**: Daily automated check

---

## Apéndice: Priorización de Backlog

### Framework de Priorización (RICE)

**Reach**: ¿Cuántos usuarios impacta?
**Impact**: ¿Qué tan grande es el impacto? (Alta=3, Media=2, Baja=1)
**Confidence**: ¿Qué tan seguros estamos? (Alta=100%, Media=80%, Baja=50%)
**Effort**: ¿Cuántos story points requiere?

**RICE Score** = (Reach * Impact * Confidence) / Effort

### Ejemplos de Scoring

| US | Reach | Impact | Confidence | Effort | RICE |
|----|-------|--------|------------|--------|------|
| US-060: Onboarding | 10 | 3 | 100% | 5 | 6.0 |
| US-070: Push Notifications | 10 | 2 | 80% | 8 | 2.0 |
| US-080: Comparación Períodos | 8 | 2 | 100% | 5 | 3.2 |
| US-090: Presupuestos Compartidos | 4 | 3 | 50% | 13 | 0.46 |
| US-100: App iOS | 5 | 3 | 80% | 34 | 0.35 |

**Conclusión**: Priorizar US-060 (Onboarding) > US-080 (Comparación) > US-070 (Push) para Q1 2026.

---

## Changelog

- **2026-02-09**: Versión inicial basada en MVP actual
- **Futuro**: Actualizaciones trimestrales después de cada sprint

---

**Fin del Documento**
