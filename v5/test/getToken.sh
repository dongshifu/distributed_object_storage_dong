#!/bin/bash

curl -v 10.29.2.1:12346/objects/test6 -XPOST -H "Digest: SHA-256=$1" -H "Size: 100000" #$1为要上传/访问的对象的hash值