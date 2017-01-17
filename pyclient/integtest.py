#!/usr/bin/python

from lockd import LockdClient

lockd_client = LockdClient()

# Lock
assert lockd_client.lock("foo")

# Dup lock should fail
assert not lockd_client.lock("foo")

# Lock another entry should work
assert lockd_client.lock("bar")

# Unlock entries
assert lockd_client.unlock("foo")
assert lockd_client.unlock("bar")

# Dup unlock should fail
assert not lockd_client.unlock("bar")
