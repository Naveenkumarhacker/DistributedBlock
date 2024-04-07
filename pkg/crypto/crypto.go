package crypto

import (
	"DistributedBlock/constants"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func Init() {

	privateKeyFile := constants.PrivateKeyFile
	publicKeyFile := constants.PublicKeyFile

	var err error
	// Read private key
	privateKey, err = readPrivateKeyFromFile(privateKeyFile)
	if err != nil {
		fmt.Println("Error reading private key:", err)
		return
	}

	// Read public key
	publicKey, err = readPublicKeyFromFile(publicKeyFile)
	if err != nil {
		fmt.Println("Error reading public key:", err)
		return
	}
}

func generateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
