#!/bin/sh

# This script generates a self signed cert for your zero-http server

# https://gist.github.com/Soarez/9688998
# https://github.com/netty/netty/blob/4.1/handler/src/test/resources/io/netty/handler/ssl/generate-certs.sh#L14\

SUBJECT=${1:-"/CN=localhost"}
DAYS=3650

CA_KEY="zero-http_ca.key"
CA_CERT="zero-http_ca.pem"

LOCAL_KEY="zero-http.key"
LOCAL_SIGNED_CERT="zero-http.pem"

ALIAS="zero-http"

# Generate a new, self-signed root CA
openssl req -extensions v3_ca -new -x509 -days $DAYS -nodes -subj "/CN=${ALIAS}" -newkey rsa:2048 -sha512 -out $CA_CERT -keyout $CA_KEY

LOCAL_KEY_TEMP="temp.key"
# Generate a certificate/key for the server to use for Hostname Verification via localhost
openssl req -new -keyout $LOCAL_KEY_TEMP -nodes -newkey rsa:2048 -subj $SUBJECT | \
    openssl x509 -req -CAkey $CA_KEY -CA $CA_CERT -days $DAYS -set_serial $RANDOM -sha512 -extfile v3.ext -out $LOCAL_SIGNED_CERT
openssl pkcs8 -topk8 -inform PEM -outform PEM -in $LOCAL_KEY_TEMP -out $LOCAL_KEY -nocrypt
rm $LOCAL_KEY_TEMP

# remove the CA key since we should never need it again
rm $CA_KEY