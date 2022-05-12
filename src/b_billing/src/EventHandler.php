<?php

namespace KomronM\OpentracingJaeger;

use Exception;
use KomronM\OpentracingJaeger\Billing\BillingService;
use KomronM\OpentracingJaeger\Billing\Repo;
use OpenTracing\GlobalTracer;
use PhpAmqpLib\Message\AMQPMessage;
use PhpAmqpLib\Wire\AMQPTable;
use const OpenTracing\Formats\TEXT_MAP;

class EventHandler
{
    public function handle(AMQPMessage $msg)
    {
        // extract the span context
        /**@var AMQPTable $table */
        $table = $msg->get_properties()["application_headers"];
        $otData = $table->getNativeData()["opentracing_data"];

        // extract the span context
        $tracer = GlobalTracer::get();
        $spanContext = $tracer->extract(TEXT_MAP, json_decode($otData, true));
        $childOf = $spanContext ? ["child_of" => $spanContext] : [];

        // create reasonable span name
        $routingKey = $msg->getRoutingKey();
        $spanName = sprintf("EventHandler.handle: %s", $routingKey);

        $scope = $tracer->startActiveSpan($spanName, $childOf);
        $span = $scope->getSpan();
        $span->setTag("span.kind", "consumer");

        try {
            $service = static::resolve($msg);

            $service->serve($msg);

            $msg->ack();
        } catch (Exception $e) {
            $span->setTag("error", true);
            $span->log([
                "error" => $e->getMessage(),
                "file" => $e->getFile(),
                "line" => $e->getLine(),
            ]);

            $msg->nack();
        } finally {
            $scope->close();
            $tracer->flush();
        }
    }

    /**
     * @throws Exception
     * build dependency graph
     * in production, one should use some dependency injection (service container)
     */
    private function resolve(AMQPMessage $msg)
    {
        $THE_ONLY_TOPIC = "a_creator.order.created";
        $event = $msg->getRoutingKey();
        if ($event != $THE_ONLY_TOPIC) {
            throw new \Exception("Handler not found for event: $event");
        }
        $repo = new Repo();
        return new BillingService($repo);
    }
}
