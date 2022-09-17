package item

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"test_db_server/pkg/logging"
)

const (
	itemsURL = "/items"
	itemURL  = "/items/:uuid"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, itemsURL, h.GetList)
	router.HandlerFunc(http.MethodGet, itemURL, h.GetOne)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	all, err := h.repository.FindAll(context.TODO())
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
	//uuid, err := strconv.Atoi(str[len(str)-1])
	h.logger.Info(str)
	//if err != nil {
	//	h.logger.Infof("%v", err)
	//	h.logger.Infof("uuidStr: %d\n", uuid)
	//}
	one, err := h.repository.FindOne(context.TODO(), str[len(str)-1])
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

func NewHandler(repo Repository, logger *logging.Logger) Handler {
	return &handler{
		repository: repo,
		logger:     logger,
	}
}
