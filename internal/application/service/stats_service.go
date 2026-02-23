package service

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
)

// StatsService defines the stats service interface
type StatsService interface {
	GetGlobalStats(ctx context.Context) (*dto.GlobalStatsResponse, error)
}

type statsService struct {
	statsRepo repository.StatsRepository
	logger    logger.Logger
}

// NewStatsService creates a new stats service
func NewStatsService(statsRepo repository.StatsRepository, logger logger.Logger) StatsService {
	return &statsService{statsRepo: statsRepo, logger: logger}
}

func (s *statsService) GetGlobalStats(ctx context.Context) (*dto.GlobalStatsResponse, error) {
	stats, err := s.statsRepo.GetGlobalStats(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError("get global stats", err)
	}

	s.logger.Info("global stats retrieved")

	return &dto.GlobalStatsResponse{
		TotalUsers:             stats.TotalUsers,
		TotalActiveUsers:       stats.TotalActiveUsers,
		TotalSchools:           stats.TotalSchools,
		TotalSubjects:          stats.TotalSubjects,
		TotalGuardianRelations: stats.TotalGuardianRelations,
	}, nil
}
