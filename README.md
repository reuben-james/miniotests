# miniotests

A test suite for a Minio Object Store implemented natively in Go

## Prerequisites

* docker
* docker compose (Go)

## ENVIRONMENT

The following environment variables are required to be set, in order to run the test suite:
```
# For all tests
export MINIO_SERVER="minio.example.com"
export MINIO_PORT="9000"
export MINIO_ACCESS_KEY="YOUR_ACCESS_KEY"
export MINIO_SECRET_KEY="YOUR_SECRET_KEY"
export MINIO_SECURE="true" # or "false" if not using SSL/TLS
# For read-only tests
export MINIO_READONLY_ACCESS_KEY="READONLY_ACCESS_KEY"
export MINIO_READONLY_SECRET_KEY="READONLY_SECRET_KEY"
```

## Run the tests

Run against an insecrue Minio
```
go test -v
```

Run against a Secure Minio
```
go test -v -args -useTLS=true
```

Example output
```
=== RUN   TestAuthentication
--- PASS: TestAuthentication (0.02s)
=== RUN   TestAuthorization
--- PASS: TestAuthorization (0.01s)
=== RUN   TestCRUDOperations
--- PASS: TestCRUDOperations (0.09s)
=== RUN   TestConnectivity
--- PASS: TestConnectivity (0.01s)
=== RUN   TestHealthChecks
--- PASS: TestHealthChecks (0.01s)
PASS
ok      github.com/reuben-james/miniotests      0.157s
```

Run just the Authentication test against a secure Minio
```
go test -v -run TestAuthentication -args -useTLS=true
```

## Local Development

### Stand up a local dev instance

#### Insecure
```
cd docker
docker compose -f docker-compose-insecure,yml up -d
```

#### Secure (TLS Enabled)
```
cd docker
docker compose up -d
```

#### MC Client

Install the MC client for manual testing
```
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
sudo mv mc /usr/local/bin/
```

Setup an mc alias
```
MINIO_URL=http://${MINIO_SERVER}:${MINIO_PORT}
# OR IF SECURE
MINIO_URL=https://${MINIO_SERVER}:${MINIO_PORT}

mc alias set test/ ${MINIO_URL}
$ Enter Access Key: 
$ Enter Secret Key: 
```

Make and verify a test bucket
```
mc mb test/test-bucket-default/

mc ls test
# OUTPUT
[YYYY-MM-DD HH:MM:SS UTC]     0B test-bucket-default/
```

Upload a manual test file over
```
mc cp resources/manual-test-file.txt test/test-bucket-default/
```

Verify it arrived
```
mc ls test/test-bucket-default
# OUTPUT
[YYYY-MM-DD HH:MM:SS UTC]    19B STANDARD manual-test-file.txt
```