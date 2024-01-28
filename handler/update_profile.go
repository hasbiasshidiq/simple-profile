package handler

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/hasbiasshidiq/simple-profile/generated"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
)

type (
	UpdateProfileValidator struct {
		FullName    *string `json:"full_name" validate:"omitempty,min=3,max=60"`
		PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=13,indonesiaCountryCodePrefix"`
	}
)

func (s *Server) PutProfile(ctx echo.Context) error {
	// Extract the token from the Authorization header
	token, err := extractToken(ctx)
	if err != nil {
		responsePayload := generated.GeneralErrorResponse{Message: "Invalid Token"}
		return ctx.JSON(http.StatusForbidden, responsePayload)
	}

	userID, err := extractUserIDFromToken(token)
	if err != nil {
		responsePayload := generated.GeneralErrorResponse{Message: "Invalid Token"}
		return ctx.JSON(http.StatusForbidden, responsePayload)
	}

	var request generated.UpdateProfileRequest

	customValidator := validator.New()
	customValidator.RegisterValidation("indonesiaCountryCodePrefix", validatePhoneWithPrefix)
	customValidator.RegisterValidation("strongPassword", validateStrongPassword)

	ctx.Echo().Validator = &CustomValidator{validator: customValidator}

	err = ctx.Bind(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if request.FullName == nil && request.PhoneNumber == nil {
		responsePayload := generated.GeneralErrorResponse{Message: "full_name or phone_number should be filled"}
		return ctx.JSON(http.StatusBadRequest, responsePayload)
	}
	updateProfileValidator := UpdateProfileValidator{}
	if request.FullName != nil {
		updateProfileValidator.FullName = request.FullName
	}
	if request.PhoneNumber != nil {
		updateProfileValidator.PhoneNumber = request.PhoneNumber
	}

	if err = ctx.Validate(updateProfileValidator); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.GeneralErrorResponse{Message: err.Error()})
	}

	if request.PhoneNumber != nil {
		localPhoneNumber := (*request.PhoneNumber)[3:]

		isExist, err := s.Repository.GetPhoneNumberExistenceWithExcludedID(localPhoneNumber, userID)
		if err != nil {
			return err
		}
		if isExist {
			responsePayload := generated.GeneralErrorResponse{Message: "Phone Number Already Exist"}
			return ctx.JSON(http.StatusConflict, responsePayload)
		}
	}

	profile := repository.Profile{ID: uint64(userID)}
	if request.PhoneNumber != nil {
		localPhoneNumber := (*request.PhoneNumber)[3:]
		profile.PhoneNumber = localPhoneNumber
	}
	if request.FullName != nil {
		profile.FullName = *request.FullName
	}

	err = s.Repository.UpdateProfileByID(profile)

	if err != nil {
		responsePayload := generated.GeneralErrorResponse{Message: "Can't update profile"}
		return ctx.JSON(http.StatusConflict, responsePayload)
	}

	profile, err = s.Repository.GetProfileByID(userID)
	if err != nil {
		responsePayload := generated.GeneralErrorResponse{Message: "Profile not found"}
		return ctx.JSON(http.StatusNotFound, responsePayload)
	}

	resp := generated.UpdateProfileResponse{
		FullName:    profile.FullName,
		PhoneNumber: profile.CountryCode + profile.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}
