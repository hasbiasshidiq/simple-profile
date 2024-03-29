package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "github.com/hasbiasshidiq/simple-profile/testing-init"

	"github.com/golang/mock/gomock"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestLogin(t *testing.T, requestBody string) (context echo.Context, rec *httptest.ResponseRecorder, mockRepository *repository.MockRepositoryInterface) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepository = repository.NewMockRepositoryInterface(mockCtrl)

	return context, rec, mockRepository

}

func TestLogin(t *testing.T) {
	var (
		loginSuccess = `{
			"phone_number" : "+6289627117",
			"password" : "1n19s9H88@"
		}`
		loginAccountNotFound = `{
			"phone_number" : "+6289627117",
			"password" : "1n19s9H88@"
		}`
		loginInvalidPassword = `{
			"phone_number" : "+6289627117",
			"password" : "INVALIDPASS"
		}`
	)

	//The hashed password that is stored in the repository.
	hashedPassword := hashAndSalt([]byte("1n19s9H88@"))

	var profile repository.Profile = repository.Profile{
		ID:          1,
		FullName:    "Bakri",
		CountryCode: "+62",
		PhoneNumber: "89627117",
		Password:    hashedPassword,
		CreatedAt:   time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		context, rec, mockRepository := setupTestLogin(t, loginSuccess)

		mockRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any()).Return(profile, nil).Times(1)
		mockRepository.EXPECT().UpsertProfileMetaData(gomock.Any()).Return(0, nil).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostLogin(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Phone Number Exists", func(t *testing.T) {

		context, rec, mockRepository := setupTestCreateProfile(t, loginAccountNotFound)

		profile := repository.Profile{}

		mockRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any()).Return(profile, sql.ErrNoRows).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostLogin(context)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {

		context, rec, mockRepository := setupTestCreateProfile(t, loginInvalidPassword)

		mockRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any()).Return(profile, sql.ErrNoRows).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PostLogin(context)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

}
