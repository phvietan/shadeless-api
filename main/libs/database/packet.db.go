package database

import (
	"errors"
	"fmt"
	"log"
	"shadeless-api/main/libs/finder"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Packet struct {
	mgm.DefaultModel `bson:",inline"`

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

	ResponseStatus           int      `json:"responseStatus" bson:"responseStatus"`
	ResponseContentType      string   `json:"responseContentType" bson:"responseContentType"`
	ResponseStatusText       string   `json:"responseStatusText" bson:"responseStatusText"`
	ResponseLength           int      `json:"responseLength" bson:"responseLength"`
	ResponseMimeType         string   `json:"responseMimeType" bson:"responseMimeType"`
	ResponseHttpVersion      int      `json:"responseHttpVersion" bson:"responseHttpVersion"`
	ResponseInferredMimeType string   `json:"responseInferredMimeType" bson:"responseInferredMimeType"`
	ResponseCookies          string   `json:"responseCookies" bson:"responseCookies"`
	ResponseBodyHash         string   `json:"responseBodyHash" bson:"responseBodyHash"`
	ResponseHeaders          []string `json:"responseHeaders" bson:"responseHeaders"`
	Rtt                      int      `json:"rtt"`
	ReflectedParameters      []string `json:"reflectedParameters" bson:"reflectedParameters"`
	Project                  string   `json:"project"`
	CodeName                 string   `json:"codeName"`
}

// TODO: Change this CRUD into interface for polymorphism
func CreatePacket(packet *Packet) error {
	if packet == nil {
		return errors.New("Packet object is nil")
	}
	err := mgm.Coll(packet).Create(packet)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func arrayInterfaceToArrayString(arr []interface{}) []string {
	result := make([]string, len(arr))
	for i, v := range arr {
		result[i] = fmt.Sprint(v)
	}
	return result
}

func getDistinc(name string, filterOptions bson.M) []string {
	ctx := mgm.Ctx()

	coll := mgm.Coll(&Packet{})
	results, err := coll.Distinct(ctx, name, filterOptions)
	if err != nil {
		fmt.Errorf("%v", err)
		return []string{}
	}
	return arrayInterfaceToArrayString(results)
}

func GetOrigins(projectName string) []string {
	filterOptions := bson.M{
		"project": projectName,
	}
	return getDistinc("origin", filterOptions)
}

func GetParameters(projectName string) []string {
	filterOptions := bson.M{
		"project": projectName,
	}
	return getDistinc("parameters", filterOptions)
}

func GetReflectedParameters(projectName string) []string {
	filterOptions := bson.M{
		"project": projectName,
	}
	return getDistinc("reflectedParameters", filterOptions)
}

type numPackets struct {
	result struct {
		Length int `json:"$numberInt" bson:"$numberInt"`
	}
}

func GetNumPacketsByOrigin(projectName string, origin string) int32 {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"origin": origin, "project": projectName}},
		bson.M{"$group": bson.M{"_id": "$path"}},
		bson.M{"$count": "result"},
	}

	ctx := mgm.Ctx()
	cursor, err := mgm.Coll(&Packet{}).Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var allDocs []bson.M
	if err = cursor.All(ctx, &allDocs); err != nil {
		log.Fatal(err)
		return 0
	}
	return allDocs[0]["result"].(int32)
}

func GetPacketsByOriginAndProject(projectName string, origin string, options *finder.FinderOptions) []Packet {
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"origin": origin, "project": projectName}},
		bson.M{"$group": bson.M{"_id": "$path", "doc": bson.M{"$last": "$$ROOT"}}},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$doc"}},
		bson.M{"$skip": options.Skip},
		bson.M{"$limit": options.Limit},
	}

	ctx := mgm.Ctx()
	cursor, err := mgm.Coll(&Packet{}).Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Errorf("%v", err)
		return []Packet{}
	}
	results := []Packet{}
	err = cursor.All(ctx, &results)
	if err != nil {
		fmt.Errorf("%v", err)
		return []Packet{}
	}
	return results
}
