package service

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"test_db_server/pkg/logging"
)

const (
	ordersURL = "/cache/orders/"
	orderURL  = "/cache/orders/:uuid"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	service Service
	logger  *logging.Logger
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, ordersURL, h.GetList)
	router.HandlerFunc(http.MethodGet, orderURL, h.GetOne)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	all, err := h.service.FindAll()
	if err != nil {
		w.WriteHeader(400)
		h.logger.Fatalf("%v", err)
		return
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(allBytes)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return
	}
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) {
	str := strings.Split(r.URL.Path, "/")
	one, err := h.service.FindOne(str[len(str)-1])
	if err != nil {
		w.WriteHeader(400)
		h.logger.Fatalf("%v", err)
		return
	}

	allBytes, err := json.Marshal(one)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(allBytes)
	if err != nil {
		h.logger.Fatalf("%v", err)
		return
	}
}

func NewHandler(service Service, logger *logging.Logger) Handler {
	return &handler{
		service: service,
		logger:  nil,
	}
}
