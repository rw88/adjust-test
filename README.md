# Test assignment for Analytics Software Engineer position

This repo contains Github event data for 1 hour.

Please write a CLI application that outputs:

- Top 10 active users sorted by amount of PRs created and commits pushed
- Top 10 repositories sorted by amount of commits pushed
- Top 10 repositories sorted by amount of watch events

The CLI interface definition is up to you.

This assignment must be written in Golang.
Please don't use any type of database or data processing engines (such as Apache Spark), this data should be processed by your application.

What we want to see in your solution:

Your code should be:
    - Readable
    - Testable
    - Maintainable
    - Extensible
    - Modifiable
    - Performant

Please be prepared to answer questions on how you achieved those goals and which trade offs you made.

Please also keep an eye on:

    - Tests
    - Structured, meaningful commits
    - a Readme on how to run the solution



## Prerequisites
Have GO installed in your environment. At least version 1.16.

At the project root, install the project dependencies.
```
go mod vendor
```

Unzip the data.tar.gz file so that the application has access to the test files.
```
tar -xf data.tar.gz
```

## Usage

At the project root, run the following commands:

Top 10 active users sorted by amount of PRs created and commits pushed
```
go run ./cmd/analytics/analytics.go -sort="pr,commit" -limit=10 users
```

Top 10 repositories sorted by amount of commits pushed
```
go run ./cmd/analytics/analytics.go -sort="commit" -limit=10 repositories
```

Top 10 repositories sorted by amount of watch events
```
go run ./cmd/analytics/analytics.go -sort="watch_events" -limit=10 repositories
```