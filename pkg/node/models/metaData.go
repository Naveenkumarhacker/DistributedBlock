package models

import (
	"encoding/json"
)

func GetMetaData() MetaData {
	metaData := MetaData{NodeType: "No Type"}
	return metaData
}

type MetaData struct {
	NodeType string `json:"nodeType"`
	//TODO: ADD MORE METADATA IF REQUIRED
}

func (m MetaData) Bytes() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}
func ParseMetaData(data []byte) (MetaData, bool) {
	meta := MetaData{}
	if err := json.Unmarshal(data, &meta); err != nil {
		return meta, false
	}
	return meta, true
}
