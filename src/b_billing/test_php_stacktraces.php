<?php

function main()
{
    first();
    third();
}

function first()
{
    printf("second\n");
    second();
}

function second()
{
    printf("third\n");
}

function third()
{
    printf("third\n");
    fourth();
}

function fourth()
{
    throw new \Exception("fourth");
}

main();
