package main

import (
    "log"
    "time"

    "golang-bot-template/state"

    tele "gopkg.in/telebot.v3"
)

func main() {
    config := state.State.Config
    config.LoadConfig()

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
