package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

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

	r.ParseForm()
	maxX := r.Form.Get("x")
	maxY := r.Form.Get("y")

	// resize the image
	img, err := resize(p.resizeAddr, r.Body, maxX, maxY)
	if err != nil {
		log.Printf("unable to resize: %s\n", err)
		http.Error(w, "unable to resize: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer img.Close()

	imgur := r.Form.Get("imgur")
	twilio := r.Form.Get("send")

	var url string
	if imgur != "" || twilio != "" {
		url, err = postImgur(img)
		if err != nil {
			http.Error(w, "unable to post to imgur: "+err.Error(), http.StatusInternalServerError)
		}

		// replace the image reader with the url to imgur
		img = ioutil.NopCloser(bytes.NewBuffer([]byte(url)))
	}

	if twilio != "" {
		err := sendTwilio(twilio, url)
		if err != nil {
			http.Error(w, "unable to send to twilio: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	io.Copy(w, img)
}

func resize(addr string, r io.Reader, x string, y string) (io.ReadCloser, error) {

	q := url.Values{}
	q.Add("x", x)
	q.Add("y", y)

	u := url.URL{
		Scheme:   "http",
		Host:     addr,
		Path:     "/resize",
		RawQuery: q.Encode(),
	}

	resp, err := http.Post(u.String(), "", r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// if there's an error reading, ignore it
		out, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("bad status code: %d [%s]", resp.StatusCode, string(out))
	}

	return resp.Body, nil

}
