package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/hasbiasshidiq/simple-profile/testing-init"

	"github.com/golang/mock/gomock"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestPutProfile(t *testing.T, token string, requestBody string) (context echo.Context, rec *httptest.ResponseRecorder, mockRepository *repository.MockRepositoryInterface) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/profile", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	req.Header.Set("Authorization", "Bearer "+token)

	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepository = repository.NewMockRepositoryInterface(mockCtrl)

	return context, rec, mockRepository

}

func TestPutProfile(t *testing.T) {
	var (
		updateProfileSuccess = `{
			"phone_number" : "+6289627117",
			"full_name" : "Mr Bill Brod"
		}`
		updatePhoneNumberOnly = `{
			"phone_number" : "+6289627117"
		}`
		updateNameOnly = `{
			"full_name" : "Mr Bill Brod"
		}`
		invalidPhoneNumber = `{
			"phone_number" : "08589627117"
		}`
	)

	t.Run("Success", func(t *testing.T) {
		profile := repository.Profile{ID: 1}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestPutProfile(t, token, updateProfileSuccess)

		mockRepository.EXPECT().GetPhoneNumberExistenceWithExcludedID(gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
		mockRepository.EXPECT().UpdateProfileByID(gomock.Any()).Return(nil).Times(1)
		mockRepository.EXPECT().GetProfileByID(gomock.Any()).Return(profile, nil).Times(1)

		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PutProfile(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Success Update Phone Number Only", func(t *testing.T) {
		profile := repository.Profile{ID: 1}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestPutProfile(t, token, updatePhoneNumberOnly)

		mockRepository.EXPECT().GetPhoneNumberExistenceWithExcludedID(gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
		mockRepository.EXPECT().UpdateProfileByID(gomock.Any()).Return(nil).Times(1)
		mockRepository.EXPECT().GetProfileByID(gomock.Any()).Return(profile, nil).Times(1)

		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PutProfile(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Success Update Name Only", func(t *testing.T) {
		profile := repository.Profile{ID: 1}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestPutProfile(t, token, updateNameOnly)

		mockRepository.EXPECT().GetPhoneNumberExistenceWithExcludedID(gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
		mockRepository.EXPECT().UpdateProfileByID(gomock.Any()).Return(nil).Times(1)
		mockRepository.EXPECT().GetProfileByID(gomock.Any()).Return(profile, nil).Times(1)

		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PutProfile(context)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {

		token := "INVALIDTOKEN"
		context, rec, mockRepository := setupTestPutProfile(t, token, updateProfileSuccess)

		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PutProfile(context)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("Duplicate Phone Number", func(t *testing.T) {
		profile := repository.Profile{ID: 1}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestPutProfile(t, token, updateProfileSuccess)

		mockRepository.EXPECT().GetPhoneNumberExistenceWithExcludedID(gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PutProfile(context)) {
			assert.Equal(t, http.StatusConflict, rec.Code)
		}
	})

	t.Run("Invalid Phone Number", func(t *testing.T) {
		profile := repository.Profile{ID: 1}

		token, _ := createToken(profile)
		context, rec, mockRepository := setupTestPutProfile(t, token, invalidPhoneNumber)

		mockServer := &Server{Repository: mockRepository}

		if assert.NoError(t, mockServer.PutProfile(context)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

}
