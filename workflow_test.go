package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestIndirectWorkflow_EndToEnd(t *testing.T) {
	server := httptest.NewServer(NewHTTPHandler())
	defer server.Close()

	reqBody := CreateCustomerJSONRequestBody{
		DistributorID:   "dist-42",
		DistributorName: "ACME Distribution",
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(server.URL+"/customers", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("post /customers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status code: got=%d want=%d", resp.StatusCode, http.StatusCreated)
	}

	var got CreateCustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Customer == nil {
		t.Fatal("customer is nil")
	}
	if got.Customer.Name != "ACME Distribution" {
		t.Fatalf("unexpected customer name: got=%q", got.Customer.Name)
	}
	if got.Customer.DistributorID != "dist-42" {
		t.Fatalf("unexpected distributorID: got=%q", got.Customer.DistributorID)
	}
	if got.Customer.ID == "" {
		t.Fatal("customer id is empty")
	}

	if got.Command == nil {
		t.Fatal("command result is nil")
	}
	if got.Command.Stderr == "" {
		t.Fatal("expected timing output on stderr from `time` command")
	}

	wantTrace := []string{
		"http-server",
		"chi",
		"oapi-handler",
		"oapi-middleware",
		"route-handler",
		"controller",
	}
	if !reflect.DeepEqual(got.Trace, wantTrace) {
		t.Fatalf("unexpected trace order:\n got: %#v\nwant: %#v", got.Trace, wantTrace)
	}
}
