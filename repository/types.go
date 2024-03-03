// This file contains types that are used in the repository layer.
package repository

type User struct {
	ID          int64
	PhoneNumber string
	Password    string
	FullName    string
}
