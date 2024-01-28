package handler

import (
	"net/http"
	"regexp"
	"unicode"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
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
