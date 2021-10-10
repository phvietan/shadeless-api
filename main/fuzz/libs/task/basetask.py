from os import system

from config import Config
from libs.utils import Utils
from libs.mongo import mongo

class BaseTask:
    detectPattern = None
    def detect():
        return False

    def __init__(self, packet, cmd=""):
        self.bugType = '[cmdi][example]'
        self.packet = packet
        self.bin = "/bin/cat"
        if cmd:
            self.cmd = cmd
        else:
            self.cmd = "{bin}"
        self.rawReq = " {rawReqPath}"

        self.result = "{outputPath}/{packetId}-{bugType}".format(
            outputPath = Config.outputPath,
            packetId = self.packet.requestPacketId,
            bugType = self.bugType
        )
        self.outputCmd = " 2>&1 | tee " + self.result
        self.status = "new"

    def getCmd(self):
        self.fullCmd = ''
        self.fullCmd += self.cmd.format(bin = self.bin)

        # create request tmp file, add path to command
        self.fullCmd += self.rawReq.format(rawReqPath = self.packet.getTmpPath())

        return self.fullCmd

    def getRealCmd(self):
        return self.getCmd() + self.outputCmd

    def run(self):
        readCmd = self.getRealCmd()
        Utils.spawn_process(system, (readCmd, ))
        self.pid = self.getPid()
        if self.pid != -1:
            self.status = "running"
        else:
            self.status = "done"
        mongo.saveTask(self)

    def getPid(self):
        return Utils.findProc(self.fullCmd)

