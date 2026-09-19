# Reporte de Vulnerabilidades de Seguridad y Protección de Datos (Control de Acceso / BOLA)

**Fecha:** 19 de Septiembre, 2026  
**Sistema:** Backend Sistema de Gestión de Incapacidades (`disability_system_backend`)  
**Marco Legal y Normativo:** Ley Estatutaria 1581 de 2012 (Habeas Data), Ley 23 de 1981 y Resolución 1995 de 1999 (Reserva Legal de Historia Clínica), Estándar OWASP Top 10 API Security (Broken Object Level Authorization - BOLA / IDOR).

---

## Resumen Ejecutivo

Durante la auditoría de seguridad y control de acceso basada en roles (RBAC) y atributos (ABAC) en el backend, se identificaron **6 vulnerabilidades activas** de severidad crítica, alta y media, que permiten a roles no privilegiados (como `Empleado` o usuarios con tokens válidos) acceder, descargar o alterar información médica, personal y financiera de otros colaboradores o de la organización.

Adicionalmente, se documentan al final **2 vulnerabilidades críticas ya remediadas** durante la sesión actual (radicación para terceros y exposición del directorio corporativo de usuarios).

### Matriz de Vulnerabilidades Encontradas

| Código | Vulnerabilidad | Severidad | Módulo Afectado | Endpoint(s) |
| :--- | :--- | :--- | :--- | :--- |
| **BUG-SEC-01** | Fuga Masiva y Descarga de Documentos Médicos Privados (BOLA/IDOR) | 🔴 **Crítica** | `incapacidades/documentos` | `GET /api/v1/incapacidades/:id/documentos` |
| **BUG-SEC-02** | Fuga de Información Financiera y Cartera a Empleados Rasos | 🔴 **Alta** | `cobros/cartera` | `GET /api/v1/cartera/*` |
| **BUG-SEC-03** | Carga No Autorizada de Archivos en Incapacidades Ajenas (BOLA/IDOR) | 🔴 **Alta** | `incapacidades/documentos` | `POST /api/v1/incapacidades/:id/documentos` |
| **BUG-SEC-04** | Exposición de Trazabilidad y Notas Internas Médicas (BOLA/IDOR) | 🟡 **Media / Alta** | `historial` | `GET /api/v1/incapacidades/:id/historial` |
| **BUG-SEC-05** | Fuga de Plazos y Términos Legales de Incapacidades Ajenas | 🟡 **Media** | `incapacidades` | `GET /api/v1/incapacidades/:id/plazos` |
| **BUG-SEC-06** | Consulta Global y Transcripción No Autorizada ante EPS/ARL | 🟡 **Media / Alta** | `transcripciones` | `GET /api/v1/incapacidades/transcripciones/pendientes`<br>`POST /api/v1/incapacidades/:id/transcribir` |

---

## Detalle de Defectos y Vulnerabilidades

### 🔴 BUG-SEC-01: Fuga Masiva y Descarga de Documentos Médicos Privados (BOLA / IDOR)

* **Severidad:** **CRÍTICA** (CWE-639: Authorization Bypass Through User-Controlled Key)
* **Endpoint:** `GET /api/v1/incapacidades/:id/documentos`
* **Ubicación:** [`internal/modules/incapacidades/usecase/documento_usecase.go:129-134`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/usecase/documento_usecase.go#L129-L134)
* **Descripción:**
  El método `Listar` de documentos solo evalúa:
  ```go
  if !actor.HasPermission("consultar_incapacidad") {
      return nil, 0, apperrors.ErrForbidden.WithMessage("no tienes permiso para consultar documentos")
  }
  return uc.repo.List(ctx, incapacidadID, estado, tipo, page, limit)
  ```
  Como todos los usuarios (incluidos los que tienen rol `Empleado`) poseen el permiso `consultar_incapacidad` para ver sus propios registros, cualquier usuario autenticado puede alterar el parámetro `:id` en la URL para consultar los documentos adjuntos de cualquier incapacidad de la empresa.
* **Impacto:** Fuga masiva de documentos altamente sensibles:
  * Historias clínicas completas y epicrisis.
  * Certificados médicos con diagnósticos y códigos CIE-10 reservados.
  * Reportes de presuntos accidentes laborales (FURAT).
  * Documentos de identidad (cédulas de ciudadanía).
  * URLs públicas o presignadas de descarga directa desde el bucket Cloudflare R2 / S3.
* **Solución Propuesta:**
  Consultar la entidad incapacidad y validar que pertenezca al usuario autenticado o que el actor tenga rol de gestión (`actor.CanManageIncapacidades()`):
  ```go
  incapacidad, err := uc.incapacidadRepo.FindByID(ctx, incapacidadID)
  if err != nil {
      return nil, 0, err
  }
  if incapacidad.IDUsuario != actor.UserID && !actor.CanManageIncapacidades() {
      return nil, 0, apperrors.ErrForbidden.WithMessage("solo puedes consultar documentos de tus propias incapacidades")
  }
  ```

---

### 🔴 BUG-SEC-02: Fuga de Información Financiera y Cartera a Empleados Rasos

* **Severidad:** **ALTA** (CWE-285: Improper Authorization)
* **Endpoints:**
  * `GET /api/v1/cartera/estadisticas`
  * `GET /api/v1/cartera/resumen-entidad`
  * `GET /api/v1/cartera/vencida`
  * `GET /api/v1/cartera/alertas-vencimiento`
* **Ubicación:** [`internal/modules/cobros/usecase/cobro_usecase.go:322-328`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/cobros/usecase/cobro_usecase.go#L322-L328) y [`internal/modules/cobros/usecase/cartera_usecase.go:374-385`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/cobros/usecase/cartera_usecase.go#L374-L385)
* **Descripción:**
  La función que autoriza la lectura del módulo de cartera (`canReadCobros`) incluye el permiso `consultar_incapacidad`:
  ```go
  func canReadCobros(actor ports.Actor) bool {
      return actor.HasPermission("consultar_incapacidad") || // <-- VULNERABILIDAD
          actor.HasPermission("registrar_pago") ||
          actor.HasPermission("realizar_conciliacion") ||
          actor.HasPermission("gestionar_cobro_persuasivo") ||
          actor.HasPermission("gestionar_cobro_juridico")
  }
  ```
* **Impacto:** Los empleados comunes tienen acceso irrestricto a:
  * Cifras financieras globales de la compañía (dinero total recobrado vs pendiente por EPS y ARL).
  * Reportes de deudas institucionales de entidades de salud.
  * Incapacidades en cobro jurídico o persuasivo y cartera vencida clasificada por nivel de riesgo.
* **Solución Propuesta:**
  Eliminar `actor.HasPermission("consultar_incapacidad")` de `canReadCobros`. El acceso a estadísticas y cartera debe reservarse exclusivamente a roles con permisos de contabilidad, tesorería o gestión (`registrar_pago`, `realizar_conciliacion`, `gestionar_cobro_persuasivo`, `gestionar_cobro_juridico`, `consultar_reportes`).

---

### 🔴 BUG-SEC-03: Carga No Autorizada de Archivos en Incapacidades Ajenas (BOLA / IDOR)

* **Severidad:** **ALTA** (CWE-639 / CWE-434)
* **Endpoint:** `POST /api/v1/incapacidades/:id/documentos`
* **Ubicación:** [`internal/modules/incapacidades/adapters/http/documento_handler.go:88`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/adapters/http/documento_handler.go#L88) y [`usecase/documento_usecase.go:55-72`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/usecase/documento_usecase.go#L55-L72)
* **Descripción:**
  1. En el Handler, el archivo binario se envía y almacena físicamente en el storage de Cloudflare R2 **antes** de verificar si el actor tiene derecho a interactuar con esa incapacidad.
  2. En el UseCase, solo se valida que el actor tenga `crear_incapacidad` o `editar_incapacidad` y que la incapacidad exista (`ExistsIncapacidad`), pero **no se valida que la incapacidad le pertenezca a `actor.UserID`**.
* **Impacto:** Un empleado malintencionado puede contaminar expedientes ajenos subiendo documentos falsos o archivos no autorizados en incapacidades de otros trabajadores.
* **Solución Propuesta:**
  Verificar la propiedad de la incapacidad antes de invocar la carga en storage o en el caso de uso:
  ```go
  if incapacidad.IDUsuario != actor.UserID && !actor.CanManageIncapacidades() {
      return nil, apperrors.ErrForbidden.WithMessage("no tienes autorización para subir documentos a incapacidades de otros colaboradores")
  }
  ```

---

### 🟡 BUG-SEC-04: Exposición de Trazabilidad y Notas Internas Médicas (BOLA / IDOR)

* **Severidad:** **MEDIA / ALTA** (CWE-639)
* **Endpoint:** `GET /api/v1/incapacidades/:id/historial`
* **Ubicación:** [`internal/modules/incapacidades/adapters/http/documento_handler.go:285-328`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/adapters/http/documento_handler.go#L285-L328)
* **Descripción:**
  La función `ListarHistorial` invoca directamente la consulta `h.historialListFn(incapacidadID, ...)` sin verificar la identidad del actor solicitante, sin validar permisos y sin comprobar si la incapacidad le pertenece.
* **Impacto:** Cualquier usuario autenticado puede conocer la bitácora de auditoría médica:
  * Notas de rechazo o validación emitidas por analistas de Gestión Humana.
  * Cambios de estado y fechas de radicación de cualquier expediente de la compañía.
  * Identificadores de los gestores que auditaron el caso.
* **Solución Propuesta:**
  Cargar el actor en el handler y validar que `actor.CanManageIncapacidades()` sea verdadero o que la incapacidad pertenezca al usuario en sesión.

---

### 🟡 BUG-SEC-05: Fuga de Plazos y Términos Legales de Incapacidades Ajenas

* **Severidad:** **MEDIA** (CWE-285)
* **Endpoint:** `GET /api/v1/incapacidades/:id/plazos`
* **Ubicación:** [`internal/modules/incapacidades/usecase/incapacidad_usecase.go:396-405`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/usecase/incapacidad_usecase.go#L396-L405)
* **Descripción:**
  A diferencia del endpoint `GET /incapacidades/:id` (que ejecuta `ensureCanRead`), la función `ObtenerInfoPlazos` omitió invocar la verificación de acceso.
* **Impacto:** Un empleado puede monitorear las fechas límites de radicación, plazos de transcripción ante EPS y alertas de vencimiento de cualquier colaborador de la organización.
* **Solución Propuesta:**
  Invocar `ensureCanRead(actor, incapacidad)` inmediatamente tras cargar la incapacidad por su ID.

---

### 🟡 BUG-SEC-06: Consulta Global y Transcripción No Autorizada ante EPS/ARL

* **Severidad:** **MEDIA / ALTA** (CWE-285 / CWE-862)
* **Endpoints:**
  * `GET /api/v1/incapacidades/transcripciones/pendientes`
  * `POST /api/v1/incapacidades/:id/transcribir`
* **Ubicación:** [`internal/modules/incapacidades/adapters/http/transcripcion_handler.go:37-60, 133-155`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/adapters/http/transcripcion_handler.go#L37-L60) y [`transcripcion_usecase.go:26-30`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/usecase/transcripcion_usecase.go#L26-L30)
* **Descripción:**
  1. `ListarPendientes` no extrae actor ni verifica permisos; devuelve todas las incapacidades pendientes de transcripción de la empresa.
  2. `Transcribir` no valida si el actor cuenta con el rol o permiso para transcribir incapacidades ante la EPS/ARL.
* **Impacto:**
  * Visibilidad no autorizada de incapacidades ajenas pendientes de trámite.
  * Riesgo de que un empleado marque administrativamente como "completada/transcrita" su propia incapacidad sin que la EPS haya expedido el reconocimiento.
* **Solución Propuesta:**
  Exigir `actor.CanManageIncapacidades()` o restringir ambas rutas con `jwtMiddleware.RequireRole("Administrador", "Gestión Humana", "SG-SST")`.

---

## Registro de Vulnerabilidades Ya Remediadas (Sesión Actual)

Para trazabilidad del proyecto, se registran las vulnerabilidades críticas corregidas previamente:

### ✅ BUG-SEC-00A: Radicación de Incapacidades a Nombre de Terceros (CORREGIDO)
* **Descripción:** En `IncapacidadUseCase.Crear`, cualquier empleado con permiso `crear_incapacidad` podía enviar un `id_usuario` arbitrario y registrar la incapacidad a nombre de otro trabajador.
* **Solución Aplicada:**
  Se agregó la validación estricta en [`internal/modules/incapacidades/usecase/incapacidad_usecase.go`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/incapacidades/usecase/incapacidad_usecase.go#L81-L83):
  ```go
  if input.IDUsuario != actor.UserID && !actor.CanManageIncapacidades() {
      return nil, apperrors.ErrForbidden.WithMessage("solo los administradores o gestores autorizados pueden registrar incapacidades para otros colaboradores")
  }
  ```
  Se añadió la prueba unitaria `TestCrearIncapacidad_EmpleadoNoPuedeCrearParaTercero`.

### ✅ BUG-SEC-00B: Fuga del Directorio Corporativo de Usuarios a Roles sin Privilegios (CORREGIDO)
* **Descripción:** `GET /api/v1/usuarios` estaba abierto a cualquier usuario con JWT, exponiendo nombres, cédulas, correos y números celulares de toda la empresa.
* **Solución Aplicada:**
  Se protegió el endpoint en [`internal/modules/usuarios/adapters/http/register.go`](file:///home/computer/Documentos/proyectos/disability_system_back/internal/modules/usuarios/adapters/http/register.go#L27) con:
  ```go
  jwtMiddleware.RequireRole("Administrador", "admin", "Gestión Humana", "SG-SST", "Recepcionista")
  ```
  Y en `GET /usuarios/:id` se restringió para que un empleado regular solo pueda ver su propio ID de usuario (`uid == id`).
