package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// UserService defines the user service interface
type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id string) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, filters repository.ListFilters) ([]*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
	logger   logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, logger logger.Logger) UserService {
	return &userService{userRepo: userRepo, logger: logger}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("check user", err)
	}
	if exists {
		return nil, errors.NewAlreadyExistsError("user").WithField("email", req.Email)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewValidationError("invalid password: " + err.Error())
	}

	now := time.Now()
	user := &entities.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("create user", err)
	}

	s.logger.Info("entity created", "entity_type", "user", "entity_id", user.ID.String(), "email", user.Email)
	return dto.ToUserResponse(user), nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user ID")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("find user", err)
	}
	if user == nil {
		return nil, errors.NewNotFoundError("user")
	}
	return dto.ToUserResponse(user), nil
}

func (s *userService) ListUsers(ctx context.Context, filters repository.ListFilters) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.List(ctx, filters)
	if err != nil {
		return nil, errors.NewDatabaseError("list users", err)
	}
	return dto.ToUserResponseList(users), nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user ID")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("find user", err)
	}
	if user == nil {
		return nil, errors.NewNotFoundError("user")
	}

	if req.FirstName != nil && *req.FirstName != "" {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != "" {
		user.LastName = *req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("update user", err)
	}

	s.logger.Info("entity updated", "entity_type", "user", "entity_id", id)
	return dto.ToUserResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid user ID")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewDatabaseError("find user", err)
	}
	if user == nil {
		return errors.NewNotFoundError("user")
	}
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errors.NewDatabaseError("delete user", err)
	}
	s.logger.Info("entity deleted", "entity_type", "user", "entity_id", id)
	return nil
}
