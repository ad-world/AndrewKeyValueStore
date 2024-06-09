package main

import (
	"akv/akv"
	"fmt"
	"net/rpc"
)

func main() {
	client := new(akv.AndrewKeyValueClient)
	// Connect to the server which is listening on port 1234
	client.Client, _ = rpc.Dial("tcp", "localhost:1234")

	_, err := client.Put("key1", "value1")
	if err != nil {
		fmt.Println("Error in Put:", err)
		return
	}
	value, _ := client.Get("key1")

	fmt.Println("Value of key1 is", value)

	_, err = client.Get("key2")
	if err != nil {
		fmt.Println("Error in Get:", err)
		return
	}
}