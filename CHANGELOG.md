# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-02-24


## [0.9.0] - 2026-03-07

### Tipo de Release: patch

- update (#18)

---

## [0.8.0] - 2026-03-06

### Tipo de Release: patch

- update (#16)
- feat(admin-api): authorization guards + audit logging (#13) (#14)
- swagger
- Add Golang Pro skill with concurrency, generics, and testing guides

---

## [0.7.0] - 2026-03-04

### Tipo de Release: patch

- chore: bump repository to v0.3.3 (ApplyPagination + int64 totals)
- chore(deps): bump shared/common to v0.52.0
- fix(pagination): address code review comments from PR #11
- feat(pagination): implement real pagination with COUNT for admin endpoints
- perf(docker): eliminate Go compilation from Docker, reduce image time ~80%

---

## [0.6.0] - 2026-03-03

### Tipo de Release: patch

- chore(deps): bump edugo-shared/auth to v0.52.0 and repository to v0.3.2
- chore(deps): bump edugo-infrastructure/postgres to v0.58.0

---

## [0.5.0] - 2026-03-02

### Tipo de Release: patch

- fix: validar createdBy como UUID e incluir tests de cobertura
- Update subject API docs and dependencies
- chore: actualizar edugo-infrastructure/postgres a v0.57.0
- fix: corregir tipo CreatedBy en GuardianRelation DTO y service

---

## [0.4.0] - 2026-03-02

### Tipo de Release: patch

- fix: update subject service tests for school_id parameter
- chore: bump edugo-infrastructure/postgres to v0.55.0
- feat: add school_id filtering to subjects entity (multi-tenancy fix)

---

## [0.3.0] - 2026-02-26

### Tipo de Release: patch

- Add search and filter support to list endpoints and update dependencies

---

## [0.2.0] - 2026-02-25

### Tipo de Release: patch

- Update dependencies and clean up struct formatting and imports
- fix: use GITHUB_TOKEN instead of GHCR_TOKEN for registry auth

---

## [0.1.1] - 2026-02-24

### Tipo de Release: patch



---

## [0.1.0] - 2026-02-24

### Tipo de Release: patch

- Change exposed port from 8081 to 8080 in Dockerfile
- chore: release v0.1.0
- Update Azure deploy workflow env vars and secrets
- refactor: consolida IAM en iam-platform y agrega clientes HTTP para auth y roles
- feat: Extiende los mocks del repositorio de usuarios, actualiza las pruebas del servicio de autenticación para incluir dependencias de membresía y escuela, mejora la descripción de seguridad en la documentación y añade un archivo de espacio de trabajo de VS Code.
- Refactor repositories to use GORM for database operations
- feat: initial commit - clean API Admin rebuild

---
### Tipo de Release: patch

- Update Azure deploy workflow env vars and secrets
- refactor: consolida IAM en iam-platform y agrega clientes HTTP para auth y roles
- feat: Extiende los mocks del repositorio de usuarios, actualiza las pruebas del servicio de autenticación para incluir dependencias de membresía y escuela, mejora la descripción de seguridad en la documentación y añade un archivo de espacio de trabajo de VS Code.
- Refactor repositories to use GORM for database operations
- feat: initial commit - clean API Admin rebuild

---
