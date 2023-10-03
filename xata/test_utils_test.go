package xata_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

func prepareConfigFile() error {
	sourcePath := "xatarc_test"
	destPath := ".xatarc"
	// Open the source file.
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file.
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents of the source file to the destination file.
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func cleanConfigFile() error {
	// Check if the old file exists.
	if _, err := os.Stat(".xatarc"); err != nil {
		return err
	}

	// Rename the file.
	err := os.Remove(".xatarc")
	if err != nil {
		return err
	}

	return nil
}
