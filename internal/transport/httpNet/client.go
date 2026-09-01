package httpNet

import "net/http"

type ClientHTTP struct {
	hc *http.Client
}

func NewClient(hc *http.Client) *ClientHTTP {
	return &ClientHTTP{hc: hc}
}
