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

func TestGuardianService_CreateRelation(t *testing.T) {
	guardianID := uuid.New().String()
	studentID := uuid.New().String()

	tests := []struct {
		name        string
		request     dto.CreateGuardianRelationRequest
		createdBy   string
		setupMock   func(m *mock.MockGuardianRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:      "success",
			request:   dto.CreateGuardianRelationRequest{GuardianID: guardianID, StudentID: studentID, RelationshipType: "father"},
			createdBy: uuid.New().String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.ExistsActiveRelationFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, _ *entities.GuardianRelation) error { return nil }
			},
			wantErr: false,
		},
		{
			name:        "error - invalid guardian_id",
			request:     dto.CreateGuardianRelationRequest{GuardianID: "bad", StudentID: studentID, RelationshipType: "father"},
			setupMock:   func(_ *mock.MockGuardianRepository) {},
			wantErr:     true,
			errContains: "invalid guardian_id",
		},
		{
			name:        "error - invalid student_id",
			request:     dto.CreateGuardianRelationRequest{GuardianID: guardianID, StudentID: "bad", RelationshipType: "father"},
			setupMock:   func(_ *mock.MockGuardianRepository) {},
			wantErr:     true,
			errContains: "invalid student_id",
		},
		{
			name:        "error - empty relationship type",
			request:     dto.CreateGuardianRelationRequest{GuardianID: guardianID, StudentID: studentID, RelationshipType: ""},
			setupMock:   func(_ *mock.MockGuardianRepository) {},
			wantErr:     true,
			errContains: "relationship_type is required",
		},
		{
			name:    "error - relation already exists",
			request: dto.CreateGuardianRelationRequest{GuardianID: guardianID, StudentID: studentID, RelationshipType: "father"},
			setupMock: func(m *mock.MockGuardianRepository) {
				m.ExistsActiveRelationFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) { return true, nil }
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:    "error - database error on create",
			request: dto.CreateGuardianRelationRequest{GuardianID: guardianID, StudentID: studentID, RelationshipType: "mother"},
			setupMock: func(m *mock.MockGuardianRepository) {
				m.ExistsActiveRelationFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, _ *entities.GuardianRelation) error { return fmt.Errorf("db error") }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockGuardianRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewGuardianService(mockRepo, mock.NewMockLogger())
			result, err := svc.CreateRelation(context.Background(), tt.request, tt.createdBy)

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

func TestGuardianService_GetRelation(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockGuardianRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.GuardianRelation, error) {
					return &entities.GuardianRelation{ID: validID, GuardianID: uuid.New(), StudentID: uuid.New(), RelationshipType: "father"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockGuardianRepository) {},
			wantErr:   true,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.GuardianRelation, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockGuardianRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewGuardianService(mockRepo, mock.NewMockLogger())
			result, err := svc.GetRelation(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestGuardianService_DeleteRelation(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setupMock func(m *mock.MockGuardianRepository)
		wantErr bool
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.DeleteFn = func(_ context.Context, _ uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockGuardianRepository) {},
			wantErr:   true,
		},
		{
			name: "error - database error",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.DeleteFn = func(_ context.Context, _ uuid.UUID) error { return fmt.Errorf("db error") }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockGuardianRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewGuardianService(mockRepo, mock.NewMockLogger())
			err := svc.DeleteRelation(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGuardianService_GetGuardianRelations(t *testing.T) {
	guardianID := uuid.New()

	tests := []struct {
		name      string
		id        string
		setupMock func(m *mock.MockGuardianRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success",
			id:   guardianID.String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.FindByGuardianFn = func(_ context.Context, _ uuid.UUID) ([]*entities.GuardianRelation, error) {
					return []*entities.GuardianRelation{
						{ID: uuid.New(), GuardianID: guardianID, StudentID: uuid.New(), RelationshipType: "father"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockGuardianRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockGuardianRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewGuardianService(mockRepo, mock.NewMockLogger())
			result, err := svc.GetGuardianRelations(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestGuardianService_GetStudentGuardians(t *testing.T) {
	studentID := uuid.New()

	tests := []struct {
		name      string
		id        string
		setupMock func(m *mock.MockGuardianRepository)
		wantErr   bool
	}{
		{
			name: "success",
			id:   studentID.String(),
			setupMock: func(m *mock.MockGuardianRepository) {
				m.FindByStudentFn = func(_ context.Context, _ uuid.UUID) ([]*entities.GuardianRelation, error) {
					return []*entities.GuardianRelation{}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockGuardianRepository) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockGuardianRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewGuardianService(mockRepo, mock.NewMockLogger())
			_, err := svc.GetStudentGuardians(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
