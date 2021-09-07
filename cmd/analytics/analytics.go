package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lana/go-commandbus"
	"github.com/rw88/adjust-coding-challenge/internal/handlers"
	"log"
	"os"
	"strings"
)

const usage = `analytics

Usage:
  analytics -sort="pr,commit" -limit=10 users
  analytics -sort="commit" -limit=10 repositories

Options:

  -sort <sort>                          Sort option [default: ""]
  -limit <int>                          Limit; how many items to retrieve [default: 10]
  -help 
`

var commandNameFlag = flag.String("command_name", "", "Command name (required)")
var sortFlag = flag.String("sort", "", "Sort by a field (or multiple fields), separated by | (optional)")
var limitFlag = flag.Int("limit", 10, "Limit the results. Default 10")
var helpFlag = flag.Bool("help", false, "Show help")

var argToNewCommand = map[string]func() interface{} {
	"users": NewListUsersCommand,
	"repositories": NewListRepositoriesCommand,
}

func main()  {
	flag.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if *helpFlag {
		showUsageHelp()
	}

	args := flag.Args()

	if len(args) < 1 {
		showUsageHelp()
		return
	}


	bus := commandbus.New()
	registerHandlers(bus)


	var command interface{}
	if newCommandCallback, ok := argToNewCommand[args[0]]; ok {
		command = newCommandCallback()
	}
	if command == nil {
		log.Fatal("Please provide a valid command name")
	}

	err := bus.Execute(context.Background(), command)
	if err != nil {
		log.Fatal(err)
	}
}

// registerHandlers registers handlers to specific command objects
func registerHandlers(bus commandbus.CommandBus)  {

	err := bus.Register(&handlers.ListUsersCommand{}, handlers.ListUsersHandler)
	if err != nil {
		log.Fatal(err)
	}

	err = bus.Register(&handlers.ListRepositoriesCommand{}, handlers.ListRepositoriesHandler)
	if err != nil {
		log.Fatal(err)
	}
}

// NewListRepositoriesCommand creates a new ListRepositoriesCommand
func NewListRepositoriesCommand() interface{} {

	return &handlers.ListRepositoriesCommand{
		Name: "list-repositories",
		Sort: getSortOption(),
		Limit: getLimitOption(),
	}
}
// NewListUsersCommand creates a new ListUsersCommand
func NewListUsersCommand() interface{} {

	return &handlers.ListUsersCommand{
		Name: "list-users",
		Sort: getSortOption(),
		Limit: getLimitOption(),
	}
}

// getSortOption retrieves the cli -sort option
func getSortOption() []string  {
	return strings.Split(*sortFlag, ",")
}

// getLimitOption retrieves the cli -limit option
func getLimitOption() int  {
	return *limitFlag
}

// showUsageHelp it shows a usage help and stop the application
func showUsageHelp() {
	flag.Usage()
	os.Exit(1)
}