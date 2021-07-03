# asynq-example
This project is a simple example with asynq package which implements a distributed task queue in Go.

There have three parts.
* Worker: When a task is sent to a queue, consume the task and handle it.
* Client: Simulate real users to send asynchronous task.
* Beat: Period task, send the specific tasks automatically.

# How to use

Run the redis server via docker

`docker-compose up -d`

Run the each server respectively.

`TYPE=worker go run main.go`

`TYPE=client go run main.go`

`TYPE=beat go run main.go`
