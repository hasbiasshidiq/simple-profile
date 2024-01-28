package handler

import (
	"log"
	"net/http"
	"regexp"
	"unicode"

	"github.com/go-playground/validator"
	"github.com/hasbiasshidiq/simple-profile/generated"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	ProfileValidator struct {
		FullName    string `json:"full_name" validate:"required,min=3,max=60"`
		PhoneNumber string `json:"phone_number" validate:"required,min=10,max=13,indonesiaCountryCodePrefix"`
		Password    string `json:"password" validate:"required,min=6,max=64,strongPassword"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// validatePhoneWithPrefix is a custom validation function for phone numbers with prefix "+62"
func validatePhoneWithPrefix(fl validator.FieldLevel) bool {

	phoneNumber := fl.Field().String()

	// Use a regular expression to check if the phone number starts with "+62"
	match, _ := regexp.MatchString(`^\+62[0-9]+$`, phoneNumber)

	return match
}

// validateStrongPassword is a custom validation function for strong passwords
func validateStrongPassword(fl validator.FieldLevel) bool {

	password := fl.Field().String()

	// Check if the password contains at least 1 uppercase letter, 1 number, and 1 special character
	hasUppercase := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsNumber(char):
			hasNumber = true
		case !unicode.IsLetter(char) && !unicode.IsNumber(char):
			hasSpecial = true
		}
	}

	return hasUppercase && hasNumber && hasSpecial
}

func (s *Server) PostRegister(ctx echo.Context) error {

	var request generated.RegisterRequest

	customValidator := validator.New()
	customValidator.RegisterValidation("indonesiaCountryCodePrefix", validatePhoneWithPrefix)
	customValidator.RegisterValidation("strongPassword", validateStrongPassword)

	ctx.Echo().Validator = &CustomValidator{validator: customValidator}

	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	profileValidator := ProfileValidator{
		FullName:    request.FullName,
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	}
	if err = ctx.Validate(profileValidator); err != nil {
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

	resp := generated.RegisterResponse{CreatedId: &createdID, Message: "Profile is successfully created"}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) PostLogin(ctx echo.Context) error {

	resp := generated.LoginResponse{}

	return ctx.JSON(http.StatusCreated, resp)
}

func hashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
