package server

import (
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/busyLambda/bbkk/internal/util"
)

type McServer struct {
  Cmd *exec.Cmd
  Stdout io.ReadCloser
  Wg sync.WaitGroup
}

func NewMcServer(dir string, jar string, flags string) McServer {
  c := util.JavaCmd(dir, jar, flags)
  
	return McServer{
    Cmd: c,
	}
}

func (ms *McServer) Start(wg *sync.WaitGroup) {
  defer wg.Done()

  sp, err := ms.Cmd.StdoutPipe()
  ms.Stdout = sp
  if err != nil {
    println("Error creating stdout pipe.")
  }

	ms.Cmd.Start()

  outchan := make(chan string)

  go ms.ReadStdout(outchan)

  var wg_internal sync.WaitGroup
  wg_internal.Add(1)

  go func() {
    defer wg_internal.Done()
    for {
      select {
        case data := <-outchan:
        fmt.Printf(data)
      }
    }
  }()

  ms.Cmd.Wait()
  wg_internal.Wait()
}

func (ms *McServer) ReadStdout(output chan<-string) {
  if ms.Cmd.ProcessState != nil {
    if ms.Cmd.ProcessState.Exited() {
      return
    }
  }

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

// TODO: Write the function :3
func WriteStdin(s string) {}
