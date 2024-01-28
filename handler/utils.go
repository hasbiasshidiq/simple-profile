package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func extractUserIDFromToken(token string) (profileID int, err error) {

	// read public key from .key.pub file
	pubKey, err := ioutil.ReadFile("cert/jwtRS256.key.pub")
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

	return err == nil
}

func createToken(profile repository.Profile) (tokenString string, err error) {

	prvKey, err := ioutil.ReadFile("cert/jwtRS256.key")
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
