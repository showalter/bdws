// This file contains the main routine for workers.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/showalter/bdws/internal/data"
)

type codeFunction func([]byte, string, *int, []string) []byte

// Map various extension names to their code
var extensionMap = map[string]codeFunction{
	"sh":             script,
	"py":             pythonScript,
	"java":           javaFile,
	"class":          javaClass,
	"jar":            jarFile,
	"none":           script,
	"system program": system_program,
}

var workerDirectory string

// run the code given an extension
func runCode(e string, code []byte, fn string, num *int, args []string) []byte {
	f, found := extensionMap[e]
	if found {
		return f(code, fn, num, args)
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
func run(command string, args ...string) []byte {

	output, _ := exec.Command(command, args...).CombinedOutput()
	return output
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

	var num *int = nil

	if job.ParameterEnd >= job.ParameterStart {
		num = &job.ParameterStart
	}

	var args []string
	if job.Args[0] != "NONE" {
		args = job.Args
	}

	// Run the code and get []byte output
	output := runCode(job.Extension, job.Code, job.FileName, num, args)

	// Send a response back.
	w.Write(output)
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

	resp, err := http.Post(args[1]+"/register", "text/plain", strings.NewReader(args[2]))
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// This gives what the supervisor thinks the worker is, which is useful for debugging.
	_ = data.JsonToWorker(buf.Bytes())

	// Make a directory for this worker, to avoid IO errors from workers writing and reading to
	// the same file.
	workerDirectory = args[2]
	if _, err = os.Stat(workerDirectory); os.IsNotExist(err) {
		err = os.Mkdir(args[2], 0777)
		check(err)
	}

	// If there is a request for /newjob,
	// the new_job routine will handle it.
	http.HandleFunc("/newjob", new_job)

	// Listen on a port.
	log.Fatal(http.ListenAndServe(":"+args[2], nil))
}

/* Code Strategies */

// create a tmp file with given name and write to it
func createFile(name string, code []byte) {
	// Create a temporary file
	file, err := os.Create(name)
	check(err)

	// Write to file
	_, err = file.Write(code)
	check(err)
	file.Sync()
	file.Close()
}

// Run a bash script / script
func script(code []byte, fileName string, num *int, args []string) []byte {

	var output []byte

	fullName := workerDirectory + "/" + fileName

	// Create a temporary file
	createFile(fullName, code)

	// Make temp file executable.
	check(os.Chmod(fullName, 0700))

	// Execute temp file.
	if num != nil {
		args = append([]string{strconv.Itoa(*num)}, args...)
	}
	output = run(fullName, args...)

	return output
}

// Run a .class file
func javaClass(code []byte, fileName string, num *int, args []string) []byte {

	var output []byte

	fullName := workerDirectory + "/" + fileName

	// Create temporary file
	createFile(fullName, code)

	// Execute temp file.
	if num != nil {
		args = append([]string{strconv.Itoa(*num)}, args...)
	}

	args = append([]string{"-cp", workerDirectory, strings.Split(fileName, ".")[0]}, args...)
	output = run("java", args...)

	return output
}

// Run a .java file
func javaFile(code []byte, fileName string, num *int, args []string) []byte {

	fullName := workerDirectory + "/" + fileName

	// Create temporary java file
	className := strings.Split(fileName, ".")[0] + ".class"

	_, err := os.Stat(fullName)
	if os.IsNotExist(err) {
		createFile(fullName, code)

		// compile java file
		run("javac", fullName)

	} else {
		existingCode, err := ioutil.ReadFile(fullName)
		check(err)

		if !bytes.Equal(existingCode, code) {
			createFile(fullName, code)

			// compile java file
			run("javac", fullName)
		}
	}

	// get []byte code from class file
	classCode, err := ioutil.ReadFile(workerDirectory + "/" + className)
	if err != nil {
		panic(err)
	}

	// Return output
	return (javaClass(classCode, className, num, args))
}

// Run a jar file
func jarFile(code []byte, fileName string, num *int, args []string) []byte {

	var output []byte

	fullName := workerDirectory + "/" + fileName

	// Create temporary file
	createFile(fullName, code)

	// Execute temp file.

	if num != nil {
		args = append([]string{strconv.Itoa(*num)}, args...)
	}

	args = append([]string{"-jar", fullName}, args...)
	output = run("java", args...)

	return output
}

// Run a python script
func pythonScript(code []byte, fileName string, num *int, args []string) []byte {

	var output []byte

	fullName := workerDirectory + "/" + fileName

	// Create temporary file
	createFile(fullName, code)

	// Execute temp script.
	if num != nil {
		args = append([]string{strconv.Itoa(*num)}, args...)
	}
	args = append([]string{fullName}, args...)
	output = run("python3", args...)

	return output
}

// Run a system program
func system_program(code []byte, fileName string, num *int, args []string) []byte {

	var output []byte

	if num != nil {
		args = append([]string{strconv.Itoa(*num)}, args...)
	}

	output = run(fileName, args...)

	return output
}
