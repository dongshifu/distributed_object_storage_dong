#!/bin/bash

echo -n "this object will have only 1 instance" | openssl dgst -sha256 -binary | base64

curl -v 10.29.2.1:12346/objects/test4_1 -XPUT -d "this object will have only 1 instance" -H "Digest: SHA-256=incorrecthash"

curl -v 10.29.2.1:12346/objects/test4_1 -XPUT -d "this object will have only 1 instance" -H "Digest: SHA-256=aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="

curl -v 10.29.2.1:12346/objects/test4_2 -XPUT -d "this object will have only 1 instance" -H "Digest: SHA-256=aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="

ls -ltr /tmp/?/objects
echo
curl 10.29.2.1:12346/objects/test4_1
echo
curl 10.29.2.1:12346/objects/test4_2
echo
curl 10.29.2.1:12346/locate/aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=
echo
curl 10.29.2.1:12346/versions/test4_1
echo
curl 10.29.2.1:12346/versions/test4_2
