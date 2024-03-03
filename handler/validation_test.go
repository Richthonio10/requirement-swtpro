package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Richthonio10/requirement-swtpro/generated"
	utilsHelper "github.com/Richthonio10/requirement-swtpro/utils"
	"github.com/Richthonio10/requirement-swtpro/repository"
	"github.com/golang/mock/gomock"
)


func Test_validateFullName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		detailRes []string
	}{
		{
			name: "too short",
			args: args{
				input: "!a",
			},
			detailRes: []string{
				"Full name must be at minimum 3 characters and maximum 60 characters",
			},
		},
		{
			name: "too long",
			args: args{
				input: "1234567891123456789212345678931234567894123456789512345678961234567897123",
			},
			detailRes: []string{
				"Full name must be at minimum 3 characters and maximum 60 characters",
			},
		},
		{
			name: "passed",
			args: args{
				input: "1234567890",
			},
			detailRes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateFullName(tt.args.input)
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When validateFullName() %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_checkPhoneNumber(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		detailRes []string
	}{
		{
			name: "too short",
			args: args{
				input: "+62822334",
			},
			detailRes: []string{
				"Phone numbers must be at minimum 10 characters and maximum 13 characters",
			},
		},
		{
			name: "too long",
			args: args{
				input: "+6282233445566",
			},
			detailRes: []string{
				"Phone numbers must be at minimum 10 characters and maximum 13 characters",
			},
		},
		{
			name: "invalid prefix",
			args: args{
				input: "1234567890",
			},
			detailRes: []string{
				"Phone numbers must start with the Indonesia country code “+62”",
			},
		},
		{
			name: "passed",
			args: args{
				input: "+628223344556",
			},
			detailRes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := checkPhoneNumber(tt.args.input)
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When checkPhoneNumber() %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_ValidatePhoneNumber(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		detailRes bool
		detailErr error
	}{
		{
			name: "error GetUserByPhoneNumber",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344556").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			detailRes: false,
			detailErr: errors.New("expected GetUserByPhoneNumber error"),
		},
		{
			name: "not a new phone number",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344556").
					Return(repository.User{
						ID: 1,
					}, nil).
					Times(1)
			},
			detailRes: false,
			detailErr: nil,
		},
		{
			name: "a new phone number",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+628223344556").
					Return(repository.User{}, nil).
					Times(1)
			},
			detailRes: true,
			detailErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			res, err := ValidatePhoneNumber(tt.args.ctx, s, tt.args.phoneNumber)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When ValidatePhoneNumber() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When ValidatePhoneNumber() %+v, detailRes = %+v", res, tt.detailRes)
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_validatePassword(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		detailRes bool
	}{
		{
			name: "too short",
			args: args{
				input: "!a",
			},
			detailRes: false,
		},
		{
			name: "too long",
			args: args{
				input: "!a",
			},
			detailRes: false,
		},
		{
			name: "no capital",
			args: args{
				input: "sawit@pr0",
			},
			detailRes: false,
		},
		{
			name: "no numeric",
			args: args{
				input: "sawit@pro",
			},
			detailRes: false,
		},
		{
			name: "no non alpha-numeric",
			args: args{
				input: "sawitpr0",
			},
			detailRes: false,
		},
		{
			name: "passed",
			args: args{
				input: "SawitPro!2344",
			},
			detailRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validatePassword(tt.args.input)
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When validatePassword() %t, detailRes = %t", res, tt.detailRes)
			}
		})
	}
}

func Test_validateRegistration(t *testing.T) {
	type args struct {
		request generated.RegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		detailRes []string
	}{
		{
			name: "invalid phone number",
			args: args{
				request: generated.RegistrationRequest{
					PhoneNumber: "+62812345",
					FullName:    "Some Full Name",
					Password:    "SawitPro!2344",
				},
			},
			detailRes: []string{
				"Phone numbers must be at minimum 10 characters and maximum 13 characters",
			},
		},
		{
			name: "invalid full name",
			args: args{
				request: generated.RegistrationRequest{
					PhoneNumber: "+6298080980",
					FullName:    "SP",
					Password:    "SawitPro!2344",
				},
			},
			detailRes: []string{
				"Full name must be at minimum 3 characters and maximum 60 characters",
			},
		},
		{
			name: "invalid password",
			args: args{
				request: generated.RegistrationRequest{
					PhoneNumber: "+6298080980",
					FullName:    "Some Full Name",
					Password:    "sawitpro",
				},
			},
			detailRes: []string{
				"Passwords must be minimum 6 characters and maximum 64 characters, containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateRegistration(tt.args.request)
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When validateRegistration() %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}