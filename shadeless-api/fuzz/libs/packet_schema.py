# import pymongo

from typing import List
from libs.utils import Utils

class Packet:
    _id: str
    requestPacketId: str
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
    rawReq: bytes
    tmpPath: str
    URL: str

    def getTmpPath(self):
        if not self.tmpPath:
            self.tmpPath = Utils.createTmpRaw(self.rawReq, self.requestPacketId)
        return self.tmpPath

    def constructRawRequest(self):
        self.rawReq = b"\r\n".join([bytes(i, encoding="utf8") for i in self.requestHeaders])
        self.rawReq += b"\r\n"*2 + self.requestBody

    def __init__(self, packet: any):
        for key in packet:
            setattr(self, key, packet[key])
        self.requestBody = open("../files/%s/%s" % (self.project, self.requestBodyHash), "rb").read()
        self.constructRawRequest()
        self.tmpPath = ""
        self.URL = self.origin

    def getHeaders(self) -> dict:
        result = {}
        for i in range(1, len(self.requestHeaders)):
            h = self.requestHeaders[i].split(": ")
            key, value = h[0], h[1]
            result[key] = value
        return result

from hashlib import md5

class PacketFile(Packet):
    def __init__(self, URL, rawReq, rawResp=""):
        self.requestPacketId = md5(rawReq).hexdigest()
        self.URL = URL
        self.rawReq = rawReq
        self.rawResp = rawResp
        self.tmpPath = ""
