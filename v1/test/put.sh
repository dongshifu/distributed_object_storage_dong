#!/bin/bash

curl -v 10.29.2.1:12346/objects/test2 -XPUT -d"this is object test2"

curl 10.29.2.2:12346/locate/test2
echo
curl 10.29.2.1:12346/objects/test2
echo
