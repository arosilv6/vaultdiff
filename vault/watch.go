package vault

import (
	"context"
	"time"
)

// WatchEvent is emitted when a secret version changes.
type WatchEvent struct {
	Path       string
	OldVersion int
	NewVersion int
	Error      error
}

// Watcher polls a Vault secret path for version changes.
type Watcher struct {
	client   *Client
	path     string
	interval time.Duration
}

// NewWatcher creates a Watcher for the given path and poll interval.
func NewWatcher(client *Client, path string, interval time.Duration) *Watcher {
	return &Watcher{
		client:   client,
		path:     path,
		interval: interval,
	}
}

// Watch polls the secret path and sends events on the returned channel.
// It stops when ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		last := -1
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				versions, err := List w.path)
				if err != nil {
					ch <- WatchEvent{: err}
					continue
				}
				latest :=(versions)
				if latest != last {
					if last != -1 {
						ch <- WatchEvent{Path: w.path, OldVersion: last, NewVersion: latest}
					}
					last = latest
				}
			}
		}
	}()
	return ch
}
