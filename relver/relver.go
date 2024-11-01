package relver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
)

const (
	// URLEndpointDefault is the default URL to fetch the latest Go release version.
	// Use MockURL for testing.
	URLEndpointDefault = "https://go.dev/dl/?mode=json"
	// RequestTimeoutDefault is the default timeout duration for the HTTP request.
	RequestTimeoutDefault = 5
)

// Variables to monkey-patch the functions for testing.
//
//nolint:gochecknoglobals // allow these global variables for testing
var (
	// urlEndpoint is the URL to fetch the latest Go release version.
	urlEndpoint = URLEndpointDefault
	// requestTimeout is the timeout duration for the HTTP request.
	requestTimeout = RequestTimeoutDefault
	// ioReadAll is a copy of io.ReadAll to be used in testing to monkey-patch the
	// function.
	ioReadAll = io.ReadAll
)

// ============================================================================
//  Functions
// ============================================================================

// ----------------------------------------------------------------------------
//  asGoSemVer
// ----------------------------------------------------------------------------

// asGoSemVer returns the given version string in Go semantic versioned string
// format. It will remove the prefix "v" or "go" and return the string with "v"
// prefix.
//
// Note that it will not validate the version string.
func asGoSemVer(input string) string {
	if input == "" {
		return ""
	}

	input = strings.TrimPrefix(input, "v")
	input = strings.TrimPrefix(input, "go")

	return "v" + input
}

// ----------------------------------------------------------------------------
//  Compare
// ----------------------------------------------------------------------------

// Compare returns an integer comparing two versions.
// It is similar to `semver.Compare()` but it accepts "go" and "v" prefixed
// versions.
//
//		 0: if a == b
//		-1: if a < b
//		+1: if a > b
//
//	  a |  b
//	 +1 0 -1
//
// An invalid semantic version string is considered less than a valid one. All
// invalid semantic version strings compare equal to each other as well.
func Compare(a, b string) int {
	a = asGoSemVer(a)
	b = asGoSemVer(b)

	return semver.Compare(a, b)
}

// ----------------------------------------------------------------------------
//  cURL
// ----------------------------------------------------------------------------

// cURL simply fetches the target URL content. Similar to `http.Get()` but with
// a request timeout.
//
// Note that setting the timeout to 0 is not allowed. If the timeout duration
// is 0, it will use the default timeout duration of 5 seconds.
func cURL(targetURL string, timeOut int) ([]byte, error) {
	if timeOut == 0 {
		timeOut = RequestTimeoutDefault
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new HTTP request")
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to GET HTTP request")
	}

	defer response.Body.Close()

	// Read response body. To force an error during testing, mock the ioReadAll
	respBody, err := ioReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read response")
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"fail to GET response from: %s\nStatus: %s\nResponse body: %s",
			targetURL,
			response.Status,
			string(respBody),
		)
	}

	return respBody, nil
}

// ----------------------------------------------------------------------------
//  Get
// ----------------------------------------------------------------------------

// Get fetches the latest Go release version from the URL endpoint and returns
// the version string.
//
// Note that it will timeout after 5 seconds. To change the timeout duration,
// use GetWithTimeout.
func Get() (string, error) {
	return GetWithTimeout(requestTimeout)
}

// ----------------------------------------------------------------------------
//  GetWithTimeout
// ----------------------------------------------------------------------------

// GetWithTimeout fetches the latest Go release version from the URL endpoint
// with a custom timeout duration.
//
// Note that setting the timeout to 0 is not allowed. If the timeout duration
// is 0, it will use the default timeout duration of 5 seconds.
func GetWithTimeout(timeout int) (string, error) {
	respBody, err := cURL(urlEndpoint, timeout)
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve the released version info")
	}

	// parse the JSON response body to Releases struct
	var releases Releases
	if err := json.Unmarshal(respBody, &releases); err != nil {
		return "", errors.Wrap(err, "failed to parse/unmarshal the JSON response body")
	}

	for _, unit := range releases {
		if unit.Stable && unit.Version != "" {
			return unit.Version, nil
		}
	}

	return "", errors.New("no stable releases found")
}

// ----------------------------------------------------------------------------
// MockURL
// ----------------------------------------------------------------------------

// MockURL is a helper function to mock the URL endpoint for testing. It returns
// a function to cleanup ( a function to defer restore the original URL endpoint).
func MockURL(url string) func() {
	originalURL := urlEndpoint
	urlEndpoint = url

	return func() {
		urlEndpoint = originalURL
	}
}

// ----------------------------------------------------------------------------
//  SetTimeout
// ----------------------------------------------------------------------------

// SetTimeout sets the timeout duration for the HTTP request. The default
// timeout duration is 5 seconds.
//
// Note that even if the timeout duration is set to 0, it will use the default
// timeout during the request (see cURL function).
func SetTimeout(timeoutSec int) {
	requestTimeout = timeoutSec
}

// ============================================================================
//  Types and Methods
// ============================================================================

// ----------------------------------------------------------------------------
//  Type: Releases
// ----------------------------------------------------------------------------

// Releases is a struct to hold each release information from the JSON response.
type Release struct {
	FileName  string `json:"filename"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Version   string `json:"version"`
	SHA256SUM string `json:"sha256"`
	Size      int64  `json:"size"`
	Kind      string `json:"kind"`
}

// Releases is a struct to hold the JSON response from the URL endpoint.
type Releases []struct {
	Version string    `json:"version"`
	Stable  bool      `json:"stable"`
	Files   []Release `json:"files"`
}
