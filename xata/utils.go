// SPDX-License-Identifier: Apache-2.0

package xata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	personalAPIKeyLocation    = "~/.config/xata/key"
	defaultControlPlaneDomain = "api.xata.io"
	xataAPIKeyEnvVar          = "XATA_API_KEY"
	xataWsIDEnvVar            = "XATA_WORKSPACE_ID" // TODO: not in use yet
	dbURLFormat               = "https://{workspace_id}.{region}.xata.sh/db/{db_name}:{branch_name}"
	defaultBranchName         = "main"
	configFileName            = ".xatarc"
	branchNameEnvVar          = "XATA_BRANCH"
	defaultDataPlaneDomain    = "xata.sh"
	defaultRegion             = "us-east-1"
)

var errAPIKey = fmt.Errorf("no API key found. Searched in `%s` env, %s, and .env", xataAPIKeyEnvVar, personalAPIKeyLocation)

// assignAPIkey add the API key to the ClientOptions by going through the following options in order:
//   - In env vars by the xataAPIKeyEnvVar dbName.
//   - In .env file with the xataAPIKeyEnvVar dbName.
//   - In .xatarc config file (TODO: not ready!)
//
// See: https://xata.io/docs/python-sdk/overview#authorization
func getAPIKey() (string, error) {
	if key, found := os.LookupEnv(xataAPIKeyEnvVar); found {
		return key, nil
	}

	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			return "", err
		}
	}

	if key, found := myEnv[xataAPIKeyEnvVar]; found {
		return key, nil
	}

	// Look for the in file created by the xata CLI
	// python code: https://github.com/xataio/xata-py/blob/main/xata/client.py#L44
	//	looks for "~/.config/xata/key" file.
	// But the documents and xata cli tests shows that the generated file is actually
	//	 .config/xata/credentials
	// https://xata.io/docs/getting-started/cli#authentication-profiles
	// clarify this!

	return "", errAPIKey
}

func String(in string) *string {
	if in == "" {
		return nil
	}
	return &in
}

func Bool(in bool) *bool {
	return &in
}

func Int(in int) *int {
	return &in
}

func Float64(in float64) *float64 {
	return &in
}

func Uint8(in uint8) *uint8 {
	return &in
}

type databaseConfig struct {
	workspaceID     string
	region          string
	dbName          string
	branchName      string
	domainWorkspace string
}

// parseDatabaseURL parses a given DB URL.
//
//	Branch dbName is optional.
//	Format: https://{workspace_id}.{region}.xata.sh/db/{db_name}:{branch_name}
func parseDatabaseURL(rawURL string) (databaseConfig, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return databaseConfig{}, err
	}

	host := strings.Split(parsedURL.Host, ".")
	if len(host) != 4 {
		return databaseConfig{}, fmt.Errorf("invalid databaseConfig URL: %s, expected format %s", rawURL, dbURLFormat)
	}

	path := strings.Split(parsedURL.Path, "/")
	if len(path) != 3 {
		return databaseConfig{}, fmt.Errorf("invalid databaseConfig URL: %s, expected format %s", rawURL, dbURLFormat)
	}

	db := databaseConfig{
		workspaceID:     host[0],
		region:          host[1],
		domainWorkspace: fmt.Sprintf("%s.%s", host[2], host[3]),
		dbName:          path[2],
	}

	if strings.Contains(db.dbName, ":") {
		names := strings.Split(db.dbName, ":")
		if names[1] == "" {
			return databaseConfig{}, fmt.Errorf("invalid databaseConfig URL: %s, expected format %s", rawURL, dbURLFormat)
		}

		db.dbName = names[0]
		db.branchName = names[1]
	}

	if db.branchName == "" {
		db.branchName = getBranchName()
	}

	return db, err
}

// config represents the JSON configuration.
type config struct {
	DatabaseURL string `json:"databaseURL"`
}

// loadConfig reads the JSON file and extracts the databaseURL.
func loadConfig(fieName string) (config, error) {
	// Read the JSON file
	data, err := os.ReadFile(fieName)
	if err != nil {
		return config{}, err
	}

	// Parse the JSON data into a cfg struct
	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config{}, err
	}

	return cfg, nil
}

// getBranchName retrieves the branch name.
// If not found, falls back to defaultBranchName
func getBranchName() string {
	if branchName, found := os.LookupEnv(branchNameEnvVar); found {
		return branchName
	}

	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			return ""
		}
	}

	if branchName, found := myEnv[branchNameEnvVar]; found {
		return branchName
	}

	return defaultBranchName
}

// Get the region if the corresponding env var `XATA_REGION` is set
// otherwise return the default region: us-east-1
func getRegion() string {
	if region, found := os.LookupEnv("XATA_REGION"); found {
		return region
	}
	return defaultRegion
}

// loadDatabaseConfig will return config with defaults if the error is not nil.
func loadDatabaseConfig() (databaseConfig, error) {
	defaultDBConfig := databaseConfig{
		region:          defaultRegion,
		branchName:      defaultBranchName,
		domainWorkspace: defaultDataPlaneDomain,
	}
	// Setup with env var
	// XATA_WORKSPACE_ID to set the workspace Id
	if wsID, found := os.LookupEnv("XATA_WORKSPACE_ID"); found {
		region := getRegion()
		branch := getBranchName()
		db := databaseConfig{
			workspaceID: wsID,
			region:      region,
			branchName:  branch,
		}
		return db, nil
	}

	cfg, err := loadConfig(configFileName)
	if err != nil {
		return defaultDBConfig, err
	}

	dbCfg, err := parseDatabaseURL(cfg.DatabaseURL)
	if err != nil {
		return defaultDBConfig, err
	}

	return dbCfg, nil
}
