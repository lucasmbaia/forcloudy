#!/bin/bash

c=("curl http://192.168.204.128:32769")
d=("curl http://192.168.204.128:32773")

while :
do
	eval $c
	eval $d
done
