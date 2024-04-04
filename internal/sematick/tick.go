// Sematick is a ticker that sends a single message at specified intervals
// to every subscribed channel. This ticker can be stopped, paused or resumed.
package sematick

import (
	"sync"
	"sync/atomic"
	"time"
)

type Ticker struct {
	mux      sync.Mutex       // mutex to access the channels
	channels []chan time.Time // channels to store the subscribed timers

	count       int
	state       uint32       // 0 = stopped, 1 = running, 2 = paused
	ticker      *time.Ticker // principal timer
	tickerMux   sync.Mutex   // mutex for the underlying ticker
	interval    time.Duration
	stopChannel chan struct{}
}

// Creates a new sematick, pushes time.Time messages at the
// desired interval and can start either as 0 which equals stopped
// or 1 which equals running.
func NewTicker(interval time.Duration, initialState int) *Ticker {
	t := &Ticker{
		interval: interval,
	}

	go func() {
		t.tickerMux.Lock()
		atomic.StoreUint32(&t.state, uint32(initialState))

		t.stopChannel = make(chan struct{})
		t.ticker = time.NewTicker(interval)
		t.tickerMux.Unlock()

		t.tick()
	}()

	return t
}

// Creates and returns a new channel that recieves time messages
// at the specified interval.
func (t *Ticker) Subscribe() <-chan time.Time {
	t.mux.Lock()
	defer t.mux.Unlock()

	// creating new channel and appending it to the subscription list
	new_channel := make(chan time.Time, 1)
	t.channels = append(t.channels, new_channel)

	return new_channel
}

// Gets the amount of times the main ticker has sent a message.
func (t *Ticker) Count() int {
	return t.count
}

// Changes the state to pause, ticks still happen but are not
// sent to the subscribed channels.
func (t *Ticker) Pause() {
	atomic.StoreUint32(&t.state, 2)
}

// Changes state to resume, this works as a play button when
// the ticker was initialized in the stopped state.
func (t *Ticker) Resume() {
	atomic.StoreUint32(&t.state, 1)
}

// Stops ticking and quits the main goroutine.
func (t *Ticker) Stop() {
	t.tickerMux.Lock()
	defer t.tickerMux.Unlock()

	stopped := atomic.LoadUint32(&t.state) == 0
	if !stopped && t.stopChannel != nil {
		t.ticker.Stop()
		t.stopChannel <- struct{}{}
	}

	atomic.StoreUint32(&t.state, 0)
}

func (t *Ticker) tick() {
	for {
		select {
		case tick := <-t.ticker.C:
			t.mux.Lock()
			if atomic.LoadUint32(&t.state) != 1 {
				t.mux.Unlock()
				continue
			}
			for i := range t.channels {
				select {
				case t.channels[i] <- tick:
				default:
					t.count++
				}
			}
			t.mux.Unlock()
		case <-t.stopChannel:
			return
		default:
			t.count++
		}
	}
}
