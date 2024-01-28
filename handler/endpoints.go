package handler

import (
	"database/sql"
	"errors"
	"fmt"
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
	CreateProfileValidator struct {
		FullName    string `json:"full_name" validate:"required,min=3,max=60"`
		PhoneNumber string `json:"phone_number" validate:"required,min=10,max=13,indonesiaCountryCodePrefix"`
		Password    string `json:"password" validate:"required,min=6,max=64,strongPassword"`
	}

	UpdateProfileValidator struct {
		FullName    *string `json:"full_name" validate:"omitempty,min=3,max=60"`
		PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=13,indonesiaCountryCodePrefix"`
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
		log.Println("error fetch profile by phone number : ", err)
		return err
	}

	isPasswordValid := comparePasswords(existingProfile.Password, []byte(request.Password))
	if !isPasswordValid {
		responsePayload := generated.GeneralErrorResponse{Message: "Password doesn't match"}
		return ctx.JSON(http.StatusBadRequest, responsePayload)
	}

	token, err := createToken(existingProfile)
	if err != nil {
		log.Println("error create token : ", err)
		return err
	}
	resp := generated.LoginResponse{JwtToken: token, UserId: int(existingProfile.ID)}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetProfile(ctx echo.Context) error {
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

	profile, err := s.Repository.GetProfileByID(userID)
	if err != nil {
		responsePayload := generated.GeneralErrorResponse{Message: "Profile not found"}
		return ctx.JSON(http.StatusNotFound, responsePayload)
	}

	resp := generated.GetProfileResponse{
		FullName:    profile.FullName,
		PhoneNumber: profile.CountryCode + profile.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}

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

func extractUserIDFromToken(token string) (profileID int, err error) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return profileID, errors.New("could not determine current file")
	}
	// Get the directory where current file reside
	dir := filepath.Dir(filename)
	// Get the parent directory
	parentDir := filepath.Dir(dir)

	// read public key from .key.pub file
	pubKey, err := ioutil.ReadFile(parentDir + "/cert/jwtRS256.key.pub")
	if err != nil {
		log.Println("Can't open public key ", err)
		return profileID, err
	}

	// validate token based on public key and extract claims
	claims, validated := validateToken(pubKey, token)
	if !validated {
		return profileID, errors.New("not valid token")
	}

	// get subject
	sub, ok := claims["sub"].(float64)
	if !ok {
		return profileID, errors.New("unable to extract user ID from token")
	}

	return int(sub), nil
}

func extractToken(c echo.Context) (token string, err error) {
	// Retrieve the Authorization header from the request
	authHeader := c.Request().Header.Get("Authorization")

	// Check if the Authorization header is present
	if authHeader == "" {
		return token, errors.New("no authorization header found")
	}

	// Check if the Authorization header starts with "Bearer"
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return token, errors.New("authorization header format is invalid")
	}

	// Extract the token from the Authorization header
	token = strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func validateToken(publicKey []byte, token string) (claims jwt.MapClaims, validated bool) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		log.Println("err validate: parse key: ", err)
		return nil, false
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		log.Println("token validation error : ", err)
		return nil, false
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		log.Printf("token validation error, can't create map of claims")
		return nil, false
	}

	return claims, true
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
		log.Println("can't generate hash : ", err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println("error on comparing password and hashed password : ", err)
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
		return tokenString, err
	}

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return tokenString, err
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": profile.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign the token with the secret key
	tokenString, err = token.SignedString(parsedKey)
	if err != nil {
		return tokenString, err
	}
	return tokenString, err
}
