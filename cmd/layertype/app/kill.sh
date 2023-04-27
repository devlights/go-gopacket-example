#!/usr/bin/env bash

kill $(cat ./tcpdump.pid)
kill $(cat ./server.pid)
