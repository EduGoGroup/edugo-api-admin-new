package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// ConceptTypeService defines the concept type service interface
type ConceptTypeService interface {
	CreateConceptType(ctx context.Context, req *dto.CreateConceptTypeRequest) (*dto.ConceptTypeResponse, error)
	ListConceptTypes(ctx context.Context) ([]dto.ConceptTypeResponse, error)
	GetConceptType(ctx context.Context, id uuid.UUID) (*dto.ConceptTypeResponse, error)
	UpdateConceptType(ctx context.Context, id uuid.UUID, req *dto.UpdateConceptTypeRequest) (*dto.ConceptTypeResponse, error)
	DeleteConceptType(ctx context.Context, id uuid.UUID) error

	// Definitions
	CreateDefinition(ctx context.Context, typeID uuid.UUID, req *dto.ConceptDefinitionRequest) (*dto.ConceptDefinitionResponse, error)
	ListDefinitions(ctx context.Context, typeID uuid.UUID) ([]dto.ConceptDefinitionResponse, error)
	UpdateDefinition(ctx context.Context, typeID uuid.UUID, defID uuid.UUID, req *dto.ConceptDefinitionRequest) (*dto.ConceptDefinitionResponse, error)
	DeleteDefinition(ctx context.Context, typeID uuid.UUID, defID uuid.UUID) error

	// School concepts (read + personalize)
	GetSchoolConcepts(ctx context.Context, schoolID uuid.UUID) ([]dto.SchoolConceptResponse, error)
	GetSchoolConcept(ctx context.Context, schoolID uuid.UUID, conceptID uuid.UUID) (*dto.SchoolConceptResponse, error)
	UpdateSchoolConcept(ctx context.Context, schoolID uuid.UUID, conceptID uuid.UUID, req *dto.UpdateSchoolConceptRequest) (*dto.SchoolConceptResponse, error)
}

type conceptTypeService struct {
	conceptTypeRepo   repository.ConceptTypeRepository
	conceptDefRepo    repository.ConceptDefinitionRepository
	schoolConceptRepo repository.SchoolConceptRepository
	logger            logger.Logger
	auditLogger       audit.AuditLogger
}

// NewConceptTypeService creates a new concept type service
func NewConceptTypeService(
	conceptTypeRepo repository.ConceptTypeRepository,
	conceptDefRepo repository.ConceptDefinitionRepository,
	schoolConceptRepo repository.SchoolConceptRepository,
	logger logger.Logger,
	auditLogger audit.AuditLogger,
) ConceptTypeService {
	return &conceptTypeService{
		conceptTypeRepo:   conceptTypeRepo,
		conceptDefRepo:    conceptDefRepo,
		schoolConceptRepo: schoolConceptRepo,
		logger:            logger,
		auditLogger:       auditLogger,
	}
}

func (s *conceptTypeService) CreateConceptType(ctx context.Context, req *dto.CreateConceptTypeRequest) (*dto.ConceptTypeResponse, error) {
	existing, err := s.conceptTypeRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, errors.NewDatabaseError("check concept type", err)
	}
	if existing != nil {
		return nil, errors.NewAlreadyExistsError("concept_type").WithField("code", req.Code)
	}

	if req.Name == "" || len(req.Name) < 3 {
		return nil, errors.NewValidationError("name must be at least 3 characters")
	}
	if req.Code == "" || len(req.Code) < 3 {
		return nil, errors.NewValidationError("code must be at least 3 characters")
	}

	now := time.Now()
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	ct := &entities.ConceptType{
		ID:          uuid.New(),
		Name:        req.Name,
		Code:        req.Code,
		Description: description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.conceptTypeRepo.Create(ctx, ct); err != nil {
		actorID, actorEmail, actorRole := actorFromContext(ctx)
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "create", ResourceType: "concept_type",
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return nil, errors.NewDatabaseError("create concept type", err)
	}

	s.logger.Info("entity created", "entity_type", "concept_type", "entity_id", ct.ID.String(), "name", ct.Name)
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "create", ResourceType: "concept_type", ResourceID: ct.ID.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	response := dto.ToConceptTypeResponse(ct)
	return &response, nil
}

func (s *conceptTypeService) ListConceptTypes(ctx context.Context) ([]dto.ConceptTypeResponse, error) {
	types, err := s.conceptTypeRepo.FindAll(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError("list concept types", err)
	}
	return dto.ToConceptTypeResponseList(types), nil
}

func (s *conceptTypeService) GetConceptType(ctx context.Context, id uuid.UUID) (*dto.ConceptTypeResponse, error) {
	ct, err := s.conceptTypeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return nil, errors.NewNotFoundError("concept_type")
	}
	response := dto.ToConceptTypeResponse(ct)
	return &response, nil
}

func (s *conceptTypeService) UpdateConceptType(ctx context.Context, id uuid.UUID, req *dto.UpdateConceptTypeRequest) (*dto.ConceptTypeResponse, error) {
	ct, err := s.conceptTypeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return nil, errors.NewNotFoundError("concept_type")
	}

	if req.Name != nil && *req.Name != "" {
		if len(*req.Name) < 3 {
			return nil, errors.NewValidationError("name must be at least 3 characters")
		}
		ct.Name = *req.Name
	}
	if req.Description != nil {
		ct.Description = req.Description
	}

	ct.UpdatedAt = time.Now()
	if err := s.conceptTypeRepo.Update(ctx, ct); err != nil {
		actorID, actorEmail, actorRole := actorFromContext(ctx)
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "update", ResourceType: "concept_type", ResourceID: id.String(),
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return nil, errors.NewDatabaseError("update concept type", err)
	}

	s.logger.Info("entity updated", "entity_type", "concept_type", "entity_id", id.String())
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "update", ResourceType: "concept_type", ResourceID: id.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	response := dto.ToConceptTypeResponse(ct)
	return &response, nil
}

func (s *conceptTypeService) DeleteConceptType(ctx context.Context, id uuid.UUID) error {
	ct, err := s.conceptTypeRepo.FindByID(ctx, id)
	if err != nil {
		return errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return errors.NewNotFoundError("concept_type")
	}

	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.conceptTypeRepo.SoftDelete(ctx, id); err != nil {
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "delete", ResourceType: "concept_type", ResourceID: id.String(),
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return errors.NewDatabaseError("delete concept type", err)
	}
	s.logger.Info("entity deleted", "entity_type", "concept_type", "entity_id", id.String())
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "delete", ResourceType: "concept_type", ResourceID: id.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	return nil
}

// ==================== Definitions ====================

func (s *conceptTypeService) CreateDefinition(ctx context.Context, typeID uuid.UUID, req *dto.ConceptDefinitionRequest) (*dto.ConceptDefinitionResponse, error) {
	ct, err := s.conceptTypeRepo.FindByID(ctx, typeID)
	if err != nil {
		return nil, errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return nil, errors.NewNotFoundError("concept_type")
	}

	category := req.Category
	if category == "" {
		category = "general"
	}

	now := time.Now()
	def := &entities.ConceptDefinition{
		ID:            uuid.New(),
		ConceptTypeID: typeID,
		TermKey:       req.TermKey,
		TermValue:     req.TermValue,
		Category:      category,
		SortOrder:     req.SortOrder,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.conceptDefRepo.Create(ctx, def); err != nil {
		actorID, actorEmail, actorRole := actorFromContext(ctx)
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "create", ResourceType: "concept_definition",
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return nil, errors.NewDatabaseError("create concept definition", err)
	}

	s.logger.Info("entity created", "entity_type", "concept_definition", "entity_id", def.ID.String())
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "create", ResourceType: "concept_definition", ResourceID: def.ID.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	response := dto.ToConceptDefinitionResponse(def)
	return &response, nil
}

func (s *conceptTypeService) ListDefinitions(ctx context.Context, typeID uuid.UUID) ([]dto.ConceptDefinitionResponse, error) {
	ct, err := s.conceptTypeRepo.FindByID(ctx, typeID)
	if err != nil {
		return nil, errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return nil, errors.NewNotFoundError("concept_type")
	}

	defs, err := s.conceptDefRepo.FindByTypeID(ctx, typeID)
	if err != nil {
		return nil, errors.NewDatabaseError("list concept definitions", err)
	}
	return dto.ToConceptDefinitionResponseList(defs), nil
}

func (s *conceptTypeService) UpdateDefinition(ctx context.Context, typeID uuid.UUID, defID uuid.UUID, req *dto.ConceptDefinitionRequest) (*dto.ConceptDefinitionResponse, error) {
	ct, err := s.conceptTypeRepo.FindByID(ctx, typeID)
	if err != nil {
		return nil, errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return nil, errors.NewNotFoundError("concept_type")
	}

	// Fetch the target definition directly and verify it belongs to the given concept type
	target, err := s.conceptDefRepo.FindByID(ctx, defID)
	if err != nil {
		return nil, errors.NewDatabaseError("find concept definition", err)
	}
	if target == nil || target.ConceptTypeID != typeID {
		return nil, errors.NewNotFoundError("concept_definition")
	}

	target.TermKey = req.TermKey
	target.TermValue = req.TermValue
	if req.Category != "" {
		target.Category = req.Category
	}
	target.SortOrder = req.SortOrder
	target.UpdatedAt = time.Now()

	if err := s.conceptDefRepo.Update(ctx, target); err != nil {
		actorID, actorEmail, actorRole := actorFromContext(ctx)
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "update", ResourceType: "concept_definition", ResourceID: defID.String(),
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return nil, errors.NewDatabaseError("update concept definition", err)
	}

	s.logger.Info("entity updated", "entity_type", "concept_definition", "entity_id", defID.String())
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "update", ResourceType: "concept_definition", ResourceID: defID.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	response := dto.ToConceptDefinitionResponse(target)
	return &response, nil
}

func (s *conceptTypeService) DeleteDefinition(ctx context.Context, typeID uuid.UUID, defID uuid.UUID) error {
	ct, err := s.conceptTypeRepo.FindByID(ctx, typeID)
	if err != nil {
		return errors.NewDatabaseError("find concept type", err)
	}
	if ct == nil {
		return errors.NewNotFoundError("concept_type")
	}

	// Fetch the target definition directly and verify it belongs to the given concept type
	target, err := s.conceptDefRepo.FindByID(ctx, defID)
	if err != nil {
		return errors.NewDatabaseError("find concept definition", err)
	}
	if target == nil || target.ConceptTypeID != typeID {
		return errors.NewNotFoundError("concept_definition")
	}

	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.conceptDefRepo.Delete(ctx, defID); err != nil {
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "delete", ResourceType: "concept_definition", ResourceID: defID.String(),
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return errors.NewDatabaseError("delete concept definition", err)
	}

	s.logger.Info("entity deleted", "entity_type", "concept_definition", "entity_id", defID.String())
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "delete", ResourceType: "concept_definition", ResourceID: defID.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	return nil
}

// ==================== School Concepts ====================

func (s *conceptTypeService) GetSchoolConcepts(ctx context.Context, schoolID uuid.UUID) ([]dto.SchoolConceptResponse, error) {
	concepts, err := s.schoolConceptRepo.FindBySchoolID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewDatabaseError("list school concepts", err)
	}
	return dto.ToSchoolConceptResponseList(concepts), nil
}

func (s *conceptTypeService) GetSchoolConcept(ctx context.Context, schoolID uuid.UUID, conceptID uuid.UUID) (*dto.SchoolConceptResponse, error) {
	concept, err := s.schoolConceptRepo.FindByID(ctx, conceptID)
	if err != nil {
		return nil, errors.NewDatabaseError("find school concept", err)
	}
	if concept == nil || concept.SchoolID != schoolID {
		return nil, errors.NewNotFoundError("school_concept")
	}
	response := dto.ToSchoolConceptResponse(concept)
	return &response, nil
}

func (s *conceptTypeService) UpdateSchoolConcept(ctx context.Context, schoolID uuid.UUID, conceptID uuid.UUID, req *dto.UpdateSchoolConceptRequest) (*dto.SchoolConceptResponse, error) {
	concept, err := s.schoolConceptRepo.FindByID(ctx, conceptID)
	if err != nil {
		return nil, errors.NewDatabaseError("find school concept", err)
	}
	if concept == nil {
		return nil, errors.NewNotFoundError("school_concept")
	}
	if concept.SchoolID != schoolID {
		return nil, errors.NewNotFoundError("school_concept")
	}

	concept.TermValue = req.TermValue
	concept.UpdatedAt = time.Now()

	if err := s.schoolConceptRepo.Update(ctx, concept); err != nil {
		actorID, actorEmail, actorRole := actorFromContext(ctx)
		if logErr := s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "update", ResourceType: "school_concept", ResourceID: conceptID.String(),
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityWarning, Category: audit.CategoryAdmin,
		}); logErr != nil {
			s.logger.Error("failed to write audit log", "error", logErr)
		}
		return nil, errors.NewDatabaseError("update school concept", err)
	}

	s.logger.Info("entity updated", "entity_type", "school_concept", "entity_id", conceptID.String())
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "update", ResourceType: "school_concept", ResourceID: conceptID.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityInfo, Category: audit.CategoryAdmin,
	}); err != nil {
		s.logger.Error("failed to write audit log", "error", err)
	}
	response := dto.ToSchoolConceptResponse(concept)
	return &response, nil
}
