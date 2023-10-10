package xata_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testService(t *testing.T, method string, basePath string, statusCode int, shouldErr bool, want any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.HasPrefix(request.URL.Path, basePath) {
			t.Fatalf("unexpected path: %s, expected: %s", request.URL.Path, basePath)
		}

		if request.Method != method {
			t.Fatalf("expected method %s, got %s", method, request.Method)
		}

		writer.WriteHeader(statusCode)

		var response []byte
		var err error
		if shouldErr {
			response, err = json.Marshal(testErrBody)
			if err != nil {
				t.Fatal(err)
			}
		} else if want != nil {
			response, err = json.Marshal(want)
			if err != nil {
				t.Fatal(err)
			}
		}

		_, err = writer.Write(response)
		if err != nil {
			t.Fatal(err)
		}
	}))
}
