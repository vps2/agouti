package bus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Connect(url string, capabilities map[string]interface{}) (*Client, error) {
	requestBody, err := capabilitiesToJSON(capabilities)
	if err != nil {
		return nil, err
	}

	sessionID, err := openSession(url, requestBody)
	if err != nil {
		return nil, err
	}

	sessionURL := fmt.Sprintf("%s/session/%s", url, sessionID)
	return &Client{sessionURL}, nil
}

func capabilitiesToJSON(capabilities map[string]interface{}) (io.Reader, error) {
	if capabilities == nil {
		capabilities = map[string]interface{}{}
	}
	desiredCapabilities := struct {
		DesiredCapabilities map[string]interface{} `json:"desiredCapabilities"`
	}{capabilities}

	capabiltiesJSON, err := json.Marshal(desiredCapabilities)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(capabiltiesJSON), err
}

func openSession(url string, body io.Reader) (sessionID string, err error) {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/session", url), body)
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	var sessionResponse struct{ SessionID string }

	if err := json.NewDecoder(response.Body).Decode(&sessionResponse); err != nil {
		return "", err
	}

	if sessionResponse.SessionID == "" {
		return "", errors.New("failed to retrieve a session ID")
	}

	return sessionResponse.SessionID, nil
}