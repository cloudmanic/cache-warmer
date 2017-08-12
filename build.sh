#!/bin/bash

cd src

env GOOS=linux GOARCH=amd64 go build -o ../build/cache-warm

cd ../

docker-compose build