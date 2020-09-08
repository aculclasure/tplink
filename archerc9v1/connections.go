package archerc9v1 

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const loginPageIndicator = "<title>tp-link archer c9"

// Represents a node that is connected to the router.
type Connection struct {
	MacAddress string `json:"mac_addr"`
	IPAddress  string `json:"ip_addr"`
	Name       string `json:"name"`
}

// Represents a response from the query to get wireless/wired
// connections from the router.
type connections struct {
	success bool
	timeout bool
	Data    []*Connection `json:"data"`
}

type connectionsResponse struct {
	response  *http.Response
	transport string
}

// GetWiredConnections returns a slice of Connections representing wired
// connections to the router or returns an error otherwise.
func (c *Client) GetWiredConnections() ([]*Connection, error) {
	req, err := c.NewRequest("POST", "data/map_access_wire_client_grid.json", nil)
	if err != nil {
		return nil, errors.Wrap(err, "got error creating request to get wired connections")
	}

	c.logger.Printf("sending request to get wired connections as (%s %s) ...",
		req.Method, req.URL)
	resp, err := c.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "got error doing request to get wired connections")
	}

	if err = CheckResponse(resp); err != nil {
		return nil, errors.Wrap(err, "got error in response to get wired connections")
	}
	connectionsResponse := &connectionsResponse{response: resp, transport: "wired"}
	return connectionsResponse.getConnections()
}

// GetWirelessConnections returns a slice of Connection representing wireless
// connections to the router or returns an error otherwise.
func (c *Client) GetWirelessConnections() ([]*Connection, error) {
	req, err := c.NewRequest("POST", "data/map_access_wireless_client_grid.json", nil)
	if err != nil {
		return nil, errors.Wrap(err, "got error creating request to get wireless connections")
	}

	c.logger.Printf("sending request to get wireless connections as (%s %s) ...",
		req.Method, req.URL)
	resp, err := c.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "got error doing request to get wireless connections")
	}

	if err = CheckResponse(resp); err != nil {
		return nil, errors.Wrap(err, "got error in response to get wireless connections")
	}
	connectionsResponse := &connectionsResponse{response: resp, transport: "wireless"}
	return connectionsResponse.getConnections()
}

// getConnections creates a slice of Connections from the response body
// or returns an error otherwise.
func (r *connectionsResponse) getConnections() ([]*Connection, error) {
	data, err := ioutil.ReadAll(r.response.Body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("got error reading body of %s connections response", r.transport))
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("got empty body in %s connections response", r.transport)
	}

	if strings.Contains(strings.ToLower(string(data)), loginPageIndicator) {
		return nil,
			fmt.Errorf("got the TP Link Archer C9 login webpage as %s connections response (want JSON), check your login credentials",
				r.transport)
	}

	c := new(connections)
	if err = json.Unmarshal(data, c); err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("got error trying to decode %s connections response into JSON (tried to decode %s)",
				r.transport, string(data)))
	}

	return c.Data, nil
}
