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
	"strings"

	"github.com/showalter/bdws/internal/data"
)

type codeFunction func([]byte, string) []byte

// Map various extension names to their code
var extensionMap = map[string]codeFunction{
	"sh":         script,
	"py":         pythonScript,
	"java":       javaFile,
	"class":      javaClass,
	"jar":        jarFile,
	"none":       script,
	"executable": execute,
}

// run the code given an extension
func runCode(e string, code []byte, fn string) []byte {
	f, found := extensionMap[e]
	if found {
		return f(code, fn)
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
	output, err := exec.Command(command, args...).Output()
	check(err)
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

	// Print out the json.
	fmt.Println(jobJson)

	// Convert string json to job struct
	job := data.JsonToJob([]byte(jobJson))

	// Run the code and get []byte output
	output := runCode(job.Extension, job.Code, job.FileName)

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
func script(code []byte, fileName string) []byte {
	// Create a temporary file
	createFile(fileName, code)

	// Make temp file executable.
	check(os.Chmod(fileName, 0700))

	// Execute temp file.
	output := run(("./" + fileName), "")

	// Remove temp file.
	os.Remove(fileName)

	return output
}

// Run a .class file
func javaClass(code []byte, fileName string) []byte {

	// Create temporary file
	createFile(fileName, code)

	// Execute temp file.
	output := run("java", strings.Split(fileName, ".")[0])

	// Remove temp file
	os.Remove(fileName)

	return output
}

// Run a .java file
func javaFile(code []byte, fileName string) []byte {

	// Create temporary java file
	className := strings.Split(fileName, ".")[0] + ".class"
	createFile(fileName, code)

	// compile java file
	run("javac", fileName)

	// get []byte code from class file
	classCode, err := ioutil.ReadFile(className)
	if err != nil {
		panic(err)
	}

	// Remove the temp files
	os.Remove(fileName)
	os.Remove(className)

	// Return output
	return (javaClass(classCode, className))
}

// Run a jar file
func jarFile(code []byte, fileName string) []byte {

	// Create temporary file
	createFile(fileName, code)

	// Execute temp file.
	output := run("java", "-jar "+fileName)

	// Remove temp file
	os.Remove(fileName)

	return output
}

// Run a python script
func pythonScript(code []byte, fileName string) []byte {

	// Create temporary file
	createFile(fileName, code)

	// Execute temp script.
	output := run("python3", fileName)

	// Remove temp script
	os.Remove(fileName)

	return output
}

func execute(code []byte, fileName string) []byte {
	output := run(fileName)

	return output
}
