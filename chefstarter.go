package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var (
	username *string
	password *string
)

func handler(w http.ResponseWriter, r *http.Request) {
	user, pass, authorized := r.BasicAuth()

	if user != *username || pass != *password {
		log.Println("Client", r.RemoteAddr, "failed HTTP basic authentication. Username:", user, "Password:", pass)
		return
	}

	if authorized {
		log.Println("Launched chef-client by request of", r.RemoteAddr)
		cmd := exec.Command("sudo", "chef-client")
		err := cmd.Start()
		if err != nil {
			log.Println("chefstarter ERROR:", err)
		}
		fmt.Fprintf(w, "chef-client launched")
	}
}

func main() {
	username = flag.String("u", "chefstarter", "HTTP Basic Auth username")
	password = flag.String("p", "", "HTTP Basic Auth password")

	flag.Parse()

	if len(*password) == 0 {
		log.Fatalln("Must provide a password via -p flag.  Use -h for help.")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":7100", nil)
}
