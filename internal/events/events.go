package events

import (
	"github.com/rw88/adjust-coding-challenge/internal/configuration"
	"github.com/rw88/adjust-coding-challenge/internal/eventbus"
	"github.com/rw88/adjust-coding-challenge/internal/readers"
	"log"
	"os"
)

// ReadEvents reads event's file, and publishes event actions to the eventbus
func ReadEvents(eb *eventbus.EventBus) <-chan struct{} {

	done := make(chan struct{})

	go func() {

		defer close(done)

		csvFile, err := os.Open(configuration.DataDirectory + "events.csv")
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
