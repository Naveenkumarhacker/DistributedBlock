package node

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/serialx/hashring"
)

func NewEventDelegate(nodeEventCallback NodeEventCallback, metaDataCallBack MetaDataCallback) *EventDelegate {
	return &EventDelegate{eventCB: nodeEventCallback, metaDataCB: metaDataCallBack}
}

type EventDelegate struct {
	consistent *hashring.HashRing
	eventCB    NodeEventCallback
	metaDataCB MetaDataCallback
}

func (d *EventDelegate) NotifyJoin(node *memberlist.Node) {

	d.eventCB(node, true)

	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	if d.consistent == nil {
		d.consistent = hashring.New([]string{hostPort})
	} else {
		d.consistent = d.consistent.AddNode(hostPort)
	}
}
func (d *EventDelegate) NotifyLeave(node *memberlist.Node) {
	d.eventCB(node, false)
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	if d.consistent != nil {
		d.consistent = d.consistent.RemoveNode(hostPort)
	}
}
func (d *EventDelegate) NotifyUpdate(node *memberlist.Node) {
	d.metaDataCB(node)
}
