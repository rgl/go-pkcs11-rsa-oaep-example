package main

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	_ "crypto/sha256"

	"github.com/miekg/pkcs11"
	"github.com/miekg/pkcs11/p11"
	oaep "github.com/rgl/go-pkcs11-rsa-oaep"
)

func openSession(module p11.Module, userPin string, tokenLabel string) (p11.Session, error) {
	slots, err := module.Slots()
	if err != nil {
		return nil, fmt.Errorf("failed to get slot list: %v", err)
	}
	for _, slot := range slots {
		tokenInfo, err := slot.TokenInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to get slot %d token info: %v", slot.ID(), err)
		}
		if tokenInfo.Flags&pkcs11.CKF_TOKEN_INITIALIZED == 0 {
			continue
		}
		if strings.TrimRight(tokenInfo.Label, "\x00") != tokenLabel {
			continue
		}
		session, err := slot.OpenSession()
		if err != nil {
			return nil, fmt.Errorf("failed to open session into slot %d token %s: %v", slot.ID(), tokenInfo.Label, err)
		}
		err = session.Login(userPin)
		if err != nil {
			session.Close()
			return nil, fmt.Errorf("failed to login into slot %d token %s: %v", slot.ID(), tokenInfo.Label, err)
		}
		return session, nil
	}
	return nil, fmt.Errorf("token %s not found", tokenLabel)
}

func main() {
	pkcs11LibraryPath := "/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
	if p := os.Getenv("TEST_PKCS11_LIBRARY_PATH"); p != "" {
		pkcs11LibraryPath = p
	}

	userPin := ""
	if p := os.Getenv("TEST_PKCS11_USER_PIN"); p != "" {
		userPin = p
	}

	tokenLabel := "test"
	if p := os.Getenv("TEST_PKCS11_TOKEN_LABEL"); p != "" {
		tokenLabel = p
	}

	keyLabel := "test-rsa-2048"
	if p := os.Getenv("TEST_PKCS11_KEY_LABEL"); p != "" {
		keyLabel = p
	}

	log.Printf("Loading %s...", pkcs11LibraryPath)
	module, err := p11.OpenModule(pkcs11LibraryPath)
	if err != nil {
		log.Fatalf("failed to open pkcs11 module %s: %v", pkcs11LibraryPath, err)
	}

	log.Printf("Opening session to %s...", tokenLabel)
	session, err := openSession(module, userPin, tokenLabel)
	if err != nil {
		log.Fatalf("failed to open session to token %s: %v", tokenLabel, err)
	}
	defer session.Close()

	log.Printf("Getting the %s key...", keyLabel)
	publicKey, privateKey, err := oaep.GetKey(session, keyLabel)
	if err != nil {
		log.Fatalf("Key %s not found in HSM: %v", keyLabel, err)
	}

	random := rand.Reader

	plaintext := []byte("abracadabra")

	log.Printf("Encrypting %s...", plaintext)
	ciphertext, err := oaep.Encrypt(crypto.SHA256.New(), random, publicKey, plaintext, nil)
	if err != nil {
		log.Fatalf("Failed to encrypt: %v", err)
	}

	log.Printf("Decrypting %s...", hex.EncodeToString(ciphertext))
	result, err := oaep.Decrypt(crypto.SHA256.New(), random, publicKey.Size(), privateKey, ciphertext, nil)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}

	log.Printf("Descrypted as %s", result)
}
