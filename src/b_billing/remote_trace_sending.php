<?php

require_once __DIR__ . '/vendor/autoload.php';

use Dotenv\Dotenv;
use Jaeger\Config;
use OpenTracing\GlobalTracer;

$dotenv = Dotenv::createUnsafeImmutable(__DIR__);
$dotenv->load();

$tracer = (new Config(
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

function main()
{
    echo "main fn called\n";

    $scope = GlobalTracer::get()->startActiveSpan("main");

    $span = $scope->getSpan();
    $span->log(["fn", "main"]);

    other();

    $scope->close();
}

function other()
{
    echo "other fn called\n";

    $scope = GlobalTracer::get()->startActiveSpan("other");

    $span = $scope->getSpan();

    $span->setTag("error", true);

    $scope->close();
}

main();

$tracer->flush();
