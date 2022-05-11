package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"time"
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
	o := new(Order)
	o.OrderID = uuid.NewString()
	o.ProductID = req.ProductID
	o.CustomerID = req.CustomerID
	o.CreatedAt = time.Now()

	// store in database
	if err := repo.Store(ctx, o); err != nil {
		return nil, err
	}

	// emit event
	if err := p.Publish(ctx, "a_creator.order.created", o); err != nil {
		return nil, err
	}

	return o, nil
}
