package users

import (
	"fmt"
	"github.com/rw88/adjust-coding-challenge/internal/configuration"
	"strconv"
	"testing"
)

func TestUsersSortedByPrsAndCommits(t *testing.T) {

	// Given a predefined set of testdata (actors, events, and repos)
	configuration.DataDirectory = "../../tests/testdata/"

	// Given a set of expected users
	expectedUsers := userCollection{
		{
			Username: "awesomekling",
			AmountPRs: 2,
		},
		{
			Username: "AdrianWilczynski",
			AmountPRs: 1,
		},
		{
			Username: "Apexal",
			AmountCommits: 1,
		},
		{
			Username: "ArturoCamacho0",
		},
		{
			Username: "onosendi",
			AmountCommits: 1,
		},
		{
			Username: "anggi1234",
			AmountCommits: 1,
		},
		{
			Username: "PercussiveElbow",
		},
		{
			Username: "Niceguy",
		},
		{
			Username: "John",
		},
		{
			Username: "Maria",
		},
	}

	// When I verify the top 10 users sorted by amount of PRs performed, and then sorted by commits performed
	topUsers := ProcessUsers([]string{"pr", "commit"}, 10)


	for i, topUser := range topUsers {

		if topUser.Username != expectedUsers[i].Username {
			t.Errorf("expected the Username '%s', got '%s'", expectedUsers[i].Username, topUser.Username)
		}

		if topUser.AmountCommits != expectedUsers[i].AmountCommits {
			t.Errorf("expected the #AmountCommits '%d', got '%d'",
				expectedUsers[i].AmountCommits,
				topUser.AmountCommits)
		}
	}





	for i, a := range topUsers {
		fmt.Println("#" + strconv.Itoa(i+1) + " User:" + a.Username + " " + " #PRs: " + strconv.Itoa(a.AmountPRs) + " #Commits:" + strconv.Itoa(a.AmountCommits))
	}

}
