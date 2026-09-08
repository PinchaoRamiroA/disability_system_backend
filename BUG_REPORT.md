# Reporte de Bugs y Defectos Encontrados en la Aplicación

**Fecha:** 8 de Septiembre, 2026  
**Sistema:** Backend Sistema de Gestión de Incapacidades (`disability_system_backend`)  
**Lenguaje / Framework:** Go / Gin / GORM / PostgreSQL  

---

## Resumen Ejecutivo

Durante el análisis estático y de arquitectura del código fuente backend, se identificaron **15 defectos significativos** clasificados en 5 categorías:
1. **Errores Críticos y Riesgos de Panic en Tiempo de Ejecución**
2. **Errores de Lógica y de Cálculo Financiero / Plazos**
3. **Inconsistencias en Rutas, Parámetros y Respuestas HTTP**
4. **Violaciones de Arquitectura Hexagonal y Guías de Desarrollo (`AGENTS.md`)**
5. **Defectos en Consultas SQL / GORM y Base de Datos**

---

## 1. Errores Críticos y Riesgos de Panic (Runtime Panics)

### 🔴 BUG-01: Posible `Panic` por Desreferencia `nil` en `DocumentosService.VerificarEstadoTransicion`
* **Ubicación:** [`internal/modules/incapacidades/usecase/documentos_service.go:125, 150`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/incapacidades/usecase/documentos_service.go#L124-L159)
* **Descripción:** En la función `VerificarEstadoTransicion`, si `incapacidad.Estado` es `nil` (por ejemplo, si la entidad se cargó sin precargar la relación `Estado`), la instrucción `validTransitions[incapacidad.Estado.Nombre]` provocará un **panic por puntero nulo**.
* **Impacto:** Caída inesperada del hilo de ejecución en peticiones HTTP si la estructura no tiene precargado el estado.

### 🔴 BUG-02: Posible `Panic` por Puntero Nulo en `AuditoriaUseCase.Crear`
* **Ubicación:** [`internal/modules/auditoria/usecase/auditoria_usecase.go:26-28`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/auditoria/usecase/auditoria_usecase.go#L26-L28)
* **Descripción:** En Go, una interfaz como `ports.Actor` con un valor concreto `nil` (ej. `(*ActorImpl)(nil)`) evalúa `actor != nil` como `true`. Al llamar a `actor.GetUserID()` se produce un `panic` por desreferencia de puntero nulo.
* **Impacto:** Fallos del servidor si se invoca la auditoría con una interfaz `Actor` no inicializada.

### 🔴 BUG-03: Posible `Panic` por `historialSvc` no validado en `DocumentoUseCase`
* **Ubicación:** [`internal/modules/incapacidades/usecase/documento_usecase.go:63, 98`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/incapacidades/usecase/documento_usecase.go#L63)
* **Descripción:** `DocumentoUseCase.Subir` y `Validar` ejecutan `uc.historialSvc.CreateEntry(...)` directamente sin comprobar si `uc.historialSvc` es `nil`.
* **Impacto:** Si la dependencia no es inyectada correctamente, la subida o validación de documentos crashéa con panic.

---

## 2. Errores de Lógica y Cálculo

### 🟠 BUG-04: Cálculo Erróneo de `IncapacidadesActivas` en Estadísticas de Cartera
* **Ubicación:** [`internal/modules/cobros/usecase/cartera_usecase.go:107`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/cobros/usecase/cartera_usecase.go#L107)
* **Descripción:** En `ObtenerEstadisticasGenerales`:
  ```go
  IncapacidadesActivas: total - int64(len(pagos)),
  ```
  `total` proviene de la consulta a la tabla `pago`, no de la tabla `incapacidad`. Restar la longitud del slice de pagos (`len(pagos)`) del total de pagos no calcula incapacidades activas; resulta en un valor matemáticamente incorrecto (o 0 si `total == len(pagos)`).
* **Impacto:** Los dashboards y reportes gerenciales muestran métricas de cartera distorsionadas.

### 🟠 BUG-05: Pérdida de Precisión por Conversión Reiterada `float64` $\leftrightarrow$ `string` en Cartera
* **Ubicación:** [`internal/modules/cobros/usecase/cartera_usecase.go:293-302`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/cobros/usecase/cartera_usecase.go#L293-L302)
* **Descripción:** La función `sumCurrency` convierte strings formateados a `float64`, realiza sumas flotantes y los vuelve a formatear redondeando a 2 decimales en cada iteración de un bucle.
* **Impacto:** Acumulación de errores de redondeo de punto flotante en sumatorias monetarias grandes, violando el estándar de precisión financiera que la librería `decimal` pretendía proveer.

### 🟠 BUG-06: Rama Imposible de Alcanzar en Alertas de Vencimiento de Cartera
* **Ubicación:** [`internal/modules/cobros/usecase/cartera_usecase.go:185-197`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/cobros/usecase/cartera_usecase.go#L185-L197)
* **Descripción:** La consulta filtra pagos con `pago.FechaPago.Before(fechaLimite)` donde `fechaLimite` es una fecha en el pasado (`time.Now().AddDate(0,0,-diasMinimos)`). Posteriormente evalúa `if diasVencido > 0`, lo cual siempre se cumple. La rama `else { mensaje = "Vence en X días" }` es código muerto inalcanzable.
* **Impacto:** Incoherencia en los mensajes de alerta de vencimiento cuando se consultan plazos futuros o límites negativos.

### 🟠 BUG-07: Transición de Estados de Incapacidad No Validada (Código Muerto)
* **Ubicación:** [`internal/modules/incapacidades/usecase/documentos_service.go:124-159`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/incapacidades/usecase/documentos_service.go#L124) y [`incapacidad_usecase.go:297`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/incapacidades/usecase/incapacidad_usecase.go#L297)
* **Descripción:** Se implementó `VerificarEstadoTransicion` con un mapa de transiciones válidas en `DocumentosService`, pero la función `CambiarEstado` en `IncapacidadUseCase` **nunca la invoca**.
* **Impacto:** Se permiten transiciones de estado ilegales en el dominio (por ejemplo, pasar de `Pagada` directamente a `Recibida`).

---

## 3. Inconsistencias de API, Rutas y Middleware

### 🟡 BUG-08: Parámetro de Ruta Ignorado en `GET /incapacidades/:id/documentos`
* **Ubicación:** [`internal/modules/incapacidades/adapters/http/documento_handler.go:196`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/incapacidades/adapters/http/documento_handler.go#L196)
* **Descripción:** El endpoint registrado es `/incapacidades/:id/documentos`. Sin embargo, en el handler `Listar`, no se lee `c.Param("id")`. En su lugar se exige un query param `?id_incapacidad=X`. Si no se envía el query param, responde `400 BAD REQUEST ("id_incapacidad es requerido")`, ignorando la ID de la URL.
* **Impacto:** Incompatibilidad con clientes REST que consumen la ruta RESTful `/incapacidades/123/documentos`.

### 🟡 BUG-09: Inconsistencia en la Estructura de Respuesta JSON de Errores en `JWTMiddleware`
* **Ubicación:** [`internal/modules/auth/adapters/http/jwt_middleware.go:24-49`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/auth/adapters/http/jwt_middleware.go#L24-L49)
* **Descripción:** El middleware de autenticación devuelve respuestas de error directamente con la estructura `{ "success": false, "code": "...", "message": "..." }`, mientras que el estándar global de la aplicación (definido en [`internal/shared/response/response.go`](file:///c:/dev/proyectos/disability_system_backend/internal/shared/response/response.go)) anida el código dentro de un objeto `error: { "code": "..." }`.
* **Impacto:** Dificultad para los clientes frontend al parsear errores de autenticación vs. errores generales de la API.

### 🟡 BUG-10: Expiración de Token Devuelta en `0` al Renovar Token (`Refresh`)
* **Ubicación:** [`internal/modules/auth/adapters/http/auth_handler.go:119`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/auth/adapters/http/auth_handler.go#L119)
* **Descripción:** Al solicitar `/auth/refresh`, la respuesta envía `expires_in: 0`:
  ```go
  resp := mapper.ToTokenResponse(tokens.AccessToken, tokens.RefreshToken, 0)
  ```
* **Impacto:** Los clientes de la API no pueden determinar cuándo caduca el nuevo token renovado.

---

## 4. Defectos de Arquitectura Hexagonal y Estándares (`AGENTS.md`)

### 🔵 BUG-11: Bypass del Adaptador de Persistencia / Uso de SQL Raw en UseCase de Usuarios
* **Ubicación:** [`internal/modules/usuarios/usecase/usuario_usecase.go:177, 198`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/usuarios/usecase/usuario_usecase.go#L177)
* **Descripción:** En `CambiarEstado` y `CambiarPassword`, el UseCase ejecuta SQL directo (`uc.db.Exec("UPDATE usuario SET...")`) accediendo a GORM directamente.
* **Violación:** Regla explícita en `AGENTS.md`: *"persistence belongs in adapters/postgres"*, *"All database access must stay here"*, *"Avoid raw SQL unless necessary"*.
* **Efecto secundario:** No actualiza la columna `updated_at` de la tabla `usuario`.

### 🔵 BUG-12: Hardcoding de ID de Rol en Registro de Usuarios
* **Ubicación:** [`internal/modules/auth/usecase/register_usecase.go:54`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/auth/usecase/register_usecase.go#L54)
* **Descripción:** `RegisterUseCase` asigna estáticamente `IDRol: 4` al registrar un usuario.
* **Impacto:** Si la base de datos se inicializa con otros identificadores de rol o una secuencia distinta, la creación de usuarios falla con error de clave foránea (`foreign key constraint failure`).

### 🔵 BUG-13: Uso de `log.Printf` Estándar en Lugar de Zap Logger
* **Ubicación:** [`internal/shared/storage/upload.go:23`](file:///c:/dev/proyectos/disability_system_backend/internal/shared/storage/upload.go#L23)
* **Descripción:** Se utiliza `log.Printf("R2 Upload Error: %v", err)`.
* **Violación:** Regla explícita en `AGENTS.md`: *"Use Zap logger. Do not use fmt.Println / standard log for application logging"*.

---

## 5. Defectos en Consultas SQL / GORM

### 🟣 BUG-14: Consulta SQL Inválida `IN ()` en `CobroRepository.GetIncapacidadesDetailed`
* **Ubicación:** [`internal/modules/cobros/adapters/postgres/cobro_repository.go:318-323`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/cobros/adapters/postgres/cobro_repository.go#L318-L323)
* **Descripción:** Si la lista de IDs recibida por `GetIncapacidadesDetailed` está vacía (`len(ids) == 0`), se ejecuta la consulta `Where("id_incapacidad IN ?", ids)`. En PostgreSQL/GORM esto genera la sintaxis `WHERE id_incapacidad IN (NULL)` o falla la consulta.
* **Impacto:** Error de sintaxis SQL al listar alertas cuando no hay pagos registrados.

### 🟣 BUG-15: Extracción de Extensión de Archivos Proensa a Errores en `DocumentoHandler`
* **Ubicación:** [`internal/modules/incapacidades/adapters/http/documento_handler.go:379-388`](file:///c:/dev/proyectos/disability_system_backend/internal/modules/incapacidades/adapters/http/documento_handler.go#L379-L388)
* **Descripción:** La función `getExtensionFromFilename` realiza recortes manuales de cadenas asumiendo longitudes fijas de 4 o 5 caracteres en lugar de usar `filepath.Ext(filename)`. Archivos con nombres de 4 caracteres o extensiones inusuales producen resultados erróneos.
* **Impacto:** Registro defectuoso del campo `formato` en la base de datos.

---

## Tabla Resumen de Severidad

| ID | Módulo | Severidad | Breve Descripción |
|---|---|---|---|
| BUG-01 | Incapacidades | **Alta (Critical)** | Panic por puntero nil en `VerificarEstadoTransicion`. |
| BUG-02 | Auditoría | **Alta (Critical)** | Panic por interfaz `Actor` nil en `AuditoriaUseCase`. |
| BUG-03 | Incapacidades | **Alta (Critical)** | Panic si `historialSvc` es nil al subir/validar documentos. |
| BUG-04 | Cobros | **Media (High)** | Cálculo erróneo de `IncapacidadesActivas` en Cartera. |
| BUG-05 | Cobros | **Media (High)** | Pérdida de precisión monetaria por conversiones repetidas a `float64`. |
| BUG-06 | Cobros | **Baja (Medium)** | Mensaje de alerta de vencimiento inalcanzable (código muerto). |
| BUG-07 | Incapacidades | **Media (High)** | Validación de transición de estados desconectada (código muerto). |
| BUG-08 | Incapacidades | **Media (High)** | Query param obligatorio duplica parámetro de ruta en GET documentos. |
| BUG-09 | Auth / Shared | **Baja (Medium)** | Formato JSON de error inconsistente en JWT Middleware. |
| BUG-10 | Auth | **Baja (Low)** | Refresh Token devuelve `expires_in: 0`. |
| BUG-11 | Usuarios | **Media (High)** | SQL Raw en UseCase violando arquitectura e ignorando `updated_at`. |
| BUG-12 | Auth | **Media (Medium)** | Hardcoding de `IDRol: 4` al registrar usuarios. |
| BUG-13 | Storage | **Baja (Low)** | Uso de `log.Printf` en lugar de Zap logger. |
| BUG-14 | Cobros | **Media (High)** | Consulta SQL `IN ()` sin validar slice vacío en `GetIncapacidadesDetailed`. |
| BUG-15 | Incapacidades | **Baja (Low)** | Extracción manual e incorrecta de extensiones de archivos. |

---

## Recomendaciones de Remediación

1. **Guardas de Punteros Nulos (Nil Checks):** Añadir validaciones de puntero/interfaz en UseCases y Handlers antes de invocar métodos sobre objetos o interfaces inyectadas.
2. **Validación de Transiciones de Estado:** Conectar `VerificarEstadoTransicion` dentro de `IncapacidadUseCase.CambiarEstado`.
3. **Refactorización de Cartera y Decimales:** Utilizar únicamente operaciones del paquete `shopspring/decimal` para sumas monetarias, evitando `float64`.
4. **Adherencia Arquitectónica:** Mover las consultas SQL directas de `UsuarioUseCase` al repositorio `UsuarioRepository` e integrar la actualización de la fecha `updated_at`.
5. **Estandarización de Respuestas HTTP:** Sustituir la generación manual de mapas JSON en middleware por los helper functions de `internal/shared/response`.
