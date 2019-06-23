package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab.com/rhermes/packtrack/trackers/bring"
)

func trackingGenerator(out chan<- string) {
	log.Printf("Started tracking number generator\n")
	for i := 1000001; i < 1000015; i++ {
		out <- fmt.Sprintf("%d", i)
	}
	log.Printf("Stopped tracking number generator\n")
}

func main() {
	bc, err := bring.New(bring.Config{
		Workers:      3,
		InputBuffer:  2,
		OutputBuffer: 1,
		ErrorBuffer:  1,
		RateLimitDur: 500 * time.Millisecond,
	})
	if err != nil {
		log.Fatalf("error with opening bring: %s\n", err.Error())
	}
	defer bc.Close()

	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatalf("error with opening database: %s\n", err.Error())
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("Error with pinging database: %s\n", err.Error())
	}

	cn, err := bring.NewConnector(db, bc)

	if err != nil {
		log.Fatalf("error with creating connector: %s\n", err.Error())
	}
	defer cn.Close()

	go trackingGenerator(bc.Inputs())
}
