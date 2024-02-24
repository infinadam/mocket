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
  // not really...

	log.Printf("mocket: starting on (%s)...\n", *port)
	http.HandleFunc("/", handlerFunction)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func handlerFunction(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
