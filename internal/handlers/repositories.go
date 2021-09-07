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

type ListRepositoriesCommand struct {
	Name string
	Sort []string
	Limit int
}

type repoCollection []*Repository

func (rc repoCollection) sortByFields(fields []string)  {

	if len(fields) == 0 {
		return
	}

	sort.Slice(rc, func(i, j int) bool {
		for _, sort := range fields {
			switch sort {
			case "commit":
				return rc[i].amountCommits > rc[j].amountCommits
			case "watch_event":
				return rc[i].amountWatchEvents > rc[j].amountWatchEvents
			}
		}

		return false
	})
}

var repoList repoCollection

type Repository struct {
	id string
	name string
	amountWatchEvents int
	amountCommits int
}

var repos = map[string]*Repository{}

// ListRepositoriesHandler handles ListRepositoriesCommand commands
func ListRepositoriesHandler(ctx context.Context, cmd *ListRepositoriesCommand) error {

	wg := &sync.WaitGroup{}
	wg.Add(2)
	readActors(wg)
	readRepos(wg)
	wg.Wait()

	eb := eventbus.NewEventBus()
	doneReadEvents := readEvents(eb)

	pushEventChan := make(chan eventbus.EventData)
	watchEventChan := make(chan eventbus.EventData)
	eb.Subscribe("PushEvent", pushEventChan)
	eb.Subscribe("WatchEvent", watchEventChan)


	for  {
		select {
		case <-doneReadEvents:

			repoList.sortByFields(cmd.Sort)

			for i, r := range repoList[:10] {
				fmt.Println("#" + strconv.Itoa(i+1) + " Repo:" + r.name + " " + " #Watch Events: " + strconv.Itoa(r.amountWatchEvents) + " #Commits:" + strconv.Itoa(r.amountCommits))
			}

			return nil
		case data := <- pushEventChan:
			line := data.Data.([]string)
			if repo, ok := repos[line[3]]; ok {
				repo.amountCommits++
			}
		case data := <- watchEventChan:
			line := data.Data.([]string)
			if repo, ok := repos[line[3]]; ok {
				repo.amountWatchEvents++
			}
		}
	}

	return nil
}


func readRepos(wg *sync.WaitGroup)  {
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
			id: line[0],
			name: line[1],
		}
		repoList = append(repoList, repos[line[0]])
	}

}



