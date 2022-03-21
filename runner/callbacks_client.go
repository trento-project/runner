package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=CallbacksClient

type CallbacksClient interface {
	Callback(executionID int64, event string, payload interface{}) error
}

type callbacksClient struct {
	callbacksUrl string
	httpClient   *http.Client
}

func NewCallbacksClient(callbacksUrl string) *callbacksClient {
	httpClient := &http.Client{}

	return &callbacksClient{
		callbacksUrl: callbacksUrl,
		httpClient:   httpClient,
	}
}

func (c *callbacksClient) Callback(executionID int64, event string, payload interface{}) error {
	log.Debugf("Executing callback for execution %d with event %s", executionID, event)

	requestBody, err := json.Marshal(map[string]interface{}{
		"execution_id": executionID,
		"event":        event,
		"payload":      payload,
	})
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(c.callbacksUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"something wrong happened while sending the callback data. Status: %d, Execution: %d, Event: %s",
			resp.StatusCode, executionID, event)
	}

	return nil
}
