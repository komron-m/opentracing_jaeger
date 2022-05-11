<?php

namespace KomronM\OpentracingJaeger\Utils;

use OpenTracing\GlobalTracer;
use OpenTracing\Scope;
use SplStack;

class ScopeManager
{
    public static function startActiveSelfClosingSpan(string $operationName, $options = []): Scope
    {
        $scope = static::startActiveSpan($operationName, $options);
        static::deferClose($scope);
        return $scope;
    }

    public static function startActiveSpan(string $operationName, $options = []): Scope
    {
        return GlobalTracer::get()->startActiveSpan($operationName, $options);
    }

    public static function deferClose(Scope $scope)
    {
        static::defer(function () use ($scope) {
            $scope->close();
        });
    }

    // credits: https://github.com/php-defer/php-defer/blob/5.0/src/functions.inc.php
    private static function defer(callable $callback)
    {
        $context = new SplStack();
        $context->push(
            new class($callback) {
                private $callback;

                public function __construct(callable $callback)
                {
                    $this->callback = $callback;
                }

                public function __destruct()
                {
                    \call_user_func($this->callback);
                }
            }
        );
    }
}
