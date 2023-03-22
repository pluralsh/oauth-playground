#!/bin/sh

openssl genrsa 4096 > localCA.key #generate CA key
openssl req -x509 -new -nodes -key localCA.key -days 1000 -sha256 -subj '/CN=YourCompany' > localCA.pem #generate CA cert
openssl genrsa 2048 > local.key #generate server key
openssl req -new -subj '/CN=YourCompany' -key local.key > local.csr #generate signing request
openssl x509 -req -days 1000 -sha256 -CA localCA.pem -CAkey localCA.key -CAcreateserial -in local.csr -extfile cert.conf -extensions local_san > local.pem #sign request with local CA

#limit potential for 🔥

rm localCA.key
rm localCA.srl
rm local.csr
chmod 640 local.key