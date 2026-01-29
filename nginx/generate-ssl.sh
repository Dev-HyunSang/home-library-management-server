#!/bin/bash

# Generate self-signed SSL certificate for development
# For production, use Let's Encrypt or other CA

SSL_DIR="$(dirname "$0")/ssl"
mkdir -p "$SSL_DIR"

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout "$SSL_DIR/privkey.pem" \
  -out "$SSL_DIR/fullchain.pem" \
  -subj "/C=KR/ST=Seoul/L=Seoul/O=HomeLibrary/CN=myhomelibrary.japaneast.cloudapp.azure.com"

echo "SSL certificates generated in $SSL_DIR"
echo "  - fullchain.pem (certificate)"
echo "  - privkey.pem (private key)"
