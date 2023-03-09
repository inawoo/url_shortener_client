package client

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type (
	FuncType string

	Work struct {
		ID     uuid.UUID `json:"id"`
		Input  any
		fnType FuncType
		Output any
		Error  error
	}

	WorkerFunction interface {
		ShortenURL(input ShortenURLRequest) (*URLCollection, error)
		CheckHealth() (string, error)
	}

	mutex struct {
		request         sync.Mutex
		result          sync.Mutex
		poolCount       sync.Mutex
		actualPoolCount sync.Mutex
		requestCount    sync.Mutex
	}

	Client struct {
		baseURL         string
		poolCount       int
		actualPoolCount int
		mtx             mutex
		request         chan Work
		requestCount    int
		result          map[uuid.UUID]Work
		stopSignal      chan bool
		workerFunctions WorkerFunction
	}

	Opts interface {
		Apply(*Client)
	}

	withPoolCount int
	withBaseURL   string
)

func WithPoolCount(count int) Opts {
	return withPoolCount(count)
}

func WithBaseURL(url string) Opts {
	return withBaseURL(url)
}

func (w withPoolCount) Apply(c *Client) {
	c.SetPoolCount(int(w))
}

func (w withBaseURL) Apply(c *Client) {
	c.baseURL = string(w)
}

func NewClient(opts ...Opts) *Client {

	client := &Client{
		baseURL:   "https://staging-link.inawo.live",
		poolCount: 1,
		mtx: mutex{
			request:         sync.Mutex{},
			result:          sync.Mutex{},
			poolCount:       sync.Mutex{},
			actualPoolCount: sync.Mutex{},
			requestCount:    sync.Mutex{},
		},
		request:    make(chan Work, 1),
		result:     make(map[uuid.UUID]Work),
		stopSignal: make(chan bool),
	}

	for _, opt := range opts {
		opt.Apply(client)
	}

	client.workerFunctions = NewWorkerFunctionsAdapter(client.baseURL)

	client.start()

	return client
}

func (c *Client) SetPoolCount(count int) {
	c.mtx.poolCount.Lock()
	defer c.mtx.poolCount.Unlock()
	c.poolCount = count
}

func (c *Client) clearRequestCount() {
	c.mtx.requestCount.Lock()
	defer c.mtx.requestCount.Unlock()
	c.requestCount = 0
}

func (c *Client) resetRequestCountAfterOneSecond() {

	for {
		select {
		case <-time.After(time.Second):
			c.clearRequestCount()
		}
	}
}

func (c *Client) getPoolCount() int {
	c.mtx.poolCount.Lock()
	defer c.mtx.poolCount.Unlock()
	return c.poolCount
}

func (c *Client) getActualPoolCount() int {
	c.mtx.actualPoolCount.Lock()
	defer c.mtx.actualPoolCount.Unlock()
	return c.actualPoolCount
}

func (c *Client) setRequest(work Work) {
	defer c.incrementRequestCount()
	c.request <- work
}

func (c *Client) getRequestCount() int {
	c.mtx.requestCount.Lock()
	defer c.mtx.requestCount.Unlock()
	return c.requestCount
}

func (c *Client) incrementRequestCount() {
	c.mtx.requestCount.Lock()
	defer c.mtx.requestCount.Unlock()
	c.requestCount++
}

func (c *Client) decrementRequestCount() {
	c.mtx.requestCount.Lock()
	defer c.mtx.requestCount.Unlock()
	c.requestCount--
}

func (c *Client) getResult(id uuid.UUID) (Work, bool) {
	c.mtx.result.Lock()
	defer c.mtx.result.Unlock()

	result, ok := c.result[id]
	return result, ok
}

func (c *Client) setResult(id uuid.UUID, result Work) {
	c.mtx.result.Lock()
	defer c.mtx.result.Unlock()
	c.result[id] = result
}

func (c *Client) deleteResult(id uuid.UUID) {
	c.mtx.result.Lock()
	defer c.mtx.result.Unlock()
	delete(c.result, id)
}

func (c *Client) incrementActualPoolCount(num int) {
	c.mtx.actualPoolCount.Lock()
	defer c.mtx.actualPoolCount.Unlock()
	c.actualPoolCount += num
}

func (c *Client) decrementActualPoolCount(num int) {
	c.mtx.actualPoolCount.Lock()
	defer c.mtx.actualPoolCount.Unlock()
	c.actualPoolCount -= num
}

const (
	ShortenURL  FuncType = "shorten_url"
	HealthCheck FuncType = "health_check"
)

func (c *Client) exec(raw Work) {

	switch raw.fnType {

	case HealthCheck:
		raw.Output, raw.Error = c.workerFunctions.CheckHealth()
		c.setResult(raw.ID, raw)
	case ShortenURL:
		raw.Output, raw.Error = c.workerFunctions.ShortenURL(raw.Input.(ShortenURLRequest))
		//defer c.decrementRequestCount()
		c.setResult(raw.ID, raw)
	}
}

func (c *Client) stop(num int) {
	for i := 0; i < num; i++ {
		c.stopSignal <- true
	}
}

func (c *Client) worker() {

	// shutdown
	for raw := range c.request {
		select {
		case <-c.stopSignal:
			return
		default:
			c.exec(raw)
		}
	}
}

// monitor runs every 5 seconds to check if the number of workers is correct
// if not, it will start or stop workers
func (c *Client) monitor() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:

			if c.getActualPoolCount() < c.getPoolCount() {
				go c.worker()
				c.incrementActualPoolCount(1)
			} else if len(c.request) > c.poolCount {
				c.stop(1)
				c.decrementActualPoolCount(1)
			}
		}
	}
}

func (c *Client) start() {
	for i := 0; i < c.getPoolCount(); i++ {
		go c.worker()
		c.incrementActualPoolCount(1)
	}

	go c.resetRequestCountAfterOneSecond()

	go c.monitor()

}

//
//
//			Worker Functions
//
//

func (c *Client) ShortenURL(request ShortenURLRequest) (*URLCollection, error) {

	workRequest := Work{
		ID:     uuid.New(),
		Input:  request,
		fnType: ShortenURL,
	}

	c.setRequest(workRequest)

	// wait for result
	for {
		select {
		case <-time.After(time.Millisecond * 10):

			if result, ok := c.getResult(workRequest.ID); ok {

				if result.Error != nil {
					return nil, result.Error
				}

				c.deleteResult(workRequest.ID)

				url := result.Output.(*URLCollection)
				return url, nil
			}
		}

	}

}

func (c *Client) CheckHealth() (bool, error) {

	workRequest := Work{
		ID:     uuid.New(),
		fnType: HealthCheck,
	}

	c.setRequest(workRequest)

	// wait for result
	for {
		select {
		case <-time.After(time.Millisecond * 10):

			if result, ok := c.getResult(workRequest.ID); ok {

				if result.Error != nil {
					return false, result.Error
				}

				c.deleteResult(workRequest.ID)

				return true, nil
			}
		}

	}
}
