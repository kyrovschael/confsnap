# watch

The `watch` package provides lightweight file-change detection for **confsnap**.
It works by periodically hashing watched files and comparing the result against
a baseline established at startup.

## Usage

```go
import (
    "time"
    "confsnap/internal/watch"
)

w, err := watch.New([]string{"/etc/nginx/nginx.conf", "/etc/ssh/sshd_config"}, 30*time.Second)
if err != nil {
    log.Fatal(err)
}

events, err := w.Poll()
for _, ev := range events {
    fmt.Printf("changed: %s (was %s, now %s)\n", ev.Path, ev.OldHash, ev.NewHash)
}
```

## Types

### `FileEvent`

| Field       | Type        | Description                          |
|-------------|-------------|--------------------------------------|
| `Path`      | `string`    | Absolute path of the changed file    |
| `ChangedAt` | `time.Time` | Timestamp when the change was polled |
| `OldHash`   | `string`    | SHA-256 hash before the change       |
| `NewHash`   | `string`    | SHA-256 hash after the change        |

### `Watcher`

- **`New(paths []string, interval time.Duration) (*Watcher, error)`** — creates a
  watcher and records baseline hashes. Returns an error if any file is unreadable.
- **`Poll() ([]FileEvent, error)`** — checks all paths once; returns events for
  files whose hash has changed since the last poll (or since `New`).

## Integration

The `confsnap watch` CLI command uses this package to auto-snapshot files
whenever a change is detected during a long-running session.
