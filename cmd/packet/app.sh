#!/usr/bin/env bash

./packet &
echo -n $! > app.pid
exit 0
