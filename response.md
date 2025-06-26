# Response

For a full overview and better navigation, please refer to the [README](README.md).

## A. Required Information
### A.1. Requirement Completion Rate
- [o] List all pharmacies open at a specific time and on a day of the week if requested.
  - Implemented at /api/v1/pharmacies/open
- [o] List all masks sold by a given pharmacy, sorted by mask name or price.
  - Implemented at /api/v1/pharmacies/masks
- [o] List all pharmacies with more or less than x mask products within a price range.
  - Implemented at /api/v1/pharmacies/filter
- [o] The top x users by total transaction amount of masks within a date range.
  - Implemented at /api/v1/pharmacies/users/top
- [o] The total number of masks and dollar value of transactions within a date range.
  - Implemented at /api/v1/pharmacies/transactions/summary
- [o] Search for pharmacies or masks by name, ranked by relevance to the search term.
  - Implemented at /api/v1/pharmacies/search
- [o] Process a user purchases a mask from a pharmacy, and handle all relevant data changes in an atomic transaction.
  - Implemented at /api/v1/pharmacies/purchase
### A.2. API Document

Please click [here](./api.md) to API document 

### A.3. Import Data Commands
Please run these two script commands to migrate the data into the database.

```bash
# Start the database & Run schema migrations and ETL/loading
$ docker compose -f docker-compose-db.yaml up -d
$ docker compose -f docker-compose-init.yaml up 
```
## B. Bonus Information

### B.1. Test Coverage Report

I wrote unit tests for the system function, and some integration tests for the APIs endpoint (currently first API only).

You can run the test script by using the command below:
```bash
# unit test for etl helper 
$ go test -v ./app/initial 
# integration test, please run these command with db initialized
$ go test test/integration/controllers/pharmacy_controller_test.go
$ go test test/integration/middleware/middleware_test.go
```

For test coverage, please run the command below:
```bash
$ go test -cover ./app/...
```
### B.2. Dockerized

On the local machine, please follow the commands below to run the service.

```bash
# Build the image
$ docker compose -f docker-compose.yaml build
# Start the database 
$ docker compose -f docker-compose-db.yaml up -d
# Run schema migrations and ETL/loading
$ docker compose -f docker-compose-init.yaml up --exit-code-from init-preprocess
# Start the service
$ docker compose -f docker-compose.yaml up
```


## C. Other Information

### Table of contents

Please click [`here`](./README.md).

- --
