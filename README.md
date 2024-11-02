# go-relver - Retrieve the latest Go version

[![go1.22+](https://img.shields.io/badge/Go-1.22+-blue?logo=go)](https://github.com/KEINOS/go-relver/blob/main/go.mod "Supported versions")
[![GoDoc](https://godoc.org/github.com/KEINOS/go-relver?status.svg)](https://pkg.go.dev/github.com/KEINOS/go-relver)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/KEINOS/go-relver/blob/main/LICENSE.md)

A simple Go package to **fetch the latest and stable release version of Go**.

## Usage

```shellsession
$ # Install the module
$ go get github.com/KEINOS/go-relver
...
```

```go
// Use the package
import "github.com/KEINOS/go-relver/relver"
```

```go
// Examples
import "github.com/KEINOS/go-relver/relver"

func Example() {
    // Get the latest release version of Go
    goVerLatest, err := relver.Get()
    if err != nil {
        log.Panic(err)
    }

    fmt.Println(goVerLatest)
    // Output: go1.23.2
}

func Example_set_request_timeout() {
    // Set the request timeout to 10 seconds.
    // Default is 5 seconds.
    relver.SetTimeout(10)

    goVerLatest, err := relver.Get()
    if err != nil {
        log.Panic(err)
    }

    fmt.Println(goVerLatest)
    // Output: go1.23.2
}

func Example_compare_go_versions() {
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
```

## License

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/KEINOS/go-relver/blob/main/LICENSE.md)


## Contributing

[![go1.22+](https://img.shields.io/badge/Go-1.22+-blue?logo=go)](https://github.com/KEINOS/go-relver/blob/main/go.mod "Supported versions")
[![GoDoc](https://godoc.org/github.com/KEINOS/go-relver?status.svg)](https://pkg.go.dev/github.com/KEINOS/go-relver)

- Any feedbacks for the improvement are welcome!
- Branch to PR: `main` <- `feature/your-idea`
- Bug report: Open an [issue](https://github.com/KEINOS/go-relver/issues) with a simple code that reproduces the issue. That helps a lot.

### Statuses

[![Basic Tests](https://github.com/KEINOS/go-relver/actions/workflows/basic-tests.yml/badge.svg)](https://github.com/KEINOS/go-relver/actions/workflows/basic-tests.yml)
[![CodeQL](https://github.com/KEINOS/go-relver/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/KEINOS/go-relver/actions/workflows/github-code-scanning/codeql)
[![Dependabot Updates](https://github.com/KEINOS/go-relver/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/KEINOS/go-relver/actions/workflows/dependabot/dependabot-updates)
[![Go Report Card](https://goreportcard.com/badge/github.com/KEINOS/go-relver)](https://goreportcard.com/report/github.com/KEINOS/go-relver)
[![codecov](https://codecov.io/gh/KEINOS/go-relver/graph/badge.svg?token=x0G0b1zptx)](https://codecov.io/gh/KEINOS/go-relver)
