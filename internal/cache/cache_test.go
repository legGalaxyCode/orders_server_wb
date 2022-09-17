package cache

import (
	"test_db_server/internal/delivery"
	"test_db_server/internal/item"
	"test_db_server/internal/order"
	"test_db_server/internal/payment"
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
		OrderUid:    "33489",
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
		Items:           make([]item.Item, 0),
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

func TestCache_AddOne(t *testing.T) {
	cache := New()
	cache.AddOne(&testSuit[0])
	if cache.Data[testSuit[0].OrderUid].CustomerId != testSuit[0].CustomerId {
		t.Error("Error, customerId doesn't match")
	}
}

func TestCache_AddMany(t *testing.T) {
	cache := New()
	cache.AddMany(&testSuit[0], &testSuit[1], &testSuit[2])
	if data, ok := cache.Data[testSuit[0].OrderUid]; ok {
		if data.CustomerId != testSuit[0].CustomerId {
			t.Error("Error, customerId doesn't match")
		}
	} else {
		t.Error("Error, cache doesn't contain order 0")
	}
	if data, ok := cache.Data[testSuit[1].OrderUid]; ok {
		if data.CustomerId != testSuit[1].CustomerId {
			t.Error("Error, customerId doesn't match")
		}
	} else {
		t.Error("Error, cache doesn't contain order 1")
	}
	if data, ok := cache.Data[testSuit[2].OrderUid]; ok {
		if data.CustomerId != testSuit[2].CustomerId {
			t.Error("Error, customerId doesn't match")
		}
	} else {
		t.Error("Error, cache doesn't contain order 2")
	}
}

func TestCache_GetOne(t *testing.T) {
	cache := New()
	cache.AddOne(&testSuit[0])
	orderFromCache, err := cache.GetOne(testSuit[0].OrderUid)
	if err != nil {
		t.Error(err)
	}
	if orderFromCache.DeliveryService != testSuit[0].DeliveryService {
		t.Error("Error, cache contains wrong order")
	}
}

func TestCache_GetAll(t *testing.T) {
	cache := New()
	cache.AddMany(&testSuit[0], &testSuit[1], &testSuit[2])
	orderFromCache, err := cache.GetOne(testSuit[0].OrderUid)
	if err != nil {
		t.Error(err)
	}
	if orderFromCache.DeliveryService != testSuit[0].DeliveryService {
		t.Error("Error, cache contains wrong order 0")
	}
	orderFromCache, err = cache.GetOne(testSuit[1].OrderUid)
	if err != nil {
		t.Error(err)
	}
	if orderFromCache.DeliveryService != testSuit[1].DeliveryService {
		t.Error("Error, cache contains wrong order 1")
	}
	orderFromCache, err = cache.GetOne(testSuit[2].OrderUid)
	if err != nil {
		t.Error(err)
	}
	if orderFromCache.DeliveryService != testSuit[2].DeliveryService {
		t.Error("Error, cache contains wrong order 2")
	}
}
