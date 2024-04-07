package eventHandler

import (
	"DistributedBlock/constants"
	"DistributedBlock/dao"
	"DistributedBlock/pkg/crypto"
	"DistributedBlock/pkg/node/models"
	"encoding/json"
	"github.com/hashicorp/memberlist"
	"github.com/withmandala/go-log"
	"os"
	"strings"
)

var logger = log.New(os.Stderr)

func HandlerNodeJoinAndLeaveEventCallBack(node *memberlist.Node, alive bool) {
	// TODO: if alive true then new node join or else node removed
	// TODO: Handle Onboarding logic here
	if node.Name == constants.NodeName {
		logger.Infof("Node  with Peer Id : %s is %s\n", node.Name, map[bool]string{true: "started", false: "dead"}[alive])
	} else {
		logger.Infof("Node with Peer Id : %s is %s\n", node.Name, map[bool]string{true: "joined", false: "left"}[alive])
	}

}

func HandlerNodeMessageCallBack(sender string, message *models.NodeMessage) {
	//TODO : Handle Received Message

	if message.Topic == constants.BlockTopic {

		decrypt := crypto.Decrypt(message.EncData, message.EncAESKey)
		decrypt = strings.Trim(decrypt, string('\x01'))
		var block dao.Block
		err := json.Unmarshal([]byte(decrypt), &block)
		if err != nil {
			logger.Error(err)
		}

		switch message.Category {
		case constants.Insert:
			dao.CreateBlock(constants.DbMap, &block)
		case constants.Update:
			dao.UpdateBlock(constants.DbMap, &block)
		}
	}

	logger.Infof("Received message from %s: %+v\n", sender, message)

}

func HandleMetaDataCallBack(node *memberlist.Node) {
	logger.Info("Received metadata from %s: %s\n", node.Name, node.Meta)
}
