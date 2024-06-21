package client

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

/*
 * Simple encryption and decryption of passwords
 * Replace SecretString (32 char) and bytes (16 char) for release build
 */
const SecretString string = "icouldcrysaltytearswherehaveIbee"

var bytes = []byte("howlonghasbeengo")

func cryptPassword(_ *SscClient, args *Arguments) error {
	encrypted, err := Encrypt(args.Password)
	if err != nil {
		return fmt.Errorf("could not encrypt password %v", err)
	}
	fmt.Printf("%s\n", encrypted)
	return nil
}

// test only. Do not include in a production build
func decryptPassword(_ *SscClient, args *Arguments) error {
	plain, err := Decrypt(args.Password)
	if err != nil {
		return fmt.Errorf("could not decrypt password %v", err)
	}
	fmt.Printf("%s\n", plain)
	return nil
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("cound not decode base64 %s to string %v", s, err)
	}
	return data, nil
}

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(SecretString))
	if err != nil {
		return "", fmt.Errorf("could not create cipher to encrypt %v", err)
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(SecretString))
	if err != nil {
		return "", fmt.Errorf("could not create cipher to decrypt  %v", err)
	}
	cipherText, err := Decode(text)
	if err != nil {
		return "", fmt.Errorf("could not decode %s %v", text, err)
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
