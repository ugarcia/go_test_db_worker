package queue

import (
    "fmt"
    "encoding/json"

    "github.com/streadway/amqp"

    dbModels "github.com/ugarcia/go_test_db_worker/models"
    "github.com/ugarcia/go_test_common/models"
    "github.com/ugarcia/go_test_db_worker/db"
    "github.com/ugarcia/go_test_common/mq"
)

// Global variables for queue/channel etc.
var q *mq.AMQP

// Constants for this module
const MQ_URL = "amqp://guest:guest@mq.gamewheel.local:5672/"
const QUEUE = "db_worker_q"
const ROUTE = "workers.db"
const ID = "workers.db"
const EXCHANGE = "workers"
var EXCHANGES = []string{"modules", "workers"}

/**
 * Initialize Queues, Channels and Consumer Loops
 */
func Init() {

    // Create struct
    q = new(mq.AMQP)

    // Init it
    q.Init(MQ_URL)

    // Defer closing
    defer q.Close()

    // Register exchanges
    q.RegisterExchanges(EXCHANGES)

    // Register queues
    q.RegisterQueues([]string{QUEUE})

    // Bind queues
    q.BindQueuesToExchange([]string{QUEUE}, EXCHANGE, ROUTE)

    // Start consuming
    q.Consume(QUEUE, receiveQueueMessage)
}

/**
 * Receives a message from DB Queue and adds it to queries channel
 */
func receiveQueueMessage(msg models.QueueMessage, d amqp.Delivery) {

    // Add to queries channel
    db.Worker.Ch <- msg

    // TODO: Pass this delivery object along ans send ACK only after finishing everything???
    d.Ack(false)
}

/**
 * Handles a received message from DB Queue, sending relevant info to connected client
 */
func HandleQueueMessage(msg models.QueueMessage) {

    // Validations
    // TODO: Abstract this please
    if msg.Target == "" || msg.Target != "data" {
        fmt.Println("Invalid or missing target defined in message!");
        return
    }

    if msg.Code == "" {
        fmt.Println("No entity/code defined in message!");
        return
    }

    if msg.Action == "" {
        fmt.Println("No action defined in message!");
        return
    }

    if msg.Data == nil {
        fmt.Println("No data provided in message!");
        return
    }

    // Variables for query result
    var result interface{}
    var err error

    // Check what type of entity we're dealing with
    switch msg.Code {
        case "game":
            err, result = dbModels.Games{}.HandleMessage(msg.BaseMessage)
        default:
            fmt.Printf("Unknown entity: %s\n", msg.Code)
            return
    }

    // Error from query
    if err != nil {
        fmt.Println(err.Error());
        return
    }

    // Parse result from relevant interface model to bytes
    data, err := json.Marshal(result)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Parse previous bytes to generic map
    var outData = make(map[string]interface{})
    if err := json.Unmarshal(data, &outData); err != nil {
        fmt.Println(err.Error())
        return
    }

    // Build base response message and init final one
    msg.Receiver = msg.Sender
    msg.Sender = ID
    msg.Data = outData

    // Send it
    q.SendMessage(msg)
}