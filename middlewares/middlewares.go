package middlewares

import (
	"golang-bot-template/state"

	"golang.org/x/exp/slices"
	tele "gopkg.in/telebot.v3"
)

func OwnerMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	ownerID := state.State.Config.Telegram.OwnerID
	return func(c tele.Context) error {
		sender := c.Sender()
		if sender == nil || ownerID != sender.ID {
			return nil
		}
		return next(c)
	}
}

func SudoChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	sudoChats := state.State.Config.Telegram.SudoUsers
	return func(c tele.Context) error {
		sender := c.Sender()
		if sender == nil || !slices.Contains(sudoChats, sender.ID) {
			return nil
		}
		return next(c)
	}
}

func AuthorizedChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	authChats := state.State.Config.Telegram.AuthorizedChats
	return func(c tele.Context) error {
		chat := c.Chat()
		if chat == nil || !slices.Contains(authChats, chat.ID) {
			return nil
		}
		return next(c)
	}
}
