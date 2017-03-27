
# lockd

[![Build Status](https://travis-ci.org/divtxt/lockd.svg?branch=master)](https://travis-ci.org/divtxt/lockd)

A distributed lock service.

See [API.md](API.md) for server API documentation.

## Development Quick Start

Run using:

```
go run main.go
```

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


## Development Notes

### gRPC & generated code

*lockd* uses [gRPC](http://www.grpc.io/) for cross-server communication. The generated code is checked into source control so you do not have to generate the code yourself unless you're changing the interface or the code generation tooling. The command used to generate the code is:

```
protoc -I lockapi/ lockapi/locking.proto --go_out=plugins=grpc:lockapi
```

### Python Development

These notes use *pip*, *virtualenv*, and *pipenv*. Feel free to adapt/modify for your preferred tools/style.

- Install *pip*, *virtualenv* and *pipenv* (one-time global setup)

```
sudo easy_install pip
sudo pip install virtualenv
sudo pip install pipenv --ignore-installed six
```

- Install required packages

```
cd pyclient
pipenv install
```

- Run sample python client

```
pipenv run python testclient.py
```

- To generate the Python grpc code:

```
pipenv run python -m grpc_tools.protoc -I../lockapi --python_out=. --grpc_python_out=. ../lockapi/locking.proto
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

- [ ] cluster config file & config file loading
- [ ] command line param: cluster config file name
- [ ] raft consensus module
- [ ] raft rpc endpoints
- [ ] peer rpc service
- [ ] proxy Lock & Unlock API calls to raft leader

Clients:

- [ ] Ruby client
- [ ] Python client
- [ ] Java/Scala client

Servers:

- [ ] Homebrew installable binary
- [ ] Linux binary

Misc:

- [ ] Add metrics & logging
- [ ] Expose raft details e.g. leader, term
- [x] Stop `-help` from showing "-httptest.serve"
- [ ] Vendorize dependencies
