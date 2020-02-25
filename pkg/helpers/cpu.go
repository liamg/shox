package helpers

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// CPUHelper shows the current cpu usage as a percentage
type CPUHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *CPUHelper) UpdateInterval() time.Duration {
	return time.Second
}

// Run returns the current cpu usage as a percentage
func (h *CPUHelper) Run(config string) string {
	values, err := cpu.Percent(0, false)
	if err != nil {
		return "?"
	}
	return fmt.Sprintf("%.0f%%", values[0])
}

func init() {
	Register("cpu", &CPUHelper{})
}
