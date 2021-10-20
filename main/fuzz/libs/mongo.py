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
        self.colTasks = self.db["tasks"]

     # Limit about 100, 200?
    def getPackets(self, limit: int) -> List[Packet]:
        packets = self.colPackets.find({ "fuzzed": False }).sort("created_at", -1).limit(limit)
        result: List[Packet] = []
        for v in packets: result.append(Packet(v))
        return result

    def getPacket(self, requestPacketId) -> Packet:
        packet = self.colPackets.find_one({ "requestPacketId": requestPacketId })
        return Packet(packet)

    def isPacketFuzzed():
        return False

    def saveTask(self, task):
        taskObject = {
            "pid": task.pid,
            "class": type(task).__name__,
            "fullCmd": task.fullCmd,
            "packetId": task.packet.requestPacketId,
            "result": task.result,
            "status": task.status
        }
        self.colTasks.insert_one(taskObject)
    
    def getRunningTask(self):
        tasks = self.colTasks.find({"status": "running"})
        return tasks

    def updateTaskStatus(self, mongoId, status):
        self.colTasks.update({'_id':mongoId},{"$set":{'status': status}})

mongo = MongoHandler()
