package telegram_bot

import (
	"context"
	"log/slog"
	"pandaexpress/tgbotapi"
)

func sendHelpMessage(bot *tgbotapi.BotAPI, chatID int64) {
	sessions.UpdateLastCommand(context.Background(),chatID, "/help")
	helpText := `Welcome to the Panda Express Bot!
	
	/help - Display this help message
	/wallet - Upload an image to add it to your wallet
	/search_product <query> - Search for products on Taobao.com`

	msg := tgbotapi.NewMessage(chatID, helpText)

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message:", "error", err)
	}
}
func sendUnknownCommandMessage(bot *tgbotapi.BotAPI, chatID int64) {
	unknownCommandText := `Sorry, I don't understand that command. Here are the available commands:
	
	/start - Start interacting with the bot
	/help - Display this help message
	/wallet - Upload an image to add it to your wallet
	/search_product <query> - Search for products on Taobao.com`

	msg := tgbotapi.NewMessage(chatID, unknownCommandText)

	_, err := bot.Send(msg)
	if err != nil {
		slog.Error("Error sending message:", "error", err)
	}
}
