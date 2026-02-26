package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
)

// SchoolService defines the school service interface
type SchoolService interface {
	CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error)
	GetSchoolByCode(ctx context.Context, code string) (*dto.SchoolResponse, error)
	UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
	ListSchools(ctx context.Context, filters sharedrepo.ListFilters) ([]dto.SchoolResponse, error)
	DeleteSchool(ctx context.Context, id string) error
}

type schoolService struct {
	schoolRepo sharedrepo.SchoolRepository
	logger     logger.Logger
	defaults   config.SchoolDefaults
}

// NewSchoolService creates a new school service
func NewSchoolService(schoolRepo sharedrepo.SchoolRepository, logger logger.Logger, defaults config.SchoolDefaults) SchoolService {
	return &schoolService{schoolRepo: schoolRepo, logger: logger, defaults: defaults}
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
		return nil, errors.NewDatabaseError("create school", err)
	}

	s.logger.Info("entity created", "entity_type", "school", "entity_id", school.ID.String(), "name", school.Name)
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

func (s *schoolService) ListSchools(ctx context.Context, filters sharedrepo.ListFilters) ([]dto.SchoolResponse, error) {
	schools, err := s.schoolRepo.List(ctx, filters)
	if err != nil {
		return nil, errors.NewDatabaseError("list schools", err)
	}
	return dto.ToSchoolResponseList(schools), nil
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
	if err := s.schoolRepo.Delete(ctx, schoolID); err != nil {
		return errors.NewDatabaseError("delete school", err)
	}
	s.logger.Info("entity deleted", "entity_type", "school", "entity_id", id)
	return nil
}
