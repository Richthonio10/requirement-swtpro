package handler

import (
	"context"
	"regexp"
	"strings"

	"github.com/Richthonio10/requirement-swtpro/generated"
)

func validateFullName(input string) []string {
	var res []string

	if len(input) < 3 || len(input) > 60 {
		res = append(res, "Full name must be at minimum 3 characters and maximum 60 characters")
	}

	return res
}

func ValidatePhoneNumber(
	ctx context.Context,
	s *Server,
	phoneNumber string,
) (bool, error) {
	var (
		res bool
		err error
	)

	user, err := s.Repository.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return res, err
	}

	if user.ID == 0 {
		res = true
	}

	return res, err
}

func checkPhoneNumber(input string) []string {
	var res []string

	if len(input) < 10 || len(input) > 13 {
		res = append(res, "Phone numbers must be at minimum 10 characters and maximum 13 characters")
	}
	if !strings.HasPrefix(input, "+62") {
		res = append(res, "Phone numbers must start with the Indonesia country code “+62”")
	}

	return res
}

func validatePassword(input string) bool {
	validations := []string{".{6,64}", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, validation := range validations {
		match, _ := regexp.MatchString(validation, input)
		if !match {
			return false
		}
	}

	return true
}

func validateRegistration(request generated.RegistrationRequest) []string {
	var res []string
	res = append(res, validateFullName(request.FullName)...)
	res = append(res, checkPhoneNumber(request.PhoneNumber)...)

	if !validatePassword(request.Password) {
		res = append(res, "Passwords must be minimum 6 characters and maximum 64 characters, containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters")
	}

	return res
}