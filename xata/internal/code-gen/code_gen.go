// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// OpenAPI struct for parsing the OpenAPI YAML file
type OpenAPI struct {
	OpenAPI    string           `yaml:"openapi" json:"openapi"`
	Info       map[string]any   `yaml:"info" json:"info"`
	Servers    []map[string]any `yaml:"servers" json:"servers"`
	Paths      map[string]any   `yaml:"paths" json:"paths"`
	Components map[string]any   `yaml:"components" json:"components"`
	Tags       []map[string]any `yaml:"tags" json:"tags"`
	XTagGroups []map[string]any `yaml:"x-tagGroups" json:"x-tagGroups"`
}

type scope int

const (
	core scope = iota
	workspace
)

const (
	originalCorePath            = "xata/internal/fern-core"
	coreAPIspecs                = "internal/docs/core-openapi.json"
	workspaceAPIspecs           = "internal/docs/workspace-openapi.json"
	originalWorkspacePath       = "xata/internal/fern-workspace"
	scopeCore                   = "core"
	scopeWorkspace              = "workspace"
	coreGeneratorsYamlFile      = "core-generators.yml"
	workspaceGeneratorsYamlFile = "workspace-generators.yml"
	codeGenPath                 = "xata/internal/code-gen"

	fernGenerateCmd = "fern generate --log-level debug --local"
	fernInitCmd     = "fern init --openapi ../../../"
	newSuffix       = "-new"
)

var (
	newCorePath                 = originalCorePath + newSuffix
	newWorkspacePath            = originalWorkspacePath + newSuffix
	coreGeneratorsYamlPath      = codeGenPath + "/yaml-files/" + coreGeneratorsYamlFile
	workspaceGeneratorsYamlPath = codeGenPath + "/yaml-files/" + workspaceGeneratorsYamlFile
)

func main() {
	scope := flag.String("scope", "", "scope is one of: core or workspace")
	flag.Parse()

	switch *scope {
	case scopeCore:
		err := generateFernCode(core, newCorePath, originalCorePath, coreAPIspecs, coreGeneratorsYamlPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("success generating core code")
	case scopeWorkspace:
		err := generateFernCode(workspace, newWorkspacePath, originalWorkspacePath, workspaceAPIspecs, workspaceGeneratorsYamlPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("success generating workspace code")
	default:
		log.Fatal("unknown scope: ", *scope)
	}
}

func generateFernCode(scope scope, newPath, originalPath, apiSPECS, generatorsYAML string) error {
	rootWD, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get root wd: %v", err)
	}

	log.Println("creating new folder")
	err = os.Mkdir(newPath, 0o755)
	if err != nil {
		return fmt.Errorf("unable to create %v: %v", newPath, err)
	}

	err = os.Chdir(newPath)
	if err != nil {
		return fmt.Errorf("unable to change dir to %v: %v", newPath, err)
	}

	newPathWD, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get %v wd: %v", newPathWD, err)
	}

	log.Println("initializing fern")
	output, err := executeOSCmd(fernInitCmd + apiSPECS)
	if err != nil {
		return fmt.Errorf("unable to fern init: %v", err)
	}

	log.Println(output)

	err = os.Chdir(rootWD)
	if err != nil {
		return fmt.Errorf("unable to change dir to root wd: %v", err)
	}

	log.Println("updating API specs")
	switch scope {
	case core:
		log.Println("no action needed for core")
	case workspace:
		err = updateWorkspaceAPISpecs(newWorkspacePath + "/fern/api/openapi/workspace-openapi.json")
		if err != nil {
			return fmt.Errorf("unable to update workspace API specs: %v", err)
		}
	}

	log.Println("updating the generators file")
	err = copyFile(generatorsYAML, newPath+"/fern/api/generators.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(newPathWD)
	if err != nil {
		return fmt.Errorf("unable to change dir to %v: %v", newPath, err)
	}

	log.Println("generating code")
	output, err = executeOSCmd(fernGenerateCmd)
	if err != nil {
		return fmt.Errorf("unable to generate code: %v", err)
	}

	log.Println(output)

	err = os.Chdir(rootWD)
	if err != nil {
		return fmt.Errorf("unable to get root wd: %v", err)
	}

	log.Println("updating auto gen code")
	switch scope {
	case core:
		err = copySelfFromUtils("core.go", newPath+"/generated/go/core/")
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}
	case workspace:
		newPathGenGo := newPath + "/generated/go/"
		// messed up auto-gen code
		err = copySelfFromUtils("value_booster_value.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		// https://github.com/xataio/xata-go/issues/30
		err = copySelfFromUtils("insert_record_response.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("bulk_insert_table_records_response.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("update_record_with_id_response.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("insert_record_with_id_response.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("upsert_record_with_id_response.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("record.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("files_client.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("records_client.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("get_file_response.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		err = copySelfFromUtils("transaction_success_results_item.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		// https://github.com/xataio/xata-go/issues/31
		err = copySelfFromUtils("column_type.go", newPathGenGo)
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}

		// https://github.com/xataio/xata-go/issues/22
		err = copySelfFromUtils("core.go", newPath+"/generated/go/core/")
		if err != nil {
			return fmt.Errorf("unable to copy self: %v", err)
		}
	}

	err = os.RemoveAll(originalPath)
	if err != nil {
		return fmt.Errorf("unable to remove the original folder: %v", err)
	}

	err = os.Rename(newPath, originalPath)
	if err != nil {
		return fmt.Errorf("unable to rename the new folder as the original folder")
	}

	return nil
}

func executeOSCmd(cmd string) (string, error) {
	cmdDissected := strings.Fields(cmd)

	command := exec.Command(cmdDissected[0], cmdDissected[1:]...)

	output, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func copySelfFromUtils(fileName, newPath string) error {
	suffix := "_"
	return copyFile(codeGenPath+"/go-files/"+fileName+suffix, newPath+fileName)
}

func copyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func updateWorkspaceAPISpecs(filePath string) error {
	// Read the OpenAPI YAML file
	openAPIData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read OpenAPI file: %v", err)
	}

	// Unmarshal the OpenAPI data into a struct
	var openAPI OpenAPI

	if err = json.Unmarshal(openAPIData, &openAPI); err != nil {
		log.Println("error unmarshaling OpenAPI data:", err)
		return err
	}

	log.Println("removing deprecated path definitions")
	// https://github.com/xataio/xata-go/issues/29
	delete(openAPI.Paths, "/db/{db_branch_name}/migrations")
	delete(openAPI.Paths, "/db/{db_branch_name}/migrations/execute")
	delete(openAPI.Paths, "/db/{db_branch_name}/migrations/plan")

	log.Println("updating file[] as fileMap")
	// https://github.com/xataio/xata-go/issues/31
	columnEnums := openAPI.Components["schemas"].(map[string]any)["Column"].(map[string]any)["properties"].(map[string]any)["type"].(map[string]any)["enum"].([]any)
	var newColumnEnums []string
	for _, e := range columnEnums {
		if e.(string) == "file[]" {
			newColumnEnums = append(newColumnEnums, "fileMap")
		} else {
			newColumnEnums = append(newColumnEnums, e.(string))
		}
	}

	openAPI.Components["schemas"].(map[string]any)["Column"].(map[string]any)["properties"].(map[string]any)["type"].(map[string]any)["enum"] = newColumnEnums

	for k, v := range openAPI.Components["schemas"].(map[string]any)["Column"].(map[string]any)["properties"].(map[string]any) {
		if k == "file[]" {
			openAPI.Components["schemas"].(map[string]any)["Column"].(map[string]any)["properties"].(map[string]any)["fileMap"] = v
			delete(openAPI.Components["schemas"].(map[string]any)["Column"].(map[string]any)["properties"].(map[string]any), "file[]")
		}
	}

	log.Println("removing reference of ProjectionConfig from QueryColumnsProjection")
	// https://github.com/xataio/xata-go/issues/27
	var newOneOf []any
	for _, o := range openAPI.Components["schemas"].(map[string]any)["QueryColumnsProjection"].(map[string]any)["items"].(map[string]any)["oneOf"].([]any) {
		if o.(map[string]any)["$ref"] != nil && o.(map[string]any)["$ref"] == "#/components/schemas/ProjectionConfig" {
			continue
		}

		newOneOf = append(newOneOf, o)
	}
	openAPI.Components["schemas"].(map[string]any)["QueryColumnsProjection"].(map[string]any)["items"].(map[string]any)["oneOf"] = newOneOf

	log.Println("remove object value")
	// https://github.com/xataio/xata-go/issues/27
	delete(openAPI.Components["schemas"].(map[string]any), "ObjectValue")

	log.Println("removing reference for object value")
	var newAnyOf []any
	for _, o := range openAPI.Components["schemas"].(map[string]any)["DataInputRecord"].(map[string]any)["additionalProperties"].(map[string]any)["anyOf"].([]any) {
		if o.(map[string]any)["$ref"] != nil && o.(map[string]any)["$ref"] == "#/components/schemas/ObjectValue" {
			continue
		}

		newAnyOf = append(newAnyOf, o)
	}
	openAPI.Components["schemas"].(map[string]any)["DataInputRecord"].(map[string]any)["additionalProperties"].(map[string]any)["anyOf"] = newAnyOf

	updatedOpenAPIData, err := json.Marshal(&openAPI)
	if err != nil {
		return fmt.Errorf("unable to marshal updated OpenAPI data: %v", err)
	}

	// Save the updated OpenAPI data to a new file
	err = os.WriteFile(filePath, updatedOpenAPIData, 0o644)
	if err != nil {
		return fmt.Errorf("unable to save updated OpenAPI file: %v", err)
	}

	log.Print("OpenAPI file updated and saved as", filePath)
	return nil
}
