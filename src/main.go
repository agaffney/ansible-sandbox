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

var Config struct {
	staticDir   string
	listenAddr  string
	dockerImage string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		content, _ := ioutil.ReadFile("static/index.html")
		fmt.Fprintf(w, string(content))
	} else {
		w.WriteHeader(404)
	}
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
	tmpfile.Chmod(0444)
	tmpfile.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx,
		"docker", "run", "--rm",
		// Set various env vars for Ansible
		"-e", "ANSIBLE_FORCE_COLOR=1",
		"-e", "ANSIBLE_RETRY_FILES_ENABLED=0",
		"-e", "ANSIBLE_LOCAL_TEMP=/tmp",
		// Disable networking
		"--network", "none",
		// Use non-root user
		"--user", "nobody",
		// Map temp file into container
		"-v", fmt.Sprintf("%s:%s", tmpfile.Name(), tmpfile.Name()),
		Config.dockerImage,
		// Run command through shell
		"sh", "-c",
		fmt.Sprintf("ansible-playbook -i localhost, -c local %s", tmpfile.Name()),
	)
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
	flag.StringVar(&Config.staticDir, "static", "./static", "path to static assets")
	flag.StringVar(&Config.listenAddr, "listen", ":8080", "address and port to listen on")
	flag.StringVar(&Config.dockerImage, "docker-image", "ansible-sandbox", "docker image to use for running ansible")
	flag.Parse()
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(Config.staticDir))))
	http.HandleFunc("/submit", submitHandler)
	fmt.Printf("Listening on %s\n", Config.listenAddr)
	log.Fatal(http.ListenAndServe(Config.listenAddr, nil))
}
