package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/busyLambda/bbkk/internal/util"
)

type McServer struct {
	Cmd       *exec.Cmd
	Stdout    *bufio.Scanner
	Stdin     io.WriteCloser
	streaming bool
	Wg        sync.WaitGroup
	isRunning bool
}

func NewMcServer(dir string, jar string, flags string) *McServer {
	c := util.JavaCmd(dir, jar, flags)

	return &McServer{
		Cmd:       c,
		isRunning: false,
		streaming: false,
	}
}

func (ms *McServer) IsStreaming() bool {
	return ms.streaming
}

func (ms *McServer) SetStdout() error {
	sp, err := ms.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	ms.Stdout = bufio.NewScanner(sp)
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

func (ms *McServer) IsRunning() bool {
	return ms.isRunning
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

	ms.isRunning = true
	err = ms.Cmd.Start()
	if err != nil {
		ms.isRunning = false
		log.Printf("Error starting java: %s\n", err.Error())
	}

	ms.Cmd.Wait()

	ms.isRunning = false
}

func (ms *McServer) StopServer() {
	ms.StopStdout()

	// TODO: Have to like time it so that we only set it to false if it's really not running.
	ms.WriteString("stop\n")

	var wg sync.WaitGroup
	wg.Add(1)

	// Wait till the Jar process actually stops.
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			if !ms.IsRunning() {
				break
			}
		}
	}(&wg)

	wg.Wait()

	// Reset the cmd to be able to run it again later.
	ms.Cmd = util.JavaCmd(ms.Cmd.Dir, "server.jar", "")
	ms.isRunning = false
}

func (ms *McServer) StopStdout() {
	ms.streaming = false
}

func (ms *McServer) ReadStdout(output chan<- string) {
	defer close(output)

	if ms.Cmd.ProcessState != nil {
		if ms.Cmd.ProcessState.Exited() {
			return
		}
	}

	ms.streaming = true

	for ms.Stdout.Scan() {
		text := ms.Stdout.Text()
		fmt.Println(text)
		// output <- text
	}

	if err := ms.Stdout.Err(); err != nil {
		log.Println("Error reading from pipe: ", err)
	}
}

func (ms *McServer) WriteRune(r rune) (err error) {
	_, err = ms.Stdin.Write([]byte{byte(r)})
	return
}

func (ms *McServer) WriteString(s string) (err error) {
	_, err = ms.Stdin.Write([]byte(s))
	return
}
