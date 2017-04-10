#!/usr/bin/python

from lockd import LockdClient

lockd_client = LockdClient()

# Initial state
assert not lockd_client.is_locked("foo")
assert not lockd_client.is_locked("bar")

# Lock
assert lockd_client.lock("foo")
assert lockd_client.is_locked("foo")

# Dup lock should fail
assert not lockd_client.lock("foo")
assert lockd_client.is_locked("foo")

# Lock another entry should work
assert lockd_client.lock("bar")
assert lockd_client.is_locked("bar")

# Unlock entries
assert lockd_client.unlock("foo")
assert lockd_client.unlock("bar")
assert not lockd_client.is_locked("foo")
assert not lockd_client.is_locked("bar")

# Dup unlock should fail
assert not lockd_client.unlock("bar")
assert not lockd_client.is_locked("bar")
