package index

import (
	"fmt"
	"strings"
	"sync"
)

type PebbleIndex struct {
	Stores []*PebbleStore
}

func NewPebbleIndex(paths []string) (*PebbleIndex, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no pebble paths provided")
	}

	stores := make([]*PebbleStore, len(paths))

	var (
		mu       sync.Mutex
		opened   []*PebbleStore
		firstErr error
		once     sync.Once
	)

	closeAll := func() {
		mu.Lock()
		defer mu.Unlock()

		for _, s := range opened {
			if s != nil {
				_ = s.Close()
			}
		}
		opened = nil
	}

	setErr := func(err error) {
		if err == nil {
			return
		}
		once.Do(func() {
			firstErr = err
			closeAll()
		})
	}

	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup

	for i, rawPath := range paths {
		i, rawPath := i, strings.TrimSpace(rawPath)

		if rawPath == "" {
			return nil, fmt.Errorf("empty pebble path at index %d", i)
		}

		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			ps, err := EnsurePebble(path)
			if err != nil {
				setErr(fmt.Errorf("ensure pebble for %q: %w", path, err))
				return
			}

			mu.Lock()
			stores[idx] = ps
			opened = append(opened, ps)
			mu.Unlock()
		}(i, rawPath)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return &PebbleIndex{Stores: stores}, nil
}

func (pi *PebbleIndex) Close() {
	for _, s := range pi.Stores {
		if s != nil {
			_ = s.Close()
		}
	}
}

func (pi *PebbleIndex) IsValid2of3(code string) (bool, error) {
	hits := 0
	for _, s := range pi.Stores {
		ok, err := s.Has(code)
		if err != nil {
			return false, err
		}
		if ok {
			hits++
			if hits >= 2 {
				return true, nil
			}
		}
	}
	return false, nil
}
