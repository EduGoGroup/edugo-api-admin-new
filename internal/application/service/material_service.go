package service

import (
	"context"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// MaterialService defines the material service interface
type MaterialService interface {
	DeleteMaterial(ctx context.Context, id string) error
}

type materialService struct {
	materialRepo repository.MaterialRepository
	logger       logger.Logger
}

// NewMaterialService creates a new material service
func NewMaterialService(materialRepo repository.MaterialRepository, logger logger.Logger) MaterialService {
	return &materialService{materialRepo: materialRepo, logger: logger}
}

func (s *materialService) DeleteMaterial(ctx context.Context, id string) error {
	materialID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid material ID")
	}
	exists, err := s.materialRepo.Exists(ctx, materialID)
	if err != nil {
		return errors.NewDatabaseError("find material", err)
	}
	if !exists {
		return errors.NewNotFoundError("material")
	}
	if err := s.materialRepo.Delete(ctx, materialID); err != nil {
		return errors.NewDatabaseError("delete material", err)
	}
	s.logger.Info("entity deleted", "entity_type", "material", "entity_id", id)
	return nil
}
