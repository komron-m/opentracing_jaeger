package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"sync"
)

type Repo struct {
	orders []*Order
	sync.RWMutex
}

func NewRepo() *Repo {
	return &Repo{
		orders: []*Order{},
	}
}

func (r *Repo) Store(ctx context.Context, o *Order) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Repo.Store")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	r.orders = append(r.orders, o)

	return nil
}
