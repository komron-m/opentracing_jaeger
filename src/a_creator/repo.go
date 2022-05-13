package main

import (
	"context"
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
	r.Lock()
	defer r.Unlock()

	r.orders = append(r.orders, o)

	return nil
}
