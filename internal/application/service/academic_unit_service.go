package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
)

// AcademicUnitService defines the academic unit service interface
type AcademicUnitService interface {
	CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	ListUnitsBySchool(ctx context.Context, schoolID string) ([]dto.AcademicUnitResponse, error)
	GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error)
	ListUnitsByType(ctx context.Context, schoolID, unitType string) ([]dto.AcademicUnitResponse, error)
	UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error)
	DeleteUnit(ctx context.Context, id string) error
	RestoreUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error)
	GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error)
}

type academicUnitService struct {
	unitRepo   repository.AcademicUnitRepository
	schoolRepo sharedrepo.SchoolRepository
	logger     logger.Logger
}

// NewAcademicUnitService creates a new academic unit service
func NewAcademicUnitService(unitRepo repository.AcademicUnitRepository, schoolRepo sharedrepo.SchoolRepository, logger logger.Logger) AcademicUnitService {
	return &academicUnitService{unitRepo: unitRepo, schoolRepo: schoolRepo, logger: logger}
}

func (s *academicUnitService) CreateUnit(ctx context.Context, schoolID string, req dto.CreateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}

	// Verify school exists
	school, err := s.schoolRepo.FindByID(ctx, sid)
	if err != nil || school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	var parentID *uuid.UUID
	if req.ParentUnitID != nil && *req.ParentUnitID != "" {
		pid, err := uuid.Parse(*req.ParentUnitID)
		if err != nil {
			return nil, errors.NewValidationError("invalid parent_unit_id")
		}
		parentID = &pid
	}

	// Generate code if not provided
	code := req.Code
	if code == "" {
		code = uuid.New().String()[:8]
	}

	// Check code uniqueness within school
	exists, err := s.unitRepo.ExistsBySchoolIDAndCode(ctx, sid, code)
	if err != nil {
		return nil, errors.NewDatabaseError("check unit code", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("academic_unit").WithField("code", code)
	}

	metadataJSON := []byte("{}")
	if req.Metadata != nil {
		metadataJSON, _ = json.Marshal(req.Metadata)
	}

	now := time.Now()
	var desc *string
	if req.Description != "" {
		desc = &req.Description
	}

	unit := &entities.AcademicUnit{
		ID:           uuid.New(),
		ParentUnitID: parentID,
		SchoolID:     sid,
		Name:         req.DisplayName,
		Code:         code,
		Type:         req.Type,
		Description:  desc,
		AcademicYear: 0,
		Metadata:     metadataJSON,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, errors.NewDatabaseError("create academic unit", err)
	}

	s.logger.Info("entity created", "entity_type", "academic_unit", "entity_id", unit.ID.String())
	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) GetUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}
	unit, err := s.unitRepo.FindByID(ctx, uid, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("academic_unit")
	}
	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) ListUnitsBySchool(ctx context.Context, schoolID string) ([]dto.AcademicUnitResponse, error) {
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}
	units, err := s.unitRepo.FindBySchoolID(ctx, sid, false)
	if err != nil {
		return nil, errors.NewDatabaseError("list units", err)
	}
	return dto.ToAcademicUnitResponseList(units), nil
}

func (s *academicUnitService) GetUnitTree(ctx context.Context, schoolID string) ([]*dto.UnitTreeNode, error) {
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}
	units, err := s.unitRepo.FindBySchoolID(ctx, sid, false)
	if err != nil {
		return nil, errors.NewDatabaseError("get unit tree", err)
	}
	return dto.BuildUnitTree(units), nil
}

func (s *academicUnitService) ListUnitsByType(ctx context.Context, schoolID, unitType string) ([]dto.AcademicUnitResponse, error) {
	sid, err := uuid.Parse(schoolID)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}
	if unitType == "" {
		return nil, errors.NewValidationError("type query parameter is required")
	}
	units, err := s.unitRepo.FindByType(ctx, sid, unitType, false)
	if err != nil {
		return nil, errors.NewDatabaseError("list units by type", err)
	}
	return dto.ToAcademicUnitResponseList(units), nil
}

func (s *academicUnitService) UpdateUnit(ctx context.Context, id string, req dto.UpdateAcademicUnitRequest) (*dto.AcademicUnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}
	unit, err := s.unitRepo.FindByID(ctx, uid, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("academic_unit")
	}

	if req.DisplayName != nil && *req.DisplayName != "" {
		unit.Name = *req.DisplayName
	}
	if req.Description != nil {
		unit.Description = req.Description
	}
	if req.ParentUnitID != nil {
		if *req.ParentUnitID == "" {
			unit.ParentUnitID = nil
		} else {
			pid, err := uuid.Parse(*req.ParentUnitID)
			if err != nil {
				return nil, errors.NewValidationError("invalid parent_unit_id")
			}
			unit.ParentUnitID = &pid
		}
	}
	if req.Metadata != nil {
		metadataJSON, _ := json.Marshal(req.Metadata)
		unit.Metadata = metadataJSON
	}

	unit.UpdatedAt = time.Now()
	if err := s.unitRepo.Update(ctx, unit); err != nil {
		return nil, errors.NewDatabaseError("update unit", err)
	}

	s.logger.Info("entity updated", "entity_type", "academic_unit", "entity_id", id)
	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) DeleteUnit(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid unit ID")
	}
	unit, err := s.unitRepo.FindByID(ctx, uid, false)
	if err != nil {
		return errors.NewDatabaseError("find unit", err)
	}
	if unit == nil {
		return errors.NewNotFoundError("academic_unit")
	}
	if err := s.unitRepo.SoftDelete(ctx, uid); err != nil {
		return errors.NewDatabaseError("delete unit", err)
	}
	s.logger.Info("entity deleted", "entity_type", "academic_unit", "entity_id", id)
	return nil
}

func (s *academicUnitService) RestoreUnit(ctx context.Context, id string) (*dto.AcademicUnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}
	if err := s.unitRepo.Restore(ctx, uid); err != nil {
		return nil, errors.NewDatabaseError("restore unit", err)
	}
	unit, err := s.unitRepo.FindByID(ctx, uid, false)
	if err != nil {
		return nil, errors.NewDatabaseError("find restored unit", err)
	}
	s.logger.Info("entity restored", "entity_type", "academic_unit", "entity_id", id)
	response := dto.ToAcademicUnitResponse(unit)
	return &response, nil
}

func (s *academicUnitService) GetHierarchyPath(ctx context.Context, id string) ([]dto.AcademicUnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit ID")
	}
	path, err := s.unitRepo.GetHierarchyPath(ctx, uid)
	if err != nil {
		return nil, errors.NewDatabaseError("get hierarchy path", err)
	}
	return dto.ToAcademicUnitResponseList(path), nil
}
