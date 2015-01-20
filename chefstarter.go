package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var (
	path     *string
	synch    *bool
	username *string
	password *string
	keyfile  *string
	crtfile  *string
	usessl   *bool
	listen   *string
)

func handler(w http.ResponseWriter, r *http.Request) {
	var err error

	user, pass, authorized := r.BasicAuth()

	if user != *username || pass != *password {
		log.Println("Client", r.RemoteAddr, "failed HTTP basic authentication. Username:", user, "Password:", pass)
		http.Error(w, "Invalid username/password combination", 403)
	}

	if authorized {
		log.Println("Launched chef-client by request of", r.RemoteAddr)

		cmd := exec.Command("sudo", "chef-client")
		// For testing chef-client exit codes...
		// cmd := exec.Command("./error.sh")

		if *synch {
			err = cmd.Run()
		} else {
			err = cmd.Start()
		}
		if err != nil {
			friendlyErr := "Error executing chef-client: " + err.Error()
			log.Println("ERROR: chef-client execution returned error", err.Error())
			http.Error(w, friendlyErr, 500)
		}
		fmt.Fprintf(w, "chef-client launched")
		return
	}
	return
}

func main() {
	path = flag.String("path", "/", "Request path to initiate chef-client run (Default: /)")
	synch = flag.Bool("wait", false, "Wait until chef-client completes run before returning HTTP response")
	username = flag.String("user", "chefstarter", "HTTP Basic Auth username (Default: chefstarter)")
	password = flag.String("pass", "", "HTTP Basic Auth password")
	keyfile = flag.String("key", "", "HTTPS X.509 Private Key file")
	crtfile = flag.String("cert", "", "HTTPS X.509 Public Certificate file")
	usessl = flag.Bool("ssl", false, "Enable HTTPS (true/false)")
	listen = flag.String("listen", ":7100", "IP:port to listen on (Default: listen on all interfaces on port 7100)")

	flag.Parse()

	if len(*password) == 0 {
		log.Fatalln("Must provide a HTTP Basic Auth password.  Use -h for help.")
	}

	http.HandleFunc(*path, handler)

	if *usessl {
		if len(*keyfile) == 0 || len(*crtfile) == 0 {
			log.Fatalln("For HTTPS, you must specify key and cert file.  Use -h for help.")
		}

		err := http.ListenAndServeTLS(*listen, *crtfile, *keyfile, nil)
		if err != nil {
			log.Println("ERROR:", err)
		}

	} else {

		http.ListenAndServe(*listen, nil)

	}

}
