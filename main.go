package main

import (
	"flag"
	"fmt"
)

func main() {

	port := flag.String("p", "8000", "a string")
	server := flag.Bool("s", false, "a bool")
	client := flag.Bool("c", false, "a bool")
	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()

	if *server {
		InitServer(*port)
	} else if *client {
		InitClient(*port)
	} else {
		fmt.Printf("You must enter server or client option")
	}
}
