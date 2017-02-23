package main

import (
	"context"
	"log"
	"net/http"
)

func serve(addr string, resizeAddr string) (func(context.Context) error, error) {

	p, err := NewProxy(resizeAddr)
	if err != nil {
		return nil, err
	}

	http.HandleFunc("/magician", p.magicianHandler)
	http.HandleFunc("/ping", pingHandler)

	s := http.Server{
		Addr:    addr,
		Handler: nil,
	}

	// we're ignoring the error here, which is bad
	go s.ListenAndServe()

	return s.Shutdown, nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received ping request")
	defer log.Println("Done with ping request")

	w.Write([]byte("pong"))
}
