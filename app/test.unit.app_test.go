package app

import (
	"encoding/json"
	"io"
	"net/http"
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

func TestOnStartStartsMonitorWhenConfigured(t *testing.T) {
	t.Setenv("SB_MESSENGER_MONITOR_HOST", "127.0.0.1")
	t.Setenv("SB_MESSENGER_MONITOR_PORT", "18222")

	m := New()

	payload, err := m.OnStart(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer m.OnShutdown()

	var info struct {
		MonitorHost string `json:"nats_monitor_host"`
		MonitorPort int    `json:"nats_monitor_port"`
	}
	if err := json.Unmarshal(payload, &info); err != nil {
		t.Fatal(err)
	}
	if info.MonitorHost != "127.0.0.1" || info.MonitorPort != 18222 {
		t.Fatalf("monitor payload: got host=%q port=%d", info.MonitorHost, info.MonitorPort)
	}

	resp, err := http.Get("http://127.0.0.1:18222/varz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d want %d", resp.StatusCode, http.StatusOK)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Fatal("empty varz response")
	}
}
