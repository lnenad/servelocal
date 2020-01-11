package pkg

import "net/http"

type RespPass struct {
	Body       []byte
	Headers    http.Header
	StatusCode int
}
