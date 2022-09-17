package publisher

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"math/rand"
	"os"
	"os/signal"
	"test_db_server/internal/config"
	"test_db_server/pkg/logging"
	"time"
)

type Publisher struct {
	Name    string
	Channel string
	Cluster string
	Timeout string
	logger  *logging.Logger
}

func NewPublisher(cfg config.PublisherConfig, logger *logging.Logger) *Publisher {
	return &Publisher{
		Name:    cfg.Name,
		Channel: cfg.Channel,
		Cluster: cfg.Cluster,
		Timeout: cfg.Timeout,
		logger:  logger,
	}
}

func (p *Publisher) Publish(text []byte) {
	cluster, err := stan.Connect(p.Cluster, p.Name, stan.NatsURL("localhost:4222"))
	if err != nil {
		p.logger.Infof("failed connect to cluster %v", err)
		return
	}
	defer cluster.Close()
	cluster.Publish(p.Channel, text)
	p.logger.Infof("publish text in %s channel", p.Channel)
}

func randomIntNumber(up int) int {
	return rand.Int() % up
}

func (p *Publisher) Run() {
	var i int = 0
	go func() {
		for {
			p.Publish([]byte(`{
			"order_uid": "` + fmt.Sprintf("%d", randomIntNumber(100000)) + `",
			"track_number": "WBILMTESTTRACK",
			"entry": "WBIL",
			"delivery": {
			  "name": "Test Testov",
			  "phone": "+9720000000",
			  "zip": "2639809",
			  "city": "Kiryat Mozkin",
			  "address": "Ploshad Mira 15",
			  "region": "Kraiot",
			  "email": "test@gmail.com"
			},
			"payment": {
			  "transaction": "` + fmt.Sprintf("%d", randomIntNumber(100000)) + `",
			  "request_id": "",
			  "currency": "USD",
			  "provider": "wbpay",
			  "amount": ` + fmt.Sprintf("%d", randomIntNumber(10000)) + `,
			  "payment_dt": ` + fmt.Sprintf("%d", randomIntNumber(20000)) + `,
			  "bank": "alpha",
			  "delivery_cost": ` + fmt.Sprintf("%d", randomIntNumber(3000)) + `,
			  "goods_total": ` + fmt.Sprintf("%d", randomIntNumber(1000)) + `,
			  "custom_fee": 0
			},
			"items": [
			  {
				"chrt_id": ` + fmt.Sprintf("%d", randomIntNumber(550000)) + `,
				"track_number": "WBILMTESTTRACK",
				"price": ` + fmt.Sprintf("%d", randomIntNumber(14500)) + `,
				"rid": "` + fmt.Sprintf("%d", randomIntNumber(100000)) + `",
				"name": "Mascaras",
				"sale": 30,
				"size": "0",
				"total_price": 317,
				"nm_id": ` + fmt.Sprintf("%d", randomIntNumber(99999)) + `,
				"brand": "Vivienne Sabo",
				"status": 202
			  }
			],
			"locale": "en",
			"internal_signature": "",
			"customer_id": "test",
			"delivery_service": "meest",
			"shardkey": "9",
			"sm_id": ` + fmt.Sprintf("%d", randomIntNumber(555)) + `,
			"date_created": "2021-11-26T06:22:19Z",
			"oof_shard": "1"
		  }`))
			i++
			timeout, _ := time.ParseDuration(p.Timeout)
			time.Sleep(timeout)
		}
	}()
	// Unsubscribe if receiving Ctrl+C interrupt
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Received an interrupt, end publishing...")
}
