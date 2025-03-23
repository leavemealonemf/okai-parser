package magichttp

import (
	"bytes"
	"fmt"
	"net/http"
)

func POST(url string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer resp.Body.Close()
	return resp, nil
}
