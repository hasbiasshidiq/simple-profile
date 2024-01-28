package handler

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/hasbiasshidiq/simple-profile/generated"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
)

type (
	CreateProfileValidator struct {
		FullName    string `json:"full_name" validate:"required,min=3,max=60"`
		PhoneNumber string `json:"phone_number" validate:"required,min=10,max=13,indonesiaCountryCodePrefix"`
		Password    string `json:"password" validate:"required,min=6,max=64,strongPassword"`
	}
)

func (s *Server) PostProfile(ctx echo.Context) error {

	var request generated.CreateProfileRequest

	customValidator := validator.New()
	customValidator.RegisterValidation("indonesiaCountryCodePrefix", validatePhoneWithPrefix)
	customValidator.RegisterValidation("strongPassword", validateStrongPassword)

	ctx.Echo().Validator = &CustomValidator{validator: customValidator}

	err := ctx.Bind(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	CreateProfileValidator := CreateProfileValidator{
		FullName:    request.FullName,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	}
	if err = ctx.Validate(CreateProfileValidator); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.GeneralErrorResponse{Message: err.Error()})
	}

	countryCode := request.PhoneNumber[:3]
	localPhoneNumber := request.PhoneNumber[3:]

	isExist, err := s.Repository.GetPhoneNumberExistence(localPhoneNumber)
	if err != nil {
		return err
	}
	if isExist {
		responsePayload := generated.GeneralErrorResponse{Message: "Phone Number Already Exist"}
		return ctx.JSON(http.StatusConflict, responsePayload)
	}

	hashedPassword := hashAndSalt([]byte(request.Password))

	profileCreate := repository.Profile{
		FullName:    request.FullName,
		CountryCode: countryCode,
		PhoneNumber: localPhoneNumber,
		Password:    hashedPassword,
	}

	createdID, err := s.Repository.CreateProfile(profileCreate)

	if err != nil {
		return err
	}

	resp := generated.CreateProfileResponse{CreatedId: &createdID, Message: "Profile is successfully created"}

	return ctx.JSON(http.StatusCreated, resp)
}
