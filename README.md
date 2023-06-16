# asynq-example
This project is a simple example using the asynq package, which implements a distributed task queue in Go.

There are three parts to this project:
* Worker: When a task is sent to a queue, it consumes and handles the task.
* Client: Simulate real users by sending asynchronous tasks.
* Beat: Executes periodic tasks and automatically sends specific tasks.

# How to use

Run the redis server via docker

`docker-compose up -d`

Run the each server respectively.

`TYPE=worker go run main.go`

`TYPE=client go run main.go`

`TYPE=beat go run main.go`
