package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidBlockSize        = errors.New("invalid block size")
	ErrCiphertextTooShort      = errors.New("ciphertext is too short")
	ErrCiphertextNotFullBlocks = errors.New("ciphertext is not a multiple of the block size")
	ErrInvalidPKCS7Padding     = errors.New("invalid PKCS#7 padding")
)

func Encrypt(unencrypted, cipherKey string) (string, error) {
	key := []byte(cipherKey)
	plainText := []byte(unencrypted)
	plainText, err := Pad(plainText, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("pad plaintext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("generate ciphertext IV: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

func Decrypt(encrypted, cipherKey string) (string, error) {
	key := []byte(cipherKey)
	cipherText, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}

	if len(cipherText) <= aes.BlockSize {
		return "", ErrCiphertextTooShort
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		return "", ErrCiphertextNotFullBlocks
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, err = Unpad(cipherText, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad plaintext: %w", err)
	}
	return fmt.Sprintf("%s", cipherText), nil
}

func Pad(buf []byte, size int) ([]byte, error) {
	if size <= 0 {
		return nil, ErrInvalidBlockSize
	}

	bufLen := len(buf)
	padLen := size - bufLen%size
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := range padLen {
		padded[bufLen+i] = byte(padLen)
	}
	return padded, nil
}

func Unpad(padded []byte, size int) ([]byte, error) {
	if size <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(padded) == 0 || len(padded)%size != 0 {
		return nil, ErrInvalidPKCS7Padding
	}

	padLen := int(padded[len(padded)-1])
	if padLen == 0 || padLen > size || padLen > len(padded) {
		return nil, ErrInvalidPKCS7Padding
	}
	for _, value := range padded[len(padded)-padLen:] {
		if value != byte(padLen) {
			return nil, ErrInvalidPKCS7Padding
		}
	}

	bufLen := len(padded) - padLen
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf, nil
}
