package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/lnenad/servelocal/pkg/client"
)

func main() {
	parser := argparse.NewParser("servelocal-client", "Forwards remote traffic to localhost")
	portTarget := parser.String("", "port-target", &argparse.Options{Required: true, Help: "Local port target"})
	portTCP := parser.String("", "port-server", &argparse.Options{Help: "Server port", Default: "8080"})
	server := parser.String("", "server", &argparse.Options{Help: "TCP Server", Default: "127.0.0.1"})
	host := parser.String("", "host", &argparse.Options{Help: "Local host target", Default: "localhost"})
	schema := parser.String("", "schema", &argparse.Options{Help: "Local host schema", Default: "http"})
	clientID := parser.String("", "client", &argparse.Options{Help: "Client ID", Default: "mrs"})
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
	client.Client(clientID, server, schema, host, portTarget, portTCP)
}
