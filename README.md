
# lockd

A distributed lock service.


## Development

Run using:

    go run main.go

Then, run the following commands in another terminal:

Lock an entry:

```
curl -i -H "Content-Type: application/json" -X POST \
    -d "{'name':'Foo', 'ttl':30}" \
    http://127.0.0.1:8080/lock
```

Try to lock it again - you should get an error (409 Conflict ?).

To check if an entry is locked:

```
curl -i http://127.0.0.1:8080/lock?name=Foo
```

To unlock the entry:

```
curl -i -H "Content-Type: application/json" -X POST \
    -d "{'name':'Foo'}" \
    http://127.0.0.1:8080/unlock
```
