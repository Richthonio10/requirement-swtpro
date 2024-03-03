// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

//go:generate mockgen -source=interfaces.go -destination=interfaces.mock.gen.go -package=repository
type RepositoryInterface interface {
	GetUserByID(ctx context.Context, userID int64) (user User, err error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user User, err error)
	CreateLoginCount(ctx context.Context, userID int64) (err error)
	InsertUser(ctx context.Context, data User) (userID int64, err error)
	UpdateUser(ctx context.Context, data User) (err error)
}
