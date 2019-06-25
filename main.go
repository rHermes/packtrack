package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"gitlab.com/rhermes/packtrack/store"
)

var (
	NodeID           = flag.String("nodeid", "", "the nodeid for this node")
	Tracker          = flag.String("tracker", "", "the name of the tracker we will be using")
	InsertRange      = flag.Bool("range", false, "shall we use range mode?")
	InsertRangeStart = flag.Int64("rangeStart", -1, "The start of the insert range")
	InsertRangeEnd   = flag.Int64("rangeEnd", -1, "The end of the insert range")
	PerformMode      = flag.Bool("perform", false, "Shall we use perform mode")
)

func insertJob(s *store.Store, tracker int, start, stop int64) error {
	trackers := make([]int, 0)
	args := make([][]byte, 0)
	createdAt := make([]time.Time, 0)

	for i := start; i < stop; i++ {
		trackers = append(trackers, tracker)
		args = append(args, []byte(fmt.Sprintf(`{"q":"%d"}`, i)))
		createdAt = append(createdAt, time.Now())
	}
	startTime := createdAt[0]
	endTime := createdAt[len(createdAt)-1]
	dur := endTime.Sub(startTime)
	log.Printf("We spent %s building internal slices.\n", dur.String())

	beforeInsert := time.Now()
	if err := s.InsertJobs(trackers, args, createdAt); err != nil {
		return err
	}
	insertDur := time.Since(beforeInsert)
	log.Printf("We spent %s inserting into postgresql.\n", insertDur.String())

	return nil
}

func performJobs(s *store.Store) error {
	for {
		err := s.PerformJob()
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("There appears to be nothing to do, waiting 3 sec\n")
				time.Sleep(3 * time.Second)
			} else if err == store.ErrRateLimit {
				log.Printf("We have been ratelimited, waiting 10 minutes\n")
				time.Sleep(10 * time.Minute)
			} else {
				log.Printf("There was some other error, waiting 1 minute\n")
				time.Sleep(1 * time.Minute)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	flag.Parse()

	if *NodeID == "" {
		log.Fatalf("NodeID is required\n")
	}

	if *InsertRange {
		if *Tracker == "" {
			log.Fatalf("A tracker is required\n")
		}

		if *InsertRangeStart == -1 || *InsertRangeEnd == -1 {
			log.Fatalf("We need a range start and range end\n")
		}
		if *InsertRangeStart > *InsertRangeEnd {
			log.Fatalf("We need a range start smaller than the range end")
		}
	}

	s, err := store.New(store.Config{
		NodeID:     *NodeID,
		ConnString: "",
	})
	if err != nil {
		log.Fatalf("Error opening store: %s\n", err.Error())
	}
	defer s.Close()

	if *InsertRange {
		trackers, err := s.Trackers()
		if err != nil {
			log.Fatalf("Error getting trackers: %s\n", err.Error())
		}

		var bt store.Tracker

		for i, tracker := range trackers {
			fmt.Printf("Tracker %d: %#v\n", i, tracker)
			if tracker.Name == "bring" {
				bt = tracker
			}
		}

		if bt.ID == 0 {
			log.Fatalf("Didn't find the tracker we needed!\n")
		}

		if err := insertJob(s, bt.ID, *InsertRangeStart, *InsertRangeEnd); err != nil {
			log.Fatalf("Couldn't insert the tracker we needed: %s\n", err.Error())
		}
	}
	if *PerformMode {
		if err := performJobs(s); err != nil {
			log.Fatalf("Couldn't perform jobs: %s\n", err.Error())
		}
	}
}
