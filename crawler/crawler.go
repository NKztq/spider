// crawler.go - crawler.

package crawler

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/baidu/go-lib/log"
	"golang.org/x/net/html"

	"github.com/NKztq/spider/conf"
	"github.com/NKztq/spider/parser"
)

var (
	taskQueueLength = 10000 // length of Crawler's task queue
)

// token bucket of length one
type hostTokenBucket chan bool

type Fetcher interface {
	// Fetch body of URL.
	Fetch(url string) ([]byte, error)
}

type Outputer interface {
	// Output content to file.
	OutputFile(fileName string, content []byte) error
}

type task struct {
	url   *url.URL
	depth int
}

type Crawler struct {
	// crawler cfg
	maxDepth      int // max crawling depth
	crawlInterval int // crawl interval, in seconds
	threadCount   int // crawling thread limit

	seeds []string // urls

	fetchedURL sync.Map

	// TODO: Use go-lib/queue instead of WaitGroup&chan as task manager
	taskManager *sync.WaitGroup // task manager
	tasks       chan *task      // task queue

	frequencyLimiter sync.Map // host => host lock

	// proxy
	fetcher  Fetcher  // fetcher for crawler
	outputer Outputer // outputer for crawler
}

func NewCrawler(cfg conf.CrawlerConf, seeds []string, fetcher Fetcher, outputer Outputer) *Crawler {
	return &Crawler{
		maxDepth:      cfg.MaxDepth,
		crawlInterval: cfg.CrawlInterval,
		threadCount:   cfg.ThreadCount,
		seeds:         seeds,
		taskManager:   &sync.WaitGroup{},
		tasks:         make(chan *task, taskQueueLength),
		fetcher:       fetcher,
		outputer:      outputer,
	}
}

// Run crawler once.
func (c *Crawler) RunOnce() error {
	if c.maxDepth < 0 {
		return fmt.Errorf("maxDepth should >= 0, but got: %d", c.maxDepth)
	}

	c.initTasks()

	// crawl
	for i := 0; i < c.threadCount; i++ {
		go c.crawl()
	}

	c.taskManager.Wait()

	return nil
}

// Add seeds to tasks.
func (c *Crawler) initTasks() {
	// parsed all seeds into *url.URL, filter out invalid ones
	validSeeds := []*url.URL{}
	for _, seed := range c.seeds {
		parsedURL, err := url.Parse(seed)
		if err != nil {
			log.Logger.Error("initTasks(): url: %s, url.Parse(): %v", seed, err)
			continue
		}

		validSeeds = append(validSeeds, parsedURL)
	}

	var syncSeeds []*url.URL  // save seeds should add to tasks synchronously
	var asyncSeeds []*url.URL // save seeds should add to tasks asynchronously

	if len(validSeeds) > taskQueueLength {
		syncSeeds = validSeeds[:taskQueueLength]
		asyncSeeds = validSeeds[taskQueueLength:]
		log.Logger.Warn("initTasks(): length of valid seeds(%d) > taskQueueLength(%d), length of seeds: %d", len(validSeeds), taskQueueLength, len(c.seeds))
	} else {
		syncSeeds = validSeeds
	}

	c.taskManager.Add(len(validSeeds))

	// add syncSeeds synchronously
	for _, u := range syncSeeds {
		c.addTask(task{u, c.maxDepth})
	}

	// add asyncSeeds asynchronously
	if len(asyncSeeds) > 0 {
		go c.addTasks(asyncSeeds, c.maxDepth)
	}
}

// Productor for c.tasks queue, add one task.
func (c *Crawler) addTask(task task) {
	c.tasks <- &task
}

// Productor for c.tasks queue, add mutiple tasks.
func (c *Crawler) addTasks(urls []*url.URL, depth int) {
	for _, url := range urls {
		c.addTask(task{url, depth})
	}
}

// Both consumer and productor for c.tasks queue.
// As consumer, deals task in task queue.
// As productor, adds further tasks to task queue asynchronously.
func (c *Crawler) crawl() {
	for {
		t := <-c.tasks
		u := t.url
		uStr := u.String()

		log.Logger.Info("crawl(): start crawling %s", uStr)

		c.limitFrequency(u.Host)

		fetchRes, err := c.fetcher.Fetch(uStr)
		if err != nil {
			log.Logger.Error("crawl(): fetch url failed: %s, fetcher.Fetch(): %v", uStr, err)
			c.taskManager.Done()
			continue
		} else {
			// output to file
			err := c.outputer.OutputFile(url.QueryEscape(uStr), fetchRes)
			if err != nil {
				log.Logger.Warn("crawl(): write url: %s to file failed, outputer.Output(): %v", uStr, err)
			}
		}

		// parse html
		r := bytes.NewReader(fetchRes)
		node, err := html.Parse(r)
		if err != nil {
			log.Logger.Error("crawl(): url: %s, html.Parse(): %v", t.url, err)
		}

		// get deeper URLs
		deeperURLs := []*url.URL{}
		parser.Parse(node, u, &deeperURLs)

		// add further tasks
		depth := t.depth - 1
		if depth >= 0 && len(deeperURLs) > 0 {
			for _, url := range deeperURLs {
				furtherTask := task{url, depth}
				if _, exist := c.fetchedURL.LoadOrStore(url.String(), true); !exist {
					c.taskManager.Add(1)
					go c.addTask(furtherTask)
				}
			}
		}

		c.taskManager.Done()
	}
}

// TODO: Optimize efficiency for limitFrequency(), at present, c.crawl() will hang out when c.limitFrequency() failed.
// Limit fetch frequency for host by tokenBucket.
func (c *Crawler) limitFrequency(host string) {
	t, ok := c.frequencyLimiter.LoadOrStore(host, make(hostTokenBucket))
	if !ok {
		// issue a token for host circularly by c.crawlInterval
		go issueTokenCircularly(t.(hostTokenBucket), c.crawlInterval)
	}

	// get token
	getToken(t.(hostTokenBucket))
}

func issueTokenCircularly(tb hostTokenBucket, interval int) {
	for {
		tb <- true
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func getToken(tb hostTokenBucket) {
	<-tb
}
