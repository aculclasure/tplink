// Package archerc9v1 provides client logic for interacting with the TP Link
// Archer C9 V1 wifi router.
package archerc9v1

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type logger interface {
	Printf(string, ...interface{})
}

// Client manages communication with a TP Link Archer
// C9 V1 wifi router.
type Client struct {
	userName, password, encodedBasicAuth string
	baseURL                              *url.URL
	httpClient                           *http.Client
	logger                               logger
}

// New returns a Client to an Archer C9 V1 wifi router given a user name, password,
// url, http.Client, and a type implementing the logger interface. If any of the user
// name, password, url, or http.Client arguments are invalid, then an error is returned.
// If the lgr argument is nil, then a default logger of type *log.Logger is created
// for the Client which logs to os.Stderr.
func New(userName, password, rawURL string, httpClient *http.Client, lgr logger) (*Client, error) {
	if len(userName) == 0 {
		return nil, errors.New("got empty value for userName (want a valid user name)")
	}
	if len(password) == 0 {
		return nil, errors.New("got empty value for password (want a non-empty password)")
	}
	// TODO: Check for "/" suffix in rawUrl and append it if not found
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "got error parsing rawUrl")
	}
	if !strings.HasPrefix(u.Scheme, "http") && !strings.HasPrefix(u.Scheme, "https") {
		return nil, errors.New("got invalid scheme in rawUrl (want http or https): " + rawURL)
	}
	if len(u.Host) == 0 {
		return nil, errors.New("got empty hostname in rawUrl (want http://hostname or https://hostname): " + rawURL)
	}
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if lgr == nil {
		lgr = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return &Client{
		userName:         userName,
		password:         password,
		baseURL:          u,
		encodedBasicAuth: base64.StdEncoding.EncodeToString([]byte(userName + ":" + password)),
		httpClient:       httpClient,
		logger:           lgr,
	}, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, and if so,
// it should always be specified without a preceding slash. The items specified in body will
// be encoded as request body parameters.
func (c *Client) NewRequest(method, urlStr string, body map[string]string) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "got error creating new request URL")
	}

	data := url.Values{}
	for k, v := range body {
		data.Set(k, v)
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "got error creating new "+method+" request to: "+u.String())
	}

	req.Header.Set("Referer", c.baseURL.String())
	req.Header.Set("Cookie", "Authorization=Basic "+c.encodedBasicAuth)
	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// CheckResponse checks the response for errors and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if r.StatusCode/100 == 2 {
		return nil
	}

	return fmt.Errorf("got status code %d (want 200-299) when doing %s request to %s",
		r.StatusCode, r.Request.Method, r.Request.URL.String())
}
