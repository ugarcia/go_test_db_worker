package main

import (
    "github.com/ugarcia/go_test_db_worker/db"
    "github.com/ugarcia/go_test_db_worker/queue"
)

/**
    Main execution
 */
func main() {

    // Init amqp connection
    go queue.Init()

    // Init DB connection
    // Note: Worker is a db package global variable (in db/db.go)
    db.Worker.Init()

    // Defer the DB connection close
    defer db.Worker.Db.Close()

    // Blocked loop by queries channel, will consume them when available
    for {
        go queue.HandleQueueMessage(<-db.Worker.Ch)
    }
}
