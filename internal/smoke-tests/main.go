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

	workspacesClient, err := xata.NewWorkspacesClient(
		// xata.WithAPIKey("wrong token"),
		httpCli,
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := workspacesClient.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, ws := range resp.Workspaces {
		fmt.Printf("%#v\n", *ws)
		fmt.Printf("%s\n", ws.Role.String())
	}

	var workSpaceIDToBeDeleted string
	workspace, err := workspacesClient.Create(ctx, &xata.WorkspaceMeta{Name: "test-ws"})
	if err != nil {
		log.Fatal(err)
	}

	workSpaceIDToBeDeleted = workspace.Id
	fmt.Println("ws id to delete", workSpaceIDToBeDeleted)

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

	recordsCli, err := xata.NewRecordsClient(
		httpCli,
	)
	if err != nil {
		log.Fatal(err)
	}

	tableName := "first-table"
	insReq := xata.InsertRecordRequest{
		RecordRequest: xata.RecordRequest{TableName: tableName},
		Columns:       []string{"user-name"},
		Body: map[string]*xata.DataInputRecordValue{
			"user-name": xata.ValueFromString("test-value-from-SDK-smoke-test"),
		},
	}
	recordResp, err := recordsCli.Insert(ctx, insReq)
	if err != nil {
		log.Fatal(err)
	}

	if recordResp.Data["user-name"] != "test-value-from-SDK-smoke-test" {
		log.Fatal("unexpected response")
	}

	record, err := recordsCli.Get(ctx, xata.GetRecordRequest{
		RecordRequest: xata.RecordRequest{TableName: tableName},
		RecordID:      recordResp.Id,
	})
	if err != nil {
		log.Fatal(err)
	}

	if record.Id != recordResp.Id {
		log.Fatal("unexpected ID")
	}

	tableCli, err := xata.NewTableClient(
		httpCli,
	)
	if err != nil {
		log.Fatal(err)
	}

	testTableName := "my-test-table-smoke-test"
	createTableResponse, err := tableCli.Create(ctx, xata.TableRequest{
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

	columnName := "test-column"
	_, err = tableCli.AddColumn(ctx, xata.AddColumnRequest{
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

	_, err = tableCli.DeleteColumn(ctx, xata.DeleteColumnRequest{
		TableRequest: xata.TableRequest{TableName: testTableName},
		ColumnName:   columnName,
	})
	if err != nil {
		log.Fatal(err)
	}

	delTableResponse, err := tableCli.Delete(ctx, xata.TableRequest{
		TableName: testTableName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table status", delTableResponse.Status.String())

	usersCli, err := xata.NewUsersClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	user, err := usersCli.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user)

	dbCli, err := xata.NewDatabasesClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	databases, err := dbCli.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases.Databases[0].Name)

	branchCli, err := xata.NewBranchClient(httpCli)
	if err != nil {
		log.Fatal(err)
	}

	branches, err := branchCli.List(ctx, databases.Databases[0].Name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("%#v\n", branches.Branches[0].Name))

	newBranchName := "new-branch-from-smoke-test-2"
	createBranchRes, err := branchCli.Create(ctx, xata.CreateBranchRequest{
		BranchName: newBranchName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(createBranchRes.BranchName)
	fmt.Println(createBranchRes.DatabaseName)
	fmt.Println(createBranchRes.Status)

}
