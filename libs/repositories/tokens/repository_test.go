package tokens

import (
	"testing"

	"github.com/google/uuid"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

func TestTokenInputsAndModelSupportNullableEncryptedDeviceID(t *testing.T) {
	encryptedDeviceID := "encrypted-device-id"

	create := CreateInput{UserID: uuid.New(), DeviceID: &encryptedDeviceID}
	if create.DeviceID == nil || *create.DeviceID != encryptedDeviceID {
		t.Fatalf("create input device ID = %#v, want %q", create.DeviceID, encryptedDeviceID)
	}

	update := UpdateTokenInput{DeviceID: &encryptedDeviceID}
	if update.DeviceID == nil || *update.DeviceID != encryptedDeviceID {
		t.Fatalf("update input device ID = %#v, want %q", update.DeviceID, encryptedDeviceID)
	}

	token := tokenmodel.Token{DeviceID: &encryptedDeviceID}
	if token.DeviceID == nil || *token.DeviceID != encryptedDeviceID {
		t.Fatalf("token device ID = %#v, want %q", token.DeviceID, encryptedDeviceID)
	}

	var existingToken tokenmodel.Token
	if existingToken.DeviceID != nil {
		t.Fatalf("existing token device ID = %#v, want nil", existingToken.DeviceID)
	}
}
