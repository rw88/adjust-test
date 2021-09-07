package users

import (
	"github.com/rw88/adjust-coding-challenge/internal/configuration"
	"github.com/rw88/adjust-coding-challenge/internal/eventbus"
	"github.com/rw88/adjust-coding-challenge/internal/events"
	"github.com/rw88/adjust-coding-challenge/internal/readers"
	"log"
	"os"
	"sort"
	"sync"
)

var users = map[string]*User{}

type userCollection []*User
var userList userCollection

type User struct {
	Id string
	Username string
	Active bool
	AmountPRs int
	AmountCommits int
}

// ProcessUsers reads users and related files, and applies sorting criteria
func ProcessUsers(sorts []string, limit int) userCollection   {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ReadUsersFile(wg)
	wg.Wait()

	pushEventChan := make(chan eventbus.EventData)
	prCreatedChan := make(chan eventbus.EventData)

	eb := eventbus.NewEventBus()
	eb.Subscribe("PushEvent", pushEventChan)
	eb.Subscribe("PullRequestEvent", prCreatedChan)

	doneReadEvents := events.ReadEvents(eb)

	for  {
		select {
		case <-doneReadEvents:

			userList.sortByFields(sorts)

			if len(userList) < limit {
				return userList
			}

			return userList[:limit]
		case data := <- pushEventChan:
			line := data.Data.([]string)
			if user, ok := users[line[2]]; ok {
				user.AmountCommits++
			}
		case data := <- prCreatedChan:
			line := data.Data.([]string)
			if user, ok := users[line[2]]; ok {
				user.AmountPRs++
			}
		}
	}
}

// sortByFields sorts the collection by multiple fields
func (uc userCollection) sortByFields(fields []string)  {

	if len(fields) == 0 {
		return
	}

	sort.Slice(uc, func(i, j int) bool {
		for _, sort := range fields {
			switch sort {
			case "commit":
				return uc[i].AmountCommits > uc[j].AmountCommits
			case "pr":
				return uc[i].AmountPRs > uc[j].AmountPRs
			}
		}

		return false
	})
}



func ReadUsersFile(wg *sync.WaitGroup)  {

	defer wg.Done()

	csvFile, err := os.Open(configuration.DataDirectory + "actors.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	ch := make(chan []string)
	go readers.CSVReader(csvFile, ch)

	for line := range ch {

		if _, ok := users[line[0]]; !ok {
			users[line[0]] = &User{
				Id:       line[0],
				Username: line[1],
			}
			userList = append(userList, users[line[0]])
		}
	}

}