// This file contains the main routine for workers.
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/showalter/bdws/internal/data"
)

type codeFunction func([]byte) []byte
type extension string

// Map various extension names to their code
var extensionMap = map[extension]codeFunction{
	".sh":    bashScript,
	".py":    nil,
	".java":  nil,
	".class": nil,
	".jar":   nil,
	"":       nil,
}

// run the code given an extension
func runCode(e extension, code []byte) []byte {
	f, found := extensionMap[e]
	if found {
		return f(code)
	} else {
		return []byte("Error: Extension not found.")
	}
}

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

	// Run the code and get []byte output
	output := runCode(".sh", job.Code)

	// Print out the json.
	fmt.Println(jobJson)

	// Send a response back.
	w.Write(output)
}

// The entry point of the program.
func main() {

	// The command line arguments. args[0] is the port to run on.
	args := os.Args

	// If there was no argument passed, ask for one and exit.
	if len(args) == 1 {
		fmt.Println("Please pass a port number. Eg. :38471")
		os.Exit(1)
	}

	// If there is a request for /newjob,
	// the new_job routine will handle it.
	http.HandleFunc("/newjob", new_job)

	// Listen on a port.
	http.ListenAndServe(args[1], nil)
}

/* Code Strategies */

func bashScript(code []byte) []byte {
	// Create a temporary file
	// TODO: Make this run with various extensions
	scriptName := "tmp.sh"
	file, err := os.Create(scriptName)
	check(err)

	_, err = file.Write(code)
	check(err)

	file.Sync()
	file.Close()

	// Make temp file executable.
	check(os.Chmod(scriptName, 0700))

	// Execute temp file and print output.
	cmd := run(("./" + scriptName), "")

	// Remove temp file.
	os.Remove(scriptName)

	return cmd
}
