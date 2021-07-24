package database

import (
	"github.com/kamva/mgm/v3"
)

type Packet struct {
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
