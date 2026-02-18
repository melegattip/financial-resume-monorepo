# Reporte de Bugs - Financial Resume
**Fecha:** 2026-02-05  
**Total de issues:** 16

---

## 🎮 Gamificación

### BUG-001: Categorización y racha de días no funcionan
**Severidad:** Alta  
**Módulo:** Gamification  
**Captura:** ![Bug 1](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/1.jpeg)

**Descripción:**  
Las funcionalidades de categorización de transacciones y el sistema de racha de días continúan sin funcionar correctamente.

**Observaciones:**
- "Expert en Organización": progreso 0/50 (0% completado)
- "Guerrero Semanal": progreso 0/7 (0% completado)  
- "Leyenda Mensual": progreso 0/30 (0% completado)

**Acción requerida:**
- Verificar lógica de tracking de categorización
- Revisar sistema de conteo de días consecutivos

---

## 👤 Gestión de Usuarios / Planes

### BUG-002: Banner de trial visible en usuario nivel 7
**Severidad:** Media  
**Módulo:** User Subscription / Presupuestos y Metas  
**Captura:** ![Bug 2](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/2.jpeg)

**Descripción:**  
Usuario con nivel 7 (Financial Guru) sigue viendo el banner de "Acceso de prueba a Presupuestos" en las secciones de Presupuestos y Metas de Ahorro.

**Comportamiento esperado:**  
El banner de trial no debería mostrarse para usuarios que ya tienen acceso completo.

**Acción requerida:**
- Revisar lógica de validación de plan/nivel de usuario
- Ajustar condiciones de visualización del banner

---

## 📅 Transacciones Recurrentes

### BUG-003: Falta opción de ejecución a futuro en transacciones recurrentes
**Severidad:** Media  
**Módulo:** Transacciones Recurrentes  
**Tipo:** Feature Request + Bug

**Descripción:**  
Se necesita incluir en Transacciones recurrentes la opción de poder ejecutar a futuro, que el usuario pueda elegir ya tener disponibles los resúmenes de meses próximos.

**Problema adicional:**  
El filtro global de mes está fijado por defecto en el último mes del último año, pero debería estar en el **mes actual del año en curso**.

**Acción requerida:**
- Agregar opción para generar transacciones recurrentes hacia meses futuros
- Modificar filtro de mes por defecto: de "último mes del último año" a "mes actual del año en curso"

---

### BUG-008: Vencimientos incorrectos en recurrentes creadas en meses anteriores
**Severidad:** Alta  
**Módulo:** Transacciones Recurrentes  
**Captura:** ![Bug 8](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/8.jpeg)

**Descripción:**  
Los vencimientos de las transacciones recurrentes creadas en meses anteriores quedan mal configurados.

**Comportamiento esperado:**
- Dejar en `nil` los vencimientos que NO estableció el usuario
- No mostrar en el frontend vencimientos que no fueron establecidos manualmente

**Acción requerida:**
- Revisar lógica de cálculo de vencimientos en recurrentes
- Ocultar vencimientos auto-generados que no fueron definidos por el usuario

---

### BUG-016: Gastos recurrentes de octubre no reconocidos
**Severidad:** Alta  
**Módulo:** Transacciones Recurrentes  
**Captura:** ![Bug 16](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/16.jpeg)

**Descripción:**  
No se reconocen los gastos recurrentes ejecutados en el mes de octubre.

**Acción requerida:**
- Investigar lógica de ejecución de recurrentes para el mes de octubre
- Verificar si hay problemas con el cambio de mes o año fiscal

---

## 💰 Metas de Ahorro

### BUG-004: Falta soporte para metas de ahorro en USD
**Severidad:** Media  
**Módulo:** Metas de Ahorro  
**Captura:** ![Bug 4](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/4.jpeg)

**Descripción:**  
Actualmente solo se pueden crear metas de ahorro en la moneda local. Se necesita agregar soporte para metas en USD.

**Acción requerida:**
- Implementar selector de moneda en creación de metas
- Adaptar cálculos y visualizaciones para múltiples monedas

---

## 📅 Visualización de Fechas

### BUG-006: Mostrar nombre del día en fechas de vencimiento
**Severidad:** Baja  
**Módulo:** Gastos / Transacciones / Visualización  
**Tipo:** Enhancement

**Descripción:**  
Mostrar las fechas de vencimiento de transacciones con el nombre del día para mejor legibilidad.

**Ejemplo:**  
En lugar de: `31/8/2025`  
Mostrar: `Domingo, 31/8/2025` o `31/8/2025 (Dom)`

**Acción requerida:**
- Agregar formato de fecha que incluya día de la semana
- Aplicar en todos los lugares donde se muestren vencimientos

---

## 🤖 IA Financiera

### BUG-007: IA no distingue entre gastos e inversiones/ahorros
**Severidad:** Alta  
**Módulo:** IA Financiera / "Puedo comprarlo"  

**Descripción:**  
La IA Financiera no está interpretando correctamente la naturaleza de ciertas transacciones. Hay "gastos" que no son realmente gastos, sino inversiones/ahorros.

**Problemas identificados:**

1. **Evaluación de salud financiera incorrecta:**
   - La IA solo mira el "balance" (que es un resto)
   - No diferencia egresos de inversión/ahorro vs gastos reales
   - El balance debería tender a 0 siempre, ya que es dinero sin destino

2. **Sección "Puedo comprarlo" con información imprecisa:**
   - **Balance:** no significa mucho por sí solo
   - **Gastos mensuales:** incluye ahorros e inversiones incorrectamente
   - **Disponible por mes:** es lo mismo que balance

**Comportamiento esperado:**
- La IA debe tener criterio para identificar qué categorías o gastos representan inversiones/ahorros
- No tomarlos como simples salidas de dinero sino como creación de capital
- A la hora de evaluar salud financiera, considerar estos activos creados

**Acción requerida:**
- Implementar categorización de egresos: Gasto vs Inversión/Ahorro
- Recopilar los activos creados por egresos de categorías "Ahorro" o "Inversiones"
- Mejorar análisis de "Puedo comprarlo" basándose en activos reales
- Ajustar recomendaciones de IA considerando esta distinción

---

## 📊 Dashboard y Visualización

### BUG-005: Meses anteriores a veces no aparecen
**Severidad:** Media  
**Módulo:** Dashboard / Selector de Período  
**Captura:** ![Bug 5](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/5.jpeg)

**Descripción:**  
En ocasiones, los meses anteriores no aparecen en el selector de período o en la visualización de datos históricos.

**Acción requerida:**
- Revisar lógica de carga de períodos históricos
- Verificar filtros y queries de base de datos

---

### BUG-008: Filtro de "filas" quita items en lugar de plegarlos
**Severidad:** Baja  
**Módulo:** Dashboard / Tabla de Transacciones  
**Captura:** ![Bug 10](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/10.jpeg)

**Descripción:**  
El filtro de "filas" está quitando items del listado en lugar de plegarlos y agregar una barra para scroll.

**Comportamiento esperado:**
- Los items filtrados deberían mantenerse en el DOM pero colapsados
- Mostrar scrollbar cuando haya overflow

**Acción requerida:**
- Modificar comportamiento del filtro: de ocultar a colapsar
- Implementar scroll vertical adecuado

---

### BUG-009: Mejorar UX de añadir transacciones rápidamente
**Severidad:** Media  
**Módulo:** Dashboard / Transacciones  
**Captura:** ![Bug 15](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/15.jpeg)

**Descripción:**  
Falta un método rápido para añadir transacciones sin abrir modal.

**Requerimientos:**
- Agregar un último renglón en gastos e ingresos
- Permitir añadir sin modal ni confirmación
- Crear las transacciones en el mes seleccionado actualmente (ej: septiembre)
- **Ubicación:** esta sección debe estar **arriba** de toda la pantalla, widgets de resumen debajo

**Mejora de UX adicional:**
- Los márgenes entre transacciones deben ser más chicos
- Optimizar espacio vertical para que entren más transacciones en pantalla
- **No** achicar letra ni espacios entre columnas, solo modificar espacios en blanco en eje vertical

**Acción requerida:**
- Implementar fila de "quick add" arriba del dashboard
- Reordenar layout: transacciones arriba, widgets abajo
- Ajustar espaciados verticales

---

## 💸 Gestión de Gastos

### BUG-010: Gastos necesitan campo de prioridad
**Severidad:** Media  
**Módulo:** Gastos  
**Tipo:** Feature Request

**Descripción:**  
Los gastos deben tener un nuevo campo: **prioridad**.

**Comportamiento esperado:**
- Campo numérico de prioridad
- Ordenamiento por defecto: ascendente (prioridad 0 para arriba)
- Los más importantes primero

**Acción requerida:**
- Agregar campo `prioridad` al modelo de Gasto
- Implementar ordenamiento por defecto
- Agregar UI para editar prioridad

---

### BUG-011: Comportamiento confuso del icono "pagar"
**Severidad:** Media  
**Módulo:** Gastos / Pagos  
**Captura:** ![Bug 9](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/9.jpeg)

**Descripción:**  
El icono verde "pagar" hace exactamente lo mismo que la cruz que indica pendiente de pago.

**Cambio solicitado:**
- El **tilde verde** debe ser directamente "pagar total"
- En versión de escritorio, mostrar en el renglón el **monto ya pagado** en caso de pagos parciales

**Acción requerida:**
- Modificar acción del botón verde a "pagar total"
- Agregar visualización de pagos parciales en desktop

---

### BUG-012: Fecha de vencimiento muestra un día menos
**Severidad:** Media  
**Módulo:** Gastos  
**Captura:** ![Bug 11](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/11.jpeg)

**Descripción:**  
El vencimiento de gasto no coincide con lo que muestra el ítem: aparece un día menos.

**Ejemplo:**
- Fecha establecida: 11/09/2025
- Fecha mostrada en el ítem: 10/09/2025

**Acción requerida:**
- Revisar conversión de fechas y timezone
- Verificar si hay problema con almacenamiento UTC vs local

---

### BUG-013: Edición de gasto desde dashboard es lenta
**Severidad:** Alta  
**Módulo:** Gastos / Dashboard  
**Tipo:** UX Enhancement

**Descripción:**  
Editar gasto desde el dashboard queda muy lerdo al intentar actualizar.

**Cambio solicitado:**
- Evitar el modal de edición
- La fila debe funcionar como Excel: al pasar el cursor y hacer clic, se edita
- Al modificar una transacción, mostrar un simple **check** para confirmar los cambios

**Acción requerida:**
- Implementar edición inline tipo Excel
- Simplificar confirmación de cambios
- Optimizar performance de actualización

---

### BUG-014: Pago parcial no se refleja correctamente
**Severidad:** Alta  
**Módulo:** Gastos / Pagos  
**Captura:** ![Bug 13](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/13.jpeg)

**Descripción:**  
Se registró un pago parcial, pero:
- No se descontó de "gastos pendientes"
- No se muestra ningún indicador en el ítem del gasto

**Acción requerida:**
- Verificar lógica de actualización de pagos parciales
- Agregar indicador visual de pago parcial en el ítem
- Asegurar que se descuente correctamente de pendientes

---

### BUG-015: No se puede crear gasto en mes anterior
**Severidad:** Alta  
**Módulo:** Gastos / Creación  
**Captura:** ![Bug 14](file:///c:/Users/Los Meles/Documents/Development/Financial Resume/financial-resume-monorepo/bugfixing_ctx/14.jpeg)

**Descripción:**  
Al crear un gasto, no se puede crearlo en un mes anterior al actual. No hay campo para seleccionar la fecha.

**Comportamiento esperado:**
- Debe existir un campo de fecha de transacción
- Permitir crear gastos con fecha retroactiva

**Acción requerida:**
- Agregar campo de fecha en formulario de creación
- Validar que se permita seleccionar meses anteriores

---

## 📌 Resumen por Prioridad

### 🔴 Alta Prioridad (6 bugs)
- BUG-001: Gamificación no funciona
- BUG-003: Vencimientos incorrectos en recurrentes
- BUG-007: IA no distingue gastos vs inversiones
- BUG-013: Edición de gasto lenta
- BUG-014: Pago parcial no se refleja
- BUG-015: No se puede crear gasto en mes anterior
- BUG-016: Recurrentes de octubre no reconocidos

### 🟡 Media Prioridad (8 bugs)
- BUG-002: Banner trial en usuario nivel 7
- BUG-004: Ejecución a futuro en recurrentes
- BUG-005: Metas de ahorro en USD
- BUG-009: Meses anteriores no aparecen
- BUG-010: Campo de prioridad en gastos
- BUG-011: Comportamiento confuso de icono pagar
- BUG-012: Vencimiento muestra un día menos
- BUG-015: UX de añadir transacciones

### 🟢 Baja Prioridad (2 bugs)
- BUG-006: Mostrar nombre del día en vencimientos
- BUG-008: Filtro de filas quita items

---

## 📊 Módulos Afectados

| Módulo | Cantidad de Bugs |
|--------|------------------|
| Gastos | 6 |
| Transacciones Recurrentes | 3 |
| Dashboard / Visualización | 3 |
| IA Financiera | 1 |
| Gamificación | 1 |
| Metas de Ahorro | 1 |
| Usuarios / Planes | 1 |

---

**Próximos pasos:**
1. Priorizar bugs de severidad alta
2. Asignar responsables
3. Estimar tiempos de resolución
4. Crear issues en sistema de tracking
