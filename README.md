# Twitter feed back end

Uses docker-compose and Go

## Run

The project starts with one command (bash file) without installing anything except docker

``` bash
./start.sh
```

## Functionality

- Endpoint to add message
- Endpoint to get feed (get existing messages and stream new ones - use HTTP streaming)
- Implemented back pressure for message creation with Kafka
- Uses Cockroachdb with 3-node cluster
- Implemented a bot to generate messages. Configurable speed at docker-compose.yml

### Configuration

All configurations should be added to docker-compose.yml

The default configuration is as follows.
The server service starts on port 8888.
Cockroachdb cluster exposes 3 web-ui on ports 8081, 8082, 8083
 and SQL connections on ports 26257, 26258, and 26259.
Kafka service uses KRaft and is started without exposing any ports.
The message-consumer service connects to the Kafka and Cockroachdb clusters.
The message-bot service sends requests to the server every 15 seconds which creates a new message filled with lorem ipsum.

To configure message-bot speed change SENDING_TIMEOUT_SECONDS environmental variable for this service in docker-compose.yml. By default, it creates a new message every 15 seconds.

### API

#### Get message feed

- Default URL: "<http://localhost:8888/>"
- Method: "GET"
- Response: All existing messages in the form of separate JSON objects and new ones as soon as they are added. It uses the HTTP streaming concept to establish a connection and send the data

#### Add message

- Default URL: "<http://localhost:8888/>"
- Method: "POST"
- Request Body: A JSON object with the message details.
  - "username" (String): name of the user who sent the message
  - "message" (String): the message itself

##### Message feed testing

Utilities like cURL can help with testing how streaming works.

```powershell
curl.exe --no-buffer -v "http://localhost:8888/"
```

## License

This project is licensed under the **MIT** License. See the LICENSE file for details.
