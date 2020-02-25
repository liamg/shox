package helpers

import (
	"fmt"
	"time"

	"github.com/distatus/battery"
)

// BatteryHelper shows the current battery charge level as a percentage
type BatteryHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *BatteryHelper) UpdateInterval() time.Duration {
	return time.Minute
}

// Run returns the current battery charge level as a percentage
func (h *BatteryHelper) Run(config string) string {

	var current float64
	var total float64
	batteries, err := battery.GetAll()
	if err != nil {
		return "?"
	}

	for _, bat := range batteries {
		current += bat.Current
		total += bat.Full
	}

	return fmt.Sprintf("%.0f%%", 100*current/total)
}

func init() {
	Register("battery", &BatteryHelper{})
}
