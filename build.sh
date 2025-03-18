#!/bin/bash

(docker build -t message-server -f ./cmd/message-server/Dockerfile .) &
(docker build -t message-consumer -f ./cmd/message-consumer/Dockerfile .) &
(docker build -t message-bot -f ./cmd/message-bot/Dockerfile .) &