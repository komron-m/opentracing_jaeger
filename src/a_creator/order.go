package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"time"
)

const (
	EventOrderCreated = "a_creator.order.created"
)

type Order struct {
	OrderID    string    `json:"order_id"`
	ProductID  string    `json:"product_id"`
	CustomerID string    `json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	ProductID  string `json:"product_id"`
	CustomerID string `json:"customer_id"`
}

func CreateOrder(
	ctx context.Context,
	req *CreateOrderRequest,
	repo *Repo,
	p *Publisher,
) (*Order, error) {
	// request validation skipped for simplicity
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateOrder")
	defer span.Finish()

	// make new Order instance
	order := new(Order)
	order.OrderID = uuid.NewString()
	order.ProductID = req.ProductID
	order.CustomerID = req.CustomerID
	order.CreatedAt = time.Now()

	// store in database
	if err := repo.Store(ctx, order); err != nil {
		return nil, err
	}

	// emit event
	if err := p.Publish(ctx, EventOrderCreated, order); err != nil {
		return nil, err
	}

	return order, nil
}
