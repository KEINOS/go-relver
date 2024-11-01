package relver_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"

	"github.com/KEINOS/go-relver/relver"
)

// ----------------------------------------------------------------------------
//  Example Basic Usage
// ----------------------------------------------------------------------------

func Example() {
	// Spawn a test server that responds the given file content to avoid accessing
	// the real URL (avoid E2E testing).
	cleanup := spawnTestServer("data_golden.json")
	defer cleanup()

	// The actual function call to get the latest Go version.
	goVerLatest, err := relver.Get()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(goVerLatest)
	// Output: go1.23.2
}

// By default, the request timeout is set to 5 seconds. You can change the timeout
// duration by using the `SetTimeout` function.
func Example_set_request_timeout() {
	// Spawn a test server that responds the given file content.
	cleanup := spawnTestServer("data_golden.json")
	defer cleanup()

	// Set the request timeout to 10 seconds
	relver.SetTimeout(10)

	goVerLatest, err := relver.Get()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(goVerLatest)
	// Output: go1.23.2
}

func Example_compare_go_versions() {
	// Spawn a test server that responds the given file content.
	// This file contains a future Go version such as "go5.55.55".
	cleanup := spawnTestServer("data_future_ver.json")
	defer cleanup()

	goVerLatest, err := relver.Get()
	if err != nil {
		log.Panic(err)
	}

	// Get the local Go version
	goVerLocal := runtime.Version()

	// Compare the two Go versions
	if relver.Compare(goVerLatest, goVerLocal) > 0 {
		fmt.Println("You are using an outdated Go version.")
	} else {
		fmt.Println("You are using the latest Go version.")
	}
	//
	// Output: You are using an outdated Go version.
}

// ----------------------------------------------------------------------------
//  Compare
// ----------------------------------------------------------------------------

func ExampleCompare() {
	fmt.Println("Same version with 'go' prefix:",
		relver.Compare("go1.23.2", "go1.23.2"),
	)

	fmt.Println("Same version with 'v' prefix:",
		relver.Compare("v1.23.2", "v1.23.2"),
	)

	fmt.Println("Same version with diff prefix:",
		relver.Compare("go1.23.2", "v1.23.2"),
	)

	fmt.Println("Version 'a' is greater than 'b':",
		relver.Compare("go1.23.2", "go1.23.1"),
	)

	fmt.Println("Version 'a' is less than 'b':",
		relver.Compare("go1.23.1", "go1.23.2"),
	)

	fmt.Println("Version 'a' is empty:",
		relver.Compare("", "go1.23.2"),
	)
	fmt.Println("Version 'b' is empty:",
		relver.Compare("go1.23.1", ""),
	)
	//
	// Output:
	// Same version with 'go' prefix: 0
	// Same version with 'v' prefix: 0
	// Same version with diff prefix: 0
	// Version 'a' is greater than 'b': 1
	// Version 'a' is less than 'b': -1
	// Version 'a' is empty: -1
	// Version 'b' is empty: 1
}

// ----------------------------------------------------------------------------
//  Test Helpers
// ----------------------------------------------------------------------------

// spawnTestServer creates a test server that responds the given file content
// which is located in the testdata directory.
// It returns will returns the endpoint URL and a cleanup function. It is the
// callers responsibility to call the cleanup function.
func spawnTestServer(fileName string) func() {
	// Start a test server
	server := httptest.NewServer(http.FileServer(http.Dir("testdata")))

	serverURL := server.URL + "/" + fileName
	cleanUp := relver.MockURL(serverURL) // stub the test server URL

	return func() {
		cleanUp()      // Recover the original URL after the test
		server.Close() // Close the test server
	}
}
