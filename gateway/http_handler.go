package main

import (
	"github.com/hawful70/common"
	pb "github.com/hawful70/common/api"
	"net/http"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
}

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerId,
		Items:      items,
	})
	if err != nil {
		return
	}
}
