package main

import (
	"log"
	"net/http"
	"time"
)

func MarkdownRenderHandler(w http.ResponseWriter, r *http.Request) {
	s := "dynamic data"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func main() {
	// Primary servermux routes
	routes := http.NewServeMux()

	// Static server
	fs := http.FileServer(http.Dir("static"))
	// Init routes
	routes.Handle("/", fs)
	routes.Handle("/js/", fs)
	routes.Handle("/css/", fs)
	routes.Handle("/fonts/", fs)
	routes.Handle("/images/", fs)
	routes.HandleFunc("/markdown/", MarkdownRenderHandler)
	
	log.Println("Listening on port: 80")

	// Create Server
	server := http.Server{
		Addr: ":http",
		Handler: routes,
		ReadTimeout: time.Second*10,
		WriteTimeout: time.Second*10,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
