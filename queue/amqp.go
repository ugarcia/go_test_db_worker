package queue

import (
    "fmt"
    "log"
    "encoding/json"

    "github.com/streadway/amqp"

    dbModels "github.com/ugarcia/go_test_db_worker/models"
    "github.com/ugarcia/go_test_common/models"
    "github.com/ugarcia/go_test_db_worker/db"
)

// Global reference to WS queue/channel
var wsQueue amqp.Queue
var wsChannel *amqp.Channel

/**
 * Helper function for handling errors
 */
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

/**
 * Initialize Queues, Channels and Consumer Loops
 */
func Init() {

    // Create connection for WS
    wsConn, err := amqp.Dial("amqp://guest:guest@mq.gamewheel.local:5672/")
    failOnError(err, "Failed to connect to RabbitMQ for WS")
    defer wsConn.Close()

    // Create WS Channel
    wsChannel, err = wsConn.Channel()
    failOnError(err, "Failed to open WS channel")
    defer wsChannel.Close()

    // Create WS Queue
    wsQueue, err = wsChannel.QueueDeclare(
        "ws", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare WS queue")

    // Create connection for DB
    dbConn, err := amqp.Dial("amqp://guest:guest@mq.gamewheel.local:5672/")
    failOnError(err, "Failed to connect to RabbitMQ for DB")
    defer dbConn.Close()

    // Create DB Channel
    dbChannel, err := dbConn.Channel()
    failOnError(err, "Failed to open DB channel")
    defer dbChannel.Close()

    // Create DB Queue
    dbQueue, err := dbChannel.QueueDeclare(
        "db", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare DB queue")

    err = dbChannel.Qos(
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )
    failOnError(err, "Failed to set QoS")

    // Consume DB Channel messages
    deliveries, err := dbChannel.Consume(
        dbQueue.Name, // queue
        "", // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil, // args
    )
    failOnError(err, "Failed to register WS consumer")

    // Loop forever for messages, and call handler for each
    forever := make(chan bool)
    go func() {
        for d := range deliveries {
            log.Printf("Received a message: %s", d.Body)
            d.Ack(false)
            go receiveDbQueueMessage(d)
        }
    }()
    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

/**
 * Receives a message from DB Queue and adds it to queries channel
 */
func receiveDbQueueMessage(msg amqp.Delivery) {
    // TODO: Check delivery parameters first?
    data := models.DbQueueMessage{}
    if err := json.Unmarshal([]byte(msg.Body), &data); err != nil {
        fmt.Println(err.Error())
        return
    }
    addQuery(data)
}

/**
 * Sends a message to DB Queue
 */
func SendWsQueueMessage(msg []byte) {
    sendQueueMessage(msg, wsChannel, wsQueue.Name)
}

/**
 * Sends a message to a Queue
 */
func sendQueueMessage(msg []byte, ch *amqp.Channel, queue string) {
    err := ch.Publish(
        "",     // exchange
        queue, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType: "application/json",
            Body:        []byte(msg),
        })
    failOnError(err, "Failed to publish a message")
    log.Printf(" [x] Sent %s", string(msg))
}

/**
 * Query addition to channel
 * Note: Worker is a db package global variable (in db/db.go)
 */
func addQuery(q models.DbQueueMessage) {
    db.Worker.Ch <- q
}

/**
 * Handles a received message from DB Queue, sending relevant info to connected client
 */
func HandleDbQueueMessage(msg models.DbQueueMessage) {

    // Validations
    if msg.Entity == "" {
        fmt.Println("No entity defined in queue message!");
        return
    }

    if msg.Action == "" {
        fmt.Println("No action defined in queue message!");
        return
    }

    if msg.Sender == "" {
        fmt.Println("No sender defined in queue message!");
        return
    }

    if msg.Data == nil {
        fmt.Println("No data provided in queue message!");
        return
    }

    // Variables for query result
    var result interface{}
    var err error

    // Check what type of entity we're dealing with
    switch msg.Entity {
        case "game":
            err, result = dbModels.Games{}.HandleQueueMessage(msg)
        default:
            fmt.Printf("Unknown entity: %s\n", msg.Entity)
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
    baseOutMsg := models.BaseMessage {
        Sender: "db",
        ConnectionId: msg.ConnectionId,
        Action: msg.Action,
        Data: outData,
    }

    // Check if data affected, so we need to broadcast response to all clients
    switch msg.Action {
        case "post", "delete", "update":
            baseOutMsg.Broadcast = true
    }

    // Check who was the sender, so we fill additional fields
    var outMsg interface{}
    switch msg.Sender {
        case "ws":
            outMsg = models.WsQueueMessage{
                BaseMessage: baseOutMsg,
                Code: msg.Entity,
            }
    }

    // Encode message
    resp, err := json.Marshal(outMsg)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Check who was the sender, so we send the response back
    switch msg.Sender {
        case "ws":
            SendWsQueueMessage(resp)
        default:
            fmt.Printf("Unknown sender: %s\n", msg.Sender)
            return
    }
}