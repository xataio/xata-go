package integrationtests

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xataio/xata-go/xata"
)

func Test_branchClient(t *testing.T) {
	cfg, err := setupDatabase()
	if err != nil {
		t.Fatalf("unable to setup db: %v", err)
	}

	t.Cleanup(func() {
		err := cleanup(cfg)
		if err != nil {
			t.Fatalf("unable to cleanup test setup: %v", err)
		}
	})

	branchCli, err := xata.NewBranchClient(
		xata.WithAPIKey(cfg.apiKey),
		xata.WithHTTPClient(cfg.httpCli),
		xata.WithBaseURL(cfg.workspaceCliBaseURL),
	)
	if err != nil {
		log.Fatalf("unable to construct a branch client: %v", err)
	}

	ctx := context.Background()

	t.Run("should create get list delete branch", func(t *testing.T) {
		newBranchFromMain := "new-branch-from-main_" + cfg.testID
		createBranchRes, err := branchCli.Create(ctx, xata.CreateBranchRequest{
			BranchName:   newBranchFromMain,
			DatabaseName: xata.String(cfg.databaseName),
		})
		if err != nil {
			log.Fatalf("unable to create a branch: %v", err)
		}
		assert.Equal(t, newBranchFromMain, createBranchRes.BranchName)

		t.Cleanup(func() {
			_, err = branchCli.Delete(ctx, xata.BranchRequest{
				DatabaseName: xata.String(cfg.databaseName),
				BranchName:   newBranchFromMain,
			})
			if err != nil {
				log.Fatal(err)
			}
		})

		listBranchRes, err := branchCli.List(ctx, cfg.databaseName)
		assert.NoError(t, err)
		var branchNames []string
		for _, branch := range listBranchRes.Branches {
			branchNames = append(branchNames, branch.Name)
		}
		assert.Contains(t, branchNames, newBranchFromMain)

		getBranchRes, err := branchCli.GetDetails(ctx, xata.BranchRequest{
			DatabaseName: xata.String(cfg.databaseName),
			BranchName:   newBranchFromMain,
		})
		assert.NoError(t, err)
		assert.Equal(t, newBranchFromMain, getBranchRes.BranchName)

		anotherBranchFromNewBranch := "very-new-branch"
		createBranchRes, err = branchCli.Create(ctx, xata.CreateBranchRequest{
			BranchName:   anotherBranchFromNewBranch,
			DatabaseName: xata.String(cfg.databaseName),
			From:         xata.String(newBranchFromMain),
		})
		assert.NoError(t, err)
		assert.Equal(t, anotherBranchFromNewBranch, createBranchRes.BranchName)
		t.Cleanup(func() {
			_, err = branchCli.Delete(ctx, xata.BranchRequest{
				DatabaseName: xata.String(cfg.databaseName),
				BranchName:   anotherBranchFromNewBranch,
			})
			if err != nil {
				log.Fatal(err)
			}
		})

		yetAnotherBranchFromVeryNewBranch := "super-new-branch"
		createBranchRes, err = branchCli.Create(ctx, xata.CreateBranchRequest{
			BranchName:   yetAnotherBranchFromVeryNewBranch,
			DatabaseName: xata.String(cfg.databaseName),
			Payload: &xata.CreateBranchRequestPayload{
				CreateBranchRequestFrom: xata.String(anotherBranchFromNewBranch),
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, yetAnotherBranchFromVeryNewBranch, createBranchRes.BranchName)
		t.Cleanup(func() {
			_, err = branchCli.Delete(ctx, xata.BranchRequest{
				DatabaseName: xata.String(cfg.databaseName),
				BranchName:   yetAnotherBranchFromVeryNewBranch,
			})
			if err != nil {
				log.Fatal(err)
			}
		})
	})
}
