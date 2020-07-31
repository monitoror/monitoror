package main

import (
	"context"
	"fmt"

	"github.com/jsdidierlaurent/go-bitbucket"
)

func main() {
	auth := context.WithValue(context.Background(), bitbucket.ContextBasicAuth, bitbucket.BasicAuth{
		UserName: "jsdidierlaurent",
		Password: "uZfHt4y4pFeWT2ghgSvT",
	})

	cfg := bitbucket.NewConfiguration()
	client := bitbucket.NewAPIClient(cfg)

	// Load PR
	var pullRequests []bitbucket.Pullrequest
	results, _, _ := client.PullrequestsApi.RepositoriesUsernameRepoSlugPullrequestsGet(auth, "jsdidierlaurent", "bitbucket-test", map[string]interface{}{"pagelen": int32(50)})
	pullRequests = append(pullRequests, results.Values...)

	for results.Next != "" {
		results, _, _ = client.PagingApi.PullrequestsPageGet(auth, results.Next)
		pullRequests = append(pullRequests, results.Values...)
	}

	for _, pullRequest := range pullRequests {
		fmt.Printf("PR : %s : %s\n", pullRequest.Title, pullRequest.State)

		commit, _, _ := client.CommitstatusesApi.RepositoriesUsernameRepoSlugCommitNodeStatusesGet(auth, "jsdidierlaurent", pullRequest.Source.Commit.Hash, "bitbucket-test")
		fmt.Println("  - Build: " + commit.Values[0].State)
	}

	branch, _, _ := client.RefsApi.RepositoriesUsernameRepoSlugRefsBranchesNameGet(auth, "jsdidierlaurent", "master", "bitbucket-test")
	fmt.Printf("Branch : %s\n", branch.Name)

	commit, _, _ := client.CommitstatusesApi.RepositoriesUsernameRepoSlugCommitNodeStatusesGet(auth, "jsdidierlaurent", branch.Target.Hash, "bitbucket-test")
	fmt.Println("  - Build: " + commit.Values[0].State)
}
