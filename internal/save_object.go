// main package for saving objects
package main

import (
	"os"
	"fmt"
    "encoding/json"
    "log"
)

type Client struct {
    Id int64
    Time int64
}

/**
 * Saves a client's information into json
 */
func client_to_json(id int64, time int64) []byte{

    // Create Client Object
    c := Client{id, time}
    // Save c as json byte array
    b, err := json.Marshal(c)

    // Exit on error, otherwise return b
    if err != nil {
        log.Println(err)
        os.Exit(-1)
    }
    return b;
}

/**
 * Coverts a []byte of json into a client struct
 */
func json_to_client(b []byte) Client {
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
