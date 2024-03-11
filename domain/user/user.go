package user

import (
	"github.com/busyLambda/bbkk/domain"
	"github.com/busyLambda/bbkk/internal/util"
)

type Username = string

func NewUsername(u string) (un Username, err error) {
	if l := len(u); l > 48 {
		err = &domain.TooLong{Expected: 48, Found: l}
		return
	}

	if l := len(u); l < 1 {
		err = &domain.TooShort{}
		return
	}

	un = Username(u)
	return
}

type Password = string

func NewPassword(p string) (pw Password, err error) {
	if l := len(p); l > 128 {
		err = &domain.TooLong{Expected: 128, Found: l}
		return
	}

	if len(p) < 12 {
		err = &domain.TooShort{}
		return
	}

	p, err = util.HashPassword(p)
	if err != nil {
		return
	}

	pw = Password(p)
	return
}

type Role = uint

const (
	SUPERADMIN Role = iota
	ADMIN
)
