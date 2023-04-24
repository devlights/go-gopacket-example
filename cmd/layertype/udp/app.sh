#!/usr/bin/env bash

./udp &
echo -n $! > app.pid
exit 0
