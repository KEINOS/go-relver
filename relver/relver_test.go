//nolint:paralleltest // these tests are not parallel safe
package relver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Get
// ----------------------------------------------------------------------------

func TestGet_bad_server_response(t *testing.T) {
	// Requesting a non-existing file
	urlTestSrv := spawnTestServer(t, "unknown.json", 0)

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	ver, err := Get()

	require.Error(t, err,
		"if all the releases does not contain the version, it should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "failed to retrieve the released version info",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "404 Not Found",
		"error message should contain the underlying error reason")
}

func TestGet_bad_url_to_request(t *testing.T) {
	badURL := "http://192.168.0.%31/" // invalid URL

	cleanUp := MockURL(badURL)
	t.Cleanup(cleanUp)

	ver, err := Get()

	require.Error(t, err,
		"malformed URL should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "failed to retrieve the released version info",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "failed to create a new HTTP request",
		"error message should contain the underlying error reason")
	require.Contains(t, err.Error(), "invalid URL escape",
		"error message should contain the underlying error reason")
}

func TestGet_empty_version(t *testing.T) {
	// Test with empty version releases
	urlTestSrv := spawnTestServer(t, "data_versions_are_empty.json", 0)

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	ver, err := Get()

	require.Error(t, err,
		"if all the releases does not contain the version, it should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "no stable releases found",
		"error message should contain the error reason")
}

func TestGet_non_json_response(t *testing.T) {
	// Test with non-JSON response
	urlTestSrv := spawnTestServer(t, "data_not_a_json.json", 0)

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	ver, err := Get()

	require.Error(t, err,
		"if all the releases does not contain the version, it should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "failed to parse/unmarshal the JSON response body",
		"error message should contain the error reason")
}

func TestGet_reading_resp_body_failure(t *testing.T) {
	// Test with empty version releases
	urlTestSrv := spawnTestServer(t, "data_versions_are_empty.json", 0)

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	// Mock ioReadAll to force reading the response body to fail
	oldIoReadAll := ioReadAll

	t.Cleanup(func() {
		ioReadAll = oldIoReadAll
	})

	ioReadAll = func(_ io.Reader) ([]byte, error) {
		return nil, errors.New("forced error to read body")
	}

	ver, err := Get()

	require.Error(t, err,
		"if all the releases does not contain the version, it should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "failed to retrieve the released version info",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "forced error to read body",
		"error message should contain the underlying error reason")
}

func TestGet_server_delay(t *testing.T) {
	// Test with slow server response
	urlTestSrv := spawnTestServer(t, "data_golden.json", 10)

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	SetTimeout(1) // set request timeout to 1 second
	t.Cleanup(func() {
		SetTimeout(0)
	})

	ver, err := Get()

	require.Error(t, err,
		"if the server response exceeds the timeout, it should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "failed to retrieve the released version info",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "context deadline exceeded",
		"error message should contain the underlying error reason")
}

func TestGet_stable_and_unstable_mix_releases(t *testing.T) {
	urlTestSrv := spawnTestServer(t, "data_stable_and_unstable_mix.json", 0)
	expectVer := "go1.22.8"

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	actualVer, err := Get()

	require.NoError(t, err,
		"well formatted JSON response should not return an error")
	require.Equal(t, expectVer, actualVer,
		"stable and unstable mix releases should return the latest stable release")
}

func TestGet_unstable_releases(t *testing.T) {
	// Test with all unstable releases
	urlTestSrv := spawnTestServer(t, "data_unstable.json", 0)

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	ver, err := Get()

	require.Error(t, err,
		"if all the releases are unstable, it should return an error")
	require.Empty(t, ver,
		"returned version should be empty on error")
	require.Contains(t, err.Error(), "no stable releases found",
		"error message should contain the error reason")
}

func TestGet_zeroTimeout(t *testing.T) {
	urlTestSrv := spawnTestServer(t, "data_golden.json", 0)
	expectVer := "go1.23.2"

	cleanUp := MockURL(urlTestSrv)
	t.Cleanup(cleanUp)

	// Setting the timeout to 0 should use the default timeout during the
	// request.
	SetTimeout(0)
	t.Cleanup(func() {
		SetTimeout(RequestTimeoutDefault)
	})

	actualVer, err := Get()

	require.NoError(t, err,
		"unexpected error on zero timeout")
	require.Equal(t, expectVer, actualVer,
		"unexpected response body content")
}

// ----------------------------------------------------------------------------
//  Test Helpers
// ----------------------------------------------------------------------------

// spawnTestServer creates a test server that responds the given file content
// which is located in the testdata directory.
// It returns the endpoint URL.
func spawnTestServer(t *testing.T, fileName string, delay int) string {
	t.Helper()

	// Create an HTTP handler with a delay
	//nolint:varnamelen // w,r are common names for http handlers
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add delay before serving the file
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Second)
		}

		// Serve the requested file
		http.FileServer(http.Dir("testdata")).ServeHTTP(w, r)
	})

	// Start a test server
	server := httptest.NewServer(handler)
	serverURL := server.URL + "/" + fileName

	t.Cleanup(func() {
		server.Close()
	})

	return serverURL
}
