package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"math/rand"
)

type ProcessBillRequest struct {
	BillID     string  `json:"bill_id"`
	OrderID    string  `json:"order_id"`
	ProductID  string  `json:"product_id"`
	CustomerID string  `json:"customer_id"`
	Price      float64 `json:"price"`
}

func ProcessBill(ctx context.Context, req *ProcessBillRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProcessBill")
	defer span.Finish()

	paymentMethod := getPreferredPaymentMethod(ctx, req.CustomerID)

	if err := withdraw(ctx, req.Price, req.CustomerID, paymentMethod); err != nil {
		ext.LogError(span, err)
		return err
	}

	return nil
}

// dummy func for demo purpose
func getPreferredPaymentMethod(ctx context.Context, customerID string) string {
	if rand.Int()%2 == 0 {
		return "CREDIT_CARD"
	}
	return "CRYPTO"
}

// dummy func for demo purpose
func withdraw(ctx context.Context, amount float64, customerID string, method string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "withdraw")
	defer span.Finish()

	span.LogFields(log.Float64("amount", amount))
	span.LogFields(log.String("payment_method", method))
	span.LogFields(log.String("customer_id", customerID))

	if amount > 100. {
		return fmt.Errorf("not enough balance")
	}
	return nil
}
