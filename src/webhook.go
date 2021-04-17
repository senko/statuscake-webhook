/* Simple webhook that executes a command on incoming StatusCake DOWN notification
 *
 * Written by Senko Rasic <senko@senko.net>
 * Released to Public Domain. Copyright is for chums.
 *
 * See https://www.statuscake.com/kb/knowledge-base/how-to-use-the-web-hook-url/
 *
 * Usage:
 *
 *   1. Set environment variable PORT to a HTTP port to listen to (eg. 9000). Make
 *      sure the service is reachable (not firewalled) on that port.
 *
 *   2. Set environment variable TOKEN to the token StatusCake will send along
 *      the notification, to verify the sender. See above URL for details on
 *      how the token is generated, or send a test notification to webhook.site and
 *      copy/paste the token from there.
 *
 *   3. Run the webhook service with command that needs to be executed when the
 *      webhook is triggered
 *
 *   4. Relax and enjoy
 *
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const DEFAULT_PORT = 9000
var statuscakeToken = ""
var command = []string{}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	r.ParseForm()

	token := r.PostFormValue("Token")
	if token != statuscakeToken {
		fmt.Println("Warning: received request with incorrect token:", token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	name := r.PostFormValue("Name")
	ip := r.PostFormValue("IP")

	status := strings.ToLower(r.PostFormValue("Status"))
	if status != "down" {
		fmt.Printf("[%s/%s] Host status %s, nothing to do\n", name, ip, status)
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Printf("[%s/%s] Host status is DOWN, executing command: %s\n", name, ip, command)

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		fmt.Printf("[%s/%s] Error executing command: %s\n", name, ip, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getConfig() (port int, token string) {
	port = DEFAULT_PORT
	if providedPort, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		port = providedPort
	}

	token = os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is not set")
	}

	command = os.Args[1:]
	if len(command) == 0 {
		panic("Command should be provided")
	}
	return
}


func main() {
	port, token := getConfig()

	statuscakeToken = token

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

