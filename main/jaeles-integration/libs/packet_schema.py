# import pymongo

from typing import List

class Packet:
    _id: str
    origin: str
    port: int
    path: str
    method: str
    fuzzed: bool
    querystring: str
    parameters: List[str]
    project: str
    requestBody: bytes
    requestBodyHash: str
    requestHeaders: List[str]

    def __init__(self, packet: any):
        for key in packet:
            setattr(self, key, packet[key])
        self.requestBody = open("../files/%s/%s" % (self.project, self.requestBodyHash), "rb").read()

    def getHeaders(self) -> dict:
        result = {}
        for i in range(1, len(self.requestHeaders)):
            h = self.requestHeaders[i].split(": ")
            key, value = h[0], h[1]
            result[key] = value
        return result
