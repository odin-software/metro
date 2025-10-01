package tenjin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/newspaper"
	"github.com/odin-software/metro/internal/tenjin/analysis"
	"github.com/odin-software/metro/internal/tenjin/observation"
)

// Tenjin is the central brain that observes and manages the simulation
type Tenjin struct {
	eventChannel chan interface{}
	observation  *observation.Collector
	analysis     *analysis.MetricsEngine
	logger       *analysis.MetricsLogger
	newspaper    *newspaper.Newspaper
	ticker       *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewTenjin creates a new Tenjin brain
func NewTenjin(totalTrains int) (*Tenjin, error) {
	// Create event channel with buffer of 500
	eventChannel := make(chan interface{}, 500)

	// Create observation layer
	collector := observation.NewCollector(eventChannel)

	// Create analysis layer
	metricsEngine := analysis.NewMetricsEngine(totalTrains)

	// Create metrics logger
	metricsDir := control.DefaultConfig.LogsDirectory + "tenjin/"
	logger, err := analysis.NewMetricsLogger(metricsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics logger: %w", err)
	}

	// Create newspaper
	news, err := newspaper.NewNewspaper()
	if err != nil {
		return nil, fmt.Errorf("failed to create newspaper: %w", err)
	}

	// Create context for lifecycle management
	ctx, cancel := context.WithCancel(context.Background())

	tenjin := &Tenjin{
		eventChannel: eventChannel,
		observation:  collector,
		analysis:     metricsEngine,
		logger:       logger,
		newspaper:    news,
		ticker:       time.NewTicker(control.DefaultConfig.TenjinTickRate),
		ctx:          ctx,
		cancel:       cancel,
	}

	return tenjin, nil
}

// GetEventChannel returns the send-only channel for trains to emit events
func (t *Tenjin) GetEventChannel() chan<- interface{} {
	return t.eventChannel
}

// Start begins Tenjin's operations
// Should be called after all trains have been initialized
func (t *Tenjin) Start() {
	control.Log("Tenjin: Starting brain operations")

	// Start observation collector in its own goroutine
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.observation.Start()
	}()

	// Start main processing loop
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.run()
	}()
}

// run is the main processing loop - runs every second
func (t *Tenjin) run() {
	control.Log("Tenjin: Main loop started")

	for {
		select {
		case <-t.ctx.Done():
			control.Log("Tenjin: Shutting down")
			return

		case <-t.ticker.C:
			// Collect events from observation layer
			events := t.observation.Collect()

			// Process events through analysis layer
			if len(events) > 0 {
				t.analysis.ProcessEvents(events)
			}

			// Get formatted metrics output
			output := t.analysis.GetFormattedOutput()

			// Log to file
			if err := t.logger.Log(output); err != nil {
				control.Log(fmt.Sprintf("Tenjin: Error logging metrics: %v", err))
			}

			// Also print to stdout if configured
			if control.DefaultConfig.StdLogs {
				fmt.Print(output)
			}

			// Check if newspaper needs new edition (daily + on-demand)
			if t.newspaper.NeedsNewEdition() && !t.newspaper.IsGenerating() {
				control.Log("Tenjin: Triggering newspaper generation...")

				// Generate in background goroutine (non-blocking)
				go func() {
					ctx, cancel := context.WithTimeout(t.ctx, 2*time.Minute)
					defer cancel()

					metrics := t.analysis.GetMetrics()
					if err := t.newspaper.GenerateEdition(ctx, &metrics); err != nil {
						control.Log(fmt.Sprintf("Tenjin: Newspaper generation error: %v", err))
					}
				}()
			}
		}
	}
}

// Stop gracefully shuts down Tenjin
func (t *Tenjin) Stop() {
	control.Log("Tenjin: Stopping...")

	// Cancel context to stop goroutines
	t.cancel()

	// Stop ticker
	t.ticker.Stop()

	// Close event channel (trains should already be stopped)
	close(t.eventChannel)

	// Wait for goroutines to finish
	t.wg.Wait()

	// Close logger
	if err := t.logger.Close(); err != nil {
		control.Log(fmt.Sprintf("Tenjin: Error closing logger: %v", err))
	}

	control.Log("Tenjin: Stopped successfully")
}

// GetMetrics returns current metrics (useful for UI or external queries)
func (t *Tenjin) GetMetrics() analysis.Metrics {
	return t.analysis.GetMetrics()
}

// GetNewspaper returns the newspaper instance (for UI access)
func (t *Tenjin) GetNewspaper() *newspaper.Newspaper {
	return t.newspaper
}
