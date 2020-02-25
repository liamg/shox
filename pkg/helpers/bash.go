package helpers

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

// BashHelper runs a bash command and displays the output
type BashHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *BashHelper) UpdateInterval() time.Duration {
	return time.Second * 5
}

// Run runs a bash command and displays the output
func (h *BashHelper) Run(config string) string {
	path, err := exec.LookPath("bash")
	if err != nil {
		return "bash not in PATH"
	}
	cmd := exec.Command(path, "-c", config)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "command failed"
	}
	return strings.ReplaceAll(out.String(), "\n", "")
}

func init() {
	Register("bash", &BashHelper{})
}
