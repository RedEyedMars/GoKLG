package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"

	"./src/databasing"
	"./src/events"
	"./src/networking"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	args := os.Args
	if len(args) <= 1 {
		MainStart("main.Run", func(Shutdown chan bool) {
			databasing.Run(Shutdown)
			networking.Run(Shutdown)
		},
			func(msg string) bool {
				if !networking.HandleAdminCommand(msg) {
					return databasing.HandleAdminCommand(msg)
				}
				return true

			}, func() {
				databasing.End()
				networking.End()
			})
	} else {
		switch args[1] {
		case "web":
			MainStart("networking.Run", networking.Run, networking.HandleAdminCommand, networking.End)
		case "database":
			MainStart("databasing.Setup", databasing.Run, databasing.HandleAdminCommand, databasing.End)
		}
	}
}

func MainStart(name string, f func(chan bool), adminCommand func(string) bool, end func()) {
	Shutdown := make(chan bool, 1)
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if text[:4] == "exit" {
				Shutdown <- true
				break
			} else {
				log.Printf(text)
				adminCommand(text)
			}
		}
	}()
	events.DoneFuncEvent(name, f, Shutdown)
	log.Printf(" Waiting for done")
	<-Shutdown
	end()
}
