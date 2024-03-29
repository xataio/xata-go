// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientOptions_getAPIKey(t *testing.T) {
	apiKeyFromEnv := "test-API-key-from-env"
	err := os.Setenv(EnvXataAPIKey, apiKeyFromEnv)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Unsetenv(EnvXataAPIKey) })

	t.Run("should assign the API key from the env vars", func(t *testing.T) {
		apiKey, err := getAPIKey()
		assert.NoError(t, err)

		assert.Equal(t, apiKey, apiKeyFromEnv)
	})
}

func Test_getBranchName(t *testing.T) {
	// default state
	t.Run("should be default branch name", func(t *testing.T) {
		gotBranchName := getBranchName(nil)
		assert.Equal(t, gotBranchName, defaultBranchName)
	})

	setBranchName := "feature-042"
	setEnvForTests(t, EnvXataBranch, setBranchName)

	// from env var
	t.Run("should be branch name from env var", func(t *testing.T) {
		gotBranchName := getBranchName(nil)
		assert.Equal(t, setBranchName, gotBranchName)
	})

	// from ClientOptions
	t.Run("from ClientOptions", func(t *testing.T) {
		want := "branch-from-opts"
		got := getBranchName(&ClientOptions{Branch: want})
		assert.Equal(t, want, got)
	})
}

func Test_getRegion(t *testing.T) {
	// default state
	t.Run("should be default region", func(t *testing.T) {
		gotRegion := getRegion(nil)
		assert.Equal(t, gotRegion, defaultRegion)
	})

	setRegion := "eu-west-3"
	err := os.Setenv(EnvXataRegion, setRegion)
	if err != nil {
		t.Fatal(err)
	}

	// from env var
	t.Run("should be region from the env var", func(t *testing.T) {
		gotRegion := getRegion(nil)
		assert.Equal(t, setRegion, gotRegion)
	})

	t.Run("should be region from ClientOptions", func(t *testing.T) {
		wantRegion := "region-options"
		gotRegion := getRegion(&ClientOptions{Region: wantRegion})
		assert.Equal(t, wantRegion, gotRegion)
	})

	t.Cleanup(func() { os.Unsetenv(EnvXataRegion) })
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
	// from .xatarc
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

func Test_loadDatabaseConfig_with_envvars(t *testing.T) {
	setWsId := "workspace-0lac00"

	assertClientOptions := func(t *testing.T) {
		wsID := "workspace-fco"
		opts := &ClientOptions{
			WorkspaceID: wsID,
		}
		dbCfg, err := loadDatabaseConfig(opts)
		assert.NoError(t, err)
		assert.Equal(t, wsID, dbCfg.workspaceID)
		assert.Equal(t, defaultBranchName, dbCfg.branchName)
		assert.Equal(t, defaultRegion, dbCfg.region)
		assert.Equal(t, defaultDataPlaneDomain, dbCfg.domainWorkspace)
	}

	// test workspace is from ClientOptions
	t.Run("load config from ClientOptions", assertClientOptions)

	setEnvForTests(t, EnvXataWorkspaceID, setWsId)

	// Check again after environment variable set
	t.Run("config from ClientOptions takes precedence", assertClientOptions)

	// test workspace id from env var
	t.Run("load config from WORKSPACE_ID env var", func(t *testing.T) {
		dbCfg, err := loadDatabaseConfig(nil)
		assert.NoError(t, err)
		assert.Equal(t, setWsId, dbCfg.workspaceID)
		assert.Equal(t, defaultBranchName, dbCfg.branchName)
		assert.Equal(t, defaultRegion, dbCfg.region)
		assert.Equal(t, defaultDataPlaneDomain, dbCfg.domainWorkspace)
	})

	setBranch := "branch123"
	setEnvForTests(t, EnvXataBranch, setBranch)
	setRegion := "ap-southeast-16"
	setEnvForTests(t, EnvXataRegion, setRegion)

	// with branch and region env vars
	t.Run("load config from XATA_WORKSPACE_ID, XATA_REGION and XATA_BRANCH env vars", func(t *testing.T) {
		dbCfg, err := loadDatabaseConfig(nil)
		assert.NoError(t, err)
		assert.Equal(t, setWsId, dbCfg.workspaceID)
		assert.Equal(t, setBranch, dbCfg.branchName)
		assert.Equal(t, setRegion, dbCfg.region)
	})
}

func setEnvForTests(t *testing.T, key, value string) {
	t.Helper()
	err := os.Setenv(key, value)
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, os.Unsetenv(key))
	})
}
