#!/usr/bin/python

from lockd import LockdClient

lockd_client = LockdClient()

print "Lock('foo') ->", lockd_client.Lock("foo")

print "Lock('foo') ->", lockd_client.Lock("foo")

print "IsLocked('foo') ->", lockd_client.IsLocked("foo")

print "Unlock('foo') ->", lockd_client.Unlock("foo")
