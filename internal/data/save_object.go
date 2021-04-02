// data package for saving objects
package data

import (
	"os"
    "encoding/json"
    "log"
    "time"
)

type Client struct {
    Id int64
    Time time.Time
}

/**
 * Saves a client's information into json
 */
func ClientToJson(client Client) []byte{

    // Save c as json byte array
    b, err := json.Marshal(client)

    // Exit on error, otherwise return b
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return b
}

/**
 * Saves a client's information into json
 */
 func ClientDataToJson(id int64, time time.Time) []byte{

    // Create Client Object
    c := Client{id, time}
    
    return ClientToJson(c)
}

/**
 * Coverts a []byte of json into a client struct
 */
func JsonToClient(b []byte) Client {
    var c Client

    // Unmarshall b into Client c
    err := json.Unmarshal(b, &c)

    // Exit on error, otherwise return c
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return c
}

type Job struct {
    Id int64
    Time time.Time
    Machines int64
    ParameterStart int64
    ParameterEnd int64
    Code []byte
}

/**
 * Saves a Job information into json
 */
func JobToJson(job Job) []byte{

    // Save c as json byte array
    b, err := json.Marshal(job)

    // Exit on error, otherwise return b
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return b
}

/**
 * Saves a Job information into json
 */
 func JobDataToJson(id int64, time time.Time, machines int64,
    parameterStart int64, parameterEnd int64, code []byte) []byte{

    // Create Job Object
    j := Job{id, time, machines, parameterStart, parameterEnd, code}
    
    return JobToJson(j)
}

/**
 * Coverts a []byte of json into a Job struct
 */
func JsonToJob(b []byte) Job {
    var j Job

    // Unmarshall b into Job j
    err := json.Unmarshal(b, &j)

    // Exit on error, otherwise return j
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return j
}

type Worker struct {
    Id int64
    Busy bool
	Hostname string
}

/**
 * Saves a Worker information into json
 */
 func WorkerToJson(worker Worker) []byte{

    // Save c as json byte array
    b, err := json.Marshal(worker)

    // Exit on error, otherwise return b
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return b
}

/**
 * Saves a Worker information into json
 */
 func WorkerDataToJson(id int64, busy bool, hostname string) []byte{

    // Create Worker Object
    w := Worker{id, busy, hostname}

    return WorkerToJson(w)
}

/**
 * Coverts a []byte of json into a Worker struct
 */
func JsonToWorker(b []byte) Worker {
    var w Worker

    // Unmarshall b into Worker w
    err := json.Unmarshal(b, &w)

    // Exit on error, otherwise return j
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return w
}
