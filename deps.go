package main

import (
	"context"
	"fmt"
	"sync"
)

type Customer struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	DistributorID string `json:"distributorId"`
}

type AuthorizationsClient interface {
	CheckCreateCustomerPermission(ctx context.Context, distributorID string) (bool, error)
	AddCustomerAuthorizations(ctx context.Context, customerID, distributorID string) error
}

type FakeAuthorizationsClient struct{}

func (f *FakeAuthorizationsClient) CheckCreateCustomerPermission(_ context.Context, distributorID string) (bool, error) {
	return distributorID != "", nil
}

func (f *FakeAuthorizationsClient) AddCustomerAuthorizations(_ context.Context, _ string, _ string) error {
	return nil
}

type DBClient struct {
	Customer *CustomerStore
}

func NewDBClient() *DBClient {
	return &DBClient{
		Customer: &CustomerStore{
			nextID: 1,
		},
	}
}

type CustomerStore struct {
	mu     sync.Mutex
	nextID int
	items  []Customer
}

func (s *CustomerStore) Create() *CustomerCreateBuilder {
	return &CustomerCreateBuilder{store: s}
}

type CustomerCreateBuilder struct {
	store         *CustomerStore
	name          string
	distributorID string
}

func (b *CustomerCreateBuilder) SetName(name string) *CustomerCreateBuilder {
	b.name = name
	return b
}

func (b *CustomerCreateBuilder) SetDistributorID(distributorID string) *CustomerCreateBuilder {
	b.distributorID = distributorID
	return b
}

func (b *CustomerCreateBuilder) Save(_ context.Context) (*Customer, error) {
	if b.name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if b.distributorID == "" {
		return nil, fmt.Errorf("distributorID is required")
	}

	b.store.mu.Lock()
	defer b.store.mu.Unlock()

	customer := Customer{
		ID:            fmt.Sprintf("cus-%04d", b.store.nextID),
		Name:          b.name,
		DistributorID: b.distributorID,
	}
	b.store.nextID++
	b.store.items = append(b.store.items, customer)

	c := customer
	return &c, nil
}
