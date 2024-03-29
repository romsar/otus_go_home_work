#!/bin/bash

echo "running tests..."

docker-compose \
  -f ./deployments/docker-compose.test.yml \
  -p tests \
  --env-file .env \
  up --build \
  --exit-code-from integration_tests
