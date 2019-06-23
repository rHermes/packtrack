package bring

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type CrawlResponse struct {
	Input  string
	Worker string
	Output []byte
}

type CrawlError struct {
	Input  string
	Worker string
	Error  error
}

type Config struct {
	Workers int

	InputBuffer  int
	OutputBuffer int
	ErrorBuffer  int

	RateLimitDur time.Duration
}

type Client struct {
	hc              http.Client
	wg              sync.WaitGroup
	chanInputs      chan string
	chanRateLimited chan string
	chanOutputs     chan CrawlResponse
	chanErrors      chan CrawlError
}

// New returns a new client which we can use to crawl
func New(cfg Config) (*Client, error) {
	chanInputs := make(chan string, cfg.InputBuffer)
	chanRateLimited := make(chan string)
	chanOutputs := make(chan CrawlResponse, cfg.OutputBuffer)
	chanErrors := make(chan CrawlError, cfg.ErrorBuffer)

	c := &Client{
		chanInputs:      chanInputs,
		chanOutputs:     chanOutputs,
		chanErrors:      chanErrors,
		chanRateLimited: chanRateLimited,
	}

	go c.runRateLimiter(cfg.RateLimitDur)

	for i := 0; i < cfg.Workers; i++ {
		c.wg.Add(1)
		go c.runWorker(fmt.Sprintf("worker-%02d", i))
	}

	return c, nil
}

// Inputs return a channel which can be used to
// give the crawler new requests
func (c *Client) Inputs() chan<- string { return c.chanInputs }

// Outputs returns a channel with the output from the crawlers
func (c *Client) Outputs() <-chan CrawlResponse { return c.chanOutputs }

// Errors returns a channel with errors from the crawlers
func (c *Client) Errors() <-chan CrawlError { return c.chanErrors }

// Close the scraping process
func (c *Client) Close() error {
	close(c.chanInputs)
	c.wg.Wait()
	close(c.chanErrors)
	close(c.chanOutputs)
	return nil
}

// runRateLimiter rate limits the run
func (c *Client) runRateLimiter(dur time.Duration) {
	lastSend := time.Now()
	for pid := range c.chanInputs {
		curDur := dur - time.Since(lastSend)
		time.Sleep(curDur)
		c.chanRateLimited <- pid
		lastSend = time.Now()
	}
	close(c.chanRateLimited)
}

func (c *Client) runWorker(id string) {
	log.Printf("[%s] worker started.\n", id)
	for pid := range c.chanRateLimited {
		log.Printf("[%s] Start processing package id %s\n", id, pid)

		// TODO(rHermes): Make this into something that can handle proper ids
		u := "https://tracking.bring.com/api/v2/tracking.json?q=" + pid
		resp, err := c.hc.Get(u)
		if err != nil {
			c.chanErrors <- CrawlError{pid, id, err}
			log.Printf("[%s] Stopped processing package id %s\n", id, pid)
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		err2 := resp.Body.Close()
		if err != nil {
			c.chanErrors <- CrawlError{pid, id, err}
			log.Printf("[%s] Stopped processing package id %s\n", id, pid)
			continue
		}
		if err2 != nil {
			c.chanErrors <- CrawlError{pid, id, err2}
			log.Printf("[%s] Stopped processing package id %s\n", id, pid)
			continue
		}

		cr := CrawlResponse{Input: pid, Worker: id, Output: data}
		c.chanOutputs <- cr

		log.Printf("[%s] Stopped processing package id %s\n", id, pid)
	}
	log.Printf("[%s] worker stopped.\n", id)
	c.wg.Done()
}
