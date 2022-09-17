package service

import (
	"context"
	config2 "test_db_server/internal/config"
	"test_db_server/internal/delivery"
	"test_db_server/internal/order"
	"test_db_server/internal/payment"
	"test_db_server/pkg/client/postgresql"
	"test_db_server/pkg/logging"
	"testing"
)

var testSuit = []order.Order{
	order.Order{
		OrderUid:        "123",
		TrackNumber:     "1",
		Entry:           "2",
		Delivery:        delivery.Delivery{},
		Payment:         payment.Payment{},
		Items:           nil,
		Locale:          "en",
		InternalSig:     "12",
		CustomerId:      "98123",
		DeliveryService: "wb",
		ShardKey:        "213123",
		SmId:            1,
		DateCreated:     "now",
		OofShard:        "34",
	},
	order.Order{
		OrderUid:        "124",
		TrackNumber:     "12123",
		Entry:           "289234",
		Delivery:        delivery.Delivery{},
		Payment:         payment.Payment{},
		Items:           nil,
		Locale:          "ru",
		InternalSig:     "12212",
		CustomerId:      "986151",
		DeliveryService: "sb",
		ShardKey:        "14123",
		SmId:            4,
		DateCreated:     "nlk",
		OofShard:        "12",
	},
	order.Order{
		OrderUid:    "1025",
		TrackNumber: "89ubsg123",
		Entry:       "12378asd",
		Delivery: delivery.Delivery{
			Name:    "alex",
			Phone:   "+89137642",
			Zip:     "123",
			City:    "Moscow",
			Address: "BB 14",
			Region:  "North",
			Email:   "test@gmail.com",
		},
		Payment: payment.Payment{
			Transaction:  "ksdf",
			RequestId:    "ksifj29113",
			Currency:     "rub",
			Provider:     "wb",
			Amount:       132,
			PaymentDt:    923324,
			Bank:         "sber",
			DeliveryCost: 1209,
			GoodsTotal:   12,
			CustomFee:    87,
		},
		Items:           nil,
		Locale:          "port",
		InternalSig:     "213123",
		CustomerId:      "0779123",
		DeliveryService: "wb",
		ShardKey:        "451623",
		SmId:            5,
		DateCreated:     "now",
		OofShard:        "35",
	},
}

func TestService_Create(t *testing.T) {
	logging.Init()
	logger := logging.GetLogger()
	config := config2.GetConfig()
	client, _ := postgresql.NewClient(context.TODO(), config.DatabaseConfig)
	serv := NewService(config.SubscriberConfig, client, &logger)
	_, err := serv.Create(context.TODO(), &testSuit[0])
	if err != nil {
		t.Errorf("Error in creating order %v", err)
	}
	res, _ := serv.FindOne(testSuit[0].OrderUid)
	if res.DateCreated != testSuit[0].DateCreated {
		t.Error("Error in creating order")
	}
}

func TestService_CreateMany(t *testing.T) {
	logging.Init()
	logger := logging.GetLogger()
	config := config2.GetConfig()
	client, _ := postgresql.NewClient(context.TODO(), config.DatabaseConfig)
	serv := NewService(config.SubscriberConfig, client, &logger)
	_, err := serv.CreateMany(context.TODO(), &testSuit[1], &testSuit[2])
	if err != nil {
		t.Errorf("Error in creating many orders: %v", err)
	}
	res, _ := serv.FindAll()
	for i, val := range res {
		if i == 0 {
			continue
		}
		if val.TrackNumber != testSuit[i].TrackNumber {
			t.Errorf("Error in find order %d", i)
		}
	}
}
