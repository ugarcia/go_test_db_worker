# go_test_db_worker
Testing Go Capabilities

06.01.2016:
-----------

- A DB Worker which does:
    - Maintain a connections pool during the service lifecycle
    - Consume job messages from DB Worker Queue and query the DB based on messages data
    - Publish job responses to WS (Websockets) Queue
    Pending:
        - Publish job responses to more Queues (like API Queue ...)
        - Pub/Sub? Only if relevant for our use cases ...
        - Check ORM performance against native sql/database package (which I think should be enough)

Easiest way to test it is running the process and browsing the html examples

    go run src/github.com/ugarcia/go_test_db_worker/main.go

Dependencies:
-------------
- RabbitMq >= 3.5.1
- MySQL >= 5.5.x
- Golang >= 1.5.x (and proper system config for $GOPATH and $GOROOT)
- Package github.com/ugarcia/go_test_common
