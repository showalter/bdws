// This file contains the main routine for supervisors.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/showalter/bdws/internal/data"
)

type Worker struct {
	Id       int64
	Busy     bool
	Hostname string
}

var workers []data.Worker
var workerCounter int64 = 1

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

	// array to store responses
	var responses []byte

	// Send the job to each worker
	for _, w := range workers {
		job.ParameterStart = argStart
		job.ParameterEnd = argEnd

		argStart += argInterval
		argEnd += argInterval

		jobBytes := data.JobToJson(job)

		resp, err := http.Post("http://"+w.Hostname+"/newjob",
			"text/plain", bytes.NewReader(jobBytes))
		if err != nil {
			panic(err)
		}

		// Put the bytes from the request into a file
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		responses = append(responses, buf.Bytes()...)
	}

	// Send a response back.
	w.Write(responses)
}

// Look through a list and add the item if it isn't in the list already.
// This is slow for big lists, but there won't likely be a large number of workers.
func appendIfUnique(list []data.Worker, w data.Worker) []data.Worker {
	for _, x := range list {
		if x.Hostname == w.Hostname {
			return list
		}
	}

	return append(list, w)
}

// Handle a worker registering to receive work
func register(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling registration...")

	// Parse the HTTP request.
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	// The worker will send the port they'll be listening on
	port := buf.String()

	// Replace the port number for the sender
	split := strings.Split(req.RemoteAddr, ":")
	split[len(split) - 1] = port

	worker := data.Worker{Id: workerCounter, Busy: false,
		Hostname: strings.Join(split, ":")}

	workerCounter += 1

	// We don't need multiple workers with the same hostname.
	workers = appendIfUnique(workers, worker)

	w.Write(data.WorkerToJson(worker))
}

// The entry point of the program.
func main() {

	// The command line arguments. args[0] is the port to run on.
	args := os.Args

	// If the right number of arguments weren't passed, ask for them.
	if len(args) != 2 {
		fmt.Println("Please pass the port to run on preceded by a colon. eg. :4001")
		os.Exit(1)
	}

	// If there is a request for /newjob,
	// the new_job routine will handle it.
	http.HandleFunc("/newjob", new_job)

	// Handle requests for /register with the register function
	http.HandleFunc("/register", register)

	// Listen on a port.
	http.ListenAndServe(args[1], nil)
}
