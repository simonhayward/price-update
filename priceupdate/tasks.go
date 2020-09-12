package priceupdate

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const SourcesTimeout = 15 * time.Second

func MustPopulateSecuritiesAndRows(sourcesUrl, storeUrl string) (*Securities, Rows) {
	errCh := make(chan error)
	doneCh := make(chan bool)
	defer close(errCh)

	var securities *Securities
	var rows Rows

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err error
		securities, err = GetSecurities(sourcesUrl)
		if err != nil {
			errCh <- fmt.Errorf("securities failed to load %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		rows, err = GetRows(storeUrl)
		if err != nil {
			errCh <- fmt.Errorf("rows failed to load %s", err)
		}
	}()

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
	case err := <-errCh:
		panic(fmt.Sprintf("setup failed: %s", err))
	}

	return securities, rows
}

func GetISINSourcesPrice(security *Security, rows Rows, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	newPriceCh := make(chan *Source)
	tasksCh := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_wg := sync.WaitGroup{}
		_wg.Add(len(security.Sources))
		for i := 0; i < len(security.Sources); i++ {
			go getISINPrice(ctx, &security.Sources[i], newPriceCh, &_wg)
		}
		_wg.Wait()
		close(tasksCh)
	}()

	select {
	case source := <-newPriceCh:
		LogOutput(Log{Message: fmt.Sprintf("price found: %s at: %s", source.Price, source.URL), Severity: Info})
		// Lock access to map
		mutex.Lock()
		var index int
		if row, ok := rows[security.ISIN]; ok {
			index = row.Index
		} else {
			// new row so append
			index = len(rows)
		}
		rows[security.ISIN] = &Row{
			Index:   index,
			Price:   source.Price,
			Updated: fmt.Sprintf("%s", time.Now().UTC()),
		}
		mutex.Unlock()
	case <-tasksCh:
		LogOutput(Log{Message: fmt.Sprintf("no price: %s", security.ISIN), Severity: Error})
	case <-time.After(SourcesTimeout):
		LogOutput(Log{Message: fmt.Sprintf("timed out for: %s after: %d seconds", security.ISIN, SourcesTimeout), Severity: Error})
	}
	cancel() // cancel running tasks
}

func getISINPrice(ctx context.Context, s *Source, newPriceCh chan *Source, wgPrices *sync.WaitGroup) {
	defer wgPrices.Done()

	doneCh := make(chan struct{})
	go func() {
		if err := s.SetPrice(); err == nil {
			newPriceCh <- s
		} else {
			LogOutput(Log{Message: fmt.Sprintf("price error: %s", err), Severity: Error})
		}

		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		LogOutput(Log{Message: fmt.Sprintf("cancelling: %s", s.URL), Severity: Info})
		return
	case <-doneCh:
		return
	}
}
