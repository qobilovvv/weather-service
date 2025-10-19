package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	jobs, err := initJobs(s)
	if err != nil {
		panic(err)
	}
	
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Starting server at port -> 3000")
		err := http.ListenAndServe(":3000", r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		fmt.Printf("Starting job: %v\n", jobs[0].ID())
		s.Start()
	}()
	
	wg.Wait()
}

func initJobs(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("helo world")
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return []gocron.Job{j}, nil
}

// func runCron() {
// 	// each job has a unique id
// 	fmt.Println(j.ID())

// 	// start the scheduler
// 	s.Start()
// }
