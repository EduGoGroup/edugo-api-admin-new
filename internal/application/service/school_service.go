package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
)

// SchoolActorContextKey is the unexported type used to store audit actor values
// in context.Context. Defined here so the handler and service use the same type.
type SchoolActorContextKey string

// Context keys for audit actor data, set by the school handler and read by the service.
const (
	SchoolActorIDKey    SchoolActorContextKey = "actor_id"
	SchoolActorEmailKey SchoolActorContextKey = "actor_email"
	SchoolActorRoleKey  SchoolActorContextKey = "actor_role"
)

// actorFromContext extracts the audit actor fields stored by the school handler.
func actorFromContext(ctx context.Context) (id, email, role string) {
	if v, ok := ctx.Value(SchoolActorIDKey).(string); ok {
		id = v
	}
	if v, ok := ctx.Value(SchoolActorEmailKey).(string); ok {
		email = v
	}
	if v, ok := ctx.Value(SchoolActorRoleKey).(string); ok {
		role = v
	}
	return id, email, role
}

// SchoolService defines the school service interface
type SchoolService interface {
	CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error)
	GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error)
	UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
	ListSchools(ctx context.Context, filters sharedrepo.ListFilters) ([]dto.SchoolResponse, int, error)
	DeleteSchool(ctx context.Context, id string) error
}

type schoolService struct {
	schoolRepo  sharedrepo.SchoolRepository
	logger      logger.Logger
	defaults    config.SchoolDefaults
	auditLogger audit.AuditLogger
}

// NewSchoolService creates a new school service
func NewSchoolService(schoolRepo sharedrepo.SchoolRepository, logger logger.Logger, defaults config.SchoolDefaults, auditLogger audit.AuditLogger) SchoolService {
	return &schoolService{schoolRepo: schoolRepo, logger: logger, defaults: defaults, auditLogger: auditLogger}
}

func (s *schoolService) CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
	exists, err := s.schoolRepo.ExistsByCode(ctx, req.Code)
	if err != nil {
		return nil, errors.NewDatabaseError("check school", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("school").WithField("code", req.Code)
	}

	if req.Name == "" || len(req.Name) < 3 {
		return nil, errors.NewValidationError("name must be at least 3 characters")
	}
	if req.Code == "" || len(req.Code) < 3 {
		return nil, errors.NewValidationError("code must be at least 3 characters")
	}

	metadataJSON := []byte("{}")
	if req.Metadata != nil {
		metadataJSON, _ = json.Marshal(req.Metadata)
	}

	now := time.Now()
	addr := &req.Address
	email := &req.ContactEmail
	phone := &req.ContactPhone

	country := req.Country
	if country == "" {
		country = s.defaults.Country
	}
	subscriptionTier := req.SubscriptionTier
	if subscriptionTier == "" {
		subscriptionTier = s.defaults.SubscriptionTier
	}
	maxTeachers := req.MaxTeachers
	if maxTeachers == 0 {
		maxTeachers = s.defaults.MaxTeachers
	}
	maxStudents := req.MaxStudents
	if maxStudents == 0 {
		maxStudents = s.defaults.MaxStudents
	}

	var city *string
	if req.City != "" {
		city = &req.City
	}

	school := &entities.School{
		ID:               uuid.New(),
		Name:             req.Name,
		Code:             req.Code,
		Address:          addr,
		City:             city,
		Country:          country,
		Phone:            phone,
		Email:            email,
		Metadata:         metadataJSON,
		IsActive:         true,
		SubscriptionTier: subscriptionTier,
		MaxTeachers:      maxTeachers,
		MaxStudents:      maxStudents,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.schoolRepo.Create(ctx, school); err != nil {
		actorID, actorEmail, actorRole := actorFromContext(ctx)
		s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "create", ResourceType: "school",
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityCritical, Category: audit.CategoryAdmin,
		})
		return nil, errors.NewDatabaseError("create school", err)
	}

	s.logger.Info("entity created", "entity_type", "school", "entity_id", school.ID.String(), "name", school.Name)
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "create", ResourceType: "school", ResourceID: school.ID.String(),
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityCritical, Category: audit.CategoryAdmin,
	})
	response := dto.ToSchoolResponse(school)
	return &response, nil
}

func (s *schoolService) GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error) {
	schoolID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}
	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewDatabaseError("find school", err)
	}
	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}
	response := dto.ToSchoolResponse(school)
	return &response, nil
}

func (s *schoolService) GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error) {
	school, err := s.schoolRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, errors.NewDatabaseError("find school", err)
	}
	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}
	response := dto.ToSchoolResponse(school)
	return &response, nil
}

func (s *schoolService) UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error) {
	schoolID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid school ID")
	}
	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return nil, errors.NewDatabaseError("find school", err)
	}
	if school == nil {
		return nil, errors.NewNotFoundError("school")
	}

	if req.Name != nil && *req.Name != "" {
		if len(*req.Name) < 3 {
			return nil, errors.NewValidationError("name must be at least 3 characters")
		}
		school.Name = *req.Name
	}
	if req.Address != nil {
		school.Address = req.Address
	}
	if req.ContactEmail != nil {
		school.Email = req.ContactEmail
	}
	if req.ContactPhone != nil {
		school.Phone = req.ContactPhone
	}
	if req.City != nil {
		school.City = req.City
	}
	if req.Country != nil && *req.Country != "" {
		school.Country = *req.Country
	}
	if req.SubscriptionTier != nil && *req.SubscriptionTier != "" {
		school.SubscriptionTier = *req.SubscriptionTier
	}
	if req.MaxTeachers != nil && *req.MaxTeachers > 0 {
		school.MaxTeachers = *req.MaxTeachers
	}
	if req.MaxStudents != nil && *req.MaxStudents > 0 {
		school.MaxStudents = *req.MaxStudents
	}
	if req.Metadata != nil {
		metadataJSON, _ := json.Marshal(req.Metadata)
		school.Metadata = metadataJSON
	}

	school.UpdatedAt = time.Now()
	if err := s.schoolRepo.Update(ctx, school); err != nil {
		return nil, errors.NewDatabaseError("update school", err)
	}

	s.logger.Info("entity updated", "entity_type", "school", "entity_id", id)
	response := dto.ToSchoolResponse(school)
	return &response, nil
}

func (s *schoolService) ListSchools(ctx context.Context, filters sharedrepo.ListFilters) ([]dto.SchoolResponse, int, error) {
	schools, total, err := s.schoolRepo.List(ctx, filters)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list schools", err)
	}
	return dto.ToSchoolResponseList(schools), int(total), nil
}

func (s *schoolService) DeleteSchool(ctx context.Context, id string) error {
	schoolID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid school ID")
	}
	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		return errors.NewDatabaseError("find school", err)
	}
	if school == nil {
		return errors.NewNotFoundError("school")
	}
	actorID, actorEmail, actorRole := actorFromContext(ctx)
	if err := s.schoolRepo.Delete(ctx, schoolID); err != nil {
		s.auditLogger.Log(ctx, audit.AuditEvent{
			Action: "delete", ResourceType: "school", ResourceID: id,
			ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
			ErrorMessage: err.Error(), Severity: audit.SeverityCritical, Category: audit.CategoryAdmin,
		})
		return errors.NewDatabaseError("delete school", err)
	}
	s.logger.Info("entity deleted", "entity_type", "school", "entity_id", id)
	s.auditLogger.Log(ctx, audit.AuditEvent{
		Action: "delete", ResourceType: "school", ResourceID: id,
		ActorID: actorID, ActorEmail: actorEmail, ActorRole: actorRole,
		Severity: audit.SeverityCritical, Category: audit.CategoryAdmin,
	})
	return nil
}
