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
        $scope = GlobalTracer::get()->startActiveSpan("BillingService.serve");
        ScopeManager::close($_, function () use ($scope) {
            $scope->close();
        });

        $body = $message->body;
        $params = json_decode($body, true);

        $bill = $this->createBill($params);

        $this->processBill($bill);
    }

    public function createBill(array $params): Bill
    {
        $scope = GlobalTracer::get()->startActiveSpan("BillingService.createBill");
        ScopeManager::close($_, function () use ($scope) {
            $scope->close();
        });

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

        try {
            $client = new Client();

            $headers = [];
            $tracer->inject(
                $span->getContext(),
                Formats\HTTP_HEADERS,
                $headers
            );

            $headers["Accept"] = "application/json";
            $headers["x-api-key"] = getenv("BILL_PROCESSOR_SECRET");

            $resp = $client->post(getenv("BILL_PROCESSOR_URL"), [
                "headers" => $headers,
                "json" => $bill->toArray(),
            ]);

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
