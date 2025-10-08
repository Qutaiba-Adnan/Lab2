package messaging

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// Connect connects to the NATS server and returns a connection
func Connect() *nats.Conn {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = "nats://127.0.0.1:4222"
	}

	nc, err := nats.Connect(url,
		nats.Name("FirefightingSystem"),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(10),
	)
	if err != nil {
		log.Fatalf("Could not connect to NATS: %v", err)
	}

	log.Printf("Connected to NATS at %s", nc.ConnectedUrl())
	return nc
}

// PublishJSON publishes a struct as JSON to a NATS subject
func PublishJSON(nc *nats.Conn, subject string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return nc.Publish(subject, data)
}

// SubscribeJSON subscribes to a subject and unmarshals messages to a struct
func SubscribeJSON[T any](nc *nats.Conn, subject string, handler func(T)) (*nats.Subscription, error) {
	return nc.Subscribe(subject, func(msg *nats.Msg) {
		var data T
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("JSON unmarshal error on %s: %v", subject, err)
			return
		}
		handler(data)
	})
}

// Gracefully close the connection
func Drain(nc *nats.Conn) {
	if nc != nil {
		log.Println("Draining NATS connection...")
		nc.Drain()
		nc.Close()
	}
}
