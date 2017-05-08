
# lockd

[![Build Status](https://travis-ci.org/divtxt/lockd.svg?branch=master)](https://travis-ci.org/divtxt/lockd)

A distributed lock service.


## Development

Run using:

```
go get -t ./...
go run main.go -cluster integtests/config/1node.json -id 1
```


## Lock API

The lock service has a REST-ful interface: every lock is a resource at the path `/lock/:name`:

### GET /lock/:name

Check if the given entry is locked.

This will return one of the following status codes:

- `200 OK` - Entry is locked.
- `404 Not Found` - Entry is unlocked.

Example:

```
curl -i http://localhost:2081/lock/foo
```


### POST /lock/:name

Lock the given entry.

This will return one of the following status codes:

- `200 OK` - Success - entry is now locked.
- `409 Conflict` - Failed - entry is locked.

Example:

```
curl -i -X POST http://localhost:2081/lock/foo
```


### DELETE /lock/:name

Unlock the given entry.

This will return one of the following status codes:

- `200 OK` - Success - entry is now locked.
- `409 Conflict` - Failed - entry is locked.

Example:

```
curl -i -X DELETE http://localhost:2081/lock/foo
```


## TODO

Basic Daemon:

- [x] command line param: server listen address
- [x] logging setup
- [x] gin logging to golang standard logging bridge

Basic Single-Node Locking:

- [x] Lock & Unlock API endpoints
- [x] Lock state persistence API
- [x] In-memory state persistence
- [x] Implement single-node locking service
- [x] Refactor to use raft as single-node cluster
- [x] Lock wait for raft commit

Lock Features:

- [x] Lock status query API
- [ ] Lock request/client id field
- [ ] Request/client id override
- [ ] Lock acquire wait time

Error Handling:

- [ ] Internal error shows original stack trace
- [ ] API timeouts
- [ ] Field content and size limit checks

Integration Testing:

- [x] External test script
- [x] Travis runs server and test script

Lock State Persistence:

- [ ] Boltdb state persistence
- [ ] Command line param: choice of persistence mode
- [ ] Command line param: Boltdb state file name

Lock TTL:

- [ ] Lock TTL params
- [ ] Expired Lock Unlocker

Lock Admin:

- [ ] List Locks API
- [ ] Web UI: basic assets
- [ ] Web UI: list locks
- [ ] Web UI: lock
- [ ] Web UI: unlock

Multi-Node Locking:

- [x] raft consensus module
- [ ] cluster support
- [ ] raft snapshotting
- [ ] raft rpc endpoints
- [ ] peer rpc service
- [ ] proxy Lock & Unlock API calls to raft leader

Clients:

- [ ] Ruby client
- [x] Python client
- [ ] Java/Scala client

Servers:

- [ ] Homebrew installable binary
- [ ] Linux binary

Misc:

- [ ] Add metrics & logging
- [ ] Expose raft details e.g. leader, term
- [x] Stop `-help` from showing "-httptest.serve"
- [ ] Vendorize dependencies
