package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Check(err error, format string, args ...interface{}) {
	if err != nil {
		if format != "" {
			log.Fatalf(format, args...)
		} else {
			log.Println(err.Error())
		}
	}
}

func SendHTTPRequest(method, url string, headers map[string]string, payload interface{}) (*http.Response, error) {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}
