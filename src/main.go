package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/buildkite/terminal"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadFile("static/index.html")
	fmt.Fprintf(w, string(content))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	tmpfile, err := ioutil.TempFile("", "ansible-sandbox")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to create temp file for playbook: %s", err)
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(r.FormValue("content")))
	tmpfile.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "env", "ANSIBLE_FORCE_COLOR=1", "ansible-playbook", "-i", "localhost,", "-c", "local", tmpfile.Name())
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		w.WriteHeader(500)
		fmt.Fprintf(w, string(terminal.Render(out)))
		fmt.Fprintf(w, "\n\n***** Timeout exceeded *****\n")
	} else {
		fmt.Fprintf(w, string(terminal.Render(out)))
	}
}

func main() {
	static_dir := flag.String("static", "./static", "path to static assets")
	listen_addr := flag.String("listen", ":8080", "address and port to listen on")
	flag.Parse()
	http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir(*static_dir))))
	http.HandleFunc("/index.html", indexHandler)
	http.HandleFunc("/submit", submitHandler)
	fmt.Printf("Listening on %s\n", *listen_addr)
	log.Fatal(http.ListenAndServe(*listen_addr, nil))
}
