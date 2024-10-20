# Risala

gRPC based Message Broker written in Go for self learning purpose.

## Roadmap

### MVP

#### Broker

* [x] ~~Manage exchanges (fanout), queues (transient) & bindings (to queue).~~
* [x] ~~Publish messages.~~
* [x] ~~Retrieve messages.~~
* [ ] Consume & acknowledge messages.

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
