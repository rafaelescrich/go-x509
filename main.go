package main

import (
	"fmt"
	"os"

	"github.com/simonleung8/flags"
)

func main() {
	fc := flags.New()
	fc.NewStringFlagWithDefault("port", "p", "flag for port number", "8000")
	fc.NewStringFlag("gen-keys", "g", "flag for generate new keys")
	fc.NewBoolFlag("server", "s", "initializes software in server mode")
	fc.NewBoolFlag("client", "c", "initializes software in client mode")
	fc.Parse(os.Args...)

	if fc.IsSet("s") {
		InitServer(fc.String("p"))
	} else if fc.IsSet("c") {
		InitClient(fc.String("p"))
	} else if fc.IsSet("g") {
		genKeys(fc.String("g"))
	} else {
		fmt.Println("You must enter server or client option")
		os.Exit(1)
	}
}
