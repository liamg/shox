package helpers

import (
	"strconv"
	"time"
)

// TimeHelper shows the current time
type UntilHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *UntilHelper) UpdateInterval() time.Duration {
	return time.Second
}

// Run returns the current time
func (h *UntilHelper) Run(config string) string {
	ts, err := strconv.ParseInt(config, 10, 64)
	if err != nil {
		return "bad timestamp"
	}
	remaining := time.Until(time.Unix(ts, 0)).Truncate(time.Second)
	if remaining < 0 {
		return "0s"
	}
	return remaining.String()
}

func init() {
	Register("until", &UntilHelper{})
}
