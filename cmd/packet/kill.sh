#!/usr/bin/env bash

kill $(cat ./app.pid)
kill $(cat ./server.pid)
