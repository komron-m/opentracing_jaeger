<?php

require_once __DIR__ . '/vendor/autoload.php';

use Dotenv\Dotenv;
use Jaeger\Config;
use KomronM\OpentracingJaeger\EventHandler;
use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Message\AMQPMessage;

$dotenv = Dotenv::createUnsafeImmutable(__DIR__);
$dotenv->load();

// establish rabbitmq connection
$rabbitmqHost = getenv("RABBITMQ_HOST");
$rabbitmqPort = getenv("RABBITMQ_PORT");
$rabbitmqUser = getenv("RABBITMQ_USER");
$rabbitmqPass = getenv("RABBITMQ_PASSWORD");

$connection = new AMQPStreamConnection($rabbitmqHost, $rabbitmqPort, $rabbitmqUser, $rabbitmqPass);
$channel = $connection->channel();

// declare new queue
$queueName = getenv("RABBITMQ_QUEUE_NAME");
$exchangeName = getenv("RABBITMQ_EXCHANGE_NAME");
$queue = $channel->queue_declare($queueName, false, true, false, false);
$channel->queue_bind($queueName, $exchangeName, "#");

(new Config(
    [
        "logging" => getenv("JAEGER_LOGS_ENABLED"),
        "dispatch_mode" => Config::JAEGER_OVER_BINARY_UDP,
        "sampler" => [
            "type" => getenv("JAEGER_SAMPLER_TYPE"),
            "param" => getenv("JAEGER_SAMPLER_PARAM"),
        ],
    ],
    getenv("JAEGER_SERVICE_NAME")
))->initializeTracer();

// this callback is executed each time an event is received
$callback = function (AMQPMessage $msg) {
    echo ' [x] Received ', $msg->getRoutingKey(), "\n";
    (new EventHandler())->handle($msg);
};

$channel->basic_consume($queueName, '', false, false, false, false, $callback);

echo " [*] Waiting for messages. To exit press CTRL+C\n";
while ($channel->is_open()) {
    $channel->wait();
}
