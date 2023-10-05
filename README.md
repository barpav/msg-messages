# Messages (messenger)

This is messaging component of the [Messenger](https://github.com/barpav/messenger) pet-project. <br>

## Functions

* Message CRUD operations.

* Storing message data (PostgreSQL).

* Message syncing (including between multiple clients). 

* Updating file usage statistics (RabbitMQ).

See microservice [REST API](https://barpav.github.io/msg-api-spec/#/messages) and [deployment diagram](https://github.com/barpav/messenger#deployment-diagram) for details.

## Message syncing

> Note: since this is an educational project and the main client app is Swagger UI (or another tool like curl), the service does not currently provide any server-to-client notifications (like WebSocket API or others). This fact doesn't affect logic described below, because server-to-client notifications should be considered as trigger "something changed, time to make synchronization".

Before synchronization, it needs to be clarified *what* exactly is syncing. The service has its global timeline, and every [create](https://barpav.github.io/msg-api-spec/#/messages/post_messages), [modify](https://barpav.github.io/msg-api-spec/#/messages/patch_messages__id_) or [delete](https://barpav.github.io/msg-api-spec/#/messages/delete_messages__id_) operation result a new point on that timeline occur - called `timestamp`. Timestamps are integers, so they can be easily compared on "happened earlier" and "happened later". All timestamps are unique, so "happened at the same time" situation is impossible. Thus, every message in the service has two attributes:

* `id` - message identifier, unique among all messages in the service, assigned by the service only once when the message is initially created.

* `timestamp` - message modification point in time, unique among all timestamps in the service, assigned by the service each time this message is created, updated or deleted.

So timestamps reflects *changes*, and *synchronization* is just *receiving all changes that happened after the last known change* (or from the beginning in case of the first sync). <br>

At the same time, the service does not store anything but timestamps from previous changes, so there is no such thing as "previous versions" of the message: the message data is always as it was last time modified (or created) no matter how much changes (timestamps) were made. To provide data integrity and consistency while changing messages, [modify](https://barpav.github.io/msg-api-spec/#/messages/patch_messages__id_) and [delete](https://barpav.github.io/msg-api-spec/#/messages/delete_messages__id_) operations require correct current message timestamp.

### GET /messages

To receive changes, the service provides [GET /messages](https://barpav.github.io/msg-api-spec/#/messages/get_messages) operation, that has `after` optional query parameter, which is typically timestamp of the last successfully synced message that client must store (single value).

Received changes must be processed in the same order which they were returned (sorted by `timestamp` in ascendent order) and each successfully processed `timestamp` becomes new `after` parameter, whereas message data should only be [received](https://barpav.github.io/msg-api-spec/#/messages/get_messages__id_):

* If client has no message with specified `id`.
* If client already has a message with specified `id` _and_ received `timestamp` is _bigger_.

> **Important:** timestamp of the last successfully synced message that client store (`after` parameter) must be received from [GET /messages](https://barpav.github.io/msg-api-spec/#/messages/get_messages) operation, not [get message data](https://barpav.github.io/msg-api-spec/#/messages/get_messages__id_) in order to ensure correct pagination while syncing.
