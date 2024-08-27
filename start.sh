#!/bin/bash

docker compose up -d

docker compose exec roach1 ./cockroach init --insecure --host=roach1:26357

echo "Connecting to CockroachDB cluster..."

until docker compose exec roach1 ./cockroach sql --insecure --host=roach2:26258 --execute="CREATE DATABASE IF NOT EXISTS twitterdb;" > /dev/null 2>&1
do
  echo "Waiting for CockroachDB cluster to be ready..."
  sleep 2
done

echo "Database 'twitterdb' has been created successfully(if not existed)"