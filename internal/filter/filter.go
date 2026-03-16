package filter

import "github.com/notwillk/swimlane/internal/ticket"

// Filter holds include/exclude values per dimension.
// Empty include means "no filter on this dimension"; non-empty means ticket must match one of include.
// Exclude: ticket must not match any of exclude.
type Filter struct {
	PriorityInclude []string
	PriorityExclude []string
	TagInclude      []string
	TagExclude      []string
	StatusInclude   []string
	StatusExclude   []string
	ReadyInclude    []string // "true" or "false"
	ReadyExclude    []string
}

// Apply returns tickets that match the filter.
func Apply(tickets []*ticket.Ticket, f *Filter) []*ticket.Ticket {
	if f == nil {
		return tickets
	}
	var out []*ticket.Ticket
	for _, t := range tickets {
		if match(t, f) {
			out = append(out, t)
		}
	}
	return out
}

func match(t *ticket.Ticket, f *Filter) bool {
	if !matchSlice([]string{t.Priority}, f.PriorityInclude, f.PriorityExclude) {
		return false
	}
	if !matchSlice(t.Tags, f.TagInclude, f.TagExclude) {
		return false
	}
	if !matchSlice([]string{t.Status}, f.StatusInclude, f.StatusExclude) {
		return false
	}
	readyStr := "false"
	if t.Ready {
		readyStr = "true"
	}
	if !matchSlice([]string{readyStr}, f.ReadyInclude, f.ReadyExclude) {
		return false
	}
	return true
}

// matchSlice returns true if:
// - include is empty or values contains one of include
// - and values does not contain any of exclude
func matchSlice(values []string, include, exclude []string) bool {
	for _, e := range exclude {
		for _, v := range values {
			if v == e {
				return false
			}
		}
	}
	if len(include) == 0 {
		return true
	}
	for _, inc := range include {
		for _, v := range values {
			if v == inc {
				return true
			}
		}
	}
	return false
}
