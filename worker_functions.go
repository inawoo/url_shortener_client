package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type (
	WorkerFunctionsAdapter struct {
		BaseURL string
	}
)

func NewWorkerFunctionsAdapter(baseURL string) *WorkerFunctionsAdapter {

	return &WorkerFunctionsAdapter{
		BaseURL: baseURL,
	}
}

func (w *WorkerFunctionsAdapter) CheckHealth() (string, error) {
	req, err := http.NewRequest(http.MethodGet, w.BaseURL+"/api/health", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return "", errors.New(string(bodyBytes))
	}

	return "success", nil
}

func (w *WorkerFunctionsAdapter) ShortenURL(input ShortenURLRequest) (*URLCollection, error) {

	byteData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, w.BaseURL+"/api/save", bytes.NewBuffer(byteData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//convert body to string

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(bodyBytes))
	}

	defer resp.Body.Close()

	var result URLCollection
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
