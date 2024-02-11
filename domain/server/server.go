package server

import "fmt"

type ServerName = string

func NewServerName(s string) (sn ServerName, err error) {
	if len(s) > 128 {
		err = fmt.Errorf("TOO_LONG")
		return
	}

	if len(s) < 1 {
		err = fmt.Errorf("TOO_SHORT")
		return
	}

	sn = ServerName(s)
	return
}
