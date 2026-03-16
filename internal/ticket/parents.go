package ticket

import (
	"fmt"
	"os"
	"strings"

	"github.com/notwillk/swimlane/internal/config"
)

// CheckParentsWhenSubtaskDone runs after a ticket is marked done (or deleted).
// It finds parents that list completedULID in subtasks; if all of a parent's
// subtasks are done, it applies the configured close_parent_when_subtasks_done policy.
func CheckParentsWhenSubtaskDone(cfg *config.Config, completedULID string) error {
	policy := cfg.CloseParentWhenSubtasksDone
	if policy == "" {
		policy = config.CloseParentNever
	}
	if policy == config.CloseParentNever {
		return nil
	}

	tickets, err := Discover(cfg)
	if err != nil {
		return err
	}
	byULID := make(map[string]*Ticket)
	for _, t := range tickets {
		byULID[t.ULID] = t
	}

	for _, t := range tickets {
		if len(t.Subtasks) == 0 {
			continue
		}
		hasCompleted := false
		for _, s := range t.Subtasks {
			if s == completedULID {
				hasCompleted = true
				break
			}
		}
		if !hasCompleted {
			continue
		}

		allDone := true
		for _, subULID := range t.Subtasks {
			if subULID == completedULID {
				continue
			}
			sub, ok := byULID[subULID]
			if !ok {
				continue // deleted or missing => treat as done
			}
			if sub.Status != "done" {
				allDone = false
				break
			}
		}
		if !allDone {
			continue
		}

		// All subtasks done; apply policy
		shouldClose := false
		switch policy {
		case config.CloseParentAlways:
			shouldClose = true
		case config.CloseParentWhenEmpty:
			_, body, err := ReadFrontmatterAndBody(t.Path)
			if err != nil {
				return fmt.Errorf("parent %s: %w", t.ULID, err)
			}
			shouldClose = len(bytesTrim(body)) == 0
		case config.CloseParentWhenMatches:
			_, parentBody, err := ReadFrontmatterAndBody(t.Path)
			if err != nil {
				return fmt.Errorf("parent %s: %w", t.ULID, err)
			}
			combined, err := combinedSubtaskBodies(t, byULID)
			if err != nil {
				return fmt.Errorf("parent %s: %w", t.ULID, err)
			}
			if descriptionsMatch(parentBody, combined) {
				shouldClose = true
			} else {
				fmt.Fprintf(os.Stderr, "swimlane: parent %s not auto-closed: description does not match combined subtasks (when-matches)\n", t.ULID)
			}
		}
		if shouldClose {
			fm, body, err := ReadFrontmatterAndBody(t.Path)
			if err != nil {
				return err
			}
			fm.Status = "done"
			if err := WriteFrontmatterAndBody(t.Path, fm, body); err != nil {
				return err
			}
		}
	}
	return nil
}

func bytesTrim(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

func combinedSubtaskBodies(parent *Ticket, byULID map[string]*Ticket) ([]byte, error) {
	var parts []string
	for _, subULID := range parent.Subtasks {
		sub, ok := byULID[subULID]
		if !ok {
			continue
		}
		_, body, err := ReadFrontmatterAndBody(sub.Path)
		if err != nil {
			return nil, err
		}
		parts = append(parts, strings.TrimSpace(string(body)))
	}
	return []byte(strings.Join(parts, "\n\n")), nil
}

// descriptionsMatch normalizes both and compares. For when-matches we could use an LLM later.
func descriptionsMatch(parentBody, combinedSubtaskBody []byte) bool {
	a := normalizeBody(parentBody)
	b := normalizeBody(combinedSubtaskBody)
	return a == b
}

func normalizeBody(b []byte) string {
	s := string(b)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
