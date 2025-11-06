package index

import (
	"strings"
	"sync"

	erwp "github.com/PerumallaGiridhar/oolio/internal/errorwrap"
)

type PebbleIndex struct {
	Stores []*PebbleStore
}

func NewPebbleIndex(paths []string) *PebbleIndex {
	stores := make([]*PebbleStore, len(paths))
	opened := make([]*PebbleStore, len(paths))

	var mu sync.Mutex
	sem := make(chan struct{}, 4)
	wg := sync.WaitGroup{}

	closeAll := func() {
		mu.Lock()
		defer mu.Unlock()
		for _, s := range opened {
			_ = s.Close()
		}
		opened = opened[:0]
	}

	for i, p := range paths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			ps := EnsurePebble(strings.TrimSpace(p), closeAll)
			mu.Lock()
			defer mu.Unlock()
			stores[i] = ps
			opened = append(opened, ps)
		}()
	}
	wg.Wait()
	return &PebbleIndex{Stores: stores}
}

func (pi *PebbleIndex) Close() {
	for _, s := range pi.Stores {
		_ = s.Close()
	}
}

func (pi *PebbleIndex) IsValid2of3(code string) bool {
	hits := 0
	for _, s := range pi.Stores {
		ok := erwp.LetReturn(erwp.Try(s.Has(code)))
		if ok {
			hits++
			if hits >= 2 {
				return true
			}
		}
	}
	return false
}
