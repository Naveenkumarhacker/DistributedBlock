package constants

import (
	"github.com/go-gorp/gorp"
	"gitlab.com/avarf/getenvs"
)

var DbMap *gorp.DbMap

var (
	Port                = getenvs.GetEnvString("PORT", "8080")
	BlockServicePort    = getenvs.GetEnvString("PORT", "50051")
	PrivateKeyFile      = getenvs.GetEnvString("PRIVATE_KEY_FILE", "private.pem")
	PublicKeyFile       = getenvs.GetEnvString("PUBLIC_KEY_FILE", "public.pem")
	BlockGRPCServiceURL = getenvs.GetEnvString("BLOCK_GRPC_SERVICE_URL", "localhost:50051")
)

type MessageType int

const (
	Insert MessageType = iota
	Update
)
