package testaservice_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jslang/testaservice"
	"github.com/stretchr/testify/require"
)

var _ = func() bool {
	gofakeit.Seed(0)
	return true
}

type TestServiceData struct {
	*testaservice.TestService
	*http.Client
	t *testing.T
}

func NewTestServiceData(t *testing.T) *TestServiceData {
	td := new(TestServiceData)
	td.t = t
	td.TestService = testaservice.NewTestService(t)
	td.Client = &http.Client{}
	return td
}

func (td *TestServiceData) ServiceURL(p string) string {
	u, err := url.Parse(td.TestService.URL)
	require.NoError(td.t, err)

	u.Path = path.Join(u.Path, p)
	return u.String()
}

func TestTestService(t *testing.T) {

	t.Run("AssertReceivedBasicAuth", func(t *testing.T) {
		td := NewTestServiceData(t)

		username := "username"
		password := "password"
		request, _ := http.NewRequest(http.MethodGet, td.ServiceURL("/"), nil)
		request.SetBasicAuth(username, password)
		_, err := td.Client.Do(request)
		require.NoError(t, err)

		td.TestService.AssertReceivedBasicAuth(username, password)
	})

	t.Run("AssertNotCalled", func(t *testing.T) {
		td := NewTestServiceData(t)
		td.TestService.AssertNotCalled()
	})

	t.Run("AssertCalled", func(t *testing.T) {
		td := NewTestServiceData(t)
		_, err := td.Client.Get(td.ServiceURL("/"))
		require.NoError(t, err)

		td.TestService.AssertCalled()
	})

	t.Run("AssertReceivedPath", func(t *testing.T) {
		td := NewTestServiceData(t)

		path := "/" + gofakeit.BeerName()
		_, err := td.Client.Get(td.ServiceURL(path))
		require.NoError(t, err)

		td.TestService.AssertReceivedPath(path)
	})

	t.Run("AssertReceivedHeader", func(t *testing.T) {
		td := NewTestServiceData(t)

		header := gofakeit.HackerNoun()
		value := gofakeit.HackerAdjective()

		request, _ := http.NewRequest(http.MethodGet, td.ServiceURL("/"), nil)
		request.Header.Set(header, value)

		_, err := td.Client.Do(request)
		require.NoError(t, err)
		td.TestService.AssertReceivedHeader(header, value)
	})

	t.Run("AssertReceivedParam", func(t *testing.T) {
		td := NewTestServiceData(t)

		param := gofakeit.HackerNoun()
		value := gofakeit.HackerAdjective()

		url := fmt.Sprintf("%s?%s=%s", td.ServiceURL("/"), param, value)
		request, _ := http.NewRequest(http.MethodGet, url, nil)

		_, err := td.Client.Do(request)
		require.NoError(t, err)
		td.TestService.AssertReceivedParam(param, value)
	})

	t.Run("AssertReceivedBody", func(t *testing.T) {
		td := NewTestServiceData(t)

		body := []byte(gofakeit.HackerPhrase())
		request, _ := http.NewRequest(http.MethodGet, td.ServiceURL("/"), bytes.NewReader(body))

		_, err := td.Client.Do(request)
		require.NoError(t, err)

		td.TestService.AssertReceivedBody(body)
	})

	t.Run("AssertReceivedJSON", func(t *testing.T) {
		td := NewTestServiceData(t)

		body, _ := json.Marshal(struct {
			AField string
			BField string
		}{
			AField: gofakeit.BS(),
			BField: gofakeit.HackerPhrase(),
		})
		request, _ := http.NewRequest(http.MethodGet, td.ServiceURL("/"), bytes.NewReader(body))

		_, err := td.Client.Do(request)
		require.NoError(t, err)

		td.TestService.AssertReceivedJSON(string(body))
	})

	t.Run("AssertReceivedMethod", func(t *testing.T) {
		td := NewTestServiceData(t)

		method := gofakeit.HTTPMethod()
		request, _ := http.NewRequest(method, td.ServiceURL("/"), nil)

		_, err := td.Client.Do(request)
		require.NoError(t, err)

		td.TestService.AssertReceivedMethod(method)
	})

	t.Run("AssertReceivedAs", func(t *testing.T) {
		td := NewTestServiceData(t)

		type testType struct {
			FieldA string
			FieldB int
			FieldC float32
		}

		var requestData testType
		gofakeit.Struct(&requestData)

		body, _ := json.Marshal(requestData)

		request, _ := http.NewRequest(http.MethodPost, td.ServiceURL("/"), bytes.NewReader(body))

		_, err := td.Client.Do(request)
		require.NoError(t, err)

		var testData testType
		td.TestService.AssertReceivedAs(&testData)
	})
}
