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

func TestAcademicUnitService_CreateUnit(t *testing.T) {
	validSchoolID := uuid.New()

	tests := []struct {
		name        string
		schoolID    string
		request     dto.CreateAcademicUnitRequest
		setupMock   func(ur *mock.MockAcademicUnitRepository, sr *mock.MockSchoolRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:     "success - creates unit",
			schoolID: validSchoolID.String(),
			request:  dto.CreateAcademicUnitRequest{Type: "grade", DisplayName: "Grade 1", Code: "G1"},
			setupMock: func(ur *mock.MockAcademicUnitRepository, sr *mock.MockSchoolRepository) {
				sr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: validSchoolID}, nil
				}
				ur.ExistsBySchoolIDAndCodeFn = func(_ context.Context, _ uuid.UUID, _ string) (bool, error) { return false, nil }
				ur.CreateFn = func(_ context.Context, _ *entities.AcademicUnit) error { return nil }
			},
			wantErr: false,
		},
		{
			name:     "error - invalid school ID",
			schoolID: "bad-id",
			request:  dto.CreateAcademicUnitRequest{Type: "grade", DisplayName: "Grade 1"},
			setupMock: func(_ *mock.MockAcademicUnitRepository, _ *mock.MockSchoolRepository) {},
			wantErr:     true,
			errContains: "invalid school ID",
		},
		{
			name:     "error - school not found",
			schoolID: validSchoolID.String(),
			request:  dto.CreateAcademicUnitRequest{Type: "grade", DisplayName: "Grade 1"},
			setupMock: func(_ *mock.MockAcademicUnitRepository, sr *mock.MockSchoolRepository) {
				sr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:     "error - duplicate code in school",
			schoolID: validSchoolID.String(),
			request:  dto.CreateAcademicUnitRequest{Type: "grade", DisplayName: "Grade 1", Code: "G1"},
			setupMock: func(ur *mock.MockAcademicUnitRepository, sr *mock.MockSchoolRepository) {
				sr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: validSchoolID}, nil
				}
				ur.ExistsBySchoolIDAndCodeFn = func(_ context.Context, _ uuid.UUID, _ string) (bool, error) { return true, nil }
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:     "error - invalid parent unit ID",
			schoolID: validSchoolID.String(),
			request: func() dto.CreateAcademicUnitRequest {
				badParent := "not-uuid"
				return dto.CreateAcademicUnitRequest{Type: "class", DisplayName: "Class A", ParentUnitID: &badParent}
			}(),
			setupMock: func(_ *mock.MockAcademicUnitRepository, sr *mock.MockSchoolRepository) {
				sr.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.School, error) {
					return &entities.School{ID: validSchoolID}, nil
				}
			},
			wantErr:     true,
			errContains: "invalid parent_unit_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitRepo := &mock.MockAcademicUnitRepository{}
			schoolRepo := &mock.MockSchoolRepository{}
			if tt.setupMock != nil {
				tt.setupMock(unitRepo, schoolRepo)
			}

			svc := service.NewAcademicUnitService(unitRepo, schoolRepo, mock.NewMockLogger())
			result, err := svc.CreateUnit(context.Background(), tt.schoolID, tt.request)

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

func TestAcademicUnitService_GetUnit(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockAcademicUnitRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - found",
			id:   validID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.FindByIDFn = func(_ context.Context, id uuid.UUID, _ bool) (*entities.AcademicUnit, error) {
					return &entities.AcademicUnit{ID: id, Name: "Unit", SchoolID: uuid.New()}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid UUID",
			id:        "bad",
			setupMock: func(_ *mock.MockAcademicUnitRepository) {},
			wantErr:   true,
			errContains: "invalid unit ID",
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID, _ bool) (*entities.AcademicUnit, error) {
					return nil, nil
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitRepo := &mock.MockAcademicUnitRepository{}
			if tt.setupMock != nil {
				tt.setupMock(unitRepo)
			}

			svc := service.NewAcademicUnitService(unitRepo, &mock.MockSchoolRepository{}, mock.NewMockLogger())
			result, err := svc.GetUnit(context.Background(), tt.id)

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

func TestAcademicUnitService_DeleteUnit(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockAcademicUnitRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - soft deletes unit",
			id:   validID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID, _ bool) (*entities.AcademicUnit, error) {
					return &entities.AcademicUnit{ID: validID, SchoolID: uuid.New()}, nil
				}
				m.SoftDeleteFn = func(_ context.Context, _ uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID, _ bool) (*entities.AcademicUnit, error) {
					return nil, nil
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitRepo := &mock.MockAcademicUnitRepository{}
			if tt.setupMock != nil {
				tt.setupMock(unitRepo)
			}

			svc := service.NewAcademicUnitService(unitRepo, &mock.MockSchoolRepository{}, mock.NewMockLogger())
			err := svc.DeleteUnit(context.Background(), tt.id)

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

func TestAcademicUnitService_ListUnitsBySchool(t *testing.T) {
	schoolID := uuid.New()

	tests := []struct {
		name      string
		schoolID  string
		setupMock func(m *mock.MockAcademicUnitRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name:     "success - returns units",
			schoolID: schoolID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.FindBySchoolIDFn = func(_ context.Context, _ uuid.UUID, _ bool) ([]*entities.AcademicUnit, error) {
					return []*entities.AcademicUnit{
						{ID: uuid.New(), SchoolID: schoolID, Name: "Unit 1"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:      "error - invalid school ID",
			schoolID:  "bad",
			setupMock: func(_ *mock.MockAcademicUnitRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitRepo := &mock.MockAcademicUnitRepository{}
			if tt.setupMock != nil {
				tt.setupMock(unitRepo)
			}

			svc := service.NewAcademicUnitService(unitRepo, &mock.MockSchoolRepository{}, mock.NewMockLogger())
			result, err := svc.ListUnitsBySchool(context.Background(), tt.schoolID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestAcademicUnitService_ListUnitsByType(t *testing.T) {
	schoolID := uuid.New()

	tests := []struct {
		name        string
		schoolID    string
		unitType    string
		setupMock   func(m *mock.MockAcademicUnitRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:     "success",
			schoolID: schoolID.String(),
			unitType: "grade",
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.FindByTypeFn = func(_ context.Context, _ uuid.UUID, _ string, _ bool) ([]*entities.AcademicUnit, error) {
					return []*entities.AcademicUnit{}, nil
				}
			},
			wantErr: false,
		},
		{
			name:     "error - empty type",
			schoolID: schoolID.String(),
			unitType: "",
			setupMock: func(_ *mock.MockAcademicUnitRepository) {},
			wantErr:     true,
			errContains: "type query parameter is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitRepo := &mock.MockAcademicUnitRepository{}
			if tt.setupMock != nil {
				tt.setupMock(unitRepo)
			}

			svc := service.NewAcademicUnitService(unitRepo, &mock.MockSchoolRepository{}, mock.NewMockLogger())
			_, err := svc.ListUnitsByType(context.Background(), tt.schoolID, tt.unitType)

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

func TestAcademicUnitService_RestoreUnit(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockAcademicUnitRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - restores unit",
			id:   validID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.RestoreFn = func(_ context.Context, _ uuid.UUID) error { return nil }
				m.FindByIDFn = func(_ context.Context, id uuid.UUID, _ bool) (*entities.AcademicUnit, error) {
					return &entities.AcademicUnit{ID: id, SchoolID: uuid.New(), Name: "Restored"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockAcademicUnitRepository) {},
			wantErr:   true,
		},
		{
			name: "error - restore fails",
			id:   validID.String(),
			setupMock: func(m *mock.MockAcademicUnitRepository) {
				m.RestoreFn = func(_ context.Context, _ uuid.UUID) error { return fmt.Errorf("db error") }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unitRepo := &mock.MockAcademicUnitRepository{}
			if tt.setupMock != nil {
				tt.setupMock(unitRepo)
			}

			svc := service.NewAcademicUnitService(unitRepo, &mock.MockSchoolRepository{}, mock.NewMockLogger())
			result, err := svc.RestoreUnit(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
