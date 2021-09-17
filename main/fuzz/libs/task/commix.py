from libs.task.basetask import BaseTask

class Commix(BaseTask):
    detectPattern = "injectable"
    def detect(output):
        with open(output, 'r') as f:
            if Commix.detectPattern in f.read():
                return True
        return False
        

    def __init__(self, packet, cmd="", proxy = None):
        super().__init__(packet, cmd)
        self.bugType = '[cmdi][commix]'
        self.bin = "python2 /root/tools/commix/commix.py --ignore-session --batch"
        self.rawReq = " -r {rawReqPath}"
        if proxy:
            self.proxy = " --proxy %s" % proxy
        else:
            self.proxy = ""
        self.ssl = " --force-ssl"

    def getCmd(self):
        self.fullCmd = ''
        self.fullCmd += self.cmd.format(bin = self.bin)

        # proxy for debug
        self.fullCmd += self.proxy

        # create request tmp file, add path to command
        self.fullCmd += self.rawReq.format(rawReqPath = self.packet.getTmpPath())

        if "https" in self.packet.URL:
            self.fullCmd += " --force-ssl"

        return self.fullCmd




