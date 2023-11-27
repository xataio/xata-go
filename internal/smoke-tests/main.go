// SPDX-License-Identifier: Apache-2.0

package main // nolint: typecheck

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/xataio/xata-go/xata"
)

func main() {
	ctx := context.Background()

	httpCli := xata.WithHTTPClient(retryablehttp.NewClient().StandardClient())

	//
	// User
	//
	usersClient, err := xata.NewUsersClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(("# fetch user"))
	user, err := usersClient.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user)

	//
	// Workspace
	//
	workspacesClient, err := xata.NewWorkspacesClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# List Workspaces")
	resp, err := workspacesClient.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, ws := range resp.Workspaces {
		fmt.Printf("%#v\n", *ws)
		fmt.Printf("%s\n", ws.Role.String())
	}

	log.Print("# Create new Workspace")
	var workSpaceIDToBeDeleted string
	workspace, err := workspacesClient.Create(ctx, &xata.WorkspaceMeta{Name: "xata-go-smoke-test"})
	if err != nil {
		log.Fatal(err)
	}

	workSpaceIDToBeDeleted = workspace.Id
	fmt.Println("ws id to delete", workSpaceIDToBeDeleted)

	log.Print("# Delete created Workspace")
	err = workspacesClient.Delete(ctx, workSpaceIDToBeDeleted)
	if err != nil {
		log.Fatal(err)
	}

	var isNotDeleted bool
	resp, err = workspacesClient.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, ws := range resp.Workspaces {
		if ws.Id == workSpaceIDToBeDeleted {
			isNotDeleted = true
		}
	}

	if isNotDeleted {
		log.Println("expected to be deleted but not")
	}

	//
	// Database
	//
	dbClient, err := xata.NewDatabasesClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# List all Databases of Workspace")
	databases, err := dbClient.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases.Databases[0].Name)

	log.Print("# Create new Database")
	// TODO

	log.Print("# Rename Database")
	// TODO

	log.Print("# Delete Database")
	// TODO

	//
	// Table
	//
	tableClient, err := xata.NewTableClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# Create new Table")
	testTableName := "my-test-table-smoke-test"
	createTableResponse, err := tableClient.Create(ctx, xata.TableRequest{
		TableName: testTableName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table name", createTableResponse.TableName)
	fmt.Println("table status", createTableResponse.Status.String())
	fmt.Println("table branch name", createTableResponse.BranchName)

	if createTableResponse.TableName != testTableName {
		log.Fatalf("unexpected table name: %v", createTableResponse.TableName)
	}

	log.Print("# Add column to table")
	columnName := "test-column"
	_, err = tableClient.AddColumn(ctx, xata.AddColumnRequest{
		TableRequest: xata.TableRequest{TableName: testTableName},
		Column: &xata.Column{
			Name:         columnName,
			Type:         xata.ColumnTypeString,
			NotNull:      xata.Bool(true),
			DefaultValue: xata.String("defaultValue"),
			Unique:       xata.Bool(false),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# Get Columns")
	// TODO

	log.Print("# Get Schema")
	// TODO

	log.Print("# Delete new column")
	_, err = tableClient.DeleteColumn(ctx, xata.DeleteColumnRequest{
		TableRequest: xata.TableRequest{TableName: testTableName},
		ColumnName:   columnName,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# Delete table")
	delTableResponse, err := tableClient.Delete(ctx, xata.TableRequest{
		TableName: testTableName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table status", delTableResponse.Status.String())

	//
	// Records
	//
	recordsClient, err := xata.NewRecordsClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# Insert new Record")
	tableName := "first-table"
	insReq := xata.InsertRecordRequest{
		RecordRequest: xata.RecordRequest{TableName: tableName},
		Columns:       []string{"user-name"},
		Body: map[string]*xata.DataInputRecordValue{
			"user-name": xata.ValueFromString("test-value-from-SDK-smoke-test"),
		},
	}
	recordResp, err := recordsClient.Insert(ctx, insReq)
	if err != nil {
		log.Fatal(err)
	}

	if recordResp.Data["user-name"] != "test-value-from-SDK-smoke-test" {
		log.Fatal("unexpected response")
	}

	log.Print("# Get the new Record")
	record, err := recordsClient.Get(ctx, xata.GetRecordRequest{
		RecordRequest: xata.RecordRequest{TableName: tableName},
		RecordID:      recordResp.Id,
	})
	if err != nil {
		log.Fatal(err)
	}

	if record.Id != recordResp.Id {
		log.Fatal("unexpected ID")
	}

	log.Print("# Delete the Record")
	err = recordsClient.Delete(ctx, xata.DeleteRecordRequest{
		RecordRequest: xata.RecordRequest{TableName: tableName},
		RecordID:      recordResp.Id,
	})
	if err != nil {
		log.Fatal(err)
	}

	//
	// Files
	//

	log.Print("# Add a File")
	// TODO

	log.Print("# Get the File")
	// TODO

	log.Print("# Delete the File")
	// TODO

	//
	// Branch
	//
	branchClient, err := xata.NewBranchClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# List all Branches")
	branches, err := branchClient.List(ctx, databases.Databases[0].Name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", branches.Branches[0].Name)

	log.Print("# Create new Branch of main branch")
	newBranchName := "new-branch-from-smoke-test-3"
	createBranchRes, err := branchClient.Create(ctx, xata.CreateBranchRequest{
		BranchName: newBranchName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(createBranchRes.BranchName)
	fmt.Println(createBranchRes.DatabaseName)
	fmt.Println(createBranchRes.Status)

	log.Print("# Create new Branch of specific branch")
	newBranchName2 := "new-branch-from-smoke-test-4"
	_, err = branchClient.Create(ctx, xata.CreateBranchRequest{
		BranchName: newBranchName2,
		From:       xata.String(newBranchName),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# Create new Branch with branch name in payload")
	newBranchName3 := "new-branch-from-smoke-test-5"
	_, err = branchClient.Create(ctx, xata.CreateBranchRequest{
		BranchName: newBranchName3,
		Payload: &xata.CreateBranchRequestPayload{
			CreateBranchRequestFrom: xata.String(newBranchName2),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Print("# Delete branch 1")
	delBranchRes, err := branchClient.Delete(ctx, xata.BranchRequest{
		BranchName: newBranchName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(delBranchRes.Status)

	log.Print("# Delete branch 2")
	delBranchRes, err = branchClient.Delete(ctx, xata.BranchRequest{
		BranchName: newBranchName2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(delBranchRes.Status)

	log.Print("# List branch 3")
	branchDetails, err := branchClient.GetDetails(ctx, xata.BranchRequest{
		BranchName: newBranchName3,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(branchDetails.BranchName)
	fmt.Println(branchDetails.Metadata)
	fmt.Println(branchDetails.DatabaseName)
	fmt.Println(branchDetails.Id)
	fmt.Println(branchDetails.StartedFrom)
	fmt.Println(branchDetails.Schema)

	log.Print("# Delete branch 3")
	delBranchRes, err = branchClient.Delete(ctx, xata.BranchRequest{
		BranchName: newBranchName3,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(delBranchRes.Status)
}
