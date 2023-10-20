package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/models"
	"github.com/iNDicat0r/company/internal/app/services"
	"github.com/stretchr/testify/assert"
)

func TestNewCompanyHandler(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		companyService services.CompanyGetCreateUpdateDeleter
		expErr         string
	}{
		"no company service": {
			companyService: nil,
			expErr:         "company service is nil",
		},
		"success": {
			companyService: &mockCompanyService{},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, err := NewCompanyHandler(tt.companyService)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Nil(t, h)
			} else {
				assert.NotNil(t, h)
			}
		})
	}
}

func TestHandleGetCompany(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		companyService services.CompanyGetCreateUpdateDeleter
		params         gin.Params
		responseStatus int
		responseBody   string
		expErr         string
	}{
		"invalid company id": {
			companyService: &mockCompanyService{},
			params: gin.Params{gin.Param{
				Key:   "companyID",
				Value: "invalidid2738",
			}},
			responseStatus: http.StatusBadRequest,
			responseBody:   "{\"error\":\"invalid UUID length: 13\"}",
		},
		"internal service error": {
			companyService: &mockCompanyService{err: errors.New("internal error")},
			params: gin.Params{gin.Param{
				Key:   "companyID",
				Value: "ca8fc620-509a-40ac-8cc0-525c37c9c4b9",
			}},
			responseStatus: http.StatusInternalServerError,
			responseBody:   "{\"error\":\"internal error\"}",
		},
		"success": {
			companyService: &mockCompanyService{singleCompany: models.Company{
				Description: "description 1",
			}},
			params: gin.Params{gin.Param{
				Key:   "companyID",
				Value: "ca8fc620-509a-40ac-8cc0-525c37c9c4b9",
			}},
			responseStatus: http.StatusOK,
			responseBody:   "{\"ID\":\"00000000-0000-0000-0000-000000000000\",\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"Name\":\"\",\"Description\":\"description 1\",\"EmployeesAmount\":0,\"Registered\":false,\"Type\":\"\",\"UserID\":\"00000000-0000-0000-0000-000000000000\"}",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = append(c.Params, tt.params...)
			handler, _ := NewCompanyHandler(tt.companyService)
			handler.HandleGetCompany(c)
			assert.Equal(t, tt.responseStatus, c.Writer.Status())
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandleCreateCompany(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		companyService   services.CompanyGetCreateUpdateDeleter
		params           gin.Params
		setUserIDContext string
		requestBody      string
		responseStatus   int
		responseBody     string
		expErr           string
	}{
		"invalid request": {
			companyService:   &mockCompanyService{},
			responseStatus:   http.StatusBadRequest,
			requestBody:      "{",
			responseBody:     "{\"error\":\"unexpected EOF\"}",
			setUserIDContext: "ca8fc620-509a-40ac-8cc0-525c37c9c4b9",
		},
		"success": {
			companyService:   &mockCompanyService{},
			responseStatus:   http.StatusCreated,
			requestBody:      `{"name":"company1"}`,
			responseBody:     "{\"ID\":\"00000000-0000-0000-0000-000000000000\",\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"Name\":\"\",\"Description\":\"\",\"EmployeesAmount\":0,\"Registered\":false,\"Type\":\"\",\"UserID\":\"00000000-0000-0000-0000-000000000000\"}",
			setUserIDContext: "ca8fc620-509a-40ac-8cc0-525c37c9c4b9",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("userID", tt.setUserIDContext)
			c.Params = append(c.Params, tt.params...)
			req, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler, _ := NewCompanyHandler(tt.companyService)
			handler.HandleCreateCompany(c)
			assert.Equal(t, tt.responseStatus, c.Writer.Status())
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

type mockCompanyService struct {
	singleCompany models.Company
	err           error
}

func (m *mockCompanyService) Get(_ context.Context, _ uuid.UUID) (models.Company, error) {
	return m.singleCompany, m.err
}

func (m *mockCompanyService) Create(_ context.Context, _ uuid.UUID, _ services.CreateUpdateCompanyPayload) (models.Company, error) {
	return m.singleCompany, m.err
}

func (m *mockCompanyService) Update(_ context.Context, _ uuid.UUID, _ services.CreateUpdateCompanyPayload) (models.Company, error) {
	return m.singleCompany, m.err
}

func (m *mockCompanyService) Delete(_ context.Context, _, _ uuid.UUID) error {
	return m.err
}
