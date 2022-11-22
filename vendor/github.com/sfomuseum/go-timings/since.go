package timings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"
	"time"
)

// type SinceMonitor implements the `Monitor` interface providing a background timings mechanism that tracks the duration of time
// between events.
type SinceMonitor struct {
	Monitor
	done_ch   chan bool
	since_ch  chan *SinceResponse
	start     time.Time
	lastevent time.Time
	mu        *sync.RWMutex
}

// SinceResponse is a struct containing information related to a "since" timing event.
type SinceResponse struct {
	// Message is an optional string that was included with a `Signal` event
	Message string `json:"message"`
	// Duration is the string representation of a `time.Duuration` which is the amount of time that elapsed between `Signal` events
	Duration string `json:"duration"`
	// Timestamp is the Unix timestamp when the `SinceResponse` was created
	Timestamp int64 `json:"timestamp"`
}

func init() {
	ctx := context.Background()
	RegisterMonitor(ctx, "since", NewSinceMonitor)
}

// NewSinceMonitor creates a new `SinceMonitor` instance that will dispatch notifications using a time.Ticker configured
// by 'uri' which is expected to take the form of:
//
//	since://
func NewSinceMonitor(ctx context.Context, uri string) (Monitor, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	done_ch := make(chan bool)
	since_ch := make(chan *SinceResponse)

	mu := new(sync.RWMutex)

	t := &SinceMonitor{
		done_ch:  done_ch,
		since_ch: since_ch,
		mu:       mu,
	}

	return t, nil
}

// Start() will cause background monitoring to begin, dispatching notifications to wr in
// the form of JSON-encoded `SinceResponse` values.
func (t *SinceMonitor) Start(ctx context.Context, wr io.Writer) error {

	if !t.start.IsZero() {
		return fmt.Errorf("Monitor has already been started")
	}

	now := time.Now()
	t.start = now
	t.lastevent = now

	go func() {

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.done_ch:
				return
			case rsp := <-t.since_ch:

				enc := json.NewEncoder(wr)
				err := enc.Encode(rsp)

				if err != nil {
					log.Printf("Failed to encode response, %v", err)
				}
			}
		}
	}()

	return nil
}

// Stop() will cause background monitoring to be halted.
func (t *SinceMonitor) Stop(ctx context.Context) error {
	t.done_ch <- true
	return nil
}

// Signal will cause the background monitors since to be incremented by one.
func (t *SinceMonitor) Signal(ctx context.Context, args ...interface{}) error {

	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	duration := time.Since(t.lastevent)

	rsp := &SinceResponse{
		Timestamp: now.Unix(),
		Duration:  duration.String(),
	}

	if len(args) > 0 {

		switch args[0].(type) {
		case string:
			rsp.Message = args[0].(string)
		default:
			rsp.Message = fmt.Sprintf("%v", args[0])
		}

	}

	t.since_ch <- rsp

	t.lastevent = now
	return nil
}
