package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/middleware"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	sharedrepo "github.com/EduGoGroup/edugo-shared/repository"

	"github.com/EduGoGroup/edugo-api-admin-new/test/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestRouter creates a gin router with error handler middleware for handler tests.
func newTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.ErrorHandler(mock.NewMockLogger()))
	return r
}

func TestSchoolHandler_CreateSchool(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		setupMock  func(m *mock.MockSchoolService)
		wantStatus int
	}{
		{
			name: "success - returns 201",
			body: dto.CreateSchoolRequest{Name: "Test School", Code: "TST001"},
			setupMock: func(m *mock.MockSchoolService) {
				m.CreateSchoolFn = func(_ context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
					return &dto.SchoolResponse{
						ID: uuid.New().String(), Name: req.Name, Code: req.Code,
						IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
					}, nil
				}
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "error - invalid body returns 400",
			body:       "invalid json",
			setupMock:  func(_ *mock.MockSchoolService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - service returns validation error",
			body: dto.CreateSchoolRequest{Name: "AB", Code: "TST"},
			setupMock: func(m *mock.MockSchoolService) {
				m.CreateSchoolFn = func(_ context.Context, _ dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
					return nil, errors.NewValidationError("name too short")
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - service returns already exists",
			body: dto.CreateSchoolRequest{Name: "Test School", Code: "DUP"},
			setupMock: func(m *mock.MockSchoolService) {
				m.CreateSchoolFn = func(_ context.Context, _ dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
					return nil, errors.NewAlreadyExistsError("school")
				}
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSchoolService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSchoolHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.POST("/schools", h.CreateSchool)

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/schools", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSchoolHandler_GetSchool(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mock.MockSchoolService)
		wantStatus int
	}{
		{
			name: "success - returns 200",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSchoolService) {
				m.GetSchoolFn = func(_ context.Context, _ string) (*dto.SchoolResponse, error) {
					return &dto.SchoolResponse{ID: uuid.New().String(), Name: "School"}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error - not found returns 404",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSchoolService) {
				m.GetSchoolFn = func(_ context.Context, _ string) (*dto.SchoolResponse, error) {
					return nil, errors.NewNotFoundError("school")
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "error - invalid ID returns 400",
			id:   "bad-id",
			setupMock: func(m *mock.MockSchoolService) {
				m.GetSchoolFn = func(_ context.Context, _ string) (*dto.SchoolResponse, error) {
					return nil, errors.NewValidationError("invalid school ID")
				}
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSchoolService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSchoolHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.GET("/schools/:id", h.GetSchool)

			req, _ := http.NewRequest(http.MethodGet, "/schools/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSchoolHandler_ListSchools(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(m *mock.MockSchoolService)
		wantStatus int
	}{
		{
			name: "success - returns 200",
			setupMock: func(m *mock.MockSchoolService) {
				m.ListSchoolsFn = func(_ context.Context, _ sharedrepo.ListFilters) ([]dto.SchoolResponse, error) {
					return []dto.SchoolResponse{{ID: uuid.New().String(), Name: "School 1"}}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error - database error returns 500",
			setupMock: func(m *mock.MockSchoolService) {
				m.ListSchoolsFn = func(_ context.Context, _ sharedrepo.ListFilters) ([]dto.SchoolResponse, error) {
					return nil, errors.NewDatabaseError("list", nil)
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSchoolService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSchoolHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.GET("/schools", h.ListSchools)

			req, _ := http.NewRequest(http.MethodGet, "/schools", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSchoolHandler_DeleteSchool(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setupMock  func(m *mock.MockSchoolService)
		wantStatus int
	}{
		{
			name: "success - returns 204",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSchoolService) {
				m.DeleteSchoolFn = func(_ context.Context, _ string) error { return nil }
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "error - not found returns 404",
			id:   uuid.New().String(),
			setupMock: func(m *mock.MockSchoolService) {
				m.DeleteSchoolFn = func(_ context.Context, _ string) error {
					return errors.NewNotFoundError("school")
				}
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSchoolService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSchoolHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.DELETE("/schools/:id", h.DeleteSchool)

			req, _ := http.NewRequest(http.MethodDelete, "/schools/"+tt.id, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSchoolHandler_UpdateSchool(t *testing.T) {
	newName := "Updated"

	tests := []struct {
		name       string
		id         string
		body       interface{}
		setupMock  func(m *mock.MockSchoolService)
		wantStatus int
	}{
		{
			name: "success - returns 200",
			id:   uuid.New().String(),
			body: dto.UpdateSchoolRequest{Name: &newName},
			setupMock: func(m *mock.MockSchoolService) {
				m.UpdateSchoolFn = func(_ context.Context, _ string, _ dto.UpdateSchoolRequest) (*dto.SchoolResponse, error) {
					return &dto.SchoolResponse{ID: uuid.New().String(), Name: newName}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - invalid body returns 400",
			id:         uuid.New().String(),
			body:       "invalid",
			setupMock:  func(_ *mock.MockSchoolService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSchoolService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSchoolHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.PUT("/schools/:id", h.UpdateSchool)

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPut, "/schools/"+tt.id, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestSchoolHandler_GetSchoolByCode(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		setupMock  func(m *mock.MockSchoolService)
		wantStatus int
	}{
		{
			name: "success",
			code: "SCH001",
			setupMock: func(m *mock.MockSchoolService) {
				m.GetSchoolByCodeFn = func(_ context.Context, _ string) (*dto.SchoolResponse, error) {
					return &dto.SchoolResponse{ID: uuid.New().String(), Code: "SCH001"}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error - not found",
			code: "NOTEXIST",
			setupMock: func(m *mock.MockSchoolService) {
				m.GetSchoolByCodeFn = func(_ context.Context, _ string) (*dto.SchoolResponse, error) {
					return nil, errors.NewNotFoundError("school")
				}
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mock.MockSchoolService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			h := handler.NewSchoolHandler(mockSvc, mock.NewMockLogger())
			r := newTestRouter()
			r.GET("/schools/code/:code", h.GetSchoolByCode)

			req, _ := http.NewRequest(http.MethodGet, "/schools/code/"+tt.code, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
