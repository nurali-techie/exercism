package gigasecond

import (
	"time"
)

const testVersion = 4

func AddGigasecond(in time.Time) time.Time {
	oneGigaSecond := time.Second * 1000000000
	out := in.Add(time.Duration(oneGigaSecond))
	return out
}
