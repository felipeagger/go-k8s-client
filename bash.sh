#!/bin/sh
echo starting
echo host
echo $1
echo passwd
echo $2

redis-cli -c -h $1 -p 6379 -a $2 --bigkeys -i 3 > out.log

exit
