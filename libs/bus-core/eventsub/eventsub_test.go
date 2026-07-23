package eventsub_test

import (
	"reflect"
	"testing"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func TestUnsubscribeRequestGobRoundTripPreservesBindingSnapshot(t *testing.T) {
	requestType := reflect.TypeOf(eventsub.EventsubUnsubscribeRequest{})
	bindingField, found := requestType.FieldByName("Binding")
	if !found {
		t.Fatal("unsubscribe request is missing Binding")
	}
	if bindingField.Type.Kind() != reflect.Ptr {
		t.Fatalf("Binding type = %s, want pointer", bindingField.Type)
	}
	snapshotType := bindingField.Type.Elem()

	request := reflect.New(requestType).Elem()
	request.FieldByName("ChannelID").SetString("channel-id")
	request.FieldByName("Platform").Set(reflect.ValueOf(platformentity.PlatformKick))
	snapshot := reflect.New(snapshotType)
	for name, value := range map[string]string{
		"ID":                "binding-id",
		"UserID":            "provider-user-id",
		"PlatformChannelID": "provider-channel-id",
	} {
		field := snapshot.Elem().FieldByName(name)
		if !field.IsValid() || field.Kind() != reflect.String {
			t.Fatalf("snapshot field %s is missing or is not a string", name)
		}
		field.SetString(value)
	}
	request.FieldByName("Binding").Set(snapshot)
	want := request.Interface().(eventsub.EventsubUnsubscribeRequest)

	encoded, err := buscore.EncodeGob(want)
	if err != nil {
		t.Fatalf("encode unsubscribe request: %v", err)
	}
	got, err := buscore.DecodeGob[eventsub.EventsubUnsubscribeRequest](encoded)
	if err != nil {
		t.Fatalf("decode unsubscribe request: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}
