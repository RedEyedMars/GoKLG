package main

import (
	"Events"
	"bufio"
	"databasing"
	"math/rand"
	"networking"
	"os"
	"time"
)

func Run(Shutdown chan bool) {
	Events.GoFuncEvent("Networking.StartWebClient", func() {
		Networking.StartWebClient(Shutdown)
	})
}
func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	args := os.Args
	if len(args) <= 1 {
		MainStart("main.Run", func(Shutdown chan bool) {
			databasing.Run(Shutdown)
			Run(Shutdown)
		},
			func(msg string) bool {
				if !Networking.HandleAdminCommand(msg) {
					return databasing.HandleAdminCommand(msg)
				} else {
					return true
				}
			}, func() {
				databasing.End()
				Networking.End()
			})
	} else {
		switch args[1] {
		case "chat_service":
			MainStart("main.Run", Run, Networking.HandleAdminCommand, Networking.End)
		case "setup_database":
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
				adminCommand(text)
			}
		}
	}()
	Events.DoneFuncEvent(name, f, Shutdown)
	<-Shutdown
	end()
}
