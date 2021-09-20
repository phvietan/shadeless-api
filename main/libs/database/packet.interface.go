package database

import "shadeless-api/main/libs/database/schema"

type IPacketDatabase interface {
	IDatabase

	Init() *PacketDatabase
	GetPacketByPacketId(projectName string, packetId string) *schema.Packet
	GetPacketsAsTimeTravel(projectName string, packetPrefix string, packetIndex int, number int) []schema.Packet
}
