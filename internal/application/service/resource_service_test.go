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

func TestResourceService_ListResources(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(m *mock.MockResourceRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success - returns resources",
			setupMock: func(m *mock.MockResourceRepository) {
				m.FindAllFn = func(_ context.Context) ([]*entities.Resource, error) {
					return []*entities.Resource{
						{ID: uuid.New(), Key: "schools", DisplayName: "Schools", Scope: "system"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "error - database error",
			setupMock: func(m *mock.MockResourceRepository) {
				m.FindAllFn = func(_ context.Context) ([]*entities.Resource, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockResourceRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewResourceService(mockRepo, mock.NewMockLogger())
			result, err := svc.ListResources(context.Background())

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.wantCount, result.Total)
			}
		})
	}
}

func TestResourceService_GetResource(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockResourceRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockResourceRepository) {
				m.FindByIDFn = func(_ context.Context, id uuid.UUID) (*entities.Resource, error) {
					return &entities.Resource{ID: id, Key: "schools", DisplayName: "Schools", Scope: "system"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:        "error - invalid ID",
			id:          "bad",
			setupMock:   func(_ *mock.MockResourceRepository) {},
			wantErr:     true,
			errContains: "invalid resource ID",
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockResourceRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Resource, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockResourceRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewResourceService(mockRepo, mock.NewMockLogger())
			result, err := svc.GetResource(context.Background(), tt.id)

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

func TestResourceService_CreateResource(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.CreateResourceRequest
		setupMock   func(m *mock.MockResourceRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - creates resource",
			request: dto.CreateResourceRequest{
				Key: "new_resource", DisplayName: "New Resource", Scope: "system",
			},
			setupMock: func(m *mock.MockResourceRepository) {
				m.CreateFn = func(_ context.Context, _ *entities.Resource) error { return nil }
			},
			wantErr: false,
		},
		{
			name: "success - with parent",
			request: func() dto.CreateResourceRequest {
				pid := uuid.New().String()
				return dto.CreateResourceRequest{
					Key: "child", DisplayName: "Child", Scope: "system", ParentID: &pid,
				}
			}(),
			setupMock: func(m *mock.MockResourceRepository) {
				m.CreateFn = func(_ context.Context, _ *entities.Resource) error { return nil }
			},
			wantErr: false,
		},
		{
			name: "error - invalid parent_id",
			request: func() dto.CreateResourceRequest {
				pid := "bad"
				return dto.CreateResourceRequest{
					Key: "child", DisplayName: "Child", Scope: "system", ParentID: &pid,
				}
			}(),
			setupMock:   func(_ *mock.MockResourceRepository) {},
			wantErr:     true,
			errContains: "invalid parent_id",
		},
		{
			name: "error - database error",
			request: dto.CreateResourceRequest{
				Key: "res", DisplayName: "Res", Scope: "system",
			},
			setupMock: func(m *mock.MockResourceRepository) {
				m.CreateFn = func(_ context.Context, _ *entities.Resource) error { return fmt.Errorf("db error") }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockResourceRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewResourceService(mockRepo, mock.NewMockLogger())
			result, err := svc.CreateResource(context.Background(), tt.request)

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

func TestResourceService_UpdateResource(t *testing.T) {
	validID := uuid.New()
	newName := "Updated"

	tests := []struct {
		name        string
		id          string
		request     dto.UpdateResourceRequest
		setupMock   func(m *mock.MockResourceRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - updates resource",
			id:      validID.String(),
			request: dto.UpdateResourceRequest{DisplayName: &newName},
			setupMock: func(m *mock.MockResourceRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Resource, error) {
					return &entities.Resource{ID: validID, Key: "schools", DisplayName: "Schools", Scope: "system"}, nil
				}
				m.UpdateFn = func(_ context.Context, _ *entities.Resource) error { return nil }
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			request:   dto.UpdateResourceRequest{DisplayName: &newName},
			setupMock: func(_ *mock.MockResourceRepository) {},
			wantErr:   true,
		},
		{
			name:    "error - not found",
			id:      validID.String(),
			request: dto.UpdateResourceRequest{DisplayName: &newName},
			setupMock: func(m *mock.MockResourceRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Resource, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockResourceRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewResourceService(mockRepo, mock.NewMockLogger())
			result, err := svc.UpdateResource(context.Background(), tt.id, tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
