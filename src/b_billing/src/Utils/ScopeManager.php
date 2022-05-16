<?php

namespace KomronM\OpentracingJaeger\Utils;

use OpenTracing\Scope;
use SplStack;

class ScopeManager
{
    // credits: https://github.com/php-defer/php-defer/blob/5.0/src/functions.inc.php
    public static function close(?SplStack &$context, callable $callback)
    {
        $context = $context ?? new SplStack();
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
