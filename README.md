# Outbox

Outbox is a simple example of a message broker that uses the transactional outbox pattern to send messages to Kafka.

## Installation & Running

### Docker

1. Start the services:

   ```sh
   docker compose up
   ```

2. Test the health endpoint:

   ```sh
   curl 'http://127.0.0.1/health'
   ```

### Manual

1. Load default environment variables from `example.env`.

2. Prepare Postgres and Kafka, then provide the connection details via environment variables:

   ```sh
   export OUTBOX_KAFKA_BROKERS=localhost:9094
   export OUTBOX_POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
   ```

3. Provision Postgres and Kafka:

    ```sh
    go run ./cmd/postgres-up
    go run ./cmd/kafka-up
    ```

4. Start worker in the background:

    ```sh
    go run ./cmd/worker &
    ```

5. Start server:

    ```sh
    go run ./cmd/server
    ```

6. Test the health endpoint:

    ```sh
    curl 'http://127.0.0.1/health'
    ```

## Usage

### `POST /messages`

Creates a new message, saves information about it to a table and sends it to Kafka via the outbox table.

The request is unmarshalled into the following structure:

```go
type createMessageRequest struct {
	Topic   string                       `json:"topic"`
	Key     string                       `json:"key"`
	Value   string                       `json:"value"`
	Headers []createMessageHeaderRequest `json:"headers"`
}

type createMessageHeaderRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
```

The topic must exist in the Kafka cluster, otherwise the worker will fail to send the message. During provisioning,
`kafka-up` creates a topic named `example`.

Example:

```sh
curl -X POST 'http://127.0.0.1:8080/messages' \
  -H 'Content-Type: application/json' \
  -d '{ "topic": "example", "key": "a-key", "value": "a-value", "headers": [{ "key": "Content-Type", "value": "application/json" }] }'
```

### `GET /statistics`

Returns statistics about processed messages.

Example:

```sh
curl 'http://127.0.0.1:8080/statistics'
```
