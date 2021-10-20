from time import sleep, time
from libs.mongo import mongo

def work():
    packets = mongo.getPackets(100)
    for p in packets:
        if len(p.requestBody) != 0:
            print(p._id)
            print(p.requestBody)
            exit()

if __name__ == '__main__':
    before = time()
    while True:
        now = time()
        if now - before > 2: # Increase this in production
            work()
            exit()
        sleep(1) # Increase when run production
