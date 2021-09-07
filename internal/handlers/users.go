package handlers

import (
	"context"
	"fmt"
	"github.com/rw88/adjust-coding-challenge/internal/users"
	"strconv"
)

type ListUsersCommand struct {
	Name string
	Sort []string
	Limit int
}

// ListUsersHandler handles ListUsersCommand commands
func ListUsersHandler(ctx context.Context, cmd *ListUsersCommand) error {

	topUsers := users.ProcessUsers(cmd.Sort)

	for i, a := range topUsers {
		fmt.Println("#" + strconv.Itoa(i+1) + " User:" + a.Username + " " + " #PRs: " + strconv.Itoa(a.AmountPRs) + " #Commits:" + strconv.Itoa(a.AmountCommits))
	}

	return nil
}




