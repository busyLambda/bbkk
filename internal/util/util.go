package util

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetRssByPid(pid int) (int, error) {
	c := fmt.Sprintf("cat /proc/%d/status | grep RSS | awk '{print $2}'", pid)
	cmd := exec.Command("bash", "-c", c)
	mstr, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.Replace(string(mstr), "\n", "", 1))
}

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

func ServerDirName(name string, id string) string {
	return fmt.Sprintf("%s-%s", name, id)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type RegistrationForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerForm struct {
	Name         string `json:"name"`
	DedicatedRam uint   `json:"dedicated_ram"`
}

type UserKey struct{}

type ServerStats struct {
	IsRunning bool `json:"is_running"`
	MemUse    uint `json:"mem_use"`
}
