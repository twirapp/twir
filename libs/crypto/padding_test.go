package crypto

import (
	"crypto/aes"
	"errors"
	"math"
	"testing"
)

func TestPadReturnsInvalidBlockSizeWhenSizeIsNotPKCS7Representable(t *testing.T) {
	// Given
	cases := []struct {
		name string
		size int
	}{
		{name: "negative", size: -1},
		{name: "zero", size: 0},
		{name: "larger than one byte", size: 256},
		{name: "maximum integer", size: math.MaxInt},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			panicked := false
			var padErr error

			// When
			func() {
				defer func() {
					panicked = recover() != nil
				}()

				_, padErr = Pad(nil, testCase.size)
			}()

			// Then
			if panicked {
				t.Fatal("Pad panicked for an invalid block size")
			}
			if !errors.Is(padErr, ErrInvalidBlockSize) {
				t.Fatal("Pad did not return the invalid block size error category")
			}
		})
	}
}

func TestUnpadReturnsInvalidBlockSizeWhenSizeIsNotPKCS7Representable(t *testing.T) {
	// Given
	cases := []struct {
		name string
		size int
	}{
		{name: "negative", size: -1},
		{name: "zero", size: 0},
		{name: "larger than one byte", size: 256},
		{name: "near maximum integer", size: math.MaxInt - 1},
		{name: "maximum integer", size: math.MaxInt},
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

				_, unpadErr = Unpad(nil, testCase.size)
			}()

			// Then
			if panicked {
				t.Fatal("Unpad panicked for an invalid block size")
			}
			if !errors.Is(unpadErr, ErrInvalidBlockSize) {
				t.Fatal("Unpad did not return the invalid block size error category")
			}
		})
	}
}

func TestPadAndUnpadRoundTripWhenBlockSizeIsPKCS7Representable(t *testing.T) {
	// Given
	plainText := []byte("padding fixture")
	cases := []struct {
		name string
		size int
	}{
		{name: "one", size: 1},
		{name: "AES", size: aes.BlockSize},
		{name: "maximum representable", size: 255},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			// When
			padded, padErr := Pad(plainText, testCase.size)
			unpadded, unpadErr := Unpad(padded, testCase.size)

			// Then
			if padErr != nil || unpadErr != nil {
				t.Fatal("PKCS#7-representable block size returned an error")
			}
			if len(padded)%testCase.size != 0 || string(unpadded) != string(plainText) {
				t.Fatal("Pad and Unpad did not preserve the plaintext")
			}
		})
	}
}
