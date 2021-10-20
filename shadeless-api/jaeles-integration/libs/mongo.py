from typing import List
import pymongo
from .packet_schema import Packet

class MongoHandler:
    def __init__(self):
        self.conn = pymongo.MongoClient("mongodb://localhost:27017/")
        self.db = self.conn["shadeless"]
        # collections
        self.colFiles = self.db["files"]
        self.colPackets = self.db["packets"]
        self.colProjects = self.db["projects"]

     # Limit about 100, 200?
    def getPackets(self, limit: int) -> List[Packet]:
        packets = self.colPackets.find({ "fuzzed": False }).sort("created_at", -1).limit(limit)
        result: List[Packet] = []
        for v in packets: result.append(Packet(v))
        return result

    def isPacketFuzzed():
        return False

mongo = MongoHandler()
