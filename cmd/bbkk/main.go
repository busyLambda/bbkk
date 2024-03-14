package main

import "github.com/busyLambda/bbkk/internal/api"

func main() {
	app := api.NewApiMaster()

	app.AttachRoutes()

	app.Run()
	/*
		var wg sync.WaitGroup

		mcServer := server.NewMcServer("server", "paper.jar", "")

		wg.Add(1)
		go mcServer.Start(&wg)

		outchan := make(chan string)

		if mcServer.Stdout == nil {
			fmt.Printf("No stdout, setting it...\n")
			mcServer.SetStdout()
		}

		go mcServer.ReadStdout(outchan)

		go func() {
			for {
				select {
				case data := <-outchan:
					fmt.Printf(data)
				}
			}
		}()

		time.Sleep(1 * time.Second)
		mcServer.WriteString("stop\n")

		wg.Wait()
	*/
}
