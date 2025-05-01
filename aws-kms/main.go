package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func main() {
	ctx := context.Background()

	// Load default AWS config (uses ~/.aws/credentials or env vars)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Replace this with your KMS key ID or ARN
	keyID := "<AWS ARN"

	// Plaintext to encrypt
	plaintext := []byte("Hello, AWS KMS!")

	// Encrypt
	encryptOutput, err := kmsClient.Encrypt(ctx, &kms.EncryptInput{
		KeyId:     &keyID,
		Plaintext: plaintext,
	})
	if err != nil {
		log.Fatalf("failed to encrypt: %v", err)
	}
	ciphertextBlob := encryptOutput.CiphertextBlob
	fmt.Println(*encryptOutput.KeyId)

	// Print encrypted base64 encoded
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertextBlob)
	fmt.Println("Encrypted (base64):", encodedCiphertext)

	// Decrypt
	decryptOutput, err := kmsClient.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: ciphertextBlob,
	})
	if err != nil {
		log.Fatalf("failed to decrypt: %v", err)
	}

	decryptedPlaintext := string(decryptOutput.Plaintext)
	fmt.Println("Decrypted:", decryptedPlaintext)

	// ------------------------------------------------
	// Generate Data Key

	res, err := kmsClient.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
		KeyId:   &keyID,
		KeySpec: types.DataKeySpecAes256,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.CiphertextBlob)
	fmt.Println(res.CiphertextForRecipient)
	fmt.Println(res.Plaintext)

	plaintextDataKey := res.Plaintext
	encryptedDataKey := res.CiphertextBlob

	fmt.Println("Generated encrypted data key (base64):", base64.StdEncoding.EncodeToString(encryptedDataKey))

	// Step 2: Encrypt a string locally using AES-GCM
	plaintext = []byte("Hello, this is a secret message!")
	ciphertext, nonce, err := encryptWithAESGCM(plaintextDataKey, plaintext)
	if err != nil {
		log.Fatalf("failed to encrypt data: %v", err)
	}

	fmt.Println("Encrypted message (base64):", base64.StdEncoding.EncodeToString(ciphertext))
	fmt.Println("Nonce (base64):", base64.StdEncoding.EncodeToString(nonce))

	// -------- Time passes... You now want to decrypt --------

	// Step 3: Decrypt the data key using KMS
	decryptOutput, err = kmsClient.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: encryptedDataKey,
	})
	if err != nil {
		log.Fatalf("failed to decrypt data key: %v", err)
	}
	decryptedDataKey := decryptOutput.Plaintext

	// Step 4: Decrypt the message locally using AES-GCM
	decryptedPlaintextByte, err := decryptWithAESGCM(decryptedDataKey, ciphertext, nonce)
	if err != nil {
		log.Fatalf("failed to decrypt message: %v", err)
	}

	fmt.Println("Decrypted message:", string(decryptedPlaintextByte))
}

// AES-GCM encryption helper
func encryptWithAESGCM(key, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// AES-GCM decryption helper
func decryptWithAESGCM(key, ciphertext, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
