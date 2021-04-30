// This file contains the main routine for clients.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/showalter/bdws/internal/data"
)

// The entry point of the program
func main() {

	// Declare variables
	var hostName string
	var fullFileName string
	var start int
	var end int
	var args []string
	var runs int

	// Parse command line
	parseCommandLine(&hostName, &fullFileName, &start, &end, &args, &runs)

	// Get extension and file name
	fileName, extension := getFileName(fullFileName)

	// Code is unessesary to send if executable exists
	var code []byte
	if extension == "system program" {
		code = nil
	} else {
		// File is not an binary executable, so copy code
		// Open the file whose name was passed as an argument.
		var err error
		code, err = ioutil.ReadFile(fullFileName)
		if err != nil {
			fmt.Println("Error opening file. Aborting")
			os.Exit(3)
		}
	}

	// Make a job with the given code.
	jobBytes := data.JobDataToJson(1, time.Now(), 2, start, end, fileName, extension, code, args, runs)

	// Send a post request to the supervisor.
	resp, err := http.Post(hostName+"/newjob",
		"text/plain", bytes.NewReader(jobBytes))
	if err != nil {
		fmt.Println("Error posting job. Aborting")
		os.Exit(3)
	}

	// Put the bytes from the request into a file
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	file := buf.String()

    fmt.Println(file)

}

/* ----- Helper functions ----- */
func parseCommandLine(hostname *string, fullFileName *string, start *int, end *int, args *[]string, runs *int) {
	// Optional flags
	argsPtr := flag.String("args", "NONE", "Command line args for file\nExample: -args \"-al\" when running ls")
	rangePtr := flag.String("range", "NONE", "Range for job\nExample: -range 1-10")
	runsPtr := flag.Int("runs", 1, "Number of times to run job")
	flag.Parse()

	*args = strings.Split(*argsPtr, " ")

	// Non optional command line argsgi
	tail := flag.Args()

	// If the right number of arguments weren't passed, ask for them and exit.
	if len(tail) < 2 {
		fmt.Println("Please pass the address of the supervisor and a file to run, and an optional range of parameters.")
		fmt.Println("\tExample: {optional flags} http://stu.cs.jmu.edu:4001 fun_code.py")
		fmt.Println("\tRun ./client -h for more info on optional flags")
		os.Exit(1)
	} else {
		*hostname = tail[0]
		*fullFileName = tail[1]
	}

	// Set the range
	// A start index greater than the end index indicates the program should be run once
	// with no parameters.
	*start = 0
	*end = -1

	*runs = *runsPtr

	var err error

	// Get range if specified with flag
	if *rangePtr != "NONE" {
		split := strings.Split(*rangePtr, "-")
		if len(split) != 2 {
			fmt.Println("Please give the parameter range as two dash-delimited numbers. For example, 1-100")
			os.Exit(1)
		}

		*start, err = strconv.Atoi(split[0])
		check(err)

		*end, err = strconv.Atoi(split[1])
		check(err)
	}
}

// Check for an error.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Find the absolute path of a file
func findAbsolute(fileName string) string {
	var out string
	var err error

	// Check if file is in current directory
	pwd, _ := os.Getwd()

	if _, err := os.Stat(pwd + "/" + fileName); err == nil {
		return pwd + "/" + fileName
	}

	// Check if file is a binary
	out, err = exec.LookPath(fileName)

	// If file is not a binary, try to find abs path
	if err != nil || !filepath.IsAbs(out) {
		out, err = filepath.Abs(fileName)
		check(err)
	}

	// If no error return absolute path
	return out
}

// Get the filename and extension type of a file
func getFileName(arg string) (string, string) {
	abs := findAbsolute(arg)

	// Get file name
	fullPath := strings.Split(abs, "/")
	fileName := fullPath[len(fullPath)-1]

	// Find file type
	var extension string

	if strings.Contains(fileName, ".") { // file has extension
		extension = strings.Split(fileName, ".")[1]

	} else if !strings.Contains(abs, os.Getenv("HOME")) { // file is not in home dir
		extension = "system program"
	} else { // file is in home dir
		extension = "none"
	}
	return fileName, extension
}
