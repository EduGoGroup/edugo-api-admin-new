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

func TestSubjectService_CreateSubject(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.CreateSubjectRequest
		setupMock   func(m *mock.MockSubjectRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success",
			request: dto.CreateSubjectRequest{Name: "Mathematics", Description: "Math subject"},
			setupMock: func(m *mock.MockSubjectRepository) {
				m.ExistsByNameFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, _ *entities.Subject) error { return nil }
			},
			wantErr: false,
		},
		{
			name:        "error - name too short",
			request:     dto.CreateSubjectRequest{Name: "M"},
			setupMock:   func(_ *mock.MockSubjectRepository) {},
			wantErr:     true,
			errContains: "name must be at least 2 characters",
		},
		{
			name:        "error - empty name",
			request:     dto.CreateSubjectRequest{Name: ""},
			setupMock:   func(_ *mock.MockSubjectRepository) {},
			wantErr:     true,
			errContains: "name must be at least 2 characters",
		},
		{
			name:    "error - duplicate name",
			request: dto.CreateSubjectRequest{Name: "Mathematics"},
			setupMock: func(m *mock.MockSubjectRepository) {
				m.ExistsByNameFn = func(_ context.Context, _ string) (bool, error) { return true, nil }
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:    "error - database error on create",
			request: dto.CreateSubjectRequest{Name: "Physics"},
			setupMock: func(m *mock.MockSubjectRepository) {
				m.ExistsByNameFn = func(_ context.Context, _ string) (bool, error) { return false, nil }
				m.CreateFn = func(_ context.Context, _ *entities.Subject) error { return fmt.Errorf("db error") }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSubjectRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSubjectService(mockRepo, mock.NewMockLogger())
			result, err := svc.CreateSubject(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
			}
		})
	}
}

func TestSubjectService_GetSubject(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mock.MockSubjectRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockSubjectRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Subject, error) {
					return &entities.Subject{ID: validID, Name: "Math"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockSubjectRepository) {},
			wantErr:   true,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockSubjectRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Subject, error) { return nil, nil }
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSubjectRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSubjectService(mockRepo, mock.NewMockLogger())
			result, err := svc.GetSubject(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestSubjectService_ListSubjects(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(m *mock.MockSubjectRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name: "success",
			setupMock: func(m *mock.MockSubjectRepository) {
				m.ListFn = func(_ context.Context) ([]*entities.Subject, error) {
					return []*entities.Subject{
						{ID: uuid.New(), Name: "Math"},
						{ID: uuid.New(), Name: "Science"},
					}, nil
				}
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name: "error - database error",
			setupMock: func(m *mock.MockSubjectRepository) {
				m.ListFn = func(_ context.Context) ([]*entities.Subject, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSubjectRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSubjectService(mockRepo, mock.NewMockLogger())
			result, err := svc.ListSubjects(context.Background())

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}
		})
	}
}

func TestSubjectService_DeleteSubject(t *testing.T) {
	validID := uuid.New()

	tests := []struct {
		name      string
		id        string
		setupMock func(m *mock.MockSubjectRepository)
		wantErr   bool
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(m *mock.MockSubjectRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Subject, error) {
					return &entities.Subject{ID: validID, Name: "Math"}, nil
				}
				m.DeleteFn = func(_ context.Context, _ uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name:      "error - invalid ID",
			id:        "bad",
			setupMock: func(_ *mock.MockSubjectRepository) {},
			wantErr:   true,
		},
		{
			name: "error - not found",
			id:   validID.String(),
			setupMock: func(m *mock.MockSubjectRepository) {
				m.FindByIDFn = func(_ context.Context, _ uuid.UUID) (*entities.Subject, error) { return nil, nil }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockSubjectRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			svc := service.NewSubjectService(mockRepo, mock.NewMockLogger())
			err := svc.DeleteSubject(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
