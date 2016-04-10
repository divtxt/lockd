
# lockd

A distributed lock service.


## Development

Run using:

    go run main.go

Then, run the following commands in another terminal:

Lock an entry:

```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:2080/lock \
    -d '{"name":"Foo"}'
```

Try to lock the same name again and you should get an error (409 Conflict ?).


To check if an entry is locked:

```
curl -i http://127.0.0.1:2080/lock?name=Foo
```

To unlock the entry:

```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:2080/unlock \
    -d '{"name":"Foo"}'
```


## TODO

Basic Daemon:

- [x] command line param: server listen address
- [x] logging setup
- [x] gin logging to golang standard logging bridge

Basic Single-Node Locking:

- [x] Lock & Unlock API endpoints
- [ ] Lock state persistence API
- [ ] In-memory state persistence
- [ ] Single node driver

Lock State Persistence:

- [ ] Boltdb state persistence
- [ ] Command line param: choice of persistence mode
- [ ] Command line param: Boltdb state file name

Lock TTL:

- [ ] Lock TTL params
- [ ] Expired Lock Unlocker

Lock Admin:

- [ ] List Locks API
- [ ] Web UI: bootstrap assets
- [ ] Web UI: list locks
- [ ] Web UI: lock
- [ ] Web UI: unlock

Raft Locking:

- [ ] cluster config file & config file loading
- [ ] command line param: cluster config file name
- [ ] raft consensus module
- [ ] raft rpc endpoints
- [ ] peer rpc service
- [ ] proxy Lock & Unlock API calls to raft leader

Misc:
- [ ] Stop `-help` from showing "-httptest.serve"
- [ ] Vendorize dependencies
