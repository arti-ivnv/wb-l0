package handlers

import (
	"encoding/json"
	"fmt"
	"ls-0/arti/order/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

type saferManager interface {
	Get(order_uuid string) (storage.Order, bool)
}

type OrdersHandler struct {
	manager saferManager
}

func NewOrderHandler(m saferManager) *OrdersHandler {
	return &OrdersHandler{
		manager: m,
	}
}

func (o OrdersHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	order_uid := vars["order_uid"]

	order, ok := o.manager.Get(order_uid)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Order does not exist"))
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Error marshaling string:", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
