package database

type IPacketDatabase interface {
	IDatabase

	Init() *PacketDatabase
	GetPacketByPacketId(projectName string, packetId string) *Packet
	GetPacketsAsTimeTravel(projectName string, packetPrefix string, packetIndex int, number int) []Packet
}
