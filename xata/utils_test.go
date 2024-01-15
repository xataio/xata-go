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
		gotBranchName := getBranchName()
		assert.Equal(t, gotBranchName, defaultBranchName)
	})

	setBranchName := "feature-042"
	err := os.Setenv(EnvXataBranch, setBranchName)
	if err != nil {
		t.Fatal(err)
	}

	// from env var
	t.Run("should be branch name from env var", func(t *testing.T) {
		gotBranchName := getBranchName()
		assert.Equal(t, gotBranchName, setBranchName)
	})

	t.Cleanup(func() { os.Unsetenv(EnvXataBranch) })
}

func Test_getRegion(t *testing.T) {
	// default state
	t.Run("should be default region", func(t *testing.T) {
		gotRegion := getRegion()
		assert.Equal(t, gotRegion, defaultRegion)
	})

	setRegion := "eu-west-3"
	err := os.Setenv(EnvXataRegion, setRegion)
	if err != nil {
		t.Fatal(err)
	}

	// from env var
	t.Run("should be region from the env var", func(t *testing.T) {
		gotRegion := getRegion()
		assert.Equal(t, gotRegion, setRegion)
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
	err := os.Setenv(EnvXataWorkspaceID, setWsId)
	if err != nil {
		t.Fatal(err)
	}

	// test workspace id from env var
	t.Run("load config from WORKSPACE_ID env var", func(t *testing.T) {
		dbCfg, err := loadDatabaseConfig()
		if err != nil {
			t.Fatalf("Error loading config: %v", err)
		}

		if dbCfg.workspaceID != setWsId {
			t.Fatalf("Expected Workspace ID: %s, got: %s", setWsId, dbCfg.workspaceID)
		}
		if dbCfg.branchName != defaultBranchName {
			t.Fatalf("Expected branch name: %s, got: %s", defaultBranchName, dbCfg.branchName)
		}
		if dbCfg.region != defaultRegion {
			t.Fatalf("Expected region: %s, got: %s", defaultRegion, dbCfg.region)
		}
	})

	setBranch := "branch123"
	err2 := os.Setenv(EnvXataBranch, setBranch)
	if err2 != nil {
		t.Fatal(err2)
	}
	setRegion := "ap-southeast-16"
	err3 := os.Setenv(EnvXataRegion, setRegion)
	if err3 != nil {
		t.Fatal(err3)
	}

	// with branch and region env vars
	t.Run("load config from XATA_WORKSPACE_ID, regionEnvVar and XATA_BRANCH env vars", func(t *testing.T) {
		dbCfg, err := loadDatabaseConfig()
		if err != nil {
			t.Fatalf("Error loading config: %v", err)
		}

		if dbCfg.workspaceID != setWsId {
			t.Fatalf("Expected Workspace ID: %s, got: %s", setWsId, dbCfg.workspaceID)
		}
		if dbCfg.branchName != setBranch {
			t.Fatalf("Expected branch name: %s, got: %s", setBranch, dbCfg.branchName)
		}
		if dbCfg.region != setRegion {
			t.Fatalf("Expected region: %s, got: %s", setRegion, dbCfg.region)
		}
	})

	t.Cleanup(func() {
		os.Unsetenv(EnvXataWorkspaceID)
		os.Unsetenv(EnvXataBranch)
		os.Unsetenv(EnvXataRegion)
	})
}
