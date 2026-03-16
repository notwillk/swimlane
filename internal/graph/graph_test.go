package graph

import (
	"testing"

	"github.com/notwillk/swimlane/internal/filter"
	"github.com/notwillk/swimlane/internal/ticket"
)

func mk(ulid, priority, status string, ready bool, blockedBy, subtasks []string) *ticket.Ticket {
	return &ticket.Ticket{
		ULID:      ulid,
		Priority:  priority,
		Status:    status,
		Ready:     ready,
		BlockedBy: blockedBy,
		Subtasks:  subtasks,
		Path:      "/tickets/" + ulid + "-slug.md",
	}
}

func TestIsBlocked_blockedBy(t *testing.T) {
	// A blocked by B; B not done -> A blocked
	a := mk("01A", "p1", "todo", true, []string{"01B"}, nil)
	b := mk("01B", "p2", "todo", true, nil, nil)
	g := Build([]*ticket.Ticket{a, b})
	if !g.IsBlocked(a) {
		t.Error("A should be blocked by B")
	}
	b.Status = "done"
	g = Build([]*ticket.Ticket{a, b})
	if g.IsBlocked(a) {
		t.Error("A should not be blocked when B is done")
	}
}

func TestIsBlocked_subtasks(t *testing.T) {
	parent := mk("01P", "p1", "todo", true, nil, []string{"01S1", "01S2"})
	sub1 := mk("01S1", "p1", "todo", true, nil, nil)
	sub2 := mk("01S2", "p1", "todo", true, nil, nil)
	g := Build([]*ticket.Ticket{parent, sub1, sub2})
	if !g.IsBlocked(parent) {
		t.Error("parent should be blocked when subtasks are not done")
	}
	sub1.Status = "done"
	g = Build([]*ticket.Ticket{parent, sub1, sub2})
	if !g.IsBlocked(parent) {
		t.Error("parent should still be blocked when one subtask remains todo")
	}
	sub2.Status = "done"
	g = Build([]*ticket.Ticket{parent, sub1, sub2})
	if g.IsBlocked(parent) {
		t.Error("parent should not be blocked when all subtasks are done")
	}
}

func TestIsBlocked_subtaskNotInGraph(t *testing.T) {
	// Parent lists subtask ULID that is not in the graph -> not blocked (treat as done/out-of-scope)
	parent := mk("01P", "p1", "todo", true, nil, []string{"01MISSING"})
	g := Build([]*ticket.Ticket{parent})
	if g.IsBlocked(parent) {
		t.Error("parent should not be blocked when subtask is not in graph")
	}
}

func TestIsBlocked_noSubtasks(t *testing.T) {
	tick := mk("01X", "p1", "todo", true, nil, nil)
	g := Build([]*ticket.Ticket{tick})
	if g.IsBlocked(tick) {
		t.Error("ticket with no deps should not be blocked")
	}
}

func TestNext_returnsSubtaskWhenParentBlockedBySubtasks(t *testing.T) {
	parent := mk("01P", "p1", "todo", true, nil, []string{"01S"})
	sub := mk("01S", "p1", "todo", true, nil, nil)
	g := Build([]*ticket.Ticket{parent, sub})
	path := g.Next(nil)
	if path == "" {
		t.Fatal("Next should return something")
	}
	if path != sub.Path {
		t.Errorf("Next should return subtask path, got %q", path)
	}
}

func TestNext_returnsParentWhenAllSubtasksDone(t *testing.T) {
	parent := mk("01P", "p1", "todo", true, nil, []string{"01S"})
	sub := mk("01S", "p1", "done", true, nil, nil)
	g := Build([]*ticket.Ticket{parent, sub})
	path := g.Next(nil)
	if path == "" {
		t.Fatal("Next should return something")
	}
	if path != parent.Path {
		t.Errorf("Next should return parent path when subtasks done, got %q", path)
	}
}

func TestNext_priorityOrderEligible(t *testing.T) {
	lo := mk("01L", "p2", "todo", true, nil, nil)
	hi := mk("01H", "p1", "todo", true, nil, nil)
	g := Build([]*ticket.Ticket{lo, hi})
	path := g.Next(nil)
	if path != hi.Path {
		t.Errorf("Next should return higher priority (p1), got %q", path)
	}
}

func TestNext_filterApplied(t *testing.T) {
	a := mk("01A", "p1", "todo", true, nil, nil)
	b := mk("01B", "p2", "todo", true, nil, nil)
	g := Build([]*ticket.Ticket{a, b})
	f := filter.FromFlags([]string{"p2"}, nil, nil, nil)
	path := g.Next(f)
	if path != b.Path {
		t.Errorf("Next with p2 filter should return B, got %q", path)
	}
}

func TestNext_emptyWhenNoEligible(t *testing.T) {
	// All done or not ready
	a := mk("01A", "p1", "done", true, nil, nil)
	b := mk("01B", "p1", "todo", false, nil, nil)
	g := Build([]*ticket.Ticket{a, b})
	path := g.Next(nil)
	if path != "" {
		t.Errorf("Next should return empty when no eligible, got %q", path)
	}
}
