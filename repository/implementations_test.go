package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	utilsHelper "github.com/Richthonio10/requirement-swtpro/utils"
)

func Test_Repository_GetUserByID(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_GetUserByID] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		detailRes User
		detailErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByID)).
					WithArgs(int64(1)).
					WillReturnError(errors.New("expected error"))
			},
			detailRes: User{},
			detailErr: errors.New("expected error"),
		},
		{
			name: "no data",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"})

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByID)).
					WithArgs(int64(1)).
					WillReturnRows(resultRows)
			},
			detailRes: User{},
			detailErr: nil,
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"}).
					AddRow(1, "+628223344556", "<password>", "Sawit")

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByID)).
					WithArgs(int64(1)).
					WillReturnRows(resultRows)
			},
			detailRes: User{
				ID:          1,
				PhoneNumber: "+628223344556",
				Password:    "<password>",
				FullName:    "Sawit",
			},
			detailErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			res, err := r.GetUserByID(tt.args.ctx, tt.args.userID)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When GetUserByID() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When GetUserByID() %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_Repository_GetUserByPhoneNumber(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_GetUserByPhoneNumber] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
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
		detailRes User
		detailErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByPhoneNumber)).
					WithArgs("+628223344556").
					WillReturnError(errors.New("expected error"))
			},
			detailRes: User{},
			detailErr: errors.New("expected error"),
		},
		{
			name: "no data",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"})

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByPhoneNumber)).
					WithArgs("+628223344556").
					WillReturnRows(resultRows)
			},
			detailRes: User{},
			detailErr: nil,
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+628223344556",
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id", "phone_number", "password", "full_name"}).
					AddRow(1, "+628223344556", "<password>", "Sawit")

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryGetUserByPhoneNumber)).
					WithArgs("+628223344556").
					WillReturnRows(resultRows)
			},
			detailRes: User{
				ID:          1,
				PhoneNumber: "+628223344556",
				Password:    "<password>",
				FullName:    "Sawit",
			},
			detailErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			res, err := r.GetUserByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When GetUserByPhoneNumber() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When GetUserByPhoneNumber() %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_Repository_CreateLoginCount(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_CreateLoginCount] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		detailErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(queryCreateLoginCount)).
					WithArgs(int64(1)).
					WillReturnError(errors.New("expected error"))
			},
			detailErr: errors.New("expected error"),
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(queryCreateLoginCount)).
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			detailErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			err := r.CreateLoginCount(tt.args.ctx, tt.args.userID)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When CreateLoginCount() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
		})
	}
}

func Test_Repository_InsertUser(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_InsertUser] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx  context.Context
		data User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		detailRes int64
		detailErr error
	}{
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					PhoneNumber: "+628223344556",
					Password:    "<password>",
					FullName:    "Sawit",
				},
			},
			mock: func(fields *fields) {
				sqlMock.ExpectQuery(regexp.QuoteMeta(queryInsertUser)).
					WithArgs("+628223344556", "<password>", "Sawit").
					WillReturnError(errors.New("expected error"))
			},
			detailRes: 0,
			detailErr: errors.New("expected error"),
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					PhoneNumber: "+628223344556",
					Password:    "<password>",
					FullName:    "Sawit",
				},
			},
			mock: func(fields *fields) {
				resultRows := sqlmock.
					NewRows([]string{"id"}).
					AddRow(1)

				sqlMock.ExpectQuery(regexp.QuoteMeta(queryInsertUser)).
					WithArgs("+628223344556", "<password>", "Sawit").
					WillReturnRows(resultRows)
			},
			detailRes: 1,
			detailErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			res, err := r.InsertUser(tt.args.ctx, tt.args.data)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When InsertUser() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
			if !reflect.DeepEqual(res, tt.detailRes) {
				t.Errorf("Result When InsertUser()  %+v, detailRes = %+v", res, tt.detailRes)
			}
		})
	}
}

func Test_Repository_UpdateUser(t *testing.T) {
	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Errorf("[Test_Repository_UpdateUser] %s", err.Error())
		return
	}
	defer dbMock.Close()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx  context.Context
		data User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(fields *fields)
		detailErr error
	}{
		{
			name: "no changes",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					ID:          1,
					PhoneNumber: "",
					FullName:    "",
				},
			},
			mock:    func(fields *fields) {},
			detailErr: nil,
		},
		{
			name: "error",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					ID:          1,
					PhoneNumber: "+62812345678",
					FullName:    "New Name",
				},
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(queryUpdateUser, "phone_number = $2, full_name = $3"))).
					WithArgs(int64(1), "+62812345678", "New Name").
					WillReturnError(errors.New("expected error"))
			},
			detailErr: errors.New("expected error"),
		},
		{
			name: "passed",
			fields: fields{
				Db: dbMock,
			},
			args: args{
				ctx: context.Background(),
				data: User{
					ID:          1,
					PhoneNumber: "+62812345678",
					FullName:    "New Name",
				},
			},
			mock: func(fields *fields) {
				sqlMock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(queryUpdateUser, "phone_number = $2, full_name = $3"))).
					WithArgs(int64(1), "+62812345678", "New Name").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			detailErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			tt.mock(&tt.fields)
			err := r.UpdateUser(tt.args.ctx, tt.args.data)
			if utilsHelper.ErrorMessage(err) != utilsHelper.ErrorMessage(tt.detailErr) {
				t.Errorf("Error When UpdateUser() %s, detailErr = %s", utilsHelper.ErrorMessage(err), utilsHelper.ErrorMessage(tt.detailErr))
			}
		})
	}
}
