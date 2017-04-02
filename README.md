
# lockd

[![Build Status](https://travis-ci.org/divtxt/lockd.svg?branch=master)](https://travis-ci.org/divtxt/lockd)

A distributed lock service.

See [API.md](API.md) for server API documentation.

## Development

### Run the server

- (One-time global setup:) install and configure Go

- Clone the repository

- Build and run the server:

```
go get -t ./...
go run main.go
```

### Test using the python client

In another terminal, setup python dependencies and run the python client:

- (One-time global setup:) Install *pip*, *virtualenv* and *pipenv* (or use your preferred style)

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

- Try the sample python client:

```
~/wgo/src/github.com/divtxt/lockd/pyclient$ pipenv run python
Python 2.7.10 (default, Oct 23 2015, 19:19:21)
[GCC 4.2.1 Compatible Apple LLVM 7.0.0 (clang-700.0.59.5)] on darwin
Type "help", "copyright", "credits" or "license" for more information.
>>> from lockd import LockdClient
>>> lockd_client = LockdClient()
>>> lockd_client.Lock("foo")
True
>>> lockd_client.IsLocked("foo")
True
>>> lockd_client.Unlock("foo")
True
```

### gRPC & generated code

*lockd* uses [gRPC](http://www.grpc.io/) for cross-server communication. The generated code is checked into source control so you do not have to generate the code yourself unless you're changing the interface or the code generation tooling.

If you want to change the code and regenerate code, do the following:

- Install [Protocol Buffers](https://developers.google.com/protocol-buffers/) - e.g:

```
brew install protobuf
```

- Generate the Go grpc code:

```
protoc -I lockapi/ lockapi/lockapi.proto --go_out=plugins=grpc:lockapi
```

- Generate the Python grpc code:

```
cd pyclient
pipenv run python -m grpc_tools.protoc -I../lockapi --python_out=. --grpc_python_out=. ../lockapi/lockapi.proto
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

- [ ] API timeouts
- [ ] Field content and size limit checks

Integration Testing:

- [x] External test script
- [x] Travis runs server and test script
- [ ] Travis build checks codegen

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

- [ ] Go client in separate repo
- [ ] Published Python client
- [ ] Published Ruby client

Servers:

- [ ] Homebrew package
- [ ] Ubuntu/Debian package

Misc:

- [ ] Add metrics & logging
- [ ] Expose raft details e.g. leader, term
- [x] Stop `-help` from showing "-httptest.serve"
- [ ] Vendorize dependencies
