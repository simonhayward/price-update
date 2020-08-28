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

	newPriceCh := make(chan string)
	tasksCh := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_wg := sync.WaitGroup{}
		_wg.Add(len(security.Sources))
		for i := 0; i < len(security.Sources); i++ {
			go getISINPrice(ctx, cancel, &security.Sources[i], newPriceCh, &_wg)
		}
		_wg.Wait()
		close(tasksCh)
	}()

	select {
	case price := <-newPriceCh:
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
			Price:   price,
			Updated: fmt.Sprintf("%s", time.Now().UTC()),
		}
		mutex.Unlock()
	case <-tasksCh:
		fmt.Println(fmt.Sprintf(`{"message": "no price: %s", "severity": "error"}`, security.ISIN))
	case <-time.After(SourcesTimeout):
		fmt.Println(fmt.Sprintf(`{"message": "timed out for: %s after: %d seconds", "severity": "info"}`, security.ISIN, SourcesTimeout))
		cancel() // cancel running tasks
	}
}

func getISINPrice(ctx context.Context, cancel context.CancelFunc, s *Source, newPriceCh chan string, wgPrices *sync.WaitGroup) {
	defer wgPrices.Done()

	doneCh := make(chan struct{})
	go func() {
		if err := s.SetPrice(); err == nil {
			fmt.Println(fmt.Sprintf(`{"message": "price found: %s at: %s", "severity": "info"}`, s.Price, s.URL))
			newPriceCh <- s.Price
			cancel() // cancel other tasks
		} else {
			fmt.Println(fmt.Sprintf(`{"message": "price error: %s", "severity": "error"}`, err))
		}

		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		fmt.Println(fmt.Sprintf(`{"message": "cancelling: %s", "severity": "info"}`, s.URL))
		return
	case <-doneCh:
		return
	}
}
