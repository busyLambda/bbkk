package models

import (
	"fmt"

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

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role == user.SUPERADMIN {
		return fmt.Errorf("cannot_delete_superadmin")
	}
	return
}
