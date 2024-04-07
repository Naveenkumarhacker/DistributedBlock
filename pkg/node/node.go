package node

import (
	"DistributedBlock/pkg/node/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/memberlist"
	"log"
)

type Config struct {
	Name     string
	BindAddr string
	Port     int
}

// NodeEventCallback is a callback function for handling node events.
type NodeEventCallback func(node *memberlist.Node, alive bool)

// MessageCallback is a callback function for handling incoming messages.
type MessageCallback func(sender string, message *models.NodeMessage)

// MetaDataCallback is a callback function for handling incoming metadata.
type MetaDataCallback func(node *memberlist.Node)

type NodeUtils struct {
	list          *memberlist.Memberlist
	eventDelegate *EventDelegate
	delegate      *Delegate
	eventCB       NodeEventCallback
	messageCB     MessageCallback
	metaDataCB    MetaDataCallback
}

// NewNode creates a new instance of NodeUtils.
func NewNode(config Config, eventCB NodeEventCallback, messageCB MessageCallback, metaDataCB MetaDataCallback) (*NodeUtils, error) {
	memberlistConfig := memberlist.DefaultWANConfig()
	memberlistConfig.Name = config.Name
	memberlistConfig.BindAddr = config.BindAddr
	memberlistConfig.BindPort = config.Port

	eventDelegate := NewEventDelegate(eventCB, metaDataCB)
	delegate := NewDelegate(eventDelegate, models.GetMetaData())
	memberlistConfig.Events = eventDelegate
	memberlistConfig.Delegate = delegate

	m := &NodeUtils{
		list:          nil,
		delegate:      delegate,
		eventDelegate: eventDelegate,
		eventCB:       eventCB,
		messageCB:     messageCB,
		metaDataCB:    metaDataCB,
	}

	list, err := memberlist.Create(memberlistConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %v", err)
	}
	m.list = list
	go m.startEventReceiver()

	return m, nil
}

// SendMessage sends a message to a specific node in the cluster.
func (m *NodeUtils) SendMessage(targetNode string, message models.NodeMessage) error {
	target := MemberByName(m.list, targetNode)
	if target == nil {
		return fmt.Errorf("target node %s not found", targetNode)
	}

	return m.list.SendReliable(target, message.Bytes())
}

// BroadcastMessage sends a message to all nodes in the cluster.
func (m *NodeUtils) BroadcastMessage(message models.NodeMessage) {
	for _, member := range m.list.Members() {

		//if member.Name != m.list.LocalNode().Name {
		err := m.list.SendReliable(member, message.Bytes())
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", member.Name, err)
		}
		//}
	}
}

// JoinCluster joins an existing HashiCorp Memberlist cluster.
func (m *NodeUtils) JoinCluster(clusterAddr string) error {
	_, err := m.list.Join([]string{clusterAddr})
	if err != nil {
		return fmt.Errorf("failed to join cluster: %v", err)
	}
	return nil
}

// Close shuts down the Memberlist instance.
func (m *NodeUtils) Close() {
	_ = m.list.Shutdown()
}

func (m *NodeUtils) startEventReceiver() {

	stopCtx, cancel := context.WithCancel(context.TODO())
	for {
		select {
		case <-stopCtx.Done():
			m.Close()
			cancel()
			log.Printf("node shutting down...")
			return
		case data := <-m.delegate.MsgCh:
			// Handle received messages
			message, ok := ParseNodeMessage(data)
			if ok != true {
				log.Println("Error decoding message")
				continue
			}
			m.messageCB(message.Sender.NodeName, message)
		}
	}
}

func MemberByName(list *memberlist.Memberlist, targetNodeName string) *memberlist.Node {
	for _, node := range list.Members() {
		if node.Name == targetNodeName {
			return node
		}
	}
	return nil
}

func ParseNodeMessage(data []byte) (*models.NodeMessage, bool) {
	msg := new(models.NodeMessage)
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, false
	}
	return msg, true
}
