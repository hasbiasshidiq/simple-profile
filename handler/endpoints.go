package handler

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
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

	resp := generated.CreateProfileResponse{CreatedId: &createdID, Message: "Profile is successfully created"}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) PostLogin(ctx echo.Context) error {

	var request generated.LoginRequest

	err := ctx.Bind(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	hasPrefix := strings.HasPrefix(request.PhoneNumber, "+62")
	if !hasPrefix {
		responsePayload := generated.GeneralErrorResponse{Message: "Account not found"}
		return ctx.JSON(http.StatusBadRequest, responsePayload)
	}

	localPhoneNumber := strings.Replace(request.PhoneNumber, "+62", "", -1)
	existingProfile, err := s.Repository.GetProfileByPhoneNumber(localPhoneNumber)
	if err == sql.ErrNoRows {
		responsePayload := generated.GeneralErrorResponse{Message: "Account not found"}
		return ctx.JSON(http.StatusBadRequest, responsePayload)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	isPasswordValid := comparePasswords(existingProfile.Password, []byte(request.Password))
	if !isPasswordValid {
		responsePayload := generated.GeneralErrorResponse{Message: "Password doesn't match"}
		return ctx.JSON(http.StatusBadRequest, responsePayload)
	}

	token, err := createToken(existingProfile)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	resp := generated.LoginResponse{JwtToken: token, UserId: int(existingProfile.ID)}

	return ctx.JSON(http.StatusCreated, resp)
}

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

func hashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func createToken(profile repository.Profile) (tokenString string, err error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return tokenString, errors.New("could not determine current file")
	}
	// Get the directory where current file reside
	dir := filepath.Dir(filename)
	// Get the parent directory
	parentDir := filepath.Dir(dir)

	prvKey, err := ioutil.ReadFile(parentDir + "/cert/jwtRS256.key")
	if err != nil {
		log.Println(err.Error())
		return tokenString, err
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": profile.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign the token with the secret key
	tokenString, err = token.SignedString(prvKey)
	if err != nil {
		return tokenString, err
	}

	return tokenString, err
}
