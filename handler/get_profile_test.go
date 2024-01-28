package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/hasbiasshidiq/simple-profile/testing-init"

	"github.com/golang/mock/gomock"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestGetProfile(t *testing.T, token string) (context echo.Context, rec *httptest.ResponseRecorder, mockRepository *repository.MockRepositoryInterface) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)

	req.Header.Set("Authorization", "Bearer "+token)

	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepository = repository.NewMockRepositoryInterface(mockCtrl)

	return context, rec, mockRepository

}

func TestGetProfile(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		profile := repository.Profile{ID: 1}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestGetProfile(t, token)

		mockRepository.EXPECT().GetProfileByID(gomock.Any()).Return(profile, nil).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.GetProfile(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("Forbidden", func(t *testing.T) {

		token := "INVALIDTOKEN"
		context, rec, mockRepository := setupTestGetProfile(t, token)

		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.GetProfile(context)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})
	t.Run("Profile not found", func(t *testing.T) {
		profile := repository.Profile{}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestGetProfile(t, token)

		mockRepository.EXPECT().GetProfileByID(gomock.Any()).Return(profile, sql.ErrNoRows).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.GetProfile(context)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

}
