# Risala

gRPC based Message Broker written in Go for self learning purpose.

## Roadmap

![Roadmap image](https://github.com/melyouz/risala/blob/main/roadmap.svg?raw=true)

* v0.0.1 (MVP)
  + Broker
    * ~~Manage Exchanges (create, list, get, delete).~~
    * ~~Manage Queues (create, list, get, delete).~~
    * ~~Manage Bindings (create, delete).~~
    * ~~Publish messages to queue.~~
    * ~~Publish messages to exchange.~~
    * ~~Peek queue messages.~~
    * ~~Consume & acknowledge queue messages.~~
    * ~~Purge queue messages.~~
    * ~~Dead Letter queue.~~
    * ~~Process (get) a message.~~
    * ~~Acknowledge (ack) a message.~~
    * ~~Negatively acknowledge (nack) a message.~~
  + Consumer
    * Consume messages (HTTP API).
    * Acknowledge messages (HTTP API).
    * Negatively acknowledge messages (HTTP API).
  + Producer
    * Publish message (HTTP API).
* Further versions (TBD)
  + Broker
    * Retry mechanism.
    * Persistence layer (TBD).
    * More exchange types
      * Fanout: Route messages to all bound in queues (existing from MVP).
      * Direct: Route messages to bound in queues matching exact routing key (e.g. event.product.create.v1).
      * Topic: Route messages to bound in queues matching wildcard routing key (e.g. #, event.product.#,
        event.product.*.v1, ...).
    * More queue types
      * Transient: Temporary in memory queue (existing from MVP).
      * Durable: Persisted queue.
    * More binding types:
      * To queue (existing from MVP).
      * To exchange: Route messages to another exchange.
    * Tracing (e.g. Zipkin).
    * Logging (e.g. Vector + Grafana Loki, Datadog, ...).
  + Consumer
    * Consume messages (gRPC).
    * Acknowledge messages (gRPC).
    * Negatively acknowledge messages (gRPC).
    * Tracing.
    * Logging.
  + Producer
    * Publish messages (gRPC).
    * Tracing.
    * Logging.

## Suggestions

Any suggestion/recommendation? Just let me know! I appreciate it in advance 😊
