# lockd client

import grpc
import locking_pb2
import locking_pb2_grpc

class LockdClient(object):

    def __init__(self, host="127.0.0.1", port=2080):
        host_port = "%s:%s" % (host, port)
        channel = grpc.insecure_channel(host_port)
        self._stub = locking_pb2_grpc.LockingStub(channel)

    def IsLocked(self, name):
        response = self._stub.IsLocked(locking_pb2.LockName(name=name))
        return response.is_locked

    def Lock(self, name):
        response = self._stub.Lock(locking_pb2.LockName(name=name))
        return response.success

    def Unlock(self, name):
        response = self._stub.Unlock(locking_pb2.LockName(name=name))
        return response.success
