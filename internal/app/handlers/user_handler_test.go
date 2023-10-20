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

func TestNewUserHandler(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		userService services.User
		expErr      string
	}{
		"no user service": {
			userService: nil,
			expErr:      "user service is nil",
		},
		"success": {
			userService: &services.UserService{},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, err := NewUserHandler(tt.userService)
			if tt.expErr != "" {
				assert.EqualError(t, err, tt.expErr)
				assert.Nil(t, h)
			} else {
				assert.NotNil(t, h)
			}
		})
	}
}

func TestHandleAuthenticate(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		userService    services.User
		responseStatus int
		requestBody    string
		responseBody   string
	}{
		"invalid body": {
			userService:    &mockUserService{},
			responseStatus: http.StatusBadRequest,
			requestBody:    "{",
			responseBody:   "{\"error\":\"unexpected EOF\"}",
		},
		"user service error": {
			userService:    &mockUserService{err: errors.New("internal error")},
			responseStatus: http.StatusInternalServerError,
			requestBody:    "{\"username\":\"hello\", \"password\":\"123\"}",
			responseBody:   "{\"error\":\"internal error\"}",
		},
		"success": {
			userService:    &mockUserService{jwt: "jwt123"},
			responseStatus: http.StatusOK,
			requestBody:    "{\"username\":\"hello\", \"password\":\"123\"}",
			responseBody:   "{\"token\":\"jwt123\"}",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("POST", "/auth/authenticate", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler, _ := NewUserHandler(tt.userService)
			handler.HandleAuthenticate(c)
			assert.Equal(t, tt.responseStatus, c.Writer.Status())
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandleIntrospect(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		userService      services.User
		responseStatus   int
		requestBody      string
		params           gin.Params
		responseBody     string
		setUserIDContext string
	}{
		"user service error": {
			userService:    &mockUserService{err: errors.New("internal error")},
			responseStatus: http.StatusInternalServerError,
			requestBody:    "{}",
			params: gin.Params{gin.Param{
				Key:   "userID",
				Value: "862dedcb-68c5-49f7-a94a-b7190499f16b",
			}},
			responseBody:     "{\"error\":\"internal error\"}",
			setUserIDContext: "862dedcb-68c5-49f7-a94a-b7190499f16b",
		},
		"success": {
			userService: &mockUserService{user: models.User{
				ID:       uuid.MustParse("862dedcb-68c5-49f7-a94a-b7190499f16b"),
				Name:     "user 1",
				Username: "userx",
				Password: "123",
			}},
			responseStatus: http.StatusOK,
			requestBody:    "{}",
			params: gin.Params{gin.Param{
				Key:   "userID",
				Value: "862dedcb-68c5-49f7-a94a-b7190499f16b",
			}},
			responseBody:     "{\"ID\":\"862dedcb-68c5-49f7-a94a-b7190499f16b\",\"CreatedAt\":\"0001-01-01T00:00:00Z\",\"UpdatedAt\":\"0001-01-01T00:00:00Z\",\"DeletedAt\":null,\"Name\":\"user 1\",\"Username\":\"userx\",\"Companies\":null}",
			setUserIDContext: "862dedcb-68c5-49f7-a94a-b7190499f16b",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("userID", tt.setUserIDContext)
			req, _ := http.NewRequest("POST", "/auth/introspect", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Type", "application/json")
			c.Request = req
			c.Params = append(c.Params, tt.params...)

			handler, _ := NewUserHandler(tt.userService)
			handler.HandleIntrospect(c)
			assert.Equal(t, tt.responseStatus, c.Writer.Status())
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

type mockUserService struct {
	userID uuid.UUID
	jwt    string
	user   models.User
	err    error
}

func (u *mockUserService) Save(_ context.Context, _, _, _ string) (uuid.UUID, error) {
	return u.userID, u.err
}

func (u *mockUserService) Authenticate(_ context.Context, _, _ string) (string, error) {
	return u.jwt, u.err
}

func (u *mockUserService) GetUser(_ context.Context, _ uuid.UUID) (models.User, error) {
	return u.user, u.err
}
