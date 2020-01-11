package pkg

import "net/http"

type ReqPass struct {
	URL     string
	Body    []byte
	Method  string
	Headers http.Header
}
