package models

import (
    "fmt"
    "errors"
    "strings"
    "database/sql"

    "github.com/ugarcia/go_test_common/models"
    "github.com/ugarcia/go_test_db_worker/db"
)

type Game struct {
    BaseModel
    Name            string          `form:"name" json:"name"`
    Code            string          `form:"code" json:"code"`
    Desc            string          `form:"desc" json:"desc"`
    Ordering        uint            `sql:"default:0" json:"ordering"`
    Version         string          `sql:"default:1.0.0" json:"version"`
    GameTypeId      sql.NullInt64   `form:"game_type_id" json:"game_type_id"`
}

type Games []Game

func (games Games) All() Games {
    Game{}.CheckDb()
    db.Worker.Db.Find(&games)
    return games
}

func (game Game) Add() Game {
    game.CheckDb()
    db.Worker.Db.Create(&game)
    return game
}

func (game Game) Delete() Game {
    game.CheckDb()
    db.Worker.Db.Delete(&game)
    return game
}

func (games Games) HandleMessage(msg models.BaseMessage) (error, interface{}) {
    var err error
    var out interface{}
    switch msg.Action {
        case "index":
        case "post":
            name := msg.Data["name"]
            if name == nil || name.(string) == "" {
                err = errors.New("No name provided for new game in queue message!")
            } else {
                nameSt := name.(string)
                Game{Name: nameSt, Code: strings.ToLower(nameSt)}.Add()
            }
        case "delete":
            id := msg.Data["id"]
            if id == nil || uint(id.(float64)) == 0 {
                err = errors.New("No id provided for game deletion in queue message!")
            } else {
                game := Game{}
                game.ID = uint(id.(float64))
                game.Delete()
            }
        default:
            err = fmt.Errorf("Unknown action for games: %s\n", msg.Action)
    }
    if (err == nil) {
        out = map[string]Games{
            "games": games.All(),
        }
    }
    return err, out
}