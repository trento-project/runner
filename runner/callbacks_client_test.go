package runner

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/runner/test/helpers"
)

type CallbacksTestSuite struct {
	suite.Suite
	configuredClient *callbacksClient
}

func TestCallbacksTestSuite(t *testing.T) {
	suite.Run(t, new(CallbacksTestSuite))
}

func (suite *CallbacksTestSuite) SetupSuite() {
	suite.configuredClient = NewCallbacksClient("http://192.168.1.1:8000/api/runner/callbacks")
}

func (suite *CallbacksTestSuite) Test_Callback() {
	client := suite.configuredClient

	payload := map[string]interface{}{
		"key": "value",
	}

	dummyID := uuid.New()

	client.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"execution_id": dummyID,
			"event":        "new_callback_event",
			"payload":      payload,
		})

		outgoingRequestBody, _ := ioutil.ReadAll(req.Body)

		suite.EqualValues(requestBody, outgoingRequestBody)

		suite.Equal(req.URL.String(), "http://192.168.1.1:8000/api/runner/callbacks")
		return &http.Response{
			StatusCode: 200,
		}
	})

	err := client.Callback(dummyID, "new_callback_event", payload)

	suite.NoError(err)
}
