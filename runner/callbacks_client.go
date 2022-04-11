package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

//go:generate mockery --name=CallbacksClient

type CallbacksClient interface {
	Callback(executionID uuid.UUID, event string, payload interface{}) error
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

func (c *callbacksClient) Callback(executionID uuid.UUID, event string, payload interface{}) error {
	log.Debugf("Executing callback for execution %s with event %s", executionID, event)

	requestBody, err := json.Marshal(map[string]interface{}{
		"execution_id": executionID.String(),
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

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"something wrong happened while sending the callback data. Status: %d, Execution: %s, Event: %s",
			resp.StatusCode, executionID, event)
	}

	return nil
}
