package handler

import (
	"encoding/json"
	"net/http"

	utilsHelper "github.com/Richthonio10/requirement-swtpro/utils"
	"github.com/Richthonio10/requirement-swtpro/generated"
	"github.com/Richthonio10/requirement-swtpro/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (s *Server) Login(ctx echo.Context) error {
	var (
		request  generated.LoginRequest
		response generated.LoginResponse
	)

	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		log.Errorf("Error When Decode Request: %s with request: %s", err.Error(), request)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Bad request"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), request.PhoneNumber)
	if err != nil {
		log.Errorf("Error When GetUserByPhoneNumber: %s with phone number: %s", err.Error(), request.PhoneNumber)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	if user.ID == 0 {
		response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, []string{"Phone number is not found"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	if !comparePasswords(user.Password, request.Password) {
		response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, []string{"Wrong password"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	jwtToken, err := generateToken(user)
	if err != nil {
		log.Errorf("Error When generateToken: %s with data user: %s", err.Error(), user)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	err = s.Repository.CreateLoginCount(ctx.Request().Context(), user.ID)
	if err != nil {
		log.Errorf("Error When CreateLoginCount: %s with user id: %s", err.Error(), user.ID)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = createResponseHeader(200, []string{"Successfully Login!"}, true)
	response.Data = &generated.LoginResponseData{
		Id:  user.ID,
		Jwt: jwtToken,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) Register(ctx echo.Context) error {
	var (
		request  generated.RegistrationRequest
		response generated.RegistrationResponse
	)

	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		log.Errorf("Error When Decode Request: %s with request: %s", err.Error(), request)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Bad request"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	requestValidationErrors := validateRegistration(request)
	if len(requestValidationErrors) != 0 {
		response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, requestValidationErrors, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	phoneNumber, err := ValidatePhoneNumber(ctx.Request().Context(), s, request.PhoneNumber)
	if err != nil {
		log.Errorf("Error when ValidatePhoneNumber: %s with phoneNumber: %s", err.Error(), request.PhoneNumber)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if !phoneNumber {
		response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, []string{"Phone number is already registered"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	hashedPassword, err := createHashPassword(request.Password)
	if err != nil {
		log.Errorf("Error when createHashPassword: %s", err.Error())
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"error when hashing password"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	id, err := s.Repository.InsertUser(ctx.Request().Context(), repository.User{
		PhoneNumber: request.PhoneNumber,
		Password:    hashedPassword,
		FullName:    request.FullName,
	})
	if err != nil {
		log.Errorf("Error when InsertUser: %s", err.Error())
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error when register"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = createResponseHeader(200, []string{"Successfully Create User Register!"}, true)
	response.Data = &generated.RegistrationResponseData{
		Id: id,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) GetProfile(ctx echo.Context) error {
	var (
		response generated.GetProfileResponse
	)

	sessionClaims, err := getSessionClaims(ctx)
	if err != nil {
		response.Header = createResponseHeader(utilsHelper.AuthorizationErrorCode, []string{err.Error()}, false)
		return ctx.JSON(http.StatusForbidden, response)
	}

	user, err := s.Repository.GetUserByID(ctx.Request().Context(), sessionClaims.UserID)
	if err != nil {
		log.Errorf("Error When GetUserByID: %s", err.Error())
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = createResponseHeader(200, []string{"Successfully Get User Profile!"}, true)
	response.Data = &generated.GetProfileResponseData{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	var (
		request  generated.UpdateProfileRequest
		response generated.UpdateProfileResponse
	)

	var (
		fullName     string
		phoneNumber  string
		updated bool
	)

	sessionClaims, err := getSessionClaims(ctx)
	if err != nil {
		response.Header = createResponseHeader(utilsHelper.AuthorizationErrorCode, []string{err.Error()}, false)
		return ctx.JSON(http.StatusForbidden, response)
	}

	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		log.Errorf("Error When Decode Request: %s with request: %s", err.Error(), request)
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Bad request"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	if request.FullName != nil && *request.FullName != "" {
		fullName = *request.FullName

		if errorMessages := validateFullName(fullName); len(errorMessages) != 0 {
			response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, errorMessages, false)
			return ctx.JSON(http.StatusBadRequest, response)
		}

		updated = true;
	}

	if request.PhoneNumber != nil && *request.PhoneNumber != "" {
		phoneNumber = *request.PhoneNumber

		if errorMessages := checkPhoneNumber(phoneNumber); len(errorMessages) != 0 {
			response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, errorMessages, false)
			return ctx.JSON(http.StatusBadRequest, response)
		}

		user, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), phoneNumber)
		if err != nil {
			log.Errorf("Error When GetUserByPhoneNumber: %s", err.Error())
			response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
			return ctx.JSON(http.StatusInternalServerError, response)
		}
		if user.ID != 0 && user.ID != sessionClaims.UserID {
			response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, []string{"Phone number is already registered"}, false)
			return ctx.JSON(http.StatusConflict, response)
		}

		updated = true;
	}
	if updated == false {
		response.Header = createResponseHeader(utilsHelper.ValidationErrorCode, []string{"No update"}, false)
		return ctx.JSON(http.StatusBadRequest, response)
	}

	err = s.Repository.UpdateUser(ctx.Request().Context(), repository.User{
		ID:          sessionClaims.UserID,
		PhoneNumber: phoneNumber,
		FullName:    fullName,
	})
	if err != nil {
		log.Errorf("Error When UpdateUser: %s", err.Error())
		response.Header = createResponseHeader(utilsHelper.HttpErrorCode, []string{"Internal Server Error"}, false)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header = createResponseHeader(200, []string{"Successfully Update User!"}, true)

	return ctx.JSON(http.StatusOK, response)
}
