package utils

import "encoding/json"

func AnyToJSON(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func JSONToAny(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
