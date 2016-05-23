
# lockd API

This page documents the API for the _lockd_ service.


## POST /api/lock

Lock the entry with the given name.

Parameters:

- `name` - the lock name

Returns one of the following:

- `204 No Content`: Entry successfully locked
- `409 Conflict`: Entry is already locked

Example:

```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:2080/lock \
    -d '{"name":"Foo"}'
```


## POST /api/unlock

Unlock the entry with the given name.

Parameters:

- `name` - the lock name

Returns one of the following:

- `204 No Content`: Entry successfully unlocked
- `409 Conflict`: Entry is already unlocked

Example:

```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:2080/unlock \
    -d '{"name":"Foo"}'
```

