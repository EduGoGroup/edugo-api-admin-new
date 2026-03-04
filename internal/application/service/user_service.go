package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/audit"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
)

// UserService defines the user service interface
type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id string) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, filters sharedrepo.ListFilters) ([]*dto.UserResponse, int, error)
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepo    sharedrepo.UserRepository
	logger      logger.Logger
	auditLogger audit.AuditLogger
}

// NewUserService creates a new user service
func NewUserService(userRepo sharedrepo.UserRepository, logger logger.Logger, auditLogger audit.AuditLogger) UserService {
	return &userService{userRepo: userRepo, logger: logger, auditLogger: auditLogger}
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

	// Usamos nuestra función flexible para parsear el estado
	isActive, err := parseFlexibleBool(req.IsActive, true)
	if err != nil {
		return nil, errors.NewValidationError("invalid is_active value: " + err.Error())
	}

	user := &entities.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     isActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("create user", err)
	}

	s.logger.Info("entity created", "entity_type", "user", "entity_id", user.ID.String(), "email", user.Email)

	_ = s.auditLogger.Log(ctx, audit.AuditEvent{
		Action:       "create",
		ResourceType: "user",
		ResourceID:   user.ID.String(),
		Severity:     audit.SeverityCritical,
		Category:     audit.CategoryAdmin,
		Metadata:     map[string]interface{}{"email": user.Email},
	})

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

func (s *userService) ListUsers(ctx context.Context, filters sharedrepo.ListFilters) ([]*dto.UserResponse, int, error) {
	users, total, err := s.userRepo.List(ctx, filters)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("list users", err)
	}
	return dto.ToUserResponseList(users), int(total), nil
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
		isActive, err := parseFlexibleBool(req.IsActive, user.IsActive)
		if err != nil {
			return nil, errors.NewValidationError("invalid is_active value: " + err.Error())
		}
		user.IsActive = isActive
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

	_ = s.auditLogger.Log(ctx, audit.AuditEvent{
		Action:       "delete",
		ResourceType: "user",
		ResourceID:   id,
		Severity:     audit.SeverityCritical,
		Category:     audit.CategoryAdmin,
	})

	return nil
}

// parseFlexibleBool converts a bool or string value to a boolean.
// Returns an error for unrecognized string values.
func parseFlexibleBool(val interface{}, defaultVal bool) (bool, error) {
	if val == nil {
		return defaultVal, nil
	}
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean value: %q (accepted: true/false, 1/0, yes/no)", v)
		}
	default:
		strVal := fmt.Sprintf("%v", v)
		switch strings.ToLower(strings.TrimSpace(strVal)) {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean value: %q (accepted: true/false, 1/0, yes/no)", strVal)
		}
	}
}
