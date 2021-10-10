from typing import List
from libs.utils import Utils
from libs.task.basetask import BaseTask

class TaskMongo:
    _id: str
    pid: int
    command: str
    packetId: str
    outputPath: str
    status: str

    def __init__(self, task: BaseTask):
        self.pid = pid
        self.command= command
        self.packetId= packetId
        self.outputPath= outputPath
        self.status=status
    
    def save(self):
        pass

