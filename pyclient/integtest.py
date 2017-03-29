#!/usr/bin/python

from lockd import LockdClient

lockd_client = LockdClient()

# Initial state
assert not lockd_client.IsLocked("foo")
assert not lockd_client.IsLocked("bar")

# Lock
assert lockd_client.Lock("foo")
assert lockd_client.IsLocked("foo")

# Dup lock should fail
assert not lockd_client.Lock("foo")
assert lockd_client.IsLocked("foo")

# Lock another entry should work
assert lockd_client.Lock("bar")
assert lockd_client.IsLocked("bar")

# Unlock entries
assert lockd_client.Unlock("foo")
assert lockd_client.Unlock("bar")
assert not lockd_client.IsLocked("foo")
assert not lockd_client.IsLocked("bar")

# Dup unlock should fail
assert not lockd_client.Unlock("bar")
assert not lockd_client.IsLocked("bar")
