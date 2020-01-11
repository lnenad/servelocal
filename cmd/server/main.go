package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/lnenad/servelocal/pkg/server"
)

func main() {
	parser := argparse.NewParser("servelocal-server", "Serves remote traffic to servelocal clients")
	portTCP := parser.String("", "port-tcp", &argparse.Options{Required: true, Help: "Remote tcp server port"})
	portServer := parser.String("", "port-web", &argparse.Options{Required: true, Help: "Web server port"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	//go createConnection("aljaska")
	//go createConnection("moldavija")
	server.Server(portTCP, portServer)
}
