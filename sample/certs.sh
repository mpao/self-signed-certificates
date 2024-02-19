#!/bin/bash

# CA cert
openssl req -new -newkey rsa:2048 -days 3650 -extensions v3_ca -nodes -x509 -sha256 -set_serial 0 \
-keyout localCA.key -out localCA.crt -subj "/CN=RootCA/"
# Certificate Signing Request
openssl req -new -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr \
-config request.cnf -extensions v3_req
# Self-Signing
openssl x509 -req -sha256 -CAcreateserial -in localhost.csr -days 3650 -CA localCA.crt \
-CAkey localCA.key -out localhost.crt -extfile <(printf "subjectAltName=DNS:localhost")

rm localhost.csr localCA.key localCA.srl
