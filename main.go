package main

import (
	"flag"
	"log"
	"net/http"
)

var port = flag.String("p", "80", "Port to listen on.")
var scriptDir = flag.String("s", "./scripts", "Script directory.")

func main() {
	flag.Parse()

	log.Printf("mocket: reading script directory (%s)...\n", *scriptDir)
	if server, err := MakeServer(*scriptDir); err != nil {
		log.Fatalf("mocket: error making server (%v)", err)
	} else {
		log.Printf("mocket: starting on (%s)...\n", *port)
		http.HandleFunc("/", server.HandleRequest)
		http.ListenAndServe(":"+*port, nil)
	}
}
