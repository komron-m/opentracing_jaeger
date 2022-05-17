<?php

namespace KomronM\OpentracingJaeger\Billing;

class Repo
{
    public function getProductPrice(string $productId, $date): float
    {
        return floatval(random_int(1, 110));
    }

    public function storeBill(Bill $bill)
    {
        // ...
    }
}
