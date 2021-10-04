package schema

import (
	"errors"
	"fmt"
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
}

const (
	PathStatusTodo     = "todo"
	PathStatusScanning = "scanning"
	PathStatusDone     = "done"
)

func NewParsedPath(packet *ParsedPacket, path string) *ParsedPath {
	res := &ParsedPath{
		Project:         packet.Project,
		Origin:          packet.Origin,
		Path:            path,
		Status:          PathStatusTodo,
		RequestPacketId: packet.RequestPacketId,
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
		if idx == len(paths)-1 && packet.StaticScore > 50 {
			continue
		}
		curPath += path + "/"
		result = append(result, *NewParsedPath(packet, curPath))
	}
	fmt.Println(result)
	return result, nil
}
