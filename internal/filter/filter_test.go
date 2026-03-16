package filter

import (
	"testing"

	"github.com/notwillk/swimlane/internal/ticket"
)

func TestApply_priority(t *testing.T) {
	tickets := []*ticket.Ticket{
		{ULID: "1", Priority: "p1", Status: "todo", Ready: true},
		{ULID: "2", Priority: "p2", Status: "todo", Ready: true},
	}
	f := FromFlags([]string{"p1"}, nil, nil, nil)
	out := Apply(tickets, f)
	if len(out) != 1 || out[0].ULID != "1" {
		t.Errorf("expected one ticket p1, got %v", out)
	}
}

func TestApply_negate(t *testing.T) {
	tickets := []*ticket.Ticket{
		{ULID: "1", Priority: "p0", Status: "todo", Ready: true},
		{ULID: "2", Priority: "p1", Status: "todo", Ready: true},
	}
	f := FromFlags([]string{"!p0"}, nil, nil, nil)
	out := Apply(tickets, f)
	if len(out) != 1 || out[0].ULID != "2" {
		t.Errorf("expected one ticket (not p0), got %v", out)
	}
}
