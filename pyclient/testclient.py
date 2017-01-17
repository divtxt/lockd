#!/usr/bin/python

from lockd import LockdClient

lockd_client = LockdClient()

print "lock('foo') ->", lockd_client.lock("foo")

print "lock('foo') ->", lockd_client.lock("foo")

print "unlock('foo') ->", lockd_client.unlock("foo")
