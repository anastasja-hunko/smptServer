package main

import (
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal"
	"net/http"
	"time"
)

func serverCheck(server *internal.Server) {

	ticker := time.NewTicker(5 * time.Second)

	for t := range ticker.C {

		_, err := http.Get("http://127.0.0.1" + server.Config.Port)

		if err != nil {

			fmt.Println("it doesn't work at:", t)

			go server.Restart()

			select {}

		} else {

			fmt.Println("it works at:", t)

		}

	}

	time.Sleep(30 * time.Second)

	ticker.Stop()

}

func main() {
	config := internal.NewConfig()

	server := internal.New(config)

	go server.Start()

	go serverCheck(server)

	select {}
}
