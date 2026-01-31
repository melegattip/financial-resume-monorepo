# Plan de Desacoplamiento de Gestión de Usuarios

## Objetivo
Desacoplar toda la gestión de usuarios del core financiero, creando un microservicio de usuarios independiente, y definir los contratos de integración necesarios para cumplir con los requerimientos del frontend y la arquitectura de microservicios.

---

## 1. Alcance y Motivación
- Separar la lógica de usuarios (registro, login, perfil, preferencias, notificaciones, seguridad) del core financiero.
- Permitir escalabilidad, mantenibilidad y evolución independiente de la gestión de usuarios.
- Cumplir con los requerimientos de UI/UX y seguridad mostrados en el frontend.

---

## 2. Arquitectura Propuesta
- **Nuevo Microservicio:** `user-service` (gestión de usuarios, preferencias, notificaciones, seguridad)
- **API Core Financiero:** Solo referencia `user_id` y valida JWT emitido por `user-service`.
- **API Gamificación:** Obtiene datos de usuario vía API o eventos, nunca accede directamente a datos sensibles.

---

## 3. Contratos de Endpoints del Microservicio de Usuarios

### 3.1. Autenticación y Gestión de Sesión

#### POST `/users/register`
**Request:**
```json
{
  "email": "string",
  "password": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string"
}
```
**Response:**
```json
{
  "token": "string",
  "expires_at": 1712345678,
  "user": { "id": 1, "email": "string", "first_name": "string", "last_name": "string", "phone": "string" }
}
```

#### POST `/users/login`
**Request:**
```json
{
  "email": "string",
  "password": "string"
}
```
**Response:** igual a `/users/register`

#### POST `/users/logout`
**Request:**
```json
{
  "token": "string"
}
```
**Response:**
```json
{ "success": true }
```

#### POST `/users/refresh`
**Request:**
```json
{ "refresh_token": "string" }
```
**Response:** igual a `/users/register`

---

### 3.2. Perfil de Usuario

#### GET `/users/profile`
**Headers:** `Authorization: Bearer <token>`
**Response:**
```json
{
  "id": 1,
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "created_at": "2024-06-01T12:00:00Z"
}
```

#### PUT `/users/profile`
**Request:**
```json
{
  "first_name": "string",
  "last_name": "string",
  "phone": "string"
}
```
**Response:** igual a GET `/users/profile`

---

### 3.3. Preferencias de Usuario

#### GET `/users/preferences`
**Response:**
```json
{
  "currency": "ARS",
  "language": "es",
  "theme": "dark|light",
  "date_format": "DD/MM/YYYY"
}
```

#### PUT `/users/preferences`
**Request:**
```json
{
  "currency": "ARS",
  "language": "es",
  "theme": "dark",
  "date_format": "DD/MM/YYYY"
}
```
**Response:** igual a GET `/users/preferences`

---

### 3.4. Notificaciones

#### GET `/users/notifications/settings`
**Response:**
```json
{
  "email_notifications": true,
  "push_notifications": false,
  "weekly_reports": true,
  "expense_alerts": true
}
```

#### PUT `/users/notifications/settings`
**Request:**
```json
{
  "email_notifications": true,
  "push_notifications": false,
  "weekly_reports": true,
  "expense_alerts": true
}
```
**Response:** igual a GET `/users/notifications/settings`

---

### 3.5. Seguridad y Privacidad

#### PUT `/users/security/change-password`
**Request:**
```json
{
  "current_password": "string",
  "new_password": "string"
}
```
**Response:**
```json
{ "success": true }
```

#### POST `/users/security/2fa/setup`
**Response:**
```json
{ "qr_code": "data:image/png;base64,...", "secret": "string" }
```

#### POST `/users/security/2fa/verify`
**Request:**
```json
{ "code": "string" }
```
**Response:**
```json
{ "success": true }
```

#### POST `/users/export`
**Response:**
```json
{ "download_url": "https://.../export/1234.zip" }
```

#### DELETE `/users`
**Request:**
```json
{ "password": "string" }
```
**Response:**
```json
{ "success": true }
```

---

## 4. Cambios en el Core Financiero (API principal)
- Eliminar toda lógica de autenticación, perfil y preferencias de usuario.
- Validar JWT y extraer `user_id` usando la clave pública del `user-service`.
- Para datos de usuario extendidos (nombre, preferencias), consultar el `user-service` vía REST o gRPC.
- Mantener solo la referencia a `user_id` en las entidades financieras.

---

## 5. Cambios en la API de Gamificación
- Eliminar dependencias directas de datos de usuario.
- Obtener datos de usuario (nombre, email, etc) solo vía API del `user-service`.
- Para logros y estadísticas, usar solo el `user_id`.

---

## 6. Consideraciones de Integración
- Sincronizar tokens y sesiones entre servicios usando JWT estándar.
- Definir eventos o webhooks para cambios críticos (eliminación de usuario, cambio de email, etc).
- Documentar los contratos OpenAPI/Swagger para todos los endpoints nuevos y modificados.

---

## 7. Roadmap de Implementación
1. Crear el microservicio `user-service` con los endpoints y modelos definidos.
2. Migrar datos de usuario del core financiero al nuevo servicio.
3. Refactorizar el core financiero para eliminar lógica de usuario y consumir el nuevo servicio.
4. Refactorizar la API de gamificación para consumir el nuevo servicio.
5. Actualizar el frontend para consumir los endpoints desacoplados.
6. Pruebas de integración y validación de seguridad.
7. Documentar y capacitar al equipo sobre la nueva arquitectura.

---

## 8. Apéndice: Ejemplo de JWT
```json
{
  "sub": 1,
  "email": "user@example.com",
  "exp": 1712345678,
  "scope": "user"
}
```

---

**Este plan permite desacoplar la gestión de usuarios, mejorar la seguridad y escalabilidad, y cumplir con los requerimientos del frontend y la arquitectura moderna.** 