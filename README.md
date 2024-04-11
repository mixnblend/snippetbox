# SnippetBox <!-- omit in toc -->

Small gist like web application built in Go, following along with the book Let's Go! by Alex Edwards

## Table Of Contents<!-- omit in toc -->

- [Quick Start (local)](#quick-start-local)
- [Available Commands](#available-commands)
- [Contributing](#contributing)
- [Testing](#testing)

## Quick Start (local)

1. Install all [Prerequisites](docs/prerequisites.md)
2. run `go run ./cmd/web -help` to see a list of options while running the application.
3. run `go run ./cmd/web` to start the application.
4. run `docker-compose up` to start the mysql database.
5. run `go test -v ./cmd/web -tags test_all` to run all tests

**[⬆ back to top](#table-of-contents)**

## Available Commands

to see a list of available commands run:

```bash
make help
```

**[⬆ back to top](#table-of-contents)**

## Contributing

Before contributing please read through everything in [Contributing](docs/contributing.md)

**[⬆ back to top](#table-of-contents)**

## Testing

Run all tests:

```bash
make test
```

Run all unit tests:

```bash
make test/unit
```

Run all end to end tests:

```bash
make test/e2e
```

- [Contributing](docs/CONTRIBUTING.md)
- [Prerequisites](docs/prerequisites.md)

**[⬆ back to top](#table-of-contents)**
