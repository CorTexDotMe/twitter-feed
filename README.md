# Twitter feed back end

Uses docker-compose and Go

## Run

Project starts with one command (bash file) without installing anything except docker

## Functionality

- Endpoint to add message
- Endpoint to get feed (get existing messages and stream new ones - use HTTP streaming)
- Implemented back pressure for message creation with Kafka
- Uses Cockroachdb with 2-node cluster
- Implemented a bot to generate messages. Configurable speed at .env
