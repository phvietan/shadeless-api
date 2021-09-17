from libs.mongo import mongo
import psutil
from time import sleep
from libs.task.commix import Commix
from libs.task.basetask import BaseTask


class WatchDog:
    def __init__(self):
        pass

    def start(self):
        while True:
            pids = psutil.pids()
            for task in mongo.getRunningTask():
                detect = getattr(globals()[task["class"]], "detect")
                if detect(task["result"]):
                    p = psutil.Process(task["pid"])
                    p.terminate()
                    print("done %s" % task["pid"])
                    mongo.updateTaskStatus(task["_id"], "done")

                if task["pid"] not in pids:
                    mongo.updateTaskStatus(task["_id"], "done")


            sleep(5)

wd = WatchDog()
wd.start()