package archerc9v1 

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_Reboot(t *testing.T) {
	validRebootResponse := `
	<TR>
		<TD class="h2" id="t_restart">Rebooting...</TD>
		</TR>
		<TR>
		<TD class="h2" id="t_complete" style="display:none">Completed!</TD>
		</TR>
	`
	testCases := []*struct {
		description string
		input       RoundTripFunc
		expectError bool
	}{
		{
			description: "Error getting response",
			input: func(r *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("got error while rebooting")
			},
			expectError: true,
		},
		{
			description: "404 status code in response",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
					Body:       ioutil.NopCloser(strings.NewReader("")),
					Request:    r,
				}, nil
			},
			expectError: true,
		},
		{
			description: "Empty response",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("")),
					Request:    r,
				}, nil
			},
			expectError: true,
		},
		{
			description: "No indication of reboot completion",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(loginPageIndicator)),
					Request:    r,
				}, nil
			},
			expectError: true,
		},
		{
			description: "Partial response",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("Rebooting...")),
					Request:    r,
				}, nil
			},
			expectError: true,
		},
		{
			description: "Reboot completed",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(validRebootResponse)),
					Request:    r,
				}, nil
			},
			expectError: false,
		},
	}

	for _, tt := range testCases {
		client = &Client{
			baseURL:    validURL,
			httpClient: NewTestClient(tt.input),
			logger:     defaultLogger}
		err := client.Reboot()
		if tt.expectError {
			if err == nil {
				t.Fatalf("FAIL: %s\n\t%v.Reboot() did not return an expected error",
					tt.description, client)
			}
		} else if err != nil {
			t.Fatalf("FAIL: %s\n\t%v.Reboot() returned an unexpected error: %v",
				tt.description, client, err)
		}
		t.Logf("PASS: %s", tt.description)
	}
}
