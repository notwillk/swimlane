package ticket

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notwillk/swimlane/internal/config"
)

// Valid 26-char ULIDs for test filenames.
const (
	ulidParent = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
	ulidSub1   = "01ARZ3NDEKTSV4RRFFQ69G5FAW"
	ulidSub2   = "01ARZ3NDEKTSV4RRFFQ69G5FAX"
)

func writeTicket(dir, ulid, slug, status, body string, subtasks []string) (path string, err error) {
	path = filepath.Join(dir, ulid+"-"+slug+".md")
	subtasksYaml := ""
	if len(subtasks) > 0 {
		subtasksYaml = "\nsubtasks:\n"
		for _, s := range subtasks {
			subtasksYaml += "  - " + s + "\n"
		}
	}
	content := "---\npriority: p1\nstatus: " + status + "\nready: true\n" + subtasksYaml + "---\n\n" + body
	return path, os.WriteFile(path, []byte(content), 0644)
}

func TestCheckParentsWhenSubtaskDone_never(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{Tickets: "*.md", ConfigDir: dir, CloseParentWhenSubtasksDone: config.CloseParentNever}
	_, err := writeTicket(dir, ulidParent, "parent", "todo", "parent body", []string{ulidSub1})
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub1, "sub1", "done", "sub body", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckParentsWhenSubtaskDone(cfg, ulidSub1)
	if err != nil {
		t.Fatal(err)
	}
	// Parent should still be todo
	tickets, _ := Discover(cfg)
	for _, tkt := range tickets {
		if tkt.ULID == ulidParent && tkt.Status != "todo" {
			t.Errorf("policy never: parent status should remain todo, got %q", tkt.Status)
		}
	}
}

func TestCheckParentsWhenSubtaskDone_always(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{Tickets: "*.md", ConfigDir: dir, CloseParentWhenSubtasksDone: config.CloseParentAlways}
	_, err := writeTicket(dir, ulidParent, "parent", "todo", "parent body", []string{ulidSub1})
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub1, "sub1", "done", "sub body", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckParentsWhenSubtaskDone(cfg, ulidSub1)
	if err != nil {
		t.Fatal(err)
	}
	tickets, _ := Discover(cfg)
	for _, tkt := range tickets {
		if tkt.ULID == ulidParent {
			if tkt.Status != "done" {
				t.Errorf("policy always: parent status should be done, got %q", tkt.Status)
			}
			return
		}
	}
	t.Fatal("parent ticket not found")
}

func TestCheckParentsWhenSubtaskDone_whenEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{Tickets: "*.md", ConfigDir: dir, CloseParentWhenSubtasksDone: config.CloseParentWhenEmpty}
	// Parent with empty body -> should close
	_, err := writeTicket(dir, ulidParent, "parent", "todo", "", []string{ulidSub1})
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub1, "sub1", "done", "sub body", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckParentsWhenSubtaskDone(cfg, ulidSub1)
	if err != nil {
		t.Fatal(err)
	}
	tickets, _ := Discover(cfg)
	for _, tkt := range tickets {
		if tkt.ULID == ulidParent {
			if tkt.Status != "done" {
				t.Errorf("when-empty with empty body: parent should be done, got %q", tkt.Status)
			}
			break
		}
	}
	// Parent with non-empty body -> should not close (need a second parent; use different ULID)
	ulidParent2 := "01ARZ3NDEKTSV4RRFFQ69G5FAY"
	ulidSub3 := "01ARZ3NDEKTSV4RRFFQ69G5FAZ"
	_, err = writeTicket(dir, ulidParent2, "parent2", "todo", "non-empty body", []string{ulidSub3})
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub3, "sub3", "done", "sub body", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckParentsWhenSubtaskDone(cfg, ulidSub3)
	if err != nil {
		t.Fatal(err)
	}
	tickets, _ = Discover(cfg)
	for _, tkt := range tickets {
		if tkt.ULID == ulidParent2 {
			if tkt.Status != "todo" {
				t.Errorf("when-empty with non-empty body: parent should stay todo, got %q", tkt.Status)
			}
			return
		}
	}
}

func TestCheckParentsWhenSubtaskDone_whenMatches(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{Tickets: "*.md", ConfigDir: dir, CloseParentWhenSubtasksDone: config.CloseParentWhenMatches}
	// Parent body same as combined subtask bodies (normalized) -> should close
	combined := "first part\n\nsecond part"
	_, err := writeTicket(dir, ulidParent, "parent", "todo", combined, []string{ulidSub1, ulidSub2})
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub1, "sub1", "done", "first part", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub2, "sub2", "done", "second part", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckParentsWhenSubtaskDone(cfg, ulidSub2)
	if err != nil {
		t.Fatal(err)
	}
	tickets, _ := Discover(cfg)
	for _, tkt := range tickets {
		if tkt.ULID == ulidParent {
			if tkt.Status != "done" {
				t.Errorf("when-matches with matching description: parent should be done, got %q", tkt.Status)
			}
			break
		}
	}
}

func TestCheckParentsWhenSubtaskDone_parentNotClosedUntilAllSubtasksDone(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{Tickets: "*.md", ConfigDir: dir, CloseParentWhenSubtasksDone: config.CloseParentAlways}
	_, err := writeTicket(dir, ulidParent, "parent", "todo", "", []string{ulidSub1, ulidSub2})
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub1, "sub1", "done", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = writeTicket(dir, ulidSub2, "sub2", "todo", "", nil) // not done yet
	if err != nil {
		t.Fatal(err)
	}
	err = CheckParentsWhenSubtaskDone(cfg, ulidSub1)
	if err != nil {
		t.Fatal(err)
	}
	tickets, _ := Discover(cfg)
	for _, tkt := range tickets {
		if tkt.ULID == ulidParent {
			if tkt.Status != "todo" {
				t.Errorf("parent should stay todo until all subtasks done, got %q", tkt.Status)
			}
			return
		}
	}
}
