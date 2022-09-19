package api

import (
        "net/http"
        "fmt"
        "io"
)

func handleGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err := res.Body.Close(); err != nil {
		return nil, err
	}

	// ReadAll err
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 ||
		res.StatusCode > 299 {
		err = fmt.Errorf("Status: %s", res.Status)
	}

	return body, err
}

