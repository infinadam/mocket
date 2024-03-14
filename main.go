package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.String("p", "80", "Port to listen on.")
var scriptDir = flag.String("s", "./scripts", "Script directory.")

func main() {
	flag.Parse()

	log.Printf("mocket: reading script directory (%s)...\n", *scriptDir)
	var server *Server
	var err error
	if server, err = MakeServer(*scriptDir); err != nil {
		fmt.Printf("mocket: error making server (%v)", err)
		return
	}

	log.Printf("mocket: starting on (%s)...\n", *port)
	http.HandleFunc("/", server.HandleRequest)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
