# Informe de Endpoints No Migrados

**Fecha:** 2026-02-22
**Comparacion:** APIs antiguas vs APIs nuevas (admin + mobile)
**Ultima actualizacion:** 2026-02-22 - Migracion completada

---

## Resumen Ejecutivo

| API | Endpoints Antiguos | Endpoints Nuevos | Migrados en esta sesion | Pendientes |
|-----|-------------------|-----------------|------------------------|------------|
| Admin | 71 activos + 9 fantasma | **74** | 8 nuevos | 1 (verify-bulk, evaluar) |
| Mobile | 23 | 24 | 0 (ya estaban completos) | 0 |

---

## ESTADO: Endpoints Migrados (completados)

### Auth - Migrados

| Metodo | Ruta Nueva | Estado |
|--------|-----------|--------|
| `POST` | `/api/v1/auth/switch-context` | MIGRADO - Cambia contexto de escuela, genera nuevos JWT |
| `GET` | `/api/v1/auth/contexts` | MIGRADO - Retorna todos los contextos disponibles del usuario |

### Users CRUD - Migrados (antes eran handlers fantasma)

| Metodo | Ruta Nueva | Permiso | Estado |
|--------|-----------|---------|--------|
| `POST` | `/api/v1/users` | `PermissionUsersUpdate` | MIGRADO - Crear usuario con email/password |
| `GET` | `/api/v1/users` | `PermissionUsersRead` | MIGRADO - Listar usuarios con filtros (?is_active, ?limit, ?offset) |
| `GET` | `/api/v1/users/:user_id` | `PermissionUsersRead` | MIGRADO - Obtener usuario por ID |
| `PATCH` | `/api/v1/users/:user_id` | `PermissionUsersUpdate` | MIGRADO - Actualizar first_name, last_name, is_active |
| `DELETE` | `/api/v1/users/:user_id` | `PermissionUsersUpdate` | MIGRADO - Soft delete usuario |

### Stats - Migrado (antes era handler fantasma)

| Metodo | Ruta Nueva | Permiso | Estado |
|--------|-----------|---------|--------|
| `GET` | `/api/v1/stats/global` | `PermissionPermissionsMgmtRead` | MIGRADO - Stats globales del sistema |

**Respuesta:**
```json
{
  "total_users": 150,
  "total_active_users": 142,
  "total_schools": 5,
  "total_subjects": 23,
  "total_guardian_relations": 89
}
```

### Materials Admin - Migrado (antes era handler fantasma)

| Metodo | Ruta Nueva | Permiso | Estado |
|--------|-----------|---------|--------|
| `DELETE` | `/api/v1/materials/:id` | `PermissionPermissionsMgmtUpdate` | MIGRADO - Soft delete material (moderacion admin) |

---

## PENDIENTE: Evaluar migracion

### `POST /v1/auth/verify-bulk`

**Archivo original:** `edugo-api-administracion/internal/auth/handler/verify_handler.go:131`

**Que hacia:** Verificar hasta 100 tokens JWT en una sola llamada (S2S con API key).

**Recomendacion:** Verificar en logs de produccion si recibe trafico. Si no, no migrar.

---

## NO MIGRADOS: Handlers redundantes o deprecados

Estos handlers existian en la API antigua pero son redundantes con funcionalidad ya existente:

| Handler | Razon de no migracion |
|---------|----------------------|
| `UnitHandler.CreateUnit` | Redundante: `AcademicUnitHandler` ya tiene `POST /schools/:id/units` |
| `UnitHandler.UpdateUnit` | Redundante: `AcademicUnitHandler` ya tiene `PUT /units/:id` |
| `UnitHandler.AssignMember` | Deprecado: retornaba 501 NOT_IMPLEMENTED, reemplazado por `MembershipHandler` |

---

## Archivos creados/modificados en la migracion

### Archivos nuevos (auth context)
- `internal/auth/handler/auth_handler.go` - Metodos SwitchContext + GetAvailableContexts
- `internal/auth/service/auth_service.go` - Logica de negocio SwitchContext + GetAvailableContexts
- `internal/auth/dto/auth_dto.go` - DTOs SwitchContextRequest/Response, AvailableContextsResponse

### Archivos nuevos (User CRUD)
- `internal/application/dto/user_dto.go` - CreateUserRequest, UpdateUserRequest, UserResponse
- `internal/application/service/user_service.go` - UserService interface + implementacion
- `internal/infrastructure/http/handler/user_handler.go` - 5 endpoints CRUD

### Archivos nuevos (Stats)
- `internal/application/dto/stats_dto.go` - GlobalStatsResponse
- `internal/domain/repository/stats_repository.go` - StatsRepository interface + GlobalStats
- `internal/infrastructure/persistence/postgres/repository/stats_repository.go` - Implementacion GORM
- `internal/application/service/stats_service.go` - StatsService
- `internal/infrastructure/http/handler/stats_handler.go` - Handler

### Archivos nuevos (Material Delete)
- `internal/domain/repository/material_repository.go` - MaterialRepository interface
- `internal/infrastructure/persistence/postgres/repository/material_repository.go` - Implementacion GORM
- `internal/application/service/material_service.go` - MaterialService
- `internal/infrastructure/http/handler/material_handler.go` - Handler

### Archivos modificados
- `internal/domain/repository/user_repository.go` - Agregados Create, ExistsByEmail, Delete
- `internal/infrastructure/persistence/postgres/repository/user_repository.go` - Implementaciones GORM
- `internal/container/container.go` - Wiring de nuevos servicios y handlers
- `cmd/main.go` - Registro de nuevas rutas
- `test/mock/services.go` - Mocks actualizados

---

## Notas sobre Diferencias de Prefijo

| API | Prefijo Antiguo | Prefijo Nuevo |
|-----|----------------|---------------|
| Admin | `/v1/` | `/api/v1/` |
| Mobile | `/v1/` | `/api/v1/` |
