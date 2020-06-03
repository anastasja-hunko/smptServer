package main

import (
	"fmt"
	"github.com/anastasja-hunko/smptServer/internal"
	"log"

	//"log"
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
				fmt.Println("Tick at", t)
				if !server.IsAlive() {
					log.Println("is not work")
					err := server.Start()
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("it works")
				}
			}
		}
	}()

	time.Sleep(time.Minute)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")

	//go func() {
	//	for {
	//		select {
	//		case t := <-ticker.C:
	//			fmt.Println("Tick at", t)
	//			//if !server.IsAlive() {
	//			//	log.Println("is not work")
	//			//	err := server.Start()
	//			//	if err != nil {
	//			//		log.Fatal(err)
	//			//	}
	//			//} else {
	//			//	log.Println("it works")
	//			//}
	//		}
	//
	//	}
	//
	//}()
	//init config, and server. Then start server

	//sessionStore := sessions.NewCookieStore([]byte("very-secret-key"))

}
