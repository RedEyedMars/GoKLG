package networking

import (
	"Events"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const GET = "GET"

var addr = flag.String("addr", ":8080", "http service address")
var Shutdown chan bool
var homeHtml string

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf(" networking.web.serveHome:%s", r.URL.String())
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, homeHtml)
	//http.ServeFile(w, r, "src/Networking/home.html")
}

func handleImg(imgName string) {
	http.HandleFunc(imgName, func(w http.ResponseWriter, r *http.Request) {

		log.Printf(" networking.web.handleImg:%s", r.URL.String())
		if r.Method != GET {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "assets/imgs"+r.URL.String())
	})
}
func handleJs(libName string) {
	http.HandleFunc(libName, func(w http.ResponseWriter, r *http.Request) {
		log.Printf(" networking.web.handleJs:%s", r.URL.String())
		if r.Method != GET {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "src/www/js"+r.URL.String())
	})
}

var onClose func()

func End() {
	if onClose != nil {
		Events.FuncEvent("Networking.End", onClose)
	}
}

func StartWebClient(toClose chan bool) {
	Shutdown = toClose
	setupNetworkingRegex()
	homeRaw, err := ioutil.ReadFile("src/www/home.html")
	if err != nil {
		log.Fatalf("networking.web.StartWebClient:%s", err)
	}
	homeHtml = string(homeRaw)

	flag.Parse()
	registry := newRegistry()
	go registry.run()

	setupClientCommands(registry)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(registry, w, r)
	})
	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {

		log.Printf(" networking.web.GetStyleSheet:%s", r.URL.String())
		if r.Method != GET {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "src/www"+r.URL.String())
	})
	http.HandleFunc("/forge-sha256.min.js", func(w http.ResponseWriter, r *http.Request) {
		log.Printf(" networking.web.GetSha256:%s", r.URL.String())
		if r.Method != GET {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "src/www/lib/forge-sha256-master/build/forge-sha256.min.js")
	})
	handleJs("/login.js")
	handleImg("/Pending.jpg")
	handleImg("/Fail.jpg")
	handleImg("/Success.jpg")

	srv := &http.Server{Addr: ":8080"}
	Events.GoFuncEvent("Networking.ListenAndServe", func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("networking.web.ListenAndServer:%s", err)
		}
	})
	onClose = func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("networking.web.Shutdown:%s", err)
		}
	}
	go func() {
		time.Sleep(24 * time.Hour)
		close(Shutdown)
	}()
}
