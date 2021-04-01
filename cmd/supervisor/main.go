// This file contains the main routine for supervisors.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/showalter/bdws/internal/data"
)

type Worker struct {
	Id       int64
	Busy     bool
	Hostname string
}

// TODO: Make the workers register with the supervisor on startup.
// This way, we don't have to hard code the workers.
var workers = [2]Worker{Worker{1, false, "http://127.0.0.1:39481"},
	Worker{2, false, "http://127.0.0.1:39482"}}

// Handle the submission of a new job.
func new_job(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling connection...")

	// Parse the HTTP request.
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	// Put the bytes from the request into a file
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	job := data.JsonToJob(buf)

	// Set up argument information.
	argRange := job.ParameterEnd - job.ParameterStart
	argInterval := argRange / int64(len(workers))

	argStart := job.ParameterStart
	argEnd := argStart + argInterval

	var responses []byte

	// Send the job to each worker
	for _, w := range workers {
		job.ParameterStart = argStart
		job.ParameterEnd = argEnd

		argStart += argInterval
		argEnd += argInterval

		jobBytes := data.JobToJson(job)

		resp, err := http.Post(w.Hostname+"/newjob",
			"text/plain", bytes.NewReader(jobBytes))
		if err != nil {
			panic(err)
		}

		// TODO: Collect all responses into one response to send back to the client.
		// Put the bytes from the request into a file
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		responses = append(responses, buf.Bytes()...)
	}

	// Send a response back.
	w.Write(responses)
}

// The entry point of the program.
func main() {

	// If there is a request for /newjob,
	// the new_job routine will handle it.
	http.HandleFunc("/newjob", new_job)

	// Listen on a port.
	http.ListenAndServe(":39480", nil)
}
