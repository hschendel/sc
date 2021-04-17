package sc

import "time"

type Timeout struct {
	timeoutAt time.Time
}

func NewTimeout(duration time.Duration) Timeout {
	return Timeout{timeoutAt: time.Now().Add(duration)}
}

func (t *Timeout) TimeLeft() time.Duration {
	return t.timeoutAt.Sub(time.Now())
}

func (t *Timeout) Reached() bool {
	return t.TimeLeft() <= 0
}
