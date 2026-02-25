package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMembershipService_CreateMembership(t *testing.T) {
	validUserID := uuid.New().String()
	validUnitID := uuid.New().String()

	tests := []struct {
		name        string
		request     dto.CreateMembershipRequest
		setupMock   func(m *mock.MockMembershipRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - creates membership",
			request: dto.CreateMembershipRequest{UserID: validUserID, UnitID: validUnitID, Role: "student"},
			setupMock: func(m *mock.MockMembershipRepository) {
				m.CreateFn = func(_ context.Context, _ *entities.Membership) error { return nil }
			},
			wantErr: false,
		},
		{
			name:        "error - invalid user_id",
			request:     dto.CreateMembershipRequest{UserID: "bad", UnitID: validUnitID, Role: "student"},
			setupMock:   func(_ *mock.MockMembershipRepository) {},
			wantErr:     true,
			errContains: "invalid user_id",
		},
		{
			name:        "error - invalid unit_id",
			request:     dto.CreateMembershipRequest{UserID: validUserID, UnitID: "bad", Role: "student"},
			setupMock:   func(_ *mock.MockMembershipRepository) {},
			wantErr:     true,
			errContains: "invalid unit_id",
		},
		{
			name:        "error - empty role",
			request:     dto.CreateMembershipRequest{UserID: validUserID, UnitID: validUnitID, Role: ""},
			setupMock:   func(_ *mock.MockMembershipRepository) {},
			wantErr:     true,
			errContains: "role is required",
		},
		{
			name:    "error - database error",
			request: dto.CreateMembershipRequest{UserID: validUserID, UnitID: validUnitID, Role: "student"},
			setupMock: func(m *mock.MockMembershipRepository) {
				m.CreateFn = func(_ context.Context, _ *entities.Membership) error { return fmt.Errorf("db error") }
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockMembershipRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewMembershipService(mockRepo, mock.NewMockLogger())
			result, err := svc.CreateMembership(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "student", result.Role)
			}
		})
	}
}

func TestMembershipService_GetMembership(t *testing.T) {
	validID := uuid.New()
	unitID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockMembershipRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockMembershipRepository) {
				m.FindByIDFn = func(_ context.Context, id uuid.UUID) (*entities.Membership, error) {
					return &entities.Membership{ID: id, UserID: uuid.New(), AcademicUnitID: &unitID, Role: "student"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockMembershipRepository) {},
			wantErr:   true,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockMembershipRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Membership, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockMembershipRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewMembershipService(mockRepo, mock.NewMockLogger())
			result, err := svc.GetMembership(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestMembershipService_ExpireMembership(t *testing.T) {
	validID := uuid.New()
	unitID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockMembershipRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - expires membership",
			id:   validID.String(),
			setupMock: func(m *mock.MockMembershipRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Membership, error) {
					return &entities.Membership{ID: validID, UserID: uuid.New(), AcademicUnitID: &unitID, Role: "student", IsActive: true}, nil
				}
				m.UpdateFn = func(_ context.Context, mem *entities.Membership) error {
					assert.False(t, mem.IsActive)
					assert.NotNil(t, mem.WithdrawnAt)
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockMembershipRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Membership, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockMembershipRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewMembershipService(mockRepo, mock.NewMockLogger())
			result, err := svc.ExpireMembership(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.False(t, result.IsActive)
			}
		})
	}
}

func TestMembershipService_DeleteMembership(t *testing.T) {
	validID := uuid.New()
	unitID := uuid.New()

	tests := []struct {
		name      string
		id        string
		setupMock func(m *mock.MockMembershipRepository)
		wantErr   bool
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockMembershipRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Membership, error) {
					return &entities.Membership{ID: validID, AcademicUnitID: &unitID}, nil
				}
				m.DeleteFn = func(_ context.Context, _ uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockMembershipRepository) {},
			wantErr:   true,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockMembershipRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Membership, error) { return nil, nil }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockMembershipRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewMembershipService(mockRepo, mock.NewMockLogger())
			err := svc.DeleteMembership(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
