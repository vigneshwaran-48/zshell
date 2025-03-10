package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const SEPARATER = ":=:=:=:"

func getKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

func encrypt(key []byte, token string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aesgcm.Seal(nonce, nonce, []byte(token), nil), nil
}

func decryptToken(key []byte, cipherText []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func Encrypt(password string, text string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := getKey(password, salt)
	encryptedToken, err := encrypt(key, text)
	if err != nil {
		return "", err
	}
	encryptedTokenStr := base64.StdEncoding.EncodeToString(encryptedToken)
	saltStr := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("%s%s%s", encryptedTokenStr, SEPARATER, saltStr), nil
}

func Decrypt(password string, encryptedData string) (string, error) {
	splitted := strings.Split(encryptedData, SEPARATER)
	encryptedTokenStr := splitted[0]
	saltStr := splitted[1]

	encryptedToken, err := base64.StdEncoding.DecodeString(encryptedTokenStr)
	salt, err := base64.StdEncoding.DecodeString(saltStr)

	key := getKey(password, salt)

	decryptedText, err := decryptToken(key, encryptedToken)
	if err != nil {
		return "", err
	}
	return decryptedText, nil
}
