package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func serve(addr string, resizeAddr string) (func(context.Context) error, error) {

	p, err := NewProxy(resizeAddr)
	if err != nil {
		return nil, err
	}

	http.HandleFunc("/magician", p.magicianHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/leak", leakHandler)

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

type proxy struct {
	resizeAddr string
}

func NewProxy(dest string) (*proxy, error) {

	p := proxy{
		resizeAddr: dest,
	}

	return &p, nil
}

func (p *proxy) magicianHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received magician request")
	defer log.Println("magician request done")

	r.ParseForm()
	maxX := r.Form.Get("x")
	maxY := r.Form.Get("y")

	imgur := r.Form.Get("imgur") != ""
	twilio := r.Form.Get("send")

	img, url, err := processImage(p.resizeAddr, r.Body, maxX, maxY, twilio, imgur)
	if err != nil {
		http.Error(w, "unable to process image: "+err.Error(), http.StatusBadRequest)
		return
	}

	if img != nil {
		defer img.Close()
	}

	if url != "" {
		w.Write([]byte(url))
		return
	}

	io.Copy(w, img)

}

func leakHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("leaking memory")

	for {
		a := make([]byte, 50000000)
		go func(a []byte) {
			for {
				for i, v := range a {
					a[i] = v + 1
				}
				time.Sleep(time.Second * 1)
			}

		}(a)

		time.Sleep(time.Second * 1)
	}
}
