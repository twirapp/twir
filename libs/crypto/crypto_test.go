package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

const (
	fixtureCipherKey = "0123456789abcdef0123456789abcdef"
	fixturePlaintext = "synthetic-cbc-credential"
)

func legacyCBCFixture(t *testing.T) string {
	t.Helper()

	block, err := aes.NewCipher([]byte(fixtureCipherKey))
	if err != nil {
		t.Fatal("could not construct AES fixture cipher")
	}

	plainText := []byte(fixturePlaintext)
	padLength := aes.BlockSize - len(plainText)%aes.BlockSize
	padded := make([]byte, len(plainText)+padLength)
	copy(padded, plainText)
	for index := len(plainText); index < len(padded); index++ {
		padded[index] = byte(padLength)
	}

	cipherText := make([]byte, aes.BlockSize+len(padded))
	copy(cipherText, []byte("legacy-cbc-iv!!!"))
	cipher.NewCBCEncrypter(block, cipherText[:aes.BlockSize]).CryptBlocks(cipherText[aes.BlockSize:], padded)

	return hex.EncodeToString(cipherText)
}

func TestDecryptDecryptsLegacyCBCFixture(t *testing.T) {
	// Given
	legacyCiphertext := legacyCBCFixture(t)

	// When
	plainText, err := Decrypt(legacyCiphertext, fixtureCipherKey)

	// Then
	if err != nil {
		t.Fatal("Decrypt returned an error for a legacy CBC fixture")
	}
	if plainText != fixturePlaintext {
		t.Fatal("Decrypt did not recover the legacy fixture plaintext")
	}
}

func TestEncryptProducesHexCBCCompatibleCiphertext(t *testing.T) {
	// Given
	const expectedMinimumCiphertextLength = aes.BlockSize * 2

	// When
	encrypted, err := Encrypt(fixturePlaintext, fixtureCipherKey)

	// Then
	if err != nil {
		t.Fatal("Encrypt returned an error for a valid fixture")
	}
	decoded, err := hex.DecodeString(encrypted)
	if err != nil {
		t.Fatal("Encrypt did not return hexadecimal ciphertext")
	}
	if len(decoded) < expectedMinimumCiphertextLength || len(decoded)%aes.BlockSize != 0 {
		t.Fatal("Encrypt did not return an IV-prefixed CBC ciphertext")
	}

	block, err := aes.NewCipher([]byte(fixtureCipherKey))
	if err != nil {
		t.Fatal("could not construct AES fixture cipher")
	}
	cipher.NewCBCDecrypter(block, decoded[:aes.BlockSize]).CryptBlocks(decoded[aes.BlockSize:], decoded[aes.BlockSize:])
	paddedPlainText, err := Unpad(decoded[aes.BlockSize:], aes.BlockSize)
	if err != nil {
		t.Fatal("Encrypt did not apply PKCS#7 padding")
	}
	if string(paddedPlainText) != fixturePlaintext {
		t.Fatal("Encrypt ciphertext was not CBC-compatible")
	}
}

func TestDecryptReturnsRedactedErrorWhenCiphertextIsNotHex(t *testing.T) {
	// Given
	const malformedCiphertext = "not-hex-synthetic-credential-sentinel"
	const cipherKey = fixtureCipherKey

	panicked := false
	var decryptErr error

	// When
	func() {
		defer func() {
			panicked = recover() != nil
		}()

		_, decryptErr = Decrypt(malformedCiphertext, cipherKey)
	}()

	// Then
	if panicked {
		t.Fatal("Decrypt panicked for malformed ciphertext")
	}
	if decryptErr == nil {
		t.Fatal("Decrypt returned nil error for malformed ciphertext")
	}
	if strings.Contains(decryptErr.Error(), malformedCiphertext) {
		t.Fatal("Decrypt error contained sensitive ciphertext")
	}
}

func TestDecryptReturnsRedactedErrorsForMalformedInputs(t *testing.T) {
	// Given
	const invalidKey = "synthetic-key-sentinel"
	cases := []struct {
		name       string
		ciphertext string
		cipherKey  string
		secret     string
		wantErr    error
	}{
		{
			name:       "invalid AES key length",
			ciphertext: legacyCBCFixture(t),
			cipherKey:  invalidKey,
			secret:     invalidKey,
		},
		{
			name:       "ciphertext shorter than IV",
			ciphertext: hex.EncodeToString(make([]byte, aes.BlockSize-1)),
			cipherKey:  fixtureCipherKey,
			wantErr:    ErrCiphertextTooShort,
		},
		{
			name:       "ciphertext contains only an IV",
			ciphertext: hex.EncodeToString(make([]byte, aes.BlockSize)),
			cipherKey:  fixtureCipherKey,
			wantErr:    ErrCiphertextTooShort,
		},
		{
			name:       "ciphertext body is not a full block",
			ciphertext: hex.EncodeToString(make([]byte, aes.BlockSize+1)),
			cipherKey:  fixtureCipherKey,
			wantErr:    ErrCiphertextNotFullBlocks,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			panicked := false
			var decryptErr error

			// When
			func() {
				defer func() {
					panicked = recover() != nil
				}()

				_, decryptErr = Decrypt(testCase.ciphertext, testCase.cipherKey)
			}()

			// Then
			if panicked {
				t.Fatal("Decrypt panicked for malformed input")
			}
			if decryptErr == nil {
				t.Fatal("Decrypt returned nil error for malformed input")
			}
			if testCase.wantErr != nil && !errors.Is(decryptErr, testCase.wantErr) {
				t.Fatal("Decrypt did not preserve the expected error category")
			}
			if testCase.secret != "" && strings.Contains(decryptErr.Error(), testCase.secret) {
				t.Fatal("Decrypt error contained sensitive input")
			}
		})
	}
}

func TestDecryptReturnsRedactedErrorWhenPKCS7PaddingIsInvalid(t *testing.T) {
	// Given
	const malformedPaddedBlock = "secret-sentinel!"

	block, err := aes.NewCipher([]byte(fixtureCipherKey))
	if err != nil {
		t.Fatal("could not construct AES fixture cipher")
	}
	cipherText := make([]byte, aes.BlockSize*2)
	copy(cipherText, []byte("legacy-cbc-iv!!!"))
	cipher.NewCBCEncrypter(block, cipherText[:aes.BlockSize]).CryptBlocks(cipherText[aes.BlockSize:], []byte(malformedPaddedBlock))

	panicked := false
	var decryptErr error

	// When
	func() {
		defer func() {
			panicked = recover() != nil
		}()

		_, decryptErr = Decrypt(hex.EncodeToString(cipherText), fixtureCipherKey)
	}()

	// Then
	if panicked {
		t.Fatal("Decrypt panicked for invalid PKCS#7 padding")
	}
	if decryptErr == nil {
		t.Fatal("Decrypt returned nil error for invalid PKCS#7 padding")
	}
	if !errors.Is(decryptErr, ErrInvalidPKCS7Padding) {
		t.Fatal("Decrypt did not preserve the invalid padding error category")
	}
	if strings.Contains(decryptErr.Error(), malformedPaddedBlock) {
		t.Fatal("Decrypt error contained sensitive plaintext")
	}
}

func TestUnpadReturnsErrorWhenPaddingIsMalformed(t *testing.T) {
	// Given
	cases := []struct {
		name    string
		data    []byte
		size    int
		wantErr error
	}{
		{name: "zero block size", data: nil, size: 0, wantErr: ErrInvalidBlockSize},
		{name: "empty padded value", data: nil, size: aes.BlockSize, wantErr: ErrInvalidPKCS7Padding},
		{name: "misaligned padded value", data: []byte("short"), size: aes.BlockSize, wantErr: ErrInvalidPKCS7Padding},
		{name: "zero padding length", data: append(make([]byte, aes.BlockSize-1), 0), size: aes.BlockSize, wantErr: ErrInvalidPKCS7Padding},
		{name: "padding length exceeds block", data: append(make([]byte, aes.BlockSize-1), aes.BlockSize+1), size: aes.BlockSize, wantErr: ErrInvalidPKCS7Padding},
		{name: "mismatched padding", data: append(append(make([]byte, aes.BlockSize-2), 1), 2), size: aes.BlockSize, wantErr: ErrInvalidPKCS7Padding},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			panicked := false
			var unpadErr error

			// When
			func() {
				defer func() {
					panicked = recover() != nil
				}()

				_, unpadErr = Unpad(testCase.data, testCase.size)
			}()

			// Then
			if panicked {
				t.Fatal("Unpad panicked for malformed padding")
			}
			if unpadErr == nil {
				t.Fatal("Unpad returned nil error for malformed padding")
			}
			if !errors.Is(unpadErr, testCase.wantErr) {
				t.Fatal("Unpad did not preserve the expected error category")
			}
		})
	}
}

func TestEncryptReturnsRedactedErrorWhenAESKeyIsInvalid(t *testing.T) {
	// Given
	const plaintext = "synthetic-plaintext-sentinel"
	const invalidKey = "short"

	// When
	_, encryptErr := Encrypt(plaintext, invalidKey)

	// Then
	if encryptErr == nil {
		t.Fatal("Encrypt returned nil error for an invalid AES key")
	}
	if strings.Contains(encryptErr.Error(), plaintext) {
		t.Fatal("Encrypt error contained sensitive plaintext")
	}
}
