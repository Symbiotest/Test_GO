package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RouteHandler struct {
	controller *CustomerController
}

func NewRouteHandler(controller *CustomerController) *RouteHandler {
	return &RouteHandler{controller: controller}
}

type CreateCustomerResponse struct {
	Customer *Customer       `json:"customer"`
	Command  *CommandOutputs `json:"command"`
	Trace    []string        `json:"trace"`
}

func (h *RouteHandler) CreateCustomer(w http.ResponseWriter, r *http.Request, body CreateCustomerJSONRequestBody) {
	addTrace(r.Context(), "route-handler")

	customer, cmdOutputs, err := h.controller.CreateCustomer(r.Context(), body.DistributorID, body.DistributorName)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrPermissionNotGranted) {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	resp := CreateCustomerResponse{
		Customer: customer,
		Command:  cmdOutputs,
		Trace:    traceSteps(r.Context()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
