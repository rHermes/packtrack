// Copyright (c) 2019 Teodor Spæren
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	ErrRateLimit = errors.New("RateLimited")
)

const sqlGetTrackers = `SELECT id, name, description, url FROM trackers`

const sqlCreateScrapeJob = `
INSERT INTO
	scrape_jobs (
		tracker,
		args,
		created_at
	)
VALUES
	($1, $2, $3)
`

const sqlGetJobForUpdateByTracker = `
SELECT
	id,
	args
FROM
	scrape_jobs
WHERE
	status = 'created'
	AND
	tracker = $1
ORDER BY
	id ASC
FOR UPDATE
SKIP LOCKED
LIMIT 1
`

const sqlGetJobForUpdate = `
SELECT
	id,
	tracker,
	args
FROM
	scrape_jobs
WHERE
	status = 'created'
ORDER BY
	id ASC
FOR UPDATE
SKIP LOCKED
LIMIT 1
`

const sqlUpdateJob = `
UPDATE
	scrape_jobs
SET 
	status = $2,
	start_time = $3,
	end_time = $4,
	stats = $5,
	resp = $6
WHERE
	id = $1
`

type Tracker struct {
	ID          int
	Name        string
	Description string
	URL         string
}

// PerformStats is the statics we give back, when we finish a perform job
type PerformStats struct {
	NodeID string
}

// Config is used to configure a Store
type Config struct {
	NodeID     string
	ConnString string
}

// Store gives the ability to create, get and perform work
type Store struct {
	db *sql.DB
	id string

	prepGetTrackers              *sql.Stmt
	prepGetJobForUpdateByTracker *sql.Stmt
	prepGetJobForUpdate          *sql.Stmt
	prepUpdateJob                *sql.Stmt
	prepCreateScrapeJob          *sql.Stmt
}

// New creates a new store
func New(cfg Config) (*Store, error) {
	db, err := sql.Open("postgres", cfg.ConnString)
	if err != nil {
		return nil, err
	}

	prepGetTrackers, err := db.PrepareContext(context.Background(), sqlGetTrackers)
	if err != nil {
		return nil, err
	}

	prepGetJobForUpdateByTracker, err := db.PrepareContext(context.Background(), sqlGetJobForUpdateByTracker)
	if err != nil {
		return nil, err
	}
	prepGetJobForUpdate, err := db.PrepareContext(context.Background(), sqlGetJobForUpdate)
	if err != nil {
		return nil, err
	}

	prepUpdateJob, err := db.PrepareContext(context.Background(), sqlUpdateJob)
	if err != nil {
		return nil, err
	}

	prepCreateScrapeJob, err := db.PrepareContext(context.Background(), sqlCreateScrapeJob)
	if err != nil {
		return nil, err
	}

	s := &Store{
		id: cfg.NodeID,
		db: db,

		prepGetTrackers:              prepGetTrackers,
		prepGetJobForUpdateByTracker: prepGetJobForUpdateByTracker,
		prepGetJobForUpdate:          prepGetJobForUpdate,
		prepUpdateJob:                prepUpdateJob,
		prepCreateScrapeJob:          prepCreateScrapeJob,
	}
	return s, nil
}

// Close shuts the store down.
func (s *Store) Close() error {
	// TODO(rHermes): Report on the multi error that can occur here
	s.prepGetTrackers.Close()
	s.prepGetJobForUpdateByTracker.Close()
	s.prepGetJobForUpdate.Close()
	s.prepUpdateJob.Close()
	s.prepCreateScrapeJob.Close()
	return s.db.Close()
}

func (s *Store) Trackers() ([]Tracker, error) {
	rows, err := s.prepGetTrackers.QueryContext(context.Background())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trackers := make([]Tracker, 0)
	for rows.Next() {
		var track Tracker

		if err := rows.Scan(&(track.ID), &(track.Name), &(track.Description), &(track.URL)); err != nil {
			return nil, err
		}
		trackers = append(trackers, track)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return trackers, nil
}

// InsertJob creates a single job
func (s *Store) InsertJob(tracker int, args []byte, createdAt time.Time) error {
	_, err := s.prepCreateScrapeJob.ExecContext(context.Background(), tracker, args, createdAt)
	return err
}

// InsertJobs inserts all the jobs or non of them at all into the queue
func (s *Store) InsertJobs(tracker []int, args [][]byte, createdAt []time.Time) error {
	if len(tracker) != len(args) || len(args) != len(createdAt) {
		return errors.New("All arrays must be equally long")
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn("scrape_jobs", "tracker", "args", "created_at"))
	if err != nil {
		return err
	}

	for i := 0; i < len(tracker); i++ {
		if i%100 == 0 {
			log.Printf("On insert %d of %d aka %.2f%%.\n", i, len(tracker), float64(i)/float64(len(tracker))*100)
		}
		_, err = stmt.ExecContext(context.Background(), tracker[i], args[i], createdAt[i])
		if err != nil {
			return err
		}
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	if err := stmt.Close(); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (s *Store) PerformJob() error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	pGetJobForUpdate := tx.StmtContext(context.Background(), s.prepGetJobForUpdate)
	pUpdateJob := tx.StmtContext(context.Background(), s.prepUpdateJob)

	var id int
	var tracker int
	var jargs []byte

	row := pGetJobForUpdate.QueryRowContext(context.Background())
	if err := row.Scan(&id, &tracker, &jargs); err != nil {
		return err
	}
	startedAt := time.Now()

	// TODO(rHermes): This is cheating, I know that it bring has id 1
	if tracker != 1 {
		return errors.New("There is not supposed to be an other tracker ids than 1")
	}

	var workargs struct {
		Q string
	}
	if err := json.Unmarshal(jargs, &workargs); err != nil {
		// TODO(rHermes): Fail job here?
		return err
	}

	// TODO(rHermes): Make this into something that can handle proper ids
	u := "https://tracking.bring.com/api/v2/tracking.json?q=" + workargs.Q
	log.Printf("We will work on: %s\n", u)
	resp, err := http.Get(u)
	if err != nil {
		// Update job here?
		fmt.Printf("There was an error: %s\n", err.Error())
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	err2 := resp.Body.Close()
	if err != nil {
		// TODO(rHermes): Fail the query here?
		return err
	}
	if err2 != nil {
		// TODO(rHermes): Fail the query here?
		return err2
	}

	var ape struct {
		APIVersion     string `json:"apiVersion"`
		ConsignmentSet []struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		} `json:"consignmentSet"`
	}
	if err := json.Unmarshal(data, &ape); err == nil {
		// if code is not zero then it's ok
		for _, v := range ape.ConsignmentSet {
			if v.Error.Code != 0 && v.Error.Code != 404 {
				if v.Error.Code == 503 {
					return ErrRateLimit
				}
			}
		}
	}

	stat := &PerformStats{NodeID: s.id}
	statb, err := json.Marshal(stat)
	if err != nil {
		// TODO(rHermes): Fail the query here?
		return err
	}

	completedAt := time.Now()
	_, err = pUpdateJob.ExecContext(context.Background(), id, "success", startedAt, completedAt, statb, data)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
