package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/lnenad/servelocal/pkg"
)

const CookieClientVal = "clientval"

// HandleRequest requires clients that support cookies to properly serve web pages
func HandleRequest(listeners map[string]Address) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("URL", r.URL.String())
		client := ""
		cookie, err := r.Cookie(CookieClientVal)
		if err != nil {
			urlStr := r.URL.String()
			data := strings.Split(urlStr[1:], ":")
			if len(data) <= 1 {
				fmt.Println("Invalid data in URL")
				return
			}
			client = data[0]
		} else if cookie.Value != "" {
			client = cookie.Value
		}
		if lst, ok := listeners[client]; ok {
			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{Name: CookieClientVal, Value: client, Expires: expiration, Path: ""}
			http.SetCookie(w, &cookie)
			lst.req <- generateRequestPayload(r)
			select {
			case resp := <-listeners[client].resp:
				for k, vals := range resp.Headers {
					if strings.ToLower(k) == "content-encoding" {
						continue
					}
					for _, v := range vals {
						w.Header().Add(k, v)
					}
				}
				if resp.StatusCode > 199 {
					w.WriteHeader(resp.StatusCode)
				}
				fmt.Fprint(w, string(resp.Body))
			}
		} else {
			fmt.Println("Invalid listener: ", client, listeners)
			fmt.Fprint(w, `<body style="background: black; color: white">Invalid listener</body>`)
		}
	}
}

func generateRequestPayload(r *http.Request) pkg.ReqPass {
	urlStr := r.URL.String()
	data := strings.Split(urlStr[1:], ":")
	url := data[0]
	if len(data) == 2 {
		url = data[1]
	}
	body, _ := ioutil.ReadAll(r.Body)

	rq := pkg.ReqPass{
		URL:     url,
		Body:    body,
		Method:  r.Method,
		Headers: r.Header,
	}

	return rq
}
