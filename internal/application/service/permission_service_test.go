package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissionService_ListPermissions(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(m *mock.MockPermissionRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success",
			setupMock: func(m *mock.MockPermissionRepository) {
				m.FindAllFn = func(_ context.Context) ([]*entities.Permission, error) {
					return []*entities.Permission{
						{ID: uuid.New(), Name: "schools:read", DisplayName: "Read Schools", ResourceID: uuid.New(), ResourceKey: "schools", Action: "read", Scope: "system"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "error - database error",
			setupMock: func(m *mock.MockPermissionRepository) {
				m.FindAllFn = func(_ context.Context) ([]*entities.Permission, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockPermissionRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewPermissionService(mockRepo, mock.NewMockLogger())
			result, err := svc.ListPermissions(context.Background())

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Permissions, tt.wantCount)
			}
		})
	}
}

func TestPermissionService_GetPermission(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockPermissionRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockPermissionRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Permission, error) {
					return &entities.Permission{ID: validID, Name: "schools:read", DisplayName: "Read Schools", ResourceID: uuid.New(), ResourceKey: "schools", Action: "read", Scope: "system"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:        "error - invalid ID",
			id:          "bad",
			setupMock:   func(_ *mock.MockPermissionRepository) {},
			wantErr:     true,
			errContains: "invalid permission ID",
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockPermissionRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Permission, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockPermissionRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewPermissionService(mockRepo, mock.NewMockLogger())
			result, err := svc.GetPermission(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
