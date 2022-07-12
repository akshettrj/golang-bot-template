package main

import (
	"fmt"
	"net/http"
	"time"

	"golang-bot-template/database"
	"golang-bot-template/state"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"go.uber.org/zap"
)

func main() {
	config := state.State.Config
	config.LoadConfig()

	var logger *zap.Logger
	if config.Runtime.ProductionMode {
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stderr", "./golang-bot-template.log"}
		cfg.ErrorOutputPaths = []string{"stderr", "./golang-bot-template-error.log"}
		logger, _ = cfg.Build()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()
	state.State.Logger = logger

	if config.Runtime.LoadDatabase {
		db, err := database.Connect()
		if err != nil {
			logger.Panic("failed to connect to database", zap.Error(err))
		}
		state.State.Database = db
	}

	apiURL := config.Telegram.ApiURL
	if apiURL == "" {
		apiURL = gotgbot.DefaultAPIURL
	}
	bot, err := gotgbot.NewBot(config.Telegram.BotToken, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  apiURL,
		},
	})
	if err != nil {
		logger.Panic("Could not initialize bot", zap.Error(err))
	}
	state.State.Bot = bot
	logger.Info("Bot successfully logged in",
		zap.String("name", bot.FirstName),
		zap.String("username", "@"+bot.Username),
		zap.Int64("user_id", bot.Id),
	)

	// create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		ErrorLog: nil,
		DispatcherOpts: ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				fields := []zap.Field{}
				if chat := ctx.EffectiveChat; chat != nil {
					chat_title := chat.Title
					if chat_title == "" {
						chat_title = chat.FirstName
						if chat.LastName != "" {
							chat_title += " " + chat.LastName
						}
					}
					fields = append(fields, zap.Int64("chat_id", chat.Id))
					fields = append(fields, zap.String("chat_title", chat_title))
				}
				if msg := ctx.EffectiveMessage; msg != nil {
					fields = append(fields, zap.String("message", msg.Text))
				}
				logger.Error(err.Error(), fields...)
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		},
	})
	dispatcher := updater.Dispatcher

	dispatcher.AddHandler(handlers.NewCommand("start", start))

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: config.Runtime.DropUpdates,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		logger.Panic("failed to start polling", zap.Error(err))
	}
	logger.Info("Updater has started")

	updater.Idle()
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Hoi! I am up", nil)
	if err != nil {
		return fmt.Errorf("failed to send start message, %w", err)
	}
	return nil
}
