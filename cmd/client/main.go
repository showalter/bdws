// This file contains the main routine for clients.
package main

import (
	"bytes"
	"net/http"
	"os"
	"fmt"
)

// The entry point of the program
func main() {

	// The command line arguments. args[0] is the name of the program.
	args := os.Args

	// If there was no argument passed, ask for one and exit.
	if len(args) == 1 {
		fmt.Println("Please pass a file name to send to the server.")
		os.Exit(1)
	}

	// Open the file whose name was passed as an argument.
	dat, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}

	// Send a post request to the worker.
	resp, err := http.Post("http://127.0.0.1:39485/newjob", "text/plain", dat)
	if err != nil {
		panic(err)
	}
	
	// Put the bytes from the request into a file
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	file := buf.String()

	fmt.Println(file)

}
