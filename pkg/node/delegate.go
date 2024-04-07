package node

import (
	"DistributedBlock/pkg/node/models"
	"github.com/hashicorp/memberlist"
)

func NewDelegate(e *EventDelegate, metadata models.MetaData) *Delegate {
	msgCh := make(chan []byte)
	broadcasts := new(memberlist.TransmitLimitedQueue)
	broadcasts.NumNodes = e.consistent.Size

	return &Delegate{
		MsgCh:      msgCh,
		Broadcasts: broadcasts,
		meta:       metadata,
	}
}

type Delegate struct {
	MsgCh      chan []byte
	meta       models.MetaData
	Broadcasts *memberlist.TransmitLimitedQueue
}

func (d *Delegate) NotifyMsg(msg []byte) {
	d.MsgCh <- msg
}
func (d *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return d.Broadcasts.GetBroadcasts(overhead, limit)
}
func (d *Delegate) NodeMeta(limit int) []byte {
	return d.meta.Bytes()
}
func (d *Delegate) LocalState(join bool) []byte {
	// not use, noop
	return []byte("")
}
func (d *Delegate) MergeRemoteState(buf []byte, join bool) {
	// not use
}
