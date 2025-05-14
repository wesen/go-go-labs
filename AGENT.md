# AGENT.md

## Go Commands

- Build: `make build` (runs go generate ./... and go build ./...)
- Test: `make test` (runs all tests with `go test ./...`)
- Run specific test: `go test ./path/to/package -run TestName`

## Python Commands

- Install dependencies: `pip install -r requirements.txt`
- Run Flask app: `python app.py`
- Run specific test: `python -m unittest path/to/test.py::TestClass::test_method`

## Code Style Guidelines

- Go: Uses gofmt, go 1.23+, github.com/pkg/errors for error wrapping
- Go: Uses zerolog for logging, cobra for CLI, viper for config
- Go: Follow standard naming (CamelCase for exported, camelCase for unexported)
- Python: PEP 8 formatting, uses logging module for structured logging
- Python: Try/except blocks with specific exceptions and error logging
- Use interfaces to define behavior, prefer structured concurrency
- do not try to fix linting errors or use make lint
- do not try to run the application itself unless asked to
- run goimports -w ... before building

<goGuidelines>
When implementing go interfaces, use the var _ Interface = &Foo{} to make sure the interface is always implemented correctly.
When building web applications, use htmx, bootstrap and the templ templating language.
Always use a context argument when appropriate.
Use cobra for command-line applications.
Use the "defaults" package name, instead of "default" package name, as it's reserved in go.
Use github.com/pkg/errors for wrapping errors.
When starting goroutines, use errgroup.
</goGuidelines>
