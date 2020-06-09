package main

import (
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal"
	"log"
	"net/http"
	"time"
)

func main() {

	config := internal.NewConfig()

	server := internal.New(config)

	ticker := time.NewTicker(5 * time.Second)

	done := make(chan bool)

	go func() {

		for {
			select {

			case <-done:
				return

			case t := <-ticker.C:

				_, err := http.Get("http://127.0.0.1:" + config.Port)

				if err != nil {
					fmt.Println("it doen't work at:", t)

					go func() {

						err := server.Start()
						if err != nil {
							log.Fatal(err)
						}

					}()

				} else {
					fmt.Println("it works at:", t)
				}
			}
		}
	}()

	time.Sleep(20 * time.Minute)

	ticker.Stop()

	done <- true

}
