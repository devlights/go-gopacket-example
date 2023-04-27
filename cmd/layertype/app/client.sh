#!/usr/bin/env bash

echo -n helloworld > .tmp
nc -N 127.0.0.1 22222 < .tmp

echo -n golang > .tmp
nc -N 127.0.0.1 22222 < .tmp

echo -n goroutine > .tmp
nc -N 127.0.0.1 22222 < .tmp

exit 0
