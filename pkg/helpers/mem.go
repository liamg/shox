package helpers

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/mem"
)

type MemoryHelper struct {
}

func (h *MemoryHelper) UpdateInterval() time.Duration {
	return time.Second
}

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
