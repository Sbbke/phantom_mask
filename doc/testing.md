## Unit and Integration Test

I wrote unit tests for the system function, and some integration tests for the APIs endpoint (currently first API only).

You can run the test by using the command below:
```bash
# unit test for etl helper 
$ go test -v ./app/initial 
# integration test, please run these command with db initialized
$ go test test/integration/controllers/pharmacy_controller_test.go
$ go test test/integration/middleware/middleware_test.go
```
## Test Coverage Report
(doesn't include integration test)
For test coverage, please run the command below:

```bash
$ go test -cover ./app/...
```