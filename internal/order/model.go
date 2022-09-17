package order

import (
	"test_db_server/internal/delivery"
	item2 "test_db_server/internal/item"
	"test_db_server/internal/payment"
)

type Order struct {
	OrderUid        string            `json:"order_uid"`
	TrackNumber     string            `json:"track_number"`
	Entry           string            `json:"entry"`
	Delivery        delivery.Delivery `json:"delivery"`
	Payment         payment.Payment   `json:"payment"`
	Items           []item2.Item      `json:"items"`
	Locale          string            `json:"locale"`
	InternalSig     string            `json:"internal_signature"`
	CustomerId      string            `json:"customer_id"`
	DeliveryService string            `json:"delivery_service"`
	ShardKey        string            `json:"shardkey"`
	SmId            int               `json:"sm_id"`
	DateCreated     string            `json:"date_created"`
	OofShard        string            `json:"oof_shard"`
}
