package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"gorm.io/gorm"
)

type postgresStatsRepository struct{ db *gorm.DB }

func NewPostgresStatsRepository(db *gorm.DB) repository.StatsRepository {
	return &postgresStatsRepository{db: db}
}

func (r *postgresStatsRepository) GetGlobalStats(ctx context.Context) (*repository.GlobalStats, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM auth.users WHERE deleted_at IS NULL) AS total_users,
			(SELECT COUNT(*) FROM auth.users WHERE is_active = true AND deleted_at IS NULL) AS total_active_users,
			(SELECT COUNT(*) FROM academic.schools WHERE is_active = true AND deleted_at IS NULL) AS total_schools,
			(SELECT COUNT(*) FROM academic.subjects WHERE is_active = true AND deleted_at IS NULL) AS total_subjects,
			(SELECT COUNT(*) FROM academic.guardian_relations WHERE is_active = true AND deleted_at IS NULL) AS total_guardian_relations
	`

	var stats repository.GlobalStats
	if err := r.db.WithContext(ctx).Raw(query).Scan(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}
