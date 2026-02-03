package ui

import (
	"fmt"
	"sync"
	"time"
)

// Spinner shows a loading animation
type Spinner struct {
	message string
	frames  []string
	stop    chan struct{}
	stopped bool
	mu      sync.Mutex
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				return
			default:
				s.mu.Lock()
				if !s.stopped {
					fmt.Printf("\r%s %s ", s.frames[i%len(s.frames)], s.message)
				}
				s.mu.Unlock()
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return
	}
	s.stopped = true
	close(s.stop)
	fmt.Print("\r\033[K") // Clear line
}
