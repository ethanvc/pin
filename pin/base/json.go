package base

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func StructToJson(v any) []byte {
	buf, _ := json.Marshal(v)
	return buf
}

func StructToJsonStr(v any) string {
	return BytesToStr(StructToJson(v))
}

func JsonToStruct(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	return err
}

func JsonStrToStruct(s string, v any) error {
	err := json.Unmarshal(StrToBytes(s), v)
	return err
}
