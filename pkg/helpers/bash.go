package helpers

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

type BashHelper struct {
}

func (h *BashHelper) UpdateInterval() time.Duration {
	return time.Second * 5
}

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
