package xata

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientOptions_getAPIKey(t *testing.T) {
	apiKeyFromEnv := "test-API-key-from-env"
	err := os.Setenv(xataAPIKeyEnvVar, apiKeyFromEnv)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Unsetenv(xataAPIKeyEnvVar) })

	t.Run("should assign the API key from the env vars", func(t *testing.T) {
		apiKey, err := getAPIKey()
		assert.NoError(t, err)

		assert.Equal(t, apiKey, apiKeyFromEnv)
	})
}

func Test_parseDatabaseURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    databaseConfig
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "should parse successfully without the branch dbName",
			rawURL: "https://my-workspace-id.us-east-1.xata.sh/db/test-db",
			want: databaseConfig{
				workspaceID:     "my-workspace-id",
				region:          "us-east-1",
				domainWorkspace: "xata.sh",
				dbName:          "test-db",
				branchName:      "main",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:   "should parse successfully with the branch dbName",
			rawURL: "https://my-workspace-id.us-east-1.xata.sh/db/test-db:feature-branch",
			want: databaseConfig{
				workspaceID:     "my-workspace-id",
				region:          "us-east-1",
				domainWorkspace: "xata.sh",
				dbName:          "test-db",
				branchName:      "feature-branch",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name:   "should fail when host is not in expected format",
			rawURL: "https://unexpected-field.my-workspace-id.us-east-1.xata.sh/db/test-db:feature-branch",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid databaseConfig URL")
			},
		},
		{
			name:   "should fail when path is not in expected format",
			rawURL: "https://my-workspace-id.us-east-1.xata.sh/db/unexpected-field/test-db:feature-branch",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid databaseConfig URL")
			},
		},
		{
			name:   "should fail when branch dbName is missing even there is a colon after the db dbName",
			rawURL: "https://my-workspace-id.us-east-1.xata.sh/db/test-db:",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "invalid databaseConfig URL")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDatabaseURL(tt.rawURL)
			if !tt.wantErr(t, err, fmt.Sprintf("parseDatabaseURL(%v)", tt.rawURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseDatabaseURL(%v)", tt.rawURL)
		})
	}
}

func Test_loadConfig(t *testing.T) {
	t.Run("should read database URL", func(t *testing.T) {
		// Create a temporary JSON file for testing
		tempFile, err := os.CreateTemp("", "config_test.json")
		if err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}
		t.Cleanup(func() {
			os.Remove(tempFile.Name()) // Clean up the temporary file
		})

		// Write test JSON data to the temporary file
		testData := `{"databaseURL": "test-database-url"}`
		_, err = tempFile.WriteString(testData)
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
		err = tempFile.Close()
		if err != nil {
			t.Fatal(err)
		}

		// Test loading the configuration
		cfg, err := loadConfig(tempFile.Name())
		if err != nil {
			t.Fatalf("Error loading config: %v", err)
		}

		expectedURL := "test-database-url"
		if cfg.DatabaseURL != expectedURL {
			t.Fatalf("Expected database URL: %s, got: %s", expectedURL, cfg.DatabaseURL)
		}
	})
}
