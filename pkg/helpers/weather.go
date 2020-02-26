package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// DefaultWeatherFormat is the default wttr.in weather format.
// See https://github.com/chubin/wttr.in#one-line-output for more details
const DefaultWeatherFormat = "1"

// WeatherHelper shows the current weather
type WeatherHelper struct {
}

// UpdateInterval returns the minimum time period before the helper should run again
func (h *WeatherHelper) UpdateInterval() time.Duration {
	return time.Minute * 5
}

// Run returns the current weather
func (h *WeatherHelper) Run(config string) string {
	var weatherFormat string
	if config != "" {
		weatherFormat = config
	} else {
		weatherFormat = DefaultWeatherFormat
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "http://wttr.in/", nil)
	if err != nil {
		return fmt.Sprintf("err: %s", err)
	}

	q := req.URL.Query()
	q.Add("format", weatherFormat)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("err: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("err: %s", err)
	}

	return strings.TrimRight(string(body), "\n")
}

func init() {
	Register("weather", &WeatherHelper{})
}
