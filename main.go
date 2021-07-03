package main

import (
	"asynq-example/tasks"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hibiken/asynq"
)

const (
	typeWorker = "worker"
	typeClient = "client"
	typeBeat   = "beat"
)

var redisclient asynq.RedisClientOpt

// GetEnv helper
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// WorkerConsumer function
func WorkerConsumer() {
	worker := asynq.NewServer(redisclient, asynq.Config{
		// Specify how many concurrent workers to use.
		Concurrency: 10,
		// Optionally specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})
	mux := asynq.NewServeMux()
	// Define a task handler for the welcome email task.
	mux.HandleFunc(
		tasks.TypeHeartBeat,      // task type
		tasks.HandlHeartBeatTask, // handler function
	)

	// Run worker server.
	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}
}

// ClientProducer function is to simulate client send async tasks
func ClientProducer() {
	client := asynq.NewClient(redisclient)
	for id := 1; id < 5; id++ {

		task1 := tasks.HeartBeatTask(id)
		task2 := tasks.HeartBeatTask(id * 100)

		// Process the task immediately in critical queue.
		if _, err := client.Enqueue(
			task1,                   // task payload
			asynq.Queue("critical"), // set queue for task
		); err != nil {
			log.Fatal(err)
		}

		delay := 10 * time.Second
		if _, err := client.Enqueue(
			task2,                  // task payload
			asynq.Queue("low"),     // set queue for task
			asynq.ProcessIn(delay), // set time to process task
		); err != nil {
			log.Fatal(err)
		}
	}
}

// BeatProducer function
func BeatProducer() {
	scheduler := asynq.NewScheduler(
		redisclient,
		&asynq.SchedulerOpts{},
	)
	task1 := tasks.HeartBeatTask(1)
	task2 := tasks.HeartBeatTask(2)

	// You can use cron spec string to specify the schedule.
	// The cron parameter is (minute, hour, day, month, day of week)
	entryID, err := scheduler.Register("* * * * *", task1, asynq.Queue("critical"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an every minute entry: %q\n", entryID)

	// You can use "@every <duration>" to specify the interval.
	entryID, err = scheduler.Register("@every 60s", task2, asynq.Queue("default"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an every 60s entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	host := getEnv("REDISHOST", "localhost")
	port := getEnv("REDISPORT", "6379")
	addr := fmt.Sprintf("%s:%s", host, port)

	redisclient = asynq.RedisClientOpt{
		Addr: addr,
		DB:   1,
	}

	switch getEnv("TYPE", "worker") {
	case typeWorker:
		fmt.Printf("Go worker")
		WorkerConsumer()
	case typeClient:
		fmt.Printf("Go client, now time is %v", time.Now())
		ClientProducer()
	case typeBeat:
		fmt.Printf("Go beat, now time is %v", time.Now())
		BeatProducer()
	}
}
