#!/bin/sh
cd /usr/local/bin/

echo "Running URL shortening service..."
chmod +x /usr/local/bin/go-url-shortner

/usr/local/bin/go-url-shortner
