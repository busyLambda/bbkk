package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	e := []string{"-jar", jar, "--nogui"}

	j := strings.ReplaceAll(flags, "\n", " ")

	var f []string
	if j != "" {
		f = append(strings.Split(j, " "), e...)
	} else {
		f = e
	}

	c := exec.Command("java", f...)
	c.Dir = dir

	return c
}

// Creates the server on disk
func CreateServer(name string, id uint, version string, build string) error {
	err := os.MkdirAll(fmt.Sprintf("servers/%s-%d", name, id), 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("servers/%s-%d/server.jar", name, id))
	if err != nil {
		return err
	}
	defer file.Close()

	jar_name := fmt.Sprintf("paper-%s-%s.jar", version, build)
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%s/downloads/%s", version, build, jar_name)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ServerDirName(name string, id uint) string {
	return fmt.Sprintf("servers/%s-%d", name, id)
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
	Version      string `json:"version"`
	Build        string `json:"build"`
}

type UserKey struct{}

type ServerStats struct {
	IsRunning bool `json:"is_running"`
	MemUse    uint `json:"mem_use"`
}
