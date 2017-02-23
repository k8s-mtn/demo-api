package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

var resizeHost string

func NewProxy(dest string) (*httputil.ReverseProxy, error) {

	u := url.URL{
		Scheme: "http",
		Host:   dest,
	}

	resizeHost = dest

	p := httputil.NewSingleHostReverseProxy(&u)

	p.Director = rewriteRequest
	p.ModifyResponse = rewriteResponse

	return p, nil
}

func rewriteRequest(r *http.Request) {

	r.URL.Scheme = "http"
	r.URL.Host = resizeHost
	r.URL.Path = "/resize"

}

func rewriteResponse(r *http.Response) error {
	r.Header.Add("Image-Magician", "We did it!")
	return nil
}
