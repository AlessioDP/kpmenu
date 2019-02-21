package main

import (
	"log"
	"os"

	"github.com/AlessioDP/kpmenu/kpmenulib"
)

func main() {
	menu := kpmenulib.Initialize()

	if menu != nil {
		// Start client
		err := kpmenulib.StartClient()
		if err != nil {
			// Failed to comunicate with server - start server
			err = kpmenulib.StartServer(menu)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			} else {
				log.Printf("waiting for goroutines to end")
				// Wait for any goroutine (clipboard)
				menu.WaitGroup.Wait()
			}
		}
	}
}
