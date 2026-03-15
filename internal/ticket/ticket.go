package ticket

// Ticket represents a swimlane ticket (metadata + path).
type Ticket struct {
	ULID      string
	Title     string
	Priority  string
	Status    string
	Ready     bool
	BlockedBy []string
	Tags      []string
	Path      string
}

// ValidPriorities for ticket metadata.
var ValidPriorities = []string{"p0", "p1", "p2", "p3", "p4"}

// ValidStatuses for ticket metadata.
var ValidStatuses = []string{"todo", "in-progress", "done"}

// Validate checks required fields and allowed enums. path is used in error messages.
func Validate(t *Ticket) error {
	if t.Priority == "" {
		return &ValidationError{Path: t.Path, Field: "priority", Msg: "required"}
	}
	ok := false
	for _, p := range ValidPriorities {
		if t.Priority == p {
			ok = true
			break
		}
	}
	if !ok {
		return &ValidationError{Path: t.Path, Field: "priority", Msg: "must be one of p0, p1, p2, p3, p4"}
	}
	if t.Status == "" {
		return &ValidationError{Path: t.Path, Field: "status", Msg: "required"}
	}
	ok = false
	for _, s := range ValidStatuses {
		if t.Status == s {
			ok = true
			break
		}
	}
	if !ok {
		return &ValidationError{Path: t.Path, Field: "status", Msg: "must be one of todo, in-progress, done"}
	}
	return nil
}

// ValidationError is a ticket validation error with file path.
type ValidationError struct {
	Path  string
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	if e.Path != "" {
		return e.Path + ": " + e.Field + " " + e.Msg
	}
	return e.Field + " " + e.Msg
}
