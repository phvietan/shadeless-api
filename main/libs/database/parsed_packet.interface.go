package database

import "shadeless-api/main/libs/finder"

type IParsedPacketDatabase interface {
	IDatabase

	Init() *ParsedPacketDatabase
	Upsert(packet *ParsedPacket) error
	GetMetadataByProject(project *Project) ([]string, []string, map[string]string)
	GetPacketsByOriginAndProject(projectName string, origin string, options *finder.FinderOptions) []ParsedPacket
}
