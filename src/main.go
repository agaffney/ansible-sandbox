package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("static/index.html")
	fmt.Fprintf(w, string(content))
}

func main() {
	static_dir := flag.String("static", "./static", "path to static assets")
	flag.Parse()
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir(*static_dir))))
	http.HandleFunc("/index.html", indexHandler)
	fmt.Printf("Listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
