package repository

import (
	"context"
	"fmt"
	"strings"
)

func (r *Repository) GetUserByID(ctx context.Context, userID int64) (user User, err error) {
	rows, err := r.Db.QueryContext(ctx, queryGetUserByID, userID)
	if err != nil {
		return user, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.FullName)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user User, err error) {
	rows, err := r.Db.QueryContext(ctx, queryGetUserByPhoneNumber, phoneNumber)
	if err != nil {
		return user, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.FullName)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (r *Repository) CreateLoginCount(ctx context.Context, userID int64) (err error) {
	_, err = r.Db.ExecContext(ctx, queryCreateLoginCount, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) InsertUser(ctx context.Context, data User) (userID int64, err error) {
	rows, err := r.Db.QueryContext(ctx, queryInsertUser,
		data.PhoneNumber,
		data.Password,
		data.FullName)
	if err != nil {
		return userID, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return userID, err
		}
	}

	return userID, err
}

func (r *Repository) UpdateUser(ctx context.Context, data User) (err error) {
	var (
		updatedFields []string
		params        []any
	)
	params = append(params, data.ID)
	if data.PhoneNumber != "" {
		params = append(params, data.PhoneNumber)
		updatedFields = append(updatedFields, fmt.Sprintf("phone_number = $%d", len(params)))
	}
	if data.FullName != "" {
		params = append(params, data.FullName)
		updatedFields = append(updatedFields, fmt.Sprintf("full_name = $%d", len(params)))
	}
	if len(params) == 1 {
		return nil
	}
	_, err = r.Db.ExecContext(ctx,
		fmt.Sprintf(queryUpdateUser, strings.Join(updatedFields, ", ")),
		params...)
	if err != nil {
		return err
	}
	return nil
}
