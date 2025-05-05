package langfetch

import (
	"errors"
	"strings"
	"sync"
)

// Result bundles the outcome of a fetcher.
// If Err != nil, Info is zero.

type Result struct {
    Key  string    `json:"key"`
    Info LangInfo  `json:"info,omitempty"`
    Err  error     `json:"error,omitempty"`
}

// FetchMany runs fetchers *in parallel* with a bounded worker pool.
func FetchMany(keys []string) []Result {
    if len(keys) == 0 {
        return nil
    }

    const maxPar = 8
    sem := make(chan struct{}, maxPar)
    var wg sync.WaitGroup

    results := make([]Result, len(keys))

    for i, k := range keys {
        wg.Add(1)
        go func(i int, key string) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            fn, ok := registry[strings.ToLower(key)]
            if !ok {
                results[i] = Result{Key: key, Err: errors.New("unsupported item")}
                return
            }
            li, err := fn()
            results[i] = Result{Key: key, Info: li, Err: err}
        }(i, k)
    }

    wg.Wait()
    return results
}
