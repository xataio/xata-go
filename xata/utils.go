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

// Environment variables
const (
	EnvXataAPIKey      = "XATA_API_KEY"
	EnvXataWorkspaceID = "XATA_WORKSPACE_ID"
	EnvXataBranch      = "XATA_BRANCH"
	EnvXataRegion      = "XATA_REGION"
)

const (
	personalAPIKeyLocation    = "~/.config/xata/key"
	defaultControlPlaneDomain = "api.xata.io"
	dbURLFormat               = "https://{workspace_id}.{region}.xata.sh/db/{db_name}:{branch_name}"
	defaultBranchName         = "main"
	configFileName            = ".xatarc"
	defaultDataPlaneDomain    = "xata.sh"
	defaultRegion             = "us-east-1"
)

var errAPIKey = fmt.Errorf("no API key found. Searched in `%s` env, %s, and .env", EnvXataAPIKey, personalAPIKeyLocation)

// assignAPIkey add the API key to the ClientOptions by going through the following options in order:
//   - In env vars by the EnvXataAPIKey dbName.
//   - In .env file with the EnvXataAPIKey dbName.
//   - In .xatarc config file (TODO: not ready!)
//
// See: https://xata.io/docs/python-sdk/overview#authorization
func getAPIKey() (string, error) {
	if key, found := os.LookupEnv(EnvXataAPIKey); found {
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

	if key, found := myEnv[EnvXataAPIKey]; found {
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
		db.branchName = getBranchName(nil)
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

// Get value from env var with fallback to godotenv
// return default value if not found
func getEnvVar(name string, defaultValue string) string {
	if val, found := os.LookupEnv(name); found {
		return val
	}

	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		var pathError *fs.PathError
		if !errors.As(err, &pathError) {
			return ""
		}
	}

	if val, found := myEnv[name]; found {
		return val
	}
	return defaultValue
}

// getBranchName retrieves the branch name. If not found, falls back to defaultBranchName.
func getBranchName(opts *ClientOptions) string {
	if opts != nil && opts.Branch != "" {
		return opts.Branch
	}
	return getEnvVar(EnvXataBranch, defaultBranchName)
}

// getRegion gets the region if the corresponding env var `XATA_REGION` is set otherwise return
// defaultRegion.
func getRegion(opts *ClientOptions) string {
	if opts != nil && opts.Region != "" {
		return opts.Region
	}
	return getEnvVar(EnvXataRegion, defaultRegion)
}

// getWorkspaceID gets the workspace id from opts and if empty, gets it from the `XATA_WORKSPACE_ID`
// environment variable
func getWorkspaceID(opts *ClientOptions) string {
	if opts != nil && opts.WorkspaceID != "" {
		return opts.WorkspaceID
	}
	return getEnvVar(EnvXataWorkspaceID, "")
}

// loadDatabaseConfig will return config with defaults if the error is not nil.
func loadDatabaseConfig(cliOpts *ClientOptions) (databaseConfig, error) {
	defaultDBConfig := databaseConfig{
		region:          defaultRegion,
		branchName:      defaultBranchName,
		domainWorkspace: defaultDataPlaneDomain,
	}

	// Config can come from three places with differing priorities. The order from highest to lowest
	// priority is:
	// 1. Code via ClientOptions
	// 2. Environment variables
	// 3. Config files

	wsID := getWorkspaceID(cliOpts)
	if wsID != "" {
		db := databaseConfig{
			workspaceID:     wsID,
			region:          getRegion(cliOpts),
			branchName:      getBranchName(cliOpts),
			domainWorkspace: defaultDataPlaneDomain,
		}
		return db, nil
	}

	// Config not found in code or environment variables, fall back to config files
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
