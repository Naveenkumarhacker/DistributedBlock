package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

func readPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFile)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func rsaDecryptAESKey(ciphertext string, privateKey *rsa.PrivateKey) ([]byte, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedCiphertext)
	if err != nil {
		return nil, err
	}

	return aesKey, nil
}

func aesDecrypt(ciphertext string, key []byte) (string, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext length is not a multiple of the block size")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
}

func Decrypt(encryptedData string, encryptedAESKey string) string {
	// Decrypt AES key with RSA private key
	decryptedAESKey, err := rsaDecryptAESKey(encryptedAESKey, privateKey)
	if err != nil {
		fmt.Println("Error decrypting AES key:", err)
		return ""
	}

	// Decrypt message with AES key
	decryptedMessage, err := aesDecrypt(encryptedData, decryptedAESKey)
	if err != nil {
		fmt.Println("Error decrypting message:", err)
		return ""
	}
	decryptedMessage = strings.ReplaceAll(decryptedMessage, "\r", "")
	decryptedMessage = strings.ReplaceAll(decryptedMessage, "\v", "")

	return decryptedMessage
}
