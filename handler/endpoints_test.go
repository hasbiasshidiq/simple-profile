package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	createProfileSuccess = `{
		"phone_number" : "+6289627117",
		"full_name" : "Hasbi Asshidiq",
		"password" : "1n19s9H88@"
	}`
	createProfileInvalidCountryCode = `{
		"phone_number" : "+6589627117",
		"full_name" : "Hasbi Asshidiq",
		"password" : "1n19s9H88@"
	}`
	createProfileInvalidPasswordPattern = `{
		"phone_number" : "+6289627117",
		"full_name" : "Hasbi Asshidiq",
		"password" : "1n19s9H880"
	}`
)

func setupTestCreateProfile(t *testing.T, requestBody string) (context echo.Context, rec *httptest.ResponseRecorder, mockRepository *repository.MockRepositoryInterface) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepository = repository.NewMockRepositoryInterface(mockCtrl)

	return context, rec, mockRepository

}

func TestCreateProfile(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		context, rec, mockRepository := setupTestCreateProfile(t, createProfileSuccess)

		mockRepository.EXPECT().CreateProfile(gomock.Any()).Return(0, nil).Times(1)
		mockRepository.EXPECT().GetPhoneNumberExistence(gomock.Any()).Return(false, nil).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostRegister(context)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("Phone Number Exists", func(t *testing.T) {

		context, rec, mockRepository := setupTestCreateProfile(t, createProfileSuccess)

		mockRepository.EXPECT().GetPhoneNumberExistence(gomock.Any()).Return(true, nil).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostRegister(context)) {
			assert.Equal(t, http.StatusConflict, rec.Code)
		}
	})

	t.Run("Invalid Country Code", func(t *testing.T) {

		context, rec, mockRepository := setupTestCreateProfile(t, createProfileInvalidCountryCode)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostRegister(context)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Invalid Password Pattern", func(t *testing.T) {

		context, rec, mockRepository := setupTestCreateProfile(t, createProfileInvalidPasswordPattern)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostRegister(context)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

}
