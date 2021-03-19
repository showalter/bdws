package main

import (
	"bytes"
	"net/http"
	"fmt"
)

func new_job(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling connection...")

	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	file := buf.String()

	fmt.Printf(file)

}

func main() {
	http.HandleFunc("/newjob", new_job)

	http.ListenAndServe(":39485", nil)
}
