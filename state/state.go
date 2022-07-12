package state

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type state struct {
	Bot       *gotgbot.Bot
	Config    *Config
	Database  *gorm.DB
	Logger    *zap.Logger
	StartTime time.Time
}

var State state

func init() {
	State.Config = &Config{}
}
