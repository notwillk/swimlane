package filter

import "strings"

const negatePrefix = "!"

// ParseSlice parses flag values like "p1", "!p0" into include and exclude slices.
func ParseSlice(values []string) (include, exclude []string) {
	for _, v := range values {
		if strings.HasPrefix(v, negatePrefix) {
			exclude = append(exclude, strings.TrimPrefix(v, negatePrefix))
		} else {
			include = append(include, v)
		}
	}
	return include, exclude
}

// FromFlags builds a Filter from raw flag values (as returned by cobra StringSlice flags).
func FromFlags(priority, tag, status, ready []string) *Filter {
	pi, pe := ParseSlice(priority)
	ti, te := ParseSlice(tag)
	si, se := ParseSlice(status)
	ri, re := ParseSlice(ready)
	return &Filter{
		PriorityInclude: pi, PriorityExclude: pe,
		TagInclude: ti, TagExclude: te,
		StatusInclude: si, StatusExclude: se,
		ReadyInclude: ri, ReadyExclude: re,
	}
}
