#!/bin/bash

echo "################################################################################"
echo "Generating CA certificate"
echo "################################################################################"
openssl req -x509 \
    -newkey rsa:4096 \
    -days 365 -nodes \
    -keyout ca-key.pem \
    -out ca-cert.pem

echo "################################################################################"
echo "Generating server certificate"
echo "################################################################################"
openssl req \
    -newkey rsa:4096 -nodes \
    -keyout server-key.pem \
    -out server-req.pem

echo "subjectAltName=DNS:*.local,IP:0.0.0.0" > server-ext.cnf

echo "################################################################################"
echo "Sign server certificate using CA certificate"
echo "################################################################################"
openssl x509 -req \
    -in server-req.pem \
    -days 60 \
    -CA ca-cert.pem \
    -CAkey ca-key.pem \
    -CAcreateserial \
    -out server-cert.pem \
    -extfile server-ext.cnf
