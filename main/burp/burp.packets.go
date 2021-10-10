package burp

import (
	"errors"
	"fmt"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/responser"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Parse packetID with following format: b20a4f41-3e31-48ee-8cde-8f5c81adb755.142
// Into b20a4f41-3e31-48ee-8cde-8f5c81adb755 and 142
func parsePacketIndex(p *schema.Packet) (*schema.Packet, error) {
	arr := strings.Split(p.RequestPacketId, ".")
	if len(arr) != 2 {
		return nil, errors.New("Request packet index format is wrong")
	}
	var err error
	if p.RequestPacketIndex, err = strconv.Atoi(arr[1]); err != nil {
		return nil, errors.New("Request packet index format is wrong")
	}
	p.RequestPacketPrefix = arr[0]
	sort.Strings(p.Parameters)
	return p, nil
}

func decorateAndParseInputPacket(inputPacket *schema.Packet) (*schema.Packet, *schema.ParsedPacket, error) {
	if inputPacket == nil {
		return nil, nil, errors.New("Burp packet is nil")
	}
	packet, err := parsePacketIndex(inputPacket)
	if err != nil {
		return nil, nil, err
	}
	parsedPacket, err := new(schema.ParsedPacket).ParseFromPacket(packet)
	if err != nil {
		return nil, nil, err
	}
	return packet, parsedPacket, nil
}

func insertToDb(packet *schema.Packet, parsedPacket *schema.ParsedPacket) error {
	// Normal packet
	var packetDb database.IPacketDatabase = new(database.PacketDatabase).Init()
	if found := packetDb.GetPacketByPacketId(packet.Project, packet.RequestPacketId); found != nil {
		return errors.New("Packet with this packet id already existed")
	}
	if err := packetDb.Insert(packet); err != nil {
		return err
	}
	var errUser, errParsedPacket, errParsedPath error
	var userDb database.IUserDatabase = new(database.UserDatabase).Init()
	errUser = userDb.Upsert(packet.Project, packet.CodeName)

	// Parsed packet
	var parsedPacketDb database.IParsedPacketDatabase = new(database.ParsedPacketDatabase).Init()
	if errParsedPacket = parsedPacketDb.Upsert(parsedPacket); errParsedPacket != nil {
		fmt.Println("Error: ", errParsedPacket)
	}

	// Parsed path
	parsedPaths, errParsedPath := schema.GetPathsFromParsedPacket(parsedPacket)
	if errParsedPath != nil {
		fmt.Println("Error: ", errParsedPath)
	}
	var parsedPathDb database.IParsedPathDatabase = new(database.ParsedPathDatabase).Init()
	for _, p := range parsedPaths {
		if errParsedPath = parsedPathDb.Upsert(&p); errParsedPath != nil {
			fmt.Println("Error: ", errParsedPath)
		}
	}
	if errUser != nil {
		return errUser
	}
	if errParsedPacket != nil {
		return errParsedPacket
	}
	return errParsedPath
}

func postPackets(c *gin.Context) {
	inputPacket := new(schema.Packet)
	if err := c.BindJSON(inputPacket); err != nil {
		fmt.Println("Cannot bind json input packet: ", err)
		responser.ResponseError(c, err)
		return
	}
	inputPacket.RequestPacketId = strings.ToLower(inputPacket.RequestPacketId)

	packet, parsedPacket, err := decorateAndParseInputPacket(inputPacket)
	if err != nil {
		fmt.Println("Cannot decorate parse input: ", err)
		responser.ResponseError(c, err)
		return
	}

	if err := insertToDb(packet, parsedPacket); err != nil {
		fmt.Println("Cannot insert to DB: ", err)
		responser.ResponseError(c, err)
		return
	}
	responser.ResponseOk(c, "Successfully create packet")
}
