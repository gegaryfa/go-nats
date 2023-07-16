# go-nats
Example code of using the [NATS](https://nats.io/) pub/sub with Go

The example is demonstrating how to publish messages in subjects and subscribe to them either synchronously or asynchronously.

## Running nats in docker
It is required to run NATS before moving on to run the example code.
The easiest way to run NATS is to run it in docker with the command below:

```bash
docker run --name nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222
```
