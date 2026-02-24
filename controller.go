package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

var ErrPermissionNotGranted = fmt.Errorf("permission not granted")

type CommandOutputs struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type CustomerController struct {
	ac AuthorizationsClient
	db *DBClient
}

func NewCustomerController(ac AuthorizationsClient, db *DBClient) *CustomerController {
	return &CustomerController{ac: ac, db: db}
}

func (c *CustomerController) CreateCustomer(ctx context.Context, distributorID, distributorName string) (*Customer, *CommandOutputs, error) {
	addTrace(ctx, "controller")

	granted, err := c.ac.CheckCreateCustomerPermission(ctx, distributorID)
	if err != nil {
		return nil, nil, err
	}
	if !granted {
		return nil, nil, ErrPermissionNotGranted
	}

	customer, err := c.db.Customer.Create().
		SetName(distributorName).
		SetDistributorID(distributorID).
		Save(ctx)
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command("/bin/bash", "-c", "time sleep 0.01")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("command failed: %w (stderr=%s)", err, stderr.String())
	}

	if err := c.ac.AddCustomerAuthorizations(ctx, customer.ID, distributorID); err != nil {
		return nil, nil, err
	}

	return customer, &CommandOutputs{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}
