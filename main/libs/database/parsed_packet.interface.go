package database

import (
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/finder"
)

type IParsedPacketDatabase interface {
	IDatabase

	Init() *ParsedPacketDatabase
	Upsert(packet *schema.ParsedPacket) error
	GetMetadataByProject(project *schema.Project) ([]string, []string, map[string]string)
	GetPacketsByOriginAndProject(projectName string, origin string, options *finder.FinderOptions) []schema.ParsedPacket
	GetParsedByRawPackets(project string, packets []schema.Packet) []schema.ParsedPacket
}
