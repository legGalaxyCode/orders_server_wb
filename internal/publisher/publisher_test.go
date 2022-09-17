package publisher

import (
	"github.com/nats-io/stan.go"
	config2 "test_db_server/internal/config"
	"test_db_server/pkg/logging"
	"testing"
)

func TestPublisher_Publish(t *testing.T) {
	logging.Init()
	cfg := config2.GetConfig()
	logger := logging.GetLogger()
	pub := [3]*Publisher{
		NewPublisher(cfg.PublisherConfig, &logger),
		NewPublisher(cfg.PublisherConfig, &logger),
		NewPublisher(cfg.PublisherConfig, &logger),
	}
	pub[1].Channel = "failed"
	pub[2].Cluster = "failed"

	mustGetBack := [...]string{"text", "", ""}
	response := [3]string{}

	for i := range mustGetBack {
		cl, _ := stan.Connect("test-cluster", "publisher")
		sub, _ := cl.Subscribe("test_channel", func(msg *stan.Msg) {
			response[i] = string(msg.Data)
		})
		pub[i].Publish([]byte("text"))

		sub.Unsubscribe()
		cl.Close()

		if response[i] != mustGetBack[i] {
			t.Errorf("Failed data at %d", i)
		}
	}
}
