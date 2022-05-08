package state

import (
    "time"

    "golang-bot-template/config"

    "gorm.io/gorm"
    tele "gopkg.in/telebot.v3"
)

type state struct {
    Bot      *tele.Bot
    Config   *config.Config
    Database *gorm.DB

    StartTime time.Time
}

var State state

func init() {
    State.Config = &config.Config{}
}
