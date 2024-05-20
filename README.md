# Teeworlds econ server library

![golangci-lint](https://github.com/theobori/teeworlds-econ/actions/workflows/lint.yml/badge.svg)

This library is highly flexible and thread-safe by design, it allows you to interact with a Teeworlds econ server.

## 📖 Build and run

You only need the following requirements:

- [Go](https://golang.org/doc/install) 1.22.3

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## 🧪 Tests

There are some tests that require a running Teeworlds econ server, feel free to use the `econ_server.sh` script that create a DDNet server and a econ server.

It also requires the following environment variables.

| Name | Description | Optional
| - | - | - |
`ECON_DEBUG` | Enables the debug verbose | Yes
`ECON_PORT` | Specifies the econ port | Yes (7000 by default)
`ECON_PASSWORD` | Specifies the econ password | No

By default, the test are using `localhost` as host.

Below, an example of running the test.

```bash
# Override some variables for the container and the Go tests
export ECON_PORT=1234
export ECON_PASSWORD=just_a_test_password

# Run the script
./econ_test.sh

# Wait for the econ server being ready
sleep 10

# Run the Go tests
make test
```
