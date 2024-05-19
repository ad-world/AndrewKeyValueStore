package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"akv/akv"
	"net/rpc"
)

func printCommandUsage() {
	fmt.Println("Invalid command: ")
	fmt.Println("Usage: <command> <key> <value>")
	fmt.Println("Commands: get, put, delete")
}

func main() {
	var address string 
	flag.StringVar(&address, "addr", "", "The address of the server")
	flag.Parse()

	if address == "" {
		fmt.Println("Error: -addr flag is required")
		flag.Usage()
		os.Exit(1)
	}

	in := bufio.NewReader(os.Stdin)

	andrew := new(akv.AndrewKeyValueClient)
	var err error
	andrew.Client, err = rpc.Dial("tcp", address)

	if err != nil {
		fmt.Println("Error in connecting to server:", err)
	} 

	for {
		fmt.Print("akv_cli> ")
		var command string
		command, _ = in.ReadString('\n')
		stripped := strings.TrimSpace(command)
		if stripped == "exit" {
			break
		}

		parts := strings.Split(stripped, " ")
		if len(parts) < 1 {
			printCommandUsage()
			continue
		}

		switch parts[0] {
		case "get":
			if len(parts) != 2 {
				printCommandUsage()
				continue
			}
			value, err := andrew.Get(parts[1])
			if err != nil {
				fmt.Println("Error in Get:", err)
				continue
			}
			fmt.Println("Value of", parts[1], "is", value)
		case "put":
			if len(parts) != 3 {
				printCommandUsage()
				continue
			}
			_, err := andrew.Put(parts[1], parts[2])
			if err != nil {
				fmt.Println("Error in Put:", err)
				continue
			}
			fmt.Println("Put successful")
		case "delete":
			if len(parts) != 2 {
				printCommandUsage()
				continue
			}
			_, err := andrew.Delete(parts[1])
			if err != nil {
				fmt.Println("Error in Delete:", err)
				continue
			}
			fmt.Println("Delete successful")
		default:
			printCommandUsage()
		}
	}
}