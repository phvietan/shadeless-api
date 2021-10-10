import requests
from packet_schema import Packet

class JaeLesSender:
    packet: Packet
    def __init__(self, packet: Packet):
        self.jaelesIP = "localhost:5000"
        self.packet = packet

    def send(self):
        url = self.packet.origin + self.packet.path + self.packet.querystring
        requests.request(self.packet.method, url, headers=self.packet.getHeaders())
