# Local development environment and tests

## Service tools

The service tools is a set of bash scripts that provide everything you need to run, debug
and run tests for cloud uno locally.

The entry point script to run the service tools is in `scripts/service-tools.sh`.

To get a list of all the different commands and options available run the following:
```bash
bash ./scripts/service-tools.sh --help
```

## Running tests

For integration tests, all the external dependencies (e.g. emulators or databases in docker containers)
are set up by the service tools
before running the test suite and are torn down after.

To run only unit tests:
```bash
bash scripts/service-tools.sh --tests unit
```

To run only integration tests:
```bash
bash scripts/service-tools.sh --tests integration
```

To run integration and unit tests:
```bash
bash scripts/service-tools.sh --tests unit,integration
```
