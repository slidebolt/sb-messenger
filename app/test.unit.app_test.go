package app

import (
	"encoding/json"
	"testing"
	"time"

	messengersdk "github.com/slidebolt/sb-messenger-sdk"
)

func TestHelloManifest(t *testing.T) {
	h := New().Hello()
	if h.ID != "messenger" {
		t.Fatalf("id: got %q want %q", h.ID, "messenger")
	}
	if h.Kind != "service" {
		t.Fatalf("kind: got %q want %q", h.Kind, "service")
	}
}

func TestOnStartStartsReachableNATSServer(t *testing.T) {
	m := New()

	payload, err := m.OnStart(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer m.OnShutdown()

	deps := map[string]json.RawMessage{"messenger": payload}
	client, err := messengersdk.Connect(deps)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	done := make(chan string, 1)
	_, err = client.Subscribe("test.subject", func(msg *messengersdk.Message) {
		done <- string(msg.Data)
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Flush(); err != nil {
		t.Fatal(err)
	}
	if err := client.Publish("test.subject", []byte("hello")); err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-done:
		if got != "hello" {
			t.Fatalf("data: got %q want %q", got, "hello")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for pubsub")
	}
}
