# lockd client

import httplib
import json

class LockdClient(object):

    def __init__(self, host="127.0.0.1", port=2080):
        self._host_port = "%s:%s" % (host, port)

    def lock(self, name):
        return self._lockish(True, name)

    def unlock(self, name):
        return self._lockish(False, name)

    def _lockish(self, action, name):
        path = "/api/lock" if action else "/api/unlock"
        headers = {"Content-type": "application/json"}
        data = {'name': str(name)}
        datajson = json.dumps(data)
        #
        conn = httplib.HTTPConnection(self._host_port)
        conn.request("POST", path, datajson, headers)
        response = conn.getresponse()
        status = response.status
        response.read()
        conn.close()
        #
        if status == 200:
            return True
        elif status == 409:
            return False
        else:
            msg = "Unexpected response: %s %s; data: %s" % (status, response.reason, data)
            raise Exception(msg)
