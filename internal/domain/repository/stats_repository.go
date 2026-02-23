package repository

import "context"

// GlobalStats holds the raw global statistics from the database
type GlobalStats struct {
	TotalUsers             int
	TotalActiveUsers       int
	TotalSchools           int
	TotalSubjects          int
	TotalGuardianRelations int
}

// StatsRepository defines persistence operations for statistics
type StatsRepository interface {
	GetGlobalStats(ctx context.Context) (*GlobalStats, error)
}
