package models

import (
	"DistributedBlock/constants"
	"DistributedBlock/pkg/crypto"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/memberlist"
)

func NewNodeMessage(topic string, category constants.MessageType, data interface{}, encript bool) NodeMessage {

	var dataStr string

	switch v := data.(type) {
	case string:
		dataStr = v
	default:
		// Convert to string if not already
		jsonData, err := json.Marshal(v)
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
		}
		// Convert the JSON byte slice to a string
		dataStr = string(jsonData)
	}

	sender := Sender{NodeName: constants.NodeName}
	if encript {
		encryptedData, encryptedAESKey := crypto.Encrypt(dataStr)
		//TODO: Below 2 lines is just testing decrypt function use it later
		//dec := crypto.Decrypt(encryptedData, encryptedAESKey)
		//log.Println("dec ", dec)
		return NodeMessage{Sender: sender, Topic: topic, Category: category, EncData: encryptedData, EncAESKey: encryptedAESKey, IsDataEncrypted: encript}
	}

	return NodeMessage{Sender: sender, Topic: topic, Category: category, PlainData: dataStr, IsDataEncrypted: encript}
}

type NodeMessage struct {
	Sender          Sender                `json:"sender"`
	Topic           string                `json:"topic"`
	Category        constants.MessageType `json:"category"`
	PlainData       string                `json:"plainData"`
	EncData         string                `json:"encData"`
	EncAESKey       string                `json:"encAESKey"`
	IsDataEncrypted bool                  `json:"isDataEncrypted"`
}

type Sender struct {
	NodeName string `json:"nodeName"`
	//TODO: ADD Sender Details if required
}

func (n NodeMessage) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (n NodeMessage) Finished() {}

func (n NodeMessage) Message() []byte {
	data, err := json.Marshal(n)
	if err != nil {
		return []byte("")
	}
	return data
}

func (n *NodeMessage) Bytes() []byte {
	data, err := json.Marshal(n)
	if err != nil {
		return []byte("")
	}
	return data
}
