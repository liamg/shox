package helpers

import "time"

type TimeHelper struct {
}

func (h *TimeHelper) UpdateInterval() time.Duration {
	return time.Second
}

func (h *TimeHelper) Run(config string) string {
	return time.Now().Format("15:04:05")
}

func init() {
	Register("time", &TimeHelper{})
}
