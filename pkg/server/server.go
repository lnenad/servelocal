package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/lnenad/servelocal/pkg"
	"github.com/lnenad/servelocal/pkg/client"
)

type Address struct {
	resp chan pkg.RespPass
	req  chan pkg.ReqPass
	addr net.Addr
}

func Server(tcpPort, serverPort *string) {
	listener, err := net.Listen("tcp", "127.0.0.1:"+*tcpPort)
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	fmt.Println("listening for tcp connections on ", *tcpPort, " port")

	listeners := make(map[string]Address)

	http.HandleFunc("/", HandleRequest(listeners))
	go http.ListenAndServe(":"+*serverPort, nil)

	fmt.Println("web server listening on ", *serverPort, " port")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		dec := gob.NewDecoder(conn)
		var cl client.ClientData
		err = dec.Decode(&cl)
		if err != nil {
			fmt.Println("decode error:", err)
			conn.Close()
			return
		}

		fmt.Printf("new client %s \n", cl)

		listeners[cl.Token] = Address{
			resp: make(chan pkg.RespPass),
			req:  make(chan pkg.ReqPass),
			addr: conn.RemoteAddr(),
		}
		go handleConnection(cl, conn, &listeners)
	}
}

func handleConnection(cl client.ClientData, conn net.Conn, listeners *map[string]Address) {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	select {
	case request := <-(*listeners)[cl.Token].req:
		fmt.Println("sending message", request)

		// Create an encoder and send a value.
		err := enc.Encode(request)
		if err != nil {
			log.Println("client issues..", err)
			conn.Close()
			return
		}

		var rp pkg.RespPass
		err = dec.Decode(&rp)
		if err != nil {
			fmt.Println("decode error:", err)
			conn.Close()
			return
		}

		(*listeners)[cl.Token].resp <- rp

		handleConnection(cl, conn, listeners)
	}
}
