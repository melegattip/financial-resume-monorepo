# Financial Resume Engine - Frontend

Frontend moderno para la aplicación Financial Resume Engine, desarrollado con React y un diseño visual profesional y elegante.

## 🚀 Características

- **Diseño moderno**: Interface limpia y profesional
- **Responsive**: Optimizado para dispositivos móviles y desktop
- **Dashboard interactivo**: Métricas en tiempo real y gráficos dinámicos
- **Gestión completa**: CRUD para gastos, ingresos y categorías
- **Reportes avanzados**: Análisis financiero con visualizaciones
- **UX optimizada**: Navegación intuitiva y feedback visual

## 🛠️ Tecnologías

- **React 18**: Framework principal
- **Tailwind CSS**: Estilos y diseño
- **Recharts**: Gráficos y visualizaciones
- **React Router**: Navegación
- **Axios**: Cliente HTTP
- **React Hook Form**: Manejo de formularios
- **React Hot Toast**: Notificaciones
- **Lucide React**: Iconografía

## 📦 Instalación

```bash
# Instalar dependencias
npm install

# Iniciar servidor de desarrollo
npm start

# Construir para producción
npm run build
```

## 🎨 Estructura del Proyecto

```
src/
├── components/
│   └── Layout/
│       ├── Sidebar.jsx      # Navegación lateral
│       └── Header.jsx       # Cabecera de la aplicación
├── pages/
│   ├── Dashboard.jsx        # Página principal con métricas
│   ├── Expenses.jsx         # Gestión de gastos
│   ├── Incomes.jsx          # Gestión de ingresos
│   ├── Categories.jsx       # Gestión de categorías
│   ├── Reports.jsx          # Reportes y análisis
│   └── Settings.jsx         # Configuración de usuario
├── services/
│   └── api.js              # Cliente API y servicios
├── App.jsx                 # Componente principal
├── index.js               # Punto de entrada
└── index.css              # Estilos globales
```

## 🎯 Funcionalidades

### Dashboard
- Métricas principales (balance, ingresos, gastos)
- Gráficos de tendencias mensuales
- Distribución de gastos por categoría
- Transacciones recientes
- Toggle para ocultar/mostrar montos

### Gestión de Gastos
- Lista completa de gastos con filtros
- Crear, editar y eliminar gastos
- Marcar como pagado/pendiente
- Asociar categorías
- Fechas de vencimiento
- Cálculo automático de porcentajes

### Gestión de Ingresos
- Lista de ingresos con búsqueda
- CRUD completo de ingresos
- Categorización de ingresos
- Actualización automática de porcentajes de gastos

### Categorías
- Gestión visual de categorías
- Crear y editar categorías con descripciones
- Vista en tarjetas organizadas

### Reportes
- Filtros por rango de fechas
- Gráficos de tendencias y distribución
- Tabla detallada por categorías
- Métricas consolidadas
- Exportación de datos

### Configuración
- Perfil de usuario
- Preferencias de notificaciones
- Configuración de idioma y moneda
- Opciones de seguridad
- Exportación de datos

## 🎨 Diseño

El frontend utiliza un sistema de diseño profesional con:

- **Colores**: Paleta moderna y profesional con azul primario, verde y naranja de acento
- **Tipografía**: Inter como fuente principal
- **Componentes**: Cards, botones, inputs y navegación con diseño limpio
- **Animaciones**: Transiciones suaves y feedback visual
- **Iconografía**: Lucide React para iconos consistentes

## 🔧 Configuración

### Variables de Entorno

Crear un archivo `.env` en la raíz del proyecto:

```env
REACT_APP_API_URL=http://localhost:8080/api/v1
```

### Proxy de Desarrollo

El `package.json` incluye un proxy para desarrollo que redirige las llamadas API al backend:

```json
"proxy": "http://localhost:8080"
```

## 📱 Responsive Design

La aplicación está optimizada para:

- **Desktop**: Layout completo con sidebar fijo
- **Tablet**: Sidebar colapsable y grid adaptativo
- **Mobile**: Navegación móvil y componentes apilados

## 🚀 Despliegue

```bash
# Construir para producción
npm run build

# Los archivos estáticos se generan en la carpeta 'build'
```

## 🔗 Integración con Backend

El frontend se conecta con la API de Financial Resume Engine a través de:

- **Autenticación**: Header `x-caller-id` para identificación de usuario
- **Endpoints**: CRUD completo para todas las entidades
- **Manejo de errores**: Interceptores de Axios con notificaciones
- **Formato de datos**: Consistente con las respuestas del backend

## 📊 Características Avanzadas

- **Cálculo automático de porcentajes**: Los gastos se actualizan automáticamente cuando cambian los ingresos
- **Filtros inteligentes**: Búsqueda y filtrado en tiempo real
- **Validación de formularios**: Validación client-side con feedback visual
- **Estados de carga**: Spinners y estados de loading para mejor UX
- **Notificaciones**: Toast notifications para todas las acciones

## 🎯 Próximas Funcionalidades

- [ ] Modo oscuro
- [ ] PWA (Progressive Web App)
- [ ] Notificaciones push
- [ ] Exportación a PDF/Excel
- [ ] Gráficos más avanzados
- [ ] Filtros de fecha más granulares
- [ ] Búsqueda global
- [ ] Atajos de teclado

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.