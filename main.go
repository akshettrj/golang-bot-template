package main

import (
    "log"
    "time"

    "golang-bot-template/state"
    "golang-bot-template/database"

    tele "gopkg.in/telebot.v3"
)

func main() {
    config := state.State.Config
    config.LoadConfig()

    db, err := database.Connect()
    if err != nil {
        log.Fatalln("could not connect to the database : ", err)
    }
    state.State.Database = db

    bot, err := tele.NewBot(tele.Settings{
        URL: config.Telegram.ApiURL,
        Token: config.Telegram.BotToken,
        Poller: &tele.LongPoller{Timeout: 10*time.Second},
    })
    if err != nil {
        log.Fatalln("could not initialize bot : ", err)
    }
    state.State.Bot = bot
    state.State.StartTime = time.Now()

    bot.Start()
}
