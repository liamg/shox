package helpers

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/mem"
)

// MemoryHelper shows the current memory usage as a percentage
type MemoryHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *MemoryHelper) UpdateInterval() time.Duration {
	return time.Second
}

// Run returns the current memory usage as a percentage
func (h *MemoryHelper) Run(config string) string {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "?"
	}
	return fmt.Sprintf("%.0f%%", v.UsedPercent)
}

func init() {
	Register("memory", &MemoryHelper{})
}
