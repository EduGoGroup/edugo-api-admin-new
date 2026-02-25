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

// SubjectService defines the subject service interface
type SubjectService interface {
	CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error)
	GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error)
	ListSubjects(ctx context.Context) ([]dto.SubjectResponse, error)
	UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error)
	DeleteSubject(ctx context.Context, id string) error
}

type subjectService struct {
	subjectRepo repository.SubjectRepository
	logger      logger.Logger
}

// NewSubjectService creates a new subject service
func NewSubjectService(subjectRepo repository.SubjectRepository, logger logger.Logger) SubjectService {
	return &subjectService{subjectRepo: subjectRepo, logger: logger}
}

func (s *subjectService) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	if req.Name == "" || len(req.Name) < 2 {
		return nil, errors.NewValidationError("name must be at least 2 characters")
	}

	exists, err := s.subjectRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, errors.NewDatabaseError("check subject", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("subject").WithField("name", req.Name)
	}

	now := time.Now()
	subject := &entities.Subject{
		ID:        uuid.New(),
		Name:      req.Name,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if req.Description != "" {
		subject.Description = &req.Description
	}
	if req.Metadata != "" {
		subject.Metadata = &req.Metadata
	}

	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		return nil, errors.NewDatabaseError("create subject", err)
	}

	s.logger.Info("entity created", "entity_type", "subject", "entity_id", subject.ID.String())
	response := dto.ToSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error) {
	sid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid subject ID")
	}
	subject, err := s.subjectRepo.FindByID(ctx, sid)
	if err != nil {
		return nil, errors.NewDatabaseError("find subject", err)
	}
	if subject == nil {
		return nil, errors.NewNotFoundError("subject")
	}
	response := dto.ToSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) ListSubjects(ctx context.Context) ([]dto.SubjectResponse, error) {
	subjects, err := s.subjectRepo.List(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError("list subjects", err)
	}
	return dto.ToSubjectResponseList(subjects), nil
}

func (s *subjectService) UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	sid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid subject ID")
	}
	subject, err := s.subjectRepo.FindByID(ctx, sid)
	if err != nil {
		return nil, errors.NewDatabaseError("find subject", err)
	}
	if subject == nil {
		return nil, errors.NewNotFoundError("subject")
	}

	if req.Name != nil && *req.Name != "" {
		subject.Name = *req.Name
	}
	if req.Description != nil {
		subject.Description = req.Description
	}
	if req.Metadata != nil {
		subject.Metadata = req.Metadata
	}
	subject.UpdatedAt = time.Now()

	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		return nil, errors.NewDatabaseError("update subject", err)
	}

	s.logger.Info("entity updated", "entity_type", "subject", "entity_id", id)
	response := dto.ToSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) DeleteSubject(ctx context.Context, id string) error {
	sid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid subject ID")
	}
	subject, err := s.subjectRepo.FindByID(ctx, sid)
	if err != nil {
		return errors.NewDatabaseError("find subject", err)
	}
	if subject == nil {
		return errors.NewNotFoundError("subject")
	}
	if err := s.subjectRepo.Delete(ctx, sid); err != nil {
		return errors.NewDatabaseError("delete subject", err)
	}
	s.logger.Info("entity deleted", "entity_type", "subject", "entity_id", id)
	return nil
}
