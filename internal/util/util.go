package util

import (
	"log"
	"os/exec"
	"strings"
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
