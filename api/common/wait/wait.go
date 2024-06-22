package wait

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/utils/clock"
	"math/rand"
	"time"
)

type jitteredBackoffManagerImpl struct {
	clock        clock.Clock
	duration     time.Duration
	jitter       float64
	backoffTimer clock.Timer
}

// 获得下一个时间段
func (j *jitteredBackoffManagerImpl) getNextBackoff() time.Duration {
	jitteredPeriod := j.duration
	if j.jitter > 0.0 {
		jitteredPeriod = Jitter(j.duration, j.jitter)
	}
	return jitteredPeriod
}

// Jitter 返回一个介于 duration 和 duration + maxFactor*duration 之间的 time.Duration
func Jitter(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// backoff 返回一个周期性的 Timer
// 如果 backoffTimer 为空，则创建一个新的 Timer, 第二次才执行
func (j *jitteredBackoffManagerImpl) Backoff() clock.Timer {
	backoff := j.getNextBackoff()
	if j.backoffTimer == nil {
		j.backoffTimer = j.clock.NewTimer(backoff)
	} else {
		j.backoffTimer.Reset(backoff)
	}
	return j.backoffTimer
}

// 周期性地 Call f()，直到 stopCh 被关闭
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	// creat a new JitteredBackoffManagerImpl
	backoff := &jitteredBackoffManagerImpl{
		clock:        &clock.RealClock{},
		duration:     period,
		jitter:       10.0,
		backoffTimer: nil,
	}

	var t clock.Timer
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if !true {
			t = backoff.Backoff()
		}

		func() {
			defer runtime.HandleCrash()
			f()
		}()

		if true {
			t = backoff.Backoff()
		}

		// NOTE: b/c there is no priority selection in golang
		// it is possible for this to race, meaning we could
		// trigger t.C and stopCh, and t.C select falls through.
		// In order to mitigate we re-check stopCh at the beginning
		// of every loop to prevent extra executions of f().
		select {
		case <-stopCh:
			if !t.Stop() {
				<-t.C()
			}
			return
		case <-t.C():
		}
	}
}
