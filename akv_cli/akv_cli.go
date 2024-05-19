package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

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

	for {
		fmt.Print("akv_cli> ")
		var command string
		command, _ = in.ReadString('\n')
		stripped := strings.TrimSpace(command)
		if stripped == "exit" {
			break
		} else {
			fmt.Println(command)
		}
	}
}