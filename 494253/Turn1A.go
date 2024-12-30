package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt function using AES GCM
func encrypt(plaintext []byte, key []byte) (string, error) {
	// Create AES cipher block using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Generate a nonce with the correct size for AES GCM (12 bytes)
	nonceSize := gcm.NonceSize() // Typically 12 bytes for GCM
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext and append the nonce to the ciphertext
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Return the base64-encoded string of the nonce + ciphertext
	combined := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt function using AES GCM
func decrypt(ciphertextBase64 string, key []byte) ([]byte, error) {
	// Decode the base64 encoded ciphertext
	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Split the nonce and ciphertext
	nonce, ciphertext := data[:12], data[12:]

	// Create AES cipher block using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Decrypt the ciphertext using the nonce
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

func main() {
	// 32-byte key for AES-256 encryption (must be 32 bytes for AES-256)
	key := []byte("your_32_byte_aes_key_1234567890a") // 32-byte key (256 bits)
	plaintext := []byte("This is a secure notification.")

	// Encrypt the plaintext
	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Println("Encrypted:", ciphertext)

	// Decrypt the ciphertext
	decryptedText, err := decrypt(ciphertext, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Decrypted:", string(decryptedText))
}
