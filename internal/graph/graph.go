package graph

import (
	"github.com/notwillk/swimlane/internal/filter"
	"github.com/notwillk/swimlane/internal/ticket"
)

var priorityOrder = []string{"p0", "p1", "p2", "p3", "p4"}

// Graph holds tickets indexed by ULID for next selection.
type Graph struct {
	ByULID  map[string]*ticket.Ticket
	Tickets []*ticket.Ticket
}

// Build builds a graph from tickets (must be already filtered if filter was applied).
func Build(tickets []*ticket.Ticket) *Graph {
	byULID := make(map[string]*ticket.Ticket, len(tickets))
	for _, t := range tickets {
		byULID[t.ULID] = t
	}
	return &Graph{ByULID: byULID, Tickets: tickets}
}

// IsBlocked returns true if t has any blocker (blocked_by or subtask) that is not done.
func (g *Graph) IsBlocked(t *ticket.Ticket) bool {
	for _, ulid := range t.BlockedBy {
		if dep, ok := g.ByULID[ulid]; ok && dep.Status != "done" {
			return true
		}
	}
	for _, ulid := range t.Subtasks {
		if sub, ok := g.ByULID[ulid]; ok && sub.Status != "done" {
			return true
		}
	}
	return false
}

// Next returns the path of the next ticket to implement, or empty string if none.
// Eligibility: ready=true, status=todo, not blocked. Apply f before selection.
func (g *Graph) Next(f *filter.Filter) string {
	tickets := g.Tickets
	if f != nil {
		tickets = filter.Apply(tickets, f)
	}
	// Restrict to eligible: ready, todo, not blocked
	var eligible []*ticket.Ticket
	for _, t := range tickets {
		if t.Ready && t.Status == "todo" && !g.IsBlocked(t) {
			eligible = append(eligible, t)
		}
	}

	for _, p := range priorityOrder {
		// Step 1: at this priority, any eligible?
		var atPriority []*ticket.Ticket
		for _, t := range eligible {
			if t.Priority == p {
				atPriority = append(atPriority, t)
			}
		}
		if len(atPriority) > 0 {
			// return earliest ULID (min string = earliest time)
			best := atPriority[0]
			for _, t := range atPriority[1:] {
				if t.ULID < best.ULID {
					best = t
				}
			}
			return best.Path
		}
		// Step 2: at this priority, any blocked? then return highest-priority blocker
		var blockedAtPriority []*ticket.Ticket
		for _, t := range tickets {
			if t.Priority == p && t.Ready && t.Status == "todo" && g.IsBlocked(t) {
				blockedAtPriority = append(blockedAtPriority, t)
			}
		}
		if len(blockedAtPriority) == 0 {
			continue
		}
		// collect blockers (blocked_by and subtasks that are not done)
		blockerSet := make(map[string]*ticket.Ticket)
		for _, t := range blockedAtPriority {
			for _, ulid := range t.BlockedBy {
				if dep, ok := g.ByULID[ulid]; ok && dep.Status != "done" {
					blockerSet[ulid] = dep
				}
			}
			for _, ulid := range t.Subtasks {
				if dep, ok := g.ByULID[ulid]; ok && dep.Status != "done" {
					blockerSet[ulid] = dep
				}
			}
		}
		if len(blockerSet) == 0 {
			continue
		}
		// return highest-priority blocker (by priority order), then earliest ULID
		var blockers []*ticket.Ticket
		for _, t := range blockerSet {
			blockers = append(blockers, t)
		}
		for _, p2 := range priorityOrder {
			var atP []*ticket.Ticket
			for _, t := range blockers {
				if t.Priority == p2 {
					atP = append(atP, t)
				}
			}
			if len(atP) > 0 {
				best := atP[0]
				for _, t := range atP[1:] {
					if t.ULID < best.ULID {
						best = t
					}
				}
				return best.Path
			}
		}
	}
	return ""
}
