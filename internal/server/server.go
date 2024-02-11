package server

import (
	"bufio"
	"fmt"
	"os/exec"
)

type McServer struct {
	Path  string
	Flags string
}

func NewMcServer(path string, flags string) McServer {
	return McServer{
		Path:  path,
		Flags: flags,
	}
}

// TODO: Make this multithreaded.
func (ms *McServer) Start() {
	cmd := exec.Command("java", "-jar", ms.Flags)
	cmd.Dir = ms.Path

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	cmd.Wait()
}
