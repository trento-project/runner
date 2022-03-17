package runner

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/trento-project/runner/runner/mocks"
)

type ApiTestCase struct {
	suite.Suite
	config *Config
}

func TestApiTestCase(t *testing.T) {
	suite.Run(t, new(ApiTestCase))
}

func (suite *ApiTestCase) SetupTest() {
	suite.config = &Config{}
}

func (suite *ApiTestCase) Test_ApiHealthTest() {
	app, err := NewApp(suite.config)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/health", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedJson, _ := json.Marshal(map[string]string{"status": "ok"})
	suite.Equal(200, resp.Code)
	suite.JSONEq(string(expectedJson), resp.Body.String())
}

func (suite *ApiTestCase) Test_ApiReadyTest() {
	mockRunnerService := new(mocks.RunnerService)
	mockRunnerService.On("IsCatalogReady").Return(true)

	deps := Dependencies{
		mockRunnerService,
	}

	app, err := NewAppWithDeps(suite.config, deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/ready", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedJson, _ := json.Marshal(map[string]bool{"ready": true})
	suite.Equal(200, resp.Code)
	suite.JSONEq(string(expectedJson), resp.Body.String())
}
