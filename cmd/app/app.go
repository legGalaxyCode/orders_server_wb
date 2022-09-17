package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"os/signal"
	config2 "test_db_server/internal/config"
	model2 "test_db_server/internal/order"
	order "test_db_server/internal/order/db"
	"test_db_server/internal/publisher"
	service2 "test_db_server/internal/service"
	"test_db_server/pkg/client/postgresql"
	"test_db_server/pkg/logging"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("init config")
	config := config2.GetConfig()
	router := httprouter.New()

	logger.Info("Prepare database data from json models")
	jsonModel, err := os.ReadFile("../internal/order/model.json")
	jsonModelTest1, err := os.ReadFile("../internal/order/modelTest1.json")
	jsonModelTest2, err := os.ReadFile("../internal/order/modelTest2.json")
	if err != nil {
		logger.Infof("read failed %v", err)
	}
	var model model2.Order
	err = json.Unmarshal(jsonModel, &model)
	if err != nil {
		logger.Infof("model.json decoding failed %v", err)
	}
	var modelTest1 model2.Order
	err = json.Unmarshal(jsonModelTest1, &modelTest1)
	if err != nil {
		logger.Infof("modelTest1 decoding failed %v", err)
	}
	var modelTest2 model2.Order
	err = json.Unmarshal(jsonModelTest2, &modelTest2)
	if err != nil {
		logger.Infof("modelTest2 decoding failed %v", err)
	}

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), config.DatabaseConfig)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	service := service2.NewService(config.SubscriberConfig, postgreSQLClient, &logger)

	logger.Info("Fill database data with 3 orders")
	_, err = service.CreateMany(context.TODO(), &model, &modelTest1, &modelTest2)
	if err != nil {
		logger.Infof("%v", err)
	}

	// this is only for testing the difference between cache and db paths
	ordersHandler := service2.NewHandler(service, &logger)
	ordersHandler.Register(router)
	ordersDbHandler := model2.NewHandler(order.NewRepository(postgreSQLClient, &logger), &logger)
	ordersDbHandler.Register(router)

	go service.Run()
	pub := publisher.NewPublisher(config.PublisherConfig, &logger)
	go pub.Run()

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	go func() {
		logger.Fatal(server.ListenAndServe())
	}()

	// Stop serving if receiving Ctrl+C interrupt
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Received an interrupt, end serving...")
	time.Sleep(5 * time.Second)
}
