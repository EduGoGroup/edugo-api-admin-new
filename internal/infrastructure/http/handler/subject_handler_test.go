package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"

	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
)

func TestSubjectHandler_CreateSubject(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMock  func(m *mock.MockSubjectService)
		wantStatus int
	}{
		{
			name: "success - returns 201",
			body: dto.CreateSubjectRequest{Name: "Mathematics"},
			setupMock: func(m *mock.MockSubjectService) {
				m.CreateSubjectFn = func(_ context.Context, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
					return &dto.SubjectResponse{ID: uuid.New().String(), Name: req.Name, IsActive: true}, nil
				}
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "error - invalid body",
			body:       "bad",
			setupMock:  func(_ *mock.MockSubjectService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - duplicate name",
			body: dto.CreateSubjectRequest{Name: "Math"},
			setupMock: func(m *mock.MockSubjectService) {
				m.CreateSubjectFn = func(_ context.Context, _ dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
					return nil, errors.NewAlreadyExistsError("subject")
				}
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSubjectService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSubjectHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.POST("/subjects", h.CreateSubject)

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/subjects", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSubjectHandler_GetSubject(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mock.MockSubjectService)
		wantStatus int
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSubjectService) {
				m.GetSubjectFn = func(_ context.Context, _ string) (*dto.SubjectResponse, error) {
					return &dto.SubjectResponse{ID: uuid.New().String(), Name: "Math"}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error - not found",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSubjectService) {
				m.GetSubjectFn = func(_ context.Context, _ string) (*dto.SubjectResponse, error) {
					return nil, errors.NewNotFoundError("subject")
				}
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSubjectService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSubjectHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.GET("/subjects/:id", h.GetSubject)

			req, _ := http.NewRequest(http.MethodGet, "/subjects/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSubjectHandler_ListSubjects(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(m *mock.MockSubjectService)
		wantStatus int
	}{
		{
			name: "success",
			setupMock: func(m *mock.MockSubjectService) {
				m.ListSubjectsFn = func(_ context.Context, _ sharedrepo.ListFilters) ([]dto.SubjectResponse, error) {
					return []dto.SubjectResponse{}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSubjectService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSubjectHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.GET("/subjects", h.ListSubjects)

			req, _ := http.NewRequest(http.MethodGet, "/subjects", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSubjectHandler_DeleteSubject(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mock.MockSubjectService)
		wantStatus int
	}{
		{
			name: "success",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSubjectService) {
				m.DeleteSubjectFn = func(_ context.Context, _ string) error { return nil }
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "error - not found",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSubjectService) {
				m.DeleteSubjectFn = func(_ context.Context, _ string) error {
					return errors.NewNotFoundError("subject")
				}
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSubjectService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSubjectHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.DELETE("/subjects/:id", h.DeleteSubject)

			req, _ := http.NewRequest(http.MethodDelete, "/subjects/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
