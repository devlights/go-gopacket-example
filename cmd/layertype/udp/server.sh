#!/usr/bin/env bash

nc -u -k -l 22222 1>/dev/null &
echo -n $! > server.pid
exit 0