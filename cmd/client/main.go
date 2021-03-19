package main

import (
	"net/http"
	"os"
	"fmt"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		fmt.Println("Please pass a file name to send to the server.")
		os.Exit(1)
	}

	dat, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}

	_, err = http.Post("http://127.0.0.1:39485/newjob", "text/plain", dat)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done?")

}