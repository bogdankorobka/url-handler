package services

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL      string `json:"url"`
	PageSize int64  `json:"page_size"`
}

type URLService struct {
	wg          sync.WaitGroup
	numWorkers  int
	urlsList    []string
	urlsChan    chan string
	resultsChan chan *Result
	errors      bool
}

func NewURLServie(numWorkers int, urlsList []string) *URLService {
	return &URLService{
		numWorkers:  numWorkers,
		urlsList:    urlsList,
		urlsChan:    make(chan string, len(urlsList)),
		resultsChan: make(chan *Result, len(urlsList)),
		errors:      false,
	}
}

func (s *URLService) Start(ctx context.Context) ([]*Result, bool) {
	results := make([]*Result, 0, len(s.urlsList))

	// start wurkers
	for i := 0; i < s.numWorkers; i++ {
		s.wg.Add(1)

		go s.requester(ctx)
	}

	// send urls to chan
	go func() {
		for _, url := range s.urlsList {
			s.urlsChan <- url
		}
		close(s.urlsChan)
	}()

	// wait workers and close chan
	s.wg.Wait()
	close(s.resultsChan)

	// return values
	if s.errors {
		return nil, false
	}

	for r := range s.resultsChan {
		results = append(results, r)
	}

	return results, true
}

func (s *URLService) requester(ctx context.Context) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			s.errors = true

			return
		case <-time.Tick(100 * time.Millisecond):
			url, ok := <-s.urlsChan
			if !ok {
				return
			}

			// send request url
			httpResp, err := http.Get(url)
			if err != nil || httpResp.StatusCode != http.StatusOK {
				// find error request
				_, cancel := context.WithCancel(ctx)
				cancel()

				return
			}

			// calculate page size
			pageSize := httpResp.ContentLength
			if pageSize == -1 {
				content, err := ioutil.ReadAll(httpResp.Body)
				defer httpResp.Body.Close()

				if err == nil {
					pageSize = int64(len(content))
				}
			}

			// return resut to chan
			s.resultsChan <- &Result{
				URL:      url,
				PageSize: pageSize,
			}
		}
	}
}
