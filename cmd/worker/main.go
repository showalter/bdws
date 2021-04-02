// This file contains the main routine for workers.
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/showalter/bdws/internal/data"
)

// Check for an error.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Run a given command.
func run(command string, args string) []byte {
	cmd, err := exec.Command(command, args).Output()
	check(err)
	return cmd
}

// Handle the submission of a new job.
func new_job(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling connection...")

	// Parse the HTTP request.
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	// Put the bytes from the request into a file
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	jobJson := buf.String()

	// Convert string json to job struct
	job := data.JsonToJob([]byte(jobJson))

	// Create a temporary file
	// TODO: Make this run with various extensions
	scriptName := "tmp.sh"
	file, err := os.Create(scriptName)
	check(err)

	_, err = file.Write(job.Code)
	check(err)

	file.Sync()
	file.Close()

	// Make temp file executable.
	check(os.Chmod(scriptName, 0700))

	// Execute temp file and print output.
	cmd := run(("./" + scriptName), "")

	// Remove temp file.
	os.Remove(scriptName)

	// Print out the json.
	fmt.Println(jobJson)

	// Send a response back.
	w.Write(cmd)
}

// The entry point of the program.
func main() {

	// The command line arguments. args[1] is the supervisor address,
	// args[2] is the port to run on
	args := os.Args

	// If the right number of arguments weren't passed, ask for them.
	if len(args) != 3 {
		fmt.Println("Please pass the hostname of the supervisor and the outgoing port." +
			"eg. http://stu.cs.jmu.edu:4001 4031")
		os.Exit(1)
	}

	resp, err := http.Post(args[1] + "/register", "text/plain", strings.NewReader(args[2]))
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// This gives what the supervisor thinks the worker is, which is useful for debugging.
	_ = data.JsonToWorker(buf.Bytes())

	// If there is a request for /newjob,
	// the new_job routine will handle it.
	http.HandleFunc("/newjob", new_job)

	// Listen on a port.
	log.Fatal(http.ListenAndServe(":"+args[2], nil))
}
