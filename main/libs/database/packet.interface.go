package database

import (
	"shadeless-api/main/libs/finder"
)

type IPacketDatabase interface {
	Init() *PacketDatabase
	CreatePacket(packet *Packet) error

	GetOriginsByProjectName(projectName string) []string
	GetParametersByProjectName(projectName string) []string
	GetReflectedParametersByProjectName(projectName string) []string
	GetNumPacketsByOrigin(projectName string, origin string) int32
	GetPacketByPacketId(projectName string, packetId string) *Packet
	GetPacketsAsTimeTravel(projectName string, packetPrefix string, packetIndex int, number int) []Packet
	GetPacketsByOriginAndProject(projectName string, origin string, options *finder.FinderOptions) []Packet

	DeletePacketsByProjectName(projectName string) error
}
