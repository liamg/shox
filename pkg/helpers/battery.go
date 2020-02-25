package helpers

import (
	"fmt"
	"time"

	"github.com/distatus/battery"
)

type BatteryHelper struct {
}

func (h *BatteryHelper) UpdateInterval() time.Duration {
	return time.Minute
}

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
