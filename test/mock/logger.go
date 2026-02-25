package mock

import "github.com/EduGoGroup/edugo-shared/logger"

// MockLogger is a no-op logger for tests
type MockLogger struct{}

func NewMockLogger() *MockLogger { return &MockLogger{} }

func (l *MockLogger) Debug(_ string, _ ...interface{})    {}
func (l *MockLogger) Info(_ string, _ ...interface{})     {}
func (l *MockLogger) Warn(_ string, _ ...interface{})     {}
func (l *MockLogger) Error(_ string, _ ...interface{})    {}
func (l *MockLogger) Fatal(_ string, _ ...interface{})    {}
func (l *MockLogger) With(_ ...interface{}) logger.Logger { return l }
func (l *MockLogger) Sync() error                         { return nil }
