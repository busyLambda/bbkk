package server

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/busyLambda/bbkk/internal/util"
)

type McServer struct {
	Cmd    *exec.Cmd
	Stdout io.ReadCloser
	Stdin  io.WriteCloser
	Wg     sync.WaitGroup
}

func NewMcServer(dir string, jar string, flags string) *McServer {
	c := util.JavaCmd(dir, jar, flags)

	return &McServer{
		Cmd: c,
	}
}

func (ms *McServer) SetStdout() error {
	sp, err := ms.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	ms.Stdout = sp
	return nil
}

func (ms *McServer) SetStdin() error {
	sp, err := ms.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	ms.Stdin = sp
	return nil
}

func (ms *McServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	err := ms.SetStdout()
	if err != nil {
		fmt.Printf("Error with le stdout pipe: %s\n", err)
	}

	err = ms.SetStdin()
	if err != nil {
		fmt.Printf("Error with le stdin pipe: %s\n", err)
	}

	err = ms.Cmd.Start()
	if err != nil {
		fmt.Println("Error starting java.")
	}

	ms.Cmd.Wait()
}

func (ms *McServer) ReadStdout(output chan<- string) {
	log.Printf("Checking process...")
	if ms.Cmd.ProcessState != nil {
		if ms.Cmd.ProcessState.Exited() {
			return
		}
	}
	log.Printf("Streaming :3")

	buf := make([]byte, 1024)
	for {
		n, err := ms.Stdout.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Print("-> Error reading from stdout :<")
			}
			break
		}
		output <- string(buf[:n])
	}
}

func (ms *McServer) WriteStdin(r rune) (err error) {
	_, err = ms.Stdin.Write([]byte{byte(r)})
	return
}

func (ms *McServer) WriteString(s string) (err error) {
	_, err = ms.Stdin.Write([]byte(s))
	return
}
