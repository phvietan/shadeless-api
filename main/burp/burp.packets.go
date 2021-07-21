package burp

import (
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func postPackets(c *gin.Context) {
	packet := new(database.Packet)
	c.BindJSON(packet)

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	err := packetDb.CreatePacket(packet)
	if err != nil {
		responser.ResponseJson(c, 500, "", "Cannot create packet: "+err.Error())
		return
	}

	responser.ResponseJson(c, 200, "Successfully create packet", "")
}
