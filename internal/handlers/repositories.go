package handlers

import (
	"context"
	"fmt"
	"github.com/rw88/adjust-coding-challenge/internal/repositories"
	"strconv"
)

type ListRepositoriesCommand struct {
	Name string
	Sort []string
	Limit int
}
// ListRepositoriesHandler handles ListRepositoriesCommand commands
func ListRepositoriesHandler(ctx context.Context, cmd *ListRepositoriesCommand) error {

	topRepos := repositories.ProcessRepositories(cmd.Sort, cmd.Limit)

	for i, r := range topRepos {
		fmt.Println("#" + strconv.Itoa(i+1) + " Repo:" + r.Name + " " + " #Watch Events: " + strconv.Itoa(r.AmountWatchEvents) + " #Commits:" + strconv.Itoa(r.AmountCommits))
	}

	return nil
}



