#!/usr/bin/env bash

nc -l -k 127.0.0.1 22222 1>/dev/null &
echo -n $! > server.pid
exit 0
