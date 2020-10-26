#!/bin/sh
ab -n100000 -c50 -k -H "Content-Type: application/json" -p single-event.json http://localhost:5050/s
