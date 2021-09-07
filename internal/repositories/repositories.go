package repositories

import (
	"github.com/rw88/adjust-coding-challenge/internal/eventbus"
	"github.com/rw88/adjust-coding-challenge/internal/events"
	"github.com/rw88/adjust-coding-challenge/internal/readers"
	"github.com/rw88/adjust-coding-challenge/internal/users"
	"log"
	"os"
	"sort"
	"sync"
)

type repoCollection []*Repository
var repoList repoCollection
var repos = map[string]*Repository{}

type Repository struct {
	Id string
	Name string
	AmountWatchEvents int
	AmountCommits int
}

func ProcessRepositories(sort []string) repoCollection  {

	wg := &sync.WaitGroup{}
	wg.Add(2)
	users.ReadUsersFile(wg)
	ReadRepos(wg)
	wg.Wait()

	eb := eventbus.NewEventBus()
	doneReadEvents := events.ReadEvents(eb)

	pushEventChan := make(chan eventbus.EventData)
	watchEventChan := make(chan eventbus.EventData)
	eb.Subscribe("PushEvent", pushEventChan)
	eb.Subscribe("WatchEvent", watchEventChan)


	for  {
		select {
		case <-doneReadEvents:

			repoList.sortByFields(sort)

			return repoList[:10]
		case data := <- pushEventChan:
			line := data.Data.([]string)
			if repo, ok := repos[line[3]]; ok {
				repo.AmountCommits++
			}
		case data := <- watchEventChan:
			line := data.Data.([]string)
			if repo, ok := repos[line[3]]; ok {
				repo.AmountWatchEvents++
			}
		}
	}
}



func (rc repoCollection) sortByFields(fields []string)  {

	if len(fields) == 0 {
		return
	}

	sort.Slice(rc, func(i, j int) bool {
		for _, sort := range fields {
			switch sort {
			case "commit":
				return rc[i].AmountCommits > rc[j].AmountCommits
			case "watch_event":
				return rc[i].AmountWatchEvents > rc[j].AmountWatchEvents
			}
		}

		return false
	})
}

func ReadRepos(wg *sync.WaitGroup)  {
	defer wg.Done()

	csvFile, err := os.Open("data/repos.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	ch := make(chan []string)
	go readers.CSVReader(csvFile, ch)


	for line := range ch {
		repos[line[0]] = &Repository{
			Id: line[0],
			Name: line[1],
		}
		repoList = append(repoList, repos[line[0]])
	}

}

