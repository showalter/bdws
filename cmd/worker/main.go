// This file contains the main routine for workers.
package main

import (
	"os"
	"bytes"
	"net/http"
	"fmt"
)

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
	file := buf.String()

	// TODO: Run the job with the given arguments, sending the results back

	// Print out the file.
	fmt.Printf(file)

	// Send a response back.
	w.Write([]byte("Done"))
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
