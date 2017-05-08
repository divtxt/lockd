# lockd client

import httplib
import json

class LockdClient(object):

    def __init__(self, host="127.0.0.1", port=2081):
        self._host_port = "%s:%s" % (host, port)

    def is_locked(self, name):
        return self._lockish("GET", name, 404)

    def lock(self, name):
        return self._lockish("POST", name, 409)

    def unlock(self, name):
        return self._lockish("DELETE", name, 404)

    def _lockish(self, method, name, false_code):
        # FIXME: does name need to escaped here?
        path = "/lock/%s" % name
        #
        conn = httplib.HTTPConnection(self._host_port)
        conn.request(method, path)
        response = conn.getresponse()
        status = response.status
        response.read()
        conn.close()
        #
        if status == 200:
            return True
        elif status == false_code:
            return False
        else:
            msg = "Unexpected response: %s %s" % (status, response.reason)
            raise Exception(msg)
