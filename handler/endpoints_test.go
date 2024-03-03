package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	utilsHelper "github.com/Richthonio10/requirement-swtpro/utils"
	"github.com/Richthonio10/requirement-swtpro/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func Test_Register(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		statusCode int
		detailErr        error
	}{
		{
			name: "no request body",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(``)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "invalid request",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "",
						"full_name": "",
						"password": ""
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "error ValidatePhoneNumber",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "phone number already registered",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID: 1,
					}, nil).
					Times(1)
			},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "error InsertUser",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID: 0,
					}, nil).
					Times(1)

				fields.Repository.EXPECT().InsertUser(context.Background(), gomock.AssignableToTypeOf(repository.User{})).
					Return(int64(0), errors.New("expected InsertUser error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID: 0,
					}, nil).
					Times(1)

				fields.Repository.EXPECT().InsertUser(context.Background(), gomock.AssignableToTypeOf(repository.User{})).
					Return(int64(1), nil).
					Times(1)
			},
			statusCode: http.StatusOK,
			detailErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			err := s.Register(tt.args.ctx)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When Register() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if err == nil {
				if tt.args.ctx.Response().Status != tt.statusCode {
					t.Errorf("Result When Register() %d, statusCode = %d", tt.args.ctx.Response().Status, tt.statusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_Login(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		statusCode int
		detailErr        error
	}{
		{
			name: "no request body",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(``)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
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
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "user not found",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID: 0,
					}, nil).
					Times(1)
			},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "wrong password",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"password": "Random@123"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID:       1,
						Password: "$2a$04$1IjAa.80dLp2uNt.ls0pGe7JKv5QpPCo.qYwGPZjYQrK/BFL2ZDwG",
					}, nil).
					Times(1)
			},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "error CreateLoginCount",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID:       1,
						Password: "$2a$04$1IjAa.80dLp2uNt.ls0pGe7JKv5QpPCo.qYwGPZjYQrK/BFL2ZDwG",
					}, nil).
					Times(1)

				fields.Repository.EXPECT().CreateLoginCount(context.Background(), int64(1)).
					Return(errors.New("expected CreateLoginCount error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPost, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"password": "SawitPro123$"
					}`)))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID:       1,
						Password: "$2a$04$1IjAa.80dLp2uNt.ls0pGe7JKv5QpPCo.qYwGPZjYQrK/BFL2ZDwG",
					}, nil).
					Times(1)

				fields.Repository.EXPECT().CreateLoginCount(context.Background(), int64(1)).
					Return(nil).
					Times(1)
			},
			statusCode: http.StatusOK,
			detailErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			err := s.Login(tt.args.ctx)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When Login() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if err == nil {
				if tt.args.ctx.Response().Status != tt.statusCode {
					t.Errorf("Result When Register() %d, statusCode = %d", tt.args.ctx.Response().Status, tt.statusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_GetProfile(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		statusCode int
		detailErr        error
	}{
		{
			name: "invalid authorization",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusForbidden,
			detailErr:        nil,
		},
		{
			name: "error GetUserByID",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByID(context.Background(), int64(1)).
					Return(repository.User{}, errors.New("expected GetUserByID error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodGet, "url", nil)
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByID(context.Background(), int64(1)).
					Return(repository.User{}, nil).
					Times(1)
			},
			statusCode: http.StatusOK,
			detailErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			err := s.GetProfile(tt.args.ctx)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When GetProfile() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if err == nil {
				if tt.args.ctx.Response().Status != tt.statusCode {
					t.Errorf("Result When Register() %d, statusCode = %d", tt.args.ctx.Response().Status, tt.statusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}

func Test_UpdateProfile(t *testing.T) {
	type fields struct {
		mockCtrl   *gomock.Controller
		Repository *repository.MockRepositoryInterface
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mock           func(fields *fields)
		statusCode int
		detailErr        error
	}{
		{
			name: "invalid authorization",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", nil)
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusForbidden,
			detailErr:        nil,
		},
		{
			name: "no request body",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(``)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "no params at all",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "all params have empty value",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "",
						"full_name": ""
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "all params have empty value",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "",
						"full_name": ""
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
		{
			name: "invalid phone number",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "123",
						"full_name": "Some Full Name"
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock:           func(fields *fields) {},
			statusCode: http.StatusBadRequest,
			detailErr:        nil,
		},
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
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name"
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{}, errors.New("expected GetUserByPhoneNumber error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "phone number is already registered",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name"
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{
						ID: 2,
					}, nil).
					Times(1)
			},
			statusCode: http.StatusConflict,
			detailErr:        nil,
		},
		{
			name: "error UpdateUser",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name"
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{}, nil).
					Times(1)

				fields.Repository.EXPECT().UpdateUser(context.Background(),
					repository.User{
						ID:          1,
						PhoneNumber: "+62821232342",
						FullName:    "Some Full Name",
					}).
					Return(errors.New("expected UpdateUser error")).
					Times(1)
			},
			statusCode: http.StatusInternalServerError,
			detailErr:        nil,
		},
		{
			name: "passed",
			fields: func() fields {
				mockCtrl := gomock.NewController(t)
				return fields{
					mockCtrl:   mockCtrl,
					Repository: repository.NewMockRepositoryInterface(mockCtrl),
				}
			}(),
			args: args{
				ctx: func() echo.Context {
					req, _ := http.NewRequest(http.MethodPatch, "url", bytes.NewBuffer([]byte(`{
						"phone_number": "+62821232342",
						"full_name": "Some Full Name"
					}`)))
					jwt, _ := generateToken(repository.User{
						ID: 1,
					})
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
					res := httptest.NewRecorder()
					c := echo.New().NewContext(req, res)
					return c
				}(),
			},
			mock: func(fields *fields) {
				fields.Repository.EXPECT().GetUserByPhoneNumber(context.Background(), "+62821232342").
					Return(repository.User{}, nil).
					Times(1)

				fields.Repository.EXPECT().UpdateUser(context.Background(),
					repository.User{
						ID:          1,
						PhoneNumber: "+62821232342",
						FullName:    "Some Full Name",
					}).
					Return(nil).
					Times(1)
			},
			statusCode: http.StatusOK,
			detailErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			tt.mock(&tt.fields)
			err := s.UpdateProfile(tt.args.ctx)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When UpdateProfile() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if err == nil {
				if tt.args.ctx.Response().Status != tt.statusCode {
					t.Errorf("Result When Register() %d, statusCode = %d", tt.args.ctx.Response().Status, tt.statusCode)
				}
			}
			tt.fields.mockCtrl.Finish()
		})
	}
}
