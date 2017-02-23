package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const imgurURL = "https://api.imgur.com/3/image"

var (
	imgurID string
)

func setupImgur(id string) {
	imgurID = id
}

type imgurPostResponse struct {
	Data    Data `json:"data"`
	Success bool `json:"success"`
}

type Data struct {
	Link string `json:"link"`
}

func postImgur(b io.Reader) (string, error) {

	if imgurID == "" {
		return "", fmt.Errorf("imgur is not configured")
	}

	r, err := http.NewRequest("POST", imgurURL, b)
	if err != nil {
		return "", err
	}

	r.Header.Add("Authorization", "Client-ID "+imgurID)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var data imgurPostResponse
	err = dec.Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Data.Link, nil
}
