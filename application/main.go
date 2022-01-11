package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	dir := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	log.Printf("Serving %s on HTTP port: %s\n", *dir, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
