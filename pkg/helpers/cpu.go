package helpers

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type CPUHelper struct {
}

func (h *CPUHelper) UpdateInterval() time.Duration {
	return time.Second
}

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
