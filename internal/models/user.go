package models

import (
	"github.com/busyLambda/bbkk/domain/user"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username user.Username
	Password user.Password
	Role     user.Role
}

func NewUser(un user.Username, pw user.Password) User {
	return User{
		Username: un,
		Password: pw,
	}
}
