package internal

import "encoding/json"

var (
	Revision string
	Version  string
)

func GetVersionJSONString() string {
	serialized, _ := json.Marshal(map[string]string{
		"version":  Version,
		"revision": Revision,
	})
	return string(serialized)
}
