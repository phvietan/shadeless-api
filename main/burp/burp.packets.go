package burp

import (
	"errors"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/responser"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Parse packetID with following format: b20a4f41-3e31-48ee-8cde-8f5c81adb755.142
// Into b20a4f41-3e31-48ee-8cde-8f5c81adb755 and 142
func decorateInputPacket(p *database.Packet) (*database.Packet, error) {
	if p == nil {
		return nil, errors.New("Burp packet is nil")
	}
	arr := strings.Split(p.RequestPacketId, ".")
	if len(arr) != 2 || len(arr[0]) != 36 {
		return nil, errors.New("Request packet id format is wrong")
	}
	var err error
	if p.RequestPacketIndex, err = strconv.Atoi(arr[1]); err != nil {
		return nil, errors.New("Request packet id format is wrong")
	}
	p.RequestPacketPrefix = arr[0]
	return p, nil
}

func postPackets(c *gin.Context) {
	inputPacket := new(database.Packet)
	if err := c.BindJSON(inputPacket); err != nil {
		responser.ResponseError(c, err)
		return
	}

	packet, err := decorateInputPacket(inputPacket)
	if err != nil {
		responser.ResponseError(c, err)
		return
	}

	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	if err := packetDb.CreatePacket(packet); err != nil {
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully create packet")
}
