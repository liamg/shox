package terminal

import (
	"io"
	"time"
)

const BackOffReadInitialSleepDuration = time.Millisecond
const BackOffReadMaxSleepDuration = time.Millisecond * 16

// similar to io.Copy with sleep when no data is received
func lazyCopy(dst io.Writer, src io.Reader, hideFuncs ...func() bool) error {

	buffer := make([]byte, 4096)

	backOffDelay := BackOffReadInitialSleepDuration

	var hide bool

	for {
		size, err := src.Read(buffer)
		if size > 0 {
			hide = false
			for _, f := range hideFuncs {
				if f() {
					hide = true
					break
				}
			}
			if !hide {
				if _, err := dst.Write(buffer[:size]); err != nil {
					return err
				}
			}
			backOffDelay = BackOffReadInitialSleepDuration
		}
		if err != nil {
			return err
		}
		if size == 0 {
			// if there was no data to read, wait a little while before trying again
			time.Sleep(backOffDelay)
			backOffDelay = backOffDelay * 2
			if backOffDelay > BackOffReadMaxSleepDuration {
				backOffDelay = BackOffReadMaxSleepDuration
			}
		}
	}
}
