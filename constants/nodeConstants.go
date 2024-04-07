package constants

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"gitlab.com/avarf/getenvs"
)

var (
	NodeName        = getenvs.GetEnvString("NODE_NAME", generatePeerId())
	NodeBindAddress = getenvs.GetEnvString("NODE_BIND_ADDRESS", "0.0.0.0")
	NodeBindPort, _ = getenvs.GetEnvInt("NODE_BIND_PORT", 9042)

	EnableNodeJoinRequest, _ = getenvs.GetEnvBool("ENABLE_NODE_JOIN_REQUEST", false)
	NodeJoinAddress          = getenvs.GetEnvString("NODE_JOIN_ADDRESS", "172.20.0.5:7041")
)

const BlockTopic string = "BLOCK"

func generatePeerId() string {
	// Generate a random byte slice
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// Hash the random bytes using SHA-256
	hash := sha256.New()
	hash.Write(randomBytes)
	hashedBytes := hash.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
