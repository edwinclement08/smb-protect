package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

func GenerateSeed() []byte {
	seed := make([]byte, 16)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}
	return seed
}

func Encrypt(secretStr, passwordStr string) ([]byte, []byte, []byte) {
	secret := []byte(secretStr)
	password := []byte(passwordStr)
	salt := GenerateSeed()

	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		log.Fatalf("encrypt: Failed to derive Key: %s", err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("encrypt: Failed to generate cipher: %s", err.Error())
	}
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatalf("encrypt: Failed to generate GCM block: %s", err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		log.Fatalf("encrypt: Failed to generate nonce: %s", err.Error())
	}
	ciphertext := gcm.Seal(nil, nonce, secret, nil)

	fmt.Printf("Ciphertext: %x\n", ciphertext)
	return ciphertext, nonce, salt
}

func Decrypt(cipherText []byte, passwordStr string, nonce, salt []byte) string {
	password := []byte(passwordStr)

	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		log.Fatalf("encrypt: Failed to derive Key: %s", err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("encrypt: Failed to generate cipher: %s", err.Error())
	}
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		log.Fatalf("encrypt: Failed to generate GCM block: %s", err.Error())
	}
	// Decrypt secret
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		panic(err)
	}
	return string(plainText)
}

func Mainenc() {
	password := []byte("mypassword")
	salt := []byte("mysalt")

	// Derive key from password and salt
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}

	// Encrypt secret
	secret := []byte("mysecret")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}
	ciphertext := gcm.Seal(nil, nonce, secret, nil)

	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// Decrypt secret
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Plaintext: %s\n", plaintext)
}
