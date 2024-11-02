# Risala

gRPC based Message Broker written in Go for self learning purpose.

## Roadmap

### MVP

#### Broker

* [x] ~~Manage Exchanges (create, list, get, delete).~~
* [x] ~~Manage Queues (create, list, get, delete).~~
* [x] ~~Manage Bindings (create, delete).~~
* [x] ~~Publish messages to queue.~~
* [x] ~~Publish messages to exchange.~~
* [x] ~~Peek queue messages.~~
* [x] ~~Consume & acknowledge queue messages.~~
* [x] ~~Purge queue messages.~~
* [ ] Dead Letter queue.

#### Producer

* [ ] Publish messages (gRPC, in memory).

#### Consumer

* [ ] Consume & acknowledge messages (gRPC, in memory).

### Further versions

* [ ] Retry mechanism.
* [ ] Persistence layer (TBD).
* [ ] More exchange types
    * [ ] Fanout: Route messages to all bound in queues (existing from MVP).
    * [ ] Direct: Route messages to bound in queues matching exact routing key (e.g. event.product.create.v1).
    * [ ] Topic: Route messages to bound in queues matching wildcard routing key (e.g. #, event.product.#,
      event.product.*.v1, ...).
* [ ] More queue types
    * [ ] Transient: Temporary in memory queue (existing from MVP).
    * [ ] Durable: Persisted queue.
* [ ] More binding types:
    * [ ] To queue (existing from MVP).
    * [ ] To exchange: Route messages to another exchange.
* [ ] Tracing (e.g. Zipkin).
* [ ] Logging (e.g. Vector + Grafana Loki, Datadog, ...).

## Suggestions

Any suggestion/recommendation? Just let me know! I appreciate it in advance 😊
