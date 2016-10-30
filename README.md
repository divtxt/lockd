
# lockd

A distributed lock service.

See [API.md](API.md) for server API documentation.

## Development

Run using:

    go run main.go

Then, run the following commands in another terminal:

Lock an entry:

```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:2080/api/lock \
    -d '{"name":"Foo"}'
```

Try to lock the same name again and you should get a 409 Conflict error.


To check if an entry is locked:

```
curl -i http://127.0.0.1:2080/api/lock?name=Foo
```

To unlock the entry:

```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:2080/api/unlock \
    -d '{"name":"Foo"}'
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
- [ ] Lock wait for raft commit

Lock Features:

- [ ] Lock status query API
- [ ] Lock request/client id field
- [ ] Request/client id override
- [ ] Field size limit checks
- [ ] External test script

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

- [ ] cluster config file & config file loading
- [ ] command line param: cluster config file name
- [ ] raft consensus module
- [ ] raft rpc endpoints
- [ ] peer rpc service
- [ ] proxy Lock & Unlock API calls to raft leader

Misc:
- [ ] Add metrics & logging
- [ ] Expose raft details e.g. leader, term
- [x] Stop `-help` from showing "-httptest.serve"
- [ ] Vendorize dependencies
