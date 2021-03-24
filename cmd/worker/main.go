// This file contains the main routine for workers.
package main

import (
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

	// Print out the file.
	fmt.Printf(file)

	// Send a response back.
	w.Write([]byte("Done"))
}

// The entry point of the program.
func main() {

	// If there is a request for /newjob,
	// the new_job routine will handle it.
	http.HandleFunc("/newjob", new_job)

	// Listen on a port.
	http.ListenAndServe(":39485", nil)
}
