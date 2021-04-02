// This file contains the main routine for clients.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/showalter/bdws/internal/data"
)

// The entry point of the program
func main() {

	// The command line arguments. args[0] is the name of the program.
	args := os.Args

	// If the right number of arguments weren't passed, ask for them and exit.
	if len(args) != 3 {
		fmt.Println("Please pass the address of the supervisor and a file to run.")
		fmt.Println("Example: http://stu.cs.jmu.edu:4001 fun_code.py")
		os.Exit(1)
	}

	// Open the file whose name was passed as an argument.
	code, err := ioutil.ReadFile(args[2])
	if err != nil {
		panic(err)
	}

	// Make a job with the given code.
	jobBytes := data.JobDataToJson(1, time.Now(), 2, 1, 10, code)

	// Send a post request to the worker.
	resp, err := http.Post(args[1]+"/newjob",
		"text/plain", bytes.NewReader(jobBytes))
	if err != nil {
		panic(err)
	}

	// Put the bytes from the request into a file
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	file := buf.String()

	fmt.Println(file)

}
