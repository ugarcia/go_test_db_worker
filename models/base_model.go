package models

import (
    "github.com/ugarcia/go_test_common/models"
    "github.com/ugarcia/go_test_db_worker/db"
)

type BaseModel struct {
    models.DbModel
}

/**
 * Initializes DB Worker if not done yet
 * Note: Worker is a db package global variable (in db/db.go)
 */
func (bm BaseModel) CheckDb() {
    if db.Worker == nil || &db.Worker.Db == nil || db.Worker.Ch == nil {
        db.Worker.Init()
    }
}
