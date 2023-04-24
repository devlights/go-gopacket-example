#!/usr/bin/env bash

echo helloworld > .tmp
nc -u -N 127.0.0.1 22222 < .tmp &

exit 0
