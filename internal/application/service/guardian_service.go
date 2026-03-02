package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// GuardianService defines the guardian service interface
type GuardianService interface {
	CreateRelation(ctx context.Context, req dto.CreateGuardianRelationRequest, createdBy string) (*dto.GuardianRelationResponse, error)
	GetRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error)
	UpdateRelation(ctx context.Context, id string, req dto.UpdateGuardianRelationRequest) (*dto.GuardianRelationResponse, error)
	DeleteRelation(ctx context.Context, id string) error
	GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error)
	GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error)
}

type guardianService struct {
	guardianRepo repository.GuardianRepository
	logger       logger.Logger
}

// NewGuardianService creates a new guardian service
func NewGuardianService(guardianRepo repository.GuardianRepository, logger logger.Logger) GuardianService {
	return &guardianService{guardianRepo: guardianRepo, logger: logger}
}

func (s *guardianService) CreateRelation(ctx context.Context, req dto.CreateGuardianRelationRequest, createdBy string) (*dto.GuardianRelationResponse, error) {
	guardianID, err := uuid.Parse(req.GuardianID)
	if err != nil {
		return nil, errors.NewValidationError("invalid guardian_id")
	}
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid student_id")
	}
	if req.RelationshipType == "" {
		return nil, errors.NewValidationError("relationship_type is required")
	}

	exists, err := s.guardianRepo.ExistsActiveRelation(ctx, guardianID, studentID)
	if err != nil {
		return nil, errors.NewDatabaseError("check guardian relation", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("guardian_relation")
	}

	var createdByUUID *uuid.UUID
	if parsed, err := uuid.Parse(createdBy); err == nil {
		createdByUUID = &parsed
	}

	now := time.Now()
	relation := &entities.GuardianRelation{
		ID:               uuid.New(),
		GuardianID:       guardianID,
		StudentID:        studentID,
		RelationshipType: req.RelationshipType,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        createdByUUID,
	}

	if err := s.guardianRepo.Create(ctx, relation); err != nil {
		return nil, errors.NewDatabaseError("create guardian relation", err)
	}

	s.logger.Info("entity created", "entity_type", "guardian_relation", "entity_id", relation.ID.String())
	return dto.ToGuardianRelationResponse(relation), nil
}

func (s *guardianService) GetRelation(ctx context.Context, id string) (*dto.GuardianRelationResponse, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid relation ID")
	}
	relation, err := s.guardianRepo.FindByID(ctx, rid)
	if err != nil {
		return nil, errors.NewDatabaseError("find guardian relation", err)
	}
	if relation == nil {
		return nil, errors.NewNotFoundError("guardian_relation")
	}
	return dto.ToGuardianRelationResponse(relation), nil
}

func (s *guardianService) UpdateRelation(ctx context.Context, id string, req dto.UpdateGuardianRelationRequest) (*dto.GuardianRelationResponse, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid relation ID")
	}
	relation, err := s.guardianRepo.FindByID(ctx, rid)
	if err != nil {
		return nil, errors.NewDatabaseError("find guardian relation", err)
	}
	if relation == nil {
		return nil, errors.NewNotFoundError("guardian_relation")
	}

	if req.RelationshipType != nil {
		relation.RelationshipType = *req.RelationshipType
	}
	if req.IsActive != nil {
		relation.IsActive = *req.IsActive
	}
	relation.UpdatedAt = time.Now()

	if err := s.guardianRepo.Update(ctx, relation); err != nil {
		return nil, errors.NewDatabaseError("update guardian relation", err)
	}

	s.logger.Info("entity updated", "entity_type", "guardian_relation", "entity_id", id)
	return dto.ToGuardianRelationResponse(relation), nil
}

func (s *guardianService) DeleteRelation(ctx context.Context, id string) error {
	rid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid relation ID")
	}
	if err := s.guardianRepo.Delete(ctx, rid); err != nil {
		return errors.NewDatabaseError("delete guardian relation", err)
	}
	s.logger.Info("entity deleted", "entity_type", "guardian_relation", "entity_id", id)
	return nil
}

func (s *guardianService) GetGuardianRelations(ctx context.Context, guardianID string) ([]*dto.GuardianRelationResponse, error) {
	gid, err := uuid.Parse(guardianID)
	if err != nil {
		return nil, errors.NewValidationError("invalid guardian_id")
	}
	relations, err := s.guardianRepo.FindByGuardian(ctx, gid)
	if err != nil {
		return nil, errors.NewDatabaseError("find guardian relations", err)
	}
	return dto.ToGuardianRelationResponseList(relations), nil
}

func (s *guardianService) GetStudentGuardians(ctx context.Context, studentID string) ([]*dto.GuardianRelationResponse, error) {
	sid, err := uuid.Parse(studentID)
	if err != nil {
		return nil, errors.NewValidationError("invalid student_id")
	}
	relations, err := s.guardianRepo.FindByStudent(ctx, sid)
	if err != nil {
		return nil, errors.NewDatabaseError("find student guardians", err)
	}
	return dto.ToGuardianRelationResponseList(relations), nil
}
