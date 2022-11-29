#!/bin/sh
echo starting
echo host
echo $1
echo passwd
echo $2

echo PING | redis-cli -c -h $1 -p 6379 -a $2 -n 0

exit
