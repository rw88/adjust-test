package handlers

import (
	"context"
	"fmt"
	"github.com/rw88/adjust-coding-challenge/internal/eventbus"
	"github.com/rw88/adjust-coding-challenge/internal/readers"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)
type userCollection []*User

func (uc userCollection) sortByFields(fields []string)  {

	if len(fields) == 0 {
		return
	}

	sort.Slice(uc, func(i, j int) bool {
		for _, sort := range fields {
			switch sort {
			case "commit":
				return uc[i].amountCommits > uc[j].amountCommits
			case "pr":
				return uc[i].amountPRs > uc[j].amountPRs
			}
		}

		return false
	})
}

var userList userCollection


type ListUsersCommand struct {
	Name string
	Sort []string
	Limit int
}

// ListUsersHandler handles ListUsersCommand commands
func ListUsersHandler(ctx context.Context, cmd *ListUsersCommand) error {

	eb := eventbus.NewEventBus()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	readActors(wg)
	wg.Wait()

	pushEventChan := make(chan eventbus.EventData)
	prCreatedChan := make(chan eventbus.EventData)
	eb.Subscribe("PushEvent", pushEventChan)
	eb.Subscribe("PullRequestEvent", prCreatedChan)

	doneReadEvents := readEvents(eb)

	for  {
		select {
		case <-doneReadEvents:

			userList.sortByFields(cmd.Sort)

			for i, a := range userList[:10] {
				fmt.Println("#" + strconv.Itoa(i+1) + " User:" + a.username + " " + " #PRs: " + strconv.Itoa(a.amountPRs) + " #Commits:" + strconv.Itoa(a.amountCommits))
			}

			return nil
		case data := <- pushEventChan:
			line := data.Data.([]string)
			if user, ok := users[line[2]]; ok {
				user.amountCommits++
			}
		case data := <- prCreatedChan:
			line := data.Data.([]string)
			if user, ok := users[line[2]]; ok {
				user.amountPRs++
			}
		}
	}


	return nil
}


func readActors(wg *sync.WaitGroup)  {
	defer wg.Done()

	csvFile, err := os.Open("data/actors.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	ch := make(chan []string)
	go readers.CSVReader(csvFile, ch)


	for line := range ch {

		if _, ok := users[line[0]]; !ok {
			users[line[0]] = &User{
				id:       line[0],
				username: line[1],
			}
			userList = append(userList, users[line[0]])
		}
	}

}


func readEvents(eb *eventbus.EventBus) <-chan struct{} {

	done := make(chan struct{})

	go func() {

		defer close(done)

		csvFile, err := os.Open("data/events.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()

		ch := make(chan []string)
		go readers.CSVReader(csvFile, ch)


		for line := range ch {
			eb.Publish(line[1], line)
		}
	}()

	return done
}


type User struct {
	id string
	username string
	active bool
	amountPRs int
	amountCommits int
}

func(u *User) amountActivity() int {

	return u.amountPRs + u.amountCommits
}


var users = map[string]*User{}