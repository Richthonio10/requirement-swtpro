package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Richthonio10/requirement-swtpro/generated"
	utilsHelper "github.com/Richthonio10/requirement-swtpro/utils"
	"github.com/Richthonio10/requirement-swtpro/repository"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)


func Test_generateToken(t *testing.T) {
	type args struct {
		user repository.User
	}
	tests := []struct {
		name    string
		args    args
		detailRes string
		Err error
	}{
		{
			name: "passed",
			args: args{
				user: repository.User{
					ID:          1,
					PhoneNumber: "+628223344556",
				},
			},
			detailRes: "let's say this is a jwt token",
			Err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := generateToken(tt.args.user)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.Err) {
				t.Errorf("Error When generateToken() %s, Err = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.Err))
			}
			if (len(res) != 0) != (len(tt.detailRes) != 0) {
				t.Errorf("Result When generateToken() %s, detailRes = %s", res, tt.detailRes)
			}
		})
	}
}

func Test_getSessionClaims(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		args    args
		detailRes SessionClaims
		Err error
	}{
		{
			name: "passed",
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					jwt, _ := generateToken(repository.User{
						ID:          1,
						PhoneNumber: "+62821232342",
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					return echo.New().NewContext(req, res)
				}(),
			},
			detailRes: SessionClaims{
				StandardClaims: jwt.StandardClaims{
					Issuer: "some-issuer",
				},
				UserID:      1,
				PhoneNumber: "+62821232342",
			},
			Err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := getSessionClaims(tt.args.ctx)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.Err) {
				t.Errorf("Error When getSessionClaims() = %s, Err = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.Err))
			}
			res.ExpiresAt = 0
			if res != tt.detailRes {
				t.Errorf("Result When getSessionClaims() = %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_createResponseHeader(t *testing.T) {
	type args struct {
		statusCode     int
		messages []string
		successful    bool
	}
	tests := []struct {
		name    string
		args    args
		detailRes generated.ResponseHeader
	}{
		{
			name: "status code",
			args: args{
				statusCode: utilsHelper.HttpErrorCode,
			},
			detailRes: generated.ResponseHeader{
				StatusCode: func() *int {
					res := utilsHelper.HttpErrorCode
					return &res
				}(),
			},
		},
		{
			name: "message",
			args: args{
				messages: []string{"expected error"},
			},
			detailRes: generated.ResponseHeader{
				Messages: func() *[]string {
					res := []string{"expected error"}
					return &res
				}(),
			},
		},
		{
			name: "successful",
			args: args{
				successful: true,
			},
			detailRes: generated.ResponseHeader{
				Successful: func() *bool {
					res := true
					return &res
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := createResponseHeader(tt.args.statusCode, tt.args.messages, tt.args.successful)
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When createResponseHeader() %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_createHashPassword(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		detailRes string
		Err error
	}{
		{
			name: "password is too long",
			args: args{
				input: "1234567891123456789212345678931234567894123456789512345678961234567897123",
			},
			detailRes: "",
			Err: errors.New("bcrypt: password length exceeds 72 bytes"),
		},
		{
			name: "passed",
			args: args{
				input: "SawitPro123$",
			},
			detailRes: "token already hashed",
			Err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := createHashPassword(tt.args.input)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.Err) {
				t.Errorf("Error When createHashPassword() %s, Err = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.Err))
			}
			if (len(res) != 0) != (len(tt.detailRes) != 0) {
				t.Errorf("Result When createHashPassword() %s, detailRes = %s", res, tt.detailRes)
			}
		})
	}
}

func Test_comparePasswords(t *testing.T) {
	type args struct {
		hashedPassword string
		plainPassword  string
	}
	tests := []struct {
		name    string
		args    args
		detailRes bool
	}{
		{
			name: "not match",
			args: args{
				hashedPassword: "$2a$04$1IjAa.80dLp2uNt.ls0pGe7JKv5QpPCo.qYwGPZjYQrK/BFL2ZDwG",
				plainPassword:  "SawitPro@1231",
			},
			detailRes: false,
		},
		{
			name: "match",
			args: args{
				hashedPassword: "$2a$04$1IjAa.80dLp2uNt.ls0pGe7JKv5QpPCo.qYwGPZjYQrK/BFL2ZDwG",
				plainPassword:  "SawitPro123$",
			},
			detailRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := comparePasswords(tt.args.hashedPassword, tt.args.plainPassword)
			if res != tt.detailRes {
				t.Errorf("Result When comparePasswords() %t, detailRes = %t", res, tt.detailRes)
			}
		})
	}
}
