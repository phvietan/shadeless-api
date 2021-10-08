package schema

import (
	"errors"
	"strings"
	"time"

	"github.com/kamva/mgm/v3"
)

type ParsedPath struct {
	mgm.DefaultModel `bson:",inline"`
	RequestPacketId  string `json:"requestPacketId" bson:"requestPacketId"`
	Origin           string `json:"origin"`
	Path             string `json:"path"`
	Status           string `json:"status"`
	Project          string `json:"project"`
	Type             string `json:"type" bson:"type"`
	Force            bool   `json:"force" bson:"force"`
	Error            string `json:"error" bson:"error"`
}

const (
	PathStatusTodo     = "todo"
	PathStatusScanning = "scanning"
	PathStatusDone     = "done"
)

const (
	PathTypeFile      = "file"
	PathTypeDirectory = "directory"
)

func NewParsedPath(packet *ParsedPacket, path string, pathType string) *ParsedPath {
	res := &ParsedPath{
		Project:         packet.Project,
		Origin:          packet.Origin,
		Path:            path,
		Status:          PathStatusTodo,
		RequestPacketId: packet.RequestPacketId,
		Type:            pathType,
		Force:           false,
		Error:           "",
	}
	res.DefaultModel.CreatedAt = time.Now()
	res.DefaultModel.UpdatedAt = time.Now()
	return res
}

func GetPathsFromParsedPacket(packet *ParsedPacket) ([]ParsedPath, error) {
	if packet == nil {
		return nil, errors.New("Cannot parse: parsedPacket should not be nil")
	}
	paths := strings.Split(packet.Path, "/")
	curPath := ""
	result := []ParsedPath{}
	for idx, path := range paths {
		if idx == len(paths)-1 && path == "" {
			continue
		}
		t := PathTypeDirectory
		if idx == len(paths)-1 && packet.StaticScore > 50 {
			t = PathTypeFile
		}
		curPath += path + "/"
		result = append(result, *NewParsedPath(packet, curPath, t))
	}
	return result, nil
}
