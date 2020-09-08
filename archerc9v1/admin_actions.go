package archerc9v1

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

// Reboot reboots the router and returns nil if the reboot is successful.
// Otherwise, an error is returned.
func (c *Client) Reboot() error {
	req, err := c.NewRequest(
		"GET", "userRpm/SysRebootRpm.htm", nil)
	if err != nil {
		return errors.Wrap(err, "got error creating request to do reboot")
	}
	q := req.URL.Query()
	q.Set("Reboot", "Reboot")
	req.URL.RawQuery = q.Encode()

	c.logger.Printf("sending reboot request as (%s %s) ...",
		req.Method, req.URL)
	resp, err := c.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, "got error doing request to reboot")
	}

	if err = CheckResponse(resp); err != nil {
		return errors.Wrap(err, "got error in response to reboot")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "got error reading body in response to reboot")
	}

	if len(data) == 0 {
		return errors.New("got empty body in response to reboot")
	}

	if !strings.Contains(string(data), "Rebooting...") ||
		!strings.Contains(string(data), "Completed!") {
		return fmt.Errorf("got invalid response body (want response indicating rebooting has completed): %s",
			string(data))
	}

	c.logger.Printf("reboot completed successfully...response from reboot call: %s", string(data))
	return nil
}
