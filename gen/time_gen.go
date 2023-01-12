package gen

import "time"

type timeBetween struct {
	start       time.Time
	durationGen Generator[int64]
}

func (t timeBetween) GenerateOne() time.Time {
	newDuration := t.durationGen.GenerateOne()
	return t.start.Add(time.Duration(newDuration))
}

func (t timeBetween) GenerateN(n uint) []time.Time {
	res := make([]time.Time, n)
	for i := uint(0); i < n; i++ {
		res[i] = t.GenerateOne()
	}
	return res
}

func TimeBetween(start time.Time, end time.Time) Generator[time.Time] {
	if start.Equal(end) {
		return Only(start)
	}
	actualStart := start
	if start.After(end) {
		actualStart = end
	}
	actualEnd := end
	if end.Before(start) {
		actualEnd = start
	}

	dur := actualEnd.Sub(actualStart)
	return timeBetween{actualStart, Between(int64(0), int64(dur))}
}
