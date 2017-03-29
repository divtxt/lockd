# lockd client

import grpc
import lockapi_pb2
import lockapi_pb2_grpc

class LockdClient(object):

    def __init__(self, host="127.0.0.1", port=2080):
        host_port = "%s:%s" % (host, port)
        channel = grpc.insecure_channel(host_port)
        self._stub = lockapi_pb2_grpc.LockingStub(channel)

    def IsLocked(self, name):
        response = self._stub.IsLocked(lockapi_pb2.LockName(name=name))
        return response.is_locked

    def Lock(self, name):
        response = self._stub.Lock(lockapi_pb2.LockName(name=name))
        return response.success

    def Unlock(self, name):
        response = self._stub.Unlock(lockapi_pb2.LockName(name=name))
        return response.success
