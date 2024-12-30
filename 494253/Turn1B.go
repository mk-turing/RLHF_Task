
package main
import (
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"
)
// Encrypts the given data using AES-256 encryption and a random key.
// The key is returned along with the ciphertext.
func encryptData(data []byte) ([]byte, []byte, error) {
	key := make([]byte, 32) // 256 bits key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, key, nil
}
// Decrypts the given ciphertext using the provided key.
func decryptData(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
// Saves the encrypted notification data to a file.
func saveNotificationToFile(encryptedData []byte, key []byte, filename string) error {
	keyFile, err := os.Create(filename + ".key")
	if err != nil {
		return err
	}
	defer keyFile.Close()
	if _, err := keyFile.Write(key); err != nil {
		return err
	}
	dataFile, err := os.Create(filename + ".data")
	if err != nil {
		return err
	}
	defer dataFile.Close()
	if _, err := dataFile.Write(encryptedData); err != nil {
		return err
	}
	return nil
}
// Reads the encrypted notification data from a file and decrypts it.
func readNotificationFromFile(filename string) ([]byte, error) {
	keyFile, err := os.Open(filename + ".key")
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()
	key, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, err
	}
	dataFile, err := os.Open(filename + ".data")
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()
	encryptedData, err := ioutil.ReadAll(dataFile)
	if err != nil {
		return nil, err
	}
	decryptedData, err := decryptData(encryptedData, key)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}
func main() {
	notificationData := []byte("This is a sensitive notification")
	encryptedData, key, err := encryptData(notificationData)
	if err != nil {
		panic(err)
	}