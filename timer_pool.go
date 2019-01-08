package smux

import (
	"sync"
	"time"
)

// TimerPool is sync.Pool that cache time.Timer. Since there are some sublte details
// when reusing time.Timer, this package make life easier.
type TimerPool struct {
	pool sync.Pool
}

// NewTimerPool create TimerPool
func NewTimerPool() *TimerPool {
	return &TimerPool{
		pool: sync.Pool{
			New: func() interface{} {
				timer := time.NewTimer(time.Hour)
				if !timer.Stop() {
					<-timer.C
				}
				return timer
			},
		},
	}
}

// Get cached time.Timer in pool, use Reset method of time.Timer to set new deadline
// Example:
//     timer := timerPool.Get()
//     timer.Reset(1 * time.Second)
func (p *TimerPool) Get() *time.Timer {
	return p.pool.Get().(*time.Timer)
}

// Put time.Timer back to pool. consumed must be set true if time.Timer.C has been consumed
// Example:
//     select {
//     case <-timer.C:
//         timerPool.Put(timer, true)
//     default:
//         timerPool.Put(timer, false)
//     }
func (p *TimerPool) Put(timer *time.Timer, consumed bool) {
	if !consumed {
		if !timer.Stop() {
			<-timer.C
		}
	}
	p.pool.Put(timer)
}
