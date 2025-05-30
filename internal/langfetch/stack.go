package langfetch

import "fmt"

type cmdSpec struct {
	Label string
	Bin   string
	Args  []string
}

// stackFetcher builds a fetchFn that queries each component command and
// aggregates their version strings into LangInfo.Details.
func stackFetcher(name string, specs []cmdSpec) fetchFn {
	return func() (LangInfo, error) {
		var details []string
		for _, sp := range specs {
			v := runSilently(sp.Bin, sp.Args...)
			if v != "" {
				details = append(details, fmt.Sprintf("%s: %s", sp.Label, v))
			}
		}
		if len(details) == 0 {
			return LangInfo{}, fmt.Errorf("%s components not found", name)
		}
		return LangInfo{Name: name, Details: details}, nil
	}
}
