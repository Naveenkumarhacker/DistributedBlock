package servers

import (
	"DistributedBlock/constants"
	"DistributedBlock/pkg/node"
	eventHandler "DistributedBlock/pkg/node/eventhandler"
	"DistributedBlock/pkg/node/models"
	"github.com/hashicorp/memberlist"
	"log"
)

func InitNode() *node.NodeUtils {

	config := node.Config{
		Name:     constants.NodeName,
		BindAddr: constants.NodeBindAddress,
		Port:     constants.NodeBindPort,
	}

	m, err := node.NewNode(config,
		func(node *memberlist.Node, alive bool) {
			eventHandler.HandlerNodeJoinAndLeaveEventCallBack(node, alive)
		},
		func(sender string, message *models.NodeMessage) {
			eventHandler.HandlerNodeMessageCallBack(sender, message)
		},
		func(node *memberlist.Node) {
			eventHandler.HandleMetaDataCallBack(node)
		})

	if err != nil {
		log.Fatal(err)
	}

	if constants.EnableNodeJoinRequest {
		_ = m.JoinCluster(constants.NodeJoinAddress)
	}

	return m
}
