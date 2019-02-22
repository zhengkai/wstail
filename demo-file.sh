#!/bin/bash

FILE='/tmp/fortune.txt'

while true
do
	echo >> "$FILE"
	fortune tang300 | sed 's/\x1b\[[0-9;]*m//g' >> "$FILE"
	sleep 3
done
