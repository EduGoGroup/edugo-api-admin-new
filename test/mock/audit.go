package mock

import "github.com/EduGoGroup/edugo-shared/audit"

// NewNoopAuditLogger returns a no-op audit logger for tests
func NewNoopAuditLogger() audit.AuditLogger {
	return audit.NewNoopAuditLogger()
}
