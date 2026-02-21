package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultSchoolDefaults = config.SchoolDefaults{
	Country:          "CO",
	SubscriptionTier: "free",
	MaxTeachers:      50,
	MaxStudents:      500,
}

func TestSchoolService_CreateSchool(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.CreateSchoolRequest
		setupMock   func(m *mock.MockSchoolRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - creates school with valid data",
			request: dto.CreateSchoolRequest{
				Name: "Test School",
				Code: "TST001",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, _ *entities.School) error { return nil }
			},
			wantErr: false,
		},
		{
			name: "success - applies defaults when fields are empty",
			request: dto.CreateSchoolRequest{
				Name: "Test School",
				Code: "TST002",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, s *entities.School) error {
					assert.Equal(t, "CO", s.Country)
					assert.Equal(t, "free", s.SubscriptionTier)
					assert.Equal(t, 50, s.MaxTeachers)
					assert.Equal(t, 500, s.MaxStudents)
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "error - duplicate code",
			request: dto.CreateSchoolRequest{
				Name: "Test School",
				Code: "DUP001",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return true, nil }
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name: "error - name too short",
			request: dto.CreateSchoolRequest{
				Name: "AB",
				Code: "TST003",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
			},
			wantErr:     true,
			errContains: "name must be at least 3 characters",
		},
		{
			name: "error - code too short",
			request: dto.CreateSchoolRequest{
				Name: "Test School",
				Code: "AB",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
			},
			wantErr:     true,
			errContains: "code must be at least 3 characters",
		},
		{
			name: "error - empty name",
			request: dto.CreateSchoolRequest{
				Name: "",
				Code: "TST004",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
			},
			wantErr:     true,
			errContains: "name must be at least 3 characters",
		},
		{
			name: "error - database error on exists check",
			request: dto.CreateSchoolRequest{
				Name: "Test School",
				Code: "TST005",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) {
					return false, fmt.Errorf("db connection error")
				}
			},
			wantErr:     true,
			errContains: "database error",
		},
		{
			name: "error - database error on create",
			request: dto.CreateSchoolRequest{
				Name: "Test School",
				Code: "TST006",
			},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ExistsByCodeFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, _ *entities.School) error {
					return fmt.Errorf("insert error")
				}
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSchoolService(mockRepo, mock.NewMockLogger(), defaultSchoolDefaults)
			result, err := svc.CreateSchool(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.Code, result.Code)
			}
		})
	}
}

func TestSchoolService_GetSchool(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockSchoolRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - school found",
			id:   validID.String(),
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, id uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: id, Name: "Found School", Code: "FND"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:        "error - invalid UUID",
			id:          "not-a-uuid",
			setupMock:   func(m *mock.MockSchoolRepository) {},
			wantErr:     true,
			errContains: "invalid school ID",
		},
		{
			name: "error - school not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return nil, nil
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "error - database error",
			id:   validID.String(),
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return nil, fmt.Errorf("connection refused")
				}
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSchoolService(mockRepo, mock.NewMockLogger(), defaultSchoolDefaults)
			result, err := svc.GetSchool(context.Background(), tt.id)

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

func TestSchoolService_ListSchools(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(m *mock.MockSchoolRepository)
		wantErr     bool
		wantCount   int
	}{
		{
			name: "success - returns list",
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ListFn = func(_ context.Context, _ repository.ListFilters) ([]*entities.School, error) {
					return []*entities.School{
						{ID: uuid.New(), Name: "School 1", Code: "S1"},
						{ID: uuid.New(), Name: "School 2", Code: "S2"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "success - empty list",
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ListFn = func(_ context.Context, _ repository.ListFilters) ([]*entities.School, error) {
					return []*entities.School{}, nil
				}
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "error - database error",
			setupMock: func(m *mock.MockSchoolRepository) {
				m.ListFn = func(_ context.Context, _ repository.ListFilters) ([]*entities.School, error) {
					return nil, fmt.Errorf("timeout")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSchoolService(mockRepo, mock.NewMockLogger(), defaultSchoolDefaults)
			result, err := svc.ListSchools(context.Background())

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestSchoolService_DeleteSchool(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockSchoolRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - deletes school",
			id:   validID.String(),
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: validID, Name: "School"}, nil
				}
				m.DeleteFn = func(_ context.Context, _ uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name:      "error - invalid UUID",
			id:        "bad-id",
			setupMock: func(_ *mock.MockSchoolRepository) {},
			wantErr:   true,
			errContains: "invalid school ID",
		},
		{
			name: "error - school not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return nil, nil
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSchoolService(mockRepo, mock.NewMockLogger(), defaultSchoolDefaults)
			err := svc.DeleteSchool(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSchoolService_UpdateSchool(t *testing.T) {
	validID := uuid.New()
	newName := "Updated School"
	shortName := "AB"

	tests := []struct {
		name        string
		id          string
		request     dto.UpdateSchoolRequest
		setupMock   func(m *mock.MockSchoolRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - updates name",
			id:      validID.String(),
			request: dto.UpdateSchoolRequest{Name: &newName},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: validID, Name: "Old Name", Code: "OLD"}, nil
				}
				m.UpdateFn = func(_ context.Context, _ *entities.School) error { return nil }
			},
			wantErr: false,
		},
		{
			name:    "error - name too short",
			id:      validID.String(),
			request: dto.UpdateSchoolRequest{Name: &shortName},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: validID, Name: "Old Name", Code: "OLD"}, nil
				}
			},
			wantErr:     true,
			errContains: "name must be at least 3 characters",
		},
		{
			name:      "error - invalid UUID",
			id:        "bad-id",
			request:   dto.UpdateSchoolRequest{Name: &newName},
			setupMock: func(_ *mock.MockSchoolRepository) {},
			wantErr:   true,
			errContains: "invalid school ID",
		},
		{
			name:    "error - school not found",
			id:      validID.String(),
			request: dto.UpdateSchoolRequest{Name: &newName},
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return nil, nil
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSchoolService(mockRepo, mock.NewMockLogger(), defaultSchoolDefaults)
			result, err := svc.UpdateSchool(context.Background(), tt.id, tt.request)

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

func TestSchoolService_GetSchoolByCode(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		setupMock   func(m *mock.MockSchoolRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - found by code",
			code: "SCH001",
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByCodeFn = func(_ context.Context, _ string) (*entities.School, error) {
					return &entities.School{ID: uuid.New(), Name: "School", Code: "SCH001"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			code: "NOTFOUND",
			setupMock: func(m *mock.MockSchoolRepository) {
				m.FindByCodeFn = func(_ context.Context, _ string) (*entities.School, error) {
					return nil, nil
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSchoolService(mockRepo, mock.NewMockLogger(), defaultSchoolDefaults)
			result, err := svc.GetSchoolByCode(context.Background(), tt.code)

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
