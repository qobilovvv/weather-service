package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
	"github.com/qobilovvv/weather_service/internal/client/geocoding"
	openmeteo "github.com/qobilovvv/weather_service/internal/client/open_meteo"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	
	geocodingClient := geocoding.NewClient(httpClient)
	openMeteoClient := openmeteo.NewClient(httpClient)

	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")
		fmt.Println("requested city: ", city)
		
		geores, err := geocodingClient.GetCoords(city)
		if err != nil {
			log.Println(err)
		}
		
		openMetRes, err := openMeteoClient.GetTemperature(geores.Latitude, geores.Longitude)
		if err != nil {
			log.Println(err)
			return
		}
		
		raw, err := json.Marshal(openMetRes)
		if err != nil {
			log.Println(err)
		}
		
		_, err = w.Write(raw)
		
		if err != nil {
			log.Println(err)
		}
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
