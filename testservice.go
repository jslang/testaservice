package testaservice

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestService is a helpful test utility that allows tests to call a real HTTP service during test runs.
type TestService struct {
	*httptest.Server
	request *http.Request

	t            *testing.T
	responseBody []byte
	responseCode int
}

// SetResponseBody will set the provided value as the response body for any requests received by the service.
//
// The method accepts bytes and strings which will be returned directly from the service. Passing structs or other types
// will result in the response being the JSON encoding of those values.
func (s *TestService) SetResponseBody(b interface{}) {
	if body, isBytes := b.([]byte); isBytes {
		s.responseBody = body
		return
	}

	if body, isString := b.(string); isString {
		s.responseBody = []byte(body)
		return
	}

	// TODO: Add support for customizing the marshal behaviour
	body, _ := json.Marshal(b)
	s.responseBody = body
}

// SetResponseCode will set the provided value as the status code the service will return for any requests received by the service.
func (s *TestService) SetResponseCode(code int) {
	s.responseCode = code
}

// AssertReceivedBasicAuth allows callers to assert that the last request received was authenticated with the provided  values.
func (s *TestService) AssertReceivedBasicAuth(user, pass string) {
	requestUser, requestPass, valid := s.request.BasicAuth()
	require.True(s.t, valid, "valid basic authentication expected")
	require.Equal(s.t, user, requestUser)
	require.Equal(s.t, pass, requestPass)
}

// AssertCalled allows callers to assert that the service has received a request.
func (s *TestService) AssertCalled() {
	require.NotNil(s.t, s.request, "test service should have been called but wasn't")
}

// AssertNotCalled allows callers to assert that the service has not received a request.
func (s *TestService) AssertNotCalled() {
	require.Nil(s.t, s.request, "test service should not have been called but was")
}

// AssertReceivedPath allows callers to assert that the last request received was for the path provided.
func (s *TestService) AssertReceivedPath(path string) {
	require.Equal(s.t, path, s.request.URL.Path)
}

// AssertReceivedHeader allows callers to assert that the last request received contained the header provided.
func (s *TestService) AssertReceivedHeader(name string, value string) {
	require.NotNil(s.t, s.request)
	require.Equal(s.t, value, s.request.Header.Get(name))
}

// AssertReceivedHeader allows callers to assert that the last request received contained the query param provided.
func (s *TestService) AssertReceivedParam(name string, value string) {
	require.NotNil(s.t, s.request)
	require.Equal(s.t, value, s.request.URL.Query().Get(name))
}

// AssertReceivedBody allows callers to assert that the last request received matches the body provided.
func (s *TestService) AssertReceivedBody(body []byte) {
	require.NotNil(s.t, s.request)
	requestBytes, err := ioutil.ReadAll(s.request.Body)
	require.NoError(s.t, err)
	require.Equal(s.t, body, requestBytes)
}

// AssertReceivedBody allows callers to assert that the last request received matches the JSON string provided.
func (s *TestService) AssertReceivedJSON(expected string) {
	require.NotNil(s.t, s.request)
	requestBytes, err := ioutil.ReadAll(s.request.Body)
	require.NoError(s.t, err)
	require.JSONEq(s.t, expected, string(requestBytes))
}

// AssertReceivedBody allows callers to assert that the last request received matches the HTTP Method provided.
func (s *TestService) AssertReceivedMethod(method string) {
	require.NotNil(s.t, s.request)
	require.Equal(s.t, method, s.request.Method)
}

// AssertReceivedAs allows callers to assert that the last request received could be unmarshalled as the value provided.
func (s *TestService) AssertReceivedAs(v interface{}) {
	requestBytes, err := ioutil.ReadAll(s.request.Body)
	require.NoError(s.t, err)
	require.NoError(s.t, json.Unmarshal(requestBytes, v))
}

// NewTestService will return a test service configured for the provided test state.
//
// This init also handles the starting and cleanup of the test service once the test is completed.
func NewTestService(t *testing.T) *TestService {
	s := new(TestService)
	s.t = t
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		s.request = r

		if s.responseCode != 0 {
			w.WriteHeader(s.responseCode)
		}

		if len(s.responseBody) != 0 {
			w.Write(s.responseBody)
		}
	}))

	t.Cleanup(s.Close)
	return s
}
