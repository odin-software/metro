package sematick

import (
	"sync"
	"sync/atomic"
	"time"
)

type Ticker struct {
	mux      sync.Mutex       // mutex to access the channels
	channels []chan time.Time // channels to store the subscribed timers

	count         int
	state         uint32       // 0 = stopped, 1 = running, 2 = paused
	ticker        *time.Ticker // principal timer
	tickerMux     sync.Mutex   // mutex for the underlying ticker
	interval      time.Duration
	stopChannel   chan struct{}
	pauseChannel  chan struct{}
	resumeChannel chan struct{}
}

func NewTicker(interval time.Duration, initialState int) *Ticker {
	t := &Ticker{
		interval: interval,
	}

	go func() {
		t.tickerMux.Lock()
		atomic.StoreUint32(&t.state, uint32(initialState))

		t.stopChannel = make(chan struct{})
		t.pauseChannel = make(chan struct{})
		t.resumeChannel = make(chan struct{})
		t.ticker = time.NewTicker(interval)
		t.tickerMux.Unlock()

		t.tick()
	}()

	return t
}

func (t *Ticker) Subscribe() <-chan time.Time {
	t.mux.Lock()
	defer t.mux.Unlock()

	// creating new channel and appending it to the subscription list
	new_channel := make(chan time.Time, 1)
	t.channels = append(t.channels, new_channel)

	return new_channel
}

func (t *Ticker) Count() int {
	return t.count
}

func (t *Ticker) Pause() {
	atomic.StoreUint32(&t.state, 2)
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

func (t *Ticker) Resume() {
	atomic.StoreUint32(&t.state, 1)
	if t.ticker == nil {
		t.ticker = time.NewTicker(t.interval)
	}
}

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
		if atomic.LoadUint32(&t.state) != 1 {
			continue
		}
		select {
		case tick := <-t.ticker.C:
			t.mux.Lock()
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
