<?php

namespace KomronM\OpentracingJaeger\Billing;

class Bill
{
    public function __construct(
        public string $billID,
        public string $orderID,
        public string $productID,
        public string $customerID,
        public float  $price,
    )
    {
    }

    public function toArray()
    {
        return [
            "bill_id" => $this->billID,
            "order_id" => $this->orderID,
            "product_id" => $this->productID,
            "customer_id" => $this->customerID,
            "price" => $this->price,
        ];
    }
}
