package feature

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// RestContext is a struct used to share common steps across the feature test
type RestContext struct {
	BaseURL string
	Resp    *http.Response
}

// RegisterRestContext ...
func RegisterRestContext(s *godog.Suite, baseURL string) *RestContext {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	c := RestContext{
		BaseURL: baseURL,
	}

	s.Step(`^I request REST endpoint with method "([^"]*)" and path "([^"]*)"$`, c.iRequestWithMethodAndPath)
	s.Step(`^I request REST endpoint with method "([^"]*)" and path "([^"]*)" and body$`,
		c.iRequestWithMethodAndPathAndBody)
	s.Step(`^I should have an internal error response$`, c.iShouldHaveAnInternalErrorResponse)
	s.Step(`^I should have a bad request response$`, c.iShouldHaveABadRequestResponse)
	s.Step(`^I should have a precondition failed response$`, c.iShouldHaveAPreconditionFailedRequestResponse)
	s.Step(`^I should have a not found response$`, c.iShouldHaveANotFoundResponse)
	s.Step(`^I should have an OK response$`, c.iShouldHaveAnOKResponse)
	s.Step(`^I should have a no content response$`, c.iShouldHaveANoContentResponse)
	s.Step(`^I should have a created response$`, c.iShouldHaveACreatedResponse)
	s.Step(`^I should have a response with following code "([^"]*)"$`, c.iShouldHaveResponseWithTheFollowingCode)
	s.Step(`^I should have a response with following JSON body$`, c.iShouldHaveResponseWithTheFollowingJSONBody)

	return &c
}

func (c *RestContext) iRequestWithMethodAndPath(method, path string) error {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	c.Resp = resp

	return err
}

func (c *RestContext) iRequestWithMethodAndPathAndBody(method, path string, body *gherkin.DocString) error {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body.Content))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	c.Resp = resp

	return err
}

func (c *RestContext) iShouldHaveResponseWithCode(statusCode int) error {
	if c.Resp.StatusCode != statusCode {
		responseDump, err := httputil.DumpResponse(c.Resp, true)
		if err != nil {
			return fmt.Errorf("invalid response code, expected: %d (%s) got: %d",
				statusCode, http.StatusText(statusCode), c.Resp.StatusCode)
		}

		return fmt.Errorf("invalid response code, expected: %d (%s) got: %d on the response:\n\n %s",
			statusCode, http.StatusText(statusCode), c.Resp.StatusCode, string(responseDump))
	}
	return nil
}

func (c *RestContext) iShouldHaveAnInternalErrorResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusInternalServerError)
}

func (c *RestContext) iShouldHaveABadRequestResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusBadRequest)
}

func (c *RestContext) iShouldHaveAPreconditionFailedRequestResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusPreconditionFailed)
}

func (c *RestContext) iShouldHaveANotFoundResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusNotFound)
}

func (c *RestContext) iShouldHaveAnOKResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusOK)
}

func (c *RestContext) iShouldHaveANoContentResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusNoContent)
}

func (c *RestContext) iShouldHaveACreatedResponse() error {
	return c.iShouldHaveResponseWithCode(http.StatusCreated)
}

func (c *RestContext) iShouldHaveResponseWithTheFollowingCode(code int) error {
	return c.iShouldHaveResponseWithCode(code)
}

type testingT struct {
	Err error
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.Err = fmt.Errorf(format, args...)
}

func (c *RestContext) iShouldHaveResponseWithTheFollowingJSONBody(expected *gherkin.DocString) error {
	receivedBytes, err := ioutil.ReadAll(c.Resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	var receivedData, expectedData interface{}

	err = json.Unmarshal([]byte(expected.Content), &expectedData)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal expected data")
	}
	err = json.Unmarshal(receivedBytes, &receivedData)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal received data (%s)", string(receivedBytes))
	}

	t := &testingT{}
	assert.Equal(t, expectedData, receivedData, "received data: %s", string(receivedBytes))
	if t.Err != nil {
		return t.Err
	}

	return nil
}
