package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	contract "github.com/slidebolt/sb-contract"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

type App struct {
	server *natsserver.Server
	port   int
}

func New() *App {
	return &App{}
}

func (a *App) Hello() contract.HelloResponse {
	return contract.HelloResponse{
		ID:              "messenger",
		Kind:            contract.KindService,
		ContractVersion: contract.ContractVersion,
	}
}

func (a *App) OnStart(deps map[string]json.RawMessage) (json.RawMessage, error) {
	port, err := freePort()
	if err != nil {
		return nil, fmt.Errorf("find free port: %w", err)
	}

	opts := &natsserver.Options{
		Host: "127.0.0.1",
		Port: port,
	}

	ns, err := natsserver.NewServer(opts)
	if err != nil {
		return nil, fmt.Errorf("create nats server: %w", err)
	}

	ns.Start()
	if !ns.ReadyForConnections(5_000_000_000) {
		return nil, fmt.Errorf("nats server failed to start")
	}

	a.server = ns
	a.port = port

	log.Printf("nats server listening on 127.0.0.1:%d", port)

	payload, _ := json.Marshal(map[string]any{
		"nats_url":  "127.0.0.1",
		"nats_port": port,
	})

	return payload, nil
}

func (a *App) OnShutdown() error {
	if a.server != nil {
		a.server.Shutdown()
		log.Println("nats server stopped")
	}
	return nil
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port, nil
}
