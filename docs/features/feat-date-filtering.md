# Overview

Date filtering for `sb list` that lets users discover tasks by their relevant date (closed_at for closed tasks, created_at for open tasks).

# Details

- Three new flags on `sb list`: `--on`, `--after`, `--before`
- Dates accepted as ISO 8601 (`YYYY-MM-DD`) and interpreted in local time
- `--on D` matches tasks whose relevant date falls on that calendar day
- `--after D` matches tasks with relevant date >= start of D (inclusive lower bound)
- `--before D` matches tasks with relevant date < start of D+1 (inclusive upper bound, DST-safe)
- `--after` + `--before` together form a date range (both must pass)
- `--on` is mutually exclusive with `--after` / `--before`; combining returns an actionable error
- Smart-by-status rule: filters on `closed_at` for closed tasks and `created_at` for open tasks
- Bounds use a half-open interval `[start, end)` so they are immune to DST transitions
- Invalid date strings return an error naming the expected format
- Composes with existing filters: `--all`, `--status`, `-p/--priority`, `--parent`
- Implementation uses a custom `pflag.Value` (`dateFlag`) that parses via `time.ParseInLocation("2006-01-02", v, time.Local)`
- Filtering happens in Go after `GetAllTasks()`; no schema or storage changes
- First test file in `cmd/` package (`cmd/list_test.go`) covers exact-day, range, mutual-exclusivity, composition with `--priority`, and local-time day-boundary behavior

# File Paths

- cmd/list.go
- cmd/list_test.go
- README.md
