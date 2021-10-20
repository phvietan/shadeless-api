from flask import Flask
from flask import request
from flask import jsonify
from libs.utils import Utils
from base64 import b64decode 
from libs.mongo import mongo
from libs.task.basetask import BaseTask
from libs.task.commix import Commix
from libs.packet_schema import PacketFile
from config import Config

app = Flask(__name__)

def default_fuzz(packet):
    # example
    # t1 = BaseTask(packet)
    # t1.run()

    t2 = Commix(packet, proxy = Config.proxy)
    t2.run()

@app.route("/api/fuzz/requestPacketId", methods=["POST"])
def fuzz_requestPacketId():
    json_data = request.json
    packetId = json_data["requestPacketId"]
    packet = mongo.getPacket(packetId)

    # call other tools with os.system
    default_fuzz(packet)

    # just return
    return jsonify(
        message="ok",
        packetId=packet.requestPacketId
    )

@app.route("/api/fuzz/raw", methods=["POST"])
def fuzz_raw():
    json_data = request.json
    url = json_data["url"]
    rawReq = b64decode(Utils.getValue(json_data, "rawReq"))
    rawResp = b64decode(Utils.getValue(json_data, "rawResp"))
    packet = PacketFile(url, rawReq, rawResp)

    # call other tools with os.system
    default_fuzz(packet)

    # just return
    return jsonify(
        message="ok",
        packetId=packet.requestPacketId
    )

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=3108, debug=True)