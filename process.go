package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func processImage(addr string, body io.Reader, maxX string, maxY string, to string, save bool) (io.ReadCloser, string, error) {
	// resize the image
	img, err := resize(addr, body, maxX, maxY)
	if err != nil {
		return nil, "", err
	}

	var url string
	if save || to != "" {
		url, err = postImgur(img)
		if err != nil {
			return nil, "", err
		}
	}

	if to != "" {
		err := sendTwilio(to, url)
		if err != nil {
			return nil, "", err
		}
	}

	return img, url, nil
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
