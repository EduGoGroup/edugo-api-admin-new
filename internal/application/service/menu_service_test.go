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

func TestMenuService_GetFullMenu(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(r *mock.MockResourceRepository, rs *mock.MockResourceScreenRepository)
		wantErr   bool
	}{
		{
			name: "success - returns full menu",
			setupMock: func(r *mock.MockResourceRepository, _ *mock.MockResourceScreenRepository) {
				r.FindMenuVisibleFn = func(_ context.Context) ([]*entities.Resource, error) {
					return []*entities.Resource{
						{ID: uuid.New(), Key: "dashboard", DisplayName: "Dashboard", Scope: "system", IsMenuVisible: true},
					}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "success - empty menu",
			setupMock: func(r *mock.MockResourceRepository, _ *mock.MockResourceScreenRepository) {
				r.FindMenuVisibleFn = func(_ context.Context) ([]*entities.Resource, error) {
					return []*entities.Resource{}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "error - database error",
			setupMock: func(r *mock.MockResourceRepository, _ *mock.MockResourceScreenRepository) {
				r.FindMenuVisibleFn = func(_ context.Context) ([]*entities.Resource, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceRepo := &mock.MockResourceRepository{}
			resourceScreenRepo := &mock.MockResourceScreenRepository{}
			if tt.setupMock != nil {
				tt.setupMock(resourceRepo, resourceScreenRepo)
			}

			svc := service.NewMenuService(resourceRepo, resourceScreenRepo, mock.NewMockLogger())
			result, err := svc.GetFullMenu(context.Background())

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestMenuService_GetMenuForUser(t *testing.T) {
	parentID := uuid.New()

	tests := []struct {
		name        string
		permissions []string
		setupMock   func(r *mock.MockResourceRepository, rs *mock.MockResourceScreenRepository)
		wantErr     bool
		wantItems   int
	}{
		{
			name:        "success - no permissions returns empty menu",
			permissions: []string{},
			setupMock:   func(_ *mock.MockResourceRepository, _ *mock.MockResourceScreenRepository) {},
			wantErr:     false,
			wantItems:   0,
		},
		{
			name:        "success - filters by permissions",
			permissions: []string{"schools:read", "schools:create"},
			setupMock: func(r *mock.MockResourceRepository, _ *mock.MockResourceScreenRepository) {
				r.FindMenuVisibleFn = func(_ context.Context) ([]*entities.Resource, error) {
					return []*entities.Resource{
						{ID: parentID, Key: "schools", DisplayName: "Schools", Scope: "system", IsMenuVisible: true},
						{ID: uuid.New(), Key: "users", DisplayName: "Users", Scope: "system", IsMenuVisible: true},
					}, nil
				}
			},
			wantErr:   false,
			wantItems: 1,
		},
		{
			name:        "error - database error",
			permissions: []string{"schools:read"},
			setupMock: func(r *mock.MockResourceRepository, _ *mock.MockResourceScreenRepository) {
				r.FindMenuVisibleFn = func(_ context.Context) ([]*entities.Resource, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceRepo := &mock.MockResourceRepository{}
			resourceScreenRepo := &mock.MockResourceScreenRepository{}
			if tt.setupMock != nil {
				tt.setupMock(resourceRepo, resourceScreenRepo)
			}

			svc := service.NewMenuService(resourceRepo, resourceScreenRepo, mock.NewMockLogger())
			result, err := svc.GetMenuForUser(context.Background(), tt.permissions)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Len(t, result.Items, tt.wantItems)
			}
		})
	}
}
