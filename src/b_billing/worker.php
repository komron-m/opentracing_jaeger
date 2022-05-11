<?php

require_once __DIR__ . '/vendor/autoload.php';

use Jaeger\Config;
use KomronM\OpentracingJaeger\EventHandler;
use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Message\AMQPMessage;

// in production these variables should be resolved from .env or some kind of configuration file
$rabbitmqHost = "localhost";
$rabbitmqPort = 5672;
$rabbitmqUser = "rabbitmq";
$rabbitmqPass = "secret";

// establish rabbitmq connection
$connection = new AMQPStreamConnection($rabbitmqHost, $rabbitmqPort, $rabbitmqUser, $rabbitmqPass);
$channel = $connection->channel();

// declare new queue
$queue = $channel->queue_declare("order_billing", false, true, false, false);
$channel->queue_bind("order_billing", "amq.topic", "#");

// in production these variables should be resolved from .env or some kind of configuration file
// initialize opentracing implementation client
$serviceName = "b_billing";
$samplerType = "const";
$samplerParam = 1.0;

(new Config(
    [
        "dispatch_mode" => Config::JAEGER_OVER_BINARY_UDP,
        "sampler" => [
            "type" => $samplerType,
            "param" => $samplerParam,
        ],
    ],
    $serviceName
))->initializeTracer();

// this callback is executed each time an event is received
$callback = function (AMQPMessage $msg) {
    echo ' [x] Received ', $msg->getRoutingKey(), "\n";
    (new EventHandler())->handle($msg);
};

$channel->basic_consume("order_billing", '', false, false, false, false, $callback);

echo " [*] Waiting for messages. To exit press CTRL+C\n";
while ($channel->is_open()) {
    $channel->wait();
}
