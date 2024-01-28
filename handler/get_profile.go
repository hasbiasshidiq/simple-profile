package handler

import (
	"net/http"

	"github.com/hasbiasshidiq/simple-profile/generated"
	"github.com/labstack/echo/v4"
)

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
