package db

import (
	"fmt"

	"github.com/busyLambda/bbkk/domain/user"
	"github.com/busyLambda/bbkk/internal/models"
)

func (d DbManager) InsertUser(u models.User) error {
	return d.Conn.Create(&u).Error
}

func (d DbManager) DeleteUser(id int) error {
	u, err := d.GetUserByID(id)
	if err != nil {
		return err
	}

	if u.Role == user.SUPERADMIN {
		return fmt.Errorf("CANNOT_DELETE_SUPERADMIN")
	}

	return d.Conn.Delete(&u).Error
}

func (d DbManager) GetUserByID(id int) (u models.User, err error) {
	err = d.Conn.Where("id = ?", id).First(&u).Error
	return
}

func (d DbManager) GetUserByUsername(un string) (u models.User, err error) {
	err = d.Conn.Where("username = ?", un).First(&u).Error
	return
}
