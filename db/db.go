package db

import (
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"

    "github.com/ugarcia/go_test_common/models"
)

// Global reference to Worker
type DbWorker struct {
    Db gorm.DB
    Ch chan models.DbQueueMessage
}
var Worker *DbWorker = new(DbWorker)

/**
 * DB and Channel initialization
 */
func (w *DbWorker) Init() {

    // Error variable needed
    var err error

    // Connect to DB
    w.Db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/gw_core?timeout=10m&parseTime=true")
    if  err != nil {
        log.Fatal(err.Error())
        os.Exit(1)
    }

    // Init underlying DB, then we could invoke `*sql.DB`'s functions with it
    w.Db.DB()

    // Check DB with a ping
    if err = w.Db.DB().Ping(); err != nil {
        log.Fatal(err.Error())
        os.Exit(1)
    }

    // Just in case set to default OS/DB max connections
    w.Db.DB().SetMaxOpenConns(0)

    // Init the channel for queries
    w.Ch = make(chan models.DbQueueMessage)
}
