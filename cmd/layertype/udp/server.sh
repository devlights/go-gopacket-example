#!/usr/bin/env bash

nc -u -k -l 22222 &
echo -n $! > server.pid
exit 0