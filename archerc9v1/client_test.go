package archerc9v1 

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	user, password   = "user", "password"
	validEncodedAuth = base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
	validRawURL      = "http://my-tplink-rtr/"
	validURL, _      = url.Parse(validRawURL)
	validRequestBody = map[string]string{"operation": "read"}
	client           *Client
	defaultLogger    = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func TestNew(t *testing.T) {
	testCases := []struct {
		description string
		userName    string
		password    string
		url         string
		httpClient  *http.Client
		lgr         logger
		expected    *Client
		expectError bool
	}{
		{
			description: "Empty user name",
			userName:    "",
			password:    "password",
			url:         validRawURL,
			httpClient:  nil,
			lgr:         nil,
			expected:    nil,
			expectError: true,
		},
		{
			description: "Empty password",
			userName:    user,
			password:    "",
			url:         validRawURL,
			httpClient:  nil,
			lgr:         nil,
			expected:    nil,
			expectError: true,
		},
		{
			description: "Empty URL",
			userName:    user,
			password:    password,
			url:         "",
			httpClient:  nil,
			lgr:         nil,
			expected:    nil,
			expectError: true,
		},
		{
			description: "Bad URL (no scheme)",
			userName:    user,
			password:    password,
			url:         "my-tp-link-rtr",
			httpClient:  nil,
			lgr:         nil,
			expected:    nil,
			expectError: true,
		},
		{
			description: "Bad URL (no hostname)",
			userName:    user,
			password:    password,
			url:         "http://",
			httpClient:  nil,
			lgr:         nil,
			expected:    nil,
			expectError: true,
		},
		{
			description: "Valid arguments",
			userName:    "user",
			password:    "password",
			url:         validRawURL,
			httpClient:  nil,
			lgr:         defaultLogger,
			expected: &Client{
				userName:         "user",
				password:         "password",
				baseURL:          validURL,
				encodedBasicAuth: validEncodedAuth,
				httpClient:       &http.Client{},
				logger:           defaultLogger,
			},
			expectError: false,
		},
	}

	for _, test := range testCases {
		got, err := New(test.userName, test.password, test.url, test.httpClient, test.lgr)
		if test.expectError {
			if err == nil {
				t.Fatalf("FAIL: %s\n\tNew(%s, %s, %s) expected an error, got %v",
					test.description, test.userName, test.password, test.url, got)
			}
		} else {
			if err != nil {
				t.Fatalf("FAIL: %s\n\tNew(%s, %s, %s) returns an unexpected error %s",
					test.description, test.userName, test.password, test.url, err.Error())
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("FAIL: %s\n\tNew(%s, %s, %s) expected %v, got %v",
					test.description, test.userName, test.password, test.url, test.expected, got)
			}
		}
		t.Logf("PASS: %s", test.description)
	}
}

func TestNewRequest(t *testing.T) {
	client, _ = New(user, password, validRawURL, nil, nil)
	testDescription := "Valid request"
	method, relativeURL, expectedURL := "POST", "/foo", client.baseURL.String()+"foo"
	failureMsgPrefix := fmt.Sprintf("FAIL: %s\n\tNewRequest(%s, %s, %v)",
		testDescription, method, relativeURL, validRequestBody)
	expectedCookieHeader := fmt.Sprintf("Authorization=Basic %s",
		base64.StdEncoding.EncodeToString([]byte(user+":"+password)))
	req, err := client.NewRequest(method, relativeURL, validRequestBody)
	if err != nil {
		t.Fatalf("%s returned an unexpected error: %s", failureMsgPrefix, err.Error())
	}
	if req.URL.String() != expectedURL {
		t.Fatalf("%s expected URL %s, got %s", failureMsgPrefix, expectedURL, req.URL.String())
	}
	if req.Header.Get("Cookie") != expectedCookieHeader {
		t.Fatalf("%s does not contain expected header (want %s=%s)",
			failureMsgPrefix, "Cookie", expectedCookieHeader)
	}
	if req.Header.Get("Referer") != client.baseURL.String() {
		t.Fatalf("%s does not contain expected header (want %s=%s)",
			failureMsgPrefix, "Referer", client.baseURL.String())
	}

	t.Logf("PASS: %s", testDescription)
}

func TestNewRequest_withInvalidParameters(t *testing.T) {
	client, _ = New(user, password, validRawURL, nil, nil)

	errorCases := []*struct {
		description, method, relativeURL string
		body                             map[string]string
	}{
		{
			"Invalid method",
			"bad method",
			"/foo",
			validRequestBody,
		},
		{
			"Invalid relative URL",
			"POST",
			":foo",
			validRequestBody,
		},
	}
	for _, tt := range errorCases {
		_, err := client.NewRequest(tt.method, tt.relativeURL, tt.body)
		if err == nil {
			t.Fatalf("FAIL: %s\n\tNewRequest(%s, %s, %v) expected an error",
				tt.description, tt.method, tt.relativeURL, tt.body)
		}
		t.Logf("PASS: %s", tt.description)
	}
}

type RoundTripFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestDo(t *testing.T) {
	type foo struct {
		A string
	}
	testCases := []*struct {
		description, method, relativeURL string
		fn                               RoundTripFunc
		expected                         *foo
		expectError                      bool
	}{
		{
			description: "Valid request",
			method:      "POST",
			relativeURL: "/foo",
			fn: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(`{"A": "a"}`)),
					Header:     make(http.Header),
				}, nil
			},
			expected:    &foo{A: "a"},
			expectError: false,
		},
		{
			description: "HTTP Error",
			method:      "POST",
			relativeURL: "/bad-uri",
			fn: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("404 error not found")
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range testCases {
		func() {
			c := NewTestClient(tt.fn)
			client, _ = New(user, password, validRawURL, c, nil)
			req, _ := client.NewRequest(tt.method, tt.relativeURL, validRequestBody)
			resp, err := client.Do(req)
			if resp != nil {
				defer resp.Body.Close()
			}
			if tt.expectError {
				if err == nil {
					t.Fatalf("FAIL: %s\n\tDo(%s, %s) did not return an expected error",
						tt.description, tt.method, tt.relativeURL)
				}
			} else {
				if err != nil {
					t.Fatalf("FAIL: %s\n\tDo(%s, %s) returned an unexpected error: %v",
						tt.description, tt.method, tt.relativeURL, err)
				} else {
					body := new(foo)
					err = json.NewDecoder(resp.Body).Decode(body)
					if err != nil {
						t.Fatalf("FAIL: %s\n\tDo(%s, %s) got error when decoding JSON response: %v",
							tt.description, tt.method, tt.relativeURL, err)
					}
					if !reflect.DeepEqual(body, tt.expected) {
						t.Fatalf("FAIL: %s\n\tDo(%s, %s) returned %v, want %v",
							tt.description, tt.method, tt.relativeURL, body, tt.expected)
					}
				}
			}
			t.Logf("PASS: %s", tt.description)
		}()
	}
}

func TestCheckResponse(t *testing.T) {
	testCases := []*struct {
		description string
		resp        *http.Response
		expectError bool
	}{
		{
			description: "200 status code in response",
			resp:        &http.Response{StatusCode: 200},
			expectError: false,
		},
		{
			description: "500 status code in response",
			resp: &http.Response{
				StatusCode: 500,
				Request:    &http.Request{Method: http.MethodPost, URL: validURL},
			},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		err := CheckResponse(tt.resp)
		if tt.expectError {
			if err == nil {
				t.Fatalf("FAIL: %s\n\tCheckResponse(%v) did not return an expected error",
					tt.description, tt.resp)
			}
		} else if err != nil {
			t.Fatalf("FAIL: %s\n\tCheckResponse(%v) returned an unexpected error %s",
				tt.description, tt.resp, err.Error())
		}
		t.Logf("PASS: %s", tt.description)
	}
}
