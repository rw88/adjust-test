package events

import (
	"github.com/rw88/adjust-coding-challenge/internal/eventbus"
	"github.com/rw88/adjust-coding-challenge/internal/readers"
	"log"
	"os"
)

func ReadEvents(eb *eventbus.EventBus) <-chan struct{} {

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
