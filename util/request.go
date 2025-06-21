package util

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const RequestUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
const RequestTimeout = 60 * time.Second

func SendGetRequest(endpointUrl string) (string, error) {

	client := &http.Client{
		Timeout: RequestTimeout,
	}

	req, err := http.NewRequest("GET", endpointUrl, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", RequestUserAgent)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil

}

func SendPostRequest(endpointUrl string, data map[string]string) (string, error) {

	client := &http.Client{
		Timeout: RequestTimeout,
	}

	form := url.Values{}

	for key, value := range data {
		form.Add(key, value)
	}

	req, err := http.NewRequest("POST", endpointUrl, strings.NewReader(form.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", RequestUserAgent)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil

}
