import psutil
from multiprocessing import Process
from time import sleep

class Utils:
    def getValue(object, key, default=""):
        try:
            return object[key]
        except KeyError:
            return default
    
    def createTmpRaw(rawReq, requestPacketId):
        path = '/tmp/%s' % requestPacketId
        with open(path, 'wb') as f:
            f.write(rawReq)
        return path
    
    def findProc(cmd):
        print(cmd)
        for proc in psutil.process_iter():
            if cmd == " ".join(proc.cmdline()):
                return proc.pid
        return -1
    
    def spawn_process(function, args):
        p = Process(target=function, args = args)
        p.start()
        sleep(0.5)

