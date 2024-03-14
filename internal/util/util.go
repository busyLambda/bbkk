package util

import (
	"log"
	"os/exec"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Formats the flags and the jar args to java so that it's executed properly.
func JavaCmd(dir string, jar string, flags string) *exec.Cmd {
	e := []string{"-jar", jar}

	j := strings.ReplaceAll(flags, "\n", " ")

	var f []string
	if j != "" {
		f = append(strings.Split(j, " "), e...)
	} else {
		f = e
	}
	log.Printf("CMD_FLAGS: %s\n", f)

	c := exec.Command("java", f...)
	c.Dir = dir

	return c
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type RegistrationForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
