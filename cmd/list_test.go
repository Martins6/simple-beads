package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Martins6/simple-beads/internal/models"
	"github.com/Martins6/simple-beads/internal/storage"
)

func resetListFlags() {
	listOn = dateFlag{}
	listAfter = dateFlag{}
	listBefore = dateFlag{}
	listAll = false
	listPriority = -1
	listStatus = ""
	listParent = ""

	flags := listCmd.Flags()
	for _, name := range []string{"priority", "status", "parent", "all", "on", "after", "before"} {
		if f := flags.Lookup(name); f != nil {
			f.Changed = false
		}
	}
}

func setupListTest(t *testing.T) (cleanup func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "tasks.db")
	testStore := storage.NewStorageWithPath(dbPath)
	if err := testStore.Init(); err != nil {
		t.Fatalf("init storage: %v", err)
	}

	origStore := store
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	store = testStore
	resetListFlags()

	cleanup = func() {
		store.Close()
		store = origStore
		_ = os.Chdir(origCwd)
		resetListFlags()
	}
	return cleanup
}

func runList(t *testing.T, args ...string) (stdout string, err error) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	origStdout := os.Stdout
	os.Stdout = w

	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	rootCmd.SetArgs(append([]string{"list"}, args...))
	execErr := rootCmd.Execute()

	w.Close()
	<-done
	os.Stdout = origStdout
	_ = r.Close()

	return buf.String(), execErr
}

func mustParse(t *testing.T, s string) time.Time {
	t.Helper()
	loc, err := time.LoadLocation("Local")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	tm, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		t.Fatalf("parse %q: %v", s, err)
	}
	return tm
}

func saveTaskWithTimes(t *testing.T, title string, status models.Status, createdAt, closedAt time.Time) *models.Task {
	t.Helper()
	task := models.NewTask(title)
	task.CreatedAt = createdAt
	task.UpdatedAt = createdAt
	if status == models.StatusClosed {
		task.ClosedAt = closedAt
		task.UpdatedAt = closedAt
	}
	task.Status = status
	if err := store.SaveTask(task); err != nil {
		t.Fatalf("save task: %v", err)
	}
	return task
}

func TestListDateFilters_OnExactDate(t *testing.T) {
	defer setupListTest(t)()

	day := mustParse(t, "2026-08-03")
	other := mustParse(t, "2026-08-01")

	closedOn := saveTaskWithTimes(t, "Closed on day", models.StatusClosed, other, day)
	openOn := saveTaskWithTimes(t, "Open created on day", models.StatusOpen, day, time.Time{})
	openOff := saveTaskWithTimes(t, "Open off day", models.StatusOpen, other, time.Time{})
	closedOff := saveTaskWithTimes(t, "Closed off day", models.StatusClosed, other, other)

	tests := []struct {
		name    string
		args    []string
		wantIDs []string
		wantErr bool
	}{
		{
			name:    "all + on day matches both open and closed on that day",
			args:    []string{"--all", "--on", "2026-08-03"},
			wantIDs: []string{closedOn.ID, openOn.ID},
		},
		{
			name:    "closed + on day filters by closed_at",
			args:    []string{"--all", "--status", "closed", "--on", "2026-08-03"},
			wantIDs: []string{closedOn.ID},
		},
		{
			name:    "open + on day filters by created_at",
			args:    []string{"--on", "2026-08-03"},
			wantIDs: []string{openOn.ID},
		},
		{
			name:    "off-day tasks excluded",
			args:    []string{"--all", "--on", "2026-08-03"},
			wantIDs: []string{closedOn.ID, openOn.ID},
		},
		{
			name:    "on day with no matches",
			args:    []string{"--all", "--on", "2026-01-01"},
			wantIDs: nil,
		},
		{
			name:    "invalid date format",
			args:    []string{"--all", "--on", "not-a-date"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listOn = dateFlag{}
			listAfter = dateFlag{}
			listBefore = dateFlag{}
			listAll = false
			listStatus = ""

			out, err := runList(t, tt.args...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil; output=%q", out)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v; output=%q", err, out)
			}
			gotIDs := extractIDs(out)
			if !equalSlices(gotIDs, tt.wantIDs) {
				t.Errorf("ids mismatch: got %v, want %v (output=%q)", gotIDs, tt.wantIDs, out)
			}
		})
	}

	_ = openOff
	_ = closedOff
}

func TestListDateFilters_AfterBefore(t *testing.T) {
	defer setupListTest(t)()

	d1 := mustParse(t, "2026-08-01")
	d2 := mustParse(t, "2026-08-02")
	d3 := mustParse(t, "2026-08-03")
	d4 := mustParse(t, "2026-08-04")

	openD1 := saveTaskWithTimes(t, "Open d1", models.StatusOpen, d1, time.Time{})
	openD2 := saveTaskWithTimes(t, "Open d2", models.StatusOpen, d2, time.Time{})
	openD3 := saveTaskWithTimes(t, "Open d3", models.StatusOpen, d3, time.Time{})
	closedD1 := saveTaskWithTimes(t, "Closed d1", models.StatusClosed, d1, d1)
	closedD3 := saveTaskWithTimes(t, "Closed d3", models.StatusClosed, d2, d3)
	closedD4 := saveTaskWithTimes(t, "Closed d4", models.StatusClosed, d3, d4)

	tests := []struct {
		name    string
		args    []string
		wantIDs []string
	}{
		{
			name:    "after d2 inclusive matches d2 d3 d4",
			args:    []string{"--all", "--after", "2026-08-02"},
			wantIDs: []string{openD2.ID, openD3.ID, closedD3.ID, closedD4.ID},
		},
		{
			name:    "before d2 inclusive matches d1 d2",
			args:    []string{"--all", "--before", "2026-08-02"},
			wantIDs: []string{openD1.ID, openD2.ID, closedD1.ID},
		},
		{
			name:    "after d1 and before d3 matches d1 d2 d3",
			args:    []string{"--all", "--after", "2026-08-01", "--before", "2026-08-03"},
			wantIDs: []string{openD1.ID, openD2.ID, openD3.ID, closedD1.ID, closedD3.ID},
		},
		{
			name:    "after and before same day matches that day",
			args:    []string{"--all", "--after", "2026-08-02", "--before", "2026-08-02"},
			wantIDs: []string{openD2.ID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listOn = dateFlag{}
			listAfter = dateFlag{}
			listBefore = dateFlag{}
			listAll = false
			listStatus = ""

			out, err := runList(t, tt.args...)
			if err != nil {
				t.Fatalf("unexpected error: %v; output=%q", err, out)
			}
			gotIDs := extractIDs(out)
			if !equalSlices(gotIDs, tt.wantIDs) {
				t.Errorf("ids mismatch: got %v, want %v (output=%q)", gotIDs, tt.wantIDs, out)
			}
		})
	}
}

func TestListDateFilters_MutualExclusivity(t *testing.T) {
	defer setupListTest(t)()

	d1 := mustParse(t, "2026-08-01")
	d3 := mustParse(t, "2026-08-03")
	saveTaskWithTimes(t, "Task", models.StatusOpen, d1, time.Time{})
	saveTaskWithTimes(t, "Closed", models.StatusClosed, d1, d3)

	out, err := runList(t, "--all", "--on", "2026-08-02", "--after", "2026-08-01")
	if err == nil {
		t.Fatalf("expected error from --on + --after, got nil; output=%q", out)
	}
	if !strings.Contains(err.Error(), "--on cannot be combined") {
		t.Errorf("expected mutual-exclusivity error, got: %v", err)
	}
}

func TestListDateFilters_CombinesWithPriority(t *testing.T) {
	defer setupListTest(t)()

	day := mustParse(t, "2026-08-03")
	other := mustParse(t, "2026-08-01")

	p0 := saveTaskWithTimes(t, "P0 on day", models.StatusOpen, day, time.Time{})
	p0.Priority = models.PriorityP0
	if err := store.UpdateTask(p0); err != nil {
		t.Fatalf("update p0: %v", err)
	}
	saveTaskWithTimes(t, "P2 on day", models.StatusOpen, day, time.Time{})
	saveTaskWithTimes(t, "P0 off day", models.StatusOpen, other, time.Time{})

	out, err := runList(t, "--all", "--on", "2026-08-03", "-p", "0")
	if err != nil {
		t.Fatalf("unexpected error: %v; output=%q", err, out)
	}
	gotIDs := extractIDs(out)
	if !equalSlices(gotIDs, []string{p0.ID}) {
		t.Errorf("ids mismatch: got %v, want [%s] (output=%q)", gotIDs, p0.ID, out)
	}
}

func TestListDateFilters_DayBoundaryLocalTime(t *testing.T) {
	defer setupListTest(t)()

	loc, err := time.LoadLocation("Local")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	dayStart := time.Date(2026, 8, 3, 0, 0, 0, 0, loc)
	dayMid := time.Date(2026, 8, 3, 23, 59, 59, 999999999, loc)
	dayEndExclusive := time.Date(2026, 8, 4, 0, 0, 0, 0, loc)
	justAfter := time.Date(2026, 8, 4, 0, 0, 0, 1, loc)

	openMid := saveTaskWithTimes(t, "Open mid", models.StatusOpen, dayMid, time.Time{})
	openStart := saveTaskWithTimes(t, "Open start", models.StatusOpen, dayStart, time.Time{})
	openAfter := saveTaskWithTimes(t, "Open just after", models.StatusOpen, justAfter, time.Time{})
	_ = dayEndExclusive

	out, err := runList(t, "--all", "--on", "2026-08-03")
	if err != nil {
		t.Fatalf("unexpected error: %v; output=%q", err, out)
	}
	gotIDs := extractIDs(out)
	want := []string{openStart.ID, openMid.ID}
	if !equalSlices(gotIDs, want) {
		t.Errorf("ids mismatch: got %v, want %v (output=%q)", gotIDs, want, out)
	}
	_ = openAfter
}

func extractIDs(out string) []string {
	var ids []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ID") || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "No tasks found.") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			ids = append(ids, fields[0])
		}
	}
	return ids
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[string]int)
	for _, s := range a {
		seen[s]++
	}
	for _, s := range b {
		seen[s]--
		if seen[s] < 0 {
			return false
		}
	}
	return true
}
