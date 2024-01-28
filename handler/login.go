package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/hasbiasshidiq/simple-profile/generated"
	"github.com/hasbiasshidiq/simple-profile/repository"
	"github.com/labstack/echo/v4"
)

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
	profileMetadata := repository.ProfileMetaData{ProfileID: existingProfile.ID}
	_, err = s.Repository.UpsertProfileMetaData(profileMetadata)
	if err != nil {
		log.Println("error Upserting MetaData : ", err)
		responsePayload := generated.GeneralErrorResponse{Message: "Internal Server Error"}
		return ctx.JSON(http.StatusInternalServerError, responsePayload)
	}

	resp := generated.LoginResponse{JwtToken: token, UserId: int(existingProfile.ID)}

	return ctx.JSON(http.StatusOK, resp)
}
