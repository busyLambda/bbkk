package models

import (
	"fmt"

	"github.com/busyLambda/bbkk/domain/user"
	"github.com/busyLambda/bbkk/internal/util"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username user.Username
	Password user.Password
	Sessions []Session `gorm:"foreignKey:UserID"`
	Role     user.Role
}

func NewUser(rf util.RegistrationForm, role user.Role) (u User, err error) {
	pw, err := user.NewPassword(rf.Password)
	if err != nil {
		return
	}
	u.Password = pw

	un, err := user.NewUsername(rf.Username)
	if err != nil {
		return
	}
	u.Username = un

	u.Role = role

	return
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role == user.SUPERADMIN {
		return fmt.Errorf("cannot_delete_superadmin")
	}
	return
}
