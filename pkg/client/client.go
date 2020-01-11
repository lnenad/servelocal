package client

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/lnenad/servelocal/pkg"
)

type ClientData struct {
	Token string
}

func Client(client, server, schema, host, portTarget, portServer *string) {
	//go createConnection("aljaska")
	//go createConnection("moldavija")
	createConnection(*client, *server, *schema, *host, *portTarget, *portServer)
}

func createConnection(client, server, schema, host, portTarget, portServer string) {
	conn, err := net.Dial("tcp", server+":"+portServer)
	if err != nil {
		fmt.Print("Unable to connect to server")
		conn.Close()
		return
	}
	defer conn.Close()
	encoderC := gob.NewEncoder(conn)
	err = encoderC.Encode(ClientData{
		Token: client,
	})
	if err != nil {
		log.Println("client issues..", err)
		conn.Close()
	}

	fmt.Fprintf(conn, "reg:"+client+"\n")
	for {
		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)
		fmt.Println("Waiting for messages")
		var req pkg.ReqPass
		err := decoder.Decode(&req)
		if err != nil {
			log.Fatal("decode:", err)
		}
		fmt.Printf("Request from server: %s\n", req.URL)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		url := schema + "://" + host + ":" + portTarget + "/" + req.URL
		//url := "http://webhook.site/948f2d09-071e-466e-a384-530ef477d509/" + req.URL
		fmt.Println("URL: ", url, "METHOD:", req.Method)
		body := bytes.NewBuffer(req.Body)
		aReq, err := http.NewRequest(req.Method, url, body)
		aReq.Header = req.Headers
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(aReq)
		if err != nil {
			panic(err)
		}
		respPass := generateResponsePayload(resp)

		err = encoder.Encode(respPass)
		if err != nil {
			log.Fatal("encode:", err)
		}
	}
}

func generateResponsePayload(r *http.Response) pkg.RespPass {
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("content-encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			fmt.Println("error gzip", err)
			reader = r.Body
		}
		defer reader.Close()
	default:
		reader = r.Body
	}
	body, _ := ioutil.ReadAll(reader)

	clh := r.Header.Get("content-length")
	if clh != "" {
		cl, err := strconv.Atoi(clh)
		if err != nil {
			panic(err)
		}
		if cl != len(body) {
			r.Header.Set("content-length", strconv.Itoa(len(body)))
		}
	}

	fmt.Println("STATUS CODE", r.StatusCode)
	fmt.Println("STATUS CODE", r.Request.URL.String())
	if strings.Contains(r.Request.URL.String(), "admin") {
		fmt.Println("BODY", r.Header)
		fmt.Println("BODY", string(body))
	}

	rq := pkg.RespPass{
		Body:       body,
		Headers:    r.Header,
		StatusCode: r.StatusCode,
	}

	return rq
}
