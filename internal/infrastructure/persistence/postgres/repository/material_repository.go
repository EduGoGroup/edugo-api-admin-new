package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresMaterialRepository struct{ db *gorm.DB }

func NewPostgresMaterialRepository(db *gorm.DB) repository.MaterialRepository {
	return &postgresMaterialRepository{db: db}
}

func (r *postgresMaterialRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Material{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *postgresMaterialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Material{}, "id = ?", id).Error
}
