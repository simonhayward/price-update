package priceupdate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	RequestTimeout = 10 * time.Second
	UserAgent      = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"
)

func GetResponse(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return makeRequest(req, 200)
}

func PatchResponse(url, token string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	return makeRequest(req, 200)
}

func makeRequest(req *http.Request, expected int) ([]byte, error) {
	client := &http.Client{Timeout: RequestTimeout}
	req.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expected {
		return nil, fmt.Errorf("want: %d got: %d", expected, resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
