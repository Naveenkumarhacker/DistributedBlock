package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
)

func readPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	keyFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFile)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to parse RSA public key")
	}

	return rsaPublicKey, nil
}

func rsaEncryptAESKey(aesKey []byte, publicKey *rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func aesEncrypt(message string, key []byte) (string, error) {
	plaintext := []byte(message)

	// Pad the plaintext to make its length a multiple of the block size
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padText...)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Encrypt(plainText string) (string, string) {
	// Generate AES key
	aesKey, err := generateAESKey()
	if err != nil {
		fmt.Println("Error generating AES key:", err)
		return "", ""
	}

	// Encrypt AES key with RSA public key
	encryptedAESKey, err := rsaEncryptAESKey(aesKey, publicKey)
	if err != nil {
		fmt.Println("Error encrypting AES key:", err)
		return "", ""
	}

	// Encrypt message with AES key
	encryptedMessage, err := aesEncrypt(plainText, aesKey)
	if err != nil {
		fmt.Println("Error encrypting message:", err)
		return "", ""
	}

	return encryptedMessage, encryptedAESKey

}
