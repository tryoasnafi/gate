package crypto

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/term"
)

const (
	saltSize      = 16
	aeadKeySize   = chacha20poly1305.KeySize
	aeadNonceSize = chacha20poly1305.NonceSize
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id
	argonTime    = 1
	argonMemory  = 12 * 1024 // 12 MB
	argonThreads = 3
	argonKeyLen  = aeadKeySize
	maxRetries   = 3
)

func Encrypt(password, plaintext []byte) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key := argon2.IDKey(password, salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aeadNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	return append(salt, append(nonce, ciphertext...)...), nil
}

func Decrypt(password, encrypted []byte) ([]byte, error) {
	if len(encrypted) < saltSize+aeadNonceSize {
		return nil, errors.New("invalid encrypted data")
	}
	salt := encrypted[:saltSize]
	nonce := encrypted[saltSize : saltSize+aeadNonceSize]
	ciphertext := encrypted[saltSize+aeadNonceSize:]

	key := argon2.IDKey(password, salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, nonce, ciphertext, nil)
}

func ReadPassword() ([]byte, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return nil, errors.New("stdin is not a terminal")
	}
	pwd, err := term.ReadPassword(fd)
	fmt.Println()
	return bytes.TrimSpace(pwd), err
}

func PromptPasswordWithRetry(prompt string, decryptFunc func([]byte) error) error {
	for range maxRetries {
		fmt.Print(prompt)
		pwd, err := ReadPassword()
		if err != nil {
			return err
		}
		if err := decryptFunc(pwd); err == nil {
			return nil
		}
		fmt.Println("Wrong password. Try again.")
	}
	return errors.New("maximum retry attempts exceeded")
}
