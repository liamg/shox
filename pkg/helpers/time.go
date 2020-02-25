package helpers

import "time"

// TimeHelper shows the current time
type TimeHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *TimeHelper) UpdateInterval() time.Duration {
	return time.Second
}

// Run returns the current time
func (h *TimeHelper) Run(config string) string {
	return time.Now().Format("15:04:05")
}

func init() {
	Register("time", &TimeHelper{})
}
