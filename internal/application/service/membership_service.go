package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
)

// MembershipService defines the membership service interface
type MembershipService interface {
	CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error)
	GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error)
	ListMembershipsByUnit(ctx context.Context, unitID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error)
	ListMembershipsByRole(ctx context.Context, unitID, role string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error)
	ListMembershipsByUser(ctx context.Context, userID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error)
	UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error)
	DeleteMembership(ctx context.Context, id string) error
	ExpireMembership(ctx context.Context, id string) (*dto.MembershipResponse, error)
}

type membershipService struct {
	membershipRepo sharedrepo.MembershipRepository
	logger         logger.Logger
}

// NewMembershipService creates a new membership service
func NewMembershipService(membershipRepo sharedrepo.MembershipRepository, logger logger.Logger) MembershipService {
	return &membershipService{membershipRepo: membershipRepo, logger: logger}
}

func (s *membershipService) CreateMembership(ctx context.Context, req dto.CreateMembershipRequest) (*dto.MembershipResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}
	unitID, err := uuid.Parse(req.UnitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}
	if req.Role == "" {
		return nil, errors.NewValidationError("role is required")
	}

	now := time.Now()
	membership := &entities.Membership{
		ID:             uuid.New(),
		UserID:         userID,
		SchoolID:       uuid.Nil, // Will be set by the DB trigger or parent context
		AcademicUnitID: &unitID,
		Role:           req.Role,
		Metadata:       []byte("{}"),
		IsActive:       true,
		EnrolledAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		return nil, errors.NewDatabaseError("create membership", err)
	}

	s.logger.Info("entity created", "entity_type", "membership", "entity_id", membership.ID.String())
	response := dto.ToMembershipResponse(membership)
	return &response, nil
}

func (s *membershipService) GetMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	mid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}
	m, err := s.membershipRepo.FindByID(ctx, mid)
	if err != nil {
		return nil, errors.NewDatabaseError("find membership", err)
	}
	if m == nil {
		return nil, errors.NewNotFoundError("membership")
	}
	response := dto.ToMembershipResponse(m)
	return &response, nil
}

func (s *membershipService) ListMembershipsByUnit(ctx context.Context, unitID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error) {
	uid, err := uuid.Parse(unitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}
	memberships, err := s.membershipRepo.FindByUnit(ctx, uid, filters)
	if err != nil {
		return nil, errors.NewDatabaseError("list memberships", err)
	}
	return dto.ToMembershipResponseList(memberships), nil
}

func (s *membershipService) ListMembershipsByRole(ctx context.Context, unitID, role string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error) {
	uid, err := uuid.Parse(unitID)
	if err != nil {
		return nil, errors.NewValidationError("invalid unit_id")
	}
	memberships, err := s.membershipRepo.FindByUnitAndRole(ctx, uid, role, true, filters)
	if err != nil {
		return nil, errors.NewDatabaseError("list memberships by role", err)
	}
	return dto.ToMembershipResponseList(memberships), nil
}

func (s *membershipService) ListMembershipsByUser(ctx context.Context, userID string, filters sharedrepo.ListFilters) ([]dto.MembershipResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id")
	}
	memberships, err := s.membershipRepo.FindByUser(ctx, uid, filters)
	if err != nil {
		return nil, errors.NewDatabaseError("list user memberships", err)
	}
	return dto.ToMembershipResponseList(memberships), nil
}

func (s *membershipService) UpdateMembership(ctx context.Context, id string, req dto.UpdateMembershipRequest) (*dto.MembershipResponse, error) {
	mid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}
	m, err := s.membershipRepo.FindByID(ctx, mid)
	if err != nil {
		return nil, errors.NewDatabaseError("find membership", err)
	}
	if m == nil {
		return nil, errors.NewNotFoundError("membership")
	}

	if req.Role != nil && *req.Role != "" {
		m.Role = *req.Role
	}
	m.UpdatedAt = time.Now()

	if err := s.membershipRepo.Update(ctx, m); err != nil {
		return nil, errors.NewDatabaseError("update membership", err)
	}

	s.logger.Info("entity updated", "entity_type", "membership", "entity_id", id)
	response := dto.ToMembershipResponse(m)
	return &response, nil
}

func (s *membershipService) DeleteMembership(ctx context.Context, id string) error {
	mid, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid membership ID")
	}
	m, err := s.membershipRepo.FindByID(ctx, mid)
	if err != nil {
		return errors.NewDatabaseError("find membership", err)
	}
	if m == nil {
		return errors.NewNotFoundError("membership")
	}
	if err := s.membershipRepo.Delete(ctx, mid); err != nil {
		return errors.NewDatabaseError("delete membership", err)
	}
	s.logger.Info("entity deleted", "entity_type", "membership", "entity_id", id)
	return nil
}

func (s *membershipService) ExpireMembership(ctx context.Context, id string) (*dto.MembershipResponse, error) {
	mid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid membership ID")
	}
	m, err := s.membershipRepo.FindByID(ctx, mid)
	if err != nil {
		return nil, errors.NewDatabaseError("find membership", err)
	}
	if m == nil {
		return nil, errors.NewNotFoundError("membership")
	}

	now := time.Now()
	m.WithdrawnAt = &now
	m.IsActive = false
	m.UpdatedAt = now

	if err := s.membershipRepo.Update(ctx, m); err != nil {
		return nil, errors.NewDatabaseError("expire membership", err)
	}

	s.logger.Info("membership expired", "entity_type", "membership", "entity_id", id)
	response := dto.ToMembershipResponse(m)
	return &response, nil
}
