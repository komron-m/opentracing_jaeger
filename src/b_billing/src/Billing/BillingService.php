<?php

namespace KomronM\OpentracingJaeger\Billing;

use Exception;
use GuzzleHttp\Client;
use KomronM\OpentracingJaeger\Utils\ScopeManager;
use OpenTracing\GlobalTracer;
use PhpAmqpLib\Message\AMQPMessage;
use OpenTracing\Formats;

class BillingService
{
    public function __construct(private Repo $repo)
    {
    }

    public function serve(AMQPMessage $message)
    {
        $scope = ScopeManager::startActiveSpan("BillingService.serve");
        ScopeManager::deferClose($scope);

        $body = $message->body;
        $params = json_decode($body, true);

        $bill = $this->createBill($params);

        $this->processBill($bill);
    }

    public function createBill(array $params): Bill
    {
        ScopeManager::startActiveSelfClosingSpan("BillingService.createBill");

        $billID = uniqid();
        $orderID = $params["order_id"];
        $productId = $params["product_id"];
        $customerID = $params["customer_id"];
        $orderCreated = $params["created_at"];

        $productPrice = $this->repo->getProductPrice($productId, $orderCreated);

        $bill = new Bill(
            $billID,
            $orderID,
            $productId,
            $customerID,
            $productPrice,
        );

        $this->repo->storeBill($bill);
        return $bill;
    }

    public function processBill(Bill $bill)
    {
        $tracer = GlobalTracer::get();
        $scope = $tracer->startActiveSpan("BillingService.processBill");
        $span = $scope->getSpan();

        $headers = [];
        $tracer->inject(
            $span->getContext(),
            Formats\HTTP_HEADERS,
            $headers
        );

        $headers["Accept"] = "application/json";
        $headers["x-api-key"] = "secret";

        try {
            $client = new Client();
            $resp = $client->post("http://localhost:4001/process_bill", [
                "headers" => $headers,
                "json" => $bill->toArray(),
            ]);

            $span = $scope->getSpan();
            $span->setTag("status_code", $resp->getStatusCode());

        } catch (Exception $e) {
            $span = $scope->getSpan();
            $span->setTag("error", true);
            $span->log([
                "error" => $e->getMessage(),
                "file" => $e->getFile(),
                "line" => $e->getLine(),
            ]);

            // ... further error handling
        } finally {
            $scope->close();
        }
    }
}
