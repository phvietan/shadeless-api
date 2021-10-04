package schema

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/kamva/mgm/v3"
)

type ParsedPacket struct {
	mgm.DefaultModel `bson:",inline"`

	RequestPacketId    string   `json:"requestPacketId" bson:"requestPacketId"`
	ToolName           string   `json:"toolName" bson:"toolName"`
	Method             string   `json:"method"`
	RequestLength      int      `json:"requestLength" bson:"requestLength"`
	RequestHttpVersion string   `json:"requestHttpVersion" bson:"requestHttpVersion"`
	RequestContentType string   `json:"requestContentType" bson:"requestContentType"`
	Referer            string   `json:"referer"`
	Protocol           string   `json:"protocol"`
	Origin             string   `json:"origin"`
	Port               int      `json:"port"`
	Path               string   `json:"path"`
	RequestCookies     string   `json:"requestCookies" bson:"requestCookies"`
	HasBodyParam       bool     `json:"hasBodyParam" bson:"hasBodyParam"`
	Querystring        string   `json:"querystring"`
	RequestBodyHash    string   `json:"requestBodyHash" bson:"requestBodyHash"`
	Parameters         []string `json:"parameters"`
	RequestHeaders     []string `json:"requestHeaders" bson:"requestHeaders"`

	Hash        string            `json:"hash" bson:"hash"`
	Fuzzed      map[string]string `json:"fuzzed" bson:"fuzzed"`
	StaticScore float64           `json:"staticScore" bson:"staticScore"`

	ResponseStatus           int               `json:"responseStatus" bson:"responseStatus"`
	ResponseContentType      string            `json:"responseContentType" bson:"responseContentType"`
	ResponseStatusText       string            `json:"responseStatusText" bson:"responseStatusText"`
	ResponseLength           int               `json:"responseLength" bson:"responseLength"`
	ResponseMimeType         string            `json:"responseMimeType" bson:"responseMimeType"`
	ResponseHttpVersion      string            `json:"responseHttpVersion" bson:"responseHttpVersion"`
	ResponseInferredMimeType string            `json:"responseInferredMimeType" bson:"responseInferredMimeType"`
	ResponseCookies          string            `json:"responseCookies" bson:"responseCookies"`
	ResponseBodyHash         string            `json:"responseBodyHash" bson:"responseBodyHash"`
	ResponseHeaders          []string          `json:"responseHeaders" bson:"responseHeaders"`
	Rtt                      int               `json:"rtt"`
	ReflectedParameters      map[string]string `json:"reflectedParameters" bson:"reflectedParameters"`
	Project                  string            `json:"project"`
	CodeName                 string            `json:"codeName"`

	RequestPacketIndex  int    `json:"requestPacketIndex" bson:"requestPacketIndex"`
	RequestPacketPrefix string `json:"requestPacketPrefix" bson:"requestPacketPrefix"`
}

const (
	fuzzNew     = "new"
	fuzzRunning = "running"
	fuzzDone    = "done"
)

func makeDefaultFuzzMap() map[string]string {
	m := make(map[string]string)
	m["commix"] = fuzzNew
	m["jaeles"] = fuzzNew
	return m
}

func (this *ParsedPacket) ParseFromPacket(packet *Packet) (*ParsedPacket, error) {
	if packet == nil {
		return nil, errors.New("Cannot parse: Packet should not be nil")
	}
	bytesPacket, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Cannot parse input packet to JSON")
	}
	result := new(ParsedPacket)
	if err := json.Unmarshal(bytesPacket, result); err != nil {
		return nil, err
	}
	result.Fuzzed = makeDefaultFuzzMap()
	result.Hash = CalculatePacketHash(result.ResponseStatus, result.Origin, result.Path, result.Parameters)
	result.setStaticScore()
	return result, nil
}

func (this *ParsedPacket) setStaticScore() {
	score := NewStaticScorer(this).GetScore()
	this.StaticScore = score
}

func CalculatePacketHash(responseStatus int, origin string, path string, parameters []string) string {
	s := "responseStatus:" + strconv.Itoa(responseStatus) + ";origin:" + origin + ";path:" + path + "parameters:"
	for idx, val := range parameters {
		delimiter := ","
		if idx == len(parameters)-1 {
			delimiter = ";"
		}
		s += val + delimiter
	}
	b := md5.Sum([]byte(s))
	return hex.EncodeToString(b[:])
}
