package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
	"os/signal"
	"test_db_server/internal/cache"
	"test_db_server/internal/config"
	"test_db_server/internal/order"
	order2 "test_db_server/internal/order/db"
	"test_db_server/pkg/client/postgresql"
	"test_db_server/pkg/logging"
)

type service struct {
	Name    string
	Channel string
	Cluster string
	Timeout string
	repo    order.Repository
	cache   cache.Cache
}

func (s *service) Create(ctx context.Context, order *order.Order) (bool, error) {
	one, _ := s.cache.GetOne(order.OrderUid)

	// if default value
	if one.OrderUid == "" {
		isCreated, err := s.repo.Create(ctx, order)
		if err != nil || !isCreated {
			return false, err
		}
		s.cache.AddOne(order)
		return true, nil
	}
	return false, fmt.Errorf("order has already exists")
}

func (s *service) CreateMany(ctx context.Context, orders ...*order.Order) (bool, error) {
	for _, ord := range orders {
		_, err := s.Create(ctx, ord)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s *service) FindAll() ([]order.Order, error) {
	return s.cache.GetAll(), nil
}

func (s *service) FindOne(orderId string) (order.Order, error) {
	return s.cache.GetOne(orderId)
}

func (s *service) SyncData() {
	ordersDb, _ := s.repo.FindAll(context.TODO())
	for _, ord := range ordersDb {
		s.cache.AddOne(&ord)
	}
}

// subscribe functional

func (s *service) Run() {
	s.SyncData()
	cluster, err := stan.Connect(s.Cluster, s.Name)
	if err != nil {
		fmt.Printf("Cannot connect to cluster %s\n", s.Cluster)
	}

	queryHandler := func(msg *stan.Msg) {
		var ord order.Order
		err = json.Unmarshal(msg.Data, &ord)
		fmt.Printf("received msg order with order_ud: %s", ord.OrderUid)

		if err == nil {
			_, err = s.Create(context.TODO(), &ord)
			if err != nil {
				fmt.Printf("error create an order %v\n", err)
			}
		} else {
			fmt.Printf("invalid msg from channel %s\n", msg.Data)
		}
	}

	sub, err := cluster.QueueSubscribe(s.Channel, "service-chan", queryHandler, stan.DurableName("durName"))
	if err != nil {
		fmt.Printf("Cannot subscribe to channel %s\n", s.Channel)
	}

	fmt.Printf("Connected to clusterID: %s clientID: %s\n", s.Cluster, s.Name)

	// Unsubscribe if receiving Ctrl+C interrupt
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Received an interrupt, unsubscribing and closing connection...")
	sub.Unsubscribe()
	cluster.Close()
}

type Service interface {
	Create(ctx context.Context, order *order.Order) (bool, error)
	CreateMany(ctx context.Context, orders ...*order.Order) (bool, error)
	FindAll() ([]order.Order, error)
	FindOne(orderId string) (order.Order, error)
	Run()
}

func NewService(cfg config.SubscriberConfig, client postgresql.Client, logger *logging.Logger) Service {
	s := &service{
		Name:    cfg.Name,
		Channel: cfg.Channel,
		Cluster: cfg.Cluster,
		Timeout: cfg.Timeout,
		repo:    order2.NewRepository(client, logger),
		cache:   *cache.New(),
	}
	return s
}
