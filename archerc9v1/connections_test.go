package archerc9v1 

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var (
	macAddress                   = "00-00-00-00-00-00"
	ipAddress                    = "10.100.100.1"
	hostName                     = "my-tp-link-rtr"
	validConnectionsJSONResponse = fmt.Sprintf(
		`{"data":[{"mac_addr":"%s","ip_addr":"%s","name":"%s"}]}`,
		macAddress, ipAddress, hostName)
	validConnectionsResult      = []*Connection{{MacAddress: macAddress, IPAddress: ipAddress, Name: hostName}}
	connectionResponseTestCases = []*struct {
		description string
		input       RoundTripFunc
		expected    []*Connection
		expectError bool
	}{
		{
			description: "Error getting response",
			input: func(r *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("error")
			},
			expected:    nil,
			expectError: true,
		},
		{
			description: "403 status code in response",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 403,
					Body:       ioutil.NopCloser(strings.NewReader("")),
					Request:    r,
				}, nil
			},
			expected:    nil,
			expectError: true,
		},
		{
			description: "Valid response",
			input: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(validConnectionsJSONResponse)),
					Request:    r,
				}, nil
			},
			expected:    validConnectionsResult,
			expectError: false,
		},
	}
)

func TestConnectionsResponse_getConnections(t *testing.T) {
	testCases := []*struct {
		description string
		input       *connectionsResponse
		expected    []*Connection
		expectError bool
	}{
		{
			description: "Empty response body",
			input: &connectionsResponse{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("")),
				},
				transport: "wired",
			},
			expected:    nil,
			expectError: true,
		},
		{
			description: "Login page returned",
			input: &connectionsResponse{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(loginPageIndicator)),
				},
				transport: "wired",
			},
			expected:    nil,
			expectError: true,
		},
		{
			description: "Non-JSON response",
			input: &connectionsResponse{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("this is not JSON")),
				},
				transport: "wired",
			},
			expected:    nil,
			expectError: true,
		},
		{
			description: "Invalid JSON response",
			input: &connectionsResponse{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`{"not": "expected json"}`)),
				},
				transport: "wired",
			},
		},
		{
			description: "Valid JSON response",
			input: &connectionsResponse{
				response: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(validConnectionsJSONResponse)),
				},
				transport: "wired",
			},
			expected: []*Connection{
				{MacAddress: macAddress, IPAddress: ipAddress, Name: hostName},
			},
			expectError: false,
		},
	}

	for _, tt := range testCases {
		func() {
			if tt.input.response != nil && tt.input.response.Body != nil {
				defer func() {
					if err := tt.input.response.Body.Close(); err != nil {
						t.Logf("got error closing response body: %v", err)
					}
				}()
			}
			got, err := tt.input.getConnections()
			if tt.expectError {
				if err == nil {
					t.Fatalf("FAIL: %s\n\t%v.getConnections() did not return an expected error",
						tt.description, tt.input)
				}
			} else {
				if err != nil {
					t.Fatalf("FAIL: %s\n\t%v.getConnections() returned an unexpected error: %v",
						tt.description, tt.input, err)
				}
				if !reflect.DeepEqual(got, tt.expected) {
					t.Fatalf("FAIL: %s\n\t%v.getConnections() returned %v, want %v",
						tt.description, tt.input, got, tt.expected)
				}
			}
			t.Logf("PASS: %s", tt.description)
		}()
	}
}

func TestClient_GetWiredConnections(t *testing.T) {
	for _, tt := range connectionResponseTestCases {
		client = &Client{
			userName:         user,
			password:         password,
			encodedBasicAuth: validEncodedAuth,
			baseURL:          validURL,
			httpClient:       NewTestClient(tt.input),
			logger:           defaultLogger,
		}
		got, err := client.GetWiredConnections()
		if tt.expectError {
			if err == nil {
				t.Fatalf("FAIL: %s\n\t%v.GetWiredConnections() did not return an expected error",
					tt.description, client)
			}
		} else {
			if err != nil {
				t.Fatalf("FAIL: %s\n\t%v.GetWiredConnections() returned an unexpected error: %v",
					tt.description, client, err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("FAIL: %s\n\t%v.GetWiredConnections() returned %v, want %v",
					tt.description, client, got, tt.expected)
			}
		}
		t.Logf("PASS: %s", tt.description)
	}
}

func TestClient_GetWirelessConnections(t *testing.T) {
	for _, tt := range connectionResponseTestCases {
		client = &Client{
			userName:         user,
			password:         password,
			encodedBasicAuth: validEncodedAuth,
			baseURL:          validURL,
			httpClient:       NewTestClient(tt.input),
			logger:           defaultLogger,
		}
		got, err := client.GetWirelessConnections()
		if tt.expectError {
			if err == nil {
				t.Fatalf("FAIL: %s\n\t%v.GetWirelessConnections() did not return an expected error",
					tt.description, client)
			}
		} else {
			if err != nil {
				t.Fatalf("FAIL: %s\n\t%v.GetWirelessConnections() returned an unexpected error: %v",
					tt.description, client, err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("FAIL: %s\n\t%v.GetWirelessConnections() returned %v, want %v",
					tt.description, client, got, tt.expected)
			}
		}
		t.Logf("PASS: %s", tt.description)
	}
}
