package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	// Pull creds that were passed through HTTP Basic Authentication
	user, pass, _ := r.BasicAuth()

	if user != *username || pass != *password {
		log.Println("chefstarter ERROR: Client", r.RemoteAddr, "failed HTTP basic authentication. Username:", user, "Password:", pass)
		http.Error(w, "Invalid username/password combination", 403)
	} else {
		// HTTP Basic Auth succeeded, so we continue...

		log.Println("Launching chef-client by request of", r.RemoteAddr)

		cmd := exec.Command("sudo", "chef-client")
		// For testing chef-client exit codes...
		// cmd := exec.Command("./error.sh")

		// Send chef-client's stdout and stderr to chefstarter's stdout and stderr.
		// If running from supervisord, for example, this can be collected and logged.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if *synch {
			// If we're running synchronously, block on execution and wait for chef-client to finish
			err = cmd.Run()
		} else {
			// Otherwise, just execute it and move on
			err = cmd.Start()
		}
		if err != nil {
			// We were unable to execute chef-client or it executed but returned non-zero exit code
			friendlyErr := "Error executing chef-client: " + err.Error()
			log.Println("chefstarter ERROR: chef-client execution returned error", err.Error())
			http.Error(w, friendlyErr, 500)
		}

		// chef-client ran properly so return a friendly message to the client
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

	// Set up our handler at the designated path
	http.HandleFunc(*path, handler)

	// If -ssl=true flag was passed, we use TLS
	if *usessl {
		if len(*keyfile) == 0 || len(*crtfile) == 0 {
			log.Fatalln("For HTTPS, you must specify key and cert file.  Use -h for help.")
		}

		err := http.ListenAndServeTLS(*listen, *crtfile, *keyfile, nil)
		if err != nil {
			log.Println("chefstarter ERROR:", err)
		}

	} else {

		// No TLS so we use regular old HTTP
		http.ListenAndServe(*listen, nil)

	}

}
