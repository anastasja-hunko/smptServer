package main

import (
	"github.com/anastasja-hunko/smptServer/internal"
	"log"
)

func main() {
	config := internal.NewConfig()

	server := internal.New(config)

	log.Fatal(server.Start())

	//ticker := time.NewTicker(20 * time.Second)
	//done := make(chan bool)
	//
	//go func() {
	//	for {
	//		select {
	//		case <-done:
	//			return
	//		case t := <-ticker.C:
	//			_, err := http.Get("localhost:" + config.Port)
	//
	//			if err != nil {
	//				log.Println("it doen't work at:", t)
	//
	//				go server.Start()
	//
	//			} else {
	//				log.Println("it works at:", t)
	//			}
	//		}
	//	}
	//}()
	//
	//time.Sleep(5 * time.Minute)
	//ticker.Stop()
	//done <- true
	//fmt.Println("Ticker stopped")
	//os.Exit(0)
}
